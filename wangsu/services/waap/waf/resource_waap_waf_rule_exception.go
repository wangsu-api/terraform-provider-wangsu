package waf

import (
	"context"
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wangsuCommon "github.com/wangsu-api/terraform-provider-wangsu/wangsu/common"
	securityPolicy "github.com/wangsu-api/wangsu-sdk-go/wangsu/securitypolicy"
	"log"
	"time"
)

func ResourceWaapWafRuleException() *schema.Resource {
	return &schema.Resource{
		CreateContext: createWafRuleException,
		ReadContext:   readWafRuleException,
		UpdateContext: updateWafRuleException,
		DeleteContext: deleteWafRuleException,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Exception ID.",
			},
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Domain.",
			},
			"rule_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "WAF rule ID.",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Matching conditions.<br/>ip: IP<br/>path: Path<br/>uri: URI<br/>urlParamName: URI Parameter Name<br/>urlParamValue: URI Parameter Value<br/>userAgent: User Agent<br/>httpHeaderName: Request Header Name<br/>httpHeaderValue: Request Header Value<br/>cookie: Cookie<br/>body: Body<br/>bodyParamName: Body Parameter Name<br/>bodyParamValue: Body Parameter Value",
			},
			"match_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Match type,IP can only be EQUAL.<br/>EQUAL: Equal<br/>CONTAIN: Contains<br/>REGEX: Regular match",
			},
			"content_list": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Rule exceptions.<br/>When matchType=EQUAL, case-sensitive, path and uri must start with \"/\", and body can only pass one value;<br/>When matchType=REGEX, only one value can be passed.",
				Elem: &schema.Schema{
					Type:        schema.TypeString,
					Description: "Rule exceptions. When matchType=EQUAL, case-sensitive, path and uri must start with \"/\", and body can only pass one value; When matchType=REGEX, only one value can be passed.",
				},
			},
		},
	}
}

func createWafRuleException(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.waap_waf_rule_exception.create")

	var diags diag.Diagnostics
	request := &securityPolicy.CreateExceptionToWAFManagedRulesRequest{}
	if domain, ok := data.Get("domain").(string); ok && domain != "" {
		request.Domain = &domain
	}

	if ruleId, ok := data.Get("rule_id").(int); ok {
		request.RuleId = &ruleId
	}

	if exceptionType, ok := data.Get("type").(string); ok && exceptionType != "" {
		request.Type = &exceptionType
	}

	if matchType, ok := data.Get("match_type").(string); ok && matchType != "" {
		request.MatchType = &matchType
	}

	if contentList, ok := data.GetOk("content_list"); ok {
		contents := contentList.([]interface{})
		contentsStr := make([]*string, len(contents))
		for i, v := range contents {
			str := v.(string)
			contentsStr[i] = &str
		}
		request.ContentList = contentsStr
	}

	var response *securityPolicy.CreateExceptionToWAFManagedRulesResponse
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseSecurityPolicyClient().CreateExceptionToWAFManagedRules(request)
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
	data.SetId(*response.Data)
	//set status
	return readWafRuleException(context, data, meta)
}

func readWafRuleException(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.waap_waf_rule_exception.read")

	var response *securityPolicy.ListNonSharedWAFRuleExceptionsForWAFRulesResponse
	var err error
	var diags diag.Diagnostics
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		domain := data.Get("domain").(string)
		ruleId := data.Get("rule_id").(int)

		request := &securityPolicy.ListNonSharedWAFRuleExceptionsForWAFRulesRequest{
			DomainList: []*string{&domain},
			RuleIdList: []*int{&ruleId},
		}
		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseSecurityPolicyClient().ListNonSharedWAFRuleExceptionsForWAFRules(request)
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
			_ = data.Set("rule_id", item.RuleId)
			_ = data.Set("type", item.Type)
			_ = data.Set("match_type", item.MatchType)
			_ = data.Set("content_list", item.ContentList)
		}
	}

	return nil
}

func updateWafRuleException(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.waap_waf_rule_exception.update")
	var diags diag.Diagnostics
	var canNotChange = false
	if data.HasChange("domain") {
		// 把domain强制刷回旧值，否则会有权限问题
		oldDomain, _ := data.GetChange("domain")
		_ = data.Set("domain", oldDomain)
		err := errors.New("Domain cannot be changed.")
		diags = append(diags, diag.FromErr(err)...)
		canNotChange = true
	}
	if data.HasChange("rule_id") {
		// 把rule_id强制刷回旧值
		oldRuleId, _ := data.GetChange("rule_id")
		_ = data.Set("rule_id", oldRuleId)
		err := errors.New("WAF Rule ID cannot be changed.")
		diags = append(diags, diag.FromErr(err)...)
		canNotChange = true

	}
	if data.HasChange("type") {
		// 把type强制刷回旧值
		oldType, _ := data.GetChange("type")
		_ = data.Set("type", oldType)
		err := errors.New("Type cannot be changed.")
		diags = append(diags, diag.FromErr(err)...)
		canNotChange = true
	}
	if canNotChange {
		return diags
	}
	request := &securityPolicy.UpdateExceptionForWAFManagedRulesRequest{}
	if id, ok := data.Get("id").(string); ok && id != "" {
		request.Id = &id
	}
	if matchType, ok := data.Get("match_type").(string); ok && matchType != "" {
		request.MatchType = &matchType
	}

	if contentList, ok := data.GetOk("content_list"); ok {
		contents := contentList.([]interface{})
		contentsStr := make([]*string, len(contents))
		for i, v := range contents {
			str := v.(string)
			contentsStr[i] = &str
		}
		request.ContentList = contentsStr
	}

	var response *securityPolicy.UpdateExceptionForWAFManagedRulesResponse
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseSecurityPolicyClient().UpdateExceptionForWAFManagedRules(request)
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
	log.Printf("resource.waap_waf_rule_exception.update success")
	return readWafRuleException(context, data, meta)
}

func deleteWafRuleException(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.waap_waf_rule_exception.delete")
	request := &securityPolicy.DeleteExceptionForWAFManagedRulesRequest{}
	var response *securityPolicy.DeleteExceptionForWAFManagedRulesResponse
	var err error
	var diags diag.Diagnostics

	// 填充request
	dtoList := make([]*securityPolicy.DeleteExceptionForWAFManagedRulesRequestDelDTOList, 0)
	dto := &securityPolicy.DeleteExceptionForWAFManagedRulesRequestDelDTOList{}
	if domain, ok := data.Get("domain").(string); ok && domain != "" {
		dto.Domain = &domain
	}
	if ruleId, ok := data.Get("rule_id").(int); ok {
		dto.RuleId = &ruleId
	}
	exceptionId := data.Id()
	dto.SetExceptionIdList([]*string{&exceptionId})
	dtoList = append(dtoList, dto)
	request.DelDTOList = dtoList

	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseSecurityPolicyClient().DeleteExceptionForWAFManagedRules(request)
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
