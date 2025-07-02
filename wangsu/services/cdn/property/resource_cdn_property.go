package property

import (
	"bytes"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	wangsuCommon "github.com/wangsu-api/terraform-provider-wangsu/wangsu/common"
	propertyConfig "github.com/wangsu-api/wangsu-sdk-go/wangsu/propertyconfig"
	"golang.org/x/net/context"
	"log"
	"reflect"
	"strconv"
	"time"
)

func ResourceCdnProperty() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCdnPropertyCreate,
		ReadContext:   resourceCdnPropertyRead,
		UpdateContext: resourceCdnPropertyUpdate,
		DeleteContext: resourceCdnPropertyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"service_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Product Service Type related to your contract. Optional values include wsa, wsa-https.",
			},
			"property_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the property. The length must not exceed 256 characters.",
			},
			"property_comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Property description. The length must not exceed 256 characters.",
			},
			"version_comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Property version description. The length must not exceed 256 characters.",
			},
			"hostnames": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "List of hostnames.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hostname": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Hostname, the length must not exceed 128 characters. A wildcard hostname must start with an asterisk (*).",
						},
						"default_origin": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Default origin.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"servers": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "Origin servers.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"host": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Origin host.",
									},
									"http_port": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "HTTP port.",
									},
									"https_port": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "HTTPS port.",
									},
									"ip_version": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "IP version. Data range: dual, ipv4, ipv6.",
									},
								},
							},
						},
						"certificates": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Hostname association SSL configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"certificate_id": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Certificate ID.",
									},
									"certificate_usage": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Certificate usage. Data range: default_sni, ssl_bk, gm_sm2_enc, gm_sm2_sign, client_mtls.",
									},
								},
							},
						},
						"edge_hostname": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Hostname edge-hostname config.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"edge_hostname_prefix": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Edge-hostname prefix.",
									},
									"edge_hostname_suffix": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Edge-hostname suffix.",
									},
									"comment": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Edge-hostname comment.",
									},
								},
							},
						},
					},
				},
			},
			"origins": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Specify the hostname and settings used to contact the origin once service begins.",
				DiffSuppressFunc: suppressEquivalentJsonDiffs,
				ValidateFunc:     validation.StringIsJSON,
			},
			"variables": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "The Variables feature allows you to define variables, assign values to them, and reuse them in functional test cases. A variable consists of a name and a value.",
				DiffSuppressFunc: suppressEquivalentJsonDiffs,
				ValidateFunc:     validation.StringIsJSON,
			},
			"rules": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Rules.",
				DiffSuppressFunc: suppressEquivalentJsonDiffs,
				ValidateFunc:     validation.StringIsJSON,
			},
			//computed
			"property_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Property ID, which is the unique identifier for the property.",
			},
			"current_version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Current property version, which is the version number of the property.",
			},
			"latest_version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Latest version of the property.",
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
		},
	}
}

