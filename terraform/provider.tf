# Configure the digitalocean provider with it's token
variable "DO_TOKEN" {}

provider "digitalocean" {
  token = "${var.DO_TOKEN}"
}
