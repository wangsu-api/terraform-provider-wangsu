terraform {
  required_providers {
    wangsu = {
      source = "registry.terraform.io/wangsu-api/wangsu"
    }
  }
}

provider "wangsu" {
  secret_id  = "my-secret-id"
  secret_key = "my-secret-key"
}

resource "wangsu_ssl_certificate" "cert-example" {
  name = "cert-example-name"
  cert = var.cert2
  key  = var.key2
}