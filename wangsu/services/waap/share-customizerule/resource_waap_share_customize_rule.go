package share_customizerule

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wangsuCommon "github.com/wangsu-api/terraform-provider-wangsu/wangsu/common"
	waapShareCustomizerule "github.com/wangsu-api/wangsu-sdk-go/wangsu/waap/share-customizerule"
	"log"
	"time"
)

func ResourceWaapShareCustomizeRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWaapShareCustomizeRuleCreate,
		ReadContext:   resourceWaapShareCustomizeRuleRead,
		UpdateContext: resourceWaapShareCustomizeRuleUpdate,
		DeleteContext: resourceWaapShareCustomizeRuleDelete,

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
				Description: "Rule Name, maximum 50 characters.<br/>Does not support special characters and spaces.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description, maximum 200 characters.",
			},
			"act": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Action.<br/>NO_USE:Not Used<br/>LOG:Log<br/>DELAY:Delay<br/>BLOCK:Deny<br/>RESET:Reset Connection",
			},
			"condition": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Match Conditions.",
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
							Description: "URI Parameter, match type cannot be repeated.<br/>When the business scenario is API, this matching condition is not supported.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"match_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Match type.<br/>EQUAL:Equals<br/>NOT_EQUAL:Does not equal<br/>CONTAIN:Contains<br/>NOT_CONTAIN:Does not contains<br/>REGEX:Regex match<br/>NONE:Empty or non-existent",
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
						"area_conditions": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Geo, match type cannot be repeated.",
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

