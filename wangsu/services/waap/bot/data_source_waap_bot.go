package bot

import (
	"context"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wangsuCommon "github.com/wangsu-api/terraform-provider-wangsu/wangsu/common"
	waapBot "github.com/wangsu-api/wangsu-sdk-go/wangsu/waap/bot"
	"log"
	"time"
)

func DataSourceWaapBot() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWaapBotRead,

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
						"config_list": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Configuration list.",
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
										Description: "Policy switch.<br/>ON: Enable.<br/>OFF: Disable.",
									},
									"general_strategy": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "General policies.",
										MaxItems:    1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ai_bots": {
													Type:        schema.TypeList,
													Required:    true,
													Description: "AI Bots configuration.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"bot_category": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "AI Bots categories.",
															},
															"action": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Actions.<br/>NO_USE: Not Used.<br/>LOG: Log.<br/>BLOCK: Deny.",
															},
														},
													},
												},
												"public_bots": {
													Type:        schema.TypeList,
													Required:    true,
													Description: "Public Bots configuration.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"bot_category": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Public Bots categories.",
															},
															"action": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Actions.<br/>NO_USE: Not Used.<br/>LOG: Log.<br/>BLOCK: Deny.<br/>ACCEPT: Release.",
															},
														},
													},
												},
												"absolute_bots_act": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Definite Bots actions.<br/>NO_USE: Not Used.<br/>LOG: Log.<br/>BLOCK: Deny.",
												},
												"bot_tagging": {
													Type:        schema.TypeList,
													Required:    true,
													Description: "Bot Tagging.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"request_header_key": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Customizing HTTP request headers.",
															},
															"traffic_characteristics": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Tagging Bot traffic characteristics.<br/>AI_BOT: AI Bots.<br/>PUBLIC_BOT: Public Bots.<br/>CUSTOMIZE_BOT: Custom Bots.",
															},
														},
													},
												},
											},
										},
									},
									"web_config": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "Web risk detection.",
										MaxItems:    1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"act": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Action.<br/>NO_USE: Not used.<br/>LOG: Log.<br/>BLOCK: Deny.",
												},
												"browser_analyse_switch": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Browser feature switch.<br/>ON: Enable<br/>OFF: Disable",
												},
												"auto_tool_switch": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Automated tool detection function switch.<br/>ON: Enable<br/>OFF: Disable",
												},
												"crack_analyse_switch": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Cracking the behavior detection switch.<br/>ON: Enable<br/>OFF: Disable",
												},
												"page_debug_switch": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Page anti-debugging switch.<br/>ON: Enable<br/>OFF: Disable",
												},
												"interaction_analyse_switch": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Interactive behavior verification switch.<br/>ON: Enable<br/>OFF: Disable",
												},
												"js_exception_list": {
													Type:        schema.TypeList,
													Required:    true,
													Description: "Rule list of html pages without embedding JS.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Rule ID of html pages without embedding JS.",
															},
															"path": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Path.",
															},
															"description": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Description.",
															},
															"regex": {
																Type:        schema.TypeBool,
																Required:    true,
																Description: "Whether it is REGEX.<br/>false: No<br/>true: Yes",
															},
														},
													},
												},
												"interaction_rule_list": {
													Type:        schema.TypeList,
													Required:    true,
													Description: "Interactive behavior validation rule list.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Interactive behavior validation rule ID.",
															},
															"path": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Path.",
															},
															"trigger_times": {
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Minimum number of triggers.",
															},
															"regex": {
																Type:        schema.TypeBool,
																Required:    true,
																Description: "Whether it is REGEX.<br/>false: No<br/>true: Yes",
															},
														},
													},
												},
												"ajax_exception_switch": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Ajax request exception function switch.<br/>ON: Enable<br/>OFF: Disable",
												},
											},
										},
									},
									"scene_whitelist": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "Application request whitelist.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "ID.",
												},
												"name": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Whitelist name.",
												},
												"description": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Description.",
												},
												"conditions": {
													Type:        schema.TypeList,
													Required:    true,
													Description: "List of matching conditions.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"match_name": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Matching condition name.<br/>IP_IPS: IP/CIDR<br/>PATH: Path<br/>URI: Path with parameters<br/>HEADER: Request Header<br/>UA: User Agent<br/>REQUEST_METHOD: Request Method<br/>REFERER: Referer",
															},
															"match_type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "When matchName is IP_IPS, maximum 300 IP/CIDR in match value list, the optional value of matchType is:<br/>EQUAL: Equals<br/>NOT_EQUAL: Does not equal<br/>When matchName is a URI, the optional value of matchType is:<br/>EQUAL: Equals, the matching value is case-sensitive and needs to start with \"/\" and include parameters<br/>NOT_EQUAL: Does not equal, the matching value is case-sensitive, needs to start with \"/\", and contains parameters<br/>CONTAIN: Contains, match values are case insensitive<br/>NOT_CONTAIN: Does not contains, match values are case insensitive<br/>REGEX: Regex match, only one value is allowed for the match value<br/>NOT_REGEX: regular does not match<br/>START_WITH: starts with<br/>END_WITH: ends with<br/>WILDCARD: wildcard matches<br/>NOT_WILDCARD: wildcard does not match<br/>When matchName is PATH, the optional value of matchType is:<br/>EQUAL: Equals, the matching value is case-sensitive and needs to start with \"/\" , and does not contain parameters<br/>NOT_EQUAL: Does not equal, the matching value is case-sensitive, needs to start with \"/\", and does not contain parameters<br/>CONTAIN: Contains, match values are case insensitive<br/>NOT_CONTAIN: Does not contains, match values are case insensitive<br/>REGEX: Regex match, match values are case insensitive and only one value is allowed<br/>NOT_REGEX: regular does not match<br/>START_WITH: starts with<br/>END_WITH: ends with<br/>WILDCARD: wildcard matches<br/>NOT_WILDCARD: wildcard does not match<br/>When matchName is HEADER, the optional value of matchType is:<br/>EQUAL: Equals, match values are case sensitive<br/>NOT_EQUAL: Does not equal, the matching value is case-sensitive<br/>CONTAIN: Contains, match values are case insensitive<br/>NOT_CONTAIN: Does not contains, match values are case insensitive<br/>REGEX: Regex match, match values are case insensitive and only one value is allowed<br/>NONE: Empty or does not exist<br/>NOT_REGEX: regular does not match<br/>START_WITH: starts with<br/>END_WITH: ends with<br/>WILDCARD: wildcard matches<br/>NOT_WILDCARD: wildcard does not match<br/>When matchName is UA, the optional value of matchType is:<br/>EQUAL: Equals, match values are case sensitive<br/>NOT_EQUAL: Does not equal, the matching value is case-sensitive<br/>CONTAIN: Contains, match values are case insensitive<br/>NOT_CONTAIN: Does not contains, match values are case insensitive<br/>REGEX: Regex match, match values are case insensitive and only one value is allowed<br/>NONE: Empty or does not exist<br/>NOT_REGEX: regular does not match<br/>START_WITH: starts with<br/>END_WITH: ends with<br/>WILDCARD: wildcard matches<br/>NOT_WILDCARD: wildcard does not match<br/>When matchName is REFERER, the optional value of matchType is:<br/>EQUAL: Equals, match values are case sensitive<br/>NOT_EQUAL: Does not equal, the matching value is case-sensitive<br/>CONTAIN: Contains, match values are case insensitive<br/>NOT_CONTAIN: Does not contains, match values are case insensitive<br/>REGEX: Regex match, match values are case insensitive and only one value is allowed<br/>NONE: Empty or does not exist<br/>NOT_REGEX: regular does not match<br/>START_WITH: starts with<br/>END_WITH: ends with<br/>WILDCARD: wildcard matches<br/>NOT_WILDCARD: wildcard does not match<br/>When matchName is REQUEST_METHOD, the optional value of matchType is:<br/>EQUAL: Equals, match values are case sensitive<br/>NOT_EQUAL: Does not equal, the matching value is case-sensitive<br/>",
															},
															"match_key": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Mathing key, this value is not empty and valid only when matchName=HEADER.<br/>Maximum 100 characters, case insensitive.",
															},
															"match_value_list": {
																Type:        schema.TypeList,
																Required:    true,
																Description: "List of matching values.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
														},
													},
												},
											},
										},
									},
									"traffic_detection": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "Abnormal traffic detection.",
										MaxItems:    1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"start_time": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Start time, format: HH:mm.",
												},
												"end_time": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "End time, format: HH:mm.",
												},
												"action": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Action.<br/>RESET: Reset connection.<br/>LOG: Log.<br/>BLOCK: Block.<br/>NO_USE: Not Used.",
												},
												"whitelist": {
													Type:        schema.TypeList,
													Required:    true,
													Description: "Whitelist.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
									"custom_bots": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "Custom bots.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Rule ID.",
												},
												"bot_name": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Rule Name.",
												},
												"bot_description": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Description.",
												},
												"bot_act": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Actions:<br/>BLOCK: block<br/>LOG: log<br/>ACCEPT: release",
												},
												"condition_list": {
													Type:        schema.TypeList,
													Required:    true,
													Description: "Matching conditions.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"condition_name": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Matching condition name. <br/>IP_IPS: IP/IP segment <br/>JA3: JA3 Fingerprint<br/>JA4: JA4 Fingerprint<br/>UA: User-agent <br/>HEADER: Request Header <br/>ASN: AS Number <br/>CLIENT_GROUP: Client Group <br/>PUBLIC_BOT: Public Bots",
															},
															"condition_value_list": {
																Type:        schema.TypeList,
																Required:    true,
																Description: "Condition value list.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"condition_func": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Matching condition function.<br/>EQUAL: Equals<br/>NOT_EQUAL: Does not equal<br/>CONTAIN: Contains<br/>NOT_CONTAIN: Does not Contains<br/>NONE:Empty or non-existent<br/>REGEX: Regex match<br/>NOT_REGEX: Regular does not match<br/>START_WITH: Starts with<br/>END_WITH: Ends with<br/>WILDCARD: Wildcard matches, * represents zero or more arbitrary characters, ? represents any single character<br/>NOT_WILDCARD: Wildcard does not match, * represents zero or more arbitrary characters, ? represents any single character<br/><br/>When conditionName is IP_IPS, the value can be:EQUAL,NOT_EQUAL <br/>When conditionName is JA3, the value can be:EQUAL,NOT_EQUAL<br/>When conditionName is JA4, the value can be:EQUAL,NOT_EQUAL,CONTAIN,NOT_CONTAIN,START_WITH,END_WITH,WILDCARD,NOT_WILDCARD<br/>When conditionName is UA, the value can be:EQUAL,NOT_EQUAL,CONTAIN,NOT_CONTAIN,NONE,REGEX,NOT_REGEX,START_WITH,END_WITH,WILDCARD,NOT_WILDCARD <br/>When conditionName is HEADER, the value can be:EQUAL,NOT_EQUAL,CONTAIN,NOT_CONTAIN,NONE,REGEX,NOT_REGEX,START_WITH,END_WITH,WILDCARD,NOT_WILDCARD <br/>When conditionName is ASN, the value can be:EQUAL,NOT_EQUAL <br/>When conditionName is CLIENT_GROUP, the value can be:EQUAL,NOT_EQUAL <br/>When conditionName is PUBLIC_BOT, the value can be:EQUAL,NOT_EQUAL",
															},
															"condition_key": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Request header name.case insensitive.",
															},
														},
													},
												},
												"creator": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Creator.",
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

