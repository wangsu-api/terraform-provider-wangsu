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

data "wangsu_cdn_properties" "example_properties" {
  hostname     = ["www.example.com", "www.example2.com"]
  service_type = "wsa"
  target       = ["staging", "production"]
  limit        = 10
  offset       = 0
  sort_by      = "lastUpdateTime"
  sort_order   = "desc"
}

output "properties" {
  value = data.wangsu_cdn_properties.example_properties
}