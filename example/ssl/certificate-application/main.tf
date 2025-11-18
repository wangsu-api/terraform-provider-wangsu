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

resource "wangsu_ssl_certificate_application" "example" {
  algorithm           = "RSA2048"
  auto_deploy         = "true"
  auto_renew          = "false"
  auto_validate       = "true"
  certificate_brand   = "LE"
  certificate_spec    = "LetsEncryptDVFree"
  certificate_type    = "DV"
  description         = "This is description"
  domain_type         = "single"
  org_validate_method = "self_validate"
  validate_method     = "DNS"
  validity_days       = 90
  admin {
    email      = "test1@email.com"
    first_name = "firstName1"
    last_name  = "lastName1"
    phone      = "12345678901"
    title      = "title1"
  }
  dns_provider_infos {
    dns_api_access = jsonencode({
      accessKey = "123456"
      secretKey = "123456"
    })
    dns_provider_code     = "CloudDNS"
    domain                = "www.example.com"
    enable_dns_alias_mode = "true"
    validate_alias_domain = "alias.example.com"
  }
  identification_info {
    city                      = "your city"
    common_name               = "www.example.com"
    company                   = "your company"
    country                   = "CN"
    department                = "cdn"
    email                     = "test@email.com"
    phone                     = "12345678901"
    postal_code               = "123456"
    state                     = "FJ"
    street                    = "software part"
    street1                   = "#18-13"
    subject_alternative_names = ["www.example.com"]
  }
  tech {
    email      = "test2@email.com"
    first_name = "firstName2"
    last_name  = "lastName2"
    phone      = "12345678902"
    title      = "title2"
  }
}
