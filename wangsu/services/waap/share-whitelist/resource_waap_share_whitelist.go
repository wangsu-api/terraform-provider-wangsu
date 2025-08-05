package share_whitelist

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wangsuCommon "github.com/wangsu-api/terraform-provider-wangsu/wangsu/common"
	waapShareWhitelist "github.com/wangsu-api/wangsu-sdk-go/wangsu/waap/share-whitelist"
	"log"
	"time"
)

func ResourceWaapShareWhitelist() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWaapShareWhitelistCreate,
		ReadContext:   resourceWaapShareWhitelistRead,
		UpdateContext: resourceWaapShareWhitelistUpdate,
		DeleteContext: resourceWaapShareWhitelistDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Rule ID.",
			},
			"relation_domain_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Associated hostname.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"rule_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Rule name, maximum 50 characters.<br/> does not support # and & .",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description, maximum 200 characters.",
			},
			"conditions": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Match conditions, at least one, at most five.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_or_ips_conditions": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "IP/CIDR match conditions, match type cannot be repeated.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"match_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Match type.<br/>EQUAL: Equals<br/>NOT_EQUAL: Does not equal",
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
							Description: "Path match conditions, match type cannot be repeated.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"match_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Match type.<br/>EQUAL: Equals, path case sensitive<br/>NOT_EQUAL: Does not equal, path case sensitive<br/>CONTAIN: Contains, path case insensitive<br/>NOT_CONTAIN: Does not Contains, path case insensitive<br/>REGEX: Regex match, path case insensitive<br/>NOT_REGEX: Regular does not match, path case sensitive<br/>START_WITH: Starts with, path case sensitive<br/>END_WITH: Ends with, path case sensitive<br/>WILDCARD: Wildcard matches, path case sensitive, * represents zero or more arbitrary characters, ? represents any single character.<br/>NOT_WILDCARD: Wildcard does not match, path case sensitive, * represents zero or more arbitrary characters, ? represents any single character",
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
							Description: "URI match conditions, match type cannot be repeated.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"match_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Match type.<br/>EQUAL: Equals, URI case sensitive<br/>NOT_EQUAL: Does not equal, URI case sensitive<br/>CONTAIN: Contains, URI case insensitive<br/>NOT_CONTAIN: Does not Contains, URI case insensitive<br/>REGEX: Regex match, URI case insensitive<br/>NOT_REGEX: Regular does not match, URI case insensitive<br/>START_WITH: Starts with, URI case insensitive<br/>END_WITH: Ends with, URI case insensitive<br/>WILDCARD: Wildcard matches, URI case insensitive, * represents zero or more arbitrary characters, ? represents any single character<br/>NOT_WILDCARD: Wildcard does not match, URI case insensitive, * represents zero or more arbitrary characters, ? represents any single character",
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
						"ua_conditions": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "User agent match conditions, match type cannot be repeated.",
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
							Description: "Referer match conditions, match type cannot be repeated.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"match_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Match type.<br/>EQUAL: Equals, referer case sensitive<br/>NOT_EQUAL: Does not equal, referer case sensitive<br/>CONTAIN: Contains, referer case insensitive<br/>NOT_CONTAIN: Does not Contains, referer case insensitive<br/>NONE:Empty or non-existent<br/>REGEX: Regex match, referer case insensitive<br/>NOT_REGEX: Regular does not match, referer case insensitive<br/>START_WITH: Starts with, referer case insensitive<br/>END_WITH: Ends with, referer case insensitive<br/>WILDCARD: Wildcard matches, referer case insensitive, * represents zero or more arbitrary characters, ? represents any single characte<br/>NOT_WILDCARD: Wildcard does not match, referer case insensitive, * represents zero or more arbitrary characters, ? represents any single character",
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
							Description: "Request header match conditions.",
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
										Description: "Header name,case insensitive,up to 100 characters.<br/>Example: Accept.",
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
					},
				},
			},
		},
	}
}

