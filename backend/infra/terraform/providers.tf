# Configuration Terraform requise (versions des providers)
terraform {
  required_providers {
    azurerm = {                        # Provider Azure RM
      source  = "hashicorp/azurerm"
      version = "~>3.0"                # Version du provider AzureRM
    }
    random = {
      source  = "hashicorp/random"     # Provider pour générer des valeurs aléatoires
      version = ">= 3.4.0"
    }
  }
}

# Configuration du provider AzureRM (connexion à Azure)
provider "azurerm" {
  features {}                          # Bloc requis (fonctionnalités par défaut du provider)
  # (Optionnel) configuration d'authentification si nécessaire
  # Par défaut, Terraform utilisera les credentials Azure CLI actifs (az login).
}
