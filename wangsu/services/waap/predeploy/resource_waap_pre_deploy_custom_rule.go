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

func ResourceWaapPreDeployCustomRule() *schema.Resource {
	return &schema.Resource{

		CreateContext: resourceWaapPreDeployCustomRuleCreate,
		ReadContext:   resourceWaapPreDeployCustomRuleRead,
		UpdateContext: resourceWaapPreDeployCustomRuleCreate,
		DeleteContext: resourceWaapPreDeployCustomRuleRead,

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
														Description: "URI.When match type is EQUAL/NOT_EQUAL/START_WITH/END_WITH, uri needs to start with \"/\", and includes parameters.When the match type is REGEX/NOT_REGEX, only one value is allowed.Example: /test.html?id=1.",
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
														Description: "Path.When match type is EQUAL/NOT_EQUAL/START_WITH/END_WITH, path needs to start with \"/\", and no parameters.When the match type is REGEX/NOT_REGEX, only one value is allowed.Example: /test.html.",
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
														Type:        schema.TypeString,
														Description: "Param value.",
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
														Description: "User agent.When the match type is REGEX/NOT_REGEX, only one value is allowed.Example: go-Http-client/1.1.",
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
														Description: "Header value.When the match type is REGEX/NOT_REGEX, only one value is allowed.",
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
												"referer": {
													Type:        schema.TypeList,
													Required:    true,
													Description: "Referer.<br/>When the match type is REGEX/NOT_REGEX, only one value is allowed.<br/>Example: http://test.com.",
													Elem: &schema.Schema{
														Type:        schema.TypeString,
														Description: "Referer.When the match type is REGEX/NOT_REGEX, only one value is allowed.Example: http://test.com.",
													},
												},
												"match_type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Match type.<br/>EQUAL: Equals, referer case sensitive<br/>NOT_EQUAL: Does not equal, referer case sensitive<br/>CONTAIN: Contains, referer case insensitive<br/>NOT_CONTAIN: Does not Contains, referer case insensitive<br/>NONE:Empty or non-existent<br/>REGEX: Regex match, referer case insensitive<br/>NOT_REGEX: Regular does not match, referer case insensitive<br/>START_WITH: Starts with, referer case insensitive<br/>END_WITH: Ends with, referer case insensitive<br/>WILDCARD: Wildcard matches, referer case insensitive,* represents zero or more arbitrary characters, ? represents any single characte<br/>NOT_WILDCARD: Wildcard does not match, referer case insensitive,* represents zero or more arbitrary characters, ? represents any single character",
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
														Description: "JA3 Fingerprint List, maximum 300 JA3 Fingerprint.When the match type is EQUAL/NOT_EQUAL, each item's character length must be 32 and can only include numbers and lowercase letters.",
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
													Description: "Match type. <br/>EQUAL: Equals<br/>NOT_EQUAL: Does not equal<br/>CONTAIN: Contains<br/>NOT_CONTAIN: Does not Contains<br/>START_WITH: Starts with<br/>END_WITH: Ends with<br/>WILDCARD: Wildcard matches, ** represents zero or more arbitrary characters, ? represents any single character<br/>NOT_WILDCARD: Wildcard does not match, ** represents zero or more arbitrary characters, ? represents any single character",
												},
												"ja4_list": {
													Type:        schema.TypeList,
													Required:    true,
													Description: "JA4 Fingerprint List, maximum 300 JA4 Fingerprint.<br/>When the match type is EQUAL/NOT_EQUAL, each item's format must be 10 characters + 12 characters + 12 characters, separated by underscores, and can only include underscores, numbers, and lowercase letters.<br/>When the match type is CONTAIN/NOT_CONTAIN/START_WITH/END_WITH, each item is only allowed to include underscores, numbers, and lowercase letters.<br/>When the match type is WILDCARD/NOT_WILDCARD, each item, aside from  ** and ?, is only allowed to include underscores, numbers, and lowercase letters.",
													Elem: &schema.Schema{
														Type:        schema.TypeString,
														Description: "JA4 Fingerprint List, maximum 300 JA4 Fingerprint.When the match type is EQUAL/NOT_EQUAL, each item's format must be 10 characters + 12 characters + 12 characters, separated by underscores, and can only include underscores, numbers, and lowercase letters.When the match type is CONTAIN/NOT_CONTAIN/START_WITH/END_WITH, each item is only allowed to include underscores, numbers, and lowercase letters.When the match type is WILDCARD/NOT_WILDCARD, each item, aside from  ** and ?, is only allowed to include underscores, numbers, and lowercase letters.",
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

func resourceWaapPreDeployCustomRuleCreate(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("resource.wangsu_pre_deploy_custom_rule.create")

	var diags diag.Diagnostics

	request := &preDeploy.PreDeployCustomRuleConfigurationRequest{}
	if domain, ok := data.Get("domain").(string); ok {
		request.Domain = &domain
	}
	if configSwitch, ok := data.Get("config_switch").(string); ok {
		request.ConfigSwitch = &configSwitch
	}

	ruleList := data.Get("rule_list").([]interface{})
	ruleListRequest := make([]*preDeploy.CustomRule, len(ruleList))
	for i, rule := range ruleList {
		ruleMap := rule.(map[string]interface{})
		ruleRequest := &preDeploy.CustomRule{}

		if ruleName, ok := ruleMap["rule_name"].(string); ok && ruleName != "" {
			ruleRequest.RuleName = &ruleName
		}
		if description, ok := ruleMap["description"].(string); ok && description != "" {
			ruleRequest.Description = &description
		}
		if scene, ok := ruleMap["scene"].(string); ok && scene != "" {
			ruleRequest.Scene = &scene
		}
		if apiId, ok := ruleMap["api_id"].(string); ok && apiId != "" {
			ruleRequest.ApiId = &apiId
		}
		if act, ok := ruleMap["act"].(string); ok && act != "" {
			ruleRequest.Act = &act
		}

		// Parse conditions
		if ruleMap["condition"] != nil {
			conditions := ruleMap["condition"].([]interface{})
			conditionsRequest := &preDeploy.CustomRuleCondition{}
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
			}
			ruleRequest.Condition = conditionsRequest
		}
		ruleListRequest[i] = ruleRequest
	}
	request.RuleList = ruleListRequest

	var response *preDeploy.PreDeployCustomRuleConfigurationResponse
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseWaapPreDeployClient().PreDeployCustomRuleConfiguration(request)
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

func resourceWaapPreDeployCustomRuleRead(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("resource.wangsu_pre_deploy_custom_rule.read")

	data.SetId("")
	return nil
}
