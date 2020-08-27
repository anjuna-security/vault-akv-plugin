package vault_akv_plugin

import (
	"github.com/hashicorp/go-hclog"
	"os/exec"
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

func TestListSecretsUsingAz(t *testing.T) {

	cmdAz := exec.Command("az", "keyvault", "secret", "list", "--vault-name", "anjuna-key-vault")

	t.Logf("Running command %s", cmdAz.String())

	output, err := cmdAz.Output()
	if err != nil {
		t.Errorf("Failed listing secrets using az")
	}

	t.Logf("%s", output)
}

func TestListSecrets(t *testing.T) {
	secrets, err := akvClient.ListSecrets()
	if err != nil {
		t.Errorf("Failed listing secrets")
	}

	t.Logf("%v", secrets)
}