func resourceCdnPropertyRead(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_cdn_property.read")

	var diags diag.Diagnostics
	propertyId, err := strconv.Atoi(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	var response *propertyConfig.QueryPropertyConfigForTerrformResponse
	err = resource.RetryContext(context, time.Duration(1)*time.Minute, func() *resource.RetryError {
		response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UsePropertyConfigClient().QueryProperty(propertyId)
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

	if err := data.Set("service_type", response.Data.ServiceType); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("property_name", response.Data.PropertyName); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("property_comment", response.Data.PropertyComment); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("version_comment", response.Data.VersionComment); err != nil {
		return diag.FromErr(err)
	}
	hostnameList := make([]interface{}, 0, len(response.Data.Hostnames))
	for _, hostname := range response.Data.Hostnames {
		hostnamesMap := map[string]interface{}{
			"hostname":       hostname.Hostname,
			"default_origin": buildDefaultOriginForRead(hostname.DefaultOrigin),
			"certificates":   buildCertificatesForRead(hostname.Certificates),
			"edge_hostname":  buildEdgeHostnameForRead(hostname.EdgeHostname),
		}
		hostnameList = append(hostnameList, hostnamesMap)
	}
	if err := data.Set("hostnames", hostnameList); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("origins", response.Data.Origins); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("variables", response.Data.Variables); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("rules", response.Data.Rules); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("property_id", response.Data.PropertyId); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("current_version", response.Data.CurrentVersion); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("latest_version", response.Data.LatestVersion); err != nil {
		return diag.FromErr(err)
	}
	//response.Data.StagingVersion不为空时才设置
	if response.Data.StagingVersion != nil {
		if err := data.Set("staging_version", []interface{}{map[string]interface{}{
			"version":   response.Data.StagingVersion.Version,
			"hostnames": response.Data.StagingVersion.Hostnames,
		}}); err != nil {
			return diag.FromErr(err)
		}
	} else {
		if err := data.Set("staging_version", nil); err != nil {
			return diag.FromErr(err)
		}
	}
	if response.Data.ProductionVersion != nil {
		if err := data.Set("production_version", []interface{}{map[string]interface{}{
			"version":   response.Data.ProductionVersion.Version,
			"hostnames": response.Data.ProductionVersion.Hostnames,
		}}); err != nil {
			return diag.FromErr(err)
		}
	} else {
		if err := data.Set("production_version", nil); err != nil {
			return diag.FromErr(err)
		}
	}
	log.Printf("resource.wangsu_cdn_property.read success")
	return nil
}

func buildDefaultOriginForRead(defaultOrigin *propertyConfig.QueryPropertyConfigForTerrformResponseDataHostnamesDefaultOrigin) []interface{} {
	if defaultOrigin == nil {
		return nil
	}
	return []interface{}{map[string]interface{}{
		"servers":    defaultOrigin.Servers,
		"host":       defaultOrigin.Host,
		"http_port":  defaultOrigin.HttpPort,
		"https_port": defaultOrigin.HttpsPort,
		"ip_version": defaultOrigin.IpVersion,
	}}
}

func buildCertificatesForRead(certificates []*propertyConfig.QueryPropertyConfigForTerrformResponseDataHostnamesCertificates) []interface{} {
	if certificates == nil || len(certificates) == 0 {
		return nil
	}
	certificatesList := make([]interface{}, 0, len(certificates))
	for _, cert := range certificates {
		certMap := map[string]interface{}{
			"certificate_id":    cert.CertificateId,
			"certificate_usage": cert.CertificateUsage,
		}
		certificatesList = append(certificatesList, certMap)
	}
	return certificatesList
}

func buildEdgeHostnameForRead(edgeHostname *propertyConfig.QueryPropertyConfigForTerrformResponseDataHostnamesEdgeHostname) []interface{} {
	if edgeHostname == nil {
		return nil
	}
	return []interface{}{map[string]interface{}{
		"edge_hostname_prefix": edgeHostname.EdgeHostnamePrefix,
		"edge_hostname_suffix": edgeHostname.EdgeHostnameSuffix,
	}}
}

func resourceCdnPropertyCreate(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_cdn_property.create")

	var diags diag.Diagnostics
	request := &propertyConfig.CreatePropertyForTerraformRequest{}
	if serviceType, ok := data.Get("service_type").(string); ok && serviceType != "" {
		request.ServiceType = &serviceType
	}
	if propertyName, ok := data.Get("property_name").(string); ok && propertyName != "" {
		request.PropertyName = &propertyName
	}
	if propertyComment, ok := data.Get("property_comment").(string); ok && propertyComment != "" {
		request.PropertyComment = &propertyComment
	}
	if versionComment, ok := data.Get("version_comment").(string); ok && versionComment != "" {
		request.VersionComment = &versionComment
	}
	if hostnames, ok := data.Get("hostnames").([]interface{}); ok && len(hostnames) > 0 {
		hostnamesList := make([]*propertyConfig.CreatePropertyForTerraformRequestHostnames, 0, len(hostnames))
		for _, v := range hostnames {
			if hostnameMap, ok := v.(map[string]interface{}); ok {
				hostname := &propertyConfig.CreatePropertyForTerraformRequestHostnames{}
				if hostnameStr, ok := hostnameMap["hostname"].(string); ok && hostnameStr != "" {
					hostname.Hostname = &hostnameStr
				}
				if defaultOrigin, ok := hostnameMap["default_origin"].([]interface{}); ok && len(defaultOrigin) > 0 {
					defaultOriginMap := defaultOrigin[0].(map[string]interface{})

					defaultOriginConfig := &propertyConfig.CreatePropertyForTerraformRequestHostnamesDefaultOrigin{}
					if servers, ok := defaultOriginMap["servers"].([]interface{}); ok && len(servers) > 0 {
						originServers := make([]*string, 0)
						for _, server := range servers {
							if serverStr, ok := server.(string); ok && serverStr != "" {
								originServers = append(originServers, &serverStr)
							}
						}
						defaultOriginConfig.Servers = originServers
					}
					if host, ok := defaultOriginMap["host"].(string); ok && host != "" {
						defaultOriginConfig.Host = &host
					}
					if httpPort, ok := defaultOriginMap["http_port"].(int); ok && httpPort > 0 {
						defaultOriginConfig.HttpPort = &httpPort
					}
					if httpsPort, ok := defaultOriginMap["https_port"].(int); ok && httpsPort > 0 {
						defaultOriginConfig.HttpsPort = &httpsPort
					}
					if ipVersion, ok := defaultOriginMap["ip_version"].(string); ok && ipVersion != "" {
						defaultOriginConfig.IpVersion = &ipVersion
					}
					hostname.DefaultOrigin = defaultOriginConfig
				}
				if certificates, ok := hostnameMap["certificates"].([]interface{}); ok && len(certificates) > 0 {
					certificatesList := make([]*propertyConfig.CreatePropertyForTerraformRequestHostnamesCertificates, 0)
					for _, cert := range certificates {
						certMap := cert.(map[string]interface{})
						certificateId := certMap["certificate_id"].(int)
						certificateUsage := certMap["certificate_usage"].(string)
						certificatesList = append(certificatesList, &propertyConfig.CreatePropertyForTerraformRequestHostnamesCertificates{
							CertificateId:    &certificateId,
							CertificateUsage: &certificateUsage,
						})
					}
					hostname.Certificates = certificatesList
				}
				if edgeHostnames, ok := hostnameMap["edge_hostname"].([]interface{}); ok && len(edgeHostnames) > 0 {
					for _, edgeHost := range edgeHostnames {
						edgeHostMap := edgeHost.(map[string]interface{})
						edgeHostnamePrefix := edgeHostMap["edge_hostname_prefix"].(string)
						edgeHostnameSuffix := edgeHostMap["edge_hostname_suffix"].(string)
						comment := edgeHostMap["comment"].(string)
						hostname.EdgeHostname = &propertyConfig.CreatePropertyForTerraformRequestHostnamesEdgeHostname{
							EdgeHostnamePrefix: &edgeHostnamePrefix,
							EdgeHostnameSuffix: &edgeHostnameSuffix,
							Comment:            &comment,
						}
					}
				}
				hostnamesList = append(hostnamesList, hostname)
			}
		}
		request.Hostnames = hostnamesList
	}
	if origins, ok := data.Get("origins").(string); ok && origins != "" {
		request.Origins = &origins
	}
	if variables, ok := data.Get("variables").(string); ok && variables != "" {
		request.Variables = &variables
	}
	if rules, ok := data.Get("rules").(string); ok && rules != "" {
		request.Rules = &rules
	}

	var response *propertyConfig.CreatePropertyForTerraformResponse
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UsePropertyConfigClient().CreateProperty(request)
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
	var propertyIdStr = strconv.FormatInt(*response.Data.PropertyId, 10)
	data.SetId(propertyIdStr)
	log.Printf("resource.wangsu_cdn_property.create success, propertyId: %s", propertyIdStr)
	time.Sleep(1 * time.Second)
	return resourceCdnPropertyRead(context, data, meta)
}

func resourceCdnPropertyUpdate(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_cdn_property.update")

	var diags diag.Diagnostics
	propertyId, err := strconv.Atoi(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	request := &propertyConfig.UpdatePropertyForTerraformRequest{}
	if propertyName, ok := data.Get("property_name").(string); ok && propertyName != "" {
		request.PropertyName = &propertyName
	}
	if propertyComment, ok := data.Get("property_comment").(string); ok && propertyComment != "" {
		request.PropertyComment = &propertyComment
	}
	if versionComment, ok := data.Get("version_comment").(string); ok && versionComment != "" {
		request.VersionComment = &versionComment
	}
	if hostnames, ok := data.Get("hostnames").([]interface{}); ok && len(hostnames) > 0 {
		hostnamesList := make([]*propertyConfig.UpdatePropertyForTerraformRequestHostnames, 0, len(hostnames))
		for _, v := range hostnames {
			if hostnameMap, ok := v.(map[string]interface{}); ok {
				hostname := &propertyConfig.UpdatePropertyForTerraformRequestHostnames{}
				if hostnameStr, ok := hostnameMap["hostname"].(string); ok && hostnameStr != "" {
					hostname.Hostname = &hostnameStr
				}
				if defaultOrigin, ok := hostnameMap["default_origin"].([]interface{}); ok && len(defaultOrigin) > 0 {
					defaultOriginMap := defaultOrigin[0].(map[string]interface{})

					defaultOriginConfig := &propertyConfig.UpdatePropertyForTerraformRequestHostnamesDefaultOrigin{}
					if servers, ok := defaultOriginMap["servers"].([]interface{}); ok && len(servers) > 0 {
						originServers := make([]*string, 0)
						for _, server := range servers {
							if serverStr, ok := server.(string); ok && serverStr != "" {
								originServers = append(originServers, &serverStr)
							}
						}
						defaultOriginConfig.Servers = originServers
					}
					if host, ok := defaultOriginMap["host"].(string); ok && host != "" {
						defaultOriginConfig.Host = &host
					}
					if httpPort, ok := defaultOriginMap["http_port"].(int); ok && httpPort > 0 {
						defaultOriginConfig.HttpPort = &httpPort
					}
					if httpsPort, ok := defaultOriginMap["https_port"].(int); ok && httpsPort > 0 {
						defaultOriginConfig.HttpsPort = &httpsPort
					}
					if ipVersion, ok := defaultOriginMap["ip_version"].(string); ok && ipVersion != "" {
						defaultOriginConfig.IpVersion = &ipVersion
					}
					hostname.DefaultOrigin = defaultOriginConfig
				}
				if certificates, ok := hostnameMap["certificates"].([]interface{}); ok && len(certificates) > 0 {
					certificatesList := make([]*propertyConfig.UpdatePropertyForTerraformRequestHostnamesCertificates, 0)
					for _, cert := range certificates {
						certMap := cert.(map[string]interface{})
						certificateId := certMap["certificate_id"].(int)
						certificateUsage := certMap["certificate_usage"].(string)
						certificatesList = append(certificatesList, &propertyConfig.UpdatePropertyForTerraformRequestHostnamesCertificates{
							CertificateId:    &certificateId,
							CertificateUsage: &certificateUsage,
						})
					}
					hostname.Certificates = certificatesList
				}
				if edgeHostnames, ok := hostnameMap["edge_hostname"].([]interface{}); ok && len(edgeHostnames) > 0 {
					for _, edgeHost := range edgeHostnames {
						edgeHostMap := edgeHost.(map[string]interface{})
						edgeHostnamePrefix := edgeHostMap["edge_hostname_prefix"].(string)
						edgeHostnameSuffix := edgeHostMap["edge_hostname_suffix"].(string)
						comment := edgeHostMap["comment"].(string)
						hostname.EdgeHostname = &propertyConfig.UpdatePropertyForTerraformRequestHostnamesEdgeHostname{
							EdgeHostnamePrefix: &edgeHostnamePrefix,
							EdgeHostnameSuffix: &edgeHostnameSuffix,
							Comment:            &comment,
						}
					}
				}
				hostnamesList = append(hostnamesList, hostname)
			}
		}
		request.Hostnames = hostnamesList
	}
	if origins, ok := data.Get("origins").(string); ok && origins != "" {
		request.Origins = &origins
	}
	if variables, ok := data.Get("variables").(string); ok && variables != "" {
		request.Variables = &variables
	}
	if rules, ok := data.Get("rules").(string); ok && rules != "" {
		request.Rules = &rules
	}
	var response *propertyConfig.UpdatePropertyForTerraformResponse
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UsePropertyConfigClient().UpdateProperty(propertyId, request)
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
	log.Printf("resource.wangsu_cdn_property.update success, propertyId: %s", data.Id())
	time.Sleep(1 * time.Second)
	return resourceCdnPropertyRead(context, data, meta)
}

func resourceCdnPropertyDelete(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_cdn_property.delete")

	var diags diag.Diagnostics
	propertyId, err := strconv.Atoi(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		_, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UsePropertyConfigClient().DeleteProperty(propertyId)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	log.Printf("resource.wangsu_cdn_property.delete success, propertyId: %s", data.Id())
	return nil
}

func suppressEquivalentJsonDiffs(_, old, new string, _ *schema.ResourceData) bool {
	oldBuff := bytes.NewBufferString("")
	if err := json.Compact(oldBuff, []byte(old)); err != nil {
		return false
	}

	newBuff := bytes.NewBufferString("")
	if err := json.Compact(newBuff, []byte(new)); err != nil {
		return false
	}

	return jsonBytesEqual(oldBuff.Bytes(), newBuff.Bytes())
}

func jsonBytesEqual(a, b []byte) bool {
	var j, j2 interface{}
	if err := json.Unmarshal(a, &j); err != nil {
		return false
	}
	if err := json.Unmarshal(b, &j2); err != nil {
		return false
	}
	return reflect.DeepEqual(j2, j)
}
