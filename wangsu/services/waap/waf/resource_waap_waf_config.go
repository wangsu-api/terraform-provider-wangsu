package waf

import (
	"context"
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wangsuCommon "github.com/wangsu-api/terraform-provider-wangsu/wangsu/common"
	securityPolicy "github.com/wangsu-api/wangsu-sdk-go/wangsu/securitypolicy"
	"log"
	"time"
)

func ResourceWaapWafConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: createWafConfig,
		ReadContext:   readWafConfig,
		UpdateContext: updateWafConfig,
		DeleteContext: deleteWafConfig,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Domain.",
			},
			"conf_basic": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Basic configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"defend_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Protection Mode.<br/>BLOCK: Block the attack request directly.<br/>LOG: Only log the attack request without blocking it.",
						},
						"rule_update_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Ruleset Mode.<br/>MANUAL: Check Ruleset update and all Recommendations on the Console, decide to apply them or not, all of these must be done by yourself manually.<br/>AUTO: Automatically upgrade the Ruleset to the latest version and apply the Recommendations learned from your website traffic to Exception, which can keep your website with high-level security anytime.",
						},
					},
				},
			},
			"rule_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Rule list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_id": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "WAF rule ID.",
						},
						"mode": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Rule actions.<br/>BLOCK: Deny request by a default 403 response.<br/>LOG: Log request and continue further detections.<br/>OFF: Select if you do not want a policy take effect.",
						},
					},
				},
			},
			"scan_protection": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Scan protection configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scan_tools_config": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Scanning tool detection configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Action.<br/>NO_USE:Not Used.<br/>LOG:Log.<br/>BLOCK:Deny.",
									},
								},
							},
						},
						"repeated_violation_config": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Repeated violation detection configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Action.<br/>NO_USE:Not Used.<br/>LOG:Log.<br/>BLOCK:Deny.",
									},
									"target": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Statistical subject.<br/>IP: IP.<br/>IP_JA3: IP and JA3 fingerprint.",
									},
									"period": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Time range, in seconds. Allowed values are from 5 to 1800.",
									},
									"waf_rule_type_count": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Number of WAF built-in rule triggers.must be greater than 1.",
									},
									"block_count": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Number of block actions.Allowed values are from 1 to 99999.",
									},
									"duration": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Action duration, in seconds.Allowed values are from 10 to 604800.",
									},
								},
							},
						},
						"directory_probing_config": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Directory probing detection configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Action.<br/>NO_USE:Not Used.<br/>LOG:Log.<br/>BLOCK:Deny.",
									},
									"target": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Statistical subject. <br/>IP: IP. <br/>IP_JA3: IP and JA3 fingerprint.",
									},
									"period": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Time range, in seconds. Allowed values are from 5 to 1800.",
									},
									"request_count_threshold": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Number of requests.Allowed values are from 1 to 99999.",
									},
									"non_existent_directory_threshold": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Number of non-existent directory requests.Allowed values are from 1 to 500.",
									},
									"rate404_threshold": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Proportion of 404 status codes.Allowed values are from 0 to 100.",
									},
									"duration": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Action duration, in seconds.Allowed values are from 10 to 604800.",
									},
								},
							},
						},
					},
				},
			},
		},
	}

}

func createWafConfig(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	err := errors.New("Please execute \"terraform import wangsu_waap_waf_config.<resource_name> <domain>\" before executing \"terraform apply\" for the first time.")
	diags = append(diags, diag.FromErr(err)...)
	return diags
}

func readWafConfig(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.waap_waf_config.read")
	var diags diag.Diagnostics
	// 使用导入的 ID 设置资源 ID
	if data.Id() != "" {
		_ = data.Set("domain", data.Id())
	} else if domain, ok := data.GetOk("domain"); ok {
		data.SetId(domain.(string))
	}
	diags = append(diags, readBasicConf(context, data, meta)...)
	diags = append(diags, readRule(context, data, meta)...)
	diags = append(diags, readScanProtectionConf(context, data, meta)...)
	return diags
}

func updateWafConfig(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.waap_waf_config.update")
	var diags diag.Diagnostics
	if data.HasChange("domain") {
		// 把domain强制刷回旧值，否则会有权限问题
		oldDomain, _ := data.GetChange("domain")
		_ = data.Set("domain", oldDomain)
		err := errors.New("Domain cannot be changed.")
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}
	if data.HasChange("conf_basic") {
		diags = append(diags, updateBasicConf(context, data, meta)...)
	}
	if data.HasChange("rule_list") {
		diags = append(diags, updateRule(context, data, meta)...)
	}
	if data.HasChange("scan_protection") {
		diags = append(diags, updateScanProtectionConf(context, data, meta)...)
	}
	return diags
}

func deleteWafConfig(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.waap_waf_config.delete - clearing local state only")
	// 清空本地 state
	data.SetId("")
	return nil
}

