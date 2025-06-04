package domain

import (
	"context"
	"errors"
	"fmt"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wangsuCommon "github.com/wangsu-api/terraform-provider-wangsu/wangsu/common"
	"github.com/wangsu-api/terraform-provider-wangsu/wangsu/services/waap"
	waapDomain "github.com/wangsu-api/wangsu-sdk-go/wangsu/waap/domain"
	"log"
	"time"
)

func ResourceWaapDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWaapDomainCreate,
		ReadContext:   resourceWaapDomainRead,
		UpdateContext: resourceWaapDomainUpdate,
		DeleteContext: resourceWaapDomainDelete,

		Schema: map[string]*schema.Schema{
			"waf_defend_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "WAF.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_update_mode": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Ruleset pattern. <br/>MANUAL: Manual<br/>AUTO: Automatic",
						},
						"config_switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "WAF protection switch.<br/>ON: Enabled<br/>OFF: Disabled",
						},
						"defend_mode": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "WAF protection Mode.<br/>BLOCK: Interception<br/>LOG: Observation",
						},
					},
				},
			},
			"customize_rule_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Custom rules.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"config_switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Custom rules switch.<br/>ON: Enabled<br/>OFF: Disabled",
						},
					},
				},
			},
			"api_defend_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "API security.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"config_switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "API security switch.<br/>ON: Enabled<br/>OFF: Disabled",
						},
					},
				},
			},
			"whitelist_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Whitelist.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"config_switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Whitelist switch.<br/>ON: Enabled<br/>OFF: Disabled",
						},
					},
				},
			},
			"target_domains": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Hostnames to be accessed.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"block_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "IP/Geo blocking.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"config_switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "IP/Geo switch.<br/>ON: Enabled<br/>OFF: Disabled",
						},
					},
				},
			},
			"dms_defend_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "DDoS protection.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"config_switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "DDoS protection switch.<br/>ON: Enabled<br/>OFF: Disabled",
						},
						"protection_mode": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "DDoS protection mode.<br/>LOOSE: Loose<br/>MODERATE: Moderate<br/>STRICT: Strict",
						},
						"ai_switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "DDoS AI intelligent protection switch.<br/>ON: Enabled<br/>OFF: Disabled",
						},
					},
				},
			},
			"intelligence_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Threat intelligence.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"info_cate_act": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Attack risk type action.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"attack_source": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Attack resource risk action.<br/>NO_USE: Not used<br/>BLOCK: Deny<br/>LOG: Log",
									},
									"spec_attack": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Specific attack risk action.<br/>NO_USE: Not used<br/>BLOCK: Deny<br/>LOG:Log",
									},
									"industry": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Industry attack risk action.<br/>NO_USE: Not used<br/>BLOCK: Deny<br/>LOG: Log",
									},
								},
							},
						},
						"config_switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Threat intelligence switch.<br/>ON: Enabled<br/>OFF: Disabled",
						},
					},
				},
			},
			"bot_manage_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Bot management.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"public_bots_act": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Known Bots action.<br/>NO_USE: not used<br/>BLOCK: Deny<br/>LOG: Log<br/>ACCEPT: Skip",
						},
						"config_switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Bot management switch.<br/>ON: Enabled<br/>OFF: Disabled",
						},
						"ua_bots_act": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "User-Agent based detection action.<br/>NO_USE: Not used<br/>BLOCK: Deny<br/>LOG: Log<br/>ACCEPT: Skip",
						},
						"web_risk_config": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Browser Bot defense.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"act": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Action.<br/>NO_USE: Not used<br/>BLOCK: Deny<br/>LOG: Log",
									},
								},
							},
						},
						"scene_analyse_switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Client-based detection function switch.<br/>ON: Enabled<br/>OFF: Disabled",
						},
					},
				},
			},
			"rate_limit_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Rate limiting.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"config_switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Rate limiting switch.<br/>ON: Enabled<br/>OFF: Disabled",
						},
					},
				},
			},
		},
	}
}

