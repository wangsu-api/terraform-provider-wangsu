package edgehostname

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wangsuCommon "github.com/wangsu-api/terraform-provider-wangsu/wangsu/common"
	"github.com/wangsu-api/wangsu-sdk-go/wangsu/edgehostname"
	"golang.org/x/net/context"
	"log"
	"time"
)

func DataSourceWangSuCdnEdgeHostnameDetail() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWangSuCdnEdgeHostnameDetailRead,
		Schema: map[string]*schema.Schema{
			"edge_hostname": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "edge-hostname",
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
						"edge_hostname_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Edge-Hostname ID.",
						},
						"edge_hostname": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Edge-Hostname Name.",
						},
						"comment": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Edge-Hostname comment.",
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
						"allow_china_cdn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Allow China CDN; values: [0,1].",
						},
						"gdpr_compliant": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "GDPR compliant; values: [0,1,2].",
						},
						"geo_fence": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Geo-fence; values: [global, inside_china_mainland, exclude_china_mainland].",
						},
						"creation_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time in RFC3339 format.",
						},
						"last_update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Last update time in RFC3339 format.",
						},
						"hostnames": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Associated hostnames.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hostname": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Hostname.",
									},
									"target": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Deployment target.",
									},
									"property_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Property ID.",
									},
									"property_version": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Property version.",
									},
									"property_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Property name.",
									},
								},
							},
						},
						"region_configs": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Region configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Region ID.",
									},
									"action_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Action type of configuration.",
									},
									"config_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration type.",
									},
									"config_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration value.",
									},
									"ip_protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "IP protocol.",
									},
									"ttl": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "TTL value.",
									},
									"weight": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Weight of the configuration.",
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

func dataSourceWangSuCdnEdgeHostnameDetailRead(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("data_source.wangsu_cdn_edge_hostname_detail.read")

	edgeHostname := data.Get("edge_hostname").(string)
	var response *edgehostname.QueryEdgeHostnameForTerraformResponse
	var diags diag.Diagnostics
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

	if response == nil {
		data.SetId("")
		return nil
	}
	responseData := response.Data
	if responseData == nil {
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
	var edgeHostnameDetail = map[string]interface{}{
		"edge_hostname_id":   responseData.EdgeHostnameId,
		"edge_hostname":      responseData.EdgeHostname,
		"comment":            responseData.Comment,
		"dns_service_status": responseData.DnsServiceStatus,
		"deploy_status":      responseData.DeployStatus,
		"allow_china_cdn":    responseData.AllowChinaCdn,
		"gdpr_compliant":     responseData.GdprCompliant,
		"geo_fence":          responseData.GeoFence,
		"creation_time":      responseData.CreationTime,
		"last_update_time":   responseData.LastUpdateTime,
	}
	if len(responseData.Hostnames) > 0 {
		var hostnameList []interface{}
		for _, hostname := range responseData.Hostnames {
			hostnameDetail := map[string]interface{}{
				"hostname":         hostname.Hostname,
				"target":           hostname.Target,
				"property_id":      hostname.PropertyId,
				"property_version": hostname.PropertyVersion,
				"property_name":    hostname.PropertyName,
			}
			hostnameList = append(hostnameList, hostnameDetail)
		}
		edgeHostnameDetail["hostnames"] = hostnameList
	}
	if len(responseData.RegionConfigs) > 0 {
		var regionConfigList []interface{}
		for _, regionConfig := range responseData.RegionConfigs {
			regionConfigDetail := map[string]interface{}{
				"region_id":    regionConfig.RegionId,
				"action_type":  regionConfig.ActionType,
				"config_type":  regionConfig.ConfigType,
				"config_value": regionConfig.ConfigValue,
				"ip_protocol":  regionConfig.IpProtocol,
				"ttl":          regionConfig.Ttl,
				"weight":       regionConfig.Weight,
			}
			regionConfigList = append(regionConfigList, regionConfigDetail)
		}
		edgeHostnameDetail["region_configs"] = regionConfigList
	}

	resultList = append(resultList, edgeHostnameDetail)

	_ = data.Set("data", resultList)

	data.SetId(*responseData.EdgeHostname)
	log.Printf("data_source.wangsu_cdn_edge_hostname_detail.read success")
	return nil
}
