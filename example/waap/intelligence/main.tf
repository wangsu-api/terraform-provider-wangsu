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

resource "wangsu_waap_threat_intelligence_config" "demo" {
  domain = "waap.example.com"

  config_list  {
    id = "1626496419345657871"
    action = "BLOCK"
  }

  config_list  {
    id = "1626496419345657874"
    action = "BLOCK"
  }
}

data "wangsu_waap_threat_intelligence_configs" "demo" {
  domain_list = ["waap.example.com"]
}


output "threat_intelligence_configs" {
  value = data.wangsu_waap_threat_intelligence_configs.demo
}