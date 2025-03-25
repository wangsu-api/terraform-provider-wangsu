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

data "wangsu_monitor_realtime_rules_detail" "test-rule" {
  rule_name = "test_rule_name"
}

output "show-test-rule" {
  value = data.wangsu_monitor_realtime_rules_detail.test-rule
}