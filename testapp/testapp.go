// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

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
