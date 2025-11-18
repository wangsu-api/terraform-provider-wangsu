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

data "wangsu_ssl_certificate_application_detail" "example" {
  order_id = "SO20251029194956687.199497441"
}

output "data" {
  value = data.wangsu_ssl_certificate_application_detail.example
}