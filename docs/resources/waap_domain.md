---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "wangsu_waap_domain Resource - wangsu"
subcategory: "Security"
description: |-
    Use this resource to create a WAAP domain.
    This resource allows you to create a WAAP domain.
---

# wangsu_waap_domain (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `target_domains` (List of String) Hostnames to be accessed.

### Optional

- `api_defend_config` (Block List) API security. (see [below for nested schema](#nestedblock--api_defend_config))
- `block_config` (Block List) IP/Geo blocking. (see [below for nested schema](#nestedblock--block_config))
- `bot_manage_config` (Block List) Bot management. (see [below for nested schema](#nestedblock--bot_manage_config))
- `customize_rule_config` (Block List) Custom rules. (see [below for nested schema](#nestedblock--customize_rule_config))
- `dms_defend_config` (Block List) DDoS protection. (see [below for nested schema](#nestedblock--dms_defend_config))
- `intelligence_config` (Block List) Threat intelligence. (see [below for nested schema](#nestedblock--intelligence_config))
- `rate_limit_config` (Block List) Rate limiting. (see [below for nested schema](#nestedblock--rate_limit_config))
- `waf_defend_config` (Block List) WAF. (see [below for nested schema](#nestedblock--waf_defend_config))
- `whitelist_config` (Block List) Whitelist. (see [below for nested schema](#nestedblock--whitelist_config))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--api_defend_config"></a>
### Nested Schema for `api_defend_config`

Required:

- `config_switch` (String) API security switch.
ON: Enabled
OFF: Disabled


<a id="nestedblock--block_config"></a>
### Nested Schema for `block_config`

Required:

- `config_switch` (String) IP/Geo switch.
ON: Enabled
OFF: Disabled


<a id="nestedblock--bot_manage_config"></a>
### Nested Schema for `bot_manage_config`

Required:

- `config_switch` (String) Bot management switch.
ON: Enabled
OFF: Disabled
- `public_bots_act` (String) Known Bots action.
NO_USE: not used
BLOCK: Deny
LOG: Log
ACCEPT: Skip
- `scene_analyse_switch` (String) Client-based detection function switch.
ON: Enabled
OFF: Disabled
- `ua_bots_act` (String) User-Agent based detection action.
NO_USE: Not used
BLOCK: Deny
LOG: Log
ACCEPT: Skip
- `web_risk_config` (Block List, Min: 1) Browser Bot defense. (see [below for nested schema](#nestedblock--bot_manage_config--web_risk_config))

<a id="nestedblock--bot_manage_config--web_risk_config"></a>
### Nested Schema for `bot_manage_config.web_risk_config`

Required:

- `act` (String) Action.
NO_USE: Not used
BLOCK: Deny
LOG: Log



<a id="nestedblock--customize_rule_config"></a>
### Nested Schema for `customize_rule_config`

Required:

- `config_switch` (String) Custom rules switch.
ON: Enabled
OFF: Disabled


<a id="nestedblock--dms_defend_config"></a>
### Nested Schema for `dms_defend_config`

Required:

- `ai_switch` (String) DDoS AI intelligent protection switch.
ON: Enabled
OFF: Disabled
- `config_switch` (String) DDoS protection switch.
ON: Enabled
OFF: Disabled
- `protection_mode` (String) DDoS protection mode.
AI_DEPOSIT: Managed Auto-Protect
UNDER_ATTACK: I'm Under Attack


<a id="nestedblock--intelligence_config"></a>
### Nested Schema for `intelligence_config`

Required:

- `config_switch` (String) Threat intelligence switch.
ON: Enabled
OFF: Disabled
- `info_cate_act` (Block List, Min: 1) Attack risk type action. (see [below for nested schema](#nestedblock--intelligence_config--info_cate_act))

<a id="nestedblock--intelligence_config--info_cate_act"></a>
### Nested Schema for `intelligence_config.info_cate_act`

Required:

- `attack_source` (String) Attack resource risk action.
NO_USE: Not used
BLOCK: Deny
LOG: Log
- `industry` (String) Industry attack risk action.
NO_USE: Not used
BLOCK: Deny
LOG: Log
- `spec_attack` (String) Specific attack risk action.
NO_USE: Not used
BLOCK: Deny
LOG:Log



<a id="nestedblock--rate_limit_config"></a>
### Nested Schema for `rate_limit_config`

Required:

- `config_switch` (String) Rate limiting switch.
ON: Enabled
OFF: Disabled


<a id="nestedblock--waf_defend_config"></a>
### Nested Schema for `waf_defend_config`

Required:

- `config_switch` (String) WAF protection switch.
ON: Enabled
OFF: Disabled
- `defend_mode` (String) WAF protection Mode.
BLOCK: Interception
LOG: Observation
- `rule_update_mode` (String) Ruleset pattern. 
MANUAL: Manual
AUTO: Automatic


<a id="nestedblock--whitelist_config"></a>
### Nested Schema for `whitelist_config`

Required:

- `config_switch` (String) Whitelist switch.
ON: Enabled
OFF: Disabled