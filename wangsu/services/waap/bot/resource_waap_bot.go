package bot

import (
	"context"
	"errors"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wangsuCommon "github.com/wangsu-api/terraform-provider-wangsu/wangsu/common"
	waapBot "github.com/wangsu-api/wangsu-sdk-go/wangsu/waap/bot"
	"log"
	"time"
)

func ResourceWaapBot() *schema.Resource {
	return &schema.Resource{
		ReadContext:   ResourceWaapBotRead,
		UpdateContext: ResourceWaapBotUpdate,
		CreateContext: ResourceWaapBotCreate,
		DeleteContext: ResourceWaapBotDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Domain.",
			},

			"general_strategy": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "General policies.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ai_bots": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "AI Bots configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bot_category": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "AI Bots categories.",
									},
									"action": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Actions. <br/>NO_USE: Not Used.<br/>LOG: Log.<br/>BLOCK: Deny.",
									},
								},
							},
						},
						"public_bots": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Public Bots configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bot_category": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Public Bots categories.",
									},
									"action": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Actions. <br/>NO_USE: Not Used.<br/>LOG: Log.<br/>BLOCK: Deny.<br/>ACCEPT: Release.",
									},
								},
							},
						},
						"absolute_bots_act": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Definite Bots actions. NO_USE: Not Used.<br/>LOG: Log.<br/>BLOCK: Deny.",
						},
						"bot_tagging": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Bot Tagging.Header keys are case-insensitive and must be unique.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"request_header_key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Customizing HTTP request headers.Length must be ≤128 characters, using only ASCII characters with no colons or spaces permitted.",
									},
									"traffic_characteristics": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Tagging bot traffic characteristics.<br/>AI_BOT: AI Bots.<br/>PUBLIC_BOT: Public Bots.<br/>CUSTOMIZE_BOT: Custom Bots.",
									},
								},
							},
						},
					},
				},
			},
			"web_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "Web risk detection.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"act": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Action.<br/>NO_USE: Do not use<br/>BLOCK: Intercept<br/>LOG: Monitor",
						},
						"browser_analyse_switch": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Browser feature function switch.<br/>ON: Enable<br/>OFF: Disable",
						},
						"auto_tool_switch": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Automated tool detection function switch.<br/>ON: Enable<br/>OFF: Disable",
						},
						"crack_analyse_switch": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Cracking the behavior detection function switch.<br/>ON: Enable<br/>OFF: Disable",
						},
						"page_debug_switch": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Page anti-debugging function switch.<br/>ON: Enable<br/>OFF: Disable",
						},
						"interaction_analyse_switch": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Interactive behavior verification function switch.<br/>ON: Enable<br/>OFF: Disable",
						},
						"ajax_exception_switch": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Ajax request exception function switch.<br/>ON: Enable<br/>OFF: Disable",
						},
					},
				},
			},
			"traffic_detection": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "Abnormal traffic detection.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Start time, format: HH:mm.",
						},
						"end_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "End time, format: HH:mm.",
						},
						"action": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Action.<br/>RESET: Reset connection.<br/>LOG: Log.<br/>BLOCK: Block.<br/>NO_USE: Not Used.",
						},
						"whitelist": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Whitelist.",
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

func ResourceWaapBotRead(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.waap_bot.read")
	var diags diag.Diagnostics

	// 使用导入的 ID 设置资源 ID
	if data.Id() != "" {
		_ = data.Set("domain", data.Id())
	} else if domain, ok := data.GetOk("domain"); ok {
		data.SetId(domain.(string))
	}

	// Make API call
	request := &waapBot.GetBotManagementConfigRequest{}
	var response *waapBot.GetBotManagementConfigResponse
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		domain := data.Id()
		request.SetDomainList([]*string{&domain})
		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseWaapBotClient().GetBotManagementConfig(request)
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

	// Read response data
	readWaapBotResponseData(data, response)

	return diags
}

func ResourceWaapBotUpdate(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.waap_bot.update")
	var diags diag.Diagnostics

	if data.HasChange("domain") {
		// 把domain强制刷回旧值，否则会有权限问题
		oldDomain, _ := data.GetChange("domain")
		_ = data.Set("domain", oldDomain)
		err := errors.New("Domain cannot be changed.")
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	// Make API call
	var request = parseWaapBotRequestData(data)
	var response *waapBot.UpdateBotManagementConfigResponse
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseWaapBotClient().UpdateBotManagementConfig(request)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}
	if response == nil || tea.StringValue(response.Code) != "200" {
		return nil
	}

	return diags
}

func ResourceWaapBotCreate(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.waap_bot.create")
	var diags diag.Diagnostics
	err := errors.New("please execute \"terraform import <resource>.<resource_name> <domain>\" before executing \"terraform apply\" for the first time")
	diags = append(diags, diag.FromErr(err)...)
	return diags
}

func ResourceWaapBotDelete(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.waap_bot.delete")
	// 清空本地 state
	data.SetId("")
	return nil
}

// ==========

