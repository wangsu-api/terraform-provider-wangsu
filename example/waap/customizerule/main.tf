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

resource "wangsu_waap_customizerule" "demo" {
  rule_name   = "tf_test_u"
  domain      = "waap.test30.com"
  description = "terraform test update"
  #   scene = "WEB"
  scene  = "API"
  api_id = "1800805524845170689"
  act    = "BLOCK"

  condition {
    path_conditions {
      match_type = "EQUAL"
      paths      = ["/p11", "/p21"]
    }
    uri_param_conditions {
      match_type  = "NOT_EQUAL"
      param_name  = "p1"
      param_value = ["pv1", "pv2"]

    }
    area_conditions {
      match_type = "NOT_EQUAL"
      areas      = ["AI", "AU"]
    }
    method_conditions {
      match_type     = "NOT_EQUAL"
      request_method = ["GET", "POST"]
    }
    header_conditions {
      match_type = "NOT_EQUAL"
      key        = "hk1"
      value_list = ["h1", "h2"]
    }
    ja3_conditions {
      match_type = "EQUAL"
      ja3_list   = ["ja332345678901234567890123456788", "ja342345678901234567890123456788"]
    }
    ja4_conditions {
      match_type = "EQUAL"
      ja4_list   = ["ja43740600_c43983326036_1b2d6ce873a3", "ja44740600_c43983326036_1b2d6ce873a3"]
    }
  }
}


data "wangsu_waap_customizerules" "demo" {
  rule_name   = wangsu_waap_customizerule.demo.rule_name
  domain_list = [wangsu_waap_customizerule.demo.domain]
}

output "customizerule_list" {
  value = data.wangsu_waap_customizerules.demo
}