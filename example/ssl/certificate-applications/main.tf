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

data "wangsu_ssl_certificate_applications" "example" {
  page_number      = 1
  page_size        = 100
  order_id         = "SO20251029174850285.474936729"
  order_status     = ["ACCEPT_SUCCESS", "APPLYING", "CANCELED"]
  certificate_name = "20251029001.conftest.com"
  domain           = "20251029001.conftest.com"
  start_time       = "2025-10-29 10:20:00"
  end_time         = "2025-10-29 19:40:00"
}

output "data" {
  value = data.wangsu_ssl_certificate_applications.example
}