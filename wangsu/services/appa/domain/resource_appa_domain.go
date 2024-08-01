package appadomain

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wangsuCommon "github.com/wangsu-api/terraform-provider-wangsu/wangsu/common"
	appadomain "github.com/wangsu-api/wangsu-sdk-go/wangsu/appa/domain"
	cdn "github.com/wangsu-api/wangsu-sdk-go/wangsu/cdn/domain"
	"golang.org/x/net/context"
	"log"
	"time"
)

func ResourceAppaDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppaDomainCreate,
		ReadContext:   resourceAppaDomainRead,
		UpdateContext: resourceAppaDomainUpdate,
		DeleteContext: resourceAppaDomainDelete,

		Schema: map[string]*schema.Schema{
			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Domain name you want to accelerate. A generic domain name is supported, starting with the symbol '.', such as .example.com.",
			},
			"service_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The service type of the accelerated domain name. The value can be: appa: Application Acceleration",
			},
			"origin_config": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Origin configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"level": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The level of the origin, which value can be an integer ranging from 1 to 5.",
						},
						"strategy": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Origin selection strategy supports fast, robin and hash.",
						},
						"origin": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Origin information of a certain level.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"origin_ip": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Origin address, which can be an IP or domain name.",
									},
									"weight": {
										Type:        schema.TypeInt,
										Optional:    true,
										Default:     10,
										Description: "Weight, which is only useful for robin strategy. If this parameter is not specified, the default value is 10.",
									},
								},
							},
						},
					},
				},
			},
			"http_ports": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "HTTP port. Multiple ports are supported.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"https_ports": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "HTTPS port. Multiple ports are supported.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceAppaDomainCreate(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_appa_domain.create")

	var diags diag.Diagnostics
	request := &appadomain.AddAppaDomainForTerraformRequest{}
	if domainName, ok := data.Get("domain_name").(string); ok && domainName != "" {
		request.DomainName = &domainName
	}
	if serviceType, ok := data.Get("service_type").(string); ok && serviceType != "" {
		request.ServiceType = &serviceType
	}
	if originConfig, ok := data.Get("origin_config").([]interface{}); ok && len(originConfig) > 0 {
		request.OriginConfig = make([]*appadomain.AddAppaDomainForTerraformRequestOriginConfig, 0, len(originConfig))
		for _, v := range originConfig {
			var configTemp = &appadomain.AddAppaDomainForTerraformRequestOriginConfig{}
			originConfigMap := v.(map[string]interface{})

			if level, ok := originConfigMap["level"].(int); ok {
				var levelTemp = int32(level)
				configTemp.Level = &levelTemp
			}
			if strategy, ok := originConfigMap["strategy"].(string); ok {
				configTemp.Strategy = &strategy
			}
			buildOriginList(originConfigMap, configTemp)
			request.OriginConfig = append(request.OriginConfig, configTemp)
		}
	}

	if httpPorts, ok := data.Get("http_ports").([]interface{}); ok && len(httpPorts) > 0 {
		HttpPortList := make([]*string, 0, len(httpPorts))
		for _, v := range httpPorts {
			if v == nil {
				diags = append(diags, diag.FromErr(errors.New("The http port could not be empty."))...)
				return diags
			}
			port := v.(string)
			HttpPortList = append(HttpPortList, &port)
		}
		request.HttpPorts = HttpPortList
	}
	if httpsPorts, ok := data.Get("https_ports").([]interface{}); ok && len(httpsPorts) > 0 {
		HttpsPortList := make([]*string, 0, len(httpsPorts))
		for _, v := range httpsPorts {
			if v == nil {
				diags = append(diags, diag.FromErr(errors.New("The https port could not be empty."))...)
				return diags
			}
			port := v.(string)
			HttpsPortList = append(HttpsPortList, &port)
		}
		request.HttpsPorts = HttpsPortList
	}

	//start to create a domain in 2 minutes
	var addAppaDomainResponse *appadomain.AddAppaDomainForTerraformResponse
	var requestId string
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		requestId, addAppaDomainResponse, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseAppaDomainClient().AddAppaDomain(request)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}
	if addAppaDomainResponse == nil {
		data.SetId("")
		return nil
	}

	data.SetId(*request.DomainName)

	time.Sleep(3 * time.Second)
	//query domain deployment status
	var response *cdn.QueryDeployResultForTerraformResponse
	err = resource.RetryContext(context, time.Duration(5)*time.Minute, func() *resource.RetryError {
		response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseCdnClient().QueryDomainDeployStatus(requestId)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if response != nil && response.Data != nil && *response.Data.DeployResult != "SUCCESS" {
			return resource.RetryableError(fmt.Errorf("domain deployment status is in progress, retrying"))
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

	log.Printf("resource.wangsu_appa_domain.create success")
	_ = data.Set("xCncRequestId", response.Data.RequestId)
	//set status
	return resourceAppaDomainRead(context, data, meta)
}

func buildOriginList(originConfigMap map[string]interface{}, configTemp *appadomain.AddAppaDomainForTerraformRequestOriginConfig) {
	originList := originConfigMap["origin"].([]interface{})
	if len(originList) > 0 {
		configTemp.Origin = make([]*appadomain.AddAppaDomainForTerraformRequestOriginConfigOrigin, 0, len(originList))
		for _, item := range originList {
			originMap := item.(map[string]interface{})
			originIp := originMap["origin_ip"].(string)
			weight := originMap["weight"].(int)
			weightTemp := int32(weight)
			origin := appadomain.AddAppaDomainForTerraformRequestOriginConfigOrigin{
				OriginIp: &originIp,
				Weight:   &weightTemp,
			}
			configTemp.Origin = append(configTemp.Origin, &origin)
		}

	}
}

func resourceAppaDomainRead(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_appa_domain.read")

	var diags diag.Diagnostics
	var response *appadomain.QueryAppaDomainForTerraformResponse
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseAppaDomainClient().QueryAppaDomain(data.Id())
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
	_ = data.Set("domain_id", response.Data.DomainId)
	_ = data.Set("domain_name", response.Data.DomainName)
	_ = data.Set("cname", response.Data.Cname)
	_ = data.Set("service_type", response.Data.ServiceType)
	_ = data.Set("origin_config", flattenOriginConfig(response.Data.OriginConfig))
	_ = data.Set("http_ports", response.Data.HttpPorts)
	_ = data.Set("https_ports", response.Data.HttpsPorts)

	log.Printf("resource.wangsu_appa_domain.read success")
	return nil
}

func flattenOriginConfig(config []*appadomain.QueryAppaDomainForTerraformResponseDataOriginConfig) interface{} {
	var result = make([]interface{}, 0)
	for _, v := range config {
		var item = make(map[string]interface{})
		item["level"] = v.Level
		item["strategy"] = v.Strategy
		var originList = make([]interface{}, 0)
		for _, origin := range v.Origin {
			var originItem = make(map[string]interface{})
			originItem["origin_ip"] = origin.OriginIp
			originItem["weight"] = origin.Weight
			originList = append(originList, originItem)
		}
		item["origin"] = originList
		result = append(result, item)
	}
	return result
}

func resourceAppaDomainUpdate(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_appa_domain.update")
	domainName := data.Id()
	var diags diag.Diagnostics
	request := &appadomain.UpdateAppaDomainForTerraformRequest{}
	if data.HasChanges("origin_config") {
		if originConfig, ok := data.Get("origin_config").([]interface{}); ok && len(originConfig) > 0 {
			request.OriginConfig = make([]*appadomain.UpdateAppaDomainForTerraformRequestOriginConfig, 0, len(originConfig))
			for _, v := range originConfig {
				var configTemp = &appadomain.UpdateAppaDomainForTerraformRequestOriginConfig{}
				originConfigMap := v.(map[string]interface{})

				if level, ok := originConfigMap["level"].(int); ok {
					var levelTemp = int32(level)
					configTemp.Level = &levelTemp
				}
				if strategy, ok := originConfigMap["strategy"].(string); ok {
					configTemp.Strategy = &strategy
				}
				originList := originConfigMap["origin"].([]interface{})
				if len(originList) > 0 {
					configTemp.Origin = make([]*appadomain.UpdateAppaDomainForTerraformRequestOriginConfigOrigin, 0, len(originList))
					for _, item := range originList {
						originMap := item.(map[string]interface{})
						originIp := originMap["origin_ip"].(string)
						weight := originMap["weight"].(int)
						weightTemp := int32(weight)
						origin := appadomain.UpdateAppaDomainForTerraformRequestOriginConfigOrigin{
							OriginIp: &originIp,
							Weight:   &weightTemp,
						}
						configTemp.Origin = append(configTemp.Origin, &origin)
					}
				}
				request.OriginConfig = append(request.OriginConfig, configTemp)
			}
		} else {
			request.OriginConfig = make([]*appadomain.UpdateAppaDomainForTerraformRequestOriginConfig, 0)
		}
	}

	if data.HasChanges("http_ports") {
		if httpPorts, ok := data.Get("http_ports").([]interface{}); ok && len(httpPorts) > 0 {
			HttpPortList := make([]*string, 0, len(httpPorts))
			for _, v := range httpPorts {
				if v == nil {
					diags = append(diags, diag.FromErr(errors.New("The http port could not be empty."))...)
					return diags
				}
				port := v.(string)
				HttpPortList = append(HttpPortList, &port)
			}
			request.HttpPorts = HttpPortList
		} else {
			request.HttpPorts = make([]*string, 0)
		}
	}

	if data.HasChanges("https_ports") {
		if httpsPorts, ok := data.Get("https_ports").([]interface{}); ok && len(httpsPorts) > 0 {
			HttpsPortList := make([]*string, 0, len(httpsPorts))
			for _, v := range httpsPorts {
				if v == nil {
					diags = append(diags, diag.FromErr(errors.New("The https port could not be empty."))...)
					return diags
				}
				port := v.(string)
				HttpsPortList = append(HttpsPortList, &port)
			}
			request.HttpsPorts = HttpsPortList
		} else {
			request.HttpsPorts = make([]*string, 0)
		}
	}

	var updateAppaDomainResponse *appadomain.UpdateAppaDomainForTerraformResponse
	var requestId string
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		requestId, updateAppaDomainResponse, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseAppaDomainClient().UpdateAppaDomain(request, domainName)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}
	if updateAppaDomainResponse == nil {

		return nil
	}

	time.Sleep(3 * time.Second)
	//query domain deployment status
	var response *cdn.QueryDeployResultForTerraformResponse
	err = resource.RetryContext(context, time.Duration(5)*time.Minute, func() *resource.RetryError {
		response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseCdnClient().QueryDomainDeployStatus(requestId)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if response != nil && response.Data != nil && *response.Data.DeployResult != "SUCCESS" {
			return resource.RetryableError(fmt.Errorf("domain deployment status is in progress, retrying"))
		}
		return nil
	})

	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	log.Printf("resource.wangsu_appa_domain.update success")
	return resourceAppaDomainRead(context, data, meta)
}

func resourceAppaDomainDelete(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_appa_domain.delete")

	var response *cdn.DeleteDomainForTerraformResponse
	var requestId string
	var err error
	var diags diag.Diagnostics
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		requestId, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseCdnClient().DeleteCdnDomain(data.Id())
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

	time.Sleep(3 * time.Second)
	//query domain deployment status
	var deploymentResponse *cdn.QueryDeployResultForTerraformResponse
	err = resource.RetryContext(context, time.Duration(5)*time.Minute, func() *resource.RetryError {
		deploymentResponse, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseCdnClient().QueryDomainDeployStatus(requestId)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if deploymentResponse != nil && deploymentResponse.Data != nil && *deploymentResponse.Data.DeployResult != "SUCCESS" {
			return resource.RetryableError(fmt.Errorf("domain deployment status is in progress, retrying"))
		}
		return nil
	})

	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}
	log.Printf("resource.wangsu_appa_domain.delete success")
	return nil
}
