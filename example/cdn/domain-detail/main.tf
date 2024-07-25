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

data "wangsu_cdn_domain_detail" "test-domain" {
  domain_name = "20240712001.conftest.com"
}

output "show-test-domain" {
  value = data.wangsu_cdn_domain_detail.test-domain
}