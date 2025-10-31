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

resource "wangsu_waap_bot_scene_whitelist" "demo" {
  domain      = "waap.example.com"
  name        = "test"
  description = "desc"
  conditions {
    match_name       = "IP_IPS"
    match_type       = "EQUAL"
    match_key        = ""
    match_value_list = ["1.1.1.1"]
  }
  conditions {
    match_name       = "PATH"
    match_type       = "EQUAL"
    match_key        = ""
    match_value_list = ["/path/test"]
  }
}

data "wangsu_waap_bot_scene_whitelists" "demo" {
  domain_list = [wangsu_waap_bot_scene_whitelist.demo.domain]
}

output "bot_scene_whitelist" {
  value = data.wangsu_waap_bot_scene_whitelists.demo
}