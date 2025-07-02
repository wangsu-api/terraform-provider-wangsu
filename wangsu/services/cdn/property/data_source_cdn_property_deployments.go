package property

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wangsuCommon "github.com/wangsu-api/terraform-provider-wangsu/wangsu/common"
	"github.com/wangsu-api/wangsu-sdk-go/wangsu/propertyconfig"
	"golang.org/x/net/context"
	"log"
	"time"
)

func DataSourceWangSuCdnPropertyDeployments() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWangSuCdnPropertyDeploymentsRead,
		Schema: map[string]*schema.Schema{
			"property_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Property ID.",
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: wangsuCommon.ValidateAllowedStringValue([]string{"PENDING", "PENDING_REVIEW", "IN_PROCESS", "SUCCESS", "FAIL"}),
				Description:  "Status of Deployment. Enum:PENDING,IN_PROCESS,SUCCESS,FAIL",
			},
			"target": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: wangsuCommon.ValidateAllowedStringValue([]string{"staging", "production"}),
				Description:  "The value can be 'staging', or 'production' to filter the results based on where the property has been deployed.",
			},
			"offset": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(int)
					if value < 0 {
						errors = append(errors, fmt.Errorf("%q must be a non-negative integer", k))
					}
					return
				},
				Description: "Indicates the first item to return. The default is '0'.",
			},
			"limit": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  100,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(int)
					if value < 0 || value > 200 {
						errors = append(errors, fmt.Errorf("%q must be a non-negative integer less than or equal to 200", k))
					}
					return
				},
				Description: "Maximum number of properties to return.  Default: 100 Range: <= 200",
			},
			"sort_order": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "desc",
				ValidateFunc: wangsuCommon.ValidateAllowedStringValue([]string{"asc", "desc"}),
				Description:  "Order of deployments to return. Enum: asc,desc Default: desc",
			},
			"sort_by": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "submissionTime",
				ValidateFunc: wangsuCommon.ValidateAllowedStringValue([]string{"submissionTime", "lastUpdateTime"}),
				Description:  "Returns results in sorted order. Enum: submissionTime,lastUpdateTime Default: submissionTime",
			},
			//computed
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
						"count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of deployment task.",
						},
						"deployments": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of deployment task summaries.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"deployment_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "ID representing the deployment task.",
									},
									"deployment_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name representing the deployment task.",
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Status of Deployment. Enum: PENDING, IN_PROCESS, SUCCESS, FAIL.",
									},
									"target": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The value can be 'staging', or 'production'.",
									},
									"submission_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "An RFC 3339 format date indicating when the task was submitted.",
									},
									"last_update_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "An RFC 3339 format date indicating when the task was last updated.",
									},
									"finish_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "An RFC 3339 format date indicating when the task was completed.",
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

func dataSourceWangSuCdnPropertyDeploymentsRead(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("data_source.wangsu_cdn_property_deployments.read")

	parameters := &propertyconfig.QueryDeploymentsForTerraformParameters{}
	if propertyId, ok := data.Get("property_id").(int); ok && propertyId > 0 {
		parameters.PropertyId = &propertyId
	}
	if status, ok := data.Get("status").(string); ok && status != "" {
		parameters.Status = &status
	}
	if target, ok := data.Get("target").(string); ok && target != "" {
		parameters.Target = &target
	}
	if offset, ok := data.Get("offset").(int); ok && offset >= 0 {
		parameters.Offset = &offset
	}
	if limit, ok := data.Get("limit").(int); ok && limit >= 0 && limit <= 200 {
		parameters.Limit = &limit
	}
	if sortOrder, ok := data.Get("sort_order").(string); ok && sortOrder != "" {
		parameters.SortOrder = &sortOrder
	}
	if sortBy, ok := data.Get("sort_by").(string); ok && sortBy != "" {
		parameters.SortBy = &sortBy
	}

	var response *propertyconfig.QueryDeploymentsForTerraformResponse
	var diags diag.Diagnostics
	var err error
	var requestId string
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		requestId, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UsePropertyConfigClient().QueryDeployments(parameters)
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
	if response.Data == nil {
		data.SetId("")
		return nil
	}

	if err := data.Set("code", response.Code); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("message", response.Message); err != nil {
		return diag.FromErr(err)
	}

	responseData := response.Data
	var deployments []interface{}
	for _, deployment := range responseData.Deployments {
		deploymentDetail := map[string]interface{}{
			"deployment_id":    deployment.DeploymentId,
			"deployment_name":  deployment.DeploymentName,
			"status":           deployment.Status,
			"target":           deployment.Target,
			"submission_time":  deployment.SubmissionTime,
			"last_update_time": deployment.LastUpdateTime,
			"finish_time":      deployment.FinishTime,
		}
		deployments = append(deployments, deploymentDetail)
	}
	var resultList = []interface{}{
		map[string]interface{}{
			"count":       responseData.Count,
			"deployments": deployments,
		},
	}

	_ = data.Set("data", resultList)

	data.SetId(requestId)
	log.Printf("data_source.wangsu_cdn_property_deployments.read success")
	return nil
}