func resourceWaapDomainCreate(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_waap_domain.create")

	var diags diag.Diagnostics
	request := &waapDomain.AccessDomainRequest{}

	if wafDefendConfig, ok := data.GetOk("waf_defend_config"); ok {
		request.WafDefendConfig = expandWafDefendConfig(wafDefendConfig.(interface{}))
	}

	if customizeRuleConfig, ok := data.GetOk("customize_rule_config"); ok {
		request.CustomizeRuleConfig = expandCustomizeRuleConfig(customizeRuleConfig.(interface{}))
	}

	if apiDefendConfig, ok := data.GetOk("api_defend_config"); ok {
		request.ApiDefendConfig = expandApiDefendConfig(apiDefendConfig.(interface{}))
	}

	if whitelistConfig, ok := data.GetOk("whitelist_config"); ok {
		request.WhitelistConfig = expandWhitelistConfig(whitelistConfig.(interface{}))
	}

	if targetDomains, ok := data.GetOk("target_domains"); ok {
		targetDomainsList := targetDomains.([]interface{})
		targetDomainsStr := make([]*string, len(targetDomainsList))
		for i, v := range targetDomainsList {
			str := v.(string)
			targetDomainsStr[i] = &str
		}
		request.TargetDomains = targetDomainsStr
	}

	if blockConfig, ok := data.GetOk("block_config"); ok {
		request.BlockConfig = expandBlockConfig(blockConfig.(interface{}))
	}

	if dmsDefendConfig, ok := data.GetOk("dms_defend_config"); ok {
		request.DmsDefendConfig = expandDmsDefendConfig(dmsDefendConfig.(interface{}))
	}

	if intelligenceConfig, ok := data.GetOk("intelligence_config"); ok {
		request.IntelligenceConfig = expandIntelligenceConfig(intelligenceConfig.(interface{}))
	}

	if botManageConfig, ok := data.GetOk("bot_manage_config"); ok {
		request.BotManageConfig = expandBotManageConfig(botManageConfig.(interface{}))
	}

	if rateLimitConfig, ok := data.GetOk("rate_limit_config"); ok {
		request.RateLimitConfig = expandRateLimitConfig(rateLimitConfig.(interface{}))
	}

	var response *waapDomain.AccessDomainResponse
	var err error
	err = resource.RetryContext(context, time.Duration(5)*time.Minute, func() *resource.RetryError {
		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseWaapDomainClient().AddDomain(request)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}
	if response == nil {
		data.SetId("")
		return nil
	}
	var ids = request.TargetDomains
	idsStr := make([]string, len(ids))
	for i, v := range ids {
		idsStr[i] = *v
	}
	data.SetId(wangsuCommon.DataResourceIdsHash(idsStr))
	log.Printf("resource.wangsu_waap_domain.create success")
	return diags
}

func resourceWaapDomainRead(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_waap_domain.read")

	var response *waapDomain.ListDomainInfoResponse
	var err error
	var diags diag.Diagnostics
	request := &waapDomain.ListDomainInfoRequest{}
	if v, ok := data.GetOk("target_domains"); ok {
		targetDomainsList := v.([]interface{})
		targetDomainsStr := make([]*string, len(targetDomainsList))
		for i, v := range targetDomainsList {
			str := v.(string)
			targetDomainsStr[i] = &str
		}
		request.DomainList = targetDomainsStr
	}
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseWaapDomainClient().GetDomainList(request)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	if response == nil {
		return nil
	}
	if len(response.Data) == 0 || nil == response.Data {
		data.SetId("")
		return nil
	}
	var item = response.Data[0]
	// 获取当前的 waf_defend_config 值
	wafDefendConfig := data.Get("waf_defend_config").([]interface{})
	if len(wafDefendConfig) > 0 && wafDefendConfig[0] != nil {
		wafDefendConfigMap := wafDefendConfig[0].(map[string]interface{})
		// 更新 config_switch 字段
		wafDefendConfigMap["config_switch"] = item.WafDefendSwitch
		// 将完整的 waf_defend_config 设置回去
		if err := data.Set("waf_defend_config", []interface{}{wafDefendConfigMap}); err != nil {
			return diag.FromErr(fmt.Errorf("error setting waf_defend_config for resource: %s", err))
		}
	}
	// Get and update api_defend_config
	apiDefendConfig := data.Get("api_defend_config").([]interface{})
	if len(apiDefendConfig) > 0 && apiDefendConfig[0] != nil {
		apiDefendConfigMap := apiDefendConfig[0].(map[string]interface{})
		apiDefendConfigMap["config_switch"] = item.ApiDefendSwitch
		if err := data.Set("api_defend_config", []interface{}{apiDefendConfigMap}); err != nil {
			return diag.FromErr(fmt.Errorf("error setting api_defend_config for resource: %s", err))
		}
	}

	// Get and update block_config
	blockConfig := data.Get("block_config").([]interface{})
	if len(blockConfig) > 0 && blockConfig[0] != nil {
		blockConfigMap := blockConfig[0].(map[string]interface{})
		blockConfigMap["config_switch"] = item.BlockSwitch
		if err := data.Set("block_config", []interface{}{blockConfigMap}); err != nil {
			return diag.FromErr(fmt.Errorf("error setting block_config for resource: %s", err))
		}
	}

	// Get and update dms_defend_config
	dmsDefendConfig := data.Get("dms_defend_config").([]interface{})
	if len(dmsDefendConfig) > 0 && dmsDefendConfig[0] != nil {
		dmsDefendConfigMap := dmsDefendConfig[0].(map[string]interface{})
		dmsDefendConfigMap["config_switch"] = item.DmsDefendSwitch
		if err := data.Set("dms_defend_config", []interface{}{dmsDefendConfigMap}); err != nil {
			return diag.FromErr(fmt.Errorf("error setting dms_defend_config for resource: %s", err))
		}
	}

	// Get and update intelligence_config
	intelligenceConfig := data.Get("intelligence_config").([]interface{})
	if len(intelligenceConfig) > 0 && intelligenceConfig[0] != nil {
		intelligenceConfigMap := intelligenceConfig[0].(map[string]interface{})
		intelligenceConfigMap["config_switch"] = item.IntelligenceSwitch
		if err := data.Set("intelligence_config", []interface{}{intelligenceConfigMap}); err != nil {
			return diag.FromErr(fmt.Errorf("error setting intelligence_config for resource: %s", err))
		}
	}

	// Get and update bot_manage_config
	botManageConfig := data.Get("bot_manage_config").([]interface{})
	if len(botManageConfig) > 0 && botManageConfig[0] != nil {
		botManageConfigMap := botManageConfig[0].(map[string]interface{})
		botManageConfigMap["config_switch"] = item.BotManageSwitch
		if err := data.Set("bot_manage_config", []interface{}{botManageConfigMap}); err != nil {
			return diag.FromErr(fmt.Errorf("error setting bot_manage_config for resource: %s", err))
		}
	}

	// Get and update rate_limit_config
	rateLimitConfig := data.Get("rate_limit_config").([]interface{})
	if len(rateLimitConfig) > 0 && rateLimitConfig[0] != nil {
		rateLimitConfigMap := rateLimitConfig[0].(map[string]interface{})
		rateLimitConfigMap["config_switch"] = item.RateLimitSwitch
		if err := data.Set("rate_limit_config", []interface{}{rateLimitConfigMap}); err != nil {
			return diag.FromErr(fmt.Errorf("error setting rate_limit_config for resource: %s", err))
		}
	}
	var ids = request.DomainList
	idsStr := make([]string, len(ids))
	for i, v := range ids {
		idsStr[i] = *v
	}
	data.SetId(wangsuCommon.DataResourceIdsHash(idsStr))
	return nil
}

func resourceWaapDomainUpdate(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_waap_domain.update")

	var diags diag.Diagnostics
	if data.HasChange("target_domains") {
		// 把domain强制刷回旧值，否则会有权限问题
		oldDomain, _ := data.GetChange("target_domains")
		_ = data.Set("target_domains", oldDomain)
		err := errors.New("Hostname cannot be changed.")
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	changeUnExcept := false
	if data.HasChange("waf_defend_config") {
		oldConfig, newConfig := data.GetChange("waf_defend_config")
		changeUnExcept = waap.CheckChangeExcept(oldConfig.([]interface{}), newConfig.([]interface{}), "config_switch")
		if changeUnExcept {
			_ = data.Set("waf_defend_config", oldConfig)
		}
	}
	if data.HasChange("customize_rule_config") {
		oldConfig, newConfig := data.GetChange("customize_rule_config")
		changeUnExcept = changeUnExcept || waap.CheckChangeExcept(oldConfig.([]interface{}), newConfig.([]interface{}), "config_switch")
		if changeUnExcept {
			_ = data.Set("customize_rule_config", oldConfig)
		}
	}
	if data.HasChange("api_defend_config") {
		oldConfig, newConfig := data.GetChange("api_defend_config")
		changeUnExcept = changeUnExcept || waap.CheckChangeExcept(oldConfig.([]interface{}), newConfig.([]interface{}), "config_switch")
		if changeUnExcept {
			_ = data.Set("api_defend_config", oldConfig)
		}
	}
	if data.HasChange("whitelist_config") {
		oldConfig, newConfig := data.GetChange("whitelist_config")
		changeUnExcept = changeUnExcept || waap.CheckChangeExcept(oldConfig.([]interface{}), newConfig.([]interface{}), "config_switch")
		if changeUnExcept {
			_ = data.Set("whitelist_config", oldConfig)
		}
	}
	if data.HasChange("block_config") {
		oldConfig, newConfig := data.GetChange("block_config")
		changeUnExcept = changeUnExcept || waap.CheckChangeExcept(oldConfig.([]interface{}), newConfig.([]interface{}), "config_switch")
		if changeUnExcept {
			_ = data.Set("block_config", oldConfig)
		}
	}
	if data.HasChange("dms_defend_config") {
		oldConfig, newConfig := data.GetChange("dms_defend_config")
		changeUnExcept = changeUnExcept || waap.CheckChangeExcept(oldConfig.([]interface{}), newConfig.([]interface{}), "config_switch")
		if changeUnExcept {
			_ = data.Set("dms_defend_config", oldConfig)
		}
	}
	if data.HasChange("intelligence_config") {
		oldConfig, newConfig := data.GetChange("intelligence_config")
		changeUnExcept = changeUnExcept || waap.CheckChangeExcept(oldConfig.([]interface{}), newConfig.([]interface{}), "config_switch")
		if oldConfig != nil && len(oldConfig.([]interface{})) != 0 && newConfig != nil && len(newConfig.([]interface{})) != 0 {
			changeUnExcept = changeUnExcept || data.HasChange("intelligence_config.0.info_cate_act")
		}
		if changeUnExcept {
			_ = data.Set("intelligence_config", oldConfig)
		}
	}
	if data.HasChange("bot_manage_config") {
		oldConfig, newConfig := data.GetChange("bot_manage_config")
		changeUnExcept = changeUnExcept || waap.CheckChangeExcept(oldConfig.([]interface{}), newConfig.([]interface{}), "config_switch")
		if oldConfig != nil && len(oldConfig.([]interface{})) != 0 && newConfig != nil && len(newConfig.([]interface{})) != 0 {
			changeUnExcept = changeUnExcept || data.HasChange("bot_manage_config.0.web_risk_config")
		}
		if changeUnExcept {
			_ = data.Set("bot_manage_config", oldConfig)
		}
	}
	if data.HasChange("rate_limit_config") {
		oldConfig, newConfig := data.GetChange("rate_limit_config")
		changeUnExcept = changeUnExcept || waap.CheckChangeExcept(oldConfig.([]interface{}), newConfig.([]interface{}), "config_switch")
		if changeUnExcept {
			_ = data.Set("rate_limit_config", oldConfig)
		}
	}
	if changeUnExcept {
		// 如果有变化，返回错误
		diags = append(diags, diag.FromErr(errors.New("Only 'config_switch' can be updated. Other changes require other resource."))...)
		return diags
	}
	request := &waapDomain.ModifyPolicyStatusRequest{}

	if targetDomains, ok := data.GetOk("target_domains"); ok {
		targetDomainsList := targetDomains.([]interface{})
		targetDomainsStr := make([]*string, len(targetDomainsList))
		for i, v := range targetDomainsList {
			str := v.(string)
			targetDomainsStr[i] = &str
		}
		request.DomainList = targetDomainsStr
	}

	if wafDefendConfig, ok := data.GetOk("waf_defend_config"); ok {
		request.WafDefendSwitch = expandWafDefendConfig(wafDefendConfig.(interface{})).ConfigSwitch
	}

	if customizeRuleConfig, ok := data.GetOk("customize_rule_config"); ok {
		request.CustomizeRuleSwitch = expandCustomizeRuleConfig(customizeRuleConfig.(interface{})).ConfigSwitch
	}

	if apiDefendConfig, ok := data.GetOk("api_defend_config"); ok {
		request.ApiDefendSwitch = expandApiDefendConfig(apiDefendConfig.(interface{})).ConfigSwitch
	}

	if whitelistConfig, ok := data.GetOk("whitelist_config"); ok {
		request.WhitelistSwitch = expandWhitelistConfig(whitelistConfig.(interface{})).ConfigSwitch
	}

	if blockConfig, ok := data.GetOk("block_config"); ok {
		request.BlockSwitch = expandBlockConfig(blockConfig.(interface{})).ConfigSwitch
	}

	if dmsDefendConfig, ok := data.GetOk("dms_defend_config"); ok {
		request.DmsDefendSwitch = expandDmsDefendConfig(dmsDefendConfig.(interface{})).ConfigSwitch
	}

	if intelligenceConfig, ok := data.GetOk("intelligence_config"); ok {
		request.IntelligenceSwitch = expandIntelligenceConfig(intelligenceConfig.(interface{})).ConfigSwitch
	}

	if botManageConfig, ok := data.GetOk("bot_manage_config"); ok {
		request.BotManageSwitch = expandBotManageConfig(botManageConfig.(interface{})).ConfigSwitch
	}

	if rateLimitConfig, ok := data.GetOk("rate_limit_config"); ok {
		request.RateLimitSwitch = expandRateLimitConfig(rateLimitConfig.(interface{})).ConfigSwitch
	}

	var response *waapDomain.ModifyPolicyStatusResponse
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseWaapDomainClient().UpdateDomainPolicy(request)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}
	if response == nil {
		return nil
	}
	log.Printf("resource.wangsu_waap_domain.update success")
	return nil
}

func resourceWaapDomainDelete(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_waap_domain.delete")

	var response *waapDomain.RemoveProtectedHostnameResponse
	var err error
	var diags diag.Diagnostics

	if targetDomains, ok := data.GetOk("target_domains"); ok {
		targetDomainsList := targetDomains.([]interface{})
		for _, v := range targetDomainsList {
			domain := v.(string)
			err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
				request := &waapDomain.RemoveProtectedHostnameParameters{
					Domain: &domain,
				}
				_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseWaapDomainClient().DeleteDomain(request)
				if err != nil {
					return resource.NonRetryableError(err)
				}
				return nil
			})
			if err != nil {
				diags = append(diags, diag.FromErr(err)...)
				return diags
			}
			if response == nil {
				return nil
			}
		}
	}
	return nil
}

