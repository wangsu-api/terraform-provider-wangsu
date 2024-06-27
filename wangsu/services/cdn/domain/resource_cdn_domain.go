package domain

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wangsuCommon "github.com/wangsu/terraform-provider-wangsu/wangsu/common"
	cdn "github.com/wangsu/wangsu-sdk-go/wangsu/cdn/domain"
	"log"
	"time"
)

func ResourceCdnDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCdnDomainCreate,
		ReadContext:   resourceCdnDomainRead,
		UpdateContext: resourceCdnDomainUpdate,
		DeleteContext: resourceCdnDomainDelete,

		Schema: map[string]*schema.Schema{
			"version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Version code, the current version is 1.0.0",
			},
			"domain_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Acceleration domain ID returned by the system.",
			},
			"domain_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Accelerated domain name.",
			},
			"created_date": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Domain name creation time,\nFormat: week, dd month yyyy hh:mm:ss GMT +8:00",
			},
			"last_modified": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Domain name last modified time,\nFormat: week, dd month yyyy hh:mm:ss GMT +8:00\nExample: Mon, 18 Feb 2019 02:54:19 GMT +8 p.m",
			},
			"service_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Speed up domain name service types, including the following:\nWeb/web-https: web acceleration / web acceleration - https\nWsa/wsa-https: Total Station Acceleration / Total Station Acceleration - https\nVodstream/vod-https: on-demand acceleration/on-demand acceleration-https\nDownload/dl-https: Download acceleration/download acceleration - https",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Remarks, up to 1000 characters",
			},
			"service_areas": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Service areas of the domain name",
			},
			"cname": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "CNAME of the domain you queried, for example: 7nt6mrh7sdkslj.cdn30.com.",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The deployment status of the domain name. Deployed indicates that the domain configuration is distributed successfully. InProgress indicates that the deployment task of the domain configuration is still in progress, and may be in queue, or failed.",
			},
			"cdn_service_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Accelerate status of the domain in our CDN, true means the CDN acceleration is normal; false means all request will back to origin directly in DNS.",
			},
			"enabled": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Activation of the domain. It is false when the domain service is disabled, and true when the domain service is enabled.",
			},
			"cname_label": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Shared first level alias",
			},
			"origin_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Return source policy setting, used to set the source station information and return source policy of the accelerated domain name",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"origin_ips": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Return source address, which can be IP or domain name.",
						},
						"default_origin_host_header": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Back to the source HOST, used to change the HOST field in the source HTTP request header.",
						},
						"adv_origin_configs": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Advanced origin config",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"detect_url": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The advanced source monitors the url, and requests <master-ips> through the url. If the response is not 2**, 3** response, it is considered that the primary source ip is faulty, and <backup-ips> is used at this time.",
									},
									"detect_period": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Advanced source monitoring period, in seconds, optional as an integer greater than or equal to 0, 0 means no monitoring",
									},
									"adv_origin_config": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Advanced origin config",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"master_ips": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The advanced source mainly returns the source IP. Multiple IPs are separated by a semicolon \";\", and the return source IP cannot be repeated.",
												},
												"backup_ips": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Advanced source backup source IP, multiple IPs are separated by semicolon \";\", and the return source IP cannot be duplicated.",
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
			"ssl": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "SSL certificate settings, used to set the SSL certificate configuration for the accelerated domain name",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"use_ssl": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Use a certificate, the optional values are true and false, true means to use the certificate, false means not to use the certificate",
						},
						"use_for_sni": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Use SNI certificate, the optional values are true and false, true means use SNI certificate, false means use non-SNI traditional certificate",
						},
						"ssl_certificate_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Certificate ID, the certificate ID returned by the system after the new certificate is successfully added.",
						},
					},
				},
			},
			"cache_behaviors": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Cache rule settings for setting cache rules for accelerated domain names",
			},
			"cache_host": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Cache file HOST (not return by default, application is required to use)",
			},
			"enable_httpdns": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Enable HTTPDNS settings (not return by default, application is required to use)",
			},
			"header_of_clientip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Pass the response header of client IP. The optional values are Cdn-Src-Ip and X-Forwarded-For. The default value is Cdn-Src-Ip",
			},
			"domain_stream_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The live push-pull stream type, the optional values are pull and push, pull means pull flow; push means push flow.",
			},
			"live_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Live domain name configuration, RTMP live acceleration domain name push-pull flow",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"stream_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The live push-pull stream type, the optional values are pull and push.",
						},
						"origin_ips": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Source station IP. When the stream-type is pull, at least one of the source station IP and the companion push stream domain name is not empty.",
						},
						"origin_push_host": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The corresponding push domain name, the RTMP live streaming domain name corresponding to the push domain name",
						},
					},
				},
			},
			"publish_points": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Set the publishing point of the live push-pull domain name\nnote:\n1. Pull flow and corresponding push flow domain name must be configured with the same publishing point.\n2. do not want to modify the publishing point, do not pass the node and the following parameters\n3. The publishing point adopts the overlay update. Each time you modify, you need to submit all the publishing points. You cannot submit only the parts that need to be modified.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uri": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Livestream domain settings. Publish point, support multiple, do not pass the system by default to generate a publishing point uri for [/]",
						},
					},
				},
			},
			// 新增接口独有参数
			"config_form_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Configuration template, if you want to add a domain using some specified configuration by default, you can specify the template id. For more detail, please contact the technical support.",
			},
			"referenced_domain_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Refer to the configuration of the specified domain. Note: 1. If the referenced domain uses a certificate, the new domain should be in the 'DNS name' of the certificate. 2. If the referenced domain has no China ICP, while the new domain name has, it may affect the cover resources and service quality. 3. If the referenced domain has China ICP, while the new domain name doesn't, then the cover resources may be re-selected if it does not meet the policy requirements. 4. It is not allowed to reference a domain which is traffic-free.",
			},
			"cname_with_customized_prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The first level of cname prefix, true, indicates that the domain cname is used as the cname prefix, otherwise the 14-bit random string (number + letter) is used as the cname prefix. Note: When the prefix is a generic domain name, a wsall is added as a prefix. Such as .baidu.com.wscloudcdn.com, which will generate wsall.baidu.com.wscloudcdn.com",
			},
			"accelerate_no_china": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Identifies whether a domain name is fully overseas accelerated. Whether the default is false True: indicates that the client domain name is a pure overseas acceleration False: Indicates that the client domain name has accelerated in China",
			},
			"upstream_host": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The live streaming domain which is pull domain, and directly returned to the source to verify the configuration. which can be an IP or a domain name. Can be IP or domain name. Ip and domain names can only be one. Multiple input parameters are not supported.",
			},
		},
	}
}

