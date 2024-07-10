package ratelimit

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wangsuCommon "github.com/wangsu-api/terraform-provider-wangsu/wangsu/common"
	waapRatelimit "github.com/wangsu-api/wangsu-sdk-go/wangsu/waap/ratelimit"
	"log"
	"time"
)

func DataSourceRateLimit() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRateLimitRead,
		Schema: map[string]*schema.Schema{
			"domain_list": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: "Hostname list.",
			},
			"rule_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Rule name, fuzzy query.",
			},
			"data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule ID.",
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Hostname.",
						},
						"rule_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule Name.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description.",
						},
						"scene": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Protected target.\nWEB:Website\nAPI:API",
						},
						"statistical_stage": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Count on.\nREQUEST:Request\nRESPONSE:Response",
						},
						"statistical_item": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Client identifier.\nIP:Client IP\nIP_UA:Client IP and User-Agent\nCOOKIE:Cookie\nIP_COOKIE:Client IP and Cookie\nHEADER:Request Header\nIP_HEADER:Client IP and Request Header",
						},
						"statistics_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Statistical key value .",
						},
						"statistical_period": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Statistics period, unit: seconds.",
						},
						"trigger_threshold": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Trigger threshold, unit: times.",
						},
						"intercept_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Action duration, unit: seconds.",
						},
						"effective_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cycle effective status.\nPERMANENT:All time\nWITHOUT:Excluded time\nWITHIN:Selected time",
						},
						"rate_limit_effective": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"effective": {
										Type:        schema.TypeList,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Computed:    true,
										Description: "Effective.\nMON:Monday\nTUE:Tuesday\nWED:Wednesday\nTHU:Thursday\nFRI:Friday\nSAT:Saturday\nSUN:Sunday",
									},
									"start": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Start time, format: HH:mm.",
									},
									"end": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "End time, format: HH:mm.",
									},
									"timezone": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Timezone,default value: GTM+8.",
									},
								},
							},
						},
						"asset_api_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API ID under API business, multiple separated by ; sign.",
						},
						"action": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Action.\nNO_USE:Not Used\nACCEPT:Skip\nLOG:Log\nCOOKIE:Cookie verification\nJS_CHECK:Javascript verification\nDELAY:Delay\nBLOCK:Deny\nRESET:Reset Connection\nCustom response ID:Custom response ID",
						},
						"rate_limit_rule_condition": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Match conditions.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip_or_ips_conditions": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"match_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Match type.\nEQUAL:Equal\nNOT_EQUAL:Does not equal",
												},
												"ip_or_ips": {
													Type:        schema.TypeList,
													Elem:        &schema.Schema{Type: schema.TypeString},
													Computed:    true,
													Description: "IP/CIDR.",
												},
											},
										},
									},
									"path_conditions": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"match_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Match type.\nEQUAL:Equals\nNOT_EQUAL:Does not equal\nCONTAIN:Contains\nNOT_CONTAIN:Does not contains\nREGEX:Regex match\nNOT_REGEX: regular does not match\nSTART_WITH: starts with\nEND_WITH: ends with\nWILDCARD: wildcard matches\nNOT_WILDCARD: wildcard does not match",
												},
												"paths": {
													Type:        schema.TypeList,
													Elem:        &schema.Schema{Type: schema.TypeString},
													Computed:    true,
													Description: "Path.",
												},
											},
										},
									},
									"uri_conditions": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"match_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Match type.\nEQUAL:Equals\nNOT_EQUAL:Does not equal\nCONTAIN:Contains\nNOT_CONTAIN:Does not contains\nREGEX:Regex match\nNOT_REGEX: regular does not match\nSTART_WITH: starts with\nEND_WITH: ends with\nWILDCARD: wildcard matches\nNOT_WILDCARD: wildcard does not match",
												},
												"uri": {
													Type:        schema.TypeList,
													Elem:        &schema.Schema{Type: schema.TypeString},
													Computed:    true,
													Description: "URI.",
												},
											},
										},
									},
									"uri_param_conditions": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"match_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Match type.\nEQUAL:Equals\nNOT_EQUAL:Does not equal\nCONTAIN:Contains\nNOT_CONTAIN:Does not contains\nREGEX:Regex match\nNONE:Empty or non-existent",
												},
												"param_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Param name.",
												},
												"param_value": {
													Type:        schema.TypeList,
													Elem:        &schema.Schema{Type: schema.TypeString},
													Computed:    true,
													Description: "Param value.",
												},
											},
										},
									},
									"ua_conditions": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "User Agent, match type cannot be repeated.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"match_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Match type.\nEQUAL:Equals,ua case sensitive\nNOT_EQUAL:Does not equal,ua case sensitive\nCONTAIN:Contains,ua case insensitive\nNOT_CONTAIN:Does not contains,ua case insensitive\nREGEX:Regex match,ua case insensitive\nNONE:Empty or non-existent",
												},
												"ua": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "User agent.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
									"referer_conditions": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Referer, match type cannot be repeated.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"match_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Match type.\nEQUAL:Equals,referer case sensitive\nNOT_EQUAL:Does not equal,referer case sensitive\nCONTAIN:Contains,referer case insensitive\nNOT_CONTAIN:Does not contains,referer case insensitive\nREGEX:Regex match,referer case insensitive\nNONE:Empty or non-existent",
												},
												"referer": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Referer.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
									"header_conditions": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Request Header, match type can be repeated.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"match_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Match type.\nEQUAL:Equals,header value case sensitive\nNOT_EQUAL:Does not equal,header value case sensitive\nCONTAIN:Contains,header value case insensitive\nNOT_CONTAIN:Does not contains,header value case insensitive\nREGEX:Regex match,header value case insensitive\nNONE:Empty or non-existent",
												},
												"key": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Header name.",
												},
												"value_list": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Header value.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
									"area_conditions": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Geo.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"match_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Match type.\nEQUAL:Equal\nNOT_EQUAL:Does not equal",
												},
												"areas": {
													Type:        schema.TypeList,
													Computed:    true,
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
										Computed:    true,
										Description: "Response Code.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"match_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Match type.\nEQUAL:Equal\nNOT_EQUAL:Does not equal",
												},
												"status_code": {
													Type:        schema.TypeList,
													Computed:    true,
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
										Computed:    true,
										Description: "Request Method.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"match_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Match type.\nEQUAL:Equal\nNOT_EQUAL:Does not equal",
												},
												"request_method": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Request method.\nSupported values: GET/POST/DELETE/PUT/HEAD/OPTIONS/COPY.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
									"scheme_conditions": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "HTTP/S, match type cannot be repeated.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"match_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Match type.\nEQUAL:Equal\nNOT_EQUAL:Does not equal",
												},
												"scheme": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "HTTP/S.\nSupported values: HTTP/HTTPS.",
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
					},
				},
			},
		},
	}
}

