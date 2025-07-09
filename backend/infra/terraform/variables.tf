# Préfixe pour nommer les ressources (doit être en minuscules et unique pour certaines ressources)
variable "prefix" {
  description = "Préfixe utilisé dans les noms des ressources Azure"
  type        = string
  default     = "pec2"
}

# Emplacement (région Azure) où déployer les ressources
variable "location" {
  description = "Région Azure pour le déploiement des ressources."
  type        = string
  default     = "francecentral"
}

variable "db_admin" {
  description = "Nom d'utilisateur admin pour PostgreSQL"
  type        = string
  default     = "Berki"
}

variable "db_password" {
  description = "Mot de passe admin pour PostgreSQL"
  type        = string
  sensitive   = true
}

variable "db_name" {
  description = "Nom de la base de données principale"
  type        = string
  default     = "postgres"
}