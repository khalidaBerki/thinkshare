resource "helm_release" "grafana-k8s-monitoring" {
  name             = "grafana-k8s-monitoring"
  repository       = "https://grafana.github.io/helm-charts"
  chart            = "k8s-monitoring"
  namespace        = var.namespace
  create_namespace = true
  atomic           = true
  timeout          = 300

  values = [file("${path.module}/values.yaml")]

  set {
    name  = "cluster.name"
    value = var.cluster_name
  }
  set {
    name  = "destinations[0].url"
    value = var.destinations_prometheus_url
  }
  set_sensitive {
    name  = "destinations[0].auth.username"
    value = var.destinations_prometheus_username
  }
  set_sensitive {
    name  = "destinations[0].auth.password"
    value = var.destinations_prometheus_password
  }
  # ...idem pour Loki et OTLP, voir le bloc généré par Grafana Cloud...
}