func dataSourceRateLimitRead(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("data_source.wangsu_waap_ratelimit.read")

	var response *waapRatelimit.ListRateLimitingRulesResponse
	var err error
	var diags diag.Diagnostics
	request := &waapRatelimit.ListRateLimitingRulesRequest{}
	if v, ok := data.GetOk("rule_name"); ok {
		request.SetRuleName(v.(string))
	}
	if v, ok := data.GetOk("domain_list"); ok {
		domainsList := v.([]interface{})
		domainsStrList := make([]*string, len(domainsList))
		for i, v := range domainsList {
			str := v.(string)
			domainsStrList[i] = &str
		}
		request.SetDomainList(domainsStrList)
	}
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
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
	ids := make([]string, 0, len(response.Data))
	if response.Data != nil {
		itemList := make([]interface{}, len(response.Data))
		for i, item := range response.Data {
			conditionList := make([]map[string]interface{}, 1)
			condition := make(map[string]interface{})
			if item.RateLimitRuleCondition != nil {
				condition["ip_or_ips_conditions"] = flattenIpOrIpsConditions(item.RateLimitRuleCondition.IpOrIpsConditions)
				condition["path_conditions"] = flattenPathConditions(item.RateLimitRuleCondition.PathConditions)
				condition["uri_conditions"] = flattenUriConditions(item.RateLimitRuleCondition.UriConditions)
				condition["uri_param_conditions"] = flattenUriParamConditions(item.RateLimitRuleCondition.UriParamConditions)
				condition["ua_conditions"] = flattenUaConditions(item.RateLimitRuleCondition.UaConditions)
				condition["method_conditions"] = flattenMethodConditions(item.RateLimitRuleCondition.MethodConditions)
				condition["referer_conditions"] = flattenRefererConditions(item.RateLimitRuleCondition.RefererConditions)
				condition["header_conditions"] = flattenHeaderConditions(item.RateLimitRuleCondition.HeaderConditions)
				condition["area_conditions"] = flattenAreaConditions(item.RateLimitRuleCondition.AreaConditions)
				condition["status_code_conditions"] = flattenStatusCodeConditions(item.RateLimitRuleCondition.StatusCodeConditions)
				condition["scheme_conditions"] = flattenSchemeConditions(item.RateLimitRuleCondition.SchemeConditions)
			}
			conditionList[0] = condition
			ids = append(ids, *item.Id)
			itemList[i] = map[string]interface{}{
				"id":                 item.Id,
				"domain":             item.Domain,
				"rule_name":          item.RuleName,
				"description":        item.Description,
				"scene":              item.Scene,
				"statistical_stage":  item.StatisticalStage,
				"statistical_item":   item.StatisticalItem,
				"statistics_key":     item.StatisticsKey,
				"statistical_period": item.StatisticalPeriod,
				"trigger_threshold":  item.TriggerThreshold,
				"intercept_time":     item.InterceptTime,
				"effective_status":   item.EffectiveStatus,
				"rate_limit_effective": []interface{}{
					map[string]interface{}{
						"effective": item.RateLimitEffective.Effective,
						"start":     item.RateLimitEffective.Start,
						"end":       item.RateLimitEffective.End,
						"timezone":  item.RateLimitEffective.Timezone,
					},
				},
				"asset_api_id":              item.AssetApiId,
				"action":                    item.Action,
				"rate_limit_rule_condition": conditionList,
			}
		}
		if err := data.Set("data", itemList); err != nil {
			return diag.FromErr(fmt.Errorf("error setting data for resource: %s", err))
		}
		data.SetId(wangsuCommon.DataResourceIdsHash(ids))
	}
	return diags
}

