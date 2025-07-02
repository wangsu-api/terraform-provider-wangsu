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

resource "wangsu_cdn_edge_hostname" "example_edge_hostname" {
  comment       = "I am a comment"
  edge_hostname = "www.example.com.wscncdn.com"
  geo_fence     = "global"
  region_configs {
    action_type  = "redirect"
    config_value = "1.1.1.1"
    region_id    = 29
    ttl          = 600
    weight       = 100
  }
  region_configs {
    action_type  = "deliver"
    ip_protocol  = "0"
    region_id    = 1
    ttl          = 60
    weight       = 100
  }
  region_configs {
    action_type  = "reject"
    region_id    = 43
    ttl          = 66
    weight       = 100
  }
}
