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

resource "wangsu_waap_bot_config" "demo" {
  domain = "www.example.com"

  general_strategy {
    ai_bots {
      bot_category = "ai_assistant"
      action       = "NO_USE"
    }
    ai_bots {
      bot_category = "ai_data_scraper"
      action       = "NO_USE"
    }
    ai_bots {
      bot_category = "ai_search_crawler"
      action       = "NO_USE"
    }
    ai_bots {
      bot_category = "undocumented_ai_agent"
      action       = "NO_USE"
    }

    public_bots {
      bot_category = "feed_fetcher"
      action       = "ACCEPT"
    }
    public_bots {
      bot_category = "marketing_analysis"
      action       = "ACCEPT"
    }
    public_bots {
      bot_category = "page_preview"
      action       = "ACCEPT"
    }
    public_bots {
      bot_category = "search_engine_bot"
      action       = "ACCEPT"
    }
    public_bots {
      bot_category = "site_monitor"
      action       = "ACCEPT"
    }
    public_bots {
      bot_category = "tool"
      action       = "ACCEPT"
    }

    absolute_bots_act = "LOG"

    bot_tagging {
      request_header_key      = "A"
      traffic_characteristics = "CUSTOMIZE_BOT"
    }
    bot_tagging {
      request_header_key      = "b"
      traffic_characteristics = "PUBLIC_BOT"
    }
  }

  web_config {
    act                        = "NO_USE"
    browser_analyse_switch     = "ON"
    auto_tool_switch           = "ON"
    crack_analyse_switch       = "ON"
    page_debug_switch          = "OFF"
    interaction_analyse_switch = "OFF"
    ajax_exception_switch      = "ON"
  }

  traffic_detection {
    start_time = "00:00"
    end_time   = "23:59"
    action     = "NO_USE"
    whitelist  = ["1.1.1.1"]
  }
}

data "wangsu_waap_bot_configs" "demo" {
  domain_list = ["waap.example.com"]
}


output "bot_configs" {
  value = data.wangsu_waap_bot_configs.demo
}