package ratelimit

import (
	"context"
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wangsuCommon "github.com/wangsu-api/terraform-provider-wangsu/wangsu/common"
	waapRatelimit "github.com/wangsu-api/wangsu-sdk-go/wangsu/waap/ratelimit"
	"log"
	"time"
)

func ResourceWaapRateLimit() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWaapRateLimitCreate,
		ReadContext:   resourceWaapRateLimitRead,
		UpdateContext: resourceWaapRateLimitUpdate,
		DeleteContext: resourceWaapRateLimitDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Rule ID.",
			},

			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Hostname.",
			},
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
				Description: "Statistics period, unit: seconds.",
			},
			"trigger_threshold": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Trigger threshold, unit: times.",
			},
			"intercept_time": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Action duration, unit: seconds.",
			},
			"effective_status": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cycle effective status.<br/>PERMANENT:All time<br/>WITHOUT:Excluded time<br/>WITHIN:Selected time",
			},
			"rate_limit_effective": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Effective time period.When the effective status is effective within the cycle or not effective within the cycle, this field must have a value.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"effective": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Effective.<br/>MON:Monday<br/>TUE:Tuesday<br/>WED:Wednesday<br/>THU:Thursday<br/>FRI:Friday<br/>SAT:Saturday<br/>SUN:Sunday",
							Elem: &schema.Schema{
								Type: schema.TypeString,
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
				Description: "Action.<br/>NO_USE:Not Used<br/>LOG:Log<br/>COOKIE:Cookie verification<br/>JS_CHECK:Javascript verification<br/>DELAY:Delay<br/>BLOCK:Deny<br/>RESET:Reset Connection<br/>JSC:Interactive Captcha<br/>Custom response ID:Custom response ID<br/>When there is a status code in the matching condition, the supported actions are Log, Deny,Not Used, and Reset Connection.",
			},
			"rate_limit_rule_condition": {
				Type:        schema.TypeList,
				Required:    true,
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
										Description: "Match type.<br/>EQUAL:Equal<br/>NOT_EQUAL:Does not equal",
									},
									"ip_or_ips": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "IP/CIDR, maximum 300 IP/CIDR.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
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
										Description: "Match type.<br/>EQUAL: Equals, user agent case sensitive<br/>NOT_EQUAL: Does not equal, user agent case sensitive<br/>CONTAIN: Contains, user agent case insensitive<br/>NOT_CONTAIN: Does not Contains, user agent case insensitive<br/>REGEX: Regex match, user agent case insensitive<br/>NOT_REGEX: Regular does not match, user agent case insensitive<br/>START_WITH: Starts with, user agent case insensitive<br/>END_WITH: Ends with, user agent case insensitive<br/>WILDCARD: Wildcard matches, user agent case insensitive, * represents zero or more arbitrary characters, ? represents any single character<br/>NOT_WILDCARD: Wildcard does not match, user agent case insensitive, * represents zero or more arbitrary characters, ? represents any single character",
									},
									"paths": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "Path.<br/>When match type is EQUAL/NOT_EQUAL/START_WITH/END_WITH, path needs to start with \"/\", and no parameters.<br/>When the match type is REGEX/NOT_REGEX, only one value is allowed. <br/>Example: /test.html.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
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
										Description: "Match type.<br/>EQUAL: Equals, user agent case sensitive<br/>NOT_EQUAL: Does not equal, user agent case sensitive<br/>CONTAIN: Contains, user agent case insensitive<br/>NOT_CONTAIN: Does not Contains, user agent case insensitive<br/>REGEX: Regex match, user agent case insensitive<br/>NOT_REGEX: Regular does not match, user agent case insensitive<br/>START_WITH: Starts with, user agent case insensitive<br/>END_WITH: Ends with, user agent case insensitive<br/>WILDCARD: Wildcard matches, user agent case insensitive, * represents zero or more arbitrary characters, ? represents any single character<br/>NOT_WILDCARD: Wildcard does not match, user agent case insensitive, * represents zero or more arbitrary characters, ? represents any single character",
									},
									"uri": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "URI.<br/>When match type is EQUAL/NOT_EQUAL/START_WITH/END_WITH, uri needs to start with \"/\", and includes parameters.<br/>When the match type is REGEX/NOT_REGEX, only one value is allowed. <br/>Example: /test.html?id=1.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"uri_param_conditions": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "URI ParameterI, match type cannot be repeated.<br/>When the business scenario is API, this matching condition is not supported.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"match_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Match type.<br/>EQUAL:Equals,param value case sensitive<br/>NOT_EQUAL:Does not equal,param value case sensitive<br/>CONTAIN:Contains,param value case insensitive<br/>NOT_CONTAIN:Does not contains,param value case insensitive<br/>REGEX:Regex match,param value case insensitive<br/>NONE:Empty or non-existent",
									},
									"param_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Param name,case sensitive,maximum 100 characters.<br/>Example: id.",
									},
									"param_value": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "Param value.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
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
										Description: "Match type.<br/>EQUAL: Equals, user agent case sensitive<br/>NOT_EQUAL: Does not equal, user agent case sensitive<br/>CONTAIN: Contains, user agent case insensitive<br/>NOT_CONTAIN: Does not Contains, user agent case insensitive<br/>NONE:Empty or non-existent<br/>REGEX: Regex match, user agent case insensitive<br/>NOT_REGEX: Regular does not match, user agent case insensitive<br/>START_WITH: Starts with, user agent case insensitive<br/>END_WITH: Ends with, user agent case insensitive<br/>WILDCARD: Wildcard matches, user agent case insensitive, * represents zero or more arbitrary characters, ? represents any single character<br/>NOT_WILDCARD: Wildcard does not match, user agent case insensitive, * represents zero or more arbitrary characters, ? represents any single character",
									},
									"ua": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "User agent.<br/>When the match type is REGEX/NOT_REGEX, only one value is allowed. <br/>Example: go-Http-client/1.1.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
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
										Description: "Match type.<br/>EQUAL: Equals, user agent case sensitive<br/>NOT_EQUAL: Does not equal, user agent case sensitive<br/>CONTAIN: Contains, user agent case insensitive<br/>NOT_CONTAIN: Does not Contains, user agent case insensitive<br/>NONE:Empty or non-existent<br/>REGEX: Regex match, user agent case insensitive<br/>NOT_REGEX: Regular does not match, user agent case insensitive<br/>START_WITH: Starts with, user agent case insensitive<br/>END_WITH: Ends with, user agent case insensitive<br/>WILDCARD: Wildcard matches, user agent case insensitive, * represents zero or more arbitrary characters, ? represents any single character<br/>NOT_WILDCARD: Wildcard does not match, user agent case insensitive, * represents zero or more arbitrary characters, ? represents any single character",
									},
									"referer": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "Referer.<br/>When the match type is REGEX/NOT_REGEX, only one value is allowed. <br/>Example: http://test.com.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
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
										Description: "Match type.<br/>EQUAL: Equals, request header values case sensitive<br/>NOT_EQUAL: Does not equal, request header values case sensitive<br/>CONTAIN: Contains, request header values case insensitive<br/>NOT_CONTAIN: Does not Contains, request header values case insensitive<br/>NONE: Empty or non-existent<br/>REGEX: Regex match, request header values case insensitive<br/>NOT_REGEX: Regular does not match, request header values case insensitive<br/>START_WITH: Starts with, request header values case insensitive<br/>END_WITH: Ends with, request header values case insensitive<br/>WILDCARD: Wildcard matches, request header values case insensitive, * represents zero or more arbitrary characters, ? represents any single character<br/>NOT_WILDCARD: Wildcard does not match, request header values case insensitive, * represents zero or more arbitrary characters, ? represents any single character",
									},
									"key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Header name,case insensitive,maximum 100 characters.<br/>Example: Accept.",
									},
									"value_list": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "Header value.<br/>When the match type is REGEX/NOT_REGEX, only one value is allowed.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"area_conditions": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Geo,match type cannot be repeated.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"match_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Match type.<br/>EQUAL:Equal<br/>NOT_EQUAL:Does not equal",
									},
									"areas": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "Geo.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
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
										Description: "Match type.<br/>EQUAL:Equal<br/>NOT_EQUAL:Does not equal",
									},
									"status_code": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "Response Code.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
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
										Description: "Match type.<br/>EQUAL:Equal<br/>NOT_EQUAL:Does not equal",
									},
									"request_method": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "Request method.<br/>Supported values: GET/POST/DELETE/PUT/HEAD/OPTIONS/COPY.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
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
										Description: "Match type.<br/>EQUAL:Equal<br/>NOT_EQUAL:Does not equal",
									},
									"scheme": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "HTTP/S.<br/>Supported values: HTTP/HTTPS.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
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
										Optional:    true,
										Description: "Match type.\nEQUAL: Equals\nNOT_EQUAL: Does not equal",
									},
									"ja3_list": {
										Type:        schema.TypeList,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "JA3 Fingerprint List, maximum 300 JA3 Fingerprint.\nWhen the match type is EQUAL/NOT_EQUAL, each item's character length must be 32 and can only include numbers and lowercase letters.",
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
										Optional:    true,
										Description: "Match type. \nEQUAL: Equals\nNOT_EQUAL: Does not equal\nCONTAIN: Contains\nNOT_CONTAIN: Does not Contains\nSTART_WITH: Starts with\nEND_WITH: Ends with\nWILDCARD: Wildcard matches, ** represents zero or more arbitrary characters, ? represents any single character\nNOT_WILDCARD: Wildcard does not match, ** represents zero or more arbitrary characters, ? represents any single character",
									},
									"ja4_list": {
										Type:        schema.TypeList,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "JA4 Fingerprint List, maximum 300 JA4 Fingerprint.\nWhen the match type is EQUAL/NOT_EQUAL, each item's format must be 10 characters + 12 characters + 12 characters, separated by underscores, and can only include underscores, numbers, and lowercase letters.\nWhen the match type is CONTAIN/NOT_CONTAIN/START_WITH/END_WITH, each item is only allowed to include underscores, numbers, and lowercase letters.\nWhen the match type is WILDCARD/NOT_WILDCARD, each item, aside from  ** and ?, is only allowed to include underscores, numbers, and lowercase letters.",
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

func resourceWaapRateLimitCreate(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_waap_ratelimit.create")

	var diags diag.Diagnostics
	request := &waapRatelimit.CreatRateLimitingRuleRequest{}
	if v, ok := data.GetOk("domain"); ok {
		request.SetDomain(v.(string))
	}
	if v, ok := data.GetOk("rule_name"); ok {
		request.SetRuleName(v.(string))
	}
	if v, ok := data.GetOk("description"); ok {
		request.SetDescription(v.(string))
	}
	if v, ok := data.GetOk("scene"); ok {
		request.SetScene(v.(string))
	}
	if v, ok := data.GetOk("statistical_item"); ok {
		request.SetStatisticalItem(v.(string))
	}
	if v, ok := data.GetOk("statistics_key"); ok {
		request.SetStatisticsKey(v.(string))
	}
	if v, ok := data.GetOk("statistical_period"); ok {
		request.SetStatisticalPeriod(v.(int))
	}
	if v, ok := data.GetOk("trigger_threshold"); ok {
		request.SetTriggerThreshold(v.(int))
	}
	if v, ok := data.GetOk("intercept_time"); ok {
		request.SetInterceptTime(v.(int))
	}
	if v, ok := data.GetOk("effective_status"); ok {
		request.SetEffectiveStatus(v.(string))
	}
	if v, ok := data.GetOk("rate_limit_effective"); ok {
		rateLimitEffectiveV := v.([]interface{})
		for _, v := range rateLimitEffectiveV {
			rateLimitEffectiveD := v.(map[string]interface{})
			effective := rateLimitEffectiveD["effective"].([]interface{})
			start := rateLimitEffectiveD["start"].(string)
			end := rateLimitEffectiveD["end"].(string)
			timezone := rateLimitEffectiveD["timezone"].(string)
			ratelimitEffective := &waapRatelimit.RateLimitEffective{}
			effectives := make([]*string, len(effective))
			for i, v := range effective {
				str := v.(string)
				effectives[i] = &str
			}
			ratelimitEffective.SetEffective(effectives)
			ratelimitEffective.SetStart(start)
			ratelimitEffective.SetEnd(end)
			ratelimitEffective.SetTimezone(timezone)

			request.SetRateLimitEffective(ratelimitEffective)
		}
	}
	if v, ok := data.GetOk("asset_api_id"); ok {
		request.SetAssetApiId(v.(string))
	}
	if v, ok := data.GetOk("action"); ok {
		request.SetAction(v.(string))
	}
	conditions := data.Get("rate_limit_rule_condition").([]interface{})
	conditionsRequest := &waapRatelimit.RateLimitRuleCondition{}
	for _, v := range conditions {
		conditionMap := v.(map[string]interface{})
		// IpOrIps Conditions
		if conditionMap["ip_or_ips_conditions"] != nil {
			ipOrIpsConditions := make([]*waapRatelimit.IpOrIpsCondition, 0)
			for _, ipOrIpsCondition := range conditionMap["ip_or_ips_conditions"].([]interface{}) {
				ipOrIpsConditionMap := ipOrIpsCondition.(map[string]interface{})
				matchType := ipOrIpsConditionMap["match_type"].(string)
				ipOrIpsInterface := ipOrIpsConditionMap["ip_or_ips"].([]interface{})
				ipOrIps := make([]*string, len(ipOrIpsInterface))
				for i, v := range ipOrIpsInterface {
					str := v.(string)
					ipOrIps[i] = &str
				}
				ipOrIpsCondition := &waapRatelimit.IpOrIpsCondition{
					MatchType: &matchType,
					IpOrIps:   ipOrIps,
				}
				ipOrIpsConditions = append(ipOrIpsConditions, ipOrIpsCondition)
			}
			conditionsRequest.IpOrIpsConditions = ipOrIpsConditions
		}

		// Path Conditions
		if conditionMap["path_conditions"] != nil {
			pathConditions := make([]*waapRatelimit.PathCondition, 0)
			for _, pathCondition := range conditionMap["path_conditions"].([]interface{}) {
				pathConditionMap := pathCondition.(map[string]interface{})
				matchType := pathConditionMap["match_type"].(string)
				pathsInterface := pathConditionMap["paths"].([]interface{})
				paths := make([]*string, len(pathsInterface))
				for i, v := range pathsInterface {
					str := v.(string)
					paths[i] = &str
				}
				pathCondition := &waapRatelimit.PathCondition{
					MatchType: &matchType,
					Paths:     paths,
				}
				pathConditions = append(pathConditions, pathCondition)
			}
			conditionsRequest.PathConditions = pathConditions
		}

		// URI Conditions
		if conditionMap["uri_conditions"] != nil {
			uriConditions := make([]*waapRatelimit.UriCondition, 0)
			for _, uriCondition := range conditionMap["uri_conditions"].([]interface{}) {
				uriConditionMap := uriCondition.(map[string]interface{})
				matchType := uriConditionMap["match_type"].(string)
				uriInterface := uriConditionMap["uri"].([]interface{})
				uri := make([]*string, len(uriInterface))
				for i, v := range uriInterface {
					str := v.(string)
					uri[i] = &str
				}
				uriCondition := &waapRatelimit.UriCondition{
					MatchType: &matchType,
					Uri:       uri,
				}
				uriConditions = append(uriConditions, uriCondition)
			}
			conditionsRequest.SetUriConditions(uriConditions)
		}

		// URI Param Conditions
		if conditionMap["uri_param_conditions"] != nil {
			uriParamConditions := make([]*waapRatelimit.UriParamCondition, 0)
			for _, uriParamCondition := range conditionMap["uri_param_conditions"].([]interface{}) {
				uriParamConditionMap := uriParamCondition.(map[string]interface{})
				matchType := uriParamConditionMap["match_type"].(string)
				paramName := uriParamConditionMap["param_name"].(string)
				paramValueInterface := uriParamConditionMap["param_value"].([]interface{})
				paramValue := make([]*string, len(paramValueInterface))
				for i, v := range paramValueInterface {
					str := v.(string)
					paramValue[i] = &str
				}
				uriParamCondition := &waapRatelimit.UriParamCondition{
					MatchType:  &matchType,
					ParamName:  &paramName,
					ParamValue: paramValue,
				}
				uriParamConditions = append(uriParamConditions, uriParamCondition)
			}
			conditionsRequest.UriParamConditions = uriParamConditions
		}

		// UA Conditions
		if conditionMap["ua_conditions"] != nil {
			uaConditions := make([]*waapRatelimit.UaCondition, 0)
			for _, uaCondition := range conditionMap["ua_conditions"].([]interface{}) {
				uaConditionMap := uaCondition.(map[string]interface{})
				matchType := uaConditionMap["match_type"].(string)
				uaInterface := uaConditionMap["ua"].([]interface{})
				ua := make([]*string, len(uaInterface))
				for i, v := range uaInterface {
					str := v.(string)
					ua[i] = &str
				}
				uaCondition := &waapRatelimit.UaCondition{
					MatchType: &matchType,
					Ua:        ua,
				}
				uaConditions = append(uaConditions, uaCondition)
			}
			conditionsRequest.UaConditions = uaConditions
		}

		// Referer Conditions
		if conditionMap["referer_conditions"] != nil {
			refererConditions := make([]*waapRatelimit.RefererCondition, 0)
			for _, refererCondition := range conditionMap["referer_conditions"].([]interface{}) {
				refererConditionMap := refererCondition.(map[string]interface{})
				matchType := refererConditionMap["match_type"].(string)
				refererInterface := refererConditionMap["referer"].([]interface{})
				referer := make([]*string, len(refererInterface))
				for i, v := range refererInterface {
					str := v.(string)
					referer[i] = &str
				}
				refererCondition := &waapRatelimit.RefererCondition{
					MatchType: &matchType,
					Referer:   referer,
				}
				refererConditions = append(refererConditions, refererCondition)
			}
			conditionsRequest.RefererConditions = refererConditions
		}

		// Header Conditions
		if conditionMap["header_conditions"] != nil {
			headerConditions := make([]*waapRatelimit.HeaderCondition, 0)
			for _, headerCondition := range conditionMap["header_conditions"].([]interface{}) {
				headerConditionMap := headerCondition.(map[string]interface{})
				matchType := headerConditionMap["match_type"].(string)
				key := headerConditionMap["key"].(string)
				valueListInterface := headerConditionMap["value_list"].([]interface{})
				valueList := make([]*string, len(valueListInterface))
				for i, v := range valueListInterface {
					str := v.(string)
					valueList[i] = &str
				}
				headerCondition := &waapRatelimit.HeaderCondition{
					MatchType: &matchType,
					Key:       &key,
					ValueList: valueList,
				}
				headerConditions = append(headerConditions, headerCondition)
			}
			conditionsRequest.HeaderConditions = headerConditions
		}

		// Area Conditions
		if conditionMap["area_conditions"] != nil {
			areaConditions := make([]*waapRatelimit.AreaCondition, 0)
			for _, areaCondition := range conditionMap["area_conditions"].([]interface{}) {
				areaConditionMap := areaCondition.(map[string]interface{})
				matchType := areaConditionMap["match_type"].(string)
				areasInterface := areaConditionMap["areas"].([]interface{})
				areas := make([]*string, len(areasInterface))
				for i, v := range areasInterface {
					str := v.(string)
					areas[i] = &str
				}
				areaCondition := &waapRatelimit.AreaCondition{
					MatchType: &matchType,
					Areas:     areas,
				}
				areaConditions = append(areaConditions, areaCondition)
			}
			conditionsRequest.AreaConditions = areaConditions
		}

		// Method Conditions
		if conditionMap["method_conditions"] != nil {
			methodConditions := make([]*waapRatelimit.RequestMethodCondition, 0)
			for _, methodCondition := range conditionMap["method_conditions"].([]interface{}) {
				methodConditionMap := methodCondition.(map[string]interface{})
				matchType := methodConditionMap["match_type"].(string)
				requestMethodInterface := methodConditionMap["request_method"].([]interface{})
				requestMethod := make([]*string, len(requestMethodInterface))
				for i, v := range requestMethodInterface {
					str := v.(string)
					requestMethod[i] = &str
				}
				methodCondition := &waapRatelimit.RequestMethodCondition{
					MatchType:     &matchType,
					RequestMethod: requestMethod,
				}
				methodConditions = append(methodConditions, methodCondition)
			}
			conditionsRequest.MethodConditions = methodConditions
		}

		// Status Code Conditions
		if conditionMap["status_code_conditions"] != nil {
			statusCodeConditions := make([]*waapRatelimit.StatusCodeCondition, 0)
			for _, statusCodeCondition := range conditionMap["status_code_conditions"].([]interface{}) {
				statusCodeConditionMap := statusCodeCondition.(map[string]interface{})
				matchType := statusCodeConditionMap["match_type"].(string)
				statusCodeInterface := statusCodeConditionMap["status_code"].([]interface{})
				statusCode := make([]*string, len(statusCodeInterface))
				for i, v := range statusCodeInterface {
					str := v.(string)
					statusCode[i] = &str
				}
				statusCodeCondition := &waapRatelimit.StatusCodeCondition{
					MatchType:  &matchType,
					StatusCode: statusCode,
				}
				statusCodeConditions = append(statusCodeConditions, statusCodeCondition)
			}
			conditionsRequest.StatusCodeConditions = statusCodeConditions
		}

		// Scheme Conditions
		if conditionMap["scheme_conditions"] != nil {
			schemeConditions := make([]*waapRatelimit.SchemeCondition, 0)
			for _, schemeCondition := range conditionMap["scheme_conditions"].([]interface{}) {
				schemeConditionMap := schemeCondition.(map[string]interface{})
				matchType := schemeConditionMap["match_type"].(string)
				schemes := make([]*string, len(schemeConditionMap["scheme"].([]interface{})))
				for i, scheme := range schemeConditionMap["scheme"].([]interface{}) {
					schemeStr := scheme.(string)
					schemes[i] = &schemeStr
				}
				schemeConditionRequest := &waapRatelimit.SchemeCondition{
					MatchType: &matchType,
					Scheme:    schemes,
				}
				schemeConditions = append(schemeConditions, schemeConditionRequest)
			}
			conditionsRequest.SchemeConditions = schemeConditions
		}

		// JA3 Conditions
		if conditionMap["ja3_conditions"] != nil {
			ja3Conditions := make([]*waapRatelimit.Ja3Condition, 0)
			for _, ja3Condition := range conditionMap["ja3_conditions"].([]interface{}) {
				ja3ConditionMap := ja3Condition.(map[string]interface{})
				matchType := ja3ConditionMap["match_type"].(string)
				ja3Interface := ja3ConditionMap["ja3_list"].([]interface{})
				ja3 := make([]*string, len(ja3Interface))
				for i, v := range ja3Interface {
					str := v.(string)
					ja3[i] = &str
				}
				ja3Condition := &waapRatelimit.Ja3Condition{
					MatchType: &matchType,
					Ja3List:   ja3,
				}
				ja3Conditions = append(ja3Conditions, ja3Condition)
			}
			conditionsRequest.Ja3Conditions = ja3Conditions
		}

		// JA4 Conditions
		if conditionMap["ja4_conditions"] != nil {
			ja4Conditions := make([]*waapRatelimit.Ja4Condition, 0)
			for _, ja4Condition := range conditionMap["ja4_conditions"].([]interface{}) {
				ja4ConditionMap := ja4Condition.(map[string]interface{})
				matchType := ja4ConditionMap["match_type"].(string)
				ja4Interface := ja4ConditionMap["ja4_list"].([]interface{})
				ja4 := make([]*string, len(ja4Interface))
				for i, v := range ja4Interface {
					str := v.(string)
					ja4[i] = &str
				}
				ja4Condition := &waapRatelimit.Ja4Condition{
					MatchType: &matchType,
					Ja4List:   ja4,
				}
				ja4Conditions = append(ja4Conditions, ja4Condition)
			}
			conditionsRequest.Ja4Conditions = ja4Conditions
		}
	}
	request.RateLimitRuleCondition = conditionsRequest

	var response *waapRatelimit.CreatRateLimitingRuleResponse
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseWaapRatelimitClient().AddRateLimit(request)
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
	_ = data.Set("id", *response.Data)
	data.SetId(*response.Data)
	//set status
	return resourceWaapRateLimitRead(context, data, meta)
}

func resourceWaapRateLimitRead(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_waap_ratelimit.read")
	var response *waapRatelimit.ListRateLimitingRulesResponse
	var err error
	var diags diag.Diagnostics
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		domain := data.Get("domain").(string)
		// 规则名称会变，不当成查询条件
		//ruleName := data.Get("rule_name").(string)

		request := &waapRatelimit.ListRateLimitingRulesRequest{
			DomainList: []*string{&domain},
			//RuleName:   &ruleName,
		}
		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseWaapRatelimitClient().GetRateLimitList(request)
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
		data.SetId("")
		return nil
	}
	if response.Data != nil {
		for _, item := range response.Data {
			// 只要对应id的数据
			if *item.Id != data.Id() {
				continue
			}
			_ = data.Set("domain", item.Domain)
			_ = data.Set("rule_name", item.RuleName)
			_ = data.Set("description", item.Description)
			_ = data.Set("scene", item.Scene)
			_ = data.Set("statistical_item", item.StatisticalItem)
			_ = data.Set("statistics_key", item.StatisticsKey)
			_ = data.Set("statistical_period", item.StatisticalPeriod)
			_ = data.Set("trigger_threshold", item.TriggerThreshold)
			_ = data.Set("intercept_time", item.InterceptTime)
			_ = data.Set("effective_status", item.EffectiveStatus)
			_ = data.Set("asset_api_id", item.AssetApiId)
			_ = data.Set("action", item.Action)

			if item.RateLimitEffective != nil {
				rateLimitEffective := map[string]interface{}{
					"effective": item.RateLimitEffective.Effective,
					"start":     item.RateLimitEffective.Start,
					"end":       item.RateLimitEffective.End,
					"timezone":  item.RateLimitEffective.Timezone,
				}
				_ = data.Set("rate_limit_effective", []interface{}{rateLimitEffective})
			}

			rateLimitRuleCondition := make([]interface{}, 0)
			if item.RateLimitRuleCondition != nil {
				condition := make(map[string]interface{})

				if item.RateLimitRuleCondition.IpOrIpsConditions != nil {
					ipOrIpsConditions := make([]interface{}, 0)
					for _, v := range item.RateLimitRuleCondition.IpOrIpsConditions {
						ipOrIpsCondition := map[string]interface{}{
							"match_type": v.MatchType,
							"ip_or_ips":  v.IpOrIps,
						}
						ipOrIpsConditions = append(ipOrIpsConditions, ipOrIpsCondition)
					}
					condition["ip_or_ips_conditions"] = ipOrIpsConditions
				}

				if item.RateLimitRuleCondition.PathConditions != nil {
					pathConditions := make([]interface{}, 0)
					for _, v := range item.RateLimitRuleCondition.PathConditions {
						pathCondition := map[string]interface{}{
							"match_type": v.MatchType,
							"paths":      v.Paths,
						}
						pathConditions = append(pathConditions, pathCondition)
					}
					condition["path_conditions"] = pathConditions
				}

				if item.RateLimitRuleCondition.UriConditions != nil {
					uriConditions := make([]interface{}, 0)
					for _, v := range item.RateLimitRuleCondition.UriConditions {
						uriCondition := map[string]interface{}{
							"match_type": v.MatchType,
							"uri":        v.Uri,
						}
						uriConditions = append(uriConditions, uriCondition)
					}
					condition["uri_conditions"] = uriConditions
				}

				if item.RateLimitRuleCondition.UriParamConditions != nil {
					uriParamConditions := make([]interface{}, 0)
					for _, v := range item.RateLimitRuleCondition.UriParamConditions {
						uriParamCondition := map[string]interface{}{
							"match_type":  v.MatchType,
							"param_name":  v.ParamName,
							"param_value": v.ParamValue,
						}
						uriParamConditions = append(uriParamConditions, uriParamCondition)
					}
					condition["uri_param_conditions"] = uriParamConditions
				}

				if item.RateLimitRuleCondition.UaConditions != nil {
					uaConditions := make([]interface{}, 0)
					for _, v := range item.RateLimitRuleCondition.UaConditions {
						uaCondition := map[string]interface{}{
							"match_type": v.MatchType,
							"ua":         v.Ua,
						}
						uaConditions = append(uaConditions, uaCondition)
					}
					condition["ua_conditions"] = uaConditions
				}

				if item.RateLimitRuleCondition.RefererConditions != nil {
					refererConditions := make([]interface{}, 0)
					for _, v := range item.RateLimitRuleCondition.RefererConditions {
						refererCondition := map[string]interface{}{
							"match_type": v.MatchType,
							"referer":    v.Referer,
						}
						refererConditions = append(refererConditions, refererCondition)
					}
					condition["referer_conditions"] = refererConditions
				}

				if item.RateLimitRuleCondition.HeaderConditions != nil {
					headerConditions := make([]interface{}, 0)
					for _, v := range item.RateLimitRuleCondition.HeaderConditions {
						headerCondition := map[string]interface{}{
							"match_type": v.MatchType,
							"key":        v.Key,
							"value_list": v.ValueList,
						}
						headerConditions = append(headerConditions, headerCondition)
					}
					condition["header_conditions"] = headerConditions
				}

				if item.RateLimitRuleCondition.AreaConditions != nil {
					areaConditions := make([]interface{}, 0)
					for _, v := range item.RateLimitRuleCondition.AreaConditions {
						areaCondition := map[string]interface{}{
							"match_type": v.MatchType,
							"areas":      v.Areas,
						}
						areaConditions = append(areaConditions, areaCondition)
					}
					condition["area_conditions"] = areaConditions
				}

				if item.RateLimitRuleCondition.StatusCodeConditions != nil {
					statusCodeConditions := make([]interface{}, 0)
					for _, v := range item.RateLimitRuleCondition.StatusCodeConditions {
						statusCodeCondition := map[string]interface{}{
							"match_type":  v.MatchType,
							"status_code": v.StatusCode,
						}
						statusCodeConditions = append(statusCodeConditions, statusCodeCondition)
					}
					condition["status_code_conditions"] = statusCodeConditions
				}

				if item.RateLimitRuleCondition.MethodConditions != nil {
					methodConditions := make([]interface{}, 0)
					for _, v := range item.RateLimitRuleCondition.MethodConditions {
						methodCondition := map[string]interface{}{
							"match_type":     v.MatchType,
							"request_method": v.RequestMethod,
						}
						methodConditions = append(methodConditions, methodCondition)
					}
					condition["method_conditions"] = methodConditions
				}
				// Scheme Conditions
				if item.RateLimitRuleCondition.SchemeConditions != nil {
					schemeConditions := make([]interface{}, 0)
					for _, v := range item.RateLimitRuleCondition.SchemeConditions {
						schemeCondition := map[string]interface{}{
							"match_type": v.MatchType,
							"scheme":     v.Scheme,
						}
						schemeConditions = append(schemeConditions, schemeCondition)
					}
					condition["scheme_conditions"] = schemeConditions
				}

				if item.RateLimitRuleCondition.Ja3Conditions != nil {
					ja3Conditions := make([]interface{}, 0)
					for _, condition := range item.RateLimitRuleCondition.Ja3Conditions {
						ja3Condition := map[string]interface{}{
							"match_type": condition.MatchType,
							"ja3_list":   condition.Ja3List,
						}
						ja3Conditions = append(ja3Conditions, ja3Condition)
					}
					condition["ja3_conditions"] = ja3Conditions
				}
				if item.RateLimitRuleCondition.Ja4Conditions != nil {
					ja4Conditions := make([]interface{}, 0)
					for _, condition := range item.RateLimitRuleCondition.Ja4Conditions {
						ja4Condition := map[string]interface{}{
							"match_type": condition.MatchType,
							"ja4_list":   condition.Ja4List,
						}
						ja4Conditions = append(ja4Conditions, ja4Condition)
					}
					condition["ja4_conditions"] = ja4Conditions
				}
				rateLimitRuleCondition = append(rateLimitRuleCondition, condition)
			}
			_ = data.Set("rate_limit_rule_condition", rateLimitRuleCondition)
		}
	}
	return nil
}

