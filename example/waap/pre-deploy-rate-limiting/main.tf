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

resource "wangsu_waap_pre_deploy_rate_limiting" "demo" {
  domain        = "waap.example.com"
  config_switch = "ON"

  rule_list {
    rule_name          = "example_rule"
    description        = "Example description"
    scene              = "WEB"
    statistical_item   = "IP"
    statistical_period = 60
    trigger_threshold  = 100
    intercept_time     = 300
    effective_status   = "PERMANENT"
    action             = "BLOCK"

    rate_limit_rule_condition {

      method_conditions {
        match_type     = "EQUAL"
        request_method = ["GET", "POST"]
      }

      area_conditions {
        match_type = "NOT_EQUAL"
        areas      = ["US", "CN"]
      }

      ip_or_ips_conditions {
        match_type = "EQUAL"
        ip_or_ips  = ["192.168.1.1", "10.0.0.0/24"]
      }

      uri_conditions {
        match_type = "CONTAIN"
        uri        = ["/api/v1/resource"]
      }

      path_conditions {
        match_type = "START_WITH"
        paths      = ["/path/to/resource"]
      }

      scheme_conditions {
        match_type = "EQUAL"
        scheme     = ["HTTPS"]
      }

      status_code_conditions {
        match_type  = "NOT_EQUAL"
        status_code = ["404", "500"]
      }
    }
  }

  rule_list {
    rule_name          = "example_rule2"
    description        = "Example description2"
    scene              = "WEB"
    statistical_item   = "IP_COOKIE"
    statistics_key     = "SESSION"
    statistical_period = 60
    trigger_threshold  = 100
    intercept_time     = 300
    effective_status   = "WITHIN"
    rate_limit_effective {
      effective = ["MON", "FRI"]
      start     = "07:00"
      end       = "18:00"
      timezone  = "17"
    }
    action = "BLOCK"

    rate_limit_rule_condition {

      ua_conditions {
        match_type = "NOT_CONTAIN"
        ua         = ["curl", "wget"]
      }

      header_conditions {
        match_type = "EQUAL"
        key        = "Content-Type"
        value_list = ["application/json"]
      }

      referer_conditions {
        match_type = "START_WITH"
        referer    = ["https://example.com"]
      }

      ja3_conditions {
        match_type = "EQUAL"
        ja3_list   = ["ja332345678901234567890123456788"]
      }

      ja4_conditions {
        match_type = "EQUAL"
        ja4_list   = ["ja43740600_c43983326036_1b2d6ce873a3"]
      }
    }
  }
}

output "wangsu_waap_pre_deploy_result" {
  value = wangsu_waap_pre_deploy_rate_limiting.demo.host_list
}