func expandWafDefendConfig(i interface{}) *waapDomain.WafConfig {
	list := i.([]interface{})
	if len(list) == 0 || list[0] == nil {
		return nil
	}

	m := list[0].(map[string]interface{})
	config := &waapDomain.WafConfig{}
	if v, ok := m["config_switch"]; ok {
		config.ConfigSwitch = tea.String(v.(string))
	}
	if v, ok := m["rule_update_mode"]; ok {
		config.RuleUpdateMode = tea.String(v.(string))
	}
	if v, ok := m["defend_mode"]; ok {
		config.DefendMode = tea.String(v.(string))
	}
	return config
}

func expandCustomizeRuleConfig(i interface{}) *waapDomain.CustomizeRuleConfig {
	list := i.([]interface{})
	if len(list) == 0 || list[0] == nil {
		return nil
	}

	m := list[0].(map[string]interface{})
	config := &waapDomain.CustomizeRuleConfig{}
	if v, ok := m["config_switch"]; ok {
		config.ConfigSwitch = tea.String(v.(string))
	}
	return config
}

func expandApiDefendConfig(i interface{}) *waapDomain.APIDefendConfig {
	list := i.([]interface{})
	if len(list) == 0 || list[0] == nil {
		return nil
	}

	m := list[0].(map[string]interface{})
	config := &waapDomain.APIDefendConfig{}
	if v, ok := m["config_switch"]; ok {
		config.ConfigSwitch = tea.String(v.(string))
	}
	return config
}

