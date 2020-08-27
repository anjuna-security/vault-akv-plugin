package main

import (
	"github.com/hashicorp/go-hclog"
	"log"
	"vault_akv_plugin"
)

const VaultName = "anjuna-keyvaukt"

func main() {
	logger := hclog.New(&hclog.LoggerOptions{})
	akvClient, err := vault_akv_plugin.InitKeyvaultClient(&logger)
	if err != nil {
		log.Fatal("Failed initializing AKV client")
	}

	secrets, err := akvClient.ListSecrets(VaultName)
	if err != nil {
		log.Fatal("Failed listing secrets")
	}
	log.Printf("%v", secrets)
}