func resourceWaapShareWhitelistCreate(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_waap_share_whitelist.create")

	var diags diag.Diagnostics
	request := &waapShareWhitelist.CreateShareWhitelistRuleRequest{}

	if domains, ok := data.GetOk("relation_domain_list"); ok {
		domainsList := domains.([]interface{})
		domainsStr := make([]*string, len(domainsList))
		for i, v := range domainsList {
			str := v.(string)
			domainsStr[i] = &str
		}
		request.RelationDomainList = domainsStr
	}

	if ruleName, ok := data.Get("rule_name").(string); ok && ruleName != "" {
		request.RuleName = &ruleName
	}
	if description, ok := data.Get("description").(string); ok && description != "" {
		request.Description = &description
	}

	conditions := data.Get("conditions").([]interface{})
	conditionsRequest := &waapShareWhitelist.ShareWhitelistRuleConditions{}
	for _, condition := range conditions {
		conditionMap := condition.(map[string]interface{})
		if conditionMap["ip_or_ips_conditions"] != nil {
			ipOrIpsConditions := make([]*waapShareWhitelist.IpOrIpsCondition, 0)
			for _, ipOrIpsCondition := range conditionMap["ip_or_ips_conditions"].([]interface{}) {
				ipOrIpsConditionMap := ipOrIpsCondition.(map[string]interface{})
				matchType := ipOrIpsConditionMap["match_type"].(string)
				ipOrIpsInterface := ipOrIpsConditionMap["ip_or_ips"].([]interface{})
				ipOrIps := make([]*string, len(ipOrIpsInterface))
				for i, v := range ipOrIpsInterface {
					str := v.(string)
					ipOrIps[i] = &str
				}
				ipOrIpsCondition := &waapShareWhitelist.IpOrIpsCondition{
					MatchType: &matchType,
					IpOrIps:   ipOrIps,
				}
				ipOrIpsConditions = append(ipOrIpsConditions, ipOrIpsCondition)
			}
			conditionsRequest.IpOrIpsConditions = ipOrIpsConditions
		}

		if conditionMap["path_conditions"] != nil {
			pathConditions := make([]*waapShareWhitelist.PathCondition, 0)
			for _, pathCondition := range conditionMap["path_conditions"].([]interface{}) {
				pathConditionMap := pathCondition.(map[string]interface{})
				matchType := pathConditionMap["match_type"].(string)
				pathsInterface := pathConditionMap["paths"].([]interface{})
				paths := make([]*string, len(pathsInterface))
				for i, v := range pathsInterface {
					str := v.(string)
					paths[i] = &str
				}
				pathCondition := &waapShareWhitelist.PathCondition{
					MatchType: &matchType,
					Paths:     paths,
				}
				pathConditions = append(pathConditions, pathCondition)
			}
			conditionsRequest.PathConditions = pathConditions
		}

		if conditionMap["uri_conditions"] != nil {
			uriConditions := make([]*waapShareWhitelist.UriCondition, 0)
			for _, uriCondition := range conditionMap["uri_conditions"].([]interface{}) {
				uriConditionMap := uriCondition.(map[string]interface{})
				matchType := uriConditionMap["match_type"].(string)
				uriInterface := uriConditionMap["uri"].([]interface{})
				uri := make([]*string, len(uriInterface))
				for i, v := range uriInterface {
					str := v.(string)
					uri[i] = &str
				}
				uriCondition := &waapShareWhitelist.UriCondition{
					MatchType: &matchType,
					Uri:       uri,
				}
				uriConditions = append(uriConditions, uriCondition)
			}
			conditionsRequest.UriConditions = uriConditions
		}

		if conditionMap["ua_conditions"] != nil {
			uaConditions := make([]*waapShareWhitelist.UaCondition, 0)
			for _, uaCondition := range conditionMap["ua_conditions"].([]interface{}) {
				uaConditionMap := uaCondition.(map[string]interface{})
				matchType := uaConditionMap["match_type"].(string)
				uaInterface := uaConditionMap["ua"].([]interface{})
				ua := make([]*string, len(uaInterface))
				for i, v := range uaInterface {
					str := v.(string)
					ua[i] = &str
				}
				uaCondition := &waapShareWhitelist.UaCondition{
					MatchType: &matchType,
					Ua:        ua,
				}
				uaConditions = append(uaConditions, uaCondition)
			}
			conditionsRequest.UaConditions = uaConditions
		}

		if conditionMap["referer_conditions"] != nil {
			refererConditions := make([]*waapShareWhitelist.RefererCondition, 0)
			for _, refererCondition := range conditionMap["referer_conditions"].([]interface{}) {
				refererConditionMap := refererCondition.(map[string]interface{})
				matchType := refererConditionMap["match_type"].(string)
				refererInterface := refererConditionMap["referer"].([]interface{})
				referer := make([]*string, len(refererInterface))
				for i, v := range refererInterface {
					str := v.(string)
					referer[i] = &str
				}
				refererCondition := &waapShareWhitelist.RefererCondition{
					MatchType: &matchType,
					Referer:   referer,
				}
				refererConditions = append(refererConditions, refererCondition)
			}
			conditionsRequest.RefererConditions = refererConditions
		}

		if conditionMap["header_conditions"] != nil {
			headerConditions := make([]*waapShareWhitelist.HeaderCondition, 0)
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
				headerCondition := &waapShareWhitelist.HeaderCondition{
					MatchType: &matchType,
					Key:       &key,
					ValueList: valueList,
				}
				headerConditions = append(headerConditions, headerCondition)
			}
			conditionsRequest.HeaderConditions = headerConditions
		}
	}
	request.Conditions = conditionsRequest

	var response *waapShareWhitelist.CreateShareWhitelistRuleResponse
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseWaapShareWhitelistClient().AddWaapShareWhitelistRule(request)
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
	return resourceWaapShareWhitelistRead(context, data, meta)
}

