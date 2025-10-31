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

resource "wangsu_waap_waf_config" "demo" {
  domain = "waap.example.com"
  conf_basic {
    defend_mode      = "BLOCK"
    rule_update_mode = "MANUAL"
  }

  rule_list {
    rule_id = 5002
    mode    = "BLOCK"
  }

  rule_list {
    rule_id = 5003
    mode    = "LOG"
  }


  scan_protection {
    scan_tools_config {
      action = "LOG"
    }

    repeated_violation_config {
      action              = "LOG"
      target              = "IP_JA3"
      period              = 10
      waf_rule_type_count = 11
      block_count         = 12
      duration            = 13
    }

    directory_probing_config {
      action                           = "BLOCK"
      target                           = "IP"
      period                           = 20
      request_count_threshold          = 21
      non_existent_directory_threshold = 22
      rate404_threshold                = 23
      duration                         = 24
    }
  }
}

data "wangsu_waap_waf_configs" "demo" {
  domain_list = ["waap.example.com"]
}

output "waf_configs" {
  value = data.wangsu_waap_waf_configs.demo
}