package waf

import (
	"context"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wangsuCommon "github.com/wangsu-api/terraform-provider-wangsu/wangsu/common"
	waapWAF "github.com/wangsu-api/wangsu-sdk-go/wangsu/waap/waf"
	"log"
	"time"
)

func DataSourceWaapWAF() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWaapWAFRead,

		Schema: map[string]*schema.Schema{
			"domain_list": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Hostname list.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"data": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Data.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"array": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Array.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domain": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Domain.",
									},
									"config_switch": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Policy switch.ON: Enable.OFF: Disable.",
									},
									"conf_basic": {
										Type:        schema.TypeList,
										Required:    true,
										MaxItems:    1,
										Description: "Basic configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"defend_mode": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Protection Mode.<br/>BLOCK: Block the attack request directly.<br/>LOG: Only log the attack request without blocking it.",
												},
												"rule_update_mode": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Ruleset Mode.<br/>MANUAL:  Check Ruleset update and all Recommendations on the Console, decide to apply them or not, all of these must be done by yourself manually.<br/>AUTO: Automatically upgrade the Ruleset to the latest version and apply the Recommendations learned from your website traffic to Exception, which can keep your website with high-level security anytime.",
												},
											},
										},
									},
									"rule_list": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "Rule list.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"rule_id": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "WAF rule ID.",
												},
												"name": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Rule name.",
												},
												"mode": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Rule actions.<br/>BLOCK: Deny request by a default 403 response.<br/>LOG: Log request and continue further detections.<br/>OFF: Select if you do not a policy take effect.",
												},
												"exception_list": {
													Type:        schema.TypeList,
													Required:    true,
													Description: "Rule exceptions.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Matching conditions.<br/>ip: IP<br/>path: Path<br/>uri: URI<br/>urlParamName: URI Parameter Name<br/>urlParamValue: URI Parameter Value<br/>userAgent: User Agent<br/>httpHeaderName: Request Header Name<br/>httpHeaderValue: Request Header Value<br/>cookie: Cookie<br/>body: Body<br/>bodyParamName: Body Parameter Name<br/>bodyParamValue: Body Parameter Value",
															},
															"match_type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Match type,IP can only be EQUAL.<br/>EQUAL: Equal<br/>CONTAIN: Contains<br/>REGEX: Regular match",
															},
															"content_list": {
																Type:        schema.TypeList,
																Required:    true,
																Description: "Rule exceptions.<br/>When matchType=EQUAL, case-sensitive, path and uri must start with \"/\", and body can only pass one value;<br/>When matchType=REGEX, only one value can be passed.",
																Elem: &schema.Schema{
																	Type:        schema.TypeString,
																	Description: "Rule exceptions.<br/>When matchType=EQUAL, case-sensitive, path and uri must start with \"/\", and body can only pass one value;<br/>When matchType=REGEX, only one value can be passed.",
																},
															},
														},
													},
												},
											},
										},
									},
									"scan_protection": {
										Type:        schema.TypeList,
										Required:    true,
										MaxItems:    1,
										Description: "Scan protection configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"scan_tools_config": {
													Type:     schema.TypeList,
													Required: true,
													MaxItems: 1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"action": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Action.<br/>NO_USE:Not Used.<br/>LOG:Log.<br/>BLOCK:Deny.",
															},
														},
													},
													Description: "Scanning tool detection configuration.",
												},
												"repeated_violation_config": {
													Type:     schema.TypeList,
													Required: true,
													MaxItems: 1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"action": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Action.<br/>NO_USE:Not Used.<br/>LOG:Log.<br/>BLOCK:Deny.",
															},
															"target": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Statistical subject.<br/>IP: IP.<br/>IP_JA3: IP and JA3 fingerprint.",
															},
															"period": {
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Time range, in seconds.",
															},
															"waf_rule_type_count": {
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Number of WAF built-in rule triggers.",
															},
															"block_count": {
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Number of block actions.",
															},
															"duration": {
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Handling action duration, in seconds.",
															},
														},
													},
													Description: "Repeated violation detection configuration.",
												},
												"directory_probing_config": {
													Type:     schema.TypeList,
													Required: true,
													MaxItems: 1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"action": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Action.<br/>NO_USE:Not Used.<br/>LOG:Log.<br/>BLOCK:Deny.",
															},
															"target": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Statistical subject.<br/>IP: IP.<br/>IP_JA3: IP and JA3 fingerprint.",
															},
															"period": {
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Time range, in seconds.",
															},
															"request_count_threshold": {
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Number of requests.",
															},
															"non_existent_directory_threshold": {
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Number of non-existent directory requests.",
															},
															"rate404_threshold": {
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Proportion of 404 status codes.",
															},
															"duration": {
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Handling action duration, in seconds.",
															},
														},
													},
													Description: "Directory probing detection configuration.",
												},
											},
										},
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

