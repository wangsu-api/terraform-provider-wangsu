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

data "wangsu_cdn_edge_hostname_detail" "example_edge_hostname_detail" {
  edge_hostname = "www.example.com.wscncdn.com"
}

output "edge_hostname_detail" {
  value = data.wangsu_cdn_edge_hostname_detail.example_edge_hostname_detail
}