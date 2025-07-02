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
	"strings"
	"time"
)

func DataSourceWangSuCdnEdgeHostnames() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWangSuCdnEdgeHostnamesRead,
		Schema: map[string]*schema.Schema{
			"edge_hostnames": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "Filter by edgeHostname.",
			},
			"hostnames": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "Filter by hostname. If specified, only edge-hostname with this hostname will be returned.",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "comment.",
			},
			"dns_service_status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: wangsuCommon.ValidateAllowedStringValue([]string{"inactive", "active"}),
				Description:  "DNS service status. Data range: [inactive, active].",
			},
			"deploy_status": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "Deploy status. Data range: [pending, deploying, success, fail].",
			},
			"is_like": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, the edge-hostname, hostname and comment will be matched using a LIKE query. If false, an exact match is performed.",
			},
			"allow_china_cdn": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: wangsuCommon.ValidateAllowedStringValue([]string{"0", "1"}),
				Description:  "Allow China CDN. Data range: [0,1]. 0 means not allowed, 1 means allowed.",
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
				Description:  "Order of edge-hostname to return. Enum: asc,desc Default: desc",
			},
			"sort_by": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "lastUpdateTime",
				ValidateFunc: wangsuCommon.ValidateAllowedStringValue([]string{"creationTime", "lastUpdateTime", "edgeHostname"}),
				Description:  "Returns results in sorted order. Enum: creationTime,lastUpdateTime Default: lastUpdateTime",
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
							Description: "Number of properties.",
						},
						"edge_hostnames": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of edge-hostname.",
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
										Description: "Edge-Hostname.",
									},
									"dns_service_status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "DNS service status. Data range: [inactive, active].",
									},
									"deploy_status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Deploy status. Data range: [pending, deploying, success, fail].",
									},
									"allow_china_cdn": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Allow China CDN. Data range: [0,1].",
									},
									"comment": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Edge-Hostname comment.",
									},
									"creation_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "RFC3339 format date indicating when the Edge-Hostname was created.",
									},
									"last_update_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "RFC3339 format date indicating when the Edge-Hostname was last updated.",
									},
									"hostnames": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "List of domain names associated with this edge-hostname.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"hostname": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Name of the domain which uses this edge-hostname.",
												},
												"target": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The deployment target of this hostname.",
												},
											},
										},
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

func dataSourceWangSuCdnEdgeHostnamesRead(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("data_source.wangsu_cdn_edge_hostnames.read")

	parameters := &edgehostname.QueryEdgeHostnamesForTerraformParameters{}
	if edgeHostnames, ok := data.Get("edge_hostnames").([]interface{}); ok && len(edgeHostnames) > 0 {
		edgeHostnameList := make([]string, 0, len(edgeHostnames))
		for _, v := range edgeHostnames {
			if t, ok := v.(string); ok && t != "" {
				edgeHostnameList = append(edgeHostnameList, t)
			}
		}
		edgeHostnamesValue := strings.Join(edgeHostnameList, ",")
		parameters.EdgeHostnames = &edgeHostnamesValue
	}
	if hostnames, ok := data.Get("hostnames").([]interface{}); ok && len(hostnames) > 0 {
		hostnameList := make([]string, 0, len(hostnames))
		for _, v := range hostnames {
			if t, ok := v.(string); ok && t != "" {
				hostnameList = append(hostnameList, t)
			}
		}
		hostnamesValue := strings.Join(hostnameList, ",")
		parameters.Hostnames = &hostnamesValue
	}
	if comment, ok := data.Get("comment").(string); ok && comment != "" {
		parameters.Comment = &comment
	}
	if dnsServiceStatus, ok := data.Get("dns_service_status").(string); ok && dnsServiceStatus != "" {
		parameters.DnsServiceStatus = &dnsServiceStatus
	}
	if deployStatuses, ok := data.Get("deploy_status").([]interface{}); ok && len(deployStatuses) > 0 {
		deployStatusList := make([]string, 0, len(deployStatuses))
		for _, v := range deployStatuses {
			if t, ok := v.(string); ok && t != "" {
				deployStatusList = append(deployStatusList, t)
			}
		}
		deployStatusValue := strings.Join(deployStatusList, ",")
		parameters.DeployStatus = &deployStatusValue
	}
	if isLike, ok := data.Get("is_like").(bool); ok && isLike {
		parameters.SetIsLike("true")
	} else {
		parameters.SetIsLike("false")
	}
	if allowChinaCdn, ok := data.Get("allow_china_cdn").(string); ok && allowChinaCdn != "" {
		parameters.AllowChinaCdn = &allowChinaCdn
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

	var response *edgehostname.QueryEdgeHostnamesForTerraformResponse
	var diags diag.Diagnostics
	var err error
	var requestId string
	err = resource.RetryContext(context, time.Duration(1)*time.Minute, func() *resource.RetryError {
		requestId, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseEdgeHostnameClient().QueryEdgeHostnames(parameters)
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

	var edgeHostnameList []interface{}
	for _, edgeHostname := range responseData.EdgeHostnames {
		edgeHostnameDetail := map[string]interface{}{
			"edge_hostname_id":   edgeHostname.EdgeHostnameId,
			"edge_hostname":      edgeHostname.EdgeHostname,
			"dns_service_status": edgeHostname.DnsServiceStatus,
			"deploy_status":      edgeHostname.DeployStatus,
			"allow_china_cdn":    edgeHostname.AllowChinaCdn,
			"comment":            edgeHostname.Comment,
			"creation_time":      edgeHostname.CreationTime,
			"last_update_time":   edgeHostname.LastUpdateTime,
			"hostnames":          buildHostnames(edgeHostname.Hostnames),
		}
		edgeHostnameList = append(edgeHostnameList, edgeHostnameDetail)
	}

	var resultList = []interface{}{
		map[string]interface{}{
			"count":          responseData.Count,
			"edge_hostnames": edgeHostnameList,
		},
	}

	_ = data.Set("data", resultList)

	data.SetId(requestId)
	log.Printf("data_source.wangsu_cdn_edge_hostnames.read success")
	return nil
}

func buildHostnames(hostnames []*edgehostname.QueryEdgeHostnamesForTerraformResponseDataEdgeHostnamesHostnames) []interface{} {
	var hostnameList []interface{}
	if hostnames == nil {
		return hostnameList
	}
	for _, hostname := range hostnames {
		hostnameDetail := map[string]interface{}{
			"hostname": hostname.Hostname,
			"target":   hostname.Target,
		}
		hostnameList = append(hostnameList, hostnameDetail)
	}
	return hostnameList
}