func resourceCdnDomainDelete(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_cdn_domain.delete")

	var response *cdn.DeleteApiDomainServiceResponse
	var requestId string
	var err error
	var diags diag.Diagnostics
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		requestId, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseCdnClient().DeleteCdnDomain(data.Id(), data.Get("domain_id").(int))
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

	//query domain deployment status
	var deploymentResponse *cdn.QueryApiDeployServiceResponse
	err = resource.RetryContext(context, time.Duration(5)*time.Minute, func() *resource.RetryError {
		deploymentResponse, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseCdnClient().QueryDomainDeployStatus(requestId)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if deploymentResponse.AsyncResult != nil && *deploymentResponse.AsyncResult != "SUCCESS" {
			return resource.RetryableError(fmt.Errorf("domain deployment status is in progress, retrying"))
		}
		return nil
	})

	return nil
}

func resourceCdnDomainRead(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_cdn_domain.read")

	//query domain information
	var response *cdn.GetBasicConfigurationOfDomainResponse
	var diags diag.Diagnostics
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseCdnClient().GetCdnDomainStatus(data.Id())
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if response == nil {
		data.SetId("")
		return nil
	}

	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	_ = data.Set("domain_id", response.DomainId)
	_ = data.Set("domain_name", response.DomainName)
	_ = data.Set("created_date", response.CreatedDate)
	_ = data.Set("last_modified", response.LastModified)
	_ = data.Set("service_type", response.ServiceType)
	_ = data.Set("comment", response.Comment)
	_ = data.Set("service_areas", response.ServiceAreas)
	_ = data.Set("cname", response.Cname)
	_ = data.Set("status", response.Status)
	_ = data.Set("cdn_service_status", response.CdnServiceStatus)
	_ = data.Set("enabled", response.Enabled)
	_ = data.Set("cname_label", response.CnameLabel)
	originConfig := make([]interface{}, 0)
	if response.OriginConfig != nil && response.OriginConfig.AdvOriginConfigs != nil && response.OriginConfig.AdvOriginConfigs.AdvOriginConfig != nil {
		originConfig = append(originConfig, map[string]interface{}{
			"origin_ips":                 response.OriginConfig.OriginIps,
			"default_origin_host_header": response.OriginConfig.DefaultOriginHostHeader,
			"adv_origin_configs": []interface{}{
				map[string]interface{}{
					"detect_url":    response.OriginConfig.AdvOriginConfigs.DetectUrl,
					"detect_period": response.OriginConfig.AdvOriginConfigs.DetectPeriod,
					"adv_origin_config": []interface{}{
						map[string]interface{}{
							"master_ips": response.OriginConfig.AdvOriginConfigs.AdvOriginConfig.MasterIps,
							"backup_ips": response.OriginConfig.AdvOriginConfigs.AdvOriginConfig.BackupIps,
						},
					},
				},
			},
		})
	}
	_ = data.Set("origin_config", originConfig)

	ssl := make([]interface{}, 0)
	if response.Ssl != nil {
		ssl = append(ssl, map[string]interface{}{
			"use_ssl":            response.Ssl.UseSsl,
			"use_for_sni":        response.Ssl.UseForSni,
			"ssl_certificate_id": response.Ssl.SslCertificateId,
		})
	}
	_ = data.Set("ssl", ssl)

	_ = data.Set("cache_behaviors", response.CacheBehaviors)
	_ = data.Set("cache_host", response.CacheHost)
	_ = data.Set("enable_httpdns", response.EnableHttpdns)
	_ = data.Set("header_of_clientip", response.HeaderOfClientip)
	_ = data.Set("stream_mode", response.StreamMode)

	liveConfig := make([]interface{}, 0)
	if response.LiveConfig != nil {
		liveConfig = append(liveConfig, map[string]interface{}{
			"stream_type":      response.LiveConfig.StreamType,
			"origin_ips":       response.LiveConfig.OriginIps,
			"origin_push_host": response.LiveConfig.OriginPushHost,
		})
	}
	_ = data.Set("live_config", liveConfig)

	publishPoints := make([]interface{}, 0)
	if response.PublishPoints != nil {
		publishPoints = append(publishPoints, map[string]interface{}{
			"uri": response.PublishPoints.Uri,
		})
	}
	_ = data.Set("publish_points", publishPoints)

	return nil
}

