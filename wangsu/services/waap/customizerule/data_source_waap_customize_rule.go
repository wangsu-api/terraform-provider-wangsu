package customizerule

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wangsuCommon "github.com/wangsu/terraform-provider-wangsu/wangsu/common"
	waapCustomizerule "github.com/wangsu/wangsu-sdk-go/wangsu/waap/customizerule"
	"log"
	"time"
)

func DataSourceCustomizeRule() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCustomizeRuleRead,
		Schema: map[string]*schema.Schema{
			"domain_list": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
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
							Description: "Rule name.",
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
						"api_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API ID, multiple separated by ; sign.",
						},
						"act": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Action.\nNO_USE:Not Used\nLOG:Log\nDELAY:Delay\nBLOCK:Deny\nRESET:Reset Connection",
						},
						"condition_list": {
							Type:     schema.TypeList,
							Computed: true,
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
													Computed:    true,
													Elem:        &schema.Schema{Type: schema.TypeString},
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
													Description: "Match type.\nEQUAL: equal to\nNOT_EQUAL: not equal to\nCONTAIN: contains\nNOT_CONTAIN: does not contain\nREGEX: regular\nNOT_REGEX: regular does not match\nSTART_WITH: starts with\nEND_WITH: ends with\nWILDCARD: wildcard matches\nNOT_WILDCARD: wildcard does not match",
												},
												"paths": {
													Type:        schema.TypeList,
													Computed:    true,
													Elem:        &schema.Schema{Type: schema.TypeString},
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
													Description: "Match type.\nEQUAL: equal to\nNOT_EQUAL: not equal to\nCONTAIN: contains\nNOT_CONTAIN: does not contain\nREGEX: regular\nNOT_REGEX: regular does not match\nSTART_WITH: starts with\nEND_WITH: ends with\nWILDCARD: wildcard matches\nNOT_WILDCARD: wildcard does not match",
												},
												"uri": {
													Type:        schema.TypeList,
													Computed:    true,
													Elem:        &schema.Schema{Type: schema.TypeString},
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
													Computed:    true,
													Elem:        &schema.Schema{Type: schema.TypeString},
													Description: "Param value.",
												},
											},
										},
									},
									"ua_conditions": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"match_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Match type.\nEQUAL: equal to\nNOT_EQUAL: not equal to\nCONTAIN: contains\nNOT_CONTAIN: does not contain\nREGEX: regular\nNONE: empty or does not exist\nNOT_REGEX: regular does not match\nSTART_WITH: starts with\nEND_WITH: ends with\nWILDCARD: wildcard matches\nNOT_WILDCARD: wildcard does not match",
												},
												"ua": {
													Type:        schema.TypeList,
													Computed:    true,
													Elem:        &schema.Schema{Type: schema.TypeString},
													Description: "User-Agent.",
												},
											},
										},
									},
									"referer_conditions": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"match_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Match type.\nEQUAL: equal to\nNOT_EQUAL: not equal to\nCONTAIN: contains\nNOT_CONTAIN: does not contain\nREGEX: regular\nNONE: empty or does not exist\nNOT_REGEX: regular does not match\nSTART_WITH: starts with\nEND_WITH: ends with\nWILDCARD: wildcard matches\nNOT_WILDCARD: wildcard does not match",
												},
												"referer": {
													Type:        schema.TypeList,
													Computed:    true,
													Elem:        &schema.Schema{Type: schema.TypeString},
													Description: "Referer.",
												},
											},
										},
									},
									"header_conditions": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"match_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Match type.\nEQUAL: equal to\nNOT_EQUAL: not equal to\nCONTAIN: contains\nNOT_CONTAIN: does not contain\nREGEX: regular\nNONE: empty or does not exist\nNOT_REGEX: regular does not match\nSTART_WITH: starts with\nEND_WITH: ends with\nWILDCARD: wildcard matches\nNOT_WILDCARD: wildcard does not match",
												},
												"key": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Request header name.",
												},
												"value_list": {
													Type:        schema.TypeList,
													Computed:    true,
													Elem:        &schema.Schema{Type: schema.TypeString},
													Description: "Header value.",
												},
											},
										},
									},
									"area_conditions": {
										Type:     schema.TypeList,
										Computed: true,
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
													Elem:        &schema.Schema{Type: schema.TypeString},
													Description: "Geo.",
												},
											},
										},
									},
									"method_conditions": {
										Type:     schema.TypeList,
										Computed: true,
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
													Elem:        &schema.Schema{Type: schema.TypeString},
													Description: "Request method.\nSupported values: GET/POST/DELETE/PUT/HEAD/OPTIONS/COPY.",
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

func dataSourceCustomizeRuleRead(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("data_source.wangsu_waap_customizerule.read")

	var response *waapCustomizerule.ListCustomRulesResponse
	var err error
	var diags diag.Diagnostics
	request := &waapCustomizerule.ListCustomRulesRequest{}
	if v, ok := data.GetOk("rule_name"); ok {
		request.SetRuleName(v.(string))
	}
	if v, ok := data.GetOk("domain_list"); ok {
		targetDomainsList := v.([]interface{})
		targetDomainsStr := make([]*string, len(targetDomainsList))
		for i, v := range targetDomainsList {
			str := v.(string)
			targetDomainsStr[i] = &str
		}
		request.SetDomainList(targetDomainsStr)
	}
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseWaapCustomizeruleClient().GetCustomRuleList(request)
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
		ids := make([]string, 0, len(response.Data))
		itemList := make([]interface{}, 0)
		for _, item := range response.Data {
			conditionList := make([]map[string]interface{}, 1)
			condition := make(map[string]interface{})
			if item.ConditionList != nil {
				condition["ip_or_ips_conditions"] = flattenIpOrIpsConditions(item.ConditionList.IpOrIpsConditions)
				condition["path_conditions"] = flattenPathConditions(item.ConditionList.PathConditions)
				condition["uri_conditions"] = flattenUriConditions(item.ConditionList.UriConditions)
				condition["uri_param_conditions"] = flattenUriParamConditions(item.ConditionList.UriParamConditions)
				condition["ua_conditions"] = flattenUaConditions(item.ConditionList.UaConditions)
				condition["referer_conditions"] = flattenRefererConditions(item.ConditionList.RefererConditions)
				condition["header_conditions"] = flattenHeaderConditions(item.ConditionList.HeaderConditions)
				condition["area_conditions"] = flattenAreaConditions(item.ConditionList.AreaConditions)
				condition["method_conditions"] = flattenMethodConditions(item.ConditionList.MethodConditions)
			}
			conditionList[0] = condition
			itemList = append(itemList, map[string]interface{}{
				"id":             item.Id,
				"domain":         item.Domain,
				"rule_name":      item.RuleName,
				"description":    item.Description,
				"scene":          item.Scene,
				"api_id":         item.ApiId,
				"act":            item.Act,
				"condition_list": conditionList,
			})
			ids = append(ids, *item.Id)
		}
		if err := data.Set("data", itemList); err != nil {
			return diag.FromErr(fmt.Errorf("error setting data for resource: %s", err))
		}
		data.SetId(wangsuCommon.DataResourceIdsHash(ids))
	}
	return diags
}

