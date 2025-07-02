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

func DataSourceWangSuCdnPropertyDetail() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWangSuCdnPropertyDetailRead,
		Schema: map[string]*schema.Schema{
			"property_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Property ID",
			},
			"version": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Property Version",
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
				Description: "Response data",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"property_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Property ID",
						},
						"property_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the property.",
						},
						"property_comment": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A description of the property.",
						},
						"service_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Service type.",
						},
						"property_creation_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time in RFC3339 format.",
						},
						"property_last_update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Last update time in RFC3339 format.",
						},
						"staging_version": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Staging version details",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"version": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Version of the property.",
									},
									"hostnames": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Description: "List of hostnames.",
									},
								},
							},
						},
						"production_version": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Production version details",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"version": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Version of the property.",
									},
									"hostnames": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Description: "List of hostnames.",
									},
								},
							},
						},
						"staging_deploying_version": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Staging deploying version details",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"version": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Version of the property.",
									},
									"hostnames": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Description: "List of hostnames.",
									},
								},
							},
						},
						"production_deploying_version": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "production deploying version details",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"version": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Version of the property.",
									},
									"hostnames": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Description: "List of hostnames.",
									},
								},
							},
						},
						"latest_version": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Latest version of the property.",
						},
						"current_version": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "A property version.",
						},
						"version_comment": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A description of the current version.",
						},
						"frozen": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Is the version frozen?",
						},
						"version_creation_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Version creation time in RFC3339 format.",
						},
						"version_last_update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Version last update time in RFC3339 format.",
						},
						"hostnames": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "List of hostnames",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hostname": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Hostname",
									},
									"icp": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Domain ICP information",
									},
									"default_origin": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Default origin",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"host": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Origin host",
												},
												"http_port": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "HTTP port",
												},
												"https_port": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "HTTPS port",
												},
												"ip_version": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "IP version. Data range: dual, ipv4, ipv6",
												},
												"servers": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Origin servers",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
									"certificates": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Hostname association certificate configuration",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"certificate_id": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Certificate ID",
												},
												"certificate_usage": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Certificate usage. Data range: default_sni, ssl_bk, gm_sm2_enc, gm_sm2_sign, client_mtls",
												},
											},
										},
									},
									"edge_hostname": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Hostname edge-hostname config",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"edge_hostname_prefix": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Edge-hostname prefix",
												},
												"edge_hostname_suffix": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Edge-hostname suffix",
												},
												"edge_hostname": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Edge-hostname",
												},
												"dns_service_status": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "DNS service status. Data range: inactive, active",
												},
											},
										},
									},
								},
							},
						},
						"rules": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Rules.",
						},
						"origins": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specify the hostname and settings used to contact the origin once service begins.",
						},
						"variables": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The Variables feature allows you to define variables, assign values to them, and reuse them in functional test cases. A variable consists of a name and a value.",
						},
					},
				},
			},
		},
	}
}