func resourceCdnDomainCreate(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_cdn_domain.create")

	var diags diag.Diagnostics
	request := &cdn.CreateDomainRequest{}
	if version, ok := data.Get("version").(string); ok && version != "" {
		request.Version = &version
	}
	if domainName, ok := data.Get("domain_name").(string); ok && domainName != "" {
		request.DomainName = &domainName
	}
	if serviceType, ok := data.Get("service_type").(string); ok && serviceType != "" {
		request.ServiceType = &serviceType
	}
	if serviceAreas, ok := data.Get("service_areas").(string); ok && serviceAreas != "" {
		request.ServiceAreas = &serviceAreas
	}
	if comment, ok := data.Get("comment").(string); ok && comment != "" {
		request.Comment = &comment
	}
	if configFormId, ok := data.Get("config_form_id").(int); ok && configFormId != 0 {
		request.ConfigFormId = &configFormId
	}
	if referencedDomainName, ok := data.Get("referenced_domain_name").(string); ok && referencedDomainName != "" {
		request.ReferencedDomainName = &referencedDomainName
	}
	if cnameLabel, ok := data.Get("cname_label").(string); ok {
		request.CnameLabel = &cnameLabel
	}
	if cnameWithCustomizedPrefix, ok := data.Get("cname_with_customized_prefix").(string); ok && cnameWithCustomizedPrefix != "" {
		request.CnameWithCustomizedPrefix = &cnameWithCustomizedPrefix
	}
	if originConfig, ok := data.Get("origin_config").([]interface{}); ok && len(originConfig) > 0 {
		for _, v := range originConfig {
			originConfigMap := v.(map[string]interface{})

			if originIps, ok := originConfigMap["origin_ips"].(string); ok {
				if defaultOriginHostHeader, ok := originConfigMap["default_origin_host_header"].(string); ok {
					originConfig := cdn.CreateDomainRequestOriginConfig{
						OriginIps:               &originIps,
						DefaultOriginHostHeader: &defaultOriginHostHeader,
					}
					request.OriginConfig = &originConfig
				}
			}
		}
	}
	if liveConfig, ok := data.Get("live_config").([]interface{}); ok && len(liveConfig) > 0 {
		for _, v := range liveConfig {
			liveConfigMap := v.(map[string]interface{})

			if streamType, ok := liveConfigMap["stream_type"].(string); ok {
				if originPushHost, ok := liveConfigMap["origin_push_host"].(string); ok {
					if liveConfigOriginIps, ok := liveConfigMap["live_config_origin_ips"].(string); ok {
						liveConfig := cdn.CreateDomainRequestLiveConfig{
							StreamType:          &streamType,
							OriginPushHost:      &originPushHost,
							LiveConfigOriginIps: &liveConfigOriginIps,
						}
						request.LiveConfig = &liveConfig
					}
				}
			}
		}
	}
	if accelerateNoChina, ok := data.Get("accelerate_no_china").(string); ok && accelerateNoChina != "" {
		request.AccelerateNoChina = &accelerateNoChina
	}
	if headerOfClientip, ok := data.Get("header_of_clientip").(string); ok && headerOfClientip != "" {
		request.HeaderOfClientip = &headerOfClientip
	}
	if upstreamHost, ok := data.Get("upstream_host").(string); ok && upstreamHost != "" {
		request.UpstreamHost = &upstreamHost
	}
	if publishPoints, ok := data.Get("publish_points").([]interface{}); ok && len(publishPoints) > 0 {
		for _, v := range publishPoints {
			publishPointMap := v.(map[string]interface{})

			if uri, ok := publishPointMap["uri"].(string); ok {
				publishPoint := &cdn.CreateDomainRequestPublishPoints{
					Uri: &uri,
				}
				request.PublishPoints = append(request.PublishPoints, publishPoint)
			}
		}
	}
	if ssl, ok := data.Get("ssl").([]interface{}); ok && len(ssl) > 0 {
		for _, v := range ssl {
			sslMap := v.(map[string]interface{})

			if useSsl, ok := sslMap["use_ssl"].(string); ok {
				if useForSni, ok := sslMap["use_for_sni"].(string); ok {
					if sslCertificateId, ok := sslMap["ssl_certificate_id"].(int); ok {
						ssl := cdn.CreateDomainRequestSsl{
							UseSsl:           &useSsl,
							UseForSni:        &useForSni,
							SslCertificateId: &sslCertificateId,
						}
						request.Ssl = &ssl
					}
				}
			}
		}
	}

	//start to create a domain in 2 minutes
	var createDomainResponse *cdn.CreateDomainResponse
	var requestId string
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		requestId, createDomainResponse, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseCdnClient().AddCdnDomain(request)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}
	if createDomainResponse == nil {
		data.SetId("")
		return nil
	}

	data.SetId(*request.DomainName)

	//query domain deployment status
	var response *cdn.QueryApiDeployServiceResponse
	err = resource.RetryContext(context, time.Duration(5)*time.Minute, func() *resource.RetryError {
		response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseCdnClient().QueryDomainDeployStatus(requestId)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if response.AsyncResult != nil && *response.AsyncResult != "SUCCESS" {
			return resource.RetryableError(fmt.Errorf("domain deployment status is in progress, retrying"))
		}
		return nil
	})

	if response == nil {
		data.SetId("")
		return nil
	}

	log.Printf("resource.wangsu_cdn_domain.create success")
	_ = data.Set("xCncRequestId", response.CncRequestId)
	//set status
	return resourceCdnDomainRead(context, data, meta)
}

