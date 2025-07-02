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

data "wangsu_cdn_edge_hostnames" "example_edge_hostnames" {
  edge_hostnames     = ["www.example.com.wscncdn.com", "www.example2.com.wscncdn.com"]
  hostnames          = ["www.example.com", "www.example2.com"]
  comment            = "test"
  is_like            = true
  dns_service_status = "active"
  deploy_status      = ["pending", "success"]
  allow_china_cdn    = "0"
  offset             = 0
  limit              = 10
  sort_order         = "desc"
  sort_by            = "lastUpdateTime"
}

output "edge_hostnames" {
  value = data.wangsu_cdn_edge_hostnames.example_edge_hostnames
}