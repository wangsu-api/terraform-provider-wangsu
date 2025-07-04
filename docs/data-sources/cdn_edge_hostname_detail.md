---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "wangsu_cdn_edge_hostname_detail Data Source - wangsu"
subcategory: "CDN"
description: |-
  Use this data source to query detailed information of edge-hostname.
---

# wangsu_cdn_edge_hostname_detail (Data Source)

Use this data source to query detailed information of edge-hostname.

## Example Usage

```hcl
data "wangsu_cdn_edge_hostname_detail" "example_edge_hostname_detail" {
  edge_hostname = "www.example.com.wscncdn.com"
}
output "edge_hostname_detail" {
  value = data.wangsu_cdn_edge_hostname_detail.example_edge_hostname_detail
}
```


<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `edge_hostname` (String) edge-hostname

### Read-Only

- `code` (String) Response code, 0 means successful.
- `data` (List of Object) Response data. (see [below for nested schema](#nestedatt--data))
- `id` (String) The ID of this resource.
- `message` (String) Response error message if failed.

<a id="nestedatt--data"></a>
### Nested Schema for `data`

Read-Only:

- `allow_china_cdn` (String) Allow China CDN; values: [0,1].
- `comment` (String) Edge-Hostname comment.
- `creation_time` (String) Creation time in RFC3339 format.
- `deploy_status` (String) Deploy status; possible values: [pending, deploying, success, fail].
- `dns_service_status` (String) DNS service status; possible values: [inactive, active].
- `edge_hostname` (String) Edge-Hostname Name.
- `edge_hostname_id` (Number) Edge-Hostname ID.
- `gdpr_compliant` (String) GDPR compliant; values: [0,1,2].
- `geo_fence` (String) Geo-fence; values: [global, inside_china_mainland, exclude_china_mainland].
- `hostnames` (List of Object) Associated hostnames. (see [below for nested schema](#nestedobjatt--data--hostnames))
- `last_update_time` (String) Last update time in RFC3339 format.
- `region_configs` (List of Object) Region configuration. (see [below for nested schema](#nestedobjatt--data--region_configs))

<a id="nestedobjatt--data--hostnames"></a>
### Nested Schema for `hostnames`

Read-Only:

- `hostname` (String) Hostname.
- `property_id` (Number) Property ID.
- `property_name` (String) Property name.
- `property_version` (Number) Property version.
- `target` (String) Deployment target.


<a id="nestedobjatt--data--region_configs"></a>
### Nested Schema for `region_configs`

Read-Only:

- `action_type` (String) Action type of configuration.
- `config_type` (String) Configuration type.
- `config_value` (String) Configuration value.
- `ip_protocol` (String) IP protocol.
- `region_id` (Number) Region ID.
- `ttl` (Number) TTL value.
- `weight` (Number) Weight of the configuration.
