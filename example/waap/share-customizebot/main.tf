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

resource "wangsu_waap_share_customizebot" "demo" {
  bot_name         = "test"
  bot_act          = "LOG"
  bot_description  = "desc"
  rela_domain_list = ["waap.example.com"]
  condition_list {
    condition_name       = "ASN"
    condition_value_list = ["1233"]
    condition_func       = "EQUAL"
    condition_key        = ""
  }
}


data "wangsu_waap_share_customizebots" "demo" {
  bot_name = wangsu_waap_share_customizebot.demo.bot_name
}

output "customizebot_list" {
  value = data.wangsu_waap_share_customizebots.demo
}