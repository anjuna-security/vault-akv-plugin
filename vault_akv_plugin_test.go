package vault_akv_plugin

import (
	"github.com/hashicorp/go-hclog"
	"testing"
)

var (
	akvClient *keyvaultClient
)

func TestInitAkvClient(t *testing.T) {
	logger := hclog.New(&hclog.LoggerOptions{})
	akvClientRet, err := InitKeyvaultClient(&logger)
	if err != nil {
		t.Errorf("Failed initializing Azure Key Vault client")
	}

	akvClient = akvClientRet
}

func TestListSecrets(t *testing.T) {
	secrets, err := akvClient.ListSecrets()
	if err != nil {
		t.Errorf("Failed listing secrets")
	}

	t.Logf("%v", secrets)
}

func TestGetSecret(t *testing.T) {
	value, err := akvClient.GetSecret("hello")
	if err != nil {
		t.Errorf("Failed getting secret")
	}

	t.Logf("hello=%s", value)
}