func dataSourceWaapBotRead(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("data_source.wangsu_waap_bot.read")

	var diags diag.Diagnostics
	request := &waapBot.GetBotManagementConfigRequest{}

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
	var response *waapBot.GetBotManagementConfigResponse
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseWaapBotClient().GetBotManagementConfig(request)
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
	var configListData = parseResponseData(response.Data.ConfigList)
	if err := data.Set("data", []interface{}{
		map[string]interface{}{
			"config_list": configListData,
		},
	}); err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	// Generate a unique ID based on domains
	var ids []string
	for _, config := range response.Data.ConfigList {
		ids = append(ids, *config.Domain)
	}
	data.SetId(wangsuCommon.DataResourceIdsHash(ids))

	return diags
}

func parseResponseData(configList []*waapBot.GetBotManagementConfigResponseDataConfigList) []interface{} {
	var configListData []interface{}

	for _, config := range configList {
		configMap := map[string]interface{}{
			"domain":        tea.StringValue(config.Domain),
			"config_switch": tea.StringValue(config.ConfigSwitch),
		}

		// Parse general strategy
		if config.GeneralStrategy != nil {
			generalStrategy := map[string]interface{}{
				"absolute_bots_act": tea.StringValue(config.GeneralStrategy.AbsoluteBotsAct),
			}

			// AI Bots
			if config.GeneralStrategy.AiBots != nil {
				aiBots := make([]interface{}, len(config.GeneralStrategy.AiBots))
				for i, aiBot := range config.GeneralStrategy.AiBots {
					aiBots[i] = map[string]interface{}{
						"bot_category": tea.StringValue(aiBot.BotCategory),
						"action":       tea.StringValue(aiBot.Action),
					}
				}
				generalStrategy["ai_bots"] = aiBots
			}

			// Public Bots
			if config.GeneralStrategy.PublicBots != nil {
				publicBots := make([]interface{}, len(config.GeneralStrategy.PublicBots))
				for i, publicBot := range config.GeneralStrategy.PublicBots {
					publicBots[i] = map[string]interface{}{
						"bot_category": tea.StringValue(publicBot.BotCategory),
						"action":       tea.StringValue(publicBot.Action),
					}
				}
				generalStrategy["public_bots"] = publicBots
			}

			// Bot Tagging
			if config.GeneralStrategy.BotTagging != nil {
				botTagging := make([]interface{}, len(config.GeneralStrategy.BotTagging))
				for i, tag := range config.GeneralStrategy.BotTagging {
					botTagging[i] = map[string]interface{}{
						"request_header_key":      tea.StringValue(tag.RequestHeaderKey),
						"traffic_characteristics": tea.StringValue(tag.TrafficCharacteristics),
					}
				}
				generalStrategy["bot_tagging"] = botTagging
			}

			configMap["general_strategy"] = []interface{}{generalStrategy}
		}

		// Parse web config
		if config.WebConfig != nil {
			webConfig := map[string]interface{}{
				"act":                        tea.StringValue(config.WebConfig.Act),
				"browser_analyse_switch":     tea.StringValue(config.WebConfig.BrowserAnalyseSwitch),
				"auto_tool_switch":           tea.StringValue(config.WebConfig.AutoToolSwitch),
				"crack_analyse_switch":       tea.StringValue(config.WebConfig.CrackAnalyseSwitch),
				"page_debug_switch":          tea.StringValue(config.WebConfig.PageDebugSwitch),
				"interaction_analyse_switch": tea.StringValue(config.WebConfig.InteractionAnalyseSwitch),
				"ajax_exception_switch":      tea.StringValue(config.WebConfig.AjaxExceptionSwitch),
			}

			// JS Exception List
			if config.WebConfig.JsExceptionList != nil {
				jsExceptions := make([]interface{}, len(config.WebConfig.JsExceptionList))
				for i, exception := range config.WebConfig.JsExceptionList {
					jsExceptions[i] = map[string]interface{}{
						"id":          tea.StringValue(exception.Id),
						"path":        tea.StringValue(exception.Path),
						"description": tea.StringValue(exception.Description),
						"regex":       tea.BoolValue(exception.Regex),
					}
				}
				webConfig["js_exception_list"] = jsExceptions
			}

			// Interaction Rule List
			if config.WebConfig.InteractionRuleList != nil {
				interactionRules := make([]interface{}, len(config.WebConfig.InteractionRuleList))
				for i, rule := range config.WebConfig.InteractionRuleList {
					interactionRules[i] = map[string]interface{}{
						"id":            tea.StringValue(rule.Id),
						"path":          tea.StringValue(rule.Path),
						"trigger_times": tea.IntValue(rule.TriggerTimes),
						"regex":         tea.BoolValue(rule.Regex),
					}
				}
				webConfig["interaction_rule_list"] = interactionRules
			}

			configMap["web_config"] = []interface{}{webConfig}
		}

		// Parse scene whitelist
		if config.SceneWhitelist != nil {
			sceneWhitelist := make([]interface{}, len(config.SceneWhitelist))
			for i, whitelist := range config.SceneWhitelist {
				whitelistMap := map[string]interface{}{
					"id":          tea.StringValue(whitelist.Id),
					"name":        tea.StringValue(whitelist.Name),
					"description": tea.StringValue(whitelist.Description),
				}

				// Conditions
				if whitelist.Conditions != nil {
					conditions := make([]interface{}, len(whitelist.Conditions))
					for j, condition := range whitelist.Conditions {
						conditions[j] = map[string]interface{}{
							"match_name":       tea.StringValue(condition.MatchName),
							"match_type":       tea.StringValue(condition.MatchType),
							"match_key":        tea.StringValue(condition.MatchKey),
							"match_value_list": tea.StringSliceValue(condition.MatchValueList),
						}
					}
					whitelistMap["conditions"] = conditions
				}

				sceneWhitelist[i] = whitelistMap
			}
			configMap["scene_whitelist"] = sceneWhitelist
		}

		// Parse traffic detection
		if config.TrafficDetection != nil {
			trafficDetection := map[string]interface{}{
				"start_time": tea.StringValue(config.TrafficDetection.StartTime),
				"end_time":   tea.StringValue(config.TrafficDetection.EndTime),
				"action":     tea.StringValue(config.TrafficDetection.Action),
				"whitelist":  tea.StringSliceValue(config.TrafficDetection.Whitelist),
			}
			configMap["traffic_detection"] = []interface{}{trafficDetection}
		}

		// Parse custom bots
		if config.CustomBots != nil {
			customBots := make([]interface{}, len(config.CustomBots))
			for i, bot := range config.CustomBots {
				botMap := map[string]interface{}{
					"id":              tea.StringValue(bot.Id),
					"bot_name":        tea.StringValue(bot.BotName),
					"bot_description": tea.StringValue(bot.BotDescription),
					"bot_act":         tea.StringValue(bot.BotAct),
					"creator":         tea.StringValue(bot.Creator),
				}

				// Parse condition list
				if bot.ConditionList != nil {
					conditionList := make([]interface{}, len(bot.ConditionList))
					for j, condition := range bot.ConditionList {
						conditionList[j] = map[string]interface{}{
							"condition_name":       tea.StringValue(condition.ConditionName),
							"condition_value_list": tea.StringSliceValue(condition.ConditionValueList),
							"condition_func":       tea.StringValue(condition.ConditionFunc),
							"condition_key":        tea.StringValue(condition.ConditionKey),
						}
					}
					botMap["condition_list"] = conditionList
				}

				customBots[i] = botMap
			}
			configMap["custom_bots"] = customBots
		}

		configListData = append(configListData, configMap)
	}

	return configListData
}
