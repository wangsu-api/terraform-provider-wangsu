package share_customizebot

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wangsuCommon "github.com/wangsu-api/terraform-provider-wangsu/wangsu/common"
	waapShareCustomizeBot "github.com/wangsu-api/wangsu-sdk-go/wangsu/waap/share-customizebot"
	"log"
	"time"
)

func ResourceWaapShareCustomizeBot() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWaapShareCustomizeBotCreate,
		ReadContext:   resourceWaapShareCustomizeBotRead,
		UpdateContext: resourceWaapShareCustomizeBotUpdate,
		DeleteContext: resourceWaapShareCustomizeBotDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Rule ID.",
			},
			"bot_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Rule name.",
			},
			"bot_description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description.",
			},
			"bot_act": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Actions.<br/>BLOCK: block<br/>LOG: log<br/>ACCEPT: release",
			},
			"condition_list": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Matching conditions.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"condition_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Matching condition name. <br/>IP_IPS: IP/IP segment <br/>JA3: JA3 Fingerprint<br/>JA4: JA4 Fingerprint<br/>UA: User-agent <br/>HEADER: Request Header <br/>ASN: AS Number <br/>CLIENT_GROUP: Client Group <br/>PUBLIC_BOT: Public Bots",
						},
						"condition_value_list": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Condition value list.",
							Elem: &schema.Schema{
								Type:        schema.TypeString,
								Description: "Condition value.",
							},
						},
						"condition_func": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Matching condition function.<br/>EQUAL: Equals<br/>NOT_EQUAL: Does not equal<br/>CONTAIN: Contains<br/>NOT_CONTAIN: Does not contain<br/>NONE: Empty or non-existent<br/>REGEX: Regex match<br/>NOT_REGEX: Does not match regex<br/>START_WITH: Starts with<br/>END_WITH: Ends with<br/>WILDCARD: Wildcard matches, * represents zero or more arbitrary characters, ? represents any single character<br/>NOT_WILDCARD: Wildcard does not match, * represents zero or more arbitrary characters, ? represents any single character",
						},
						"condition_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Request header name.",
						},
					},
				},
			},
			"rela_domain_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of associated hostnames.",
				Elem: &schema.Schema{
					Type:        schema.TypeString,
					Description: "associated hostnam.",
				},
			},
		},
	}
}

func resourceWaapShareCustomizeBotCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_waap_share_customize_bot.create")

	var diags diag.Diagnostics
	request := &waapShareCustomizeBot.AddShareCustomizeBotTFRequest{}

	// Set bot_name
	if botName, ok := data.Get("bot_name").(string); ok && botName != "" {
		request.BotName = &botName
	}

	// Set bot_description
	if botDescription, ok := data.Get("bot_description").(string); ok && botDescription != "" {
		request.BotDescription = &botDescription
	}

	// Set bot_act
	if botAct, ok := data.Get("bot_act").(string); ok && botAct != "" {
		request.BotAct = &botAct
	}

	// Set condition_list
	if conditions, ok := data.GetOk("condition_list"); ok {
		conditionList := conditions.([]interface{})
		addConditionList := make([]*waapShareCustomizeBot.AddShareCustomizeBotTFRequestConditionList, len(conditionList))
		for i, condition := range conditionList {
			conditionMap := condition.(map[string]interface{})
			addCondition := &waapShareCustomizeBot.AddShareCustomizeBotTFRequestConditionList{}
			if conditionName, ok := conditionMap["condition_name"].(string); ok && conditionName != "" {
				addCondition.ConditionName = &conditionName
			}
			if conditionFunc, ok := conditionMap["condition_func"].(string); ok && conditionFunc != "" {
				addCondition.ConditionFunc = &conditionFunc
			}
			if conditionKey, ok := conditionMap["condition_key"].(string); ok && conditionKey != "" {
				addCondition.ConditionKey = &conditionKey
			}
			if conditionValues, ok := conditionMap["condition_value_list"].([]interface{}); ok {
				conditionValueList := make([]*string, len(conditionValues))
				for i, v := range conditionValues {
					str := v.(string)
					conditionValueList[i] = &str
				}
				addCondition.ConditionValueList = conditionValueList
			}
			addConditionList[i] = addCondition
		}
		request.ConditionList = addConditionList
	}

	// Set rela_domain_list
	if domains, ok := data.Get("rela_domain_list").([]interface{}); ok {
		domainList := make([]*string, len(domains))
		for i, v := range domains {
			str := v.(string)
			domainList[i] = &str
		}
		request.RelaDomainList = domainList
	}

	var response *waapShareCustomizeBot.AddShareCustomizeBotTFResponse
	var err error
	err = resource.RetryContext(ctx, time.Duration(2)*time.Minute, func() *resource.RetryError {
		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseWaapShareCustomizeBotClient().Add(request)
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
	data.SetId(*response.Data.Id)

	return resourceWaapShareCustomizeBotRead(ctx, data, meta)
}

func resourceWaapShareCustomizeBotRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_waap_share_customize_bot.read")
	var response *waapShareCustomizeBot.ListShareCustomizeBotsResponse
	var err error
	var diags diag.Diagnostics
	err = resource.RetryContext(ctx, time.Duration(2)*time.Minute, func() *resource.RetryError {
		request := &waapShareCustomizeBot.ListShareCustomizeBotsRequest{}
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
		for _, item := range response.Data {
			// 只要对应id的数据
			if *item.Id != data.Id() {
				continue
			}
			_ = data.Set("id", item.Id)
			_ = data.Set("rela_domain_list", item.RelaDomainList)
			_ = data.Set("bot_name", item.BotName)
			_ = data.Set("bot_description", item.BotDescription)
			_ = data.Set("bot_act", item.BotAct)
			// 映射 conditions 数据
			conditions := make([]map[string]interface{}, len(item.ConditionList))
			for i, cond := range item.ConditionList {
				conditions[i] = map[string]interface{}{
					"condition_name": *cond.ConditionName,
					"condition_func": *cond.ConditionFunc,
					"condition_key": func() string {
						if cond.ConditionKey != nil {
							return *cond.ConditionKey
						}
						return ""
					}(),
					"condition_value_list": cond.ConditionValueList,
				}
			}
			_ = data.Set("condition_list", conditions)
		}
	}
	return nil
}

func resourceWaapShareCustomizeBotUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_waap_share_customize_bot.update")

	var diags diag.Diagnostics
	request := &waapShareCustomizeBot.UpdateShareCustomizeBotTFRequest{}
	if id, ok := data.Get("id").(string); ok && id != "" {
		request.Id = &id
	}
	if domains, ok := data.GetOk("rela_domain_list"); ok {
		domainsList := domains.([]interface{})
		domainsStr := make([]*string, len(domainsList))
		for i, v := range domainsList {
			str := v.(string)
			domainsStr[i] = &str
		}
		request.RelaDomainList = domainsStr
	}
	if ruleName, ok := data.Get("bot_name").(string); ok && ruleName != "" {
		request.BotName = &ruleName
	}
	if description, ok := data.Get("bot_description").(string); ok {
		request.BotDescription = &description
	}
	if act, ok := data.Get("bot_act").(string); ok && act != "" {
		request.BotAct = &act
	}
	if conditions, ok := data.GetOk("condition_list"); ok {
		conditionList := conditions.([]interface{})
		updateConditionList := make([]*waapShareCustomizeBot.UpdateShareCustomizeBotTFRequestConditionList, len(conditionList))

		for i, condition := range conditionList {
			condMap := condition.(map[string]interface{})
			updateCondition := &waapShareCustomizeBot.UpdateShareCustomizeBotTFRequestConditionList{}

			if conditionName, ok := condMap["condition_name"].(string); ok {
				updateCondition.ConditionName = &conditionName
			}
			if conditionFunc, ok := condMap["condition_func"].(string); ok {
				updateCondition.ConditionFunc = &conditionFunc
			}
			if conditionKey, ok := condMap["condition_key"].(string); ok {
				updateCondition.ConditionKey = &conditionKey
			}
			if conditionValueList, ok := condMap["condition_value_list"].([]interface{}); ok {
				values := make([]*string, len(conditionValueList))
				for j, value := range conditionValueList {
					strValue := value.(string)
					values[j] = &strValue
				}
				updateCondition.ConditionValueList = values
			}

			updateConditionList[i] = updateCondition
		}

		request.ConditionList = updateConditionList
	}

	var response *waapShareCustomizeBot.UpdateShareCustomizeBotTFResponse
	var err error
	err = resource.RetryContext(ctx, time.Duration(2)*time.Minute, func() *resource.RetryError {
		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseWaapShareCustomizeBotClient().Update(request)
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
	log.Printf("resource.wangsu_waap_share_customize_bot.update success")
	return nil
}

func resourceWaapShareCustomizeBotDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_waap_share_customize_bot.delete")

	var response *waapShareCustomizeBot.DeleteShareCustomizeBotsResponse
	var err error
	var diags diag.Diagnostics
	err = resource.RetryContext(ctx, time.Duration(2)*time.Minute, func() *resource.RetryError {
		id := data.Id()
		request := &waapShareCustomizeBot.DeleteShareCustomizeBotsRequest{
			IdList: []*string{&id},
		}
		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseWaapShareCustomizeBotClient().Delete(request)
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