func readWaapBotResponseData(data *schema.ResourceData, response *waapBot.GetBotManagementConfigResponse) {
	if len(response.Data.ConfigList) != 1 {
		return
	}
	var botData = response.Data.ConfigList[0]

	// Parse general strategy
	if botData.GeneralStrategy != nil {
		generalStrategy := make([]interface{}, 1)
		strategy := map[string]interface{}{}

		aiBotsTfIds := make(map[string]struct{})
		publicBotsTfIds := make(map[string]struct{})
		if configList, ok := data.GetOk("general_strategy"); ok {
			if configListInterface, ok := configList.([]interface{}); ok {
				for _, configItem := range configListInterface {
					if configMap, ok := configItem.(map[string]interface{}); ok {
						// Extract ids from ai_bots
						if aiBots, ok := configMap["ai_bots"].([]interface{}); ok {
							for _, bot := range aiBots {
								if botMap, ok := bot.(map[string]interface{}); ok {
									if id, ok := botMap["bot_category"].(string); ok {
										aiBotsTfIds[id] = struct{}{}
									}
								}
							}
						}
						// Extract ids from public_bots
						if publicBots, ok := configMap["public_bots"].([]interface{}); ok {
							for _, bot := range publicBots {
								if botMap, ok := bot.(map[string]interface{}); ok {
									if id, ok := botMap["bot_category"].(string); ok {
										publicBotsTfIds[id] = struct{}{}
									}
								}
							}
						}
					}
				}
			}
		}

		// AI Bots
		if botData.GeneralStrategy.AiBots != nil {
			aiBots := make([]interface{}, 0)
			for _, bot := range botData.GeneralStrategy.AiBots {
				var id = tea.StringValue(bot.BotCategory)
				if _, exists := aiBotsTfIds[id]; len(aiBotsTfIds) == 0 || exists {
					aiBots = append(aiBots, map[string]interface{}{
						"bot_category": id,
						"action":       tea.StringValue(bot.Action),
					})
				}
			}
			strategy["ai_bots"] = aiBots
		}

		// Public Bots
		if botData.GeneralStrategy.PublicBots != nil {
			publicBots := make([]interface{}, 0)
			for _, bot := range botData.GeneralStrategy.PublicBots {
				var id = tea.StringValue(bot.BotCategory)
				if _, exists := publicBotsTfIds[id]; len(publicBotsTfIds) == 0 || exists {
					publicBots = append(publicBots, map[string]interface{}{
						"bot_category": id,
						"action":       tea.StringValue(bot.Action),
					})
				}
			}
			strategy["public_bots"] = publicBots
		}

		// Definite Bots Action
		strategy["absolute_bots_act"] = tea.StringValue(botData.GeneralStrategy.AbsoluteBotsAct)

		// Bot Tagging
		if botData.GeneralStrategy.BotTagging != nil {
			botTagging := make([]interface{}, len(botData.GeneralStrategy.BotTagging))
			for i, tag := range botData.GeneralStrategy.BotTagging {
				botTagging[i] = map[string]interface{}{
					"request_header_key":      tea.StringValue(tag.RequestHeaderKey),
					"traffic_characteristics": tea.StringValue(tag.TrafficCharacteristics),
				}
			}
			strategy["bot_tagging"] = botTagging
		}

		generalStrategy[0] = strategy
		_ = data.Set("general_strategy", generalStrategy)
	}

	// Parse web config
	if botData.WebConfig != nil {
		webConfig := make([]interface{}, 1)
		config := map[string]interface{}{
			"act":                        tea.StringValue(botData.WebConfig.Act),
			"browser_analyse_switch":     tea.StringValue(botData.WebConfig.BrowserAnalyseSwitch),
			"auto_tool_switch":           tea.StringValue(botData.WebConfig.AutoToolSwitch),
			"crack_analyse_switch":       tea.StringValue(botData.WebConfig.CrackAnalyseSwitch),
			"page_debug_switch":          tea.StringValue(botData.WebConfig.PageDebugSwitch),
			"interaction_analyse_switch": tea.StringValue(botData.WebConfig.InteractionAnalyseSwitch),
			"ajax_exception_switch":      tea.StringValue(botData.WebConfig.AjaxExceptionSwitch),
		}

		webConfig[0] = config
		_ = data.Set("web_config", webConfig)
	}

	// Parse traffic detection
	if botData.TrafficDetection != nil {
		trafficDetection := make([]interface{}, 1)
		detection := map[string]interface{}{
			"start_time": tea.StringValue(botData.TrafficDetection.StartTime),
			"end_time":   tea.StringValue(botData.TrafficDetection.EndTime),
			"action":     tea.StringValue(botData.TrafficDetection.Action),
			"whitelist":  tea.StringSliceValue(botData.TrafficDetection.Whitelist),
		}
		trafficDetection[0] = detection
		_ = data.Set("traffic_detection", trafficDetection)
	}
}

