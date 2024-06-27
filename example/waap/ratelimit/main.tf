terraform {
  required_providers {
    wangsu = {
      source  = "registry.terraform.io/wangsustack/wangsu"
      version = "1.0.0"
    }
  }
}

provider "wangsu" {
  secret_id  = ""
  secret_key = ""
}

resource "wangsu_waap_ratelimit" "demo" {
  domain             = "waap.test30.com"
  rule_name          = "test_muti_condition11"
  description = "your_description11"
  #   scene = "WEB" // or "API"
  scene              = "API"
  asset_api_id       = "1800805524845170689"
  statistical_item = "IP_UA" // or other options
  statistical_period = 601
  trigger_threshold  = 1001
  intercept_time = 601
  #   effective_status = "PERMANENT" // or other options
  effective_status   = "WITHIN" // or other options
  rate_limit_effective {
    effective = ["MON", "FRI"]
    start     = "07:00"
    end       = "18:00"
    timezone  = "17"
  }
  #   action             = "BLOCK" // or other options
  action = "LOG" // or other options

  rate_limit_rule_condition {
    ip_or_ips_conditions {
      match_type = "EQUAL"
      ip_or_ips  = ["192.168.1.11", "192.168.1.21"]
    }
    ip_or_ips_conditions {
      match_type = "NOT_EQUAL"
      ip_or_ips  = ["1.1.1.1"]
    }
    # WEB 维度才可配置
    /*path_conditions {
      match_type = "EQUAL"
      paths      = ["/p111", "/p211"]
    }*/
    # WEB 维度才可配置
    /*uri_conditions {
      match_type = "EQUAL"
      uri        = ["/uri11", "/uri21"]
    }*/
    # API 维度才可配置
    uri_param_conditions {
      match_type  = "EQUAL"
      param_name  = "param11"
      param_value = ["value11", "value21"]
    }
    ua_conditions {
      match_type = "EQUAL"
      ua         = ["ua11", "ua21"]
    }
    referer_conditions {
      match_type = "EQUAL"
      referer    = ["referer11", "referer21"]
    }
    header_conditions {
      match_type = "EQUAL"
      key        = "header_key"
      value_list = ["value11", "value21"]
    }
    area_conditions {
      match_type = "EQUAL"
      areas      = ["AI", "AU"]
    }
    status_code_conditions {
      match_type  = "EQUAL"
      status_code = ["200", "500"]
    }
    # WEB 维度才可配置
    /*method_conditions {
      match_type     = "EQUAL"
      request_method = ["GET", "DELETE"]
    }*/

    scheme_conditions {
      match_type  = "EQUAL"
      scheme = ["HTTPS"]
    }
  }
}

data "wangsu_waap_ratelimit" "demo" {
  rule_name   = wangsu_waap_ratelimit.demo.rule_name
  domain_list = [wangsu_waap_ratelimit.demo.domain]
}

output "ratelimit_list" {
  value = data.wangsu_waap_ratelimit.demo
}
