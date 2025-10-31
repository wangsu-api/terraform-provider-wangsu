package bot_scene_whitelist

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wangsuCommon "github.com/wangsu-api/terraform-provider-wangsu/wangsu/common"
	waapBotSceneWhitelist "github.com/wangsu-api/wangsu-sdk-go/wangsu/waap/bot-scene-whitelist"
	"log"
	"time"
)

func DataSourceBotSceneWhitelist() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBotSceneWhitelistRead,
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
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "ID.",
						},
						"domain": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Hostname.",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Whitelist name.",
						},
						"description": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Description.",
						},
						"conditions": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "List of matching conditions.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"match_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Matching condition name.<br/>IP_IPS: IP/CIDR<br/>PATH: Path<br/>URI: Path with parameters<br/>HEADER: Request Header<br/>UA: User Agent<br/>REQUEST_METHOD: Request Method<br/>REFERER: Referer",
									},
									"match_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "When matchName is IP_IPS, maximum 300 IP/CIDR in match value list, the optional value of matchType is:<br/>EQUAL: Equals<br/>NOT_EQUAL: Does not equal<br/>When matchName is a URI, the optional value ...", // 内容较长省略
									},
									"match_key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Mathing key, this value is not empty and valid only when matchName=HEADER.<br/>Maximum 100 characters, case insensitive.",
									},
									"match_value_list": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "List of matching values.",
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

func dataSourceBotSceneWhitelistRead(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("data_source.waap_bot_scene_whitelist.read")
	var response *waapBotSceneWhitelist.ListSpecificClientTrafficBypassResponse
	var err error
	var diags diag.Diagnostics
	request := &waapBotSceneWhitelist.ListSpecificClientTrafficBypassRequest{}
	if v, ok := data.GetOk("domain_list"); ok {
		targetDomainsList := v.([]interface{})
		targetDomainsStr := make([]*string, len(targetDomainsList))
		for i, v := range targetDomainsList {
			str := v.(string)
			targetDomainsStr[i] = &str
		}
		request.SetDomainList(targetDomainsStr)
	}
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseWaapBotSceneWhiteListClient().GetList(request)
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
		ids := make([]string, 0, len(response.Data))
		itemList := make([]interface{}, 0)
		for _, item := range response.Data {
			conditions := make([]map[string]interface{}, len(item.Conditions))
			if item.Conditions != nil {
				for i, cond := range item.Conditions {
					conditions[i] = map[string]interface{}{
						"match_name": *cond.MatchName,
						"match_type": *cond.MatchType,
						"match_key": func() string {
							if cond.MatchKey != nil {
								return *cond.MatchKey
							}
							return ""
						}(),
						"match_value_list": cond.MatchValueList,
					}
				}
			}
			itemList = append(itemList, map[string]interface{}{
				"id":          item.Id,
				"domain":      item.Domain,
				"name":        item.Name,
				"description": item.Description,
				"conditions":  conditions,
			})
			ids = append(ids, *item.Id)
		}
		if err := data.Set("data", itemList); err != nil {
			return diag.FromErr(fmt.Errorf("error setting data for resource: %s", err))
		}
		data.SetId(wangsuCommon.DataResourceIdsHash(ids))
	}
	return diags
}