func resourceWaapShareWhitelistRead(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_waap_shareWhitelist.read")
	var response *waapShareWhitelist.ListShareWhitelistRulesResponse
	var err error
	var diags diag.Diagnostics
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		request := &waapShareWhitelist.ListShareWhitelistRulesRequest{}
		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseWaapShareWhitelistClient().GetWaapShareWhitelistList(request)
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
			_ = data.Set("id", item.Id)
			_ = data.Set("relation_domain_list", item.RelationDomainList)
			_ = data.Set("rule_name", item.RuleName)
			_ = data.Set("description", item.Description)
			if item.Conditions != nil {
				conditions := make(map[string]interface{})
				if item.Conditions.IpOrIpsConditions != nil {
					ipOrIpsConditions := make([]interface{}, 0)
					for _, condition := range item.Conditions.IpOrIpsConditions {
						ipOrIpsCondition := map[string]interface{}{
							"match_type": condition.MatchType,
							"ip_or_ips":  condition.IpOrIps,
						}
						ipOrIpsConditions = append(ipOrIpsConditions, ipOrIpsCondition)
					}
					conditions["ip_or_ips_conditions"] = ipOrIpsConditions
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
					conditions["path_conditions"] = pathConditions
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
					conditions["uri_conditions"] = uriConditions
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
					conditions["ua_conditions"] = uaConditions
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
					conditions["referer_conditions"] = refererConditions
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
					conditions["header_conditions"] = headerConditions
				}
				_ = data.Set("conditions", conditions)
			}
		}
	}
	return nil
}

