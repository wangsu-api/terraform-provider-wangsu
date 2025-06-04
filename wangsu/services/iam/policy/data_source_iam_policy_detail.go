package policy

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wangsuCommon "github.com/wangsu-api/terraform-provider-wangsu/wangsu/common"
	"github.com/wangsu-api/wangsu-sdk-go/wangsu/policy"
	"golang.org/x/net/context"
	"log"
	"strconv"
	"time"
)

func ResourceIamPolicyDetail() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIamPolicyDetailRead,
		Schema: map[string]*schema.Schema{
			"policy_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Policy name",
			},
			"data": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Detailed data on the results of the request",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy name",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy description",
						},
						"policy_document": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Permission policy content, as follows:\n\n[{\n        \\\"effect\\\": \\\"allow\\\",\n        \\\"action\\\": [\n           \\\"productCode:actionCode\\\"\n        ],\n       \\\"resource\\\": [\n           \\\"*\\\"\n        ]\n    }]\n\nField descriptions:\n\n-effect: The authorization effect includes two types: allow and deny.\n\n-action: Describes the specific operations allowed or denied,  format: productCode:actionCode.\n\n-resource: The specific resources authorized. For all resources use *, for specific resources refer to the format: wsc:&lt;service-name&gt;:&lt;region&gt;:&lt;account&gt;:&lt;relatice-id&gt;. Note: CDN products do not support specifying resources.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIamPolicyDetailRead(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("data_source.wangsu_iam_policy_detail.read")
	var diags diag.Diagnostics
	var policyName string
	if v, ok := data.Get("policy_name").(string); ok {
		policyName = v
	}
	request := &policy.GetPolicyRequest{}
	request.PolicyName = &policyName
	var response *policy.GetPolicyResponse
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UsePolicyClient().GetPolicy(request)
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
	var resultList []interface{}
	var policyDetail = map[string]interface{}{
		"policy_name":     response.Data.PolicyName,
		"description":     response.Data.Description,
		"policy_document": response.Data.PolicyDocument,
	}
	resultList = append(resultList, policyDetail)
	_ = data.Set("data", resultList)
	data.SetId(strconv.FormatInt(*response.Data.PolicyId, 10))
	log.Printf("data_source.wangsu_iam_policy_detail.read success")
	return diags
}
