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

resource "wangsu_waap_domain_copy" "demo" {
  source_domain  = "waap.test30.com"
  target_domains = ["waap.czp", "waap.czp2"]
}