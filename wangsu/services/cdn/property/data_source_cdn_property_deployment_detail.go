package property

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wangsuCommon "github.com/wangsu-api/terraform-provider-wangsu/wangsu/common"
	"github.com/wangsu-api/wangsu-sdk-go/wangsu/propertyconfig"
	"golang.org/x/net/context"
	"log"
	"strconv"
	"time"
)

func DataSourceWangSuCdnPropertyDeploymentDetail() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWangSuCdnPropertyDeploymentDetailRead,
		Schema: map[string]*schema.Schema{
			"deployment_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "ID of the deployment task",
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
							Description: "An RFC 3339 format date indicating when the task completed.",
						},
						"actions": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of actions related to a deployment.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Describe an action to take. Enum: deploy_property, remove_property.",
									},
									"property_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "ID of the property to deploy or remove.",
									},
									"version": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Version of the property to deploy or remove.",
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

func dataSourceWangSuCdnPropertyDeploymentDetailRead(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("data_source.wangsu_cdn_property_deployment_detail.read")

	deploymentId := data.Get("deployment_id").(int)
	var response *propertyconfig.QueryDeploymentForTerraformResponse
	var diags diag.Diagnostics
	var err error
	err = resource.RetryContext(context, time.Duration(1)*time.Minute, func() *resource.RetryError {
		response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UsePropertyConfigClient().QueryDeployment(deploymentId)
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
	var resultList []interface{}
	var deploymentDetail = map[string]interface{}{
		"deployment_id":    response.Data.DeploymentId,
		"deployment_name":  response.Data.DeploymentName,
		"status":           response.Data.Status,
		"target":           response.Data.Target,
		"submission_time":  response.Data.SubmissionTime,
		"last_update_time": response.Data.LastUpdateTime,
		"finish_time":      response.Data.FinishTime,
	}
	if len(response.Data.Actions) > 0 {
		var actions []interface{}
		for _, action := range response.Data.Actions {
			actionDetail := map[string]interface{}{
				"action":      action.Action,
				"property_id": action.PropertyId,
				"version":     action.Version,
			}
			actions = append(actions, actionDetail)
		}
		deploymentDetail["actions"] = actions
	}

	resultList = append(resultList, deploymentDetail)

	_ = data.Set("data", resultList)

	data.SetId(strconv.Itoa(*response.Data.DeploymentId))
	log.Printf("data_source.wangsu_cdn_property_deployment_detail.read success")
	return nil
}
