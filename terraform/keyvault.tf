provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy = true
    }
  }
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "azure_key_vault" {
  name     = "anjuna-key-vault"
  location = "West US"
}

resource "azurerm_key_vault" "azure_key_vault" {
  name                        = "anjuna-keyvault"
  location                    = azurerm_resource_group.azure_key_vault.location
  resource_group_name         = azurerm_resource_group.azure_key_vault.name
  enabled_for_disk_encryption = false
  tenant_id                   = data.azurerm_client_config.current.tenant_id
  soft_delete_enabled         = true
  purge_protection_enabled    = false

  sku_name = "standard"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "get",
	  "create",
	  "update",
	  "list",
	  "delete",
    ]

    secret_permissions = [
      "get",
	  "set",
	  "list",
	  "delete",
	  "purge",
    ]

    storage_permissions = [
      "get",
    ]
  }

  network_acls {
    default_action = "Allow"
    bypass         = "AzureServices"
  }

  tags = {
    environment = "AKV Plugin"
  }
}

output "azure_tenant_id" { value = "${azurerm_key_vault.azure_key_vault.tenant_id}" }
output "vault_uri" { value = "${azurerm_key_vault.azure_key_vault.vault_uri}" }
