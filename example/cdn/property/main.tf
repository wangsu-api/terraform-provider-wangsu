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

resource "wangsu_cdn_property" "example_property" {
  origins = jsonencode([{
    http_port  = "80"
    https_port = "443"
    name       = "xxxx_114.80.165.148"
    servers = [{
      priority = 1
      server   = "1.1.1.1"
      weight   = 10
    }]
  }])
  property_comment = "property comment"
  property_name    = "example_property_name"
  rules = jsonencode({
    cachePhase = [{
      action = {
        name = "cache"
        options = {
          defaultCachetime = "30d"
        }
      }
      condition = "http.request.uri.path == /"
      enabled   = true
      name      = "test7"
      priority  = 1
      }, {
      action = {
        name = "cacheKey"
        options = {
          cacheHost = "osx.dpfile.com"
        }
      }
      enabled  = true
      name     = "test3"
      priority = 1
    }]
    originPhase = [{
      action = {
        name = "originReqHeaderModify"
        options = {
          headerName  = "accept-encoding"
          headerValue = "gzip"
          operator    = "SET"
        }
      }
      enabled  = true
      name     = "test5"
      priority = 1
    }]
    requestPhase = [{
      action = {
        name = "accessControl"
        options = {
          behavior = "block"
        }
      }
      condition = "http.request.uri.path == /a.html"
      enabled   = true
      name      = "test2"
      priority  = 1
    }]
    responsePhase = [{
      action = {
        name = "clientRespHeaderModify"
        options = {
          headerName = "cache-control"
          operator   = "DEL"
        }
      }
      description = "description_test6"
      enabled     = true
      name        = "test6"
      priority    = 1
    }]
  })
  service_type = "wsa"
  variables = jsonencode([{
    action = {
      name = "setVariable"
      options = {
        value   = "copy(\"12345\")"
        varName = "var1"
      }
    }
    condition = "http.request.method == OPTION|PUT|DELETE|PATCH"
    enabled   = true
    name      = "testVariables"
    priority  = 1
  }])
  version_comment = "version comment"
  hostnames {
    hostname = "www.example.com"
    certificates {
      certificate_id    = 1614905
      certificate_usage = "default_sni"
    }
    certificates {
      certificate_id    = 1614906
      certificate_usage = "dual_sni"
    }
    certificates {
      certificate_id    = 1614907
      certificate_usage = "gm_sm2_enc"
    }
    certificates {
      certificate_id    = 1614908
      certificate_usage = "gm_sm2_sign"
    }
    default_origin {
      host       = "example.com"
      http_port  = 80
      https_port = 443
      ip_version = "dual"
      servers    = ["example1.com"]
    }
    edge_hostname {
      comment              = "here is edge-hostname comment"
      edge_hostname_prefix = "www.example.com"
      edge_hostname_suffix = "wscncdn.com"
    }
  }
}