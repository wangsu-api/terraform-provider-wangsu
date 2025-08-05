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

data "wangsu_waap_waf_configs" "demo" {
  domain_list = ["waap.example.com"]
}

output "waf_configs" {
  value = data.wangsu_waap_waf_configs.demo
}