func expandWhitelistConfig(i interface{}) *waapDomain.WhitelistConfig {
	list := i.([]interface{})
	if len(list) == 0 || list[0] == nil {
		return nil
	}

	m := list[0].(map[string]interface{})
	config := &waapDomain.WhitelistConfig{}
	if v, ok := m["config_switch"]; ok {
		config.ConfigSwitch = tea.String(v.(string))
	}
	return config
}

func expandBlockConfig(i interface{}) *waapDomain.BlockConfig {
	list := i.([]interface{})
	if len(list) == 0 || list[0] == nil {
		return nil
	}

	m := list[0].(map[string]interface{})
	config := &waapDomain.BlockConfig{}
	if v, ok := m["config_switch"]; ok {
		config.ConfigSwitch = tea.String(v.(string))
	}
	return config
}

func expandDmsDefendConfig(i interface{}) *waapDomain.DMSConfig {
	list := i.([]interface{})
	if len(list) == 0 || list[0] == nil {
		return nil
	}

	m := list[0].(map[string]interface{})
	config := &waapDomain.DMSConfig{}
	if v, ok := m["config_switch"]; ok {
		config.ConfigSwitch = tea.String(v.(string))
	}
	if v, ok := m["protection_mode"]; ok {
		config.ProtectionMode = tea.String(v.(string))
	}
	if v, ok := m["ai_switch"]; ok {
		config.AiSwitch = tea.String(v.(string))
	}
	return config
}

