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

data "wangsu_appa_domain_detail" "appa1-data" {
  domain_name = "20240710001.conftest.com"
}

output "show-appa1" {
  value = data.wangsu_appa_domain_detail.appa1-data
}