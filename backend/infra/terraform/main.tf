resource "azurerm_resource_group" "rg" {
  name     = "MonGroupeRessources"
  location = var.location
}

resource "azurerm_key_vault" "vault" {
  name                        = "${var.prefix}-kv"
  location                    = var.location
  resource_group_name         = azurerm_resource_group.rg.name
  tenant_id                   = data.azurerm_client_config.current.tenant_id
  sku_name                    = var.key_vault_sku
  purge_protection_enabled    = false
  enabled_for_disk_encryption = true
  enabled_for_deployment      = true
  enabled_for_template_deployment = true

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    secret_permissions = [
      "Get",
      "List",
      "Set"
    ]
  }
}

data "azurerm_client_config" "current" {}

resource "azurerm_postgresql_flexible_server" "db" {
  name                   = "pec-2"
  resource_group_name    = "pec2"  # <-- doit pointer ici
  location               = var.location
  administrator_login    = var.db_admin
  administrator_password = var.db_password
  version                = "13"
  storage_mb             = 32768
  sku_name               = "B_Standard_B1ms"
  zone                   = "1"
  backup_retention_days  = 7
  geo_redundant_backup_enabled = false
  public_network_access_enabled = true
}

# Ajout du cluster AKS le moins cher
resource "azurerm_kubernetes_cluster" "aks" {
  name                = "${var.prefix}-aks"
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name
  dns_prefix          = "${var.prefix}-aks"

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_B2s" # Le plus petit VM size possible
    os_disk_size_gb = 30
  }

  identity {
    type = "SystemAssigned"
  }

  network_profile {
    network_plugin = "kubenet"
    load_balancer_sku = "standard"
  }
}

# Variables complémentaires
variable "jwt_secret" {
  description = "Secret JWT pour l'API"
  type        = string
}

variable "pg_host" {
  description = "FQDN du serveur PostgreSQL"
  type        = string
}

variable "port" {
  description = "Port d'écoute de l'application"
  type        = string
  default     = "8080"
}

variable "gin_mode" {
  description = "Mode Gin (debug ou release)"
  type        = string
  default     = "debug"
}

variable "pg_port" {
  description = "Port PostgreSQL"
  type        = string
  default     = "5432"
}

variable "pg_sslmode" {
  description = "Mode SSL PostgreSQL"
  type        = string
  default     = "require"
}

variable "key_vault_sku" {
  description = "SKU du Key Vault (standard/premium)"
  type        = string
  default     = "standard"
}

variable "stripe_secret_key" {
  description = "Clé secrète Stripe"
  type        = string
}

variable "stripe_webhook_secret" {
  description = "Clé secrète du webhook Stripe"
  type        = string
}

output "aks_node_public_ip" {
  description = "IP publique du nœud AKS pour accès à Swagger via NodePort"
  value = azurerm_kubernetes_cluster.aks.default_node_pool[0].node_public_ip_prefix_id
  # Note : Pour obtenir l'IP publique réelle, il faudra utiliser 'kubectl get nodes -o wide' après déploiement
}

resource "azurerm_key_vault_secret" "stripe_secret_key" {
  name         = "stripe-secret-key"
  value        = var.stripe_secret_key
  key_vault_id = azurerm_key_vault.vault.id
}

resource "azurerm_key_vault_secret" "stripe_webhook_secret" {
  name         = "stripe-webhook-secret"
  value        = var.stripe_webhook_secret
  key_vault_id = azurerm_key_vault.vault.id
}