func resourceCdnDomainUpdate(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_cdn_domain.update")

	request := &cdn.EditDomainConfigRequest{}

	if v, ok := data.GetOk("version"); ok && v.(string) != "" {
		version := v.(string)
		request.Version = &version
	}

	if v, ok := data.GetOk("comment"); ok && v.(string) != "" {
		comment := v.(string)
		request.Comment = &comment
	}

	if v, ok := data.GetOk("service_areas"); ok && v.(string) != "" {
		serviceAreas := v.(string)
		request.ServiceAreas = &serviceAreas
	}

	if v, ok := data.GetOk("cname_label"); ok && v.(string) != "" {
		cnameLabel := v.(string)
		request.CnameLabel = &cnameLabel
	}

	if v, ok := data.GetOk("origin_config"); ok && len(v.([]interface{})) > 0 {
		originConfig := v.([]interface{})
		for _, oc := range originConfig {
			ocMap := oc.(map[string]interface{})
			originIps := ocMap["origin_ips"].(string)
			defaultOriginHostHeader := ocMap["default_origin_host_header"].(string)
			originConfig := cdn.EditDomainConfigRequestOriginConfig{
				OriginIps:               &originIps,
				DefaultOriginHostHeader: &defaultOriginHostHeader,
			}
			request.OriginConfig = &originConfig
		}
	}

	if v, ok := data.GetOk("ssl"); ok && len(v.([]interface{})) > 0 {
		sslConfig := v.(map[string]interface{})
		useSsl := sslConfig["use_ssl"].(string)
		useForSni := sslConfig["use_for_sni"].(string)
		sslCertificateId := sslConfig["ssl_certificate_id"].(int)
		ssl := cdn.EditDomainConfigRequestSsl{
			UseSsl:           &useSsl,
			UseForSni:        &useForSni,
			SslCertificateId: &sslCertificateId,
		}
		request.Ssl = &ssl
	}

	if v, ok := data.GetOk("cache_host"); ok && v.(string) != "" {
		cacheHost := v.(string)
		request.CacheHost = &cacheHost
	}

	if v, ok := data.GetOk("enable_httpdns"); ok && v.(string) != "" {
		enableHttpdns := v.(string)
		request.EnableHttpdns = &enableHttpdns
	}

	if v, ok := data.GetOk("header_of_clientip"); ok && v.(string) != "" {
		headerOfClientip := v.(string)
		request.HeaderOfClientip = &headerOfClientip
	}

	if v, ok := data.GetOk("live_config"); ok && len(v.([]interface{})) > 0 {
		liveConfig := v.(map[string]interface{})
		liveConfigOriginIps := liveConfig["live_config_origin_ips"].(string)
		originPushHost := liveConfig["origin_push_host"].(string)
		liveConfigStruct := cdn.EditDomainConfigRequestLiveConfig{
			LiveConfigOriginIps: &liveConfigOriginIps,
			OriginPushHost:      &originPushHost,
		}
		request.LiveConfig = &liveConfigStruct
	}

	if v, ok := data.GetOk("publish_points"); ok && len(v.([]interface{})) > 0 {
		publishPoints := v.([]interface{})
		var publishPointsList []*cdn.EditDomainConfigRequestPublishPoints
		for _, pp := range publishPoints {
			ppMap := pp.(map[string]interface{})
			uri := ppMap["uri"].(string)
			publishPoint := cdn.EditDomainConfigRequestPublishPoints{
				Uri: &uri,
			}
			publishPointsList = append(publishPointsList, &publishPoint)
		}
		request.PublishPoints = publishPointsList
	}

	var editResponse *cdn.EditDomainConfigResponse
	var requestId string
	var diags diag.Diagnostics
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		requestId, editResponse, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseCdnClient().EditDomainConfig(request, data.Id())
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if editResponse == nil {
		data.SetId("")
		return nil
	}

	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	//query domain deployment status
	var response *cdn.QueryApiDeployServiceResponse
	err = resource.RetryContext(context, time.Duration(5)*time.Minute, func() *resource.RetryError {
		response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseCdnClient().QueryDomainDeployStatus(requestId)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if response.AsyncResult != nil && *response.AsyncResult != "SUCCESS" {
			return resource.RetryableError(fmt.Errorf("domain deployment status is in progress, retrying"))
		}
		return nil
	})

	log.Printf("resource.wangsu_cdn_domain.update success")
	return nil
}
