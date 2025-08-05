package pre_deploy

import (
	"context"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wangsuCommon "github.com/wangsu-api/terraform-provider-wangsu/wangsu/common"
	preDeploy "github.com/wangsu-api/wangsu-sdk-go/wangsu/waap/predeploy"
	"log"
	"time"
)

func ResourceWaapPreDeployDDoSProtection() *schema.Resource {
	return &schema.Resource{

		CreateContext: resourceWaapPreDeployDDoSProtectionCreate,
		ReadContext:   resourceWaapPreDeployDDoSProtectionRead,
		UpdateContext: resourceWaapPreDeployDDoSProtectionCreate,
		DeleteContext: resourceWaapPreDeployDDoSProtectionRead,

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
			"ddos_protect_switch": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Basic switch/mode information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"l7_ddos_switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Layer 7 HTTP DDoS protection switch.<br/>ON: Enable.<br/>OFF: Disable.",
						},
						"protect_mode": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Protection mode.<br/>LOOSE: loose.<br/>MODERATE: moderate.<br/>STRICT: strict.<br/>",
						},
						"inner_switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Built in protective switch.<br/>ON: Enable.<br/>OFF: Disable.",
						},
					},
				},
			},
			"built_in_rules": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Built-In rules, unprovided rules will take effect according to current production configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "rule ID.",
						},
						"security_level": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Security level.<br/>DEFAULT_ENABLE: default enabled.<br/>ATTACK_ENABLE: enable during attack.<br/>BASE_CLOSE: basic off.<br/>CLOSE: permanently closed.",
						},
						"action": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Mode.<br/>BLOCK: Protect(Default).<br/>RR: Protect(Managed).<br/>LOG: Monitor.<br/>DENIED: Connection denied.",
						},
					},
				},
			},
		},
	}
}

func resourceWaapPreDeployDDoSProtectionCreate(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("resource.wangsu_pre_deploy_ddos_protection.create")

	var diags diag.Diagnostics

	// Prepare the request object
	request := &preDeploy.PreDeployDDoSProtectionConfigurationRequest{}
	if domain, ok := data.Get("domain").(string); ok {
		request.Domain = &domain
	}

	// Parse ddos_protect_switch
	if ddosProtectSwitchList, ok := data.Get("ddos_protect_switch").([]interface{}); ok && len(ddosProtectSwitchList) > 0 {
		ddosProtectSwitchMap := ddosProtectSwitchList[0].(map[string]interface{})
		request.DdosProtectSwitch = &preDeploy.DdosProtectSwitch{
			L7DdosSwitch: tea.String(ddosProtectSwitchMap["l7_ddos_switch"].(string)),
			ProtectMode:  tea.String(ddosProtectSwitchMap["protect_mode"].(string)),
			InnerSwitch:  tea.String(ddosProtectSwitchMap["inner_switch"].(string)),
		}
	}

	// Parse built_in_rules
	if builtInRules, ok := data.Get("built_in_rules").([]interface{}); ok {
		builtInRulesRequest := make([]*preDeploy.DDoSBuiltInRule, len(builtInRules))
		for i, rule := range builtInRules {
			ruleMap := rule.(map[string]interface{})
			builtInRulesRequest[i] = &preDeploy.DDoSBuiltInRule{
				RuleId:        tea.String(ruleMap["rule_id"].(string)),
				SecurityLevel: tea.String(ruleMap["security_level"].(string)),
				Action:        tea.String(ruleMap["action"].(string)),
			}
		}
		request.BuiltInRules = builtInRulesRequest
	}

	// Send the request
	var response *preDeploy.PreDeployDDoSProtectionConfigurationResponse
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseWaapPreDeployClient().PreDeployDDoSProtectionConfiguration(request)
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

func resourceWaapPreDeployDDoSProtectionRead(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("resource.wangsu_pre_deploy_ddos_protection.read")

	data.SetId("")
	return nil
}