func resourceWaapShareWhitelistUpdate(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_waap_shareWhitelist.update")
	var diags diag.Diagnostics
	request := &waapShareWhitelist.UpdateShareWhitelistRuleRequest{}
	if id, ok := data.Get("id").(string); ok && id != "" {
		request.Id = &id
	}

	if domains, ok := data.GetOk("relation_domain_list"); ok {
		domainsList := domains.([]interface{})
		domainsStr := make([]*string, len(domainsList))
		for i, v := range domainsList {
			str := v.(string)
			domainsStr[i] = &str
		}
		request.RelationDomainList = domainsStr
	}

	if ruleName, ok := data.Get("rule_name").(string); ok && ruleName != "" {
		request.RuleName = &ruleName
	}
	if description, ok := data.Get("description").(string); ok && description != "" {
		request.Description = &description
	}

	conditions := data.Get("conditions").([]interface{})
	conditionsRequest := &waapShareWhitelist.ShareWhitelistRuleConditions{}
	for _, condition := range conditions {
		conditionMap := condition.(map[string]interface{})
		if conditionMap["ip_or_ips_conditions"] != nil {
			ipOrIpsConditions := make([]*waapShareWhitelist.IpOrIpsCondition, 0)
			for _, ipOrIpsCondition := range conditionMap["ip_or_ips_conditions"].([]interface{}) {
				ipOrIpsConditionMap := ipOrIpsCondition.(map[string]interface{})
				matchType := ipOrIpsConditionMap["match_type"].(string)
				ipOrIpsInterface := ipOrIpsConditionMap["ip_or_ips"].([]interface{})
				ipOrIps := make([]*string, len(ipOrIpsInterface))
				for i, v := range ipOrIpsInterface {
					str := v.(string)
					ipOrIps[i] = &str
				}
				ipOrIpsCondition := &waapShareWhitelist.IpOrIpsCondition{
					MatchType: &matchType,
					IpOrIps:   ipOrIps,
				}
				ipOrIpsConditions = append(ipOrIpsConditions, ipOrIpsCondition)
			}
			conditionsRequest.IpOrIpsConditions = ipOrIpsConditions
		}

		if conditionMap["path_conditions"] != nil {
			pathConditions := make([]*waapShareWhitelist.PathCondition, 0)
			for _, pathCondition := range conditionMap["path_conditions"].([]interface{}) {
				pathConditionMap := pathCondition.(map[string]interface{})
				matchType := pathConditionMap["match_type"].(string)
				pathsInterface := pathConditionMap["paths"].([]interface{})
				paths := make([]*string, len(pathsInterface))
				for i, v := range pathsInterface {
					str := v.(string)
					paths[i] = &str
				}
				pathCondition := &waapShareWhitelist.PathCondition{
					MatchType: &matchType,
					Paths:     paths,
				}
				pathConditions = append(pathConditions, pathCondition)
			}
			conditionsRequest.PathConditions = pathConditions
		}

		if conditionMap["uri_conditions"] != nil {
			uriConditions := make([]*waapShareWhitelist.UriCondition, 0)
			for _, uriCondition := range conditionMap["uri_conditions"].([]interface{}) {
				uriConditionMap := uriCondition.(map[string]interface{})
				matchType := uriConditionMap["match_type"].(string)
				uriInterface := uriConditionMap["uri"].([]interface{})
				uri := make([]*string, len(uriInterface))
				for i, v := range uriInterface {
					str := v.(string)
					uri[i] = &str
				}
				uriCondition := &waapShareWhitelist.UriCondition{
					MatchType: &matchType,
					Uri:       uri,
				}
				uriConditions = append(uriConditions, uriCondition)
			}
			conditionsRequest.UriConditions = uriConditions
		}

		if conditionMap["ua_conditions"] != nil {
			uaConditions := make([]*waapShareWhitelist.UaCondition, 0)
			for _, uaCondition := range conditionMap["ua_conditions"].([]interface{}) {
				uaConditionMap := uaCondition.(map[string]interface{})
				matchType := uaConditionMap["match_type"].(string)
				uaInterface := uaConditionMap["ua"].([]interface{})
				ua := make([]*string, len(uaInterface))
				for i, v := range uaInterface {
					str := v.(string)
					ua[i] = &str
				}
				uaCondition := &waapShareWhitelist.UaCondition{
					MatchType: &matchType,
					Ua:        ua,
				}
				uaConditions = append(uaConditions, uaCondition)
			}
			conditionsRequest.UaConditions = uaConditions
		}

		if conditionMap["referer_conditions"] != nil {
			refererConditions := make([]*waapShareWhitelist.RefererCondition, 0)
			for _, refererCondition := range conditionMap["referer_conditions"].([]interface{}) {
				refererConditionMap := refererCondition.(map[string]interface{})
				matchType := refererConditionMap["match_type"].(string)
				refererInterface := refererConditionMap["referer"].([]interface{})
				referer := make([]*string, len(refererInterface))
				for i, v := range refererInterface {
					str := v.(string)
					referer[i] = &str
				}
				refererCondition := &waapShareWhitelist.RefererCondition{
					MatchType: &matchType,
					Referer:   referer,
				}
				refererConditions = append(refererConditions, refererCondition)
			}
			conditionsRequest.RefererConditions = refererConditions
		}

		if conditionMap["header_conditions"] != nil {
			headerConditions := make([]*waapShareWhitelist.HeaderCondition, 0)
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
				headerCondition := &waapShareWhitelist.HeaderCondition{
					MatchType: &matchType,
					Key:       &key,
					ValueList: valueList,
				}
				headerConditions = append(headerConditions, headerCondition)
			}
			conditionsRequest.HeaderConditions = headerConditions
		}
	}
	request.Conditions = conditionsRequest

	var response *waapShareWhitelist.UpdateShareWhitelistRuleResponse
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseWaapShareWhitelistClient().UpdateWaapShareWhitelist(request)
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
	log.Printf("resource.wangsu_waap_shareWhitelist.update success")
	return nil
}

func resourceWaapShareWhitelistDelete(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_waap_shareWhitelist.delete")

	var response *waapShareWhitelist.DeleteShareWhitelistRuleResponse
	var err error
	var diags diag.Diagnostics
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		id := data.Id()
		requset := &waapShareWhitelist.DeleteShareWhitelistRuleRequest{
			IdList: []*string{&id},
		}
		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseWaapShareWhitelistClient().DeleteWaapShareWhitelist(requset)
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
