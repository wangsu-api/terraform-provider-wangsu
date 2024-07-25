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

resource "wangsu_waap_whitelist" "demo" {
  rule_name   = "tf_test_update333"
  domain      = "waap.czp1"
  description = "terraform test update 11"

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

data "wangsu_waap_whitelist" "demo" {
  rule_name   = wangsu_waap_whitelist.demo.rule_name
  domain_list = [wangsu_waap_whitelist.demo.domain]
}

# output "whitelist_list" {
#   value = data.wangsu_waap_whitelist.demo
# }