func parseWaapBotRequestData(data *schema.ResourceData) *waapBot.UpdateBotManagementConfigRequest {
	request := &waapBot.UpdateBotManagementConfigRequest{}
	request.SetDomain(data.Id())

	if data.HasChange("general_strategy") {
		if v, ok := data.GetOk("general_strategy"); ok && v != nil {
			generalStrategyList := v.([]interface{})
			if len(generalStrategyList) == 1 {
				generalStrategyMap := generalStrategyList[0].(map[string]interface{})
				var generalStrategy = &waapBot.UpdateBotManagementConfigRequestGeneralStrategy{}

				// AI Bots
				if data.HasChange("general_strategy.0.ai_bots") {
					if aiBots, ok := generalStrategyMap["ai_bots"].([]interface{}); ok {
						aiBotsConfig := make([]*waapBot.UpdateBotManagementConfigRequestGeneralStrategyAiBots, len(aiBots))
						for i, bot := range aiBots {
							botMap := bot.(map[string]interface{})
							aiBotsConfig[i] = &waapBot.UpdateBotManagementConfigRequestGeneralStrategyAiBots{
								BotCategory: tea.String(botMap["bot_category"].(string)),
								Action:      tea.String(botMap["action"].(string)),
							}
						}
						generalStrategy.AiBots = aiBotsConfig
					}
				}

				// Public Bots
				if data.HasChange("general_strategy.0.public_bots") {
					if publicBots, ok := generalStrategyMap["public_bots"].([]interface{}); ok {
						publicBotsConfig := make([]*waapBot.UpdateBotManagementConfigRequestGeneralStrategyPublicBots, len(publicBots))
						for i, bot := range publicBots {
							botMap := bot.(map[string]interface{})
							publicBotsConfig[i] = &waapBot.UpdateBotManagementConfigRequestGeneralStrategyPublicBots{
								BotCategory: tea.String(botMap["bot_category"].(string)),
								Action:      tea.String(botMap["action"].(string)),
							}
						}
						generalStrategy.PublicBots = publicBotsConfig
					}
				}

				// Definite Bots Action
				if data.HasChange("general_strategy.0.absolute_bots_act") {
					if absoluteBotsAct, ok := generalStrategyMap["absolute_bots_act"].(string); ok {
						generalStrategy.AbsoluteBotsAct = tea.String(absoluteBotsAct)
					}
				}

				// Bot Tagging
				if data.HasChange("general_strategy.0.bot_tagging") {
					if botTagging, ok := generalStrategyMap["bot_tagging"].([]interface{}); ok {
						botTaggingConfig := make([]*waapBot.UpdateBotManagementConfigRequestGeneralStrategyBotTagging, len(botTagging))
						for i, tag := range botTagging {
							tagMap := tag.(map[string]interface{})
							botTaggingConfig[i] = &waapBot.UpdateBotManagementConfigRequestGeneralStrategyBotTagging{
								RequestHeaderKey:       tea.String(tagMap["request_header_key"].(string)),
								TrafficCharacteristics: tea.String(tagMap["traffic_characteristics"].(string)),
							}
						}
						generalStrategy.BotTagging = botTaggingConfig
					}
				}

				request.GeneralStrategy = generalStrategy
			}
		}
	}

	if data.HasChange("web_config") {
		if v, ok := data.GetOk("web_config"); ok && v != nil {
			webConfigList := v.([]interface{})
			if len(webConfigList) == 1 {
				webConfigMap := webConfigList[0].(map[string]interface{})
				var webConfig = &waapBot.UpdateBotManagementConfigRequestWebConfig{
					Act:                      tea.String(webConfigMap["act"].(string)),
					BrowserAnalyseSwitch:     tea.String(webConfigMap["browser_analyse_switch"].(string)),
					AutoToolSwitch:           tea.String(webConfigMap["auto_tool_switch"].(string)),
					CrackAnalyseSwitch:       tea.String(webConfigMap["crack_analyse_switch"].(string)),
					PageDebugSwitch:          tea.String(webConfigMap["page_debug_switch"].(string)),
					InteractionAnalyseSwitch: tea.String(webConfigMap["interaction_analyse_switch"].(string)),
					AjaxExceptionSwitch:      tea.String(webConfigMap["ajax_exception_switch"].(string)),
				}

				request.WebConfig = webConfig
			}
		}
	}

	if data.HasChange("traffic_detection") {
		if v, ok := data.GetOk("traffic_detection"); ok && v != nil {
			trafficDetectionList := v.([]interface{})
			if len(trafficDetectionList) == 1 {
				trafficDetectionMap := trafficDetectionList[0].(map[string]interface{})
				var trafficDetection = &waapBot.UpdateBotManagementConfigRequestTrafficDetection{
					StartTime: tea.String(trafficDetectionMap["start_time"].(string)),
					EndTime:   tea.String(trafficDetectionMap["end_time"].(string)),
					Action:    tea.String(trafficDetectionMap["action"].(string)),
				}
				if v, ok := trafficDetectionMap["whitelist"]; ok && v != nil {
					whitelist := v.([]interface{})
					whitelistStr := make([]*string, len(whitelist))
					for i, item := range whitelist {
						whitelistStr[i] = tea.String(item.(string))
					}
					trafficDetection.Whitelist = whitelistStr
				}

				request.TrafficDetection = trafficDetection
			}
		}
	}

	return request
}
