// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package vault_akv_plugin

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"strings"
)

// Factory configures and returns AKV backends
func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	if conf == nil {
		return nil, fmt.Errorf("configuration passed into backend is nil")
	}

	b := Backend(conf)
	if b == nil {
		return nil, fmt.Errorf("failed initializing backend")
	}

	err := b.Backend.Setup(ctx, conf)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// backend wraps the backend framework and adds an Azure Key Vault client
type backend struct {
	*framework.Backend
	akvClient *keyvaultClient
}

func Backend(_ *logical.BackendConfig) *backend {
	var b backend
	logger := hclog.New(&hclog.LoggerOptions{})

	akvClient, err := InitKeyvaultClient(&logger)
	if err != nil {
		logger.Error("Failed initializing AVK client", err.Error())
		return nil
	}

	b.akvClient = akvClient

	b.Backend = &framework.Backend{
		Help:        strings.TrimSpace(pluginHelp),
		BackendType: logical.TypeLogical,
	}

	b.Backend.Paths = append(b.Backend.Paths, b.paths()...)

	logger.Debug("Initialized backend for Azure Key Vault plugin")
	return &b
}

func (b *backend) paths() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: framework.MatchAllRegex("path"),

			Fields: map[string]*framework.FieldSchema{
				"path": {
					Type:        framework.TypeString,
					Description: "Specifies the path of the secret.",
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleRead,
					Summary:  "Retrieve the secret from the map.",
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleWrite,
					Summary:  "Store a secret at the specified location.",
				},
				logical.CreateOperation: &framework.PathOperation{
					Callback: b.handleWrite,
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: b.handleDelete,
					Summary:  "Deletes the secret at the specified location.",
				},
				logical.ListOperation: &framework.PathOperation{
					Callback: b.handleList,
					Summary:  "Lists the secrets at the specified location.",
				},
			},

			ExistenceCheck: b.handleExistenceCheck,
		},
	}
}

func (b *backend) handleExistenceCheck(ctx context.Context, req *logical.Request, _ *framework.FieldData) (bool, error) {
	out, err := req.Storage.Get(ctx, req.Path)
	if err != nil {
		return false, errwrap.Wrapf("existence check failed: {{err}}", err)
	}

	return out != nil, nil
}

func (b *backend) handleRead(_ context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if req.ClientToken == "" {
		return nil, fmt.Errorf("client token empty")
	}

	path := data.Get("path").(string)
	if path == "" {
		const errMsg = "no secret path specified"
		b.Logger().Error(errMsg)
		return logical.ErrorResponse(errMsg), errors.New(errMsg)
	}

	// path encodes both the key vault name and the secret name as
	// <vault name>/<secret name>
	// in order to read secret "hello" from a vault named "anjuna-key-vault",
	// you need to run
	// $ vault read vault-akv-plugin/anjuna-key-vault/hello

	pathComponents := strings.Split(path, "/")
	if len(pathComponents) != 2 {
		const errMsg = "invalid path specified"
		b.Logger().Error(errMsg)
		return logical.ErrorResponse(errMsg), errors.New(errMsg)
	}

	vaultName := pathComponents[0]
	secretName := pathComponents[1]

	b.Logger().Debug(fmt.Sprintf("Fetching secret %s from vault %s", secretName, vaultName))

	value, err := b.akvClient.GetSecret(vaultName, secretName)
	if err != nil {
		return logical.ErrorResponse(err.Error()), errors.New(err.Error())
	}
	if value == "" {
		return nil, nil
	}

	// Generate the response
	secretData := make(map[string]interface{}, 1)
	secretData[secretName] = value

	response := &logical.Response{
		Data: secretData,
	}

	return response, nil
}

func getFirstKeyValueFromMap(m map[string]interface{}) (key string, value string) {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}

	return keys[0], m[keys[0]].(string)
}

func (b *backend) handleWrite(_ context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if req.ClientToken == "" {
		return nil, fmt.Errorf("client token empty")
	}

	// The key vault name is encoded in the path. In order to write
	// a secret "hello" with value "world" to a vault named "anjuna-key-vault",
	// you need to run:
	// $ vault write vault-akv-plugin/anjuna-key-vault hello=world

	splittedPath := strings.Split(data.Get("path").(string), "/")

	if len(splittedPath) != 2 {
		const errMsg = "vault name and secret name must be specified"
		b.Logger().Error(errMsg)
		return logical.ErrorResponse(errMsg), errors.New(errMsg)
	}

	vaultName := splittedPath[0]
	if vaultName == "" {
		const errMsg = "vault name is not specified"
		b.Logger().Error(errMsg)
		return logical.ErrorResponse(errMsg), errors.New(errMsg)
	}

	_, value := getFirstKeyValueFromMap(req.Data)
	name := splittedPath[len(splittedPath)-1]
	b.Logger().Debug(fmt.Sprintf("Setting secret %s to %s in vault %s", name, value, vaultName))

	// JSON encode the data
	err := b.akvClient.SetSecret(vaultName, name, value)
	if err != nil {
		return logical.ErrorResponse(err.Error()), errors.New(err.Error())
	}

	return nil, nil
}

func (b *backend) handleDelete(_ context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if req.ClientToken == "" {
		return nil, fmt.Errorf("client token empty")
	}

	path := data.Get("path").(string)
	if path == "" {
		const errMsg = "no secret path specified"
		b.Logger().Error(errMsg)
		return logical.ErrorResponse(errMsg), errors.New(errMsg)
	}

	// path encodes both the key vault name and the secret name as
	// <vault name>/<secret name>
	// in order to delete secret "hello" from a vault named "anjuna-key-vault",
	// you need to run
	// $ vault delete vault-akv-plugin/anjuna-key-vault/hello

	pathComponents := strings.Split(path, "/")
	if len(pathComponents) != 2 {
		const errMsg = "invalid path specified"
		b.Logger().Error(errMsg)
		return logical.ErrorResponse(errMsg), errors.New(errMsg)
	}

	vaultName := pathComponents[0]
	secretName := pathComponents[1]

	b.Logger().Debug(fmt.Sprintf("Deleting secret %s from vault %s", secretName, vaultName))

	err := b.akvClient.DeleteSecret(vaultName, secretName)
	if err != nil {
		return logical.ErrorResponse(err.Error()), err
	}

	return nil, nil
}

func (b *backend) handleList(_ context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if req.ClientToken == "" {
		return nil, fmt.Errorf("client token empty")
	}

	vaultName := strings.TrimSuffix(data.Get("path").(string), "/")
	b.Logger().Debug(fmt.Sprintf("Listing secrets in vault %s", vaultName))

	secrets, err := b.akvClient.ListSecrets(vaultName)
	if err != nil {
		b.Logger().Error(err.Error())
		return logical.ErrorResponse(err.Error()), err
	}

	b.Logger().Debug("Retrieved secrets from key vault")
	return logical.ListResponse(secrets), nil
}

const pluginHelp = `
The Azure Key Vault backend is a secrets backend that stores kv pairs in an Azure Key Vault.
`
