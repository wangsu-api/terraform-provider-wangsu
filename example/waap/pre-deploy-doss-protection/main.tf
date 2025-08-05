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

resource "wangsu_waap_pre_deploy_ddos_protection" "demo" {
  domain = "waap.example.com"

  ddos_protect_switch {
    l7_ddos_switch = "ON"
    protect_mode   = "MODERATE"
    inner_switch   = "ON"
  }

  built_in_rules {
    rule_id        = "1721428087809277953"
    security_level = "DEFAULT_ENABLE"
    action         = "LOG"
  }

  built_in_rules {
    rule_id        = "1722064417224425473"
    security_level = "ATTACK_ENABLE"
    action         = "RR"
  }
}

output "wangsu_waap_pre_deploy_result" {
  value = wangsu_waap_pre_deploy_ddos_protection.demo.host_list
}