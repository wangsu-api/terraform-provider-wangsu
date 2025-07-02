package edgehostname

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wangsuCommon "github.com/wangsu-api/terraform-provider-wangsu/wangsu/common"
	"github.com/wangsu-api/wangsu-sdk-go/wangsu/edgehostname"
	"golang.org/x/net/context"
	"log"
	"time"
)

func ResourceCdnEdgeHostname() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCdnEdgeHostnameCreate,
		ReadContext:   resourceCdnEdgeHostnameRead,
		UpdateContext: resourceCdnEdgeHostnameUpdate,
		DeleteContext: resourceCdnEdgeHostnameDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"edge_hostname": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Edge-Hostname name, must be unique.",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Edge-Hostname comment.",
			},
			"geo_fence": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "GeoFence, data range: [global, inside_china_mainland, exclude_china_mainland].",
			},
			"region_configs": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Region configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region_id": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Region ID.",
						},
						"action_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Action type.",
						},
						"config_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Config value.",
						},
						"ip_protocol": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "IP protocol.",
						},
						"ttl": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "TTL (Time to Live).",
						},
						"weight": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Weight.",
						},
					},
				},
			},
			//computed
			"edge_hostname_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Edge-Hostname ID.",
			},
			"dns_service_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "DNS service status; possible values: [inactive, active].",
			},
			"deploy_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Deploy status; possible values: [pending, deploying, success, fail].",
			},
		},
	}
}

func resourceCdnEdgeHostnameRead(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_edge_hostname.read")
	var diags diag.Diagnostics
	edgeHostname := data.Id()

	var response *edgehostname.QueryEdgeHostnameForTerraformResponse
	var err error
	err = resource.RetryContext(context, time.Duration(1)*time.Minute, func() *resource.RetryError {
		response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseEdgeHostnameClient().QueryEdgeHostname(edgeHostname)
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

	if err := data.Set("edge_hostname", response.Data.EdgeHostname); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("comment", response.Data.Comment); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("geo_fence", response.Data.GeoFence); err != nil {
		return diag.FromErr(err)
	}
	regionConfigList := make([]interface{}, 0, len(response.Data.RegionConfigs))
	for _, regionConfig := range response.Data.RegionConfigs {
		regionConfigMap := map[string]interface{}{
			"region_id":    regionConfig.RegionId,
			"action_type":  regionConfig.ActionType,
			"config_value": regionConfig.ConfigValue,
			"ip_protocol":  regionConfig.IpProtocol,
			"ttl":          regionConfig.Ttl,
			"weight":       regionConfig.Weight,
		}
		regionConfigList = append(regionConfigList, regionConfigMap)
	}
	if err := data.Set("region_configs", regionConfigList); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("edge_hostname_id", response.Data.EdgeHostnameId); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("dns_service_status", response.Data.DnsServiceStatus); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("deploy_status", response.Data.DeployStatus); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceCdnEdgeHostnameCreate(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_edge_hostname.create")
	edgeHostname := data.Get("edge_hostname").(string)
	return executeUpdate(edgeHostname, context, data, meta)
}

func executeUpdate(edgeHostname string, context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	request := &edgehostname.UpdateEdgeHostnameForTerraformRequest{}
	if comment, ok := data.Get("comment").(string); ok && comment != "" {
		request.Comment = &comment
	}
	if geoFence, ok := data.Get("geo_fence").(string); ok && geoFence != "" {
		request.GeoFence = &geoFence
	}
	if regionConfigs, ok := data.Get("region_configs").([]interface{}); ok && len(regionConfigs) > 0 {

		regionConfigList := make([]*edgehostname.UpdateEdgeHostnameForTerraformRequestRegionConfigs, 0, len(regionConfigs))
		for _, regionConfig := range regionConfigs {
			regionConfigMap := regionConfig.(map[string]interface{})
			regionConfigItem := &edgehostname.UpdateEdgeHostnameForTerraformRequestRegionConfigs{}
			if regionId, ok := regionConfigMap["region_id"].(int); ok && regionId > 0 {
				regionConfigItem.RegionId = &regionId
			}
			if actionType, ok := regionConfigMap["action_type"].(string); ok && actionType != "" {
				regionConfigItem.ActionType = &actionType
			}
			if configValue, ok := regionConfigMap["config_value"].(string); ok && configValue != "" {
				regionConfigItem.ConfigValue = &configValue
			}
			if ipProtocol, ok := regionConfigMap["ip_protocol"].(string); ok && ipProtocol != "" {
				regionConfigItem.IpProtocol = &ipProtocol
			}
			if ttl, ok := regionConfigMap["ttl"].(int); ok && ttl > 0 {
				regionConfigItem.Ttl = &ttl
			}
			if weight, ok := regionConfigMap["weight"].(int); ok && weight > 0 {
				regionConfigItem.Weight = &weight
			}
			regionConfigList = append(regionConfigList, regionConfigItem)
		}
		request.RegionConfigs = regionConfigList
	}
	var response *edgehostname.UpdateEdgeHostnameForTerraformResponse
	var err error
	err = resource.RetryContext(context, time.Duration(1)*time.Minute, func() *resource.RetryError {
		response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseEdgeHostnameClient().UpdateEdgeHostname(edgeHostname, request)
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
	var deployResponse *edgehostname.DeployEdgeHostnameForTerraformResponse
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		deployResponse, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseEdgeHostnameClient().DeployEdgeHostname(edgeHostname)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}
	if deployResponse == nil {
		data.SetId("")
		return nil
	}

	var readResponse *edgehostname.QueryEdgeHostnameForTerraformResponse
	err = resource.RetryContext(context, time.Duration(12)*time.Hour, func() *resource.RetryError {
		readResponse, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseEdgeHostnameClient().QueryEdgeHostname(edgeHostname)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if readResponse != nil && readResponse.Data != nil && *readResponse.Data.DeployStatus == "fail" {
			return resource.NonRetryableError(fmt.Errorf("edge-hostname deployment failed"))
		}
		if readResponse != nil && readResponse.Data != nil && *readResponse.Data.DeployStatus != "success" {
			return resource.RetryableError(fmt.Errorf("edge-hostname is in progress, retrying"))
		}
		return nil
	})

	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	data.SetId(edgeHostname)

	return resourceCdnEdgeHostnameRead(context, data, meta)
}

func resourceCdnEdgeHostnameUpdate(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_edge_hostname.update")
	return executeUpdate(data.Id(), context, data, meta)
}

func resourceCdnEdgeHostnameDelete(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_edge_hostname.delete")
	var diags diag.Diagnostics
	edgeHostname := data.Id()

	var response *edgehostname.DeleteEdgeHostnameForTerraformResponse
	var err error
	err = resource.RetryContext(context, time.Duration(3)*time.Minute, func() *resource.RetryError {
		response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseEdgeHostnameClient().DeleteEdgeHostname(edgeHostname)
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
		return diags
	}
	log.Printf("resource.wangsu_edge_hostname.delete success, edgeHostname: %s", edgeHostname)
	return diags
}
