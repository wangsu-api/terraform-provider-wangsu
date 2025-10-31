package bot_scene_whitelist

import (
	"context"
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wangsuCommon "github.com/wangsu-api/terraform-provider-wangsu/wangsu/common"
	waapBotSceneWhitelist "github.com/wangsu-api/wangsu-sdk-go/wangsu/waap/bot-scene-whitelist"
	"log"
	"time"
)

func ResourceWaapBotSceneWhitelist() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWaapBotSceneWhitelistCreate,
		ReadContext:   resourceWaapBotSceneWhitelistRead,
		UpdateContext: resourceWaapBotSceneWhitelistUpdate,
		DeleteContext: resourceWaapBotSceneWhitelistDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
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
							Description: "When matchName is IP_IPS, maximum 300 IP/CIDR in match value list, the optional value of matchType is:<br/>EQUAL: Equals<br/>NOT_EQUAL: Does not equal<br/>When matchName is a URI, the optional value of matchType is:<br/>EQUAL: Equals, the matching value is case-sensitive and needs to start with \"/\" and include parameters<br/>NOT_EQUAL: Does not equal, the matching value is case-sensitive, needs to start with \"/\", and contains parameters<br/>CONTAIN: Contains, match values are case insensitive<br/>NOT_CONTAIN: Does not contains, match values are case insensitive<br/>REGEX: Regex match, only one value is allowed for the match value<br/>NOT_REGEX: regular does not match<br/>START_WITH: starts with<br/>END_WITH: ends with<br/>WILDCARD: wildcard matches<br/>NOT_WILDCARD: wildcard does not match<br/>When matchName is PATH, the optional value of matchType is:<br/>EQUAL: Equals, the matching value is case-sensitive and needs to start with \"/\" , and does not contain parameters<br/>NOT_EQUAL: Does not equal, the matching value is case-sensitive, needs to start with \"/\", and does not contain parameters<br/>CONTAIN: Contains, match values are case insensitive<br/>NOT_CONTAIN: Does not contains, match values are case insensitive<br/>REGEX: Regex match, match values are case insensitive and only one value is allowed<br/>NOT_REGEX: regular does not match<br/>START_WITH: starts with<br/>END_WITH: ends with<br/>WILDCARD: wildcard matches<br/>NOT_WILDCARD: wildcard does not match<br/>When matchName is HEADER, the optional value of matchType is:<br/>EQUAL: Equals, match values are case sensitive<br/>NOT_EQUAL: Does not equal, the matching value is case-sensitive<br/>CONTAIN: Contains, match values are case insensitive<br/>NOT_CONTAIN: Does not contains, match values are case insensitive<br/>REGEX: Regex match, match values are case insensitive and only one value is allowed<br/>NONE: Empty or does not exist<br/>NOT_REGEX: regular does not match<br/>START_WITH: starts with<br/>END_WITH: ends with<br/>WILDCARD: wildcard matches<br/>NOT_WILDCARD: wildcard does not match<br/>When matchName is UA, the optional value of matchType is:<br/>EQUAL: Equals, match values are case sensitive<br/>NOT_EQUAL: Does not equal, the matching value is case-sensitive<br/>CONTAIN: Contains, match values are case insensitive<br/>NOT_CONTAIN: Does not contains, match values are case insensitive<br/>REGEX: Regex match, match values are case insensitive and only one value is allowed<br/>NONE: Empty or does not exist<br/>NOT_REGEX: regular does not match<br/>START_WITH: starts with<br/>END_WITH: ends with<br/>WILDCARD: wildcard matches<br/>NOT_WILDCARD: wildcard does not match<br/>When matchName is REFERER, the optional value of matchType is:<br/>EQUAL: Equals, match values are case sensitive<br/>NOT_EQUAL: Does not equal, the matching value is case-sensitive<br/>CONTAIN: Contains, match values are case insensitive<br/>NOT_CONTAIN: Does not contains, match values are case insensitive<br/>REGEX: Regex match, match values are case insensitive and only one value is allowed<br/>NONE: Empty or does not exist<br/>NOT_REGEX: regular does not match<br/>START_WITH: starts with<br/>END_WITH: ends with<br/>WILDCARD: wildcard matches<br/>NOT_WILDCARD: wildcard does not match<br/>When matchName is REQUEST_METHOD, the optional value of matchType is:<br/>EQUAL: Equals, match values are case sensitive<br/>NOT_EQUAL: Does not equal, the matching value is case-sensitive<br/>",
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
	}
}

