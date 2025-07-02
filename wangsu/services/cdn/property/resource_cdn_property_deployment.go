package property

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wangsuCommon "github.com/wangsu-api/terraform-provider-wangsu/wangsu/common"
	propertyConfig "github.com/wangsu-api/wangsu-sdk-go/wangsu/propertyconfig"
	"golang.org/x/net/context"
	"log"
	"strconv"
	"time"
)

func ResourceCdnPropertyDeployment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCdnPropertyDeploymentCreate,
		ReadContext:   resourceCdnPropertyDeploymentRead,
		DeleteContext: resourceCdnPropertyDeploymentDelete,

		Schema: map[string]*schema.Schema{
			"deployment_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name representing the deployment task.",
			},
			"target": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Indicates whether to deploy to staging or production. Enum: staging,production",
			},
			"actions": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "This array contains all the actions related to a deployment. They can include deployment and removal of properties to the staging or production environments.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Describe an action to take. You can deploy a property, remove a property. Enum: deploy_property,remove_property",
						},
						"property_id": {
							Type:        schema.TypeInt,
							Required:    true,
							ForceNew:    true,
							Description: "ID of the property to deploy or remove from the staging or production environment.",
						},
						"version": {
							Type:        schema.TypeInt,
							Required:    true,
							ForceNew:    true,
							Description: "Indicates the version of the property to deploy or remove.",
						},
					},
				},
			},
			//computed
			"deployment_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "ID of the deployment task.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of Deployment. Enum: PENDING, IN_PROCESS, SUCCESS, FAIL.",
			},
		},
	}
}
func resourceCdnPropertyDeploymentRead(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_cdn_property_deployment.read")

	var diags diag.Diagnostics
	deploymentId, err := strconv.Atoi(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	var response *propertyConfig.QueryDeploymentForTerraformResponse
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

	if response == nil || response.Data == nil {
		data.SetId("")
		return nil
	}

	if err := data.Set("deployment_name", response.Data.DeploymentName); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("target", response.Data.Target); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("actions", flattenActions(response.Data.Actions)); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("deployment_id", response.Data.DeploymentId); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("status", response.Data.Status); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("resource.wangsu_cdn_property_deployment.read success, deploymentId: %d", response.Data.DeploymentId)
	return nil
}

func flattenActions(actions []*propertyConfig.QueryDeploymentForTerraformResponseDataActions) []interface{} {
	if actions == nil || len(actions) == 0 {
		return []interface{}{}
	}
	resultList := make([]interface{}, len(actions))
	for i, action := range actions {
		resultList[i] = map[string]interface{}{
			"action":      action.Action,
			"property_id": action.PropertyId,
			"version":     action.Version,
		}
	}
	return resultList
}

func resourceCdnPropertyDeploymentCreate(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_cdn_property_deployment.create")

	request := &propertyConfig.CreateDeploymentTaskForTerraformRequest{}
	if deploymentName, ok := data.Get("deployment_name").(string); ok && deploymentName != "" {
		request.DeploymentName = &deploymentName
	}
	if target, ok := data.Get("target").(string); ok && target != "" {
		request.Target = &target
	}
	if actions, ok := data.Get("actions").([]interface{}); ok && len(actions) > 0 {
		actionList := make([]*propertyConfig.CreateDeploymentTaskForTerraformRequestActions, 0, len(actions))
		for _, action := range actions {
			actionMap := action.(map[string]interface{})
			actionItem := &propertyConfig.CreateDeploymentTaskForTerraformRequestActions{}
			if actionStr, ok := actionMap["action"].(string); ok && actionStr != "" {
				actionItem.Action = &actionStr
			}
			if propertyId, ok := actionMap["property_id"].(int); ok && propertyId > 0 {
				actionItem.PropertyId = &propertyId
			}
			if version, ok := actionMap["version"].(int); ok && version > 0 {
				actionItem.Version = &version
			}
			actionList = append(actionList, actionItem)
		}
		request.Actions = actionList
	}

	var diags diag.Diagnostics
	var response *propertyConfig.CreateDeploymentTaskForTerraformResponse
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UsePropertyConfigClient().CreateDeployment(request)
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
	var deploymentId = *response.Data.DeploymentId

	time.Sleep(2 * time.Second)

	var deploymentResponse *propertyConfig.QueryDeploymentForTerraformResponse
	err = resource.RetryContext(context, time.Duration(12)*time.Hour, func() *resource.RetryError {
		deploymentResponse, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UsePropertyConfigClient().QueryDeployment(deploymentId)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if deploymentResponse != nil && deploymentResponse.Data != nil && *deploymentResponse.Data.Status != "SUCCESS" {
			return resource.RetryableError(fmt.Errorf("property deployment status is in progress, retrying"))
		}
		return nil
	})
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}
	data.SetId(strconv.Itoa(deploymentId))
	log.Printf("resource.wangsu_cdn_property_deployment.create success, deploymentId: %d", response.Data.DeploymentId)
	return resourceCdnPropertyDeploymentRead(context, data, meta)
}

func resourceCdnPropertyDeploymentDelete(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_cdn_property_deployment.delete do nothing")
	return nil
}
