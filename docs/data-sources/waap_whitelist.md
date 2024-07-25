---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "wangsu_waap_whitelist Data Source - wangsu"
subcategory: "Security"
description: |-
    Use this data source to query whitelist rules.
    This data source allows you to query whitelist rules
---

# wangsu_waap_whitelist (Data Source)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `domain_list` (List of String) Hostname list.

### Optional

- `rule_name` (String) Rule name, fuzzy query.

### Read-Only

- `data` (List of Object) Data. (see [below for nested schema](#nestedatt--data))
- `id` (String) The ID of this resource.

<a id="nestedatt--data"></a>
### Nested Schema for `data`

Read-Only:

- `conditions` (List of Object) (see [below for nested schema](#nestedobjatt--data--conditions))
- `description` (String)
- `domain` (String)
- `id` (String)
- `rule_name` (String)

<a id="nestedobjatt--data--conditions"></a>
### Nested Schema for `data.conditions`

Read-Only:

- `header_conditions` (List of Object) (see [below for nested schema](#nestedobjatt--data--conditions--header_conditions))
- `ip_or_ips_conditions` (List of Object) (see [below for nested schema](#nestedobjatt--data--conditions--ip_or_ips_conditions))
- `path_conditions` (List of Object) (see [below for nested schema](#nestedobjatt--data--conditions--path_conditions))
- `referer_conditions` (List of Object) (see [below for nested schema](#nestedobjatt--data--conditions--referer_conditions))
- `ua_conditions` (List of Object) (see [below for nested schema](#nestedobjatt--data--conditions--ua_conditions))
- `uri_conditions` (List of Object) (see [below for nested schema](#nestedobjatt--data--conditions--uri_conditions))

<a id="nestedobjatt--data--conditions--header_conditions"></a>
### Nested Schema for `data.conditions.header_conditions`

Read-Only:

- `key` (String)
- `match_type` (String)
- `value_list` (List of String)


<a id="nestedobjatt--data--conditions--ip_or_ips_conditions"></a>
### Nested Schema for `data.conditions.ip_or_ips_conditions`

Read-Only:

- `ip_or_ips` (List of String)
- `match_type` (String)


<a id="nestedobjatt--data--conditions--path_conditions"></a>
### Nested Schema for `data.conditions.path_conditions`

Read-Only:

- `match_type` (String)
- `paths` (List of String)


<a id="nestedobjatt--data--conditions--referer_conditions"></a>
### Nested Schema for `data.conditions.referer_conditions`

Read-Only:

- `match_type` (String)
- `referer` (List of String)


<a id="nestedobjatt--data--conditions--ua_conditions"></a>
### Nested Schema for `data.conditions.ua_conditions`

Read-Only:

- `match_type` (String)
- `ua` (List of String)


<a id="nestedobjatt--data--conditions--uri_conditions"></a>
### Nested Schema for `data.conditions.uri_conditions`

Read-Only:

- `match_type` (String)
- `uri` (List of String)