func resourceWaapBotSceneWhitelistCreate(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.waap_bot_scene_whitelist.create")

	var diags diag.Diagnostics
	request := &waapBotSceneWhitelist.AddSpecificClientTrafficBypassRequest{}

	if domain, ok := data.Get("domain").(string); ok && domain != "" {
		request.Domain = &domain
	}
	if name, ok := data.Get("name").(string); ok && name != "" {
		request.Name = &name
	}
	if description, ok := data.Get("description").(string); ok && description != "" {
		request.Description = &description
	}
	if conditions, ok := data.GetOk("conditions"); ok {
		conditionList := conditions.([]interface{})
		addConditionList := make([]*waapBotSceneWhitelist.AddSpecificClientTrafficBypassRequestConditions, len(conditionList))

		for i, condition := range conditionList {
			condMap := condition.(map[string]interface{})
			addCondition := &waapBotSceneWhitelist.AddSpecificClientTrafficBypassRequestConditions{}

			if matchName, ok := condMap["match_name"].(string); ok {
				addCondition.MatchName = &matchName
			}
			if matchType, ok := condMap["match_type"].(string); ok {
				addCondition.MatchType = &matchType
			}
			if matchKey, ok := condMap["match_key"].(string); ok {
				addCondition.MatchKey = &matchKey
			}
			if matchValueList, ok := condMap["match_value_list"].([]interface{}); ok {
				values := make([]*string, len(matchValueList))
				for j, value := range matchValueList {
					strValue := value.(string)
					values[j] = &strValue
				}
				addCondition.MatchValueList = values
			}
			addConditionList[i] = addCondition
		}
		request.Conditions = addConditionList
	}

	var response *waapBotSceneWhitelist.AddSpecificClientTrafficBypassResponse
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseWaapBotSceneWhiteListClient().Add(request)
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

	data.SetId(*response.Data)
	return resourceWaapBotSceneWhitelistRead(context, data, meta)
}

func resourceWaapBotSceneWhitelistRead(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.waap_bot_scene_whitelist.read")
	var response *waapBotSceneWhitelist.ListSpecificClientTrafficBypassResponse
	var err error
	var diags diag.Diagnostics
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		domain := data.Get("domain").(string)
		request := &waapBotSceneWhitelist.ListSpecificClientTrafficBypassRequest{
			DomainList: []*string{&domain},
		}
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
		for _, item := range response.Data {
			// 只要对应id的数据
			if *item.Id != data.Id() {
				continue
			}
			_ = data.Set("domain", item.Domain)
			_ = data.Set("name", item.Name)
			_ = data.Set("description", item.Description)
			// 映射 conditions 数据
			conditions := make([]map[string]interface{}, len(item.Conditions))
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
			_ = data.Set("conditions", conditions)
		}
	}
	return nil
}

func resourceWaapBotSceneWhitelistUpdate(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.waap_bot_scene_whitelist.update")

	var diags diag.Diagnostics
	if data.HasChange("domain") {
		// 把domain强制刷回旧值，否则会有权限问题
		oldDomain, _ := data.GetChange("domain")
		_ = data.Set("domain", oldDomain)
		err := errors.New("Hostname cannot be changed.")
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}
	request := &waapBotSceneWhitelist.UpdateSpecificClientTrafficBypassRequest{}
	if id, ok := data.Get("id").(string); ok && id != "" {
		request.Id = &id
	}
	if domain, ok := data.Get("domain").(string); ok && domain != "" {
		request.Domain = &domain
	}
	if name, ok := data.Get("name").(string); ok && name != "" {
		request.Name = &name
	}
	if description, ok := data.Get("description").(string); ok {
		request.Description = &description
	}
	if conditions, ok := data.GetOk("conditions"); ok {
		conditionList := conditions.([]interface{})
		updateConditionList := make([]*waapBotSceneWhitelist.UpdateSpecificClientTrafficBypassRequestConditions, len(conditionList))

		for i, condition := range conditionList {
			condMap := condition.(map[string]interface{})
			updateCondition := &waapBotSceneWhitelist.UpdateSpecificClientTrafficBypassRequestConditions{}

			if matchName, ok := condMap["match_name"].(string); ok {
				updateCondition.MatchName = &matchName
			}
			if matchType, ok := condMap["match_type"].(string); ok {
				updateCondition.MatchType = &matchType
			}
			if matchKey, ok := condMap["match_key"].(string); ok {
				updateCondition.MatchKey = &matchKey
			}
			if matchValueList, ok := condMap["match_value_list"].([]interface{}); ok {
				values := make([]*string, len(matchValueList))
				for j, value := range matchValueList {
					strValue := value.(string)
					values[j] = &strValue
				}
				updateCondition.MatchValueList = values
			}

			updateConditionList[i] = updateCondition
		}

		request.Conditions = updateConditionList
	}
	var response *waapBotSceneWhitelist.UpdateSpecificClientTrafficBypassResponse
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseWaapBotSceneWhiteListClient().Update(request)
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

func resourceWaapBotSceneWhitelistDelete(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.waap_bot_scene_whitelist.delete")

	var response *waapBotSceneWhitelist.DeleteSpecificClientTrafficBypassResponse
	var err error
	var diags diag.Diagnostics
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		id := data.Id()
		request := &waapBotSceneWhitelist.DeleteSpecificClientTrafficBypassRequest{
			IdList: []*string{&id},
		}
		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseWaapBotSceneWhiteListClient().Delete(request)
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