func readBasicConf(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.waap_waf_config.readBasicConf")
	var response *securityPolicy.ListWAFBasicConfigOfDomainsResponse
	request := &securityPolicy.ListWAFBasicConfigOfDomainsRequest{}
	var err error
	var diags diag.Diagnostics
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		domain := data.Id()
		request.SetDomainList([]*string{&domain})

		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseSecurityPolicyClient().GetWafConfig(request)
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
	if len(response.Data) == 0 {
		return nil
	}
	if response.Data != nil {
		basicConfig := make([]map[string]interface{}, 0)
		for _, item := range response.Data {
			configMap := make(map[string]interface{})
			if item.DefendMode != nil {
				configMap["defend_mode"] = *item.DefendMode
			}
			if item.RuleUpdateMode != nil {
				configMap["rule_update_mode"] = *item.RuleUpdateMode
			}
			basicConfig = append(basicConfig, configMap)
		}
		_ = data.Set("conf_basic", basicConfig)
	}
	return nil
}

func updateBasicConf(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.waap_waf_config.updateBasicConf")
	var diags diag.Diagnostics
	request := &securityPolicy.UpdateModeOfWAFRequest{}
	if domain, ok := data.GetOk("domain"); ok {
		domainStr := domain.(string)
		request.SetDomainList([]*string{&domainStr})
	}
	if confBasic, ok := data.GetOk("conf_basic"); ok {
		confBasicList := confBasic.([]interface{})
		if len(confBasicList) > 0 {
			confBasicMap := confBasicList[0].(map[string]interface{})
			if defendMode, ok := confBasicMap["defend_mode"]; ok {
				defendModeStr := defendMode.(string)
				request.DefendMode = &defendModeStr
			}
			if ruleUpdateMode, ok := confBasicMap["rule_update_mode"]; ok {
				ruleUpdateModeStr := ruleUpdateMode.(string)
				request.RuleUpdateMode = &ruleUpdateModeStr
			}
			var response *securityPolicy.UpdateModeOfWAFResponse
			var err error
			err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
				_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseSecurityPolicyClient().UpdateWafConfig(request)
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
			log.Printf("resource.waap_waf_config.updateBasicConf success")
		}
	}
	return nil
}

func readRule(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.waap_waf_config.readRule")
	var response *securityPolicy.ListWAFRulesResponse
	request := &securityPolicy.ListWAFRulesRequest{}
	var err error
	var diags diag.Diagnostics
	domain := data.Id()

	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		request.Domain = &domain

		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseSecurityPolicyClient().GetWafRule(request)
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
	if len(response.Data) == 0 {
		return nil
	}
	// 提取 main.tf 中的 rule_id 列表
	mainTfRuleIds := make(map[int]struct{})
	if ruleList, ok := data.GetOk("rule_list"); ok {
		ruleListInterface := ruleList.([]interface{})
		for _, ruleItem := range ruleListInterface {
			ruleMap := ruleItem.(map[string]interface{})
			if ruleId, ok := ruleMap["rule_id"].(int); ok {
				mainTfRuleIds[ruleId] = struct{}{}
			}
		}
	}
	// 过滤远程规则，只保留 main.tf 中的 rule_id
	filteredRuleList := make([]map[string]interface{}, 0)
	for _, item := range response.Data {
		if item != nil && item.RuleId != nil {
			ruleId := *item.RuleId
			if _, exists := mainTfRuleIds[ruleId]; len(mainTfRuleIds) == 0 || exists {
				ruleMap := make(map[string]interface{})
				ruleMap["rule_id"] = ruleId
				ruleMap["mode"] = *item.Mode
				filteredRuleList = append(filteredRuleList, ruleMap)
			}
		}
	}

	// 将过滤后的规则设置回状态
	_ = data.Set("rule_list", filteredRuleList)
	return nil
}

func updateRule(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.waap_waf_config.updateRule")
	var diags diag.Diagnostics
	request := &securityPolicy.UpdateActionForWAFManagedRulesRequest{}
	if domain, ok := data.GetOk("domain"); ok {
		domainStr := domain.(string)
		request.SetDomainList([]*string{&domainStr})
	}
	if ruleList, ok := data.GetOk("rule_list"); ok {
		ruleListInterface := ruleList.([]interface{})
		rules := make([]*securityPolicy.UpdateActionForWAFManagedRulesRequestRuleList, 0)
		for _, ruleItem := range ruleListInterface {
			ruleMap := ruleItem.(map[string]interface{})
			rule := &securityPolicy.UpdateActionForWAFManagedRulesRequestRuleList{}
			if ruleId, ok := ruleMap["rule_id"].(int); ok {
				rule.RuleId = &ruleId
			}
			if mode, ok := ruleMap["mode"]; ok {
				modeStr := mode.(string)
				rule.Mode = &modeStr
			}
			rules = append(rules, rule)
		}
		request.RuleList = rules
		var response *securityPolicy.UpdateActionForWAFManagedRulesResponse
		var err error
		err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
			_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseSecurityPolicyClient().UpdateWafRule(request)
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
		log.Printf("resource.waap_waf_config.updateRule success")
	}
	return nil
}

