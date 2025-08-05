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

resource "wangsu_waap_pre_deploy_whitelist" "demo" {
  domain        = "waap.example.com"
  config_switch = "ON"

  rule_list {
    rule_name   = "example_rule"
    description = "Example description"
    conditions {
      path_conditions {
        match_type = "NOT_EQUAL"
        paths      = ["/p11", "/p21"]
      }
      uri_conditions {
        match_type = "NOT_EQUAL"
        uri        = ["/uri11", "/uri21"]
      }
      ua_conditions {
        match_type = "NOT_EQUAL"
        ua         = ["ua11", "ua21"]
      }
      referer_conditions {
        match_type = "NOT_EQUAL"
        referer    = ["re11", "re21"]
      }
      header_conditions {
        match_type = "NOT_EQUAL"
        key        = "h1"
        value_list = ["h111", "h211"]
      }
    }
  }
}

output "wangsu_waap_pre_deploy_result" {
  value = wangsu_waap_pre_deploy_whitelist.demo.host_list
}