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
	"strings"
	"time"
)

func DataSourceWangsuCdnProperties() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWangSuCdnPropertiesRead,
		Schema: map[string]*schema.Schema{
			"service_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Service type. Optional values include wsa, wsa-https.",
			},
			"target": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The value can be 'staging', or 'production', or 'none' to filter the results based on where the property has been deployed.",
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
				Description:  "Order of properties to return. Enum: asc,desc Default: desc",
			},
			"sort_by": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "lastUpdateTime",
				ValidateFunc: wangsuCommon.ValidateAllowedStringValue([]string{"creationTime", "lastUpdateTime"}),
				Description:  "Returns results in sorted order. Enum: creationTime,lastUpdateTime Default: lastUpdateTime",
			},
			"hostname": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "Filter by hostname. If specified, only properties with this hostname will be returned.",
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
						"properties": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of properties.",
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
										Description: "Service type. Optional values include wsa, wsa-https.",
									},
									"creation_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "RFC3339 format date indicating when the property was created.",
									},
									"last_update_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "RFC3339 date indicating when the property was last updated.",
									},
									"latest_version": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Latest version of the property.",
									},
									"staging_version": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Describes the version of the property deployed to staging.",
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
										Description: "Describes the version of the property deployed to production.",
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
										Description: "Describes the version of the property deploying to staging.",
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
										Description: "Describes the version of the property deploying to production.",
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
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceWangSuCdnPropertiesRead(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("data_source.wangsu_cdn_properties.read")

	parameters := &propertyconfig.QueryPropertiesForTerraformParameters{}
	if serviceType, ok := data.Get("service_type").(string); ok && serviceType != "" {
		parameters.ServiceType = &serviceType
	}
	if target, ok := data.Get("target").([]interface{}); ok && len(target) > 0 {
		targets := make([]string, 0, len(target))
		for _, v := range target {
			if t, ok := v.(string); ok && t != "" {
				targets = append(targets, t)
			}
		}
		targetValue := strings.Join(targets, ",")
		parameters.Target = &targetValue
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
	if hostnameList, ok := data.Get("hostname").([]interface{}); ok && len(hostnameList) > 0 {
		hostnameArray := make([]string, 0, len(hostnameList))
		for _, v := range hostnameList {
			if t, ok := v.(string); ok && t != "" {
				hostnameArray = append(hostnameArray, t)
			}
		}
		hostnameValue := strings.Join(hostnameArray, ",")
		parameters.Hostname = &hostnameValue
	}

	var response *propertyconfig.QueryPropertiesForTerraformResponse
	var diags diag.Diagnostics
	var err error
	var requestId string
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		requestId, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UsePropertyConfigClient().QueryProperties(parameters)
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

	if err := data.Set("code", response.Code); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("message", response.Message); err != nil {
		return diag.FromErr(err)
	}

	responseData := response.Data
	var properties []interface{}
	for _, property := range responseData.Properties {
		propertyData := map[string]interface{}{
			"property_id":                  property.PropertyId,
			"property_name":                property.PropertyName,
			"property_comment":             property.PropertyComment,
			"service_type":                 property.ServiceType,
			"creation_time":                property.CreationTime,
			"last_update_time":             property.LastUpdateTime,
			"latest_version":               property.LatestVersion,
			"staging_version":              buildForStagingVersion(property.StagingVersion),
			"production_version":           buildForProductionVersion(property.ProductionVersion),
			"staging_deploying_version":    buildForStagingDeployingVersion(property.StagingDeployingVersion),
			"production_deploying_version": buildForProductionDeployingVersion(property.ProductionDeployingVersion),
		}
		properties = append(properties, propertyData)
	}

	err = data.Set("data", buildData(responseData.Count, properties))
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(requestId)
	log.Printf("data_source.wangsu_cdn_properties.read success")
	return nil
}

func buildData(count *int, properties []interface{}) []interface{} {
	var resultList []interface{}
	var config = map[string]interface{}{
		"count":      count,
		"properties": properties,
	}
	resultList = append(resultList, config)
	return resultList
}

func buildForStagingVersion(stagingVersion *propertyconfig.QueryPropertiesForTerraformResponseDataPropertiesStagingVersion) []interface{} {
	if stagingVersion == nil {
		return nil
	}
	var resultList []interface{}
	var config = map[string]interface{}{
		"version":   stagingVersion.Version,
		"hostnames": stagingVersion.Hostnames,
	}
	resultList = append(resultList, config)
	return resultList
}
func buildForProductionVersion(productionVersion *propertyconfig.QueryPropertiesForTerraformResponseDataPropertiesProductionVersion) []interface{} {
	if productionVersion == nil {
		return nil
	}
	var resultList []interface{}
	var config = map[string]interface{}{
		"version":   productionVersion.Version,
		"hostnames": productionVersion.Hostnames,
	}
	resultList = append(resultList, config)
	return resultList
}
func buildForStagingDeployingVersion(stagingVersion *propertyconfig.QueryPropertiesForTerraformResponseDataPropertiesStagingDeployingVersion) []interface{} {
	if stagingVersion == nil {
		return nil
	}
	var resultList []interface{}
	var config = map[string]interface{}{
		"version":   stagingVersion.Version,
		"hostnames": stagingVersion.Hostnames,
	}
	resultList = append(resultList, config)
	return resultList
}
func buildForProductionDeployingVersion(productionVersion *propertyconfig.QueryPropertiesForTerraformResponseDataPropertiesProductionDeployingVersion) []interface{} {
	if productionVersion == nil {
		return nil
	}
	var resultList []interface{}
	var config = map[string]interface{}{
		"version":   productionVersion.Version,
		"hostnames": productionVersion.Hostnames,
	}
	resultList = append(resultList, config)
	return resultList
}
