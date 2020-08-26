package main

import (
	"context"
	"log"
	"vault_akv_plugin"
)

func main() {
	akvClient, err := vault_akv_plugin.InitAkvClient()
	if err != nil {
		log.Fatal("Failed initializing AKV client")
	}
	akvClient.GetSecrets(context.Background(), vault_akv_plugin.KeyVaultURL, nil)
}
