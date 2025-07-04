package customizerule

import (
	"context"
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wangsuCommon "github.com/wangsu-api/terraform-provider-wangsu/wangsu/common"
	waapCustomizerule "github.com/wangsu-api/wangsu-sdk-go/wangsu/waap/customizerule"
	"log"
	"time"
)

func ResourceWaapCustomizeRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWaapCustomizeRuleCreate,
		ReadContext:   resourceWaapCustomizeRuleRead,
		UpdateContext: resourceWaapCustomizeRuleUpdate,
		DeleteContext: resourceWaapCustomizeRuleDelete,

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
				Description: "Rule Name, maximum 50 characters.<br/>Does not support special characters and spaces.",
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
			"api_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "API ID under API business, multiple separated by ; sign.<br/>When the protected target is APIThis field is required.",
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

func resourceWaapCustomizeRuleCreate(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_waap_customize_rule.create")

	var diags diag.Diagnostics
	request := &waapCustomizerule.AddCustomizeRuleRequest{}
	if domain, ok := data.Get("domain").(string); ok && domain != "" {
		request.Domain = &domain
	}
	if ruleName, ok := data.Get("rule_name").(string); ok && ruleName != "" {
		request.RuleName = &ruleName
	}
	if description, ok := data.Get("description").(string); ok && description != "" {
		request.Description = &description
	}
	if scene, ok := data.Get("scene").(string); ok && scene != "" {
		request.Scene = &scene
	}
	if apiId, ok := data.Get("api_id").(string); ok && apiId != "" {
		request.ApiId = &apiId
	}
	if act, ok := data.Get("act").(string); ok && act != "" {
		request.Act = &act
	}
	conditions := data.Get("condition").([]interface{})
	conditionsRequest := &waapCustomizerule.CommonCustomizeRuleConditionDTO{}
	for _, condition := range conditions {
		conditionMap := condition.(map[string]interface{})
		// IpOrIps Conditions
		if conditionMap["ip_or_ips_conditions"] != nil {
			ipOrIpsConditions := make([]*waapCustomizerule.IpOrIpsCondition, 0)
			for _, ipOrIpsCondition := range conditionMap["ip_or_ips_conditions"].([]interface{}) {
				ipOrIpsConditionMap := ipOrIpsCondition.(map[string]interface{})
				matchType := ipOrIpsConditionMap["match_type"].(string)
				ipOrIpsInterface := ipOrIpsConditionMap["ip_or_ips"].([]interface{})
				ipOrIps := make([]*string, len(ipOrIpsInterface))
				for i, v := range ipOrIpsInterface {
					str := v.(string)
					ipOrIps[i] = &str
				}
				ipOrIpsCondition := &waapCustomizerule.IpOrIpsCondition{
					MatchType: &matchType,
					IpOrIps:   ipOrIps,
				}
				ipOrIpsConditions = append(ipOrIpsConditions, ipOrIpsCondition)
			}
			conditionsRequest.IpOrIpsConditions = ipOrIpsConditions
		}

		// Path Conditions
		if conditionMap["path_conditions"] != nil {
			pathConditions := make([]*waapCustomizerule.PathCondition, 0)
			for _, pathCondition := range conditionMap["path_conditions"].([]interface{}) {
				pathConditionMap := pathCondition.(map[string]interface{})
				matchType := pathConditionMap["match_type"].(string)
				pathsInterface := pathConditionMap["paths"].([]interface{})
				paths := make([]*string, len(pathsInterface))
				for i, v := range pathsInterface {
					str := v.(string)
					paths[i] = &str
				}
				pathCondition := &waapCustomizerule.PathCondition{
					MatchType: &matchType,
					Paths:     paths,
				}
				pathConditions = append(pathConditions, pathCondition)
			}
			conditionsRequest.PathConditions = pathConditions
		}

		// URI Conditions
		if conditionMap["uri_conditions"] != nil {
			uriConditions := make([]*waapCustomizerule.UriCondition, 0)
			for _, uriCondition := range conditionMap["uri_conditions"].([]interface{}) {
				uriConditionMap := uriCondition.(map[string]interface{})
				matchType := uriConditionMap["match_type"].(string)
				uriInterface := uriConditionMap["uri"].([]interface{})
				uri := make([]*string, len(uriInterface))
				for i, v := range uriInterface {
					str := v.(string)
					uri[i] = &str
				}
				uriCondition := &waapCustomizerule.UriCondition{
					MatchType: &matchType,
					Uri:       uri,
				}
				uriConditions = append(uriConditions, uriCondition)
			}
			conditionsRequest.UriConditions = uriConditions
		}

		// URI Param Conditions
		if conditionMap["uri_param_conditions"] != nil {
			uriParamConditions := make([]*waapCustomizerule.UriParamCondition, 0)
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
				uriParamCondition := &waapCustomizerule.UriParamCondition{
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
			uaConditions := make([]*waapCustomizerule.UaCondition, 0)
			for _, uaCondition := range conditionMap["ua_conditions"].([]interface{}) {
				uaConditionMap := uaCondition.(map[string]interface{})
				matchType := uaConditionMap["match_type"].(string)
				uaInterface := uaConditionMap["ua"].([]interface{})
				ua := make([]*string, len(uaInterface))
				for i, v := range uaInterface {
					str := v.(string)
					ua[i] = &str
				}
				uaCondition := &waapCustomizerule.UaCondition{
					MatchType: &matchType,
					Ua:        ua,
				}
				uaConditions = append(uaConditions, uaCondition)
			}
			conditionsRequest.UaConditions = uaConditions
		}

		// Referer Conditions
		if conditionMap["referer_conditions"] != nil {
			refererConditions := make([]*waapCustomizerule.RefererCondition, 0)
			for _, refererCondition := range conditionMap["referer_conditions"].([]interface{}) {
				refererConditionMap := refererCondition.(map[string]interface{})
				matchType := refererConditionMap["match_type"].(string)
				refererInterface := refererConditionMap["referer"].([]interface{})
				referer := make([]*string, len(refererInterface))
				for i, v := range refererInterface {
					str := v.(string)
					referer[i] = &str
				}
				refererCondition := &waapCustomizerule.RefererCondition{
					MatchType: &matchType,
					Referer:   referer,
				}
				refererConditions = append(refererConditions, refererCondition)
			}
			conditionsRequest.RefererConditions = refererConditions
		}

		// Header Conditions
		if conditionMap["header_conditions"] != nil {
			headerConditions := make([]*waapCustomizerule.HeaderCondition, 0)
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
				headerCondition := &waapCustomizerule.HeaderCondition{
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
			areaConditions := make([]*waapCustomizerule.AreaCondition, 0)
			for _, areaCondition := range conditionMap["area_conditions"].([]interface{}) {
				areaConditionMap := areaCondition.(map[string]interface{})
				matchType := areaConditionMap["match_type"].(string)
				areasInterface := areaConditionMap["areas"].([]interface{})
				areas := make([]*string, len(areasInterface))
				for i, v := range areasInterface {
					str := v.(string)
					areas[i] = &str
				}
				areaCondition := &waapCustomizerule.AreaCondition{
					MatchType: &matchType,
					Areas:     areas,
				}
				areaConditions = append(areaConditions, areaCondition)
			}
			conditionsRequest.AreaConditions = areaConditions
		}

		// Method Conditions
		if conditionMap["method_conditions"] != nil {
			methodConditions := make([]*waapCustomizerule.RequestMethodCondition, 0)
			for _, methodCondition := range conditionMap["method_conditions"].([]interface{}) {
				methodConditionMap := methodCondition.(map[string]interface{})
				matchType := methodConditionMap["match_type"].(string)
				requestMethodInterface := methodConditionMap["request_method"].([]interface{})
				requestMethod := make([]*string, len(requestMethodInterface))
				for i, v := range requestMethodInterface {
					str := v.(string)
					requestMethod[i] = &str
				}
				methodCondition := &waapCustomizerule.RequestMethodCondition{
					MatchType:     &matchType,
					RequestMethod: requestMethod,
				}
				methodConditions = append(methodConditions, methodCondition)
			}
			conditionsRequest.MethodConditions = methodConditions
		}

		// JA3 Conditions
		if conditionMap["ja3_conditions"] != nil {
			ja3Conditions := make([]*waapCustomizerule.Ja3Condition, 0)
			for _, ja3Condition := range conditionMap["ja3_conditions"].([]interface{}) {
				ja3ConditionMap := ja3Condition.(map[string]interface{})
				matchType := ja3ConditionMap["match_type"].(string)
				ja3Interface := ja3ConditionMap["ja3_list"].([]interface{})
				ja3 := make([]*string, len(ja3Interface))
				for i, v := range ja3Interface {
					str := v.(string)
					ja3[i] = &str
				}
				ja3Condition := &waapCustomizerule.Ja3Condition{
					MatchType: &matchType,
					Ja3List:   ja3,
				}
				ja3Conditions = append(ja3Conditions, ja3Condition)
			}
			conditionsRequest.Ja3Conditions = ja3Conditions
		}

		// JA4 Conditions
		if conditionMap["ja4_conditions"] != nil {
			ja4Conditions := make([]*waapCustomizerule.Ja4Condition, 0)
			for _, ja4Condition := range conditionMap["ja4_conditions"].([]interface{}) {
				ja4ConditionMap := ja4Condition.(map[string]interface{})
				matchType := ja4ConditionMap["match_type"].(string)
				ja4Interface := ja4ConditionMap["ja4_list"].([]interface{})
				ja4 := make([]*string, len(ja4Interface))
				for i, v := range ja4Interface {
					str := v.(string)
					ja4[i] = &str
				}
				ja4Condition := &waapCustomizerule.Ja4Condition{
					MatchType: &matchType,
					Ja4List:   ja4,
				}
				ja4Conditions = append(ja4Conditions, ja4Condition)
			}
			conditionsRequest.Ja4Conditions = ja4Conditions
		}
	}
	request.Condition = conditionsRequest

	var response *waapCustomizerule.AddCustomizeRuleResponse
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseWaapCustomizeruleClient().AddCustomRule(request)
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
	return resourceWaapCustomizeRuleRead(context, data, meta)
}

func resourceWaapCustomizeRuleRead(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_waap_customize_rule.read")
	var response *waapCustomizerule.ListCustomRulesResponse
	var err error
	var diags diag.Diagnostics
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		domain := data.Get("domain").(string)
		// 规则名称会变，不当成查询条件
		//ruleName := data.Get("rule_name").(string)

		request := &waapCustomizerule.ListCustomRulesRequest{
			DomainList: []*string{&domain},
			//RuleName:   &ruleName,
		}
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
		for _, item := range response.Data {
			// 只要对应id的数据
			if *item.Id != data.Id() {
				continue
			}
			_ = data.Set("domain", item.Domain)
			_ = data.Set("rule_name", item.RuleName)
			_ = data.Set("description", item.Description)
			_ = data.Set("scene", item.Scene)
			_ = data.Set("api_id", item.ApiId)
			_ = data.Set("act", item.Act)
			condition := make(map[string]interface{})
			if item.ConditionList != nil {
				if item.ConditionList.IpOrIpsConditions != nil {
					ipOrIpsConditions := make([]interface{}, 0)
					for _, condition := range item.ConditionList.IpOrIpsConditions {
						ipOrIpsCondition := map[string]interface{}{
							"match_type": condition.MatchType,
							"ip_or_ips":  condition.IpOrIps,
						}
						ipOrIpsConditions = append(ipOrIpsConditions, ipOrIpsCondition)
					}
					condition["ip_or_ips_conditions"] = ipOrIpsConditions
				}
				if item.ConditionList.PathConditions != nil {
					pathConditions := make([]interface{}, 0)
					for _, condition := range item.ConditionList.PathConditions {
						pathCondition := map[string]interface{}{
							"match_type": condition.MatchType,
							"paths":      condition.Paths,
						}
						pathConditions = append(pathConditions, pathCondition)
					}
					condition["path_conditions"] = pathConditions
				}
				if item.ConditionList.UriConditions != nil {
					uriConditions := make([]interface{}, 0)
					for _, condition := range item.ConditionList.UriConditions {
						uriCondition := map[string]interface{}{
							"match_type": condition.MatchType,
							"uri":        condition.Uri,
						}
						uriConditions = append(uriConditions, uriCondition)
					}
					condition["uri_conditions"] = uriConditions
				}
				if item.ConditionList.UriParamConditions != nil {
					uriParamConditions := make([]interface{}, 0)
					for _, condition := range item.ConditionList.UriParamConditions {
						uriParamCondition := map[string]interface{}{
							"match_type":  condition.MatchType,
							"param_name":  condition.ParamName,
							"param_value": condition.ParamValue,
						}
						uriParamConditions = append(uriParamConditions, uriParamCondition)
					}
					condition["uri_param_conditions"] = uriParamConditions
				}
				if item.ConditionList.UaConditions != nil {
					uaConditions := make([]interface{}, 0)
					for _, condition := range item.ConditionList.UaConditions {
						uaCondition := map[string]interface{}{
							"match_type": condition.MatchType,
							"ua":         condition.Ua,
						}
						uaConditions = append(uaConditions, uaCondition)
					}
					condition["ua_conditions"] = uaConditions
				}
				if item.ConditionList.RefererConditions != nil {
					refererConditions := make([]interface{}, 0)
					for _, condition := range item.ConditionList.RefererConditions {
						refererCondition := map[string]interface{}{
							"match_type": condition.MatchType,
							"referer":    condition.Referer,
						}
						refererConditions = append(refererConditions, refererCondition)
					}
					condition["referer_conditions"] = refererConditions
				}
				if item.ConditionList.HeaderConditions != nil {
					headerConditions := make([]interface{}, 0)
					for _, condition := range item.ConditionList.HeaderConditions {
						headerCondition := map[string]interface{}{
							"match_type": condition.MatchType,
							"key":        condition.Key,
							"value_list": condition.ValueList,
						}
						headerConditions = append(headerConditions, headerCondition)
					}
					condition["header_conditions"] = headerConditions
				}
				if item.ConditionList.AreaConditions != nil {
					areaConditions := make([]interface{}, 0)
					for _, condition := range item.ConditionList.AreaConditions {
						areaCondition := map[string]interface{}{
							"match_type": condition.MatchType,
							"areas":      condition.Areas,
						}
						areaConditions = append(areaConditions, areaCondition)
					}
					condition["area_conditions"] = areaConditions
				}
				if item.ConditionList.MethodConditions != nil {
					methodConditions := make([]interface{}, 0)
					for _, condition := range item.ConditionList.MethodConditions {
						methodCondition := map[string]interface{}{
							"match_type":     condition.MatchType,
							"request_method": condition.RequestMethod,
						}
						methodConditions = append(methodConditions, methodCondition)
					}
					condition["method_conditions"] = methodConditions
				}
				if item.ConditionList.Ja3Conditions != nil {
					ja3Conditions := make([]interface{}, 0)
					for _, condition := range item.ConditionList.Ja3Conditions {
						ja3Condition := map[string]interface{}{
							"match_type": condition.MatchType,
							"ja3_list":   condition.Ja3List,
						}
						ja3Conditions = append(ja3Conditions, ja3Condition)
					}
					condition["ja3_conditions"] = ja3Conditions
				}
				if item.ConditionList.Ja4Conditions != nil {
					ja4Conditions := make([]interface{}, 0)
					for _, condition := range item.ConditionList.Ja4Conditions {
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

func resourceWaapCustomizeRuleUpdate(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_waap_customize_rule.update")
	var diags diag.Diagnostics
	if data.HasChange("domain") {
		// 把domain强制刷回旧值，否则会有权限问题
		oldDomain, _ := data.GetChange("domain")
		_ = data.Set("domain", oldDomain)
		err := errors.New("Hostname cannot be changed.")
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}
	request := &waapCustomizerule.UpdateCustomRuleRequest{}
	if id, ok := data.Get("id").(string); ok && id != "" {
		request.Id = &id
	}

	if ruleName, ok := data.Get("rule_name").(string); ok && ruleName != "" {
		request.RuleName = &ruleName
	}
	if description, ok := data.Get("description").(string); ok && description != "" {
		request.Description = &description
	}
	if scene, ok := data.Get("scene").(string); ok && scene != "" {
		request.Scene = &scene
	}
	if apiId, ok := data.Get("api_id").(string); ok && apiId != "" {
		request.ApiId = &apiId
	}
	if act, ok := data.Get("act").(string); ok && act != "" {
		request.Act = &act
	}
	conditions := data.Get("condition").([]interface{})
	conditionsRequest := &waapCustomizerule.CommonCustomizeRuleConditionDTO{}
	for _, condition := range conditions {
		conditionMap := condition.(map[string]interface{})
		// IpOrIps Conditions
		if conditionMap["ip_or_ips_conditions"] != nil {
			ipOrIpsConditions := make([]*waapCustomizerule.IpOrIpsCondition, 0)
			for _, ipOrIpsCondition := range conditionMap["ip_or_ips_conditions"].([]interface{}) {
				ipOrIpsConditionMap := ipOrIpsCondition.(map[string]interface{})
				matchType := ipOrIpsConditionMap["match_type"].(string)
				ipOrIpsInterface := ipOrIpsConditionMap["ip_or_ips"].([]interface{})
				ipOrIps := make([]*string, len(ipOrIpsInterface))
				for i, v := range ipOrIpsInterface {
					str := v.(string)
					ipOrIps[i] = &str
				}
				ipOrIpsCondition := &waapCustomizerule.IpOrIpsCondition{
					MatchType: &matchType,
					IpOrIps:   ipOrIps,
				}
				ipOrIpsConditions = append(ipOrIpsConditions, ipOrIpsCondition)
			}
			conditionsRequest.IpOrIpsConditions = ipOrIpsConditions
		}

		// Path Conditions
		if conditionMap["path_conditions"] != nil {
			pathConditions := make([]*waapCustomizerule.PathCondition, 0)
			for _, pathCondition := range conditionMap["path_conditions"].([]interface{}) {
				pathConditionMap := pathCondition.(map[string]interface{})
				matchType := pathConditionMap["match_type"].(string)
				pathsInterface := pathConditionMap["paths"].([]interface{})
				paths := make([]*string, len(pathsInterface))
				for i, v := range pathsInterface {
					str := v.(string)
					paths[i] = &str
				}
				pathCondition := &waapCustomizerule.PathCondition{
					MatchType: &matchType,
					Paths:     paths,
				}
				pathConditions = append(pathConditions, pathCondition)
			}
			conditionsRequest.PathConditions = pathConditions
		}

		// URI Conditions
		if conditionMap["uri_conditions"] != nil {
			uriConditions := make([]*waapCustomizerule.UriCondition, 0)
			for _, uriCondition := range conditionMap["uri_conditions"].([]interface{}) {
				uriConditionMap := uriCondition.(map[string]interface{})
				matchType := uriConditionMap["match_type"].(string)
				uriInterface := uriConditionMap["uri"].([]interface{})
				uri := make([]*string, len(uriInterface))
				for i, v := range uriInterface {
					str := v.(string)
					uri[i] = &str
				}
				uriCondition := &waapCustomizerule.UriCondition{
					MatchType: &matchType,
					Uri:       uri,
				}
				uriConditions = append(uriConditions, uriCondition)
			}
			conditionsRequest.UriConditions = uriConditions
		}

		// URI Param Conditions
		if conditionMap["uri_param_conditions"] != nil {
			uriParamConditions := make([]*waapCustomizerule.UriParamCondition, 0)
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
				uriParamCondition := &waapCustomizerule.UriParamCondition{
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
			uaConditions := make([]*waapCustomizerule.UaCondition, 0)
			for _, uaCondition := range conditionMap["ua_conditions"].([]interface{}) {
				uaConditionMap := uaCondition.(map[string]interface{})
				matchType := uaConditionMap["match_type"].(string)
				uaInterface := uaConditionMap["ua"].([]interface{})
				ua := make([]*string, len(uaInterface))
				for i, v := range uaInterface {
					str := v.(string)
					ua[i] = &str
				}
				uaCondition := &waapCustomizerule.UaCondition{
					MatchType: &matchType,
					Ua:        ua,
				}
				uaConditions = append(uaConditions, uaCondition)
			}
			conditionsRequest.UaConditions = uaConditions
		}

		// Referer Conditions
		if conditionMap["referer_conditions"] != nil {
			refererConditions := make([]*waapCustomizerule.RefererCondition, 0)
			for _, refererCondition := range conditionMap["referer_conditions"].([]interface{}) {
				refererConditionMap := refererCondition.(map[string]interface{})
				matchType := refererConditionMap["match_type"].(string)
				refererInterface := refererConditionMap["referer"].([]interface{})
				referer := make([]*string, len(refererInterface))
				for i, v := range refererInterface {
					str := v.(string)
					referer[i] = &str
				}
				refererCondition := &waapCustomizerule.RefererCondition{
					MatchType: &matchType,
					Referer:   referer,
				}
				refererConditions = append(refererConditions, refererCondition)
			}
			conditionsRequest.RefererConditions = refererConditions
		}

		// Header Conditions
		if conditionMap["header_conditions"] != nil {
			headerConditions := make([]*waapCustomizerule.HeaderCondition, 0)
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
				headerCondition := &waapCustomizerule.HeaderCondition{
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
			areaConditions := make([]*waapCustomizerule.AreaCondition, 0)
			for _, areaCondition := range conditionMap["area_conditions"].([]interface{}) {
				areaConditionMap := areaCondition.(map[string]interface{})
				matchType := areaConditionMap["match_type"].(string)
				areasInterface := areaConditionMap["areas"].([]interface{})
				areas := make([]*string, len(areasInterface))
				for i, v := range areasInterface {
					str := v.(string)
					areas[i] = &str
				}
				areaCondition := &waapCustomizerule.AreaCondition{
					MatchType: &matchType,
					Areas:     areas,
				}
				areaConditions = append(areaConditions, areaCondition)
			}
			conditionsRequest.AreaConditions = areaConditions
		}

		// Method Conditions
		if conditionMap["method_conditions"] != nil {
			methodConditions := make([]*waapCustomizerule.RequestMethodCondition, 0)
			for _, methodCondition := range conditionMap["method_conditions"].([]interface{}) {
				methodConditionMap := methodCondition.(map[string]interface{})
				matchType := methodConditionMap["match_type"].(string)
				requestMethodInterface := methodConditionMap["request_method"].([]interface{})
				requestMethod := make([]*string, len(requestMethodInterface))
				for i, v := range requestMethodInterface {
					str := v.(string)
					requestMethod[i] = &str
				}
				methodCondition := &waapCustomizerule.RequestMethodCondition{
					MatchType:     &matchType,
					RequestMethod: requestMethod,
				}
				methodConditions = append(methodConditions, methodCondition)
			}
			conditionsRequest.MethodConditions = methodConditions
		}

		// JA3 Conditions
		if conditionMap["ja3_conditions"] != nil {
			ja3Conditions := make([]*waapCustomizerule.Ja3Condition, 0)
			for _, ja3Condition := range conditionMap["ja3_conditions"].([]interface{}) {
				ja3ConditionMap := ja3Condition.(map[string]interface{})
				matchType := ja3ConditionMap["match_type"].(string)
				ja3Interface := ja3ConditionMap["ja3_list"].([]interface{})
				ja3 := make([]*string, len(ja3Interface))
				for i, v := range ja3Interface {
					str := v.(string)
					ja3[i] = &str
				}
				ja3Condition := &waapCustomizerule.Ja3Condition{
					MatchType: &matchType,
					Ja3List:   ja3,
				}
				ja3Conditions = append(ja3Conditions, ja3Condition)
			}
			conditionsRequest.Ja3Conditions = ja3Conditions
		}

		// JA4 Conditions
		if conditionMap["ja4_conditions"] != nil {
			ja4Conditions := make([]*waapCustomizerule.Ja4Condition, 0)
			for _, ja4Condition := range conditionMap["ja4_conditions"].([]interface{}) {
				ja4ConditionMap := ja4Condition.(map[string]interface{})
				matchType := ja4ConditionMap["match_type"].(string)
				ja4Interface := ja4ConditionMap["ja4_list"].([]interface{})
				ja4 := make([]*string, len(ja4Interface))
				for i, v := range ja4Interface {
					str := v.(string)
					ja4[i] = &str
				}
				ja4Condition := &waapCustomizerule.Ja4Condition{
					MatchType: &matchType,
					Ja4List:   ja4,
				}
				ja4Conditions = append(ja4Conditions, ja4Condition)
			}
			conditionsRequest.Ja4Conditions = ja4Conditions
		}
	}
	request.Condition = conditionsRequest

	var response *waapCustomizerule.UpdateCustomRuleResponse
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseWaapCustomizeruleClient().UpdateCustomRule(request)
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
	log.Printf("resource.wangsu_waap_customize_rule.update success")
	return nil
}

func resourceWaapCustomizeRuleDelete(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_waap_customize_rule.delete")

	var response *waapCustomizerule.DeleteCustomRuleResponse
	var err error
	var diags diag.Diagnostics
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		id := data.Id()
		request := &waapCustomizerule.DeleteCustomRuleRequest{
			IdList: []*string{&id},
		}
		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseWaapCustomizeruleClient().DeleteCustomRule(request)
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
