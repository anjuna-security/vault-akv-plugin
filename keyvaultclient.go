package vault_akv_plugin

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"os"
	"os/exec"
)

const (
	VaultName          = "anjuna-key-vault"
	AzCmd              = "az"
	KeyvaultSubcommand = "keyvault"
)

type keyvaultClient struct {
	vaultName string
	logger    *hclog.Logger
}

func isEnvironmentSet() bool {
	return os.Getenv("KVAULT") != ""
}

func InitKeyvaultClient(logger *hclog.Logger) (*keyvaultClient, error) {
	var kvClient keyvaultClient

	vaultName := os.Getenv("KVAULT")
	if vaultName == "" {
		return nil, errors.New("KVAULT environment variable is not defined")
	}

	kvClient.vaultName = vaultName
	kvClient.logger = logger

	return &kvClient, nil
}

func (kvClient *keyvaultClient) ListSecrets() ([]string, error) {
	logger := *kvClient.logger
	parsedJson, err := runCmdAndParseJsonArrOutput(logger, "secret", "list", "--vault-name", VaultName)
	if err != nil {
		return nil, err
	}

	secrets := make([]string, 0)
	for _, entry := range parsedJson {
		secrets = append(secrets, entry["name"].(string))
	}

	return secrets, nil
}

func (kvClient *keyvaultClient) GetSecret(name string) (string, error) {
	logger := *kvClient.logger
	parsedJson, err := runCmdAndParseJsonOutput(logger, "secret", "show",
		"--name", name, "--vault-name", VaultName)
	if err != nil {
		return "", err
	}

	return parsedJson["value"].(string), nil
}

func runCmdAndParseJsonArrOutput(logger hclog.Logger, args ...string) ([]map[string]interface{}, error) {
	output, err := runAzKeyvaultCommand(args)
	if err != nil {
		logger.Error("Failed running command")
		return nil, err
	}

	var parsedJson []map[string]interface{}
	logger.Info(fmt.Sprintf("%s", output))
	json.Unmarshal(output, &parsedJson)
	return parsedJson, nil
}

func runCmdAndParseJsonOutput(logger hclog.Logger, args ...string) (map[string]interface{}, error) {
	output, err := runAzKeyvaultCommand(args)
	if err != nil {
		logger.Error("Failed running command")
		return nil, err
	}

	var parsedJson map[string]interface{}
	// logger.Info(fmt.Sprintf("%s", output))
	json.Unmarshal(output, &parsedJson)
	return parsedJson, nil
}

func runAzKeyvaultCommand(args []string) ([]byte, error) {
	azArgs := append([]string{KeyvaultSubcommand}, args...)
	cmdAz := exec.Command(AzCmd, azArgs...)

	output, err := cmdAz.Output()
	return output, err
}
