package vault_akv_plugin

import (
	"encoding/json"
	"errors"
	"github.com/hashicorp/go-hclog"
	"os"
	"os/exec"
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
	cmdAz := exec.Command("az", "keyvault", "secret", "list", "--vault-name", "anjuna-key-vault")

	output, err := cmdAz.Output()
	if err != nil {
		logger.Error("Failed running az")
		return nil, err
	}

	var parsedJson []map[string]interface{}
	json.Unmarshal(output, &parsedJson)

	secrets := make([]string, 0)
	for _, entry := range parsedJson {
		secrets = append(secrets, entry["name"].(string))
	}

	return secrets, nil
}
