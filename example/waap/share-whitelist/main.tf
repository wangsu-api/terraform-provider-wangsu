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

resource "wangsu_waap_share_whitelist" "demo" {
  rule_name            = "tf_test_update"
  relation_domain_list = ["waap.demo.com"]
  description          = "terraform test update"

  conditions {
    path_conditions {
      match_type = "NOT_EQUAL"
      paths      = ["/p1", "/p2"]
    }
    uri_conditions {
      match_type = "NOT_EQUAL"
      uri        = ["/uri1", "/uri2"]
    }
    ua_conditions {
      match_type = "NOT_EQUAL"
      ua         = ["ua1", "ua2"]
    }
    referer_conditions {
      match_type = "NOT_EQUAL"
      referer    = ["re1", "re2"]
    }
    header_conditions {
      match_type = "NOT_EQUAL"
      key        = "h"
      value_list = ["h11", "h21"]
    }
  }
}

data "wangsu_waap_share_whitelists" "demo" {
  rule_name = wangsu_waap_share_whitelist.demo.rule_name
}