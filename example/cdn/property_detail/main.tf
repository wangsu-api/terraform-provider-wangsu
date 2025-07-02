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

data "wangsu_cdn_property_detail" "example_property" {
  property_id = 123456
  version     = 1
}

output "property_detail" {
  value = data.wangsu_cdn_property_detail.example_property
}