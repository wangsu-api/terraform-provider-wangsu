package intelligence

import (
	"context"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wangsuCommon "github.com/wangsu-api/terraform-provider-wangsu/wangsu/common"
	securityPolicy "github.com/wangsu-api/wangsu-sdk-go/wangsu/securitypolicy"
	"log"
	"time"
)

func DataSourceWaapThreatIntelligence() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWaapThreatIntelligenceRead,

		Schema: map[string]*schema.Schema{
			"domain_list": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Hostname list.",
			},
			"data": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Data.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Domain.",
						},
						"config_list": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Configuration list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Category ID.",
									},
									"info_cate": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Category.",
									},
									"second_cate": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "SubCategory.",
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
				},
			},
		},
	}
}

func dataSourceWaapThreatIntelligenceRead(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("data_source.waap_threat_intelligence.read")

	var diags diag.Diagnostics
	request := &securityPolicy.GetThreatIntelligenceDomainConfigRequest{}

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
	var response *securityPolicy.GetThreatIntelligenceDomainConfigResponse
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
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
	if response == nil || response.Data == nil {
		data.SetId("")
		return nil
	}

	// Parse response data
	groupByDomain := make(map[string][]map[string]interface{})
	for _, item := range response.Data {
		domain := tea.StringValue(item.Domain)
		parsedItem := map[string]interface{}{
			"id":          tea.StringValue(item.Id),
			"info_cate":   tea.StringValue(item.InfoCate),
			"second_cate": tea.StringValue(item.SecondCate),
			"action":      tea.StringValue(item.Action),
		}

		groupByDomain[domain] = append(groupByDomain[domain], parsedItem)
	}

	dataList := make([]interface{}, 0, len(groupByDomain))
	ids := make([]string, 0, len(groupByDomain))

	for domain, configList := range groupByDomain {
		dataList = append(dataList, map[string]interface{}{
			"domain":      domain,
			"config_list": configList,
		})
		ids = append(ids, domain)
	}

	// Set data and ID
	if err := data.Set("data", dataList); err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}
	data.SetId(wangsuCommon.DataResourceIdsHash(ids))
	return diags
}
