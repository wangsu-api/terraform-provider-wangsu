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

data "wangsu_ssl_certificate_detail" "myCert" {
  certificate_id = "1464893"
}

output "data" {
  value = data.wangsu_ssl_certificate_detail.myCert
}