variable "namespace" {
  type    = string
  default = "default"
}

variable "cluster_name" {
  type    = string
  default = "my-cluster"
}

variable "destinations_prometheus_url" {
  type    = string
  default = "https://prometheus-prod-24-prod-eu-west-2.grafana.net./api/prom/push"
}
variable "destinations_prometheus_username" {
  type    = string
  description = "Grafana Cloud Prometheus username"
}

variable "destinations_prometheus_password" {
  type    = string
  description = "Grafana Cloud Prometheus password"
  sensitive  = true
}