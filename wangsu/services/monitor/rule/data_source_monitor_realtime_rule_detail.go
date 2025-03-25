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

func DataSourceMonitorRealtimeRuleDetail() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRealTimeRuleDetailRead,
		Schema: map[string]*schema.Schema{
			"rule_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Alert rule name",
			},
			"code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Response code, 0 means successful.",
			},
			"message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Response error message if failed.",
			},
			"data": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Response data.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule ID",
						},
						"rule_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Alert rule name",
						},
						"monitor_product": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Monitor by dimension or product.",
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource type.",
						},
						"monitor_resources": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of resource names to monitor or ALL",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"statistical_method": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Statistical method. Options: CONSOLIDATED-consolidated statistics, SEPARATE-separate statistics",
						},
						"alert_frequency": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Alert frequency in minutes. Options: 0-first alert only, 2, 5, 10, 15, 20",
						},
						"restore_notice": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Whether to notify on alert recovery. Default: true",
						},
						"rule_items": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Rule items list",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"rule_item_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Time-divided sub-rules. If an id is passed in, it is considered to specify a modification of an item rule.",
									},
									"start_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Monitoring period start time, format: HH:00, example: 00:00",
									},
									"end_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Monitoring period end time, format: HH:59, example: 01:59",
									},
									"condition_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Alert condition type. Options: ANY-any condition, ALL-all conditions",
									},
									"period": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Monitoring cycle in minutes. Options: 1, 5, 10",
									},
									"period_times": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Number of cycles to meet conditions before alerting. Options: 1, 2, 3, 5, 15, 30",
									},
									"condition_rules": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Condition rules list",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"monitor_item": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Monitoring items currently only support domain-type monitoring items. Options: BANDWIDTH, FLOW, REQUEST, BTOB, FLOW_HIT_RATE, etc.",
												},
												"condition": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Condition type. Options: MAX-greater than, MIN-less than, UPRUSH-surge, SLUMPED-plunge",
												},
												"threshold": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Threshold value.",
												},
											},
										},
									},
								},
							},
						},
						"notices": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Notification methods list",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"notice_method": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Notification method. Options: MOBILE, EMAIL, ROBOT, webhook",
									},
									"notice_object": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Notification object. For MOBILE/EMAIL: contact IDs separated by ;. For ROBOT: robot IDs separated by ;. For WEBHOOK: webhook URL",
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

func dataSourceRealTimeRuleDetailRead(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("data_source.wangsu_monitor_realtime_rule_detail.read")
	var diags diag.Diagnostics
	ruleName := data.Get("rule_name").(string)

	request := &rule.QueryCloudMonitorRealTimeAlarmRuleRequest{}
	request.RuleName = &ruleName
	var response *rule.QueryCloudMonitorRealTimeAlarmRuleResponse
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseMonitorRuleClient().QueryRealTimeRule(request)
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
	_ = data.Set("message", response.Message)
	var resultList []interface{}
	var ruleDetail = map[string]interface{}{
		"rule_id":            response.Data.RuleId,
		"rule_name":          response.Data.RuleName,
		"monitor_product":    response.Data.MonitorProduct,
		"resource_type":      response.Data.ResourceType,
		"monitor_resources":  response.Data.MonitorResources,
		"statistical_method": response.Data.StatisticalMethod,
		"alert_frequency":    response.Data.AlertFrequency,
		"restore_notice":     response.Data.RestoreNotice,
		"rule_items":         buildRuleItems(response.Data.RuleItems),
		"notices":            buildNotices(response.Data.Notices),
	}
	resultList = append(resultList, ruleDetail)

	_ = data.Set("data", resultList)
	data.SetId(*response.Data.RuleName)

	log.Printf("data_source.wangsu_monitor_realtime_rule_detail.read finish")
	return nil

}

func buildNotices(dataNotices []*rule.QueryNotice) interface{} {
	if dataNotices == nil {
		return nil
	}
	var notices []interface{}
	for _, dataNotice := range dataNotices {
		var notice = map[string]interface{}{
			"notice_method": dataNotice.NoticeMethod,
			"notice_object": dataNotice.NoticeObject,
		}
		notices = append(notices, notice)
	}
	return notices
}

func buildRuleItems(items []*rule.QueryRuleItem) interface{} {
	if items == nil {
		return nil
	}
	var ruleItems []interface{}
	for _, item := range items {
		var ruleItem = map[string]interface{}{
			"rule_item_id":    item.RuleItemId,
			"start_time":      item.StartTime,
			"end_time":        item.EndTime,
			"condition_type":  item.ConditionType,
			"period":          item.Period,
			"period_times":    item.PeriodTimes,
			"condition_rules": buildConditionRules(item.ConditionRules),
		}
		ruleItems = append(ruleItems, ruleItem)
	}
	return ruleItems
}

func buildConditionRules(rules []*rule.QueryConditionRule) interface{} {
	if rules == nil {
		return nil
	}
	var conditionRules []interface{}
	for _, cond := range rules {
		var conditionRule = map[string]interface{}{
			"monitor_item": cond.MonitorItem,
			"condition":    cond.Condition,
			"threshold":    cond.Threshold,
		}
		conditionRules = append(conditionRules, conditionRule)
	}
	return conditionRules
}
