package share_customizebot

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wangsuCommon "github.com/wangsu-api/terraform-provider-wangsu/wangsu/common"
	waapShareCustomizeBot "github.com/wangsu-api/wangsu-sdk-go/wangsu/waap/share-customizebot"
	"log"
	"time"
)

func DataSourceShareCustomizeBot() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceShareCustomizeBotRead,
		Schema: map[string]*schema.Schema{
			"bot_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Bot name,fuzzy query.",
			},
			"data": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Data.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule ID.",
						},
						"bot_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule name.",
						},
						"bot_description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description.",
						},
						"bot_act": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Actions.<br/>BLOCK: block<br/>LOG: log<br/>ACCEPT: release",
						},
						"condition_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Matching conditions.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"condition_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Matching condition name. <br/>IP_IPS: IP/IP segment <br/>JA3: JA3 Fingerprint<br/>JA4: JA4 Fingerprint<br/>UA: User-agent <br/>HEADER: Request Header <br/>ASN: AS Number <br/>CLIENT_GROUP: Client Group <br/>PUBLIC_BOT: Public Bots",
									},
									"condition_value_list": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Condition value list.",
										Elem: &schema.Schema{
											Type:        schema.TypeString,
											Description: "Condition value.",
										},
									},
									"condition_func": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Matching condition function.<br/>EQUAL: Equals<br/>NOT_EQUAL: Does not equal<br/>CONTAIN: Contains<br/>NOT_CONTAIN: Does not contain<br/>NONE: Empty or non-existent<br/>REGEX: Regex match<br/>NOT_REGEX: Does not match regex<br/>START_WITH: Starts with<br/>END_WITH: Ends with<br/>WILDCARD: Wildcard matches, * represents zero or more arbitrary characters, ? represents any single character<br/>NOT_WILDCARD: Wildcard does not match, * represents zero or more arbitrary characters, ? represents any single character",
									},
									"condition_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Request header name.",
									},
								},
							},
						},
						"rela_domain_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of associated hostnames.",
							Elem: &schema.Schema{
								Type:        schema.TypeString,
								Description: "associated hostnam.",
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceShareCustomizeBotRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("data_source.wangsu_waap_share_customizebot.read")

	var response *waapShareCustomizeBot.ListShareCustomizeBotsResponse
	var err error
	var diags diag.Diagnostics
	request := &waapShareCustomizeBot.ListShareCustomizeBotsRequest{}
	if v, ok := data.GetOk("bot_name"); ok {
		request.SetBotName(v.(string))
	}
	err = resource.RetryContext(ctx, time.Duration(2)*time.Minute, func() *resource.RetryError {
		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseWaapShareCustomizeBotClient().GetList(request)
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
			conditions := make([]map[string]interface{}, len(item.ConditionList))
			if item.ConditionList != nil {
				for i, cond := range item.ConditionList {
					conditions[i] = map[string]interface{}{
						"condition_name": *cond.ConditionName,
						"condition_key": func() string {
							if cond.ConditionKey != nil {
								return *cond.ConditionKey
							}
							return ""
						}(),
						"condition_func":       *cond.ConditionFunc,
						"condition_value_list": cond.ConditionValueList,
					}
				}
			}
			itemList = append(itemList, map[string]interface{}{
				"id":               item.Id,
				"rela_domain_list": item.RelaDomainList,
				"bot_name":         item.BotName,
				"bot_description":  item.BotDescription,
				"bot_act":          item.BotAct,
				"condition_list":   conditions,
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
