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

resource "wangsu_waap_pre_deploy_waf" "demo" {
  domain        = "waap.example.com"
  config_switch = "ON"

  conf_basic {
    defend_mode      = "BLOCK"
    rule_update_mode = "AUTO"
  }

  rule_list {
    rule_id = 5002
    mode    = "BLOCK"
    exception_list {
      type         = "ip"
      match_type   = "EQUAL"
      content_list = ["192.168.1.1"]
    }
    exception_list {
      type         = "path"
      match_type   = "REGEX"
      content_list = ["/api/v1/.*"]
    }
  }

  rule_list {
    rule_id = 5003
    mode    = "LOG"
    exception_list {
      type         = "userAgent"
      match_type   = "CONTAIN"
      content_list = ["Mozilla"]
    }
  }
}

output "wangsu_waap_pre_deploy_result" {
  value = wangsu_waap_pre_deploy_waf.demo.host_list
}