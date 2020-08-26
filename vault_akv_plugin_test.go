package vault_akv_plugin

import (
	"github.com/Azure/azure-sdk-for-go/profiles/latest/keyvault/keyvault"
	"github.com/Azure/azure-sdk-for-go/services/keyvault/auth"
	"os"
	"os/exec"
	"testing"
)

var (
	akvClient *keyvault.BaseClient
)

func TestAuthorizer(t *testing.T) {
	_, err := auth.NewAuthorizerFromCLI()
	if err != nil {
		t.Errorf("Failed initializing authorizer")
	}
}

func TestInitAkvClient(t *testing.T) {
	akvClientRet, err := InitAkvClient()
	if err != nil {
		t.Errorf("Failed initializing Azure Key Vault client")
	}

	akvClient = akvClientRet
}

func TestListSecretsUsingAz(t *testing.T) {
	azExecPath, err := exec.LookPath("az")
	if err != nil {
		t.Errorf("Failed finding az")
	}

	cmdAz := &exec.Cmd {
		Path: azExecPath,
		Args: []string { azExecPath, "keyvault", "secret", "list", "--vault-name", "anjuna-key-vault" },
		Stdout: os.Stdout,
		Stderr: os.Stdout,
	}

	t.Logf("Running command %s", cmdAz.String())

	err = cmdAz.Run()
	if err != nil {
		t.Errorf("Failed listing secrets using az")
	}
}

/*
func TestGetSecret(t *testing.T) {
	_, err := akvClient.GetSecret(context.Background(), KeyVaultURL,
		"hello", "1")
	if err != nil {
		t.Errorf("Failed retrieving secret")
	}
}

func TestListSecrets(t *testing.T) {
	_, err := akvClient.GetSecrets(context.Background(), KeyVaultURL, nil)
	if err != nil {
		t.Errorf("Failed listing secrets")
	}
}
*/