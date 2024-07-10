terraform {
  required_providers {
    wangsu = {
      source  = "wangsu-api/wangsu"
      version = "1.0.2"
    }
  }
}

provider "wangsu" {
  secret_id  = ""
  secret_key = ""
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
    protection_mode = "UNDER_ATTACK"
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

data "wangsu_waap_domain" "demo" {
  domain_list = wangsu_waap_domain.example.target_domains
  #   domain_list = []
}

# output "domain_list" {
#   value = data.wangsu_waap_domain.demo
# }