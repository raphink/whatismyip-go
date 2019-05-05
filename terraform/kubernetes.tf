# Create the cluster
resource "digitalocean_kubernetes_cluster" "dev-k8s" {
  name    = "dev-k8s"
  region  = "nyc1"
  version = "1.14.1-do.2"

  node_pool {
    name       = "dev-k8s-nodes"
    size       = "s-1vcpu-2gb"
    node_count = 1
    tags       = ["dev-k8s-nodes"]
  }
}
