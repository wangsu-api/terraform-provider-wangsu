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

func ResourceWaapPreDeployWhitelist() *schema.Resource {
	return &schema.Resource{

		CreateContext: resourceWaapPreDeployWhitelistCreate,
		ReadContext:   resourceWaapPreDeployWhitelistRead,
		UpdateContext: resourceWaapPreDeployWhitelistCreate,
		DeleteContext: resourceWaapPreDeployWhitelistRead,

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
							Description: "Rule name, maximum 50 characters.<br/>does not support # and & .",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description, maximum 200 characters.",
						},
						"conditions": {
							Type:        schema.TypeList,
							Required:    true,
							MaxItems:    1,
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
										Description: "Path match conditions, match type cannot be repeated.",
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
										Description: "URI match conditions, match type cannot be repeated.",
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
										Description: "User agent match conditions, match type cannot be repeated.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"match_type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Match type.<br/>EQUAL: Equals, user agent case sensitive<br/>NOT_EQUAL: Does not equal, user agent case sensitive<br/>CONTAIN: Contains, user agent case insensitive<br/>NOT_CONTAIN: Does not Contains, user agent case insensitive<br/>REGEX: Regex match, user agent case insensitive<br/>NOT_REGEX: Regular does not match, user agent case insensitive<br/>START_WITH: Starts with, user agent case insensitive<br/>END_WITH: Ends with, user agent case insensitive<br/>WILDCARD: Wildcard matches, user agent case insensitive,* represents zero or more arbitrary characters, ? represents any single character<br/>NOT_WILDCARD: Wildcard does not match, user agent case insensitive,* represents zero or more arbitrary characters, ? represents any single character",
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
										Description: "Referer match conditions, match type cannot be repeated.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"match_type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Match type.<br/>EQUAL: Equals, referer case sensitive<br/>NOT_EQUAL: Does not equal, referer case sensitive<br/>CONTAIN: Contains, referer case insensitive<br/>NOT_CONTAIN: Does not Contains, referer case insensitive<br/>REGEX: Regex match, referer case insensitive<br/>NOT_REGEX: Regular does not match, referer case insensitive<br/>START_WITH: Starts with, referer case insensitive<br/>END_WITH: Ends with, referer case insensitive<br/>WILDCARD: Wildcard matches, referer case insensitive,* represents zero or more arbitrary characters, ? represents any single character<br/>NOT_WILDCARD: Wildcard does not match, referer case insensitive,* represents zero or more arbitrary characters, ? represents any single character",
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
										Description: "Request header match conditions.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"match_type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Match type.<br/>EQUAL: Equals, request header values case sensitive<br/>NOT_EQUAL: Does not equal, request header values case sensitive<br/>CONTAIN: Contains, request header values case insensitive<br/>NOT_CONTAIN: Does not Contains, request header values case insensitive<br/>REGEX: Regex match, request header values case insensitive<br/>NOT_REGEX: Regular does not match, request header values case insensitive<br/>START_WITH: Starts with, request header values case insensitive<br/>END_WITH: Ends with, request header values case insensitive<br/>WILDCARD: Wildcard matches, request header values case insensitive,* represents zero or more arbitrary characters, ? represents any single character<br/>NOT_WILDCARD: Wildcard does not match, request header values case insensitive,* represents zero or more arbitrary characters, ? represents any single character",
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
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceWaapPreDeployWhitelistCreate(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("resource.wangsu_pre_deploy_whitelist.create")

	var diags diag.Diagnostics

	request := &preDeploy.PreDeployWhitelistConfigurationRequest{}
	if domain, ok := data.Get("domain").(string); ok {
		request.Domain = &domain
	}
	if configSwitch, ok := data.Get("config_switch").(string); ok {
		request.ConfigSwitch = &configSwitch
	}

	ruleList := data.Get("rule_list").([]interface{})
	ruleListRequest := make([]*preDeploy.WhitelistRule, len(ruleList))
	for i, rule := range ruleList {
		ruleMap := rule.(map[string]interface{})
		ruleRequest := &preDeploy.WhitelistRule{}

		if ruleName, ok := ruleMap["rule_name"].(string); ok && ruleName != "" {
			ruleRequest.RuleName = &ruleName
		}
		if description, ok := ruleMap["description"].(string); ok && description != "" {
			ruleRequest.Description = &description
		}

		// Parse conditions
		if ruleMap["conditions"] != nil {
			conditions := ruleMap["conditions"].([]interface{})
			conditionsRequest := &preDeploy.WhitelistRuleConditions{}
			for _, condition := range conditions {
				conditionMap := condition.(map[string]interface{})

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
			}
			ruleRequest.Conditions = conditionsRequest
		}
		ruleListRequest[i] = ruleRequest
	}
	request.RuleList = ruleListRequest

	var response *preDeploy.PreDeployWhitelistConfigurationResponse
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseWaapPreDeployClient().PreDeployWhitelistConfiguration(request)
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

func resourceWaapPreDeployWhitelistRead(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("resource.wangsu_pre_deploy_whitelist.read")

	data.SetId("")
	return nil
}
