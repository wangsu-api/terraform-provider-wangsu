terraform {
  required_providers {
    wangsu = {
      source = "registry.terraform.io/wangsustack/wangsu"
      version = "1.0.0"
    }
  }
}

provider "wangsu" {
  secret_id  = ""
  secret_key = ""
}

resource "wangsu_cdn_domain" "mydomain" {
  version      = "1.0.0"
  domain_name  = "www.mydomain.com"
  service_type = "download"
  service_areas = "cn"

  origin_config {
    origin_ips = "122.22.22.221"
    default_origin_host_header = "weixin.qq.com"
  }
}