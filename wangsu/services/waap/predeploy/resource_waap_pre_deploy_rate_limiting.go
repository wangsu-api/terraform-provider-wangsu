package pre_deploy

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wangsuCommon "github.com/wangsu-api/terraform-provider-wangsu/wangsu/common"
	"github.com/wangsu-api/terraform-provider-wangsu/wangsu/services/waap"
	preDeploy "github.com/wangsu-api/wangsu-sdk-go/wangsu/waap/predeploy"
	"log"
	"time"
)

func ResourceWaapPreDeployRateLimiting() *schema.Resource {
	return &schema.Resource{

		CreateContext: resourceWaapPreDeployRateLimitingCreate,
		ReadContext:   resourceWaapPreDeployRateLimitingRead,
		UpdateContext: resourceWaapPreDeployRateLimitingCreate,
		DeleteContext: resourceWaapPreDeployRateLimitingRead,

		Schema: map[string]*schema.Schema{
			"host_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Host list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Domain.",
						},
						"host_address": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "IP address.",
						},
					},
				},
			},
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Domain list.",
			},
			"config_switch": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Policy switch.<br/>ON: Enable.<br/>OFF: Disable.",
			},
			"rule_list": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Rule list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Rule Name, maximum 50 characters.<br/>does not support # and & .",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description, maximum 200 characters.",
						},
						"scene": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Protected target.<br/>WEB:Website<br/>API:API",
						},
						"statistical_item": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Client identifier.<br/>IP:Client IP<br/>IP_UA:Client IP and User-Agent<br/>COOKIE:Cookie<br/>IP_COOKIE:Client IP and Cookie<br/>HEADER:Request Header<br/>When there is a status code in the matching condition,this client identifier is not supported.<br/>IP_HEADER:Client IP and Request Header<br/>When there is a status code in the matching condition,this client identifier is not supported .",
						},
						"statistics_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Statistical key value.<br/>When the client identifier is cookie/header value, the corresponding key value needs to be entered.",
						},
						"statistical_period": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Statistics period, unit: seconds, the range is 1 - 3600.",
						},
						"trigger_threshold": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Trigger threshold, unit: times.",
						},
						"intercept_time": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Action duration, unit: seconds, the range is 10 - 604800.",
						},
						"effective_status": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Cycle effective status.<br/>PERMANENT:All time<br/>WITHOUT:Excluded time<br/>WITHIN:Selected time",
						},
						"rate_limit_effective": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Effective time period.<br/>When the effective status is effective within the cycle or not effective within the cycle, this field must have a value.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"effective": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "Effective.<br/>MON:Monday<br/>TUE:Tuesday<br/>WED:Wednesday<br/>THU:Thursday<br/>FRI:Friday<br/>SAT:Saturday<br/>SUN:Sunday",
										Elem: &schema.Schema{
											Type:        schema.TypeString,
											Description: "Effective.<br/>MON:Monday<br/>TUE:Tuesday<br/>WED:Wednesday<br/>THU:Thursday<br/>FRI:Friday<br/>SAT:Saturday<br/>SUN:Sunday",
										},
									},
									"start": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Start time, format: HH:mm.",
									},
									"end": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "End time, format: HH:mm.",
									},
									"timezone": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Timezone,default value: GTM+8.",
									},
								},
							},
						},
						"asset_api_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "API ID under API business, multiple separated by ; sign.<br/>When the protected target is APIThis field is required.",
						},
						"action": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Action.<br/>NO_USE:Not Used<br/>LOG:Log<br/>COOKIE:Cookie verification<br/>JS_CHECK:Javascript verification<br/>DELAY:Delay<br/>BLOCK:Deny<br/>RESET:Reset Connection<br/>Custom response ID:Custom response ID<br/>When there is a status code in the matching condition, the supported actions are Log, Deny, NO_USE, and Reset, Connection.",
						},
						"rate_limit_rule_condition": {
							Type:        schema.TypeList,
							Required:    true,
							MaxItems:    1,
							Description: "Matching conditions.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip_or_ips_conditions": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "IP/CIDR, match type cannot be repeated.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"match_type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Match type.<br/>EQUAL:Equals<br/>NOT_EQUAL:Does not equal",
												},
												"ip_or_ips": {
													Type:        schema.TypeList,
													Required:    true,
													Description: "IP/CIDR, maximum 500 IP/CIDR.",
													Elem: &schema.Schema{
														Type:        schema.TypeString,
														Description: "IP/CIDR, maximum 500 IP/CIDR.",
													},
												},
											},
										},
									},
									"path_conditions": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Path, match type cannot be repeated.<br/>When the business scenario is API, this matching condition is not supported.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"match_type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Match type.<br/>EQUAL: Equals, path case sensitive<br/>NOT_EQUAL: Does not equal, path case sensitive<br/>CONTAIN: Contains, path case insensitive<br/>NOT_CONTAIN: Does not Contains, path case insensitive<br/>REGEX: Regex match, path case insensitive<br/>NOT_REGEX: Regular does not match, path case sensitive<br/>START_WITH: Starts with, path case sensitive<br/>END_WITH: Ends with, path case sensitive<br/>WILDCARD: Wildcard matches, path case sensitive,* represents zero or more arbitrary characters, ? represents any single character.<br/>NOT_WILDCARD: Wildcard does not match, path case sensitive,* represents zero or more arbitrary characters, ? represents any single character ",
												},
												"paths": {
													Type:        schema.TypeList,
													Required:    true,
													Description: "Path.<br/>When match type is EQUAL/NOT_EQUAL/START_WITH/END_WITH, path needs to start with \"/\", and no parameters.<br/>When the match type is REGEX/NOT_REGEX, only one value is allowed.<br/>Example: /test.html.",
													Elem: &schema.Schema{
														Type:        schema.TypeString,
														Description: "Path.<br/>When match type is EQUAL/NOT_EQUAL/START_WITH/END_WITH, path needs to start with \"/\", and no parameters.<br/>When the match type is REGEX/NOT_REGEX, only one value is allowed.<br/>Example: /test.html.",
													},
												},
											},
										},
									},
									"uri_conditions": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "URI, match type cannot be repeated.<br/>When the business scenario is API, this matching condition is not supported.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"match_type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Match type.<br/>EQUAL: Equals, URI case sensitive<br/>NOT_EQUAL: Does not equal, URI case sensitive<br/>CONTAIN: Contains, URI case insensitive<br/>NOT_CONTAIN: Does not Contains, URI case insensitive<br/>REGEX: Regex match, URI case insensitive<br/>NOT_REGEX: Regular does not match, URI case insensitive<br/>START_WITH: Starts with, URI case insensitive<br/>END_WITH: Ends with, URI case insensitive<br/>WILDCARD: Wildcard matches, URI case insensitive,* represents zero or more arbitrary characters, ? represents any single character<br/>NOT_WILDCARD: Wildcard does not match, URI case insensitive,* represents zero or more arbitrary characters, ? represents any single character",
												},
												"uri": {
													Type:        schema.TypeList,
													Required:    true,
													Description: "URI.<br/>When match type is EQUAL/NOT_EQUAL/START_WITH/END_WITH, uri needs to start with \"/\", and includes parameters.<br/>When the match type is REGEX/NOT_REGEX, only one value is allowed.<br/>Example: /test.html?id=1.",
													Elem: &schema.Schema{
														Type:        schema.TypeString,
														Description: "URI.<br/>When match type is EQUAL/NOT_EQUAL/START_WITH/END_WITH, uri needs to start with \"/\", and includes parameters.<br/>When the match type is REGEX/NOT_REGEX, only one value is allowed.<br/>Example: /test.html?id=1.",
													},
												},
											},
										},
									},
									"ua_conditions": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "User Agent, match type cannot be repeated.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"match_type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Match type.<br/>EQUAL: Equals, user agent case sensitive<br/>NOT_EQUAL: Does not equal, user agent case sensitive<br/>CONTAIN: Contains, user agent case insensitive<br/>NOT_CONTAIN: Does not Contains, user agent case insensitive<br/>NONE:Empty or non-existent<br/>REGEX: Regex match, user agent case insensitive<br/>NOT_REGEX: Regular does not match, user agent case insensitive<br/>START_WITH: Starts with, user agent case insensitive<br/>END_WITH: Ends with, user agent case insensitive<br/>WILDCARD: Wildcard matches, user agent case insensitive,* represents zero or more arbitrary characters, ? represents any single character<br/>NOT_WILDCARD: Wildcard does not match, user agent case insensitive,* represents zero or more arbitrary characters, ? represents any single character",
												},
												"ua": {
													Type:        schema.TypeList,
													Required:    true,
													Description: "User agent.<br/>When the match type is REGEX/NOT_REGEX, only one value is allowed.<br/>Example: go-Http-client/1.1.",
													Elem: &schema.Schema{
														Type:        schema.TypeString,
														Description: "User agent.<br/>When the match type is REGEX/NOT_REGEX, only one value is allowed.<br/>Example: go-Http-client/1.1.",
													},
												},
											},
										},
									},
									"referer_conditions": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Referer, match type cannot be repeated.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"match_type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Match type.<br/>EQUAL: Equals, referer case sensitive<br/>NOT_EQUAL: Does not equal, referer case sensitive<br/>CONTAIN: Contains, referer case insensitive<br/>NOT_CONTAIN: Does not Contains, referer case insensitive<br/>NONE:Empty or non-existent<br/>REGEX: Regex match, referer case insensitive<br/>NOT_REGEX: Regular does not match, referer case insensitive<br/>START_WITH: Starts with, referer case insensitive<br/>END_WITH: Ends with, referer case insensitive<br/>WILDCARD: Wildcard matches, referer case insensitive,* represents zero or more arbitrary characters, ? represents any single characte<br/>NOT_WILDCARD: Wildcard does not match, referer case insensitive,* represents zero or more arbitrary characters, ? represents any single character",
												},
												"referer": {
													Type:        schema.TypeList,
													Required:    true,
													Description: "Referer.<br/>When the match type is REGEX/NOT_REGEX, only one value is allowed.<br/>Example: http://test.com.",
													Elem: &schema.Schema{
														Type:        schema.TypeString,
														Description: "Referer.<br/>When the match type is REGEX/NOT_REGEX, only one value is allowed.<br/>Example: http://test.com.",
													},
												},
											},
										},
									},
									"header_conditions": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Request Header, match type can be repeated.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"match_type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Match type.<br/>EQUAL: Equals, request header values case sensitive.<br/>NOT_EQUAL: Does not equal, request header values case sensitive.<br/>CONTAIN: Contains, request header values case insensitive.<br/>NOT_CONTAIN: Does not Contains, request header values case insensitive.<br/>NONE: Empty or non-existent.<br/>REGEX: Regex match, request header values case insensitive.<br/>NOT_REGEX: Regular does not match, request header values case insensitive.<br/>START_WITH: Starts with, request header values case insensitive<br/>END_WITH: Ends with, request header values case insensitive.<br/>WILDCARD: Wildcard matches, request header values case insensitive,* represents zero or more arbitrary characters, ? represents any single character.<br/>NOT_WILDCARD: Wildcard does not match, request header values case insensitive,* represents zero or more arbitrary characters, ? represents any single character.",
												},
												"key": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Header name,case insensitive,up to 100 characters.<br/>Example: Accept.",
												},
												"value_list": {
													Type:        schema.TypeList,
													Required:    true,
													Description: "Header value.<br/>When the match type is REGEX/NOT_REGEX, only one value is allowed.",
													Elem: &schema.Schema{
														Type:        schema.TypeString,
														Description: "Header value.<br/>When the match type is REGEX/NOT_REGEX, only one value is allowed.",
													},
												},
											},
										},
									},
									"area_conditions": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Geo, match type cannot be repeated.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"match_type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Match type.<br/>EQUAL:Equals<br/>NOT_EQUAL:Does not equal",
												},
												"areas": {
													Type:        schema.TypeList,
													Required:    true,
													Description: "Geo.",
													Elem: &schema.Schema{
														Type:        schema.TypeString,
														Description: "Geo.",
													},
												},
											},
										},
									},
									"method_conditions": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Request Method.<br/>When the business scenario is API,this matching condition is not supported.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"match_type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Match type.<br/>EQUAL:Equals<br/>NOT_EQUAL:Does not equal",
												},
												"request_method": {
													Type:        schema.TypeList,
													Required:    true,
													Description: "Request method.<br/>Supported values: GET/POST/DELETE/PUT/HEAD/OPTIONS/COPY.",
													Elem: &schema.Schema{
														Type:        schema.TypeString,
														Description: "Request method.Supported values: GET/POST/DELETE/PUT/HEAD/OPTIONS/COPY.",
													},
												},
											},
										},
									},
									"scheme_conditions": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "HTTP/S, match type cannot be repeated.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"match_type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Match type.<br/>EQUAL:Equals<br/>NOT_EQUAL:Does not equal",
												},
												"scheme": {
													Type:        schema.TypeList,
													Required:    true,
													Description: "HTTP/S.<br/>Supported values: HTTP/HTTPS.",
													Elem: &schema.Schema{
														Type:        schema.TypeString,
														Description: "HTTP/S.<br/>Supported values: HTTP/HTTPS.",
													},
												},
											},
										},
									},
									"status_code_conditions": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Response Code, match type cannot be repeated.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"match_type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Match type.<br/>EQUAL:Equals<br/>NOT_EQUAL:Does not equal",
												},
												"status_code": {
													Type:        schema.TypeList,
													Required:    true,
													Description: "Response Code.",
													Elem: &schema.Schema{
														Type:        schema.TypeString,
														Description: "Response Code.",
													},
												},
											},
										},
									},
									"ja3_conditions": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "JA3 Fingerprint, match type cannot be repeated.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"match_type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Match type.<br/>EQUAL: Equals<br/>NOT_EQUAL: Does not equal",
												},
												"ja3_list": {
													Type:        schema.TypeList,
													Required:    true,
													Description: "JA3 Fingerprint List, maximum 300 JA3 Fingerprint.<br/>When the match type is EQUAL/NOT_EQUAL, each item's character length must be 32 and can only include numbers and lowercase letters.",
													Elem: &schema.Schema{
														Type:        schema.TypeString,
														Description: "JA3 Fingerprint List, maximum 300 JA3 Fingerprint.<br/>When the match type is EQUAL/NOT_EQUAL, each item's character length must be 32 and can only include numbers and lowercase letters.",
													},
												},
											},
										},
									},
									"ja4_conditions": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "JA4 Fingerprint, match type cannot be repeated.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"match_type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Match type.<br/>EQUAL: Equals<br/>NOT_EQUAL: Does not equal<br/>CONTAIN: Contains<br/>NOT_CONTAIN: Does not Contains<br/>START_WITH: Starts with<br/>END_WITH: Ends with<br/>WILDCARD: Wildcard matches, ** represents zero or more arbitrary characters, ? represents any single character<br/>NOT_WILDCARD: Wildcard does not match, ** represents zero or more arbitrary characters, ? represents any single character",
												},
												"ja4_list": {
													Type:        schema.TypeList,
													Required:    true,
													Description: "JA4 Fingerprint List, maximum 300 JA4 Fingerprint.<br/>When the match type is EQUAL/NOT_EQUAL, each item's format must be 10 characters + 12 characters + 12 characters, separated by underscores, and can only include underscores, numbers, and lowercase letters.<br/>When the match type is CONTAIN/NOT_CONTAIN/START_WITH/END_WITH, each item is only allowed to include underscores, numbers, and lowercase letters.<br/>When the match type is WILDCARD/NOT_WILDCARD, each item, aside from  ** and ?, is only allowed to include underscores, numbers, and lowercase letters.",
													Elem: &schema.Schema{
														Type:        schema.TypeString,
														Description: "JA4 Fingerprint List, maximum 300 JA4 Fingerprint.<br/>When the match type is EQUAL/NOT_EQUAL, each item's format must be 10 characters + 12 characters + 12 characters, separated by underscores, and can only include underscores, numbers, and lowercase letters.<br/>When the match type is CONTAIN/NOT_CONTAIN/START_WITH/END_WITH, each item is only allowed to include underscores, numbers, and lowercase letters.<br/>When the match type is WILDCARD/NOT_WILDCARD, each item, aside from  ** and ?, is only allowed to include underscores, numbers, and lowercase letters.",
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
		},
	}
}

func resourceWaapPreDeployRateLimitingCreate(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("resource.wangsu_pre_deploy_rate_limiting.read")

	var diags diag.Diagnostics

	request := &preDeploy.PreDeployRateLimitingConfigurationRequest{}
	if domain, ok := data.Get("domain").(string); ok {
		request.Domain = &domain
	}
	if configSwitch, ok := data.Get("config_switch").(string); ok {
		request.ConfigSwitch = &configSwitch
	}

	ruleList := data.Get("rule_list").([]interface{})
	ruleListRequest := make([]*preDeploy.RateLimitingRule, len(ruleList))
	for i, rule := range ruleList {
		ruleMap := rule.(map[string]interface{})
		ruleRequest := &preDeploy.RateLimitingRule{}

		if v, ok := ruleMap["rule_name"]; ok {
			ruleRequest.SetRuleName(v.(string))
		}
		if v, ok := ruleMap["description"]; ok {
			ruleRequest.SetDescription(v.(string))
		}
		if v, ok := ruleMap["scene"]; ok {
			ruleRequest.SetScene(v.(string))
		}
		if v, ok := ruleMap["statistical_item"]; ok {
			ruleRequest.SetStatisticalItem(v.(string))
		}
		if v, ok := ruleMap["statistics_key"]; ok {
			ruleRequest.SetStatisticsKey(v.(string))
		}
		if v, ok := ruleMap["statistical_period"]; ok {
			ruleRequest.SetStatisticalPeriod(v.(int))
		}
		if v, ok := ruleMap["trigger_threshold"]; ok {
			ruleRequest.SetTriggerThreshold(v.(int))
		}
		if v, ok := ruleMap["intercept_time"]; ok {
			ruleRequest.SetInterceptTime(v.(int))
		}
		if v, ok := ruleMap["effective_status"]; ok {
			ruleRequest.SetEffectiveStatus(v.(string))
		}
		if v, ok := ruleMap["rate_limit_effective"]; ok {
			rateLimitEffectiveV := v.([]interface{})
			for _, v := range rateLimitEffectiveV {
				rateLimitEffectiveD := v.(map[string]interface{})
				effective := rateLimitEffectiveD["effective"].([]interface{})
				start := rateLimitEffectiveD["start"].(string)
				end := rateLimitEffectiveD["end"].(string)
				timezone := rateLimitEffectiveD["timezone"].(string)
				rateLimitEffective := &preDeploy.RateLimitEffective{}
				effectiveList := make([]*string, len(effective))
				for i, v := range effective {
					str := v.(string)
					effectiveList[i] = &str
				}
				rateLimitEffective.SetEffective(effectiveList)
				rateLimitEffective.SetStart(start)
				rateLimitEffective.SetEnd(end)
				rateLimitEffective.SetTimezone(timezone)

				ruleRequest.SetRateLimitEffective(rateLimitEffective)
			}
		}
		if v, ok := ruleMap["asset_api_id"]; ok {
			ruleRequest.SetAssetApiId(v.(string))
		}
		if v, ok := ruleMap["action"]; ok {
			ruleRequest.SetAction(v.(string))
		}

		// Parse conditions
		if ruleMap["rate_limit_rule_condition"] != nil {
			conditions := ruleMap["rate_limit_rule_condition"].([]interface{})
			conditionsRequest := &preDeploy.RateLimitingRuleCondition{}
			for _, condition := range conditions {
				conditionMap := condition.(map[string]interface{})

				// URI Parameter conditions
				if conditionMap["uri_param_conditions"] != nil {
					uriParamConditions := make([]*preDeploy.UriParamConditions, 0)
					for _, uriParamCondition := range conditionMap["uri_param_conditions"].([]interface{}) {
						uriParamConditionMap := uriParamCondition.(map[string]interface{})
						matchType := uriParamConditionMap["match_type"].(string)
						paramName := uriParamConditionMap["param_name"].(string)
						paramValue := waap.ConvertToStringSlice(uriParamConditionMap["param_value"].([]interface{}))
						uriParamConditions = append(uriParamConditions, &preDeploy.UriParamConditions{
							MatchType:  &matchType,
							ParamName:  &paramName,
							ParamValue: paramValue,
						})
					}
					conditionsRequest.UriParamConditions = uriParamConditions
				}

				// IP or IPs conditions
				if conditionMap["ip_or_ips_conditions"] != nil {
					ipOrIpsConditions := make([]*preDeploy.IpOrIpsConditions, 0)
					for _, ipOrIpsCondition := range conditionMap["ip_or_ips_conditions"].([]interface{}) {
						ipOrIpsConditionMap := ipOrIpsCondition.(map[string]interface{})
						matchType := ipOrIpsConditionMap["match_type"].(string)
						ipOrIps := waap.ConvertToStringSlice(ipOrIpsConditionMap["ip_or_ips"].([]interface{}))
						ipOrIpsConditions = append(ipOrIpsConditions, &preDeploy.IpOrIpsConditions{
							MatchType: &matchType,
							IpOrIps:   ipOrIps,
						})
					}
					conditionsRequest.IpOrIpsConditions = ipOrIpsConditions
				}

				// Path conditions
				if conditionMap["path_conditions"] != nil {
					pathConditions := make([]*preDeploy.PathConditions, 0)
					for _, pathCondition := range conditionMap["path_conditions"].([]interface{}) {
						pathConditionMap := pathCondition.(map[string]interface{})
						matchType := pathConditionMap["match_type"].(string)
						paths := waap.ConvertToStringSlice(pathConditionMap["paths"].([]interface{}))
						pathConditions = append(pathConditions, &preDeploy.PathConditions{
							MatchType: &matchType,
							Paths:     paths,
						})
					}
					conditionsRequest.PathConditions = pathConditions
				}

				// URI conditions
				if conditionMap["uri_conditions"] != nil {
					uriConditions := make([]*preDeploy.UriConditions, 0)
					for _, uriCondition := range conditionMap["uri_conditions"].([]interface{}) {
						uriConditionMap := uriCondition.(map[string]interface{})
						matchType := uriConditionMap["match_type"].(string)
						uri := waap.ConvertToStringSlice(uriConditionMap["uri"].([]interface{}))
						uriConditions = append(uriConditions, &preDeploy.UriConditions{
							MatchType: &matchType,
							Uri:       uri,
						})
					}
					conditionsRequest.UriConditions = uriConditions
				}

				// UA conditions
				if conditionMap["ua_conditions"] != nil {
					uaConditions := make([]*preDeploy.UaConditions, 0)
					for _, uaCondition := range conditionMap["ua_conditions"].([]interface{}) {
						uaConditionMap := uaCondition.(map[string]interface{})
						matchType := uaConditionMap["match_type"].(string)
						ua := waap.ConvertToStringSlice(uaConditionMap["ua"].([]interface{}))
						uaConditions = append(uaConditions, &preDeploy.UaConditions{
							MatchType: &matchType,
							Ua:        ua,
						})
					}
					conditionsRequest.UaConditions = uaConditions
				}

				// Referer conditions
				if conditionMap["referer_conditions"] != nil {
					refererConditions := make([]*preDeploy.RefererConditions, 0)
					for _, refererCondition := range conditionMap["referer_conditions"].([]interface{}) {
						refererConditionMap := refererCondition.(map[string]interface{})
						matchType := refererConditionMap["match_type"].(string)
						referer := waap.ConvertToStringSlice(refererConditionMap["referer"].([]interface{}))
						refererConditions = append(refererConditions, &preDeploy.RefererConditions{
							MatchType: &matchType,
							Referer:   referer,
						})
					}
					conditionsRequest.RefererConditions = refererConditions
				}

				// Header conditions
				if conditionMap["header_conditions"] != nil {
					headerConditions := make([]*preDeploy.HeaderConditions, 0)
					for _, headerCondition := range conditionMap["header_conditions"].([]interface{}) {
						headerConditionMap := headerCondition.(map[string]interface{})
						matchType := headerConditionMap["match_type"].(string)
						key := headerConditionMap["key"].(string)
						valueList := waap.ConvertToStringSlice(headerConditionMap["value_list"].([]interface{}))
						headerConditions = append(headerConditions, &preDeploy.HeaderConditions{
							MatchType: &matchType,
							Key:       &key,
							ValueList: valueList,
						})
					}
					conditionsRequest.HeaderConditions = headerConditions
				}

				// Geo conditions
				if conditionMap["area_conditions"] != nil {
					areaConditions := make([]*preDeploy.AreaConditions, 0)
					for _, areaCondition := range conditionMap["area_conditions"].([]interface{}) {
						areaConditionMap := areaCondition.(map[string]interface{})
						matchType := areaConditionMap["match_type"].(string)
						areas := waap.ConvertToStringSlice(areaConditionMap["areas"].([]interface{}))
						areaConditions = append(areaConditions, &preDeploy.AreaConditions{
							MatchType: &matchType,
							Areas:     areas,
						})
					}
					conditionsRequest.AreaConditions = areaConditions
				}

				// Method conditions
				if conditionMap["method_conditions"] != nil {
					methodConditions := make([]*preDeploy.MethodConditions, 0)
					for _, methodCondition := range conditionMap["method_conditions"].([]interface{}) {
						methodConditionMap := methodCondition.(map[string]interface{})
						matchType := methodConditionMap["match_type"].(string)
						requestMethods := waap.ConvertToStringSlice(methodConditionMap["request_method"].([]interface{}))
						methodConditions = append(methodConditions, &preDeploy.MethodConditions{
							MatchType:     &matchType,
							RequestMethod: requestMethods,
						})
					}
					conditionsRequest.MethodConditions = methodConditions
				}

				// JA3 conditions
				if conditionMap["ja3_conditions"] != nil {
					ja3Conditions := make([]*preDeploy.Ja3Conditions, 0)
					for _, ja3Condition := range conditionMap["ja3_conditions"].([]interface{}) {
						ja3ConditionMap := ja3Condition.(map[string]interface{})
						matchType := ja3ConditionMap["match_type"].(string)
						ja3List := waap.ConvertToStringSlice(ja3ConditionMap["ja3_list"].([]interface{}))
						ja3Conditions = append(ja3Conditions, &preDeploy.Ja3Conditions{
							MatchType: &matchType,
							Ja3List:   ja3List,
						})
					}
					conditionsRequest.Ja3Conditions = ja3Conditions
				}

				// JA4 conditions
				if conditionMap["ja4_conditions"] != nil {
					ja4Conditions := make([]*preDeploy.Ja4Conditions, 0)
					for _, ja4Condition := range conditionMap["ja4_conditions"].([]interface{}) {
						ja4ConditionMap := ja4Condition.(map[string]interface{})
						matchType := ja4ConditionMap["match_type"].(string)
						ja4List := waap.ConvertToStringSlice(ja4ConditionMap["ja4_list"].([]interface{}))
						ja4Conditions = append(ja4Conditions, &preDeploy.Ja4Conditions{
							MatchType: &matchType,
							Ja4List:   ja4List,
						})
					}
					conditionsRequest.Ja4Conditions = ja4Conditions
				}

				// Schema conditions
				if conditionMap["scheme_conditions"] != nil {
					schemeConditions := make([]*preDeploy.SchemeConditions, 0)
					for _, schemeCondition := range conditionMap["scheme_conditions"].([]interface{}) {
						schemeConditionMap := schemeCondition.(map[string]interface{})
						matchType := schemeConditionMap["match_type"].(string)
						scheme := waap.ConvertToStringSlice(schemeConditionMap["scheme"].([]interface{}))
						schemeConditions = append(schemeConditions, &preDeploy.SchemeConditions{
							MatchType: &matchType,
							Scheme:    scheme,
						})
					}
					conditionsRequest.SchemeConditions = schemeConditions
				}

				// Status code conditions
				if conditionMap["status_code_conditions"] != nil {
					statusCodeConditions := make([]*preDeploy.StatusCodeConditions, 0)
					for _, statusCodeCondition := range conditionMap["status_code_conditions"].([]interface{}) {
						statusCodeConditionMap := statusCodeCondition.(map[string]interface{})
						matchType := statusCodeConditionMap["match_type"].(string)
						statusCode := waap.ConvertToStringSlice(statusCodeConditionMap["status_code"].([]interface{}))
						statusCodeConditions = append(statusCodeConditions, &preDeploy.StatusCodeConditions{
							MatchType:  &matchType,
							StatusCode: statusCode,
						})
					}
					conditionsRequest.StatusCodeConditions = statusCodeConditions
				}
			}
			ruleRequest.RateLimitRuleCondition = conditionsRequest
		}
		ruleListRequest[i] = ruleRequest
	}
	request.RuleList = ruleListRequest

	var response *preDeploy.PreDeployRateLimitingConfigurationResponse
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseWaapPreDeployClient().PreDeployRateLimitingConfiguration(request)
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

	data.SetId(*response.Data.PreId)

	// 轮询获取部署结果
	getRequest := &preDeploy.GetPreDeployResultRequest{}
	getResponse := &preDeploy.GetPreDeployResultResponse{}
	getRequest.PreId = response.Data.PreId
	for {
		err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
			_, getResponse, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseWaapPreDeployClient().GetPreDeployResult(getRequest)
			if err != nil {
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
			return diags
		}

		if getResponse == nil {
			return nil
		}

		if *getResponse.Data.DeployStatus == "SUCCESS" {
			hostList := make([]map[string]interface{}, len(getResponse.Data.HostList))
			for i, host := range getResponse.Data.HostList {
				hostList[i] = map[string]interface{}{
					"host_name":    host.HostName,
					"host_address": host.HostAddress,
				}
			}
			_ = data.Set("host_list", hostList)
			break
		} else if *getResponse.Data.DeployStatus == "FAIL" {
			log.Println("Deployment failed!")
			break
		} else {
			log.Println("Deployment in progress, retrying...")
			time.Sleep(10 * time.Second)
		}
	}

	return diags
}

func resourceWaapPreDeployRateLimitingRead(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("resource.wangsu_pre_deploy_rate_limiting.read")

	data.SetId("")
	return nil
}
