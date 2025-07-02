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

resource "wangsu_waap_share_customizerule" "demo" {
  rule_name            = "tf_test"
  relation_domain_list = ["waap.demo.com"]
  description          = "terraform test"
  act                  = "LOG"

  condition {
    path_conditions {
      match_type = "EQUAL"
      paths      = ["/p1", "/p2"]
    }
    area_conditions {
      match_type = "NOT_EQUAL"
      areas      = ["AI", "AU"]
    }
    method_conditions {
      match_type     = "NOT_EQUAL"
      request_method = ["DELETE", "POST"]
    }
    header_conditions {
      match_type = "NOT_EQUAL"
      key        = "hk"
      value_list = ["h1", "h2"]
    }
    ua_conditions {
      match_type = "NOT_EQUAL"
      ua         = ["ua1", "ua2"]
    }
    referer_conditions {
      match_type = "NOT_EQUAL"
      referer    = ["re1", "re2"]
    }
    ja3_conditions {
      match_type = "NOT_EQUAL"
      ja3_list   = ["ja312345678901234567890123456788", "ja322345678901234567890123456788"]
    }
    ja4_conditions {
      match_type = "NOT_EQUAL"
      ja4_list   = ["ja41740600_c43983326036_1b2d6ce873a3", "ja42740600_c43983326036_1b2d6ce873a3"]
    }
  }
}


data "wangsu_waap_share_customizerules" "demo" {
  rule_name   = wangsu_waap_share_customizerule.demo.rule_name
}

output "customizerule_list" {
  value = data.wangsu_waap_share_customizerules.demo
}