func resourceWaapShareCustomizeRuleCreate(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_waap_share_customize_rule.create")

	var diags diag.Diagnostics
	request := &waapShareCustomizerule.CreateSharedCustomRuleRequest{}
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
	if act, ok := data.Get("act").(string); ok && act != "" {
		request.Act = &act
	}
	conditions := data.Get("condition").([]interface{})
	conditionsRequest := &waapShareCustomizerule.ShareCustomizeRuleCondition{}
	for _, condition := range conditions {
		conditionMap := condition.(map[string]interface{})
		// IpOrIps Conditions
		if conditionMap["ip_or_ips_conditions"] != nil {
			ipOrIpsConditions := make([]*waapShareCustomizerule.IpOrIpsCondition, 0)
			for _, ipOrIpsCondition := range conditionMap["ip_or_ips_conditions"].([]interface{}) {
				ipOrIpsConditionMap := ipOrIpsCondition.(map[string]interface{})
				matchType := ipOrIpsConditionMap["match_type"].(string)
				ipOrIpsInterface := ipOrIpsConditionMap["ip_or_ips"].([]interface{})
				ipOrIps := make([]*string, len(ipOrIpsInterface))
				for i, v := range ipOrIpsInterface {
					str := v.(string)
					ipOrIps[i] = &str
				}
				ipOrIpsCondition := &waapShareCustomizerule.IpOrIpsCondition{
					MatchType: &matchType,
					IpOrIps:   ipOrIps,
				}
				ipOrIpsConditions = append(ipOrIpsConditions, ipOrIpsCondition)
			}
			conditionsRequest.IpOrIpsConditions = ipOrIpsConditions
		}

		// Path Conditions
		if conditionMap["path_conditions"] != nil {
			pathConditions := make([]*waapShareCustomizerule.PathCondition, 0)
			for _, pathCondition := range conditionMap["path_conditions"].([]interface{}) {
				pathConditionMap := pathCondition.(map[string]interface{})
				matchType := pathConditionMap["match_type"].(string)
				pathsInterface := pathConditionMap["paths"].([]interface{})
				paths := make([]*string, len(pathsInterface))
				for i, v := range pathsInterface {
					str := v.(string)
					paths[i] = &str
				}
				pathCondition := &waapShareCustomizerule.PathCondition{
					MatchType: &matchType,
					Paths:     paths,
				}
				pathConditions = append(pathConditions, pathCondition)
			}
			conditionsRequest.PathConditions = pathConditions
		}

		// URI Conditions
		if conditionMap["uri_conditions"] != nil {
			uriConditions := make([]*waapShareCustomizerule.UriCondition, 0)
			for _, uriCondition := range conditionMap["uri_conditions"].([]interface{}) {
				uriConditionMap := uriCondition.(map[string]interface{})
				matchType := uriConditionMap["match_type"].(string)
				uriInterface := uriConditionMap["uri"].([]interface{})
				uri := make([]*string, len(uriInterface))
				for i, v := range uriInterface {
					str := v.(string)
					uri[i] = &str
				}
				uriCondition := &waapShareCustomizerule.UriCondition{
					MatchType: &matchType,
					Uri:       uri,
				}
				uriConditions = append(uriConditions, uriCondition)
			}
			conditionsRequest.UriConditions = uriConditions
		}

		// UA Conditions
		if conditionMap["ua_conditions"] != nil {
			uaConditions := make([]*waapShareCustomizerule.UaCondition, 0)
			for _, uaCondition := range conditionMap["ua_conditions"].([]interface{}) {
				uaConditionMap := uaCondition.(map[string]interface{})
				matchType := uaConditionMap["match_type"].(string)
				uaInterface := uaConditionMap["ua"].([]interface{})
				ua := make([]*string, len(uaInterface))
				for i, v := range uaInterface {
					str := v.(string)
					ua[i] = &str
				}
				uaCondition := &waapShareCustomizerule.UaCondition{
					MatchType: &matchType,
					Ua:        ua,
				}
				uaConditions = append(uaConditions, uaCondition)
			}
			conditionsRequest.UaConditions = uaConditions
		}

		// Referer Conditions
		if conditionMap["referer_conditions"] != nil {
			refererConditions := make([]*waapShareCustomizerule.RefererCondition, 0)
			for _, refererCondition := range conditionMap["referer_conditions"].([]interface{}) {
				refererConditionMap := refererCondition.(map[string]interface{})
				matchType := refererConditionMap["match_type"].(string)
				refererInterface := refererConditionMap["referer"].([]interface{})
				referer := make([]*string, len(refererInterface))
				for i, v := range refererInterface {
					str := v.(string)
					referer[i] = &str
				}
				refererCondition := &waapShareCustomizerule.RefererCondition{
					MatchType: &matchType,
					Referer:   referer,
				}
				refererConditions = append(refererConditions, refererCondition)
			}
			conditionsRequest.RefererConditions = refererConditions
		}

		// Header Conditions
		if conditionMap["header_conditions"] != nil {
			headerConditions := make([]*waapShareCustomizerule.HeaderCondition, 0)
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
				headerCondition := &waapShareCustomizerule.HeaderCondition{
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
			areaConditions := make([]*waapShareCustomizerule.AreaCondition, 0)
			for _, areaCondition := range conditionMap["area_conditions"].([]interface{}) {
				areaConditionMap := areaCondition.(map[string]interface{})
				matchType := areaConditionMap["match_type"].(string)
				areasInterface := areaConditionMap["areas"].([]interface{})
				areas := make([]*string, len(areasInterface))
				for i, v := range areasInterface {
					str := v.(string)
					areas[i] = &str
				}
				areaCondition := &waapShareCustomizerule.AreaCondition{
					MatchType: &matchType,
					Areas:     areas,
				}
				areaConditions = append(areaConditions, areaCondition)
			}
			conditionsRequest.AreaConditions = areaConditions
		}

		// Method Conditions
		if conditionMap["method_conditions"] != nil {
			methodConditions := make([]*waapShareCustomizerule.MethodCondition, 0)
			for _, methodCondition := range conditionMap["method_conditions"].([]interface{}) {
				methodConditionMap := methodCondition.(map[string]interface{})
				matchType := methodConditionMap["match_type"].(string)
				requestMethodInterface := methodConditionMap["request_method"].([]interface{})
				requestMethod := make([]*string, len(requestMethodInterface))
				for i, v := range requestMethodInterface {
					str := v.(string)
					requestMethod[i] = &str
				}
				methodCondition := &waapShareCustomizerule.MethodCondition{
					MatchType:     &matchType,
					RequestMethod: requestMethod,
				}
				methodConditions = append(methodConditions, methodCondition)
			}
			conditionsRequest.MethodConditions = methodConditions
		}

		// JA3 Conditions
		if conditionMap["ja3_conditions"] != nil {
			ja3Conditions := make([]*waapShareCustomizerule.Ja3Condition, 0)
			for _, ja3Condition := range conditionMap["ja3_conditions"].([]interface{}) {
				ja3ConditionMap := ja3Condition.(map[string]interface{})
				matchType := ja3ConditionMap["match_type"].(string)
				ja3Interface := ja3ConditionMap["ja3_list"].([]interface{})
				ja3 := make([]*string, len(ja3Interface))
				for i, v := range ja3Interface {
					str := v.(string)
					ja3[i] = &str
				}
				ja3Condition := &waapShareCustomizerule.Ja3Condition{
					MatchType: &matchType,
					Ja3List:   ja3,
				}
				ja3Conditions = append(ja3Conditions, ja3Condition)
			}
			conditionsRequest.Ja3Conditions = ja3Conditions
		}

		// JA4 Conditions
		if conditionMap["ja4_conditions"] != nil {
			ja4Conditions := make([]*waapShareCustomizerule.Ja4Condition, 0)
			for _, ja4Condition := range conditionMap["ja4_conditions"].([]interface{}) {
				ja4ConditionMap := ja4Condition.(map[string]interface{})
				matchType := ja4ConditionMap["match_type"].(string)
				ja4Interface := ja4ConditionMap["ja4_list"].([]interface{})
				ja4 := make([]*string, len(ja4Interface))
				for i, v := range ja4Interface {
					str := v.(string)
					ja4[i] = &str
				}
				ja4Condition := &waapShareCustomizerule.Ja4Condition{
					MatchType: &matchType,
					Ja4List:   ja4,
				}
				ja4Conditions = append(ja4Conditions, ja4Condition)
			}
			conditionsRequest.Ja4Conditions = ja4Conditions
		}
	}
	request.Condition = conditionsRequest

	var response *waapShareCustomizerule.CreateSharedCustomRuleResponse
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseWaapShareCustomizeruleClient().Add(request)
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
	return resourceWaapShareCustomizeRuleRead(context, data, meta)
}

func resourceWaapShareCustomizeRuleRead(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_waap_share_customize_rule.read")
	var response *waapShareCustomizerule.ListSharedCustomRulesResponse
	var err error
	var diags diag.Diagnostics
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		request := &waapShareCustomizerule.ListSharedCustomRulesRequest{}
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
		for _, item := range response.Data {
			// 只要对应id的数据
			if *item.Id != data.Id() {
				continue
			}
			_ = data.Set("id", item.Id)
			_ = data.Set("relation_domain_list", item.RelationDomainList)
			_ = data.Set("rule_name", item.RuleName)
			_ = data.Set("description", item.Description)
			_ = data.Set("act", item.Act)
			condition := make(map[string]interface{})
			if item.Condition != nil {
				if item.Condition.IpOrIpsConditions != nil {
					ipOrIpsConditions := make([]interface{}, 0)
					for _, condition := range item.Condition.IpOrIpsConditions {
						ipOrIpsCondition := map[string]interface{}{
							"match_type": condition.MatchType,
							"ip_or_ips":  condition.IpOrIps,
						}
						ipOrIpsConditions = append(ipOrIpsConditions, ipOrIpsCondition)
					}
					condition["ip_or_ips_conditions"] = ipOrIpsConditions
				}
				if item.Condition.PathConditions != nil {
					pathConditions := make([]interface{}, 0)
					for _, condition := range item.Condition.PathConditions {
						pathCondition := map[string]interface{}{
							"match_type": condition.MatchType,
							"paths":      condition.Paths,
						}
						pathConditions = append(pathConditions, pathCondition)
					}
					condition["path_conditions"] = pathConditions
				}
				if item.Condition.UriConditions != nil {
					uriConditions := make([]interface{}, 0)
					for _, condition := range item.Condition.UriConditions {
						uriCondition := map[string]interface{}{
							"match_type": condition.MatchType,
							"uri":        condition.Uri,
						}
						uriConditions = append(uriConditions, uriCondition)
					}
					condition["uri_conditions"] = uriConditions
				}
				if item.Condition.UaConditions != nil {
					uaConditions := make([]interface{}, 0)
					for _, condition := range item.Condition.UaConditions {
						uaCondition := map[string]interface{}{
							"match_type": condition.MatchType,
							"ua":         condition.Ua,
						}
						uaConditions = append(uaConditions, uaCondition)
					}
					condition["ua_conditions"] = uaConditions
				}
				if item.Condition.RefererConditions != nil {
					refererConditions := make([]interface{}, 0)
					for _, condition := range item.Condition.RefererConditions {
						refererCondition := map[string]interface{}{
							"match_type": condition.MatchType,
							"referer":    condition.Referer,
						}
						refererConditions = append(refererConditions, refererCondition)
					}
					condition["referer_conditions"] = refererConditions
				}
				if item.Condition.HeaderConditions != nil {
					headerConditions := make([]interface{}, 0)
					for _, condition := range item.Condition.HeaderConditions {
						headerCondition := map[string]interface{}{
							"match_type": condition.MatchType,
							"key":        condition.Key,
							"value_list": condition.ValueList,
						}
						headerConditions = append(headerConditions, headerCondition)
					}
					condition["header_conditions"] = headerConditions
				}
				if item.Condition.AreaConditions != nil {
					areaConditions := make([]interface{}, 0)
					for _, condition := range item.Condition.AreaConditions {
						areaCondition := map[string]interface{}{
							"match_type": condition.MatchType,
							"areas":      condition.Areas,
						}
						areaConditions = append(areaConditions, areaCondition)
					}
					condition["area_conditions"] = areaConditions
				}
				if item.Condition.MethodConditions != nil {
					methodConditions := make([]interface{}, 0)
					for _, condition := range item.Condition.MethodConditions {
						methodCondition := map[string]interface{}{
							"match_type":     condition.MatchType,
							"request_method": condition.RequestMethod,
						}
						methodConditions = append(methodConditions, methodCondition)
					}
					condition["method_conditions"] = methodConditions
				}
				if item.Condition.Ja3Conditions != nil {
					ja3Conditions := make([]interface{}, 0)
					for _, condition := range item.Condition.Ja3Conditions {
						ja3Condition := map[string]interface{}{
							"match_type": condition.MatchType,
							"ja3_list":   condition.Ja3List,
						}
						ja3Conditions = append(ja3Conditions, ja3Condition)
					}
					condition["ja3_conditions"] = ja3Conditions
				}
				if item.Condition.Ja4Conditions != nil {
					ja4Conditions := make([]interface{}, 0)
					for _, condition := range item.Condition.Ja4Conditions {
						ja4Condition := map[string]interface{}{
							"match_type": condition.MatchType,
							"ja4_list":   condition.Ja4List,
						}
						ja4Conditions = append(ja4Conditions, ja4Condition)
					}
					condition["ja4_conditions"] = ja4Conditions
				}
			}
			_ = data.Set("condition", condition)
		}
	}
	return nil
}

func resourceWaapShareCustomizeRuleUpdate(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_waap_share_customize_rule.update")
	var diags diag.Diagnostics
	request := &waapShareCustomizerule.UpdateSharedCustomRulesRequest{}
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
	if act, ok := data.Get("act").(string); ok && act != "" {
		request.Act = &act
	}
	conditions := data.Get("condition").([]interface{})
	conditionsRequest := &waapShareCustomizerule.ShareCustomizeRuleCondition{}
	for _, condition := range conditions {
		conditionMap := condition.(map[string]interface{})
		// IpOrIps Conditions
		if conditionMap["ip_or_ips_conditions"] != nil {
			ipOrIpsConditions := make([]*waapShareCustomizerule.IpOrIpsCondition, 0)
			for _, ipOrIpsCondition := range conditionMap["ip_or_ips_conditions"].([]interface{}) {
				ipOrIpsConditionMap := ipOrIpsCondition.(map[string]interface{})
				matchType := ipOrIpsConditionMap["match_type"].(string)
				ipOrIpsInterface := ipOrIpsConditionMap["ip_or_ips"].([]interface{})
				ipOrIps := make([]*string, len(ipOrIpsInterface))
				for i, v := range ipOrIpsInterface {
					str := v.(string)
					ipOrIps[i] = &str
				}
				ipOrIpsCondition := &waapShareCustomizerule.IpOrIpsCondition{
					MatchType: &matchType,
					IpOrIps:   ipOrIps,
				}
				ipOrIpsConditions = append(ipOrIpsConditions, ipOrIpsCondition)
			}
			conditionsRequest.IpOrIpsConditions = ipOrIpsConditions
		}

		// Path Conditions
		if conditionMap["path_conditions"] != nil {
			pathConditions := make([]*waapShareCustomizerule.PathCondition, 0)
			for _, pathCondition := range conditionMap["path_conditions"].([]interface{}) {
				pathConditionMap := pathCondition.(map[string]interface{})
				matchType := pathConditionMap["match_type"].(string)
				pathsInterface := pathConditionMap["paths"].([]interface{})
				paths := make([]*string, len(pathsInterface))
				for i, v := range pathsInterface {
					str := v.(string)
					paths[i] = &str
				}
				pathCondition := &waapShareCustomizerule.PathCondition{
					MatchType: &matchType,
					Paths:     paths,
				}
				pathConditions = append(pathConditions, pathCondition)
			}
			conditionsRequest.PathConditions = pathConditions
		}

		// URI Conditions
		if conditionMap["uri_conditions"] != nil {
			uriConditions := make([]*waapShareCustomizerule.UriCondition, 0)
			for _, uriCondition := range conditionMap["uri_conditions"].([]interface{}) {
				uriConditionMap := uriCondition.(map[string]interface{})
				matchType := uriConditionMap["match_type"].(string)
				uriInterface := uriConditionMap["uri"].([]interface{})
				uri := make([]*string, len(uriInterface))
				for i, v := range uriInterface {
					str := v.(string)
					uri[i] = &str
				}
				uriCondition := &waapShareCustomizerule.UriCondition{
					MatchType: &matchType,
					Uri:       uri,
				}
				uriConditions = append(uriConditions, uriCondition)
			}
			conditionsRequest.UriConditions = uriConditions
		}

		// UA Conditions
		if conditionMap["ua_conditions"] != nil {
			uaConditions := make([]*waapShareCustomizerule.UaCondition, 0)
			for _, uaCondition := range conditionMap["ua_conditions"].([]interface{}) {
				uaConditionMap := uaCondition.(map[string]interface{})
				matchType := uaConditionMap["match_type"].(string)
				uaInterface := uaConditionMap["ua"].([]interface{})
				ua := make([]*string, len(uaInterface))
				for i, v := range uaInterface {
					str := v.(string)
					ua[i] = &str
				}
				uaCondition := &waapShareCustomizerule.UaCondition{
					MatchType: &matchType,
					Ua:        ua,
				}
				uaConditions = append(uaConditions, uaCondition)
			}
			conditionsRequest.UaConditions = uaConditions
		}

		// Referer Conditions
		if conditionMap["referer_conditions"] != nil {
			refererConditions := make([]*waapShareCustomizerule.RefererCondition, 0)
			for _, refererCondition := range conditionMap["referer_conditions"].([]interface{}) {
				refererConditionMap := refererCondition.(map[string]interface{})
				matchType := refererConditionMap["match_type"].(string)
				refererInterface := refererConditionMap["referer"].([]interface{})
				referer := make([]*string, len(refererInterface))
				for i, v := range refererInterface {
					str := v.(string)
					referer[i] = &str
				}
				refererCondition := &waapShareCustomizerule.RefererCondition{
					MatchType: &matchType,
					Referer:   referer,
				}
				refererConditions = append(refererConditions, refererCondition)
			}
			conditionsRequest.RefererConditions = refererConditions
		}

		// Header Conditions
		if conditionMap["header_conditions"] != nil {
			headerConditions := make([]*waapShareCustomizerule.HeaderCondition, 0)
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
				headerCondition := &waapShareCustomizerule.HeaderCondition{
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
			areaConditions := make([]*waapShareCustomizerule.AreaCondition, 0)
			for _, areaCondition := range conditionMap["area_conditions"].([]interface{}) {
				areaConditionMap := areaCondition.(map[string]interface{})
				matchType := areaConditionMap["match_type"].(string)
				areasInterface := areaConditionMap["areas"].([]interface{})
				areas := make([]*string, len(areasInterface))
				for i, v := range areasInterface {
					str := v.(string)
					areas[i] = &str
				}
				areaCondition := &waapShareCustomizerule.AreaCondition{
					MatchType: &matchType,
					Areas:     areas,
				}
				areaConditions = append(areaConditions, areaCondition)
			}
			conditionsRequest.AreaConditions = areaConditions
		}

		// Method Conditions
		if conditionMap["method_conditions"] != nil {
			methodConditions := make([]*waapShareCustomizerule.MethodCondition, 0)
			for _, methodCondition := range conditionMap["method_conditions"].([]interface{}) {
				methodConditionMap := methodCondition.(map[string]interface{})
				matchType := methodConditionMap["match_type"].(string)
				requestMethodInterface := methodConditionMap["request_method"].([]interface{})
				requestMethod := make([]*string, len(requestMethodInterface))
				for i, v := range requestMethodInterface {
					str := v.(string)
					requestMethod[i] = &str
				}
				methodCondition := &waapShareCustomizerule.MethodCondition{
					MatchType:     &matchType,
					RequestMethod: requestMethod,
				}
				methodConditions = append(methodConditions, methodCondition)
			}
			conditionsRequest.MethodConditions = methodConditions
		}

		// JA3 Conditions
		if conditionMap["ja3_conditions"] != nil {
			ja3Conditions := make([]*waapShareCustomizerule.Ja3Condition, 0)
			for _, ja3Condition := range conditionMap["ja3_conditions"].([]interface{}) {
				ja3ConditionMap := ja3Condition.(map[string]interface{})
				matchType := ja3ConditionMap["match_type"].(string)
				ja3Interface := ja3ConditionMap["ja3_list"].([]interface{})
				ja3 := make([]*string, len(ja3Interface))
				for i, v := range ja3Interface {
					str := v.(string)
					ja3[i] = &str
				}
				ja3Condition := &waapShareCustomizerule.Ja3Condition{
					MatchType: &matchType,
					Ja3List:   ja3,
				}
				ja3Conditions = append(ja3Conditions, ja3Condition)
			}
			conditionsRequest.Ja3Conditions = ja3Conditions
		}

		// JA4 Conditions
		if conditionMap["ja4_conditions"] != nil {
			ja4Conditions := make([]*waapShareCustomizerule.Ja4Condition, 0)
			for _, ja4Condition := range conditionMap["ja4_conditions"].([]interface{}) {
				ja4ConditionMap := ja4Condition.(map[string]interface{})
				matchType := ja4ConditionMap["match_type"].(string)
				ja4Interface := ja4ConditionMap["ja4_list"].([]interface{})
				ja4 := make([]*string, len(ja4Interface))
				for i, v := range ja4Interface {
					str := v.(string)
					ja4[i] = &str
				}
				ja4Condition := &waapShareCustomizerule.Ja4Condition{
					MatchType: &matchType,
					Ja4List:   ja4,
				}
				ja4Conditions = append(ja4Conditions, ja4Condition)
			}
			conditionsRequest.Ja4Conditions = ja4Conditions
		}
	}
	request.Condition = conditionsRequest

	var response *waapShareCustomizerule.UpdateSharedCustomRulesResponse
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseWaapShareCustomizeruleClient().Update(request)
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
	log.Printf("resource.wangsu_waap_share_customize_rule.update success")
	return nil
}

func resourceWaapShareCustomizeRuleDelete(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_waap_share_customize_rule.delete")

	var response *waapShareCustomizerule.DeleteSharedCustomRulesResponse
	var err error
	var diags diag.Diagnostics
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		id := data.Id()
		request := &waapShareCustomizerule.DeleteSharedCustomRulesRequest{
			IdList: []*string{&id},
		}
		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseWaapShareCustomizeruleClient().Delete(request)
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