func readScanProtectionConf(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.waap_waf_config.readScanProtectionConf")
	var response *securityPolicy.GetWAFScanProtectionConfigResponse
	request := &securityPolicy.GetWAFScanProtectionConfigRequest{}
	var err error
	var diags diag.Diagnostics
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		domain := data.Id()
		request.SetDomainList([]*string{&domain})

		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseSecurityPolicyClient().GetWAFScanProtectionConfig(request)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	if response == nil || response.Data == nil {
		return nil
	}
	if len(response.Data.Array) == 0 {
		return nil
	}
	if response.Data != nil {
		scanProtection := make([]map[string]interface{}, 0)
		for _, item := range response.Data.Array {
			configMap := make(map[string]interface{})
			if item.ScanToolsConfig != nil {
				scanToolsConfigMap := make(map[string]interface{})
				if item.ScanToolsConfig.Action != nil {
					scanToolsConfigMap["action"] = *item.ScanToolsConfig.Action
				}
				configMap["scan_tools_config"] = []map[string]interface{}{scanToolsConfigMap}
			}
			if item.RepeatedViolationConfig != nil {
				repeatedViolationConfigMap := make(map[string]interface{})
				if item.RepeatedViolationConfig.Action != nil {
					repeatedViolationConfigMap["action"] = *item.RepeatedViolationConfig.Action
				}
				if item.RepeatedViolationConfig.Target != nil {
					repeatedViolationConfigMap["target"] = *item.RepeatedViolationConfig.Target
				}
				if item.RepeatedViolationConfig.Period != nil {
					repeatedViolationConfigMap["period"] = *item.RepeatedViolationConfig.Period
				}
				if item.RepeatedViolationConfig.WafRuleTypeCount != nil {
					repeatedViolationConfigMap["waf_rule_type_count"] = *item.RepeatedViolationConfig.WafRuleTypeCount
				}
				if item.RepeatedViolationConfig.BlockCount != nil {
					repeatedViolationConfigMap["block_count"] = *item.RepeatedViolationConfig.BlockCount
				}
				if item.RepeatedViolationConfig.Duration != nil {
					repeatedViolationConfigMap["duration"] = *item.RepeatedViolationConfig.Duration
				}
				configMap["repeated_violation_config"] = []map[string]interface{}{repeatedViolationConfigMap}
			}
			if item.DirectoryProbingConfig != nil {
				directoryProbingConfigMap := make(map[string]interface{})
				if item.DirectoryProbingConfig.Action != nil {
					directoryProbingConfigMap["action"] = *item.DirectoryProbingConfig.Action
				}
				if item.DirectoryProbingConfig.Target != nil {
					directoryProbingConfigMap["target"] = *item.DirectoryProbingConfig.Target
				}
				if item.DirectoryProbingConfig.Period != nil {
					directoryProbingConfigMap["period"] = *item.DirectoryProbingConfig.Period
				}
				if item.DirectoryProbingConfig.RequestCountThreshold != nil {
					directoryProbingConfigMap["request_count_threshold"] = *item.DirectoryProbingConfig.RequestCountThreshold
				}
				if item.DirectoryProbingConfig.NonExistentDirectoryThreshold != nil {
					directoryProbingConfigMap["non_existent_directory_threshold"] = *item.DirectoryProbingConfig.NonExistentDirectoryThreshold
				}
				if item.DirectoryProbingConfig.Rate404Threshold != nil {
					directoryProbingConfigMap["rate404_threshold"] = *item.DirectoryProbingConfig.Rate404Threshold
				}
				if item.DirectoryProbingConfig.Duration != nil {
					directoryProbingConfigMap["duration"] = *item.DirectoryProbingConfig.Duration
				}
				configMap["directory_probing_config"] = []map[string]interface{}{directoryProbingConfigMap}
			}
			scanProtection = append(scanProtection, configMap)
		}
		_ = data.Set("scan_protection", scanProtection)
	}
	return nil
}

