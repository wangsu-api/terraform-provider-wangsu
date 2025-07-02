package share_customizerule

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wangsuCommon "github.com/wangsu-api/terraform-provider-wangsu/wangsu/common"
	waapShareCustomizerule "github.com/wangsu-api/wangsu-sdk-go/wangsu/waap/share-customizerule"
	"log"
	"time"
)

func DataSourceCustomizeRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCustomizeRulesRead,
		Schema: map[string]*schema.Schema{
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
						"relation_domain_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Associated hostname.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
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
						"act": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Action.<br/>NO_USE:Not Used<br/>LOG:Log<br/>DELAY:Delay<br/>BLOCK:Deny<br/>RESET:Reset Connection",
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
													Description: "Match type.<br/>EQUAL:Equal<br/>NOT_EQUAL:Does not equal",
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
													Description: "Match type.<br/>EQUAL: Equals, user agent case sensitive<br/>NOT_EQUAL: Does not equal, user agent case sensitive<br/>CONTAIN: Contains, user agent case insensitive<br/>NOT_CONTAIN: Does not Contains, user agent case insensitive<br/>REGEX: Regex match, user agent case insensitive<br/>NOT_REGEX: Regular does not match, user agent case insensitive<br/>START_WITH: Starts with, user agent case insensitive<br/>END_WITH: Ends with, user agent case insensitive<br/>WILDCARD: Wildcard matches, user agent case insensitive, * represents zero or more arbitrary characters, ? represents any single character<br/>NOT_WILDCARD: Wildcard does not match, user agent case insensitive, * represents zero or more arbitrary characters, ? represents any single character",
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
													Description: "Match type.<br/>EQUAL: Equals, user agent case sensitive<br/>NOT_EQUAL: Does not equal, user agent case sensitive<br/>CONTAIN: Contains, user agent case insensitive<br/>NOT_CONTAIN: Does not Contains, user agent case insensitive<br/>REGEX: Regex match, user agent case insensitive<br/>NOT_REGEX: Regular does not match, user agent case insensitive<br/>START_WITH: Starts with, user agent case insensitive<br/>END_WITH: Ends with, user agent case insensitive<br/>WILDCARD: Wildcard matches, user agent case insensitive, * represents zero or more arbitrary characters, ? represents any single character<br/>NOT_WILDCARD: Wildcard does not match, user agent case insensitive, * represents zero or more arbitrary characters, ? represents any single character",
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
									"ua_conditions": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"match_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Match type.<br/>EQUAL: Equals, user agent case sensitive<br/>NOT_EQUAL: Does not equal, user agent case sensitive<br/>CONTAIN: Contains, user agent case insensitive<br/>NOT_CONTAIN: Does not Contains, user agent case insensitive<br/>NONE:Empty or non-existent<br/>REGEX: Regex match, user agent case insensitive<br/>NOT_REGEX: Regular does not match, user agent case insensitive<br/>START_WITH: Starts with, user agent case insensitive<br/>END_WITH: Ends with, user agent case insensitive<br/>WILDCARD: Wildcard matches, user agent case insensitive, * represents zero or more arbitrary characters, ? represents any single character<br/>NOT_WILDCARD: Wildcard does not match, user agent case insensitive, * represents zero or more arbitrary characters, ? represents any single character",
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
													Description: "Match type.<br/>EQUAL: Equals, user agent case sensitive<br/>NOT_EQUAL: Does not equal, user agent case sensitive<br/>CONTAIN: Contains, user agent case insensitive<br/>NOT_CONTAIN: Does not Contains, user agent case insensitive<br/>NONE:Empty or non-existent<br/>REGEX: Regex match, user agent case insensitive<br/>NOT_REGEX: Regular does not match, user agent case insensitive<br/>START_WITH: Starts with, user agent case insensitive<br/>END_WITH: Ends with, user agent case insensitive<br/>WILDCARD: Wildcard matches, user agent case insensitive, * represents zero or more arbitrary characters, ? represents any single character<br/>NOT_WILDCARD: Wildcard does not match, user agent case insensitive, * represents zero or more arbitrary characters, ? represents any single character",
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
													Description: "Match type.<br/>EQUAL: Equals, user agent case sensitive<br/>NOT_EQUAL: Does not equal, user agent case sensitive<br/>CONTAIN: Contains, user agent case insensitive<br/>NOT_CONTAIN: Does not Contains, user agent case insensitive<br/>NONE:Empty or non-existent<br/>REGEX: Regex match, user agent case insensitive<br/>NOT_REGEX: Regular does not match, user agent case insensitive<br/>START_WITH: Starts with, user agent case insensitive<br/>END_WITH: Ends with, user agent case insensitive<br/>WILDCARD: Wildcard matches, user agent case insensitive, * represents zero or more arbitrary characters, ? represents any single character<br/>NOT_WILDCARD: Wildcard does not match, user agent case insensitive, * represents zero or more arbitrary characters, ? represents any single character",
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
													Description: "Match type.<br/>EQUAL:Equal<br/>NOT_EQUAL:Does not equal",
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
													Description: "Match type.<br/>EQUAL:Equal<br/>NOT_EQUAL:Does not equal",
												},
												"request_method": {
													Type:        schema.TypeList,
													Computed:    true,
													Elem:        &schema.Schema{Type: schema.TypeString},
													Description: "Request method.<br/>Supported values: GET/POST/DELETE/PUT/HEAD/OPTIONS/COPY.",
												},
											},
										},
									},
									"ja3_conditions": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"match_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Match type.\nEQUAL: Equals\nNOT_EQUAL: Does not equal",
												},
												"ja3_list": {
													Type:        schema.TypeList,
													Computed:    true,
													Elem:        &schema.Schema{Type: schema.TypeString},
													Description: "JA3 Fingerprint List.",
												},
											},
										},
									},
									"ja4_conditions": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"match_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Match type. \nEQUAL: Equals\nNOT_EQUAL: Does not equal\nCONTAIN: Contains\nNOT_CONTAIN: Does not Contains\nSTART_WITH: Starts with\nEND_WITH: Ends with\nWILDCARD: Wildcard matches, ** represents zero or more arbitrary characters, ? represents any single character\nNOT_WILDCARD: Wildcard does not match, ** represents zero or more arbitrary characters, ? represents any single character",
												},
												"ja4_list": {
													Type:        schema.TypeList,
													Computed:    true,
													Elem:        &schema.Schema{Type: schema.TypeString},
													Description: "JA4 Fingerprint List.",
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

func dataSourceCustomizeRulesRead(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("data_source.wangsu_waap_share_customizerule.read")

	var response *waapShareCustomizerule.ListSharedCustomRulesResponse
	var err error
	var diags diag.Diagnostics
	request := &waapShareCustomizerule.ListSharedCustomRulesRequest{}
	if v, ok := data.GetOk("rule_name"); ok {
		request.SetRuleName(v.(string))
	}
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseWaapShareCustomizeruleClient().GetList(request)
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
			if item.Condition != nil {
				condition["ip_or_ips_conditions"] = flattenIpOrIpsConditions(item.Condition.IpOrIpsConditions)
				condition["path_conditions"] = flattenPathConditions(item.Condition.PathConditions)
				condition["uri_conditions"] = flattenUriConditions(item.Condition.UriConditions)
				condition["ua_conditions"] = flattenUaConditions(item.Condition.UaConditions)
				condition["referer_conditions"] = flattenRefererConditions(item.Condition.RefererConditions)
				condition["header_conditions"] = flattenHeaderConditions(item.Condition.HeaderConditions)
				condition["area_conditions"] = flattenAreaConditions(item.Condition.AreaConditions)
				condition["method_conditions"] = flattenMethodConditions(item.Condition.MethodConditions)
				condition["ja3_conditions"] = flattenJa3Conditions(item.Condition.Ja3Conditions)
				condition["ja4_conditions"] = flattenJa4Conditions(item.Condition.Ja4Conditions)
			}
			conditionList[0] = condition
			itemList = append(itemList, map[string]interface{}{
				"id":                   item.Id,
				"relation_domain_list": item.RelationDomainList,
				"rule_name":            item.RuleName,
				"description":          item.Description,
				"act":                  item.Act,
				"condition_list":       conditionList,
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

func flattenIpOrIpsConditions(conditions []*waapShareCustomizerule.IpOrIpsCondition) []interface{} {
	result := make([]interface{}, 0)
	for _, condition := range conditions {
		result = append(result, map[string]interface{}{
			"match_type": condition.MatchType,
			"ip_or_ips":  condition.IpOrIps,
		})
	}
	return result
}

func flattenPathConditions(conditions []*waapShareCustomizerule.PathCondition) []interface{} {
	result := make([]interface{}, 0)
	for _, condition := range conditions {
		result = append(result, map[string]interface{}{
			"match_type": condition.MatchType,
			"paths":      condition.Paths,
		})
	}
	return result
}

func flattenUriConditions(conditions []*waapShareCustomizerule.UriCondition) []interface{} {
	result := make([]interface{}, 0)
	for _, condition := range conditions {
		result = append(result, map[string]interface{}{
			"match_type": condition.MatchType,
			"uri":        condition.Uri,
		})
	}
	return result
}

func flattenUaConditions(conditions []*waapShareCustomizerule.UaCondition) []interface{} {
	result := make([]interface{}, 0)
	for _, condition := range conditions {
		result = append(result, map[string]interface{}{
			"match_type": condition.MatchType,
			"ua":         condition.Ua,
		})
	}
	return result
}

func flattenRefererConditions(conditions []*waapShareCustomizerule.RefererCondition) []interface{} {
	result := make([]interface{}, 0)
	for _, condition := range conditions {
		result = append(result, map[string]interface{}{
			"match_type": condition.MatchType,
			"referer":    condition.Referer,
		})
	}
	return result
}

func flattenHeaderConditions(conditions []*waapShareCustomizerule.HeaderCondition) []interface{} {
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

func flattenAreaConditions(conditions []*waapShareCustomizerule.AreaCondition) []interface{} {
	result := make([]interface{}, 0)
	for _, condition := range conditions {
		result = append(result, map[string]interface{}{
			"match_type": condition.MatchType,
			"areas":      condition.Areas,
		})
	}
	return result
}

func flattenMethodConditions(conditions []*waapShareCustomizerule.MethodCondition) []interface{} {
	result := make([]interface{}, 0)
	for _, condition := range conditions {
		result = append(result, map[string]interface{}{
			"match_type":     condition.MatchType,
			"request_method": condition.RequestMethod,
		})
	}
	return result
}

func flattenJa3Conditions(conditions []*waapShareCustomizerule.Ja3Condition) []interface{} {
	result := make([]interface{}, 0)
	for _, condition := range conditions {
		result = append(result, map[string]interface{}{
			"match_type": condition.MatchType,
			"ja3_list":   condition.Ja3List,
		})
	}
	return result
}

func flattenJa4Conditions(conditions []*waapShareCustomizerule.Ja4Condition) []interface{} {
	result := make([]interface{}, 0)
	for _, condition := range conditions {
		result = append(result, map[string]interface{}{
			"match_type": condition.MatchType,
			"ja4_list":   condition.Ja4List,
		})
	}
	return result
}
