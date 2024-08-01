package whitelist

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wangsuCommon "github.com/wangsu-api/terraform-provider-wangsu/wangsu/common"
	waapWhitelist "github.com/wangsu-api/wangsu-sdk-go/wangsu/waap/whitelist"
	"log"
	"time"
)

func DataSourceWaapWhitelist() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWaapWhitelistRead,

		Schema: map[string]*schema.Schema{
			"domain_list": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Hostname list.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"rule_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Rule name, fuzzy query.",
			},
			"data": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Data.",
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
						"conditions": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Match conditions, at least one, at most five.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip_or_ips_conditions": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "IP/CIDR match conditions.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"match_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Match type.<br/>EQUAL: Equals<br/>NOT_EQUAL: Does not equal",
												},
												"ip_or_ips": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "IP/CIDR.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
									"path_conditions": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Path match conditions.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"match_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Match type.<br/>EQUAL: Equals, path case sensitive<br/>NOT_EQUAL: Does not equal, path case sensitive<br/>CONTAIN: Contains, path case insensitive<br/>NOT_CONTAIN: Does not Contains, path case insensitive<br/>REGEX: Regex match, path case insensitive<br/>NOT_REGEX: Regular does not match, path case sensitive<br/>START_WITH: Starts with, path case sensitive<br/>END_WITH: Ends with, path case sensitive<br/>WILDCARD: Wildcard matches, path case sensitive, * represents zero or more arbitrary characters, ? represents any single character.<br/>NOT_WILDCARD: Wildcard does not match, path case sensitive, * represents zero or more arbitrary characters, ? represents any single character",
												},
												"paths": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Path.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
									"uri_conditions": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "URI match conditions.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"match_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Match type.<br/>EQUAL: Equals, URI case sensitive<br/>NOT_EQUAL: Does not equal, URI case sensitive<br/>CONTAIN: Contains, URI case insensitive<br/>NOT_CONTAIN: Does not Contains, URI case insensitive<br/>REGEX: Regex match, URI case insensitive<br/>NOT_REGEX: Regular does not match, URI case insensitive<br/>START_WITH: Starts with, URI case insensitive<br/>END_WITH: Ends with, URI case insensitive<br/>WILDCARD: Wildcard matches, URI case insensitive, * represents zero or more arbitrary characters, ? represents any single character<br/>NOT_WILDCARD: Wildcard does not match, URI case insensitive, * represents zero or more arbitrary characters, ? represents any single character",
												},
												"uri": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "URI.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
									"ua_conditions": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "User agent match conditions.",
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
										Description: "Referer match conditions.",
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
										Optional:    true,
										Description: "Request header match conditions.",
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
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceWaapWhitelistRead(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("data_source.wangsu_waap_whitelist.read")

	var response *waapWhitelist.ListWhitelistRulesResponse
	var err error
	var diags diag.Diagnostics
	request := &waapWhitelist.ListWhitelistRulesRequest{}
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
		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseWaapWhitelistClient().GetWaapWhitelistList(request)
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
			if item.Conditions != nil {
				conditions := make([]map[string]interface{}, 1)
				condition := make(map[string]interface{})
				if item.Conditions.IpOrIpsConditions != nil {
					ipOrIpsConditions := make([]interface{}, 0)
					for _, condition := range item.Conditions.IpOrIpsConditions {
						ipOrIpsCondition := map[string]interface{}{
							"match_type": condition.MatchType,
							"ip_or_ips":  condition.IpOrIps,
						}
						ipOrIpsConditions = append(ipOrIpsConditions, ipOrIpsCondition)
					}
					condition["ip_or_ips_conditions"] = ipOrIpsConditions
				}
				if item.Conditions.PathConditions != nil {
					pathConditions := make([]interface{}, 0)
					for _, condition := range item.Conditions.PathConditions {
						pathCondition := map[string]interface{}{
							"match_type": condition.MatchType,
							"paths":      condition.Paths,
						}
						pathConditions = append(pathConditions, pathCondition)
					}
					condition["path_conditions"] = pathConditions
				}
				if item.Conditions.UriConditions != nil {
					uriConditions := make([]interface{}, 0)
					for _, condition := range item.Conditions.UriConditions {
						uriCondition := map[string]interface{}{
							"match_type": condition.MatchType,
							"uri":        condition.Uri,
						}
						uriConditions = append(uriConditions, uriCondition)
					}
					condition["uri_conditions"] = uriConditions
				}
				if item.Conditions.UaConditions != nil {
					uaConditions := make([]interface{}, 0)
					for _, condition := range item.Conditions.UaConditions {
						uaCondition := map[string]interface{}{
							"match_type": condition.MatchType,
							"ua":         condition.Ua,
						}
						uaConditions = append(uaConditions, uaCondition)
					}
					condition["ua_conditions"] = uaConditions
				}
				if item.Conditions.RefererConditions != nil {
					refererConditions := make([]interface{}, 0)
					for _, condition := range item.Conditions.RefererConditions {
						refererCondition := map[string]interface{}{
							"match_type": condition.MatchType,
							"referer":    condition.Referer,
						}
						refererConditions = append(refererConditions, refererCondition)
					}
					condition["referer_conditions"] = refererConditions
				}
				if item.Conditions.HeaderConditions != nil {
					headerConditions := make([]interface{}, 0)
					for _, condition := range item.Conditions.HeaderConditions {
						headerCondition := map[string]interface{}{
							"match_type": condition.MatchType,
							"key":        condition.Key,
							"value_list": condition.ValueList,
						}
						headerConditions = append(headerConditions, headerCondition)
					}
					condition["header_conditions"] = headerConditions
				}
				conditions[0] = condition
				itemList[i] = map[string]interface{}{
					"id":          item.Id,
					"domain":      item.Domain,
					"rule_name":   item.RuleName,
					"description": item.Description,
					"conditions":  conditions,
				}
				ids = append(ids, *item.Id)
			}
		}
		if err := data.Set("data", itemList); err != nil {
			return diag.FromErr(fmt.Errorf("error setting data for resource: %s", err))
		}
	}
	data.SetId(wangsuCommon.DataResourceIdsHash(ids))
	return diags
}
