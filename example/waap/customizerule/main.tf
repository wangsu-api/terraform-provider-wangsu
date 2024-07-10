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
  }
}


data "wangsu_waap_customizerule" "demo" {
  rule_name   = wangsu_waap_customizerule.demo.rule_name
  domain_list = [wangsu_waap_customizerule.demo.domain]
}

output "customizerule_list" {
  value = data.wangsu_waap_customizerule.demo
}