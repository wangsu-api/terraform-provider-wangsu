package ddosprotection

import (
	"context"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wangsuCommon "github.com/wangsu-api/terraform-provider-wangsu/wangsu/common"
	waapDDoSProtection "github.com/wangsu-api/wangsu-sdk-go/wangsu/waap/ddosprotection"
	"log"
	"time"
)

func DataSourceWaapDDoSProtection() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWaapDDoSProtectionRead,

		Schema: map[string]*schema.Schema{
			"domain_list": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Hostname list.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"data": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Data.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"array": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Array.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
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
												"l4_ddos_switch": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "3/4 layer protection switch.<br/>ON: Enable.<br/>OFF: Disable.",
												},
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
												"ai_switch": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "AI intelligent protection switch.<br/>ON: Enable.<br/>OFF: Disable.<br/>",
												},
												"ai_action": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "AI rule action.<br/>BLOCK: BLOCK.<br/>LOG: Monitor.<br/>RR: DDoS managed challenge.",
												},
											},
										},
									},
									"built_in_rules": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "Built-In rules, unprovided rules will take effect according to production configuration.",
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
													Description: "Mode.<br/>BLOCK: Protect(Default).<br/>RR: Protect(Managed).<br/>LOG: Monitor.<br/>DENIED: Connection denied.<br/>",
												},
												"rule_name_cn": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Chinese Rule Name.",
												},
												"rule_name_en": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "English Rule Name.",
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

func dataSourceWaapDDoSProtectionRead(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("data_source.wangsu_waap_ddos_protection.read")

	var diags diag.Diagnostics
	request := &waapDDoSProtection.GetDDoSProtectionConfigurationRequest{}

	// Retrieve domain_list from input
	if v, ok := data.GetOk("domain_list"); ok {
		domainList := v.([]interface{})
		domains := make([]*string, len(domainList))
		for i, domain := range domainList {
			domainStr := domain.(string)
			domains[i] = &domainStr
		}
		request.SetDomainList(domains)
	}

	// Make API call
	var response *waapDDoSProtection.GetDDoSProtectionConfigurationResponse
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseWaapDDoSProtectionClient().GetDDoSProtectionConfiguration(request)
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

	// Parse response data
	ddosDataList := make([]interface{}, len(response.Data.Array))
	ids := make([]string, len(response.Data.Array))
	for i, item := range response.Data.Array {
		parsedItem := map[string]interface{}{
			"domain": tea.StringValue(item.Domain),
			"ddos_protect_switch": []interface{}{
				map[string]interface{}{
					"l4_ddos_switch": tea.StringValue(item.DdosProtectSwitch.L4DdosSwitch),
					"l7_ddos_switch": tea.StringValue(item.DdosProtectSwitch.L7DdosSwitch),
					"protect_mode":   tea.StringValue(item.DdosProtectSwitch.ProtectMode),
					"inner_switch":   tea.StringValue(item.DdosProtectSwitch.InnerSwitch),
					"ai_switch":      tea.StringValue(item.DdosProtectSwitch.AiSwitch),
					"ai_action":      tea.StringValue(item.DdosProtectSwitch.AiAction),
				},
			},
			"built_in_rules": func() []interface{} {
				if item.BuiltInRules == nil {
					return nil
				}
				rules := make([]interface{}, len(item.BuiltInRules))
				for j, rule := range item.BuiltInRules {
					rules[j] = map[string]interface{}{
						"rule_id":        tea.StringValue(rule.RuleId),
						"security_level": tea.StringValue(rule.SecurityLevel),
						"action":         tea.StringValue(rule.Action),
						"rule_name_cn":   tea.StringValue(rule.RuleNameCn),
						"rule_name_en":   tea.StringValue(rule.RuleNameEn),
					}
				}
				return rules
			}(),
		}

		ddosDataList[i] = parsedItem
		ids[i] = tea.StringValue(item.Domain)
	}

	// Set data and ID
	if err := data.Set("data", []interface{}{
		map[string]interface{}{
			"array": ddosDataList,
		},
	}); err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}
	data.SetId(wangsuCommon.DataResourceIdsHash(ids))
	return diags
}
