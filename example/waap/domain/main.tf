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

resource "wangsu_waap_domain" "example" {
  target_domains = ["waap.czp", "waap.czp2"]

  waf_defend_config {
    rule_update_mode = "AUTO"
    config_switch    = "OFF"
    defend_mode      = "BLOCK"
  }

  customize_rule_config {
    config_switch = "OFF"
  }

  api_defend_config {
    config_switch = "OFF"
  }

  whitelist_config {
    config_switch = "OFF"
  }

  block_config {
    config_switch = "OFF"
  }

  dms_defend_config {
    config_switch   = "OFF"
    protection_mode = "MODERATE"
    ai_switch       = "ON"
  }

  intelligence_config {
    config_switch = "OFF"
    info_cate_act {
      attack_source = "BLOCK"
      spec_attack   = "LOG"
      industry      = "LOG"
    }
  }

  bot_manage_config {
    public_bots_act = "NO_USE"
    config_switch   = "OFF"
    ua_bots_act     = "LOG"
    web_risk_config {
      act = "LOG"
    }
    scene_analyse_switch = "ON"
  }

  rate_limit_config {
    config_switch = "OFF"
  }
}

data "wangsu_waap_domains" "demo" {
  domain_list = wangsu_waap_domain.example.target_domains
  #   domain_list = []
}

# output "domain_list" {
#   value = data.wangsu_waap_domain.demo
# }