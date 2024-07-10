terraform {
  required_providers {
    wangsu = {
      source  = "wangsu-api/wangsu"
      version = "1.0.3"
    }
  }
}

provider "wangsu" {
  secret_id  = ""
  secret_key = ""
}

data "wangsu_cdn_domains" "myDomainList" {
  domain_name       = ["www.mydomain.com"]
  service_type      = "download"
  total_page_number = 1
  page_size         = 10
  total_count       = 1
  page_number       = 1
}

output "domain_list" {
  value = data.wangsu_cdn_domains.myDomainList
}