func dataSourceWaapWAFRead(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("data_source.wangsu_waap_waf.read")

	var diags diag.Diagnostics
	request := &waapWAF.GetWafConfigurationRequest{}

	// Retrieve domain_list from input
	if v, ok := data.GetOk("domain_list"); ok {
		domainList := v.([]interface{})
		domains := make([]*string, len(domainList))
		for i, domain := range domainList {
			domainStr := domain.(string)
			domains[i] = &domainStr
		}
		request.SetDomainList(domains)
	}

	// Make API call
	var response *waapWAF.GetWafConfigurationResponse
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseWaapWAFClient().GetWafConfiguration(request)
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
		data.SetId("")
		return nil
	}

	// Parse response data
	wafDataList := make([]interface{}, len(response.Data.Array))
	ids := make([]string, len(response.Data.Array))
	for i, item := range response.Data.Array {
		parsedItem := map[string]interface{}{
			"domain":        tea.StringValue(item.Domain),
			"config_switch": tea.StringValue(item.ConfigSwitch),
			"conf_basic": []interface{}{
				map[string]interface{}{
					"defend_mode":      tea.StringValue(item.ConfBasic.DefendMode),
					"rule_update_mode": tea.StringValue(item.ConfBasic.RuleUpdateMode),
				},
			},
			"scan_protection": []interface{}{
				map[string]interface{}{
					"scan_tools_config": []interface{}{
						map[string]interface{}{
							"action": tea.StringValue(item.ScanProtection.ScanToolsConfig.Action),
						},
					},
					"repeated_violation_config": []interface{}{
						map[string]interface{}{
							"action":              tea.StringValue(item.ScanProtection.RepeatedViolationConfig.Action),
							"target":              tea.StringValue(item.ScanProtection.RepeatedViolationConfig.Target),
							"period":              tea.IntValue(item.ScanProtection.RepeatedViolationConfig.Period),
							"waf_rule_type_count": tea.IntValue(item.ScanProtection.RepeatedViolationConfig.WafRuleTypeCount),
							"block_count":         tea.IntValue(item.ScanProtection.RepeatedViolationConfig.BlockCount),
							"duration":            tea.IntValue(item.ScanProtection.RepeatedViolationConfig.Duration),
						},
					},
					"directory_probing_config": []interface{}{
						map[string]interface{}{
							"action":                           tea.StringValue(item.ScanProtection.DirectoryProbingConfig.Action),
							"target":                           tea.StringValue(item.ScanProtection.DirectoryProbingConfig.Target),
							"period":                           tea.IntValue(item.ScanProtection.DirectoryProbingConfig.Period),
							"request_count_threshold":          tea.IntValue(item.ScanProtection.DirectoryProbingConfig.RequestCountThreshold),
							"non_existent_directory_threshold": tea.IntValue(item.ScanProtection.DirectoryProbingConfig.NonExistentDirectoryThreshold),
							"rate404_threshold":                tea.IntValue(item.ScanProtection.DirectoryProbingConfig.Rate404Threshold),
							"duration":                         tea.IntValue(item.ScanProtection.DirectoryProbingConfig.Duration),
						},
					},
				},
			},
		}

		// Parse rule_list
		if item.RuleList != nil {
			ruleList := make([]interface{}, len(item.RuleList))
			for j, rule := range item.RuleList {
				ruleMap := map[string]interface{}{
					"rule_id": tea.IntValue(rule.RuleId),
					"name":    tea.StringValue(rule.Name),
					"mode":    tea.StringValue(rule.Mode),
				}

				// Parse exception_list
				if rule.ExceptionList != nil {
					exceptionList := make([]interface{}, len(rule.ExceptionList))
					for k, exception := range rule.ExceptionList {
						exceptionList[k] = map[string]interface{}{
							"type":         tea.StringValue(exception.Type),
							"match_type":   tea.StringValue(exception.MatchType),
							"content_list": tea.StringSliceValue(exception.ContentList),
						}
					}
					ruleMap["exception_list"] = exceptionList
				}
				ruleList[j] = ruleMap
			}
			parsedItem["rule_list"] = ruleList
		}

		wafDataList[i] = parsedItem
		ids[i] = tea.StringValue(item.Domain)
	}

	// Set data and ID
	if err := data.Set("data", []interface{}{
		map[string]interface{}{
			"array": wafDataList,
		},
	}); err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}
	data.SetId(wangsuCommon.DataResourceIdsHash(ids))
	return diags
}
