package intelligence

import (
	"context"
	"errors"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wangsuCommon "github.com/wangsu-api/terraform-provider-wangsu/wangsu/common"
	securityPolicy "github.com/wangsu-api/wangsu-sdk-go/wangsu/securitypolicy"
	"log"
	"time"
)

func ResourceWaapThreatIntelligence() *schema.Resource {
	return &schema.Resource{

		CreateContext: resourceWaapThreatIntelligenceCreate,
		UpdateContext: resourceWaapThreatIntelligenceUpdate,
		ReadContext:   resourceWaapThreatIntelligenceRead,
		DeleteContext: resourceWaapThreatIntelligenceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Domain.",
			},
			"config_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "Configuration list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Category ID.",
						},
						"action": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Action.<br/>NO_USE: Not Used<br/>LOG: Log<br/>BLOCK: Deny",
						},
					},
				},
			},
		},
	}
}

func resourceWaapThreatIntelligenceRead(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.waap_threat_intelligence.read")

	var diags diag.Diagnostics
	// 使用导入的 ID 设置资源 ID
	if data.Id() != "" {
		_ = data.Set("domain", data.Id())
	} else if domain, ok := data.GetOk("domain"); ok {
		data.SetId(domain.(string))
	}

	// Make API call
	request := &securityPolicy.GetThreatIntelligenceDomainConfigRequest{}
	var response *securityPolicy.GetThreatIntelligenceDomainConfigResponse
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		domain := data.Id()
		request.SetDomainList([]*string{&domain})
		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseSecurityPolicyClient().GetThreatIntelligenceDomainConfig(request)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}
	if response == nil || response.Data == nil || len(response.Data) == 0 {
		return nil
	}

	// Parse response data
	// 提取 main.tf 中的 id 列表
	mainTfIds := make(map[string]struct{})
	if configList, ok := data.GetOk("config_list"); ok {
		if configListInterface, ok := configList.([]interface{}); ok {
			for _, configItem := range configListInterface {
				if configMap, ok := configItem.(map[string]interface{}); ok {
					if id, ok := configMap["id"].(string); ok {
						mainTfIds[id] = struct{}{}
					}
				}
			}
		}
	}
	// 过滤远程规则，只保留 main.tf 中的 rule_id
	filteredDataList := make([]map[string]interface{}, 0)
	for _, item := range response.Data {
		var id = tea.StringValue(item.Id)
		if _, exists := mainTfIds[id]; len(mainTfIds) == 0 || exists {
			parsedItem := map[string]interface{}{
				"id":     id,
				"action": tea.StringValue(item.Action),
			}
			filteredDataList = append(filteredDataList, parsedItem)
		}
	}

	// 将数据设置到 ResourceData 中
	if err := data.Set("config_list", filteredDataList); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceWaapThreatIntelligenceUpdate(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.waap_threat_intelligence.update")

	var diags diag.Diagnostics
	// 把 domain 强制刷回旧值，否则会有权限问题
	if data.HasChange("domain") {
		oldDomain, _ := data.GetChange("domain")
		_ = data.Set("domain", oldDomain)
		err := errors.New("Domain cannot be changed.")
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	// Prepare request data
	request := &securityPolicy.UpdateThreatIntelligenceDomainConfigRequest{}
	domain := data.Id()
	request.SetDomain(domain)

	configList := data.Get("config_list").([]interface{})
	updateData := make([]*securityPolicy.UpdateThreatIntelligenceDomainConfigRequestConfigList, len(configList))
	for i, item := range configList {
		config := item.(map[string]interface{})
		updateData[i] = &securityPolicy.UpdateThreatIntelligenceDomainConfigRequestConfigList{
			Id:     tea.String(config["id"].(string)),
			Action: tea.String(config["action"].(string)),
		}
	}
	request.SetConfigList(updateData)

	// Make API call
	var response *securityPolicy.UpdateThreatIntelligenceDomainConfigResponse
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseSecurityPolicyClient().UpdateThreatIntelligenceDomainConfig(request)
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

	return diags
}

func resourceWaapThreatIntelligenceCreate(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.waap_threat_intelligence.create")
	var diags diag.Diagnostics
	err := errors.New("please execute \"terraform import <resource>.<resource_name> <domain>\" before executing \"terraform apply\" for the first time")
	diags = append(diags, diag.FromErr(err)...)
	return diags
}

func resourceWaapThreatIntelligenceDelete(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.waap_threat_intelligence.delete")
	// 清空本地 state
	data.SetId("")
	return nil
}
