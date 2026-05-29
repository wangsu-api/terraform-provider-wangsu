# CLAUDE.md — terraform-provider-wangsu

## Project Overview

This is a **Terraform Provider** for [Wangsu](https://www.wangsu.com/) (ChinaNetCenter), a Chinese CDN and cloud security company. It enables infrastructure-as-code management of Wangsu services including CDN, WAAP (Web Application and API Protection), SSL, IAM, and more.

- **Module**: `github.com/wangsu-api/terraform-provider-wangsu`
- **Language**: Go 1.22
- **Registry address**: `registry.terraform.io/wangsu-api/wangsu`
- **SDK**: `github.com/wangsu-api/wangsu-sdk-go`

---

## Repository Structure

```
terraform-provider-wangsu/
├── main.go                          # Entry point, registers the provider
├── go.mod / go.sum                  # Go module dependencies
├── wangsu/
│   ├── provider.go                  # Provider schema, resource/data-source registration
│   ├── common/
│   │   ├── common.go                # Utility helpers (IsContains, Int64ToStr)
│   │   ├── helper.go                # Additional helper functions
│   │   ├── provider.go              # Provider-level shared logic
│   │   └── validators.go            # Schema validators (e.g., ValidateAllowedStringValue)
│   ├── connectivity/
│   │   └── client.go                # WangSuClient: lazy-initialized SDK sub-clients per service
│   ├── http/
│   │   └── client.go                # Low-level HTTP client wrapper
│   └── services/                    # One subdirectory per product area
│       ├── cdn/
│       │   ├── domain/              # CDN domain resource & data sources
│       │   ├── edgehostname/        # CDN edge hostname resource & data sources
│       │   └── property/            # CDN property & deployment resources
│       ├── ssl/
│       │   ├── certificate/         # SSL certificate resource & data sources
│       │   └── certificateapplication/ # SSL cert application resource & data sources
│       ├── waap/
│       │   ├── domain/              # WAAP domain resource & data sources
│       │   ├── whitelist/           # WAAP whitelist rules
│       │   ├── customizerule/       # WAAP custom rules
│       │   ├── ratelimit/           # WAAP rate limiting
│       │   ├── waf/                 # WAF config & rule exceptions
│       │   ├── bot/                 # Bot protection config
│       │   ├── bot-scene-whitelist/ # Bot scene whitelists
│       │   ├── ddosprotection/      # DDoS protection config
│       │   ├── intelligence/        # Threat intelligence config
│       │   ├── predeploy/           # Pre-deploy rule staging resources
│       │   ├── share-whitelist/     # Shared whitelists
│       │   ├── share-customizerule/ # Shared custom rules
│       │   ├── share-customizebot/  # Shared bot custom rules
│       │   └── common.go            # WAAP-shared utilities
│       ├── iam/
│       │   ├── policy/              # IAM policy resource & data sources
│       │   └── user/                # IAM user resource & data sources
│       ├── monitor/
│       │   └── rule/                # Monitor real-time alert rules
│       └── appa/
│           └── domain/              # AppA (application acceleration) domain
├── example/                         # Example Terraform configurations per resource
│   ├── cdn/
│   ├── ssl/
│   ├── waap/
│   ├── iam/
│   ├── monitor/
│   └── appa/
```

---

## Provider Configuration

```hcl
terraform {
  required_providers {
    wangsu = {
      source = "registry.terraform.io/wangsu-api/wangsu"
    }
  }
}

provider "wangsu" {
  secret_id    = "your-secret-id"   # or env: WANGSU_SECRET_ID
  secret_key   = "your-secret-key"  # or env: WANGSU_SECRET_KEY
  protocol     = "https"            # or env: WANGSU_PROTOCOL (http/https)
  domain       = ""                 # or env: WANGSU_DOMAIN (default: open.chinanetcenter.com)
  service_type = ""                 # Required if multiple security services purchased
}
```

---

## Available Resources & Data Sources

### CDN
| Resource / Data Source | Description |
|---|---|
| `wangsu_cdn_domain` | CDN domain configuration (caching, headers, SSL, origin) |
| `wangsu_cdn_property` | CDN property config |
| `wangsu_cdn_property_deployment` | Deploy a CDN property |
| `wangsu_cdn_edge_hostname` | CDN edge hostname |
| `wangsu_cdn_domains` *(data)* | List CDN domains |
| `wangsu_cdn_domain_detail` *(data)* | Get CDN domain details |
| `wangsu_cdn_properties` *(data)* | List CDN properties |
| `wangsu_cdn_property_detail` *(data)* | Get CDN property details |
| `wangsu_cdn_property_deployments` *(data)* | List CDN deployments |
| `wangsu_cdn_property_deployment_detail` *(data)* | Get CDN deployment details |
| `wangsu_cdn_edge_hostnames` *(data)* | List edge hostnames |
| `wangsu_cdn_edge_hostname_detail` *(data)* | Get edge hostname details |

### SSL
| Resource / Data Source | Description |
|---|---|
| `wangsu_ssl_certificate` | Upload/manage SSL certificate |
| `wangsu_ssl_certificate_application` | Apply for an SSL certificate |
| `wangsu_ssl_certificates` *(data)* | List SSL certificates |
| `wangsu_ssl_certificate_detail` *(data)* | Get certificate details |
| `wangsu_ssl_certificate_applications` *(data)* | List cert applications |
| `wangsu_ssl_certificate_application_detail` *(data)* | Get cert application details |

### WAAP (Web Application and API Protection)
| Resource / Data Source | Description |
|---|---|
| `wangsu_waap_domain` | WAAP domain |
| `wangsu_waap_domain_copy` | Copy WAAP domain configuration |
| `wangsu_waap_whitelist` | WAAP whitelist rule |
| `wangsu_waap_customizerule` | WAAP custom protection rule |
| `wangsu_waap_ratelimit` | WAAP rate limiting rule |
| `wangsu_waap_waf_config` | WAF configuration |
| `wangsu_waap_waf_rule_exception` | WAF rule exception |
| `wangsu_waap_bot_config` | Bot protection config |
| `wangsu_waap_bot_scene_whitelist` | Bot scene whitelist |
| `wangsu_waap_threat_intelligence_config` | Threat intelligence config |
| `wangsu_waap_share_whitelist` | Shared whitelist |
| `wangsu_waap_share_customizerule` | Shared custom rule |
| `wangsu_waap_share_customizebot` | Shared bot custom rule |
| `wangsu_waap_pre_deploy_whitelist` | Pre-deploy whitelist staging |
| `wangsu_waap_pre_deploy_custom_rule` | Pre-deploy custom rule staging |
| `wangsu_waap_pre_deploy_rate_limiting` | Pre-deploy rate limiting staging |
| `wangsu_waap_pre_deploy_waf` | Pre-deploy WAF staging |
| `wangsu_waap_pre_deploy_ddos_protection` | Pre-deploy DDoS staging |
| `wangsu_waap_domains` *(data)* | List WAAP domains |
| `wangsu_waap_whitelists` *(data)* | List WAAP whitelists |
| `wangsu_waap_customizerules` *(data)* | List WAAP custom rules |
| `wangsu_waap_ratelimits` *(data)* | List WAAP rate limits |
| `wangsu_waap_waf_configs` *(data)* | List WAF configs |
| `wangsu_waap_ddos_protection_configs` *(data)* | List DDoS configs |
| `wangsu_waap_threat_intelligence_configs` *(data)* | List threat intelligence configs |
| `wangsu_waap_bot_configs` *(data)* | List bot configs |
| `wangsu_waap_share_whitelists` *(data)* | List shared whitelists |
| `wangsu_waap_bot_scene_whitelists` *(data)* | List bot scene whitelists |
| `wangsu_waap_share_customizebots` *(data)* | List shared bot rules |
| `wangsu_waap_share_customizerules` *(data)* | List shared custom rules |

### IAM
| Resource / Data Source | Description |
|---|---|
| `wangsu_iam_policy` | IAM policy |
| `wangsu_iam_policy_attachment` | Attach IAM policy to user |
| `wangsu_iam_user` | IAM user |
| `wangsu_iam_policy_detail` *(data)* | Get IAM policy details |
| `wangsu_iam_user_detail` *(data)* | Get IAM user details |
| `wangsu_iam_users` *(data)* | List IAM users |

### Monitor
| Resource / Data Source | Description |
|---|---|
| `wangsu_monitor_realtime_rule` | Real-time monitoring alert rule |
| `wangsu_monitor_realtime_rules_detail` *(data)* | Get monitoring rule details |

### AppA (Application Acceleration)
| Resource / Data Source | Description |
|---|---|
| `wangsu_appa_domain` | AppA domain |
| `wangsu_appa_domain_detail` *(data)* | Get AppA domain details |

---

## Architecture Patterns

### Adding a New Resource

1. Create a directory under `wangsu/services/<product>/<feature>/`
2. Implement two files:
   - `resource_<name>.go` — implement `ResourceXxx()` returning `*schema.Resource` with CRUD functions
   - `data_source_<name>.go` — implement `DataSourceXxx()` returning `*schema.Resource` with Read function
3. Add a lazy-initialized client method to `wangsu/connectivity/client.go` using the pattern:
   ```go
   func (me *WangSuClient) UseXxxClient() *xxx.Client {
       if me.xxxConn != nil {
           return me.xxxConn
       }
       me.xxxConn, _ = xxx.NewClient(me.Credential, me.HttpProfile)
       return me.xxxConn
   }
   ```
4. Register the resource and data source in `wangsu/provider.go` in `ResourcesMap` and `DataSourcesMap`
5. Add an example configuration in `example/<product>/<feature>/main.tf`

### Client Connectivity Pattern

`WangSuClient` in `wangsu/connectivity/client.go` holds lazy-initialized connections per SDK sub-client. Each `UseXxxClient()` method initializes the connection on first use. This avoids unnecessary API connections.

### Resource Implementation Pattern

Each resource file follows this structure:
- Define `ResourceXxx()` returning `*schema.Resource` with `Schema`, `CreateContext`, `ReadContext`, `UpdateContext`, `DeleteContext`
- Access the SDK client via: `meta.(*wangsu.WangSuClient).GetAPIV3Conn().UseXxxClient()`
- Data sources only implement `ReadContext`

---

## Building & Running

```bash
# Build the provider
go build -o terraform-provider-wangsu

# Run in debug mode (for use with Delve)
./terraform-provider-wangsu -debuggable

# Generate docs and format examples
go generate ./...

# Run tests
go test ./...
```

### Installing locally for development

```bash
# Build and place in local plugin cache
go build -o ~/.terraform.d/plugins/registry.terraform.io/wangsu-api/wangsu/0.0.1/linux_amd64/terraform-provider-wangsu
```

Then in your Terraform config use:
```hcl
terraform {
  required_providers {
    wangsu = {
      source  = "registry.terraform.io/wangsu-api/wangsu"
      version = "0.0.1"
    }
  }
}
```

---

## Key Dependencies

| Package | Purpose |
|---|---|
| `github.com/hashicorp/terraform-plugin-sdk/v2` | Terraform Plugin SDK (schema, resource lifecycle) |
| `github.com/wangsu-api/wangsu-sdk-go` | Official Wangsu Go SDK for all API calls |
| `github.com/alibabacloud-go/tea` | HTTP/retry utilities used by SDK |
| `golang.org/x/net` | Extended networking support |

---

## Environment Variables

| Variable | Description |
|---|---|
| `WANGSU_SECRET_ID` | API access key ID |
| `WANGSU_SECRET_KEY` | API secret key |

---

## Notes for AI Assistants

- All service sub-clients in `connectivity/client.go` follow a consistent lazy-init pattern — follow the existing pattern exactly when adding new ones.
- Resource names use prefix `wangsu_` followed by `<product>_<feature>` (e.g., `wangsu_waap_whitelist`).
- The `example/` directory mirrors the resource hierarchy and is auto-formatted via `go generate`.
- WAAP "pre-deploy" resources are staging resources that queue rule changes before they are deployed live.
- "Share" WAAP resources (e.g., `wangsu_waap_share_whitelist`) are tenant-level shared rules, distinct from domain-level rules.
