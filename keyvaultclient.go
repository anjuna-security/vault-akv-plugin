// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package vault_akv_plugin

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"os/exec"
)

const (
	AzCmd              = "az"
	KeyvaultSubcommand = "keyvault"
)

type keyvaultClient struct {
	logger *hclog.Logger
}

func InitKeyvaultClient(logger *hclog.Logger) (*keyvaultClient, error) {
	var kvClient keyvaultClient
	kvClient.logger = logger

	_, err := exec.LookPath(AzCmd)
	if err != nil {
		(*logger).Error("Can't find Azure CLI tools")
		return nil, err
	}

	return &kvClient, nil
}

func (kvClient *keyvaultClient) ListSecrets(vaultName string) ([]string, error) {
	parsedJson, err := runCmdAndParseJsonArrOutput(*kvClient.logger,
		"secret", "list", "--vault-name", vaultName)
	if err != nil {
		return nil, err
	}

	secrets := make([]string, 0)
	for _, entry := range parsedJson {
		secrets = append(secrets, entry["name"].(string))
	}

	return secrets, nil
}

func (kvClient *keyvaultClient) GetSecret(vaultName string, name string) (string, error) {
	parsedJson, err := runCmdAndParseJsonOutput(*kvClient.logger,
		"secret", "show", "--name", name, "--vault-name", vaultName)
	if err != nil {
		return "", err
	}

	return parsedJson["value"].(string), nil
}

func (kvClient *keyvaultClient) SetSecret(vaultName string, name string, value string) error {
	_, err := runCmdAndParseJsonOutput(*kvClient.logger,
		"secret", "set", "--name", name, "--value", value, "--vault-name", vaultName)
	return err
}

func (kvClient *keyvaultClient) DeleteSecret(vaultName string, name string) error {
	_, err := runCmdAndParseJsonOutput(*kvClient.logger,
		"secret", "delete", "--name", name, "--vault-name", vaultName)
	return err
}

func (kvClient *keyvaultClient) PurgeSecret(vaultName string, name string) error {
	args := []string{"secret", "purge", "--name", name, "--vault-name", vaultName}
	output, err := runAzKeyvaultCommand(args)
	if err != nil {
		(*kvClient.logger).Trace(fmt.Sprintf("%s: %s", output, err.Error()))
	}

	return nil
}

func runCmdAndParseJsonArrOutput(logger hclog.Logger, args ...string) ([]map[string]interface{}, error) {
	output, err := runAzKeyvaultCommand(args)
	if err != nil {
		logger.Error("Failed running command")
		return nil, err
	}
	logger.Trace(fmt.Sprintf("%s", output))

	var parsedJson []map[string]interface{}
	json.Unmarshal(output, &parsedJson)
	return parsedJson, nil
}

func runCmdAndParseJsonOutput(logger hclog.Logger, args ...string) (map[string]interface{}, error) {
	output, err := runAzKeyvaultCommand(args)
	if err != nil {
		logger.Error("Failed running command")
		return nil, err
	}
	logger.Trace(fmt.Sprintf("%s", output))

	var parsedJson map[string]interface{}
	json.Unmarshal(output, &parsedJson)
	return parsedJson, nil
}

func runAzKeyvaultCommand(args []string) ([]byte, error) {
	azArgs := append([]string{KeyvaultSubcommand}, args...)
	cmdAz := exec.Command(AzCmd, azArgs...)
	output, err := cmdAz.Output()
	return output, err
}
