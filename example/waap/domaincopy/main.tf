terraform {
  required_providers {
    wangsu = {
      source  = "wangsu-api/wangsu"
      version = "1.0.3"
    }
  }
}

provider "wangsu" {
  secret_id  = ""
  secret_key = ""
}

resource "wangsu_waap_domain_copy" "demo" {
  source_domain  = "waap.test30.com"
  target_domains = ["waap.czp", "waap.czp2"]
}