func resourceWaapRateLimitUpdate(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_waap_ratelimit.update")
	var diags diag.Diagnostics
	if data.HasChange("domain") {
		// 把domain强制刷回旧值，否则会有权限问题
		oldDomain, _ := data.GetChange("domain")
		_ = data.Set("domain", oldDomain)
		err := errors.New("Hostname cannot be changed.")
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}
	request := &waapRatelimit.UpdateRateLimitingRuleRequest{}
	if v, ok := data.GetOk("id"); ok {
		request.SetId(v.(string))
	}

	if v, ok := data.GetOk("rule_name"); ok {
		request.SetRuleName(v.(string))
	}
	if v, ok := data.GetOk("description"); ok {
		request.SetDescription(v.(string))
	}
	if v, ok := data.GetOk("scene"); ok {
		request.SetScene(v.(string))
	}
	if v, ok := data.GetOk("statistical_item"); ok {
		request.SetStatisticalItem(v.(string))
	}
	if v, ok := data.GetOk("statistics_key"); ok {
		request.SetStatisticsKey(v.(string))
	}
	if v, ok := data.GetOk("statistical_period"); ok {
		request.SetStatisticalPeriod(v.(int))
	}
	if v, ok := data.GetOk("trigger_threshold"); ok {
		request.SetTriggerThreshold(v.(int))
	}
	if v, ok := data.GetOk("intercept_time"); ok {
		request.SetInterceptTime(v.(int))
	}
	if v, ok := data.GetOk("effective_status"); ok {
		request.SetEffectiveStatus(v.(string))
	}
	if v, ok := data.GetOk("rate_limit_effective"); ok {
		rateLimitEffectiveV := v.([]interface{})
		for _, v := range rateLimitEffectiveV {
			rateLimitEffectiveD := v.(map[string]interface{})
			effective := rateLimitEffectiveD["effective"].([]interface{})
			start := rateLimitEffectiveD["start"].(string)
			end := rateLimitEffectiveD["end"].(string)
			timezone := rateLimitEffectiveD["timezone"].(string)
			ratelimitEffective := &waapRatelimit.RateLimitEffective{}
			effectives := make([]*string, len(effective))
			for i, v := range effective {
				str := v.(string)
				effectives[i] = &str
			}
			ratelimitEffective.SetEffective(effectives)
			ratelimitEffective.SetStart(start)
			ratelimitEffective.SetEnd(end)
			ratelimitEffective.SetTimezone(timezone)

			request.SetRateLimitEffective(ratelimitEffective)
		}
	}
	if v, ok := data.GetOk("asset_api_id"); ok {
		request.SetAssetApiId(v.(string))
	}
	if v, ok := data.GetOk("action"); ok {
		request.SetAction(v.(string))
	}
	conditions := data.Get("rate_limit_rule_condition").([]interface{})
	conditionsRequest := &waapRatelimit.RateLimitRuleCondition{}
	for _, v := range conditions {
		conditionMap := v.(map[string]interface{})
		// IpOrIps Conditions
		if conditionMap["ip_or_ips_conditions"] != nil {
			ipOrIpsConditions := make([]*waapRatelimit.IpOrIpsCondition, 0)
			for _, ipOrIpsCondition := range conditionMap["ip_or_ips_conditions"].([]interface{}) {
				ipOrIpsConditionMap := ipOrIpsCondition.(map[string]interface{})
				matchType := ipOrIpsConditionMap["match_type"].(string)
				ipOrIpsInterface := ipOrIpsConditionMap["ip_or_ips"].([]interface{})
				ipOrIps := make([]*string, len(ipOrIpsInterface))
				for i, v := range ipOrIpsInterface {
					str := v.(string)
					ipOrIps[i] = &str
				}
				ipOrIpsCondition := &waapRatelimit.IpOrIpsCondition{
					MatchType: &matchType,
					IpOrIps:   ipOrIps,
				}
				ipOrIpsConditions = append(ipOrIpsConditions, ipOrIpsCondition)
			}
			conditionsRequest.IpOrIpsConditions = ipOrIpsConditions
		}

		// Path Conditions
		if conditionMap["path_conditions"] != nil {
			pathConditions := make([]*waapRatelimit.PathCondition, 0)
			for _, pathCondition := range conditionMap["path_conditions"].([]interface{}) {
				pathConditionMap := pathCondition.(map[string]interface{})
				matchType := pathConditionMap["match_type"].(string)
				pathsInterface := pathConditionMap["paths"].([]interface{})
				paths := make([]*string, len(pathsInterface))
				for i, v := range pathsInterface {
					str := v.(string)
					paths[i] = &str
				}
				pathCondition := &waapRatelimit.PathCondition{
					MatchType: &matchType,
					Paths:     paths,
				}
				pathConditions = append(pathConditions, pathCondition)
			}
			conditionsRequest.PathConditions = pathConditions
		}

		// URI Conditions
		if conditionMap["uri_conditions"] != nil {
			uriConditions := make([]*waapRatelimit.UriCondition, 0)
			for _, uriCondition := range conditionMap["uri_conditions"].([]interface{}) {
				uriConditionMap := uriCondition.(map[string]interface{})
				matchType := uriConditionMap["match_type"].(string)
				uriInterface := uriConditionMap["uri"].([]interface{})
				uri := make([]*string, len(uriInterface))
				for i, v := range uriInterface {
					str := v.(string)
					uri[i] = &str
				}
				uriCondition := &waapRatelimit.UriCondition{
					MatchType: &matchType,
					Uri:       uri,
				}
				uriConditions = append(uriConditions, uriCondition)
			}
			conditionsRequest.SetUriConditions(uriConditions)
		}

		// URI Param Conditions
		if conditionMap["uri_param_conditions"] != nil {
			uriParamConditions := make([]*waapRatelimit.UriParamCondition, 0)
			for _, uriParamCondition := range conditionMap["uri_param_conditions"].([]interface{}) {
				uriParamConditionMap := uriParamCondition.(map[string]interface{})
				matchType := uriParamConditionMap["match_type"].(string)
				paramName := uriParamConditionMap["param_name"].(string)
				paramValueInterface := uriParamConditionMap["param_value"].([]interface{})
				paramValue := make([]*string, len(paramValueInterface))
				for i, v := range paramValueInterface {
					str := v.(string)
					paramValue[i] = &str
				}
				uriParamCondition := &waapRatelimit.UriParamCondition{
					MatchType:  &matchType,
					ParamName:  &paramName,
					ParamValue: paramValue,
				}
				uriParamConditions = append(uriParamConditions, uriParamCondition)
			}
			conditionsRequest.UriParamConditions = uriParamConditions
		}

		// UA Conditions
		if conditionMap["ua_conditions"] != nil {
			uaConditions := make([]*waapRatelimit.UaCondition, 0)
			for _, uaCondition := range conditionMap["ua_conditions"].([]interface{}) {
				uaConditionMap := uaCondition.(map[string]interface{})
				matchType := uaConditionMap["match_type"].(string)
				uaInterface := uaConditionMap["ua"].([]interface{})
				ua := make([]*string, len(uaInterface))
				for i, v := range uaInterface {
					str := v.(string)
					ua[i] = &str
				}
				uaCondition := &waapRatelimit.UaCondition{
					MatchType: &matchType,
					Ua:        ua,
				}
				uaConditions = append(uaConditions, uaCondition)
			}
			conditionsRequest.UaConditions = uaConditions
		}

		// Referer Conditions
		if conditionMap["referer_conditions"] != nil {
			refererConditions := make([]*waapRatelimit.RefererCondition, 0)
			for _, refererCondition := range conditionMap["referer_conditions"].([]interface{}) {
				refererConditionMap := refererCondition.(map[string]interface{})
				matchType := refererConditionMap["match_type"].(string)
				refererInterface := refererConditionMap["referer"].([]interface{})
				referer := make([]*string, len(refererInterface))
				for i, v := range refererInterface {
					str := v.(string)
					referer[i] = &str
				}
				refererCondition := &waapRatelimit.RefererCondition{
					MatchType: &matchType,
					Referer:   referer,
				}
				refererConditions = append(refererConditions, refererCondition)
			}
			conditionsRequest.RefererConditions = refererConditions
		}

		// Header Conditions
		if conditionMap["header_conditions"] != nil {
			headerConditions := make([]*waapRatelimit.HeaderCondition, 0)
			for _, headerCondition := range conditionMap["header_conditions"].([]interface{}) {
				headerConditionMap := headerCondition.(map[string]interface{})
				matchType := headerConditionMap["match_type"].(string)
				key := headerConditionMap["key"].(string)
				valueListInterface := headerConditionMap["value_list"].([]interface{})
				valueList := make([]*string, len(valueListInterface))
				for i, v := range valueListInterface {
					str := v.(string)
					valueList[i] = &str
				}
				headerCondition := &waapRatelimit.HeaderCondition{
					MatchType: &matchType,
					Key:       &key,
					ValueList: valueList,
				}
				headerConditions = append(headerConditions, headerCondition)
			}
			conditionsRequest.HeaderConditions = headerConditions
		}

		// Area Conditions
		if conditionMap["area_conditions"] != nil {
			areaConditions := make([]*waapRatelimit.AreaCondition, 0)
			for _, areaCondition := range conditionMap["area_conditions"].([]interface{}) {
				areaConditionMap := areaCondition.(map[string]interface{})
				matchType := areaConditionMap["match_type"].(string)
				areasInterface := areaConditionMap["areas"].([]interface{})
				areas := make([]*string, len(areasInterface))
				for i, v := range areasInterface {
					str := v.(string)
					areas[i] = &str
				}
				areaCondition := &waapRatelimit.AreaCondition{
					MatchType: &matchType,
					Areas:     areas,
				}
				areaConditions = append(areaConditions, areaCondition)
			}
			conditionsRequest.AreaConditions = areaConditions
		}

		// Method Conditions
		if conditionMap["method_conditions"] != nil {
			methodConditions := make([]*waapRatelimit.RequestMethodCondition, 0)
			for _, methodCondition := range conditionMap["method_conditions"].([]interface{}) {
				methodConditionMap := methodCondition.(map[string]interface{})
				matchType := methodConditionMap["match_type"].(string)
				requestMethodInterface := methodConditionMap["request_method"].([]interface{})
				requestMethod := make([]*string, len(requestMethodInterface))
				for i, v := range requestMethodInterface {
					str := v.(string)
					requestMethod[i] = &str
				}
				methodCondition := &waapRatelimit.RequestMethodCondition{
					MatchType:     &matchType,
					RequestMethod: requestMethod,
				}
				methodConditions = append(methodConditions, methodCondition)
			}
			conditionsRequest.MethodConditions = methodConditions
		}

		// Status Code Conditions
		if conditionMap["status_code_conditions"] != nil {
			statusCodeConditions := make([]*waapRatelimit.StatusCodeCondition, 0)
			for _, statusCodeCondition := range conditionMap["status_code_conditions"].([]interface{}) {
				statusCodeConditionMap := statusCodeCondition.(map[string]interface{})
				matchType := statusCodeConditionMap["match_type"].(string)
				statusCodeInterface := statusCodeConditionMap["status_code"].([]interface{})
				statusCode := make([]*string, len(statusCodeInterface))
				for i, v := range statusCodeInterface {
					str := v.(string)
					statusCode[i] = &str
				}
				statusCodeCondition := &waapRatelimit.StatusCodeCondition{
					MatchType:  &matchType,
					StatusCode: statusCode,
				}
				statusCodeConditions = append(statusCodeConditions, statusCodeCondition)
			}
			conditionsRequest.StatusCodeConditions = statusCodeConditions
		}

		// Scheme Conditions
		if conditionMap["scheme_conditions"] != nil {
			schemeConditions := make([]*waapRatelimit.SchemeCondition, 0)
			for _, schemeCondition := range conditionMap["scheme_conditions"].([]interface{}) {
				schemeConditionMap := schemeCondition.(map[string]interface{})
				matchType := schemeConditionMap["match_type"].(string)
				schemes := make([]*string, len(schemeConditionMap["scheme"].([]interface{})))
				for i, scheme := range schemeConditionMap["scheme"].([]interface{}) {
					schemeStr := scheme.(string)
					schemes[i] = &schemeStr
				}
				schemeConditionRequest := &waapRatelimit.SchemeCondition{
					MatchType: &matchType,
					Scheme:    schemes,
				}
				schemeConditions = append(schemeConditions, schemeConditionRequest)
			}
			conditionsRequest.SchemeConditions = schemeConditions
		}

		// JA3 Conditions
		if conditionMap["ja3_conditions"] != nil {
			ja3Conditions := make([]*waapRatelimit.Ja3Condition, 0)
			for _, ja3Condition := range conditionMap["ja3_conditions"].([]interface{}) {
				ja3ConditionMap := ja3Condition.(map[string]interface{})
				matchType := ja3ConditionMap["match_type"].(string)
				ja3Interface := ja3ConditionMap["ja3_list"].([]interface{})
				ja3 := make([]*string, len(ja3Interface))
				for i, v := range ja3Interface {
					str := v.(string)
					ja3[i] = &str
				}
				ja3Condition := &waapRatelimit.Ja3Condition{
					MatchType: &matchType,
					Ja3List:   ja3,
				}
				ja3Conditions = append(ja3Conditions, ja3Condition)
			}
			conditionsRequest.Ja3Conditions = ja3Conditions
		}

		// JA4 Conditions
		if conditionMap["ja4_conditions"] != nil {
			ja4Conditions := make([]*waapRatelimit.Ja4Condition, 0)
			for _, ja4Condition := range conditionMap["ja4_conditions"].([]interface{}) {
				ja4ConditionMap := ja4Condition.(map[string]interface{})
				matchType := ja4ConditionMap["match_type"].(string)
				ja4Interface := ja4ConditionMap["ja4_list"].([]interface{})
				ja4 := make([]*string, len(ja4Interface))
				for i, v := range ja4Interface {
					str := v.(string)
					ja4[i] = &str
				}
				ja4Condition := &waapRatelimit.Ja4Condition{
					MatchType: &matchType,
					Ja4List:   ja4,
				}
				ja4Conditions = append(ja4Conditions, ja4Condition)
			}
			conditionsRequest.Ja4Conditions = ja4Conditions
		}
	}
	request.RateLimitRuleCondition = conditionsRequest

	var response *waapRatelimit.UpdateRateLimitingRuleResponse
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseWaapRatelimitClient().UpdateRateLimit(request)
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
	log.Printf("resource.wangsu_waap_datelimit.update success")
	return nil
}

func resourceWaapRateLimitDelete(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_waap_ratelimit.delete")

	var response *waapRatelimit.DeleteRateLimitingRulesResponse
	var err error
	var diags diag.Diagnostics
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		id := data.Id()
		request := &waapRatelimit.DeleteRateLimitingRulesRequest{
			Ids: []*string{&id},
		}
		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseWaapRatelimitClient().DeleteRateLimit(request)
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
	return nil
}
