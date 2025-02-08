package rule

import (
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wangsuCommon "github.com/wangsu-api/terraform-provider-wangsu/wangsu/common"
	"github.com/wangsu-api/wangsu-sdk-go/wangsu/monitor/rule"
	"golang.org/x/net/context"
)

func ResourceMonitorRealtimeRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRealTimeRuleCreate,
		ReadContext:   resourceRealTimeRuleRead,
		UpdateContext: resourceRealTimeRuleUpdate,
		DeleteContext: resourceRealTimeRuleDelete,
		Schema: map[string]*schema.Schema{
			"rule_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Rule ID",
			},
			"rule_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Alert rule name, only supports Chinese, English, numbers, underscore, hyphen, max 100 characters",
			},
			"monitor_product": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Monitor by dimension or product. Dimensions: userDimension, DG, domainDimension. Or input product code",
			},
			"resource_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Resource type. Required when monitorProduct is specific product. Options: domain",
			},
			"monitor_resources": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of resource names to monitor or ALL",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"statistical_method": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Statistical method. Default: CONSOLIDATED. Options: CONSOLIDATED-consolidated statistics, SEPARATE-separate statistics",
				Default:     "CONSOLIDATED",
			},
			"alert_frequency": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Alert frequency in minutes. Default: 5. Options: 0-first alert only, 2, 5, 10, 15, 20",
				Default:     5,
			},
			"restore_notice": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Whether to notify on alert recovery. Default: true",
				Default:     "true",
			},
			"rule_items": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Rule items list",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_item_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Time-divided sub-rules. If an id is passed in, it is considered to specify a modification of an item rule.",
						},
						"start_time": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Monitoring period start time, format: HH:00, example: 00:00",
						},
						"end_time": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Monitoring period end time, format: HH:59, example: 01:59",
						},
						"condition_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Alert condition type. Options: ANY-any condition, ALL-all conditions",
						},
						"period": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Monitoring cycle in minutes. Options: 1, 5, 10",
						},
						"period_times": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Number of cycles to meet conditions before alerting. Options: 1, 2, 3, 5, 15, 30",
						},
						"condition_rules": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Condition rules list",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"monitor_item": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Monitoring items currently only support domain-type monitoring items. Options: BANDWIDTH, FLOW, REQUEST, BTOB, FLOW_HIT_RATE, etc.",
									},
									"condition": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Condition type. Options: MAX-greater than, MIN-less than, UPRUSH-surge, SLUMPED-plunge",
									},
									"threshold": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Threshold value, please provide a positive integer",
									},
								},
							},
						},
					},
				},
			},
			"notices": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Notification methods list",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"notice_method": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Notification method. Options: MOBILE, EMAIL, ROBOT, WEBHOOK",
						},
						"notice_object": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Notification object. For MOBILE/EMAIL: contact IDs separated by ;. For ROBOT: robot IDs separated by ;. For WEBHOOK: webhook URL",
						},
					},
				},
			},
		},
	}
}

func resourceRealTimeRuleCreate(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_monitor_realtime_rule.create")
	var diags diag.Diagnostics
	request := &rule.CreateCloudMonitorRealTimeAlarmRuleRequest{}

	// 设置必需字段
	if ruleName, ok := data.Get("rule_name").(string); ok && ruleName != "" {
		request.RuleName = &ruleName
	}

	if monitorProduct, ok := data.Get("monitor_product").(string); ok && monitorProduct != "" {
		request.MonitorProduct = &monitorProduct
	}

	if resourceType, ok := data.Get("resource_type").(string); ok && resourceType != "" {
		request.ResourceType = &resourceType
	}

	// 处理监控资源列表
	if v, ok := data.Get("monitor_resources").([]interface{}); ok && len(v) > 0 {
		resources := make([]*string, 0, len(v))
		for _, res := range v {
			r := res.(string)
			resources = append(resources, &r)
		}
		request.MonitorResources = resources
	}

	if statisticalMethod, ok := data.Get("statistical_method").(string); ok && statisticalMethod != "" {
		request.StatisticalMethod = &statisticalMethod
	}

	if alertFrequency, ok := data.Get("alert_frequency").(int); ok {
		alertFrequencyInt32 := int32(alertFrequency)
		request.AlertFrequency = &alertFrequencyInt32
	}

	if restoreNotice, ok := data.Get("restore_notice").(string); ok && restoreNotice != "" {
		request.RestoreNotice = &restoreNotice
	}

	// 处理规则项列表
	if v, ok := data.Get("rule_items").([]interface{}); ok && len(v) > 0 {
		ruleItems := make([]*rule.CreateRuleItem, 0, len(v))
		for _, item := range v {
			itemMap := item.(map[string]interface{})
			ruleItem := &rule.CreateRuleItem{}

			if startTime, ok := itemMap["start_time"].(string); ok && startTime != "" {
				ruleItem.StartTime = &startTime
			}
			if endTime, ok := itemMap["end_time"].(string); ok && endTime != "" {
				ruleItem.EndTime = &endTime
			}
			if conditionType, ok := itemMap["condition_type"].(string); ok && conditionType != "" {
				ruleItem.ConditionType = &conditionType
			}
			if period, ok := itemMap["period"].(int); ok {
				periodInt32 := int32(period)
				ruleItem.Period = &periodInt32
			}
			if periodTimes, ok := itemMap["period_times"].(int); ok {
				periodTimesInt32 := int32(periodTimes)
				ruleItem.PeriodTimes = &periodTimesInt32
			}

			// 处理条件规则
			if conditions, ok := itemMap["condition_rules"].([]interface{}); ok && len(conditions) > 0 {
				conditionRules := make([]*rule.CreateRuleCondition, 0, len(conditions))
				for _, condition := range conditions {
					condMap := condition.(map[string]interface{})
					condRule := &rule.CreateRuleCondition{}

					if monitorItem, ok := condMap["monitor_item"].(string); ok && monitorItem != "" {
						condRule.MonitorItem = &monitorItem
					}
					if cond, ok := condMap["condition"].(string); ok && cond != "" {
						condRule.Condition = &cond
					}
					if threshold, ok := condMap["threshold"].(string); ok && threshold != "" {
						condRule.Threshold = &threshold
					}
					conditionRules = append(conditionRules, condRule)
				}
				ruleItem.ConditionRules = conditionRules
			}

			ruleItems = append(ruleItems, ruleItem)
		}
		request.RuleItems = ruleItems
	}

	// 处理通知方式
	if v, ok := data.Get("notices").([]interface{}); ok && len(v) > 0 {
		notices := make([]*rule.CreateRuleNotice, 0, len(v))
		for _, item := range v {
			noticeMap := item.(map[string]interface{})
			notice := &rule.CreateRuleNotice{}

			if method, ok := noticeMap["notice_method"].(string); ok && method != "" {
				notice.NoticeMethod = &method
			}
			if object, ok := noticeMap["notice_object"].(string); ok && object != "" {
				notice.NoticeObject = &object
			}
			notices = append(notices, notice)
		}
		request.Notices = notices
	}

	// 调用 API 创建规则
	var response *rule.CreateCloudMonitorRealTimeAlarmRuleResponse
	var requestId string
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		requestId, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseMonitorRuleClient().CreateRealTimeRule(request)
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
		data.SetId("")
		return nil
	}
	_ = data.Set("code", response.Code)
	_ = data.Set("message", response.Message)
	_ = data.Set("rule_id", *response.Data.RuleId)
	data.SetId(*response.Data.RuleId)
	log.Printf("resource.wangsu_monitor_realtime_rule.create finish, requestId: %s", requestId)
	return diags

}

func resourceRealTimeRuleUpdate(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_monitor_realtime_rule.update")
	var diags diag.Diagnostics
	request := &rule.EditCloudMonitorRealTimeAlarmRuleRequest{}

	if ruleId, ok := data.Get("rule_id").(string); ok && ruleId != "" {
		request.RuleId = &ruleId
	}

	if ruleName, ok := data.Get("rule_name").(string); ok && ruleName != "" {
		request.RuleName = &ruleName
	}

	if monitorProduct, ok := data.Get("monitor_product").(string); ok && monitorProduct != "" {
		request.MonitorProduct = &monitorProduct
	}

	if resourceType, ok := data.Get("resource_type").(string); ok && resourceType != "" {
		request.ResourceType = &resourceType
	}

	// 处理监控资源列表
	if v, ok := data.Get("monitor_resources").([]interface{}); ok && len(v) > 0 {
		resources := make([]*string, 0, len(v))
		for _, res := range v {
			r := res.(string)
			resources = append(resources, &r)
		}
		request.MonitorResources = resources
	}

	if statisticalMethod, ok := data.Get("statistical_method").(string); ok && statisticalMethod != "" {
		request.StatisticalMethod = &statisticalMethod
	}

	if alertFrequency, ok := data.Get("alert_frequency").(int); ok {
		alertFrequencyInt32 := int32(alertFrequency)
		request.AlertFrequency = &alertFrequencyInt32
	}

	if restoreNotice, ok := data.Get("restore_notice").(string); ok && restoreNotice != "" {
		request.RestoreNotice = &restoreNotice
	}

	// 处理规则项列表
	if v, ok := data.Get("rule_items").([]interface{}); ok && len(v) > 0 {
		ruleItems := make([]*rule.EditRuleItem, 0, len(v))
		for _, item := range v {
			itemMap := item.(map[string]interface{})
			ruleItem := &rule.EditRuleItem{}

			if id, ok := itemMap["rule_item_id"].(string); ok && id != "" {
				ruleItem.RuleItemId = &id
			}
			if startTime, ok := itemMap["start_time"].(string); ok && startTime != "" {
				ruleItem.StartTime = &startTime
			}
			if endTime, ok := itemMap["end_time"].(string); ok && endTime != "" {
				ruleItem.EndTime = &endTime
			}
			if conditionType, ok := itemMap["condition_type"].(string); ok && conditionType != "" {
				ruleItem.ConditionType = &conditionType
			}
			if period, ok := itemMap["period"].(int); ok {
				periodInt32 := int32(period)
				ruleItem.Period = &periodInt32
			}
			if periodTimes, ok := itemMap["period_times"].(int); ok {
				periodTimesInt32 := int32(periodTimes)
				ruleItem.PeriodTimes = &periodTimesInt32
			}

			// 处理条件规则
			if conditions, ok := itemMap["condition_rules"].([]interface{}); ok && len(conditions) > 0 {
				conditionRules := make([]*rule.EditRuleCondition, 0, len(conditions))
				for _, condition := range conditions {
					condMap := condition.(map[string]interface{})
					condRule := &rule.EditRuleCondition{}

					if monitorItem, ok := condMap["monitor_item"].(string); ok && monitorItem != "" {
						condRule.MonitorItem = &monitorItem
					}
					if cond, ok := condMap["condition"].(string); ok && cond != "" {
						condRule.Condition = &cond
					}
					if threshold, ok := condMap["threshold"].(string); ok && threshold != "" {
						condRule.Threshold = &threshold
					}
					conditionRules = append(conditionRules, condRule)
				}
				ruleItem.ConditionRules = conditionRules
			}

			ruleItems = append(ruleItems, ruleItem)
		}
		request.RuleItems = ruleItems
	}

	// 处理通知方式
	if v, ok := data.Get("notices").([]interface{}); ok && len(v) > 0 {
		notices := make([]*rule.EditRuleNotice, 0, len(v))
		for _, item := range v {
			noticeMap := item.(map[string]interface{})
			notice := &rule.EditRuleNotice{}

			if method, ok := noticeMap["notice_method"].(string); ok && method != "" {
				notice.NoticeMethod = &method
			}
			if object, ok := noticeMap["notice_object"].(string); ok && object != "" {
				notice.NoticeObject = &object
			}
			notices = append(notices, notice)
		}
		request.Notices = notices
	}

	// 调用 API 创建规则
	var response *rule.EditCloudMonitorRealTimeAlarmRuleResponse
	var requestId string
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		requestId, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseMonitorRuleClient().EditRealTimeRule(request)
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
	_ = data.Set("code", response.Code)
	_ = data.Set("x_cnc_request_id", requestId)
	_ = data.Set("message", response.Message)
	data.SetId(*response.Data.RuleId)
	log.Printf("resource.wangsu_monitor_realtime_rule.update finish, requestId: %s", requestId)
	return diags

}

func resourceRealTimeRuleRead(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	request := &rule.QueryCloudMonitorRealTimeAlarmRuleRequest{}
	ruleId := data.Id()
	request.RuleId = &ruleId
	var response *rule.QueryCloudMonitorRealTimeAlarmRuleResponse
	var requestId string
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		requestId, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseMonitorRuleClient().QueryRealTimeRule(request)
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

	responseData := response.Data

	_ = data.Set("rule_id", ruleId)
	_ = data.Set("rule_name", responseData.RuleName)
	_ = data.Set("monitor_product", responseData.MonitorProduct)
	_ = data.Set("resource_type", responseData.ResourceType)
	_ = data.Set("monitor_resources", responseData.MonitorResources)
	_ = data.Set("statistical_method", responseData.StatisticalMethod)
	_ = data.Set("alert_frequency", responseData.AlertFrequency)
	_ = data.Set("restore_notice", responseData.RestoreNotice)

	// 设置规则项列表
	if responseData.RuleItems != nil {
		ruleItems := make([]interface{}, len(responseData.RuleItems))
		for i, item := range responseData.RuleItems {
			ruleItem := make(map[string]interface{})

			// 只有当 RuleItemId 不为空时才设置
			if item.RuleItemId != nil && *item.RuleItemId != "" {
				ruleItem["rule_item_id"] = *item.RuleItemId
			}

			ruleItem["start_time"] = item.StartTime
			ruleItem["end_time"] = item.EndTime
			ruleItem["condition_type"] = item.ConditionType
			ruleItem["period"] = item.Period
			ruleItem["period_times"] = item.PeriodTimes

			// 设置条件规则
			if item.ConditionRules != nil {
				conditionRules := make([]interface{}, len(item.ConditionRules))
				for j, cond := range item.ConditionRules {
					condRule := make(map[string]interface{})
					condRule["monitor_item"] = cond.MonitorItem
					condRule["condition"] = cond.Condition
					condRule["threshold"] = cond.Threshold
					conditionRules[j] = condRule
				}
				ruleItem["condition_rules"] = conditionRules
			}
			ruleItems[i] = ruleItem
		}
		_ = data.Set("rule_items", ruleItems)
	}

	// 设置通知方式
	if responseData.Notices != nil {
		notices := make([]interface{}, len(responseData.Notices))
		for i, notice := range responseData.Notices {
			noticeMap := make(map[string]interface{})
			noticeMap["notice_method"] = notice.NoticeMethod
			noticeMap["notice_object"] = notice.NoticeObject
			notices[i] = noticeMap
		}
		_ = data.Set("notices", notices)
	}
	log.Printf("resource.wangsu_monitor_realtime_rule.read finish, requestId: %s", requestId)
	return diags

}

func resourceRealTimeRuleDelete(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_monitor_realtime_rule.delete")
	var diags diag.Diagnostics
	request := &rule.DeleteCloudMonitorRealTimeAlarmRuleRequest{}
	ruleId := data.Id()
	request.RuleId = &ruleId
	var response *rule.DeleteCloudMonitorRealTimeAlarmRuleResponse
	var requestId string
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		requestId, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseMonitorRuleClient().DeleteRealTimeRule(request)
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
		data.SetId("")
		return nil
	}
	_ = data.Set("code", response.Code)
	_ = data.Set("message", response.Message)
	data.SetId(data.Id())
	log.Printf("resource.wangsu_monitor_realtime_rule.delete finish, requestId: %s", requestId)
	return diags

}