func flattenIpOrIpsConditions(conditions []*waapCustomizerule.IpOrIpsCondition) []interface{} {
	result := make([]interface{}, 0)
	for _, condition := range conditions {
		result = append(result, map[string]interface{}{
			"match_type": condition.MatchType,
			"ip_or_ips":  condition.IpOrIps,
		})
	}
	return result
}

func flattenPathConditions(conditions []*waapCustomizerule.PathCondition) []interface{} {
	result := make([]interface{}, 0)
	for _, condition := range conditions {
		result = append(result, map[string]interface{}{
			"match_type": condition.MatchType,
			"paths":      condition.Paths,
		})
	}
	return result
}

func flattenUriConditions(conditions []*waapCustomizerule.UriCondition) []interface{} {
	result := make([]interface{}, 0)
	for _, condition := range conditions {
		result = append(result, map[string]interface{}{
			"match_type": condition.MatchType,
			"uri":        condition.Uri,
		})
	}
	return result
}

func flattenUriParamConditions(conditions []*waapCustomizerule.UriParamCondition) []interface{} {
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

func flattenUaConditions(conditions []*waapCustomizerule.UaCondition) []interface{} {
	result := make([]interface{}, 0)
	for _, condition := range conditions {
		result = append(result, map[string]interface{}{
			"match_type": condition.MatchType,
			"ua":         condition.Ua,
		})
	}
	return result
}

func flattenRefererConditions(conditions []*waapCustomizerule.RefererCondition) []interface{} {
	result := make([]interface{}, 0)
	for _, condition := range conditions {
		result = append(result, map[string]interface{}{
			"match_type": condition.MatchType,
			"referer":    condition.Referer,
		})
	}
	return result
}

func flattenHeaderConditions(conditions []*waapCustomizerule.HeaderCondition) []interface{} {
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

func flattenAreaConditions(conditions []*waapCustomizerule.AreaCondition) []interface{} {
	result := make([]interface{}, 0)
	for _, condition := range conditions {
		result = append(result, map[string]interface{}{
			"match_type": condition.MatchType,
			"areas":      condition.Areas,
		})
	}
	return result
}

func flattenMethodConditions(conditions []*waapCustomizerule.RequestMethodCondition) []interface{} {
	result := make([]interface{}, 0)
	for _, condition := range conditions {
		result = append(result, map[string]interface{}{
			"match_type":     condition.MatchType,
			"request_method": condition.RequestMethod,
		})
	}
	return result
}