func expandIntelligenceConfig(i interface{}) *waapDomain.IntelligenceConfig {
	list := i.([]interface{})
	if len(list) == 0 || list[0] == nil {
		return nil
	}

	m := list[0].(map[string]interface{})
	config := &waapDomain.IntelligenceConfig{}
	if v, ok := m["config_switch"]; ok {
		config.ConfigSwitch = tea.String(v.(string))
	}
	if v, ok := m["info_cate_act"]; ok {
		list := v.([]interface{})
		if len(list) == 0 || list[0] == nil {
			return nil
		}
		art := list[0].(map[string]interface{})
		attackRiskType := &waapDomain.AttackRiskType{}
		if v, ok := art["attack_source"]; ok {
			attackRiskType.AttackSource = tea.String(v.(string))
		}
		if v, ok := art["spec_attack"]; ok {
			attackRiskType.SpecAttack = tea.String(v.(string))
		}
		if v, ok := art["industry"]; ok {
			attackRiskType.Industry = tea.String(v.(string))
		}
		config.InfoCateAct = attackRiskType
	}
	return config
}

func expandBotManageConfig(i interface{}) *waapDomain.BOTConfig {
	list := i.([]interface{})
	if len(list) == 0 || list[0] == nil {
		return nil
	}

	m := list[0].(map[string]interface{})
	config := &waapDomain.BOTConfig{}
	if v, ok := m["config_switch"]; ok {
		config.ConfigSwitch = tea.String(v.(string))
	}
	if v, ok := m["ua_bots_act"]; ok {
		config.UaBotsAct = tea.String(v.(string))
	}
	if v, ok := m["public_bots_act"]; ok {
		config.PublicBotsAct = tea.String(v.(string))
	}
	if v, ok := m["web_risk_config"]; ok {
		list := v.([]interface{})
		if len(list) == 0 || list[0] == nil {
			return nil
		}
		bwc := list[0].(map[string]interface{})
		botWebConfigDefault := &waapDomain.BotWebConfigDefault{}
		if v, ok := bwc["act"]; ok {
			botWebConfigDefault.Act = tea.String(v.(string))
		}
		config.WebRiskConfig = botWebConfigDefault
	}
	if v, ok := m["scene_analyse_switch"]; ok {
		config.SceneAnalyseSwitch = tea.String(v.(string))
	}
	return config
}

func expandRateLimitConfig(i interface{}) *waapDomain.RateLimitConfig {
	list := i.([]interface{})
	if len(list) == 0 || list[0] == nil {
		return nil
	}

	m := list[0].(map[string]interface{})
	config := &waapDomain.RateLimitConfig{}
	if v, ok := m["config_switch"]; ok {
		config.ConfigSwitch = tea.String(v.(string))
	}
	return config
}
