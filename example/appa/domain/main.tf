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


resource "wangsu_appa_domain" "appa1" {
  domain_name  = "20240710001.conftest.com"
  service_type = "appa"
  origin_config {
    level    = 1
    strategy = "robin"
    origin {
      origin_ip = "1.1.1.6"
      weight    = 11
    }
    origin {
      origin_ip = "2.2.2.6"
      weight    = 22
    }
  }
  origin_config {
    level    = 2
    strategy = "fast"
    origin {
      origin_ip = "3.3.3.6"
      weight    = 33
    }
  }
  http_ports  = ["1000", "1002"]
  https_ports = ["2000", "2002"]
}

data "wangsu_appa_domain_detail" "appa1-data" {
  domain_name = wangsu_appa_domain.appa1.domain_name
}

output "show-appa1" {
  value = data.wangsu_appa_domain_detail.appa1-data
}