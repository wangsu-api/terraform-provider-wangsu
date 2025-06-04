package policy

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wangsuCommon "github.com/wangsu-api/terraform-provider-wangsu/wangsu/common"
	policy "github.com/wangsu-api/wangsu-sdk-go/wangsu/policy"
	"log"
	"strconv"
	"time"
)

func ResourceIamPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIamPolicyCreate,
		ReadContext:   resourceIamPolicyRead,
		UpdateContext: resourceIamPolicyUpdate,
		DeleteContext: resourceIamPolicyDelete,

		Schema: map[string]*schema.Schema{
			"policy_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Policy Name. Supports Chinese, English, and underline, with no more than 150 characters",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Policy description. You may describe the strategy here, including permission details, limited to a maximum of 250 characters.",
			},
			"policy_document": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Permission policy content, please fill in the permission policy script as follows:\n\n\n[{\n        \\\"effect\\\": \\\"allow\\\",\n        \\\"action\\\": [\n           \\\"productCode:actionCode\\\"\n        ],\n       \\\"resource\\\": [\n           \\\"*\\\"\n        ]\n    }]\n\n\nA single permission policy can include permissions for multiple products, but CDN and non-CDN product permissions cannot be added to the same policy simultaneously.\n\nField descriptions:\n\n-effect: The authorization effect includes two types: allow and deny.\n\n-action: Describes the specific operations allowed or denied,  format: productCode:actionCode.\n\n-resource: The specific resources authorized. For all resources use *, for specific resources refer to the format: wsc:&lt;service-name&gt;:&lt;region&gt;:&lt;account&gt;:&lt;relatice-id&gt;. Note: CDN products do not support specifying resources.",
			},
		},
	}
}
func resourceIamPolicyCreate(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_iam_policy.create")
	var diags diag.Diagnostics
	request := &policy.CreatePolicyRequest{}
	if policyName, ok := data.Get("policy_name").(string); ok {
		request.PolicyName = &policyName
	}
	if description, ok := data.Get("description").(string); ok {
		request.Description = &description
	}
	if policyDocument, ok := data.Get("policy_document").(string); ok {
		request.PolicyDocument = &policyDocument
	}

	var response *policy.CreatePolicyResponse
	var requestId string
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		requestId, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UsePolicyClient().AddPolicy(request)
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
	data.SetId(strconv.FormatInt(*response.Data.PolicyId, 10))
	log.Printf("resource.wangsu_iam_policy.create success")
	log.Printf("requestId: %s", requestId)
	time.Sleep(2 * time.Second)
	return resourceIamPolicyRead(context, data, meta)
}

func resourceIamPolicyRead(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_iam_policy.read")
	var diags diag.Diagnostics
	request := &policy.GetPolicyRequest{}
	policyId, err1 := strconv.ParseInt(data.Id(), 10, 64)
	if err1 != nil {
		return diag.FromErr(err1)
	}
	request.PolicyId = &policyId
	var response *policy.GetPolicyResponse
	var requestId string
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		requestId, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UsePolicyClient().GetPolicy(request)
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
	_ = data.Set("policy_name", response.Data.PolicyName)
	_ = data.Set("description", response.Data.Description)
	_ = data.Set("policy_document", response.Data.PolicyDocument)
	log.Printf("resource.wangsu_iam_policy.read success, requestId: %s", requestId)
	return diags
}

func resourceIamPolicyUpdate(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_iam_policy.update")
	var diags diag.Diagnostics
	request := &policy.EditPolicyRequest{}
	policyId := data.Id()
	request.PolicyId = &policyId
	if policyName, ok := data.Get("policy_name").(string); ok {
		request.PolicyName = &policyName
	}
	// 需要支持可以传""，此处不判断是否存在
	description, _ := data.Get("description").(string)
	request.Description = &description
	if policyDocument, ok := data.Get("policy_document").(string); ok {
		request.PolicyDocument = &policyDocument
	}

	var response *policy.EditPolicyResponse
	var requestId string
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		requestId, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UsePolicyClient().EditPolicy(request)
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
	log.Printf("resource.wangsu_iam_policy.update success")
	log.Printf("requestId: %s", requestId)
	time.Sleep(2 * time.Second)
	return resourceIamPolicyRead(context, data, meta)
}

func resourceIamPolicyDelete(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_iam_policy.delete")
	var diags diag.Diagnostics
	request := &policy.DeletePolicyRequest{}
	policyId := data.Id()
	request.PolicyId = &policyId
	var response *policy.DeletePolicyResponse
	var requestId string
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		requestId, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UsePolicyClient().DeletePolicy(request)
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
	log.Printf("resource.wangsu_iam_policy.delete success")
	log.Printf("requestId: %s", requestId)
	return nil
}
