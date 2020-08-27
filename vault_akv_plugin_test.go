package vault_akv_plugin

import (
	"github.com/hashicorp/go-hclog"
	"testing"
)

var (
	akvClient *keyvaultClient
)

const VaultName = "anjuna-key-vault"

func TestInitAkvClient(t *testing.T) {
	logger := hclog.New(&hclog.LoggerOptions{})
	akvClientRet, err := InitKeyvaultClient(&logger)
	if err != nil {
		t.Errorf("Failed initializing Azure Key Vault client")
	}

	akvClient = akvClientRet
}

func TestListSecrets(t *testing.T) {
	secrets, err := akvClient.ListSecrets(VaultName)
	if err != nil {
		t.Errorf("Failed listing secrets")
	}

	t.Logf("%v", secrets)
}

func TestSetSecret(t *testing.T) {
	err := akvClient.SetSecret(VaultName, "hello", "world")
	if err != nil {
		t.Errorf("Failed setting secret")
	}
}

func TestGetSecret(t *testing.T) {
	value, err := akvClient.GetSecret(VaultName, "hello")
	if err != nil {
		t.Errorf("Failed getting secret")
	}

	t.Logf("hello=%s", value)
}
