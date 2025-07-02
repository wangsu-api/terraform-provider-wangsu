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

data "wangsu_cdn_property_deployments" "example_deployments" {
    property_id  = 123456
    status       = "FAIL"
    target       = "staging"
    offset       = 0
    limit        = 10
    sort_order   = "desc"
    sort_by      = "submissionTime"
}

output "deployments" {
  value = data.wangsu_cdn_property_deployments.example_deployments
}