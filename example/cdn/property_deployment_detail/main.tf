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

data "wangsu_cdn_property_deployment_detail" "example_deployment_detail" {
    deployment_id = 123456
}

output "deployment_detail" {
  value = data.wangsu_cdn_property_deployment_detail.example_deployment_detail
}