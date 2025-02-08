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

resource "wangsu_monitor_realtime_rule" "testrule" {
  rule_name          = "testrule"
  monitor_product    = "web"
  resource_type      = "domain"
  monitor_resources  = ["ALL"]
  statistical_method = "CONSOLIDATED"
  alert_frequency    = 5
  restore_notice     = "false"

  rule_items {
    start_time     = "00:00"
    end_time       = "00:59"
    condition_type = "ANY"
    period         = 5
    period_times   = 1

    condition_rules {
      monitor_item = "BANDWIDTH"
      condition    = "MAX"
      threshold    = "100000"
    }
    condition_rules {
      monitor_item = "FLOW"
      condition    = "MIN"
      threshold    = "1"
    }
  }
  rule_items {
    start_time     = "01:00"
    end_time       = "05:59"
    condition_type = "ALL"
    period         = 10
    period_times   = 2

    condition_rules {
      monitor_item = "REQUEST"
      condition    = "MAX"
      threshold    = "1000000"
    }
  }

  notices {
    notice_method = "WEBHOOK"
    notice_object = "https://example.com"
  }
  notices {
    notice_method = "MOBILE"
    notice_object = "1;2"
  }
}