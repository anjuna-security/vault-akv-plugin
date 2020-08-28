# Azure Key-Vault Secrets Plugin for Vault

This plugin enables storing secrets in [Azure Key Vault](https://azure.microsoft.com/en-us/services/key-vault/) 
while using [HashiCorp Vault](https://www.vaultproject.io/).

## Usage

All commands can be run using the provided [Makefile](./Makefile). However, it may be instructive to look at the commands to gain a greater understanding of how Vault registers plugins. Using the Makefile will result in running the Vault server in `dev` mode. Do not run Vault in `dev` mode in production. The `dev` server allows you to configure the plugin directory as a flag, and automatically registers plugin binaries in that directory. In production, plugin binaries must be manually registered.

This will build the plugin binary and start the Vault dev server:

```
# Build the Azure Key Vault plugin and start Vault dev server with plugin automatically registered
$ make
```

Now open a new terminal window and run the following commands:

```
# Open a new terminal window and export Vault dev server http address
$ export VAULT_ADDR='http://127.0.0.1:8200'

# Enable the AKV plugin
$ make enable

# Write a secret to the Mock secrets engine
$ vault write azure-key-vault/<keyvault-name> hello="world"
Success! Data written to: mock/test

# Retrieve secret from Mock secrets engine
$ vault read azure-key-vault/<keyvault-name>/hello
Key      Value
---      -----
hello    world
```

## Setting up an Azure Key Vault using Terraform

We provide a Terraform script under the ``terraform`` directory for automatically setting up an Azure Key Vault on Azure.

Running 

```
$ terraform init
$ terraform apply
```

in the ``terraform`` directory would create a key vault under your Azure account, assuming you have logged in using the Azure CLI tools.

## License

The Azure Key Vault plugin for Vault is released under a Mozilla Public License v2.0 (MPL 2.0). For details, check out the LICENSE file.