func updateScanProtectionConf(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.waap_waf_config.updateScanProtectionConf")
	var diags diag.Diagnostics
	request := &securityPolicy.UpdateWAFScanProtectionConfigRequest{}
	if domain, ok := data.GetOk("domain"); ok {
		domainStr := domain.(string)
		request.SetDomainList([]*string{&domainStr})
	}
	if scanProtection, ok := data.GetOk("scan_protection"); ok {
		scanProtectionList := scanProtection.([]interface{})
		if len(scanProtectionList) > 0 {
			scanProtectionMap := scanProtectionList[0].(map[string]interface{})
			if scanToolsConfig, ok := scanProtectionMap["scan_tools_config"]; ok {
				scanToolsConfigList := scanToolsConfig.([]interface{})
				if len(scanToolsConfigList) > 0 {
					scanToolsConfigMap := scanToolsConfigList[0].(map[string]interface{})
					scanToolsConfigObj := &securityPolicy.UpdateWAFScanProtectionConfigRequestScanToolsConfig{}
					if action, ok := scanToolsConfigMap["action"]; ok {
						actionStr := action.(string)
						scanToolsConfigObj.Action = &actionStr
					}
					request.ScanToolsConfig = scanToolsConfigObj
				}
			}
			if repeatedViolationConfig, ok := scanProtectionMap["repeated_violation_config"]; ok {
				repeatedViolationConfigList := repeatedViolationConfig.([]interface{})
				if len(repeatedViolationConfigList) > 0 {
					repeatedViolationConfigMap := repeatedViolationConfigList[0].(map[string]interface{})
					repeatedViolationConfigObj := &securityPolicy.UpdateWAFScanProtectionConfigRequestRepeatedViolationConfig{}
					if action, ok := repeatedViolationConfigMap["action"]; ok {
						actionStr := action.(string)
						repeatedViolationConfigObj.Action = &actionStr
					}
					if target, ok := repeatedViolationConfigMap["target"]; ok {
						targetStr := target.(string)
						repeatedViolationConfigObj.Target = &targetStr
					}
					if period, ok := repeatedViolationConfigMap["period"]; ok {
						periodInt := period.(int)
						repeatedViolationConfigObj.Period = &periodInt
					}
					if wafRuleTypeCount, ok := repeatedViolationConfigMap["waf_rule_type_count"]; ok {
						wafRuleTypeCountInt := wafRuleTypeCount.(int)
						repeatedViolationConfigObj.WafRuleTypeCount = &wafRuleTypeCountInt
					}
					if blockCount, ok := repeatedViolationConfigMap["block_count"]; ok {
						blockCountInt := blockCount.(int)
						repeatedViolationConfigObj.BlockCount = &blockCountInt
					}
					if duration, ok := repeatedViolationConfigMap["duration"]; ok {
						durationInt := duration.(int)
						repeatedViolationConfigObj.Duration = &durationInt
					}
					request.RepeatedViolationConfig = repeatedViolationConfigObj
				}
			}
			if directoryProbingConfig, ok := scanProtectionMap["directory_probing_config"]; ok {
				directoryProbingConfigList := directoryProbingConfig.([]interface{})
				if len(directoryProbingConfigList) > 0 {
					directoryProbingConfigMap := directoryProbingConfigList[0].(map[string]interface{})
					directoryProbingConfigObj := &securityPolicy.UpdateWAFScanProtectionConfigRequestDirectoryProbingConfig{}
					if action, ok := directoryProbingConfigMap["action"]; ok {
						actionStr := action.(string)
						directoryProbingConfigObj.Action = &actionStr
					}
					if target, ok := directoryProbingConfigMap["target"]; ok {
						targetStr := target.(string)
						directoryProbingConfigObj.Target = &targetStr
					}
					if period, ok := directoryProbingConfigMap["period"]; ok {
						periodInt := period.(int)
						directoryProbingConfigObj.Period = &periodInt
					}
					if requestCountThreshold, ok := directoryProbingConfigMap["request_count_threshold"]; ok {
						requestCountThresholdInt := requestCountThreshold.(int)
						directoryProbingConfigObj.RequestCountThreshold = &requestCountThresholdInt
					}
					if nonExistentDirectoryThreshold, ok := directoryProbingConfigMap["non_existent_directory_threshold"]; ok {
						nonExistentDirectoryThresholdInt := nonExistentDirectoryThreshold.(int)
						directoryProbingConfigObj.NonExistentDirectoryThreshold = &nonExistentDirectoryThresholdInt
					}
					if rate404Threshold, ok := directoryProbingConfigMap["rate404_threshold"]; ok {
						rate404ThresholdInt := rate404Threshold.(int)
						directoryProbingConfigObj.Rate404Threshold = &rate404ThresholdInt
					}
					if duration, ok := directoryProbingConfigMap["duration"]; ok {
						durationInt := duration.(int)
						directoryProbingConfigObj.Duration = &durationInt
					}
					request.DirectoryProbingConfig = directoryProbingConfigObj
				}
			}

			var response *securityPolicy.UpdateWAFScanProtectionConfigResponse
			var err error
			err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
				_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseSecurityPolicyClient().UpdateWAFScanProtectionConfig(request)
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
			log.Printf("resource.waap_waf_config.updateScanProtectionConf success")
		}
	}
	return nil
}
