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

data "wangsu_ssl_certificates" "myCertList" {
  name = "test20240625"
}

output "certList" {
  value = data.wangsu_ssl_certificates.myCertList
}