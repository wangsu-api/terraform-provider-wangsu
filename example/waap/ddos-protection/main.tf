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

data "wangsu_waap_ddos_protection_configs" "demo" {
  domain_list = ["example.waap.com"]
}

output "ddos_protection_configs" {
  value = data.wangsu_waap_ddos_protection_configs.demo
}