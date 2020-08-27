package vault_akv_plugin

import (
	"context"
	"github.com/hashicorp/go-hclog"
	"strings"

	"fmt"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	httpsPrefix = "https://"
)

var (
	KeyVaultURL string
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

	b.Backend.Setup(ctx, conf)
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

	if !isEnvironmentSet() {
		logger.Error("Environment vars not set. " +
			"You need to set KVAULT")
		return nil
	}

	akvClient, err := InitKeyvaultClient(&logger)
	if err != nil {
		logger.Error("Failed initializing AVK client")
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

func (b *backend) handleExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	out, err := req.Storage.Get(ctx, req.Path)
	if err != nil {
		return false, errwrap.Wrapf("existence check failed: {{err}}", err)
	}

	return out != nil, nil
}

func (b *backend) handleRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if req.ClientToken == "" {
		return nil, fmt.Errorf("client token empty")
	}

	//path := data.Get("path").(string)

	// Decode the data
	var rawData map[string]interface{}

	//if err := jsonutil.DecodeJSON(b.store[req.ClientToken+"/"+path], &rawData); err != nil {
	//	return nil, errwrap.Wrapf("json decoding failed: {{err}}", err)
	//}

	// Generate the response
	resp := &logical.Response{
		Data: rawData,
	}

	return resp, nil
}

func (b *backend) handleWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if req.ClientToken == "" {
		return nil, fmt.Errorf("client token empty")
	}

	// Check to make sure that kv pairs provided
	if len(req.Data) == 0 {
		return nil, fmt.Errorf("data must be provided to store in secret")
	}

	//path := data.Get("path").(string)

	// JSON encode the data
	//buf, err := json.Marshal(req.Data)
	//if err != nil {
	//	return nil, errwrap.Wrapf("json encoding failed: {{err}}", err)
	//}

	// Store kv pairs in map at specified path
	//b.store[req.ClientToken+"/"+path] = buf

	return nil, nil
}

func (b *backend) handleDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if req.ClientToken == "" {
		return nil, fmt.Errorf("client token empty")
	}

	//path := data.Get("path").(string)

	// Remove entry for specified path
	//delete(b.store, path)

	return nil, nil
}

func (b *backend) handleList(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if req.ClientToken == "" {
		return nil, fmt.Errorf("client token empty")
	}

	logger := b.Logger()
	secretsPath := data.Get("path").(string)
	logger.Debug(fmt.Sprintf("Listing secrets at path %s", secretsPath))

	akvClient := b.akvClient
	secrets, err := akvClient.ListSecrets()
	if err != nil {
		logger.Error(err.Error())
		return logical.ErrorResponse(err.Error()), err
	}

	logger.Debug("Retrieved secrets from key vault")
	return logical.ListResponse(secrets), nil
}

const pluginHelp = `
The Azure Key Vault backend is a secrets backend that stores kv pairs in an Azure Key Vault.
`
