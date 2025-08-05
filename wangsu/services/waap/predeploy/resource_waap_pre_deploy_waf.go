package pre_deploy

import (
	"context"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wangsuCommon "github.com/wangsu-api/terraform-provider-wangsu/wangsu/common"
	"github.com/wangsu-api/terraform-provider-wangsu/wangsu/services/waap"
	preDeploy "github.com/wangsu-api/wangsu-sdk-go/wangsu/waap/predeploy"
	"log"
	"time"
)

func ResourceWaapPreDeployWAF() *schema.Resource {
	return &schema.Resource{

		CreateContext: resourceWaapPreDeployWAFCreate,
		ReadContext:   resourceWaapPreDeployWAFRead,
		UpdateContext: resourceWaapPreDeployWAFCreate,
		DeleteContext: resourceWaapPreDeployWAFRead,

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
				Description: "Domain.",
			},
			"config_switch": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Policy switch.<br/>ON: Enable.<br/>OFF: Disable.",
			},
			"conf_basic": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Basic configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"defend_mode": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Protection Mode.<br/>BLOCK: Block the attack request directly.<br/>LOG: Only log the attack request without blocking it.",
						},
						"rule_update_mode": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Ruleset Mode.<br/>MANUAL: Check Ruleset update and all Recommendations on the Console, decide to apply them or not, all of these must be done by yourself manually.<br/>AUTO: Automatically upgrade the Ruleset to the latest version and apply the Recommendations learned from your website traffic to Exception, which can keep your website with high-level security anytime.",
						},
					},
				},
			},
			"rule_list": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Rule list, unprovided rules will take effect according to current production configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_id": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "WAF rule ID.",
						},
						"mode": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Rule actions.<br/>BLOCK: Deny request by a default 403 response.<br/>LOG: Log request and continue further detections.<br/>OFF: Select if you do not a policy take effect.",
						},
						"exception_list": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Rule exceptions.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Matching conditions.<br/>ip: IP<br/>path: Path<br/>uri: URI<br/>urlParamName: URI Parameter Name<br/>urlParamValue: URI Parameter Value<br/>userAgent: User Agent<br/>httpHeaderName: Request Header Name<br/>httpHeaderValue: Request Header Value<br/>cookie: Cookie<br/>body: Body<br/>bodyParamName: Body Parameter Name<br/>bodyParamValue: Body Parameter Value",
									},
									"match_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Match type,IP can only be EQUAL.<br/>EQUAL: Equal<br/>CONTAIN: Contains<br/>REGEX: Regular match",
									},
									"content_list": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "Rule exceptions.<br/>When matchType=EQUAL, case-sensitive, path and uri must start with \"/\", and body can only pass one value;<br/>When matchType=REGEX, only one value can be passed.",
										Elem: &schema.Schema{
											Type:        schema.TypeString,
											Description: "Rule exceptions.<br/>When matchType=EQUAL, case-sensitive, path and uri must start with \"/\", and body can only pass one value;<br/>When matchType=REGEX, only one value can be passed.",
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

func resourceWaapPreDeployWAFCreate(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("resource.wangsu_pre_deploy_waf.create")

	var diags diag.Diagnostics

	// Prepare the request object
	request := &preDeploy.PreDeployWAFConfigurationRequest{}
	if domain, ok := data.Get("domain").(string); ok {
		request.Domain = &domain
	}
	if configSwitch, ok := data.Get("config_switch").(string); ok {
		request.ConfigSwitch = &configSwitch
	}

	// Parse conf_basic
	if confBasicList, ok := data.Get("conf_basic").([]interface{}); ok && len(confBasicList) > 0 {
		confBasicMap := confBasicList[0].(map[string]interface{})
		request.ConfBasic = &preDeploy.WAFConfBasic{
			DefendMode:     tea.String(confBasicMap["defend_mode"].(string)),
			RuleUpdateMode: tea.String(confBasicMap["rule_update_mode"].(string)),
		}
	}

	// Parse rule_list
	if ruleList, ok := data.Get("rule_list").([]interface{}); ok {
		ruleListRequest := make([]*preDeploy.WAFRule, len(ruleList))
		for i, rule := range ruleList {
			ruleMap := rule.(map[string]interface{})
			ruleRequest := &preDeploy.WAFRule{}

			if v, ok := ruleMap["rule_id"]; ok {
				ruleRequest.SetRuleId(v.(int))
			}
			if v, ok := ruleMap["mode"]; ok {
				ruleRequest.SetMode(v.(string))
			}

			// Parse exception_list
			if exceptionList, ok := ruleMap["exception_list"].([]interface{}); ok {
				exceptionListRequest := make([]*preDeploy.WAFRuleListException, len(exceptionList))
				for j, exception := range exceptionList {
					exceptionMap := exception.(map[string]interface{})
					exceptionRequest := &preDeploy.WAFRuleListException{
						Type:        tea.String(exceptionMap["type"].(string)),
						MatchType:   tea.String(exceptionMap["match_type"].(string)),
						ContentList: waap.ConvertToStringSlice(exceptionMap["content_list"].([]interface{})),
					}
					exceptionListRequest[j] = exceptionRequest
				}
				ruleRequest.ExceptionList = exceptionListRequest
			}
			ruleListRequest[i] = ruleRequest
		}
		request.RuleList = ruleListRequest
	}

	// Send the request
	var response *preDeploy.PreDeployWAFConfigurationResponse
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseWaapPreDeployClient().PreDeployWAFConfiguration(request)
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

	// Poll for deployment result
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

func resourceWaapPreDeployWAFRead(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("resource.wangsu_pre_deploy_waf.read")

	data.SetId("")
	return nil
}