func dataSourceWangSuCdnPropertyDetailRead(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("data_source.wangsu_cdn_property_detail.read")
	propertyId := data.Get("property_id").(int)
	version := data.Get("version").(int)

	var response *propertyconfig.QueryPropertyVersionConfigForTerrformResponse
	var diags diag.Diagnostics
	var err error
	err = resource.RetryContext(context, time.Duration(1)*time.Minute, func() *resource.RetryError {
		response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UsePropertyConfigClient().QueryPropertyVersion(propertyId, version)
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
	if err := data.Set("code", response.Code); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("message", response.Message); err != nil {
		return diag.FromErr(err)
	}

	responseData := response.Data
	if responseData == nil {
		data.SetId("")
		return nil
	}

	var propertyVersion = map[string]interface{}{
		"property_id":                  responseData.PropertyId,
		"property_name":                responseData.PropertyName,
		"property_comment":             responseData.PropertyComment,
		"service_type":                 responseData.ServiceType,
		"property_creation_time":       responseData.PropertyCreationTime,
		"property_last_update_time":    responseData.PropertyLastUpdateTime,
		"staging_version":              buildStagingVersion(responseData.StagingVersion),
		"production_version":           buildProductionVersion(responseData.ProductionVersion),
		"staging_deploying_version":    buildStagingDeployingVersion(responseData.StagingDeployingVersion),
		"production_deploying_version": buildProductionDeployingVersion(responseData.ProductionDeployingVersion),
		"latest_version":               responseData.LatestVersion,
		"current_version":              responseData.CurrentVersion,
		"version_comment":              responseData.VersionComment,
		"frozen":                       responseData.Frozen,
		"version_creation_time":        responseData.VersionCreationTime,
		"version_last_update_time":     responseData.VersionLastUpdateTime,
		"hostnames":                    buildHostnames(responseData.Hostnames),
		"rules":                        responseData.Rules,
		"origins":                      responseData.Origins,
		"variables":                    responseData.Variables,
	}
	var resultList []interface{}
	resultList = append(resultList, propertyVersion)
	if err := data.Set("data", resultList); err != nil {
		return diag.FromErr(err)
	}
	data.SetId(strconv.Itoa(*responseData.PropertyId))
	log.Printf("data_source.wangsu_cdn_property_detail.read success, property_id: %d, version: %d", propertyId, version)
	return nil
}

func buildStagingVersion(versionConfig *propertyconfig.QueryPropertyVersionConfigForTerrformResponseDataStagingVersion) []interface{} {
	if versionConfig == nil {
		return nil
	}
	var resultList []interface{}
	var config = map[string]interface{}{
		"version":   versionConfig.Version,
		"hostnames": versionConfig.Hostnames,
	}
	resultList = append(resultList, config)
	return resultList
}
func buildProductionVersion(versionConfig *propertyconfig.QueryPropertyVersionConfigForTerrformResponseDataProductionVersion) []interface{} {
	if versionConfig == nil {
		return nil
	}
	var resultList []interface{}
	var config = map[string]interface{}{
		"version":   versionConfig.Version,
		"hostnames": versionConfig.Hostnames,
	}
	resultList = append(resultList, config)
	return resultList
}
func buildStagingDeployingVersion(versionConfig *propertyconfig.QueryPropertyVersionConfigForTerrformResponseDataStagingDeployingVersion) []interface{} {
	if versionConfig == nil {
		return nil
	}
	var resultList []interface{}
	var config = map[string]interface{}{
		"version":   versionConfig.Version,
		"hostnames": versionConfig.Hostnames,
	}
	resultList = append(resultList, config)
	return resultList
}
func buildProductionDeployingVersion(versionConfig *propertyconfig.QueryPropertyVersionConfigForTerrformResponseDataProductionDeployingVersion) []interface{} {
	if versionConfig == nil {
		return nil
	}
	var resultList []interface{}
	var config = map[string]interface{}{
		"version":   versionConfig.Version,
		"hostnames": versionConfig.Hostnames,
	}
	resultList = append(resultList, config)
	return resultList
}

func buildHostnames(hostnames []*propertyconfig.QueryPropertyVersionConfigForTerrformResponseDataHostnames) []interface{} {
	if hostnames == nil || len(hostnames) == 0 {
		return []interface{}{}
	}
	var resultList []interface{}
	for _, hostname := range hostnames {
		result := map[string]interface{}{
			"hostname":       hostname.Hostname,
			"icp":            hostname.Icp,
			"edge_hostname":  buildEdgeHostnameConfig(hostname.EdgeHostname),
			"default_origin": buildDefaultOrigin(hostname.DefaultOrigin),
			"certificates":   buildCertificates(hostname.Certificates),
		}
		resultList = append(resultList, result)
	}
	return resultList
}

func buildEdgeHostnameConfig(edgeHostnameConfig *propertyconfig.QueryPropertyVersionConfigForTerrformResponseDataHostnamesEdgeHostname) []interface{} {
	if edgeHostnameConfig == nil {
		return nil
	}
	return []interface{}{map[string]interface{}{
		"edge_hostname_prefix": edgeHostnameConfig.EdgeHostnamePrefix,
		"edge_hostname_suffix": edgeHostnameConfig.EdgeHostnameSuffix,
		"edge_hostname":        edgeHostnameConfig.EdgeHostname,
		"dns_service_status":   edgeHostnameConfig.DnsServiceStatus,
	}}
}

func buildDefaultOrigin(defaultOrigin *propertyconfig.QueryPropertyVersionConfigForTerrformResponseDataHostnamesDefaultOrigin) []interface{} {
	if defaultOrigin == nil {
		return nil
	}
	return []interface{}{map[string]interface{}{
		"host":       defaultOrigin.Host,
		"http_port":  defaultOrigin.HttpPort,
		"https_port": defaultOrigin.HttpsPort,
		"ip_version": defaultOrigin.IpVersion,
		"servers":    defaultOrigin.Servers,
	}}
}
func buildCertificates(certificates []*propertyconfig.QueryPropertyVersionConfigForTerrformResponseDataHostnamesCertificates) []interface{} {
	if certificates == nil || len(certificates) == 0 {
		return []interface{}{}
	}
	var resultList []interface{}
	for _, certificate := range certificates {
		result := map[string]interface{}{
			"certificate_id":    certificate.CertificateId,
			"certificate_usage": certificate.CertificateUsage,
		}
		resultList = append(resultList, result)
	}
	return resultList
}