func flattenIpOrIpsConditions(conditions []*waapRatelimit.IpOrIpsCondition) []interface{} {
	result := make([]interface{}, 0)
	for _, condition := range conditions {
		result = append(result, map[string]interface{}{
			"match_type": condition.MatchType,
			"ip_or_ips":  condition.IpOrIps,
		})
	}
	return result
}

func flattenPathConditions(conditions []*waapRatelimit.PathCondition) []interface{} {
	result := make([]interface{}, 0)
	for _, condition := range conditions {
		result = append(result, map[string]interface{}{
			"match_type": condition.MatchType,
			"paths":      condition.Paths,
		})
	}
	return result
}

func flattenUriConditions(conditions []*waapRatelimit.UriCondition) []interface{} {
	result := make([]interface{}, 0)
	for _, condition := range conditions {
		result = append(result, map[string]interface{}{
			"match_type": condition.MatchType,
			"uri":        condition.Uri,
		})
	}
	return result
}

func flattenUriParamConditions(conditions []*waapRatelimit.UriParamCondition) []interface{} {
	result := make([]interface{}, 0)
	for _, condition := range conditions {
		result = append(result, map[string]interface{}{
			"match_type":  condition.MatchType,
			"param_name":  condition.ParamName,
			"param_value": condition.ParamValue,
		})
	}
	return result
}

func flattenUaConditions(conditions []*waapRatelimit.UaCondition) []interface{} {
	result := make([]interface{}, 0)
	for _, condition := range conditions {
		result = append(result, map[string]interface{}{
			"match_type": condition.MatchType,
			"ua":         condition.Ua,
		})
	}
	return result
}

func flattenRefererConditions(conditions []*waapRatelimit.RefererCondition) []interface{} {
	result := make([]interface{}, 0)
	for _, condition := range conditions {
		result = append(result, map[string]interface{}{
			"match_type": condition.MatchType,
			"referer":    condition.Referer,
		})
	}
	return result
}

func flattenHeaderConditions(conditions []*waapRatelimit.HeaderCondition) []interface{} {
	result := make([]interface{}, 0)
	for _, condition := range conditions {
		result = append(result, map[string]interface{}{
			"match_type": condition.MatchType,
			"key":        condition.Key,
			"value_list": condition.ValueList,
		})
	}
	return result
}

func flattenAreaConditions(conditions []*waapRatelimit.AreaCondition) []interface{} {
	result := make([]interface{}, 0)
	for _, condition := range conditions {
		result = append(result, map[string]interface{}{
			"match_type": condition.MatchType,
			"areas":      condition.Areas,
		})
	}
	return result
}

func flattenMethodConditions(conditions []*waapRatelimit.RequestMethodCondition) []interface{} {
	result := make([]interface{}, 0)
	for _, condition := range conditions {
		result = append(result, map[string]interface{}{
			"match_type":     condition.MatchType,
			"request_method": condition.RequestMethod,
		})
	}
	return result
}

func flattenStatusCodeConditions(conditions []*waapRatelimit.StatusCodeCondition) []interface{} {
	result := make([]interface{}, 0)
	for _, condition := range conditions {
		result = append(result, map[string]interface{}{
			"match_type":  condition.MatchType,
			"status_code": condition.StatusCode,
		})
	}
	return result
}

func flattenSchemeConditions(conditions []*waapRatelimit.SchemeCondition) []interface{} {
	result := make([]interface{}, 0)
	for _, condition := range conditions {
		result = append(result, map[string]interface{}{
			"match_type": condition.MatchType,
			"scheme":     condition.Scheme,
		})
	}
	return result
}
