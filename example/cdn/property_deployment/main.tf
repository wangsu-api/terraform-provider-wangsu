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

resource "wangsu_cdn_property_deployment" "example_deployment" {
  deployment_name = "example_deployment"
  target          = "staging"
  actions {
    action      = "deploy_property"
    property_id = 123456
    version     = 1
  }
}