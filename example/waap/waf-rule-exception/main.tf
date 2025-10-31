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

resource "wangsu_waap_waf_rule_exception" "demo" {
  domain = "waap.example.com"
  rule_id = 5002
  type = "ip"
  match_type = "EQUAL"
  content_list = ["1.1.1.2", "2.2.2.1"]
}