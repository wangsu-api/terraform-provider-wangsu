package policy

import (
	"context"
	"errors"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wangsuCommon "github.com/wangsu-api/terraform-provider-wangsu/wangsu/common"
	"log"
	"time"

	policyAttachment "github.com/wangsu-api/wangsu-sdk-go/wangsu/usermanage"
)

func ResourceIamPolicyAttachment() *schema.Resource {

	return &schema.Resource{
		CreateContext: resourceIamPolicyAttachmentCreate,
		ReadContext:   ResourceIamPolicyAttachmentRead,
		UpdateContext: ResourceIamPolicyAttachmentUpdate,
		DeleteContext: ResourceIamPolicyAttachmentDelete,
		Schema: map[string]*schema.Schema{
			"login_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Sub account login name",
			},
			"policy_name": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Specify policy ID",
			},
			"data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "policy name",
						},
					},
				},
			},
		},
	}

}
func resourceIamPolicyAttachmentCreate(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_iam_policy_attachment.create")
	var diags diag.Diagnostics
	request := &policyAttachment.BatchAddOrRevokePolicyToSubAccountRequest{}
	if login_name, ok := data.Get("login_name").(string); ok && login_name != "" {
		request.LoginName = &login_name
	}

	if policyName, ok := data.Get("policy_name").([]interface{}); ok && len(policyName) > 0 {
		PolicyNameList := make([]*string, len(policyName))
		for i, v := range policyName {
			if name, ok := v.(string); ok {
				PolicyNameList[i] = &name
			} else {
				return append(diags, diag.FromErr(errors.New("Invalid policy name type."))...)
			}
		}
		request.PolicyName = PolicyNameList
	}
	request.Type = tea.Int(0)
	var response *policyAttachment.BatchAddOrRevokePolicyToSubAccountResponse
	var requestId string
	var err error
	err = resource.RetryContext(ctx, time.Duration(5)*time.Minute, func() *resource.RetryError {
		requestId, response, err = m.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UsePolicyAttachmentClient().BatchAddOrRevokePolicyToSubAccount(request)
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
	data.SetId(*request.LoginName)
	log.Printf("resource.wangsu_iam_policy.create success")
	log.Printf("requestId: %s", requestId)
	return ResourceIamPolicyAttachmentRead(ctx, data, m)
}
func ResourceIamPolicyAttachmentDelete(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_iam_policy_attachment.create")
	var diags diag.Diagnostics
	request := &policyAttachment.BatchAddOrRevokePolicyToSubAccountRequest{}
	if login_name, ok := data.Get("login_name").(string); ok && login_name != "" {
		request.LoginName = &login_name
	}

	if policyName, ok := data.Get("policy_name").([]interface{}); ok && len(policyName) > 0 {
		PolicyNameList := make([]*string, len(policyName))
		for i, v := range policyName {
			if name, ok := v.(string); ok {
				PolicyNameList[i] = &name
			} else {
				return append(diags, diag.FromErr(errors.New("Invalid policy name type."))...)
			}
		}
		request.PolicyName = PolicyNameList
	}
	request.Type = tea.Int(1)
	var response *policyAttachment.BatchAddOrRevokePolicyToSubAccountResponse
	var requestId string
	var err error
	err = resource.RetryContext(ctx, time.Duration(5)*time.Minute, func() *resource.RetryError {
		requestId, response, err = m.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UsePolicyAttachmentClient().BatchAddOrRevokePolicyToSubAccount(request)
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
	data.SetId(*request.LoginName)
	log.Printf("resource.wangsu_iam_policy.update success")
	log.Printf("requestId: %s", requestId)
	return diags
}
func ResourceIamPolicyAttachmentUpdate(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_iam_policy_attachment.update")
	var diags diag.Diagnostics

	response1, requestId1, err1 := getPolicyAttachment(ctx, data, m)
	log.Printf("getPolicyAttachment requestId: %s", requestId1)

	if err1 != nil {
		diags = append(diags, diag.FromErr(err1)...)
		return diags
	}
	if response1 == nil || response1.Data == nil {
		data.SetId("")
		return nil
	}
	// 获取getPolicyAttachment得到的policyNameList,与data中的policyNameList进行对比,多的做type=0,少的做type=1
	policyNameSetInDB := make(map[string]struct{}, len(response1.Data))
	for _, v := range response1.Data {
		if v.PolicyName != nil {
			policyNameSetInDB[*v.PolicyName] = struct{}{}
		} else {
			return append(diags, diag.FromErr(errors.New("Invalid policy name type."))...)
		}
	}
	policyNameSet := make(map[string]struct{})
	if policyName, ok := data.Get("policy_name").([]interface{}); ok && len(policyName) > 0 {
		for _, v := range policyName {
			if name, ok := v.(string); ok {
				policyNameSet[name] = struct{}{}
			} else {
				return append(diags, diag.FromErr(errors.New("Invalid policy name type."))...)
			}
		}
	}

	// policyNameListInDB中不在policyNameSet中的policyName
	var missingPolicyNames []*string
	for name := range policyNameSetInDB {
		if _, found := policyNameSet[name]; !found {
			missingPolicyNames = append(missingPolicyNames, tea.String(name))
		}
	}
	// policyNameSet中不在policyNameListInDB中的policyName
	var extraPolicyNames []*string
	for name := range policyNameSet {
		if _, found := policyNameSetInDB[name]; !found {
			extraPolicyNames = append(extraPolicyNames, tea.String(name))
		}
	}

	request := &policyAttachment.BatchAddOrRevokePolicyToSubAccountRequest{}
	if loginName, ok := data.Get("login_name").(string); ok && loginName != "" {
		request.LoginName = &loginName
	}

	//missingPolicyNames不为空则执行updatePolicyAttachment
	if len(missingPolicyNames) > 0 {
		response2, requestId2, err2 := updatePolicyAttachment(ctx, request, missingPolicyNames, 1, m)
		if err2 != nil {
			diags = append(diags, diag.FromErr(err2)...)
			return diags
		}
		if response2 == nil {
			data.SetId("")
			return nil
		}
		log.Printf("reclaim policy requestId: %s", requestId2)
	}
	if len(extraPolicyNames) > 0 {
		response3, requestId3, err3 := updatePolicyAttachment(ctx, request, extraPolicyNames, 0, m)
		if err3 != nil {
			diags = append(diags, diag.FromErr(err3)...)
			return diags
		}
		if response3 == nil {
			data.SetId("")
			return nil
		}
		log.Printf("attach policy requestId: %s", requestId3)
	}
	data.SetId(*request.LoginName)
	log.Printf("resource.wangsu_iam_policy.update success")
	return diags
}

func updatePolicyAttachment(ctx context.Context, request *policyAttachment.BatchAddOrRevokePolicyToSubAccountRequest, missingPolicyNames []*string, typ int, m interface{}) (*policyAttachment.BatchAddOrRevokePolicyToSubAccountResponse, string, error) {
	request.PolicyName = missingPolicyNames
	request.Type = tea.Int(typ)
	var response *policyAttachment.BatchAddOrRevokePolicyToSubAccountResponse
	var requestId string
	var err error
	err = resource.RetryContext(ctx, time.Duration(5)*time.Minute, func() *resource.RetryError {
		requestId, response, err = m.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UsePolicyAttachmentClient().BatchAddOrRevokePolicyToSubAccount(request)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	return response, requestId, err
}

func ResourceIamPolicyAttachmentRead(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_iam_policy_attachment.read")
	var diags diag.Diagnostics
	response, requestId, err := getPolicyAttachment(ctx, data, m)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}
	if response == nil || response.Data == nil {
		data.SetId("")
		return nil
	}
	dataList := make([]interface{}, 0)
	for _, item := range response.Data {
		dataList = append(dataList, map[string]interface{}{
			"policy_name": item.PolicyName,
		})
	}
	_ = data.Set("data", dataList)
	log.Printf("resource.wangsu_iam_policy_attachment.read success")
	log.Printf("requestId: %s", requestId)
	return diags
}

func getPolicyAttachment(ctx context.Context, data *schema.ResourceData, m interface{}) (*policyAttachment.QueryPolicyAttachedMainAccountOrSubAccountResponse, string, error) {
	request := &policyAttachment.QueryPolicyAttachedMainAccountOrSubAccountPaths{}
	loginName := data.Get("login_name").(string)
	request.LoginName = &loginName
	var response *policyAttachment.QueryPolicyAttachedMainAccountOrSubAccountResponse
	var requestId string
	var err error
	err = resource.RetryContext(ctx, time.Duration(2)*time.Minute, func() *resource.RetryError {
		requestId, response, err = m.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UsePolicyAttachmentClient().QueryPolicyAttachedMainAccountOrSubAccount(request)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	return response, requestId, err
}
