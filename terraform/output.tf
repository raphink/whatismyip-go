# Export the kubectl configuration file
resource "local_file" "kubernetes_config" {
  content  = "${digitalocean_kubernetes_cluster.dev-k8s.kube_config.0.raw_config}"
  filename = "kubeconfig.yaml"
}

# Print the load balancer ip
output "digitalocean_loadbalancer" {
  value       = "${digitalocean_loadbalancer.public.ip}"
  description = "The public IP address of the load balancer."
}
