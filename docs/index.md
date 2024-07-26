---
layout: "wangsu"
page_title: "Provider: wangsu"
sidebar_current: "docs-wangsu-index"
description: |-
  The Wangsu provider is used to interact with many resources supported by wangsu. The provider needs to be configured with the proper credentials before it can be used.
---

# Wangsu Provider

The Wangsu provider is used to interact with many resources supported by [Wangsu](https://www.wangsu.com).
The provider needs to be configured with the proper credentials before it can be used.

Use the navigation on the left to read about the available resources.

## Example Usage

```hcl
terraform {
  required_providers {
    wangsu = {
      source = "wangsu/wangsu"
    }
  }
}

# Configure the Wangsu Provider
resource "wangsu_cdn_domain" "mydomain" {
  version       = "1.0.2"
  domain_name   = "www.mydomain.com"
  service_type  = "download"
  service_areas = "cn"

  origin_config {
    origin_ips                 = "122.22.22.221"
    default_origin_host_header = "weixin.qq.com"
  }
}
```

## Authentication

The wangsu provider offers a flexible means of providing credentials for authentication.
The following methods are supported, in this order, and explained below:

- Static credentials
- Environment variables

### Static credentials

!> **Warning:** Hard-coding credentials into any Terraform configuration is not
recommended, and risks secret leakage should this file ever be committed to a
public version control system.

Static credentials can be provided by adding an `secret_id` `secret_key` and `region` in-line in the wangsu provider block:

Usage:

```hcl
provider "wangsu" {
  secret_id  = "my-secret-id"
  secret_key = "my-secret-key"
}
```

### Environment variables

You can provide your credentials via `WANGSU_SECRET_ID` and `WANGSU_SECRET_KEY` environment variables,

```hcl
provider "wangsu" {
}
```

Usage:

```shell
$ export WANGSU_SECRET_ID="my-secret-id"
$ export WANGSU_SECRET_KEY="my-secret-key"
$ terraform plan
```


## Argument Reference

In addition to generic provider arguments (e.g. alias and version), the following arguments are supported in the wangsu provider block:

* `secret_id` - (Optional) This is the wangsu secret id. It must be provided, but it can also be sourced from the `WANGSU_SECRET_KEY` environment variable.
* `secret_key` - (Optional) This is the wangsu secret key. It must be provided, but it can also be sourced from the `WANGSU_SECRET_KEY` environment variable.
* `protocol` - (Optional) The protocol of the API request. Valid values: `HTTP` and `HTTPS`. Default is `HTTPS`.
* `domain` - (Optional) The root domain of the API request, Default is `open.chinanetcenter.com`.
* `service_type` (Optional) The service type of the accelerated domain name. The value can be: appa: Application Acceleration