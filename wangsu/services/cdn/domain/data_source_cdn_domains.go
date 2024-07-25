package domain

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wangsuCommon "github.com/wangsu-api/terraform-provider-wangsu/wangsu/common"
	cdn "github.com/wangsu-api/wangsu-sdk-go/wangsu/cdn/domain"
	"golang.org/x/net/context"
	"log"
	"time"
)

func DataSourceWangSuCdnDomains() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWangSuCdnDomainsRead,
		Schema: map[string]*schema.Schema{
			"page_number": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Page number must be a positive integer greater than 0.If not passed, then no paging. If it is passed, pageSize is required.",
			},
			"page_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Number of domain name data items for paging, must be a positive integer greater than 0.If not passed, then no paging. If it is passed, pageSize is required.",
			},
			"service_types": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specify the service type to be queried. Multiple services are allowed. Data will be returned if any one service is satisfied. If not passed, all services will be checked by default. For example: [wsa,waf], returns all domains whose services include wsa or include waf.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"domain_names": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specify the accelerated domain name for the query. Multiple domain names are allowed. If not specified, all domain names will be searched by default.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"start_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "RFC3339 formatted date indicating the starting date. Example: 2024-01-01T22:30:00+08:00",
			},
			"end_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "RFC3339 formatted date indicating the ending date. Example: 2024-01-01T22:30:00+08:00",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Status of the accelerated domain. Optional value: enabled, disabled, deploying, checking, disabling, deployFailed, disableFailed.",
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
						"page_number": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Page number to query",
						},
						"page_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of items per page",
						},
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total count of matching domains",
						},
						"total_page_number": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total number of pages",
						},
						"result_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Domain list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cname": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cname of the accelerated domain",
									},
									"create_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Create time of the accelerated domain. Example: 2024-01-01T22:30:00+08:00",
									},
									"domain_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Corresponding domain ID",
									},
									"domain_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Accelerated domain name",
									},
									"service_types": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Service type for accelerated domain name",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Status of the accelerated domain. Optional value: enabled, disabled, deploying, checking, disabling.",
									},
									"enabled": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Accelerated domain enabling status: true indicates that it is enabled, false indicates that it is disabled.",
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

func dataSourceWangSuCdnDomainsRead(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("data_source.wangsu_cdn_domains.read")

	request := &cdn.QueryPagingDomainListForTerraformRequest{}
	if pageNumber, ok := data.Get("page_number").(int); ok && pageNumber > 0 {
		request.PageNumber = &pageNumber
	}

	if pageSize, ok := data.Get("page_size").(int); ok && pageSize > 0 {
		request.PageSize = &pageSize
	}

	if domainNameList, ok := data.Get("domain_names").([]interface{}); ok && len(domainNameList) > 0 {
		domainNames := make([]*string, 0, len(domainNameList))
		for _, v := range domainNameList {
			if domainName, ok := v.(string); ok {
				domainNames = append(domainNames, &domainName)
			}
		}
		request.DomainNames = domainNames
	}

	if serviceTypeList, ok := data.Get("service_types").([]interface{}); ok && len(serviceTypeList) > 0 {
		serviceTypes := make([]*string, 0, len(serviceTypeList))
		for _, v := range serviceTypeList {
			if serviceType, ok := v.(string); ok {
				serviceTypes = append(serviceTypes, &serviceType)
			}
		}
		request.ServiceTypes = serviceTypes
	}

	if domainStatus, ok := data.Get("status").(string); ok && domainStatus != "" {
		request.Status = &domainStatus
	}

	if startTime, ok := data.Get("start_time").(string); ok && startTime != "" {
		request.StartTime = &startTime
	}

	if endTime, ok := data.Get("end_time").(string); ok && endTime != "" {
		request.EndTime = &endTime
	}

	var response *cdn.QueryPagingDomainListForTerraformResponse
	var diags diag.Diagnostics
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseCdnClient().QueryCdnDomainList(request)
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

	ids := make([]string, 0)

	_ = data.Set("code", response.Code)
	_ = data.Set("message", response.Message)
	dataMap := map[string]interface{}{
		"page_number":       response.Data.PageNumber,
		"page_size":         response.Data.PageSize,
		"total_count":       response.Data.TotalCount,
		"total_page_number": response.Data.TotalPageNumber,
	}
	if response.Data.ResultList != nil {
		resultList := make([]map[string]interface{}, len(response.Data.ResultList))
		for i, item := range response.Data.ResultList {
			resultList[i] = map[string]interface{}{
				"domain_id":     item.DomainId,
				"domain_name":   item.DomainName,
				"service_types": item.ServiceTypes,
				"cname":         item.Cname,
				"create_time":   item.CreateTime,
				"status":        item.Status,
				"enabled":       item.Enabled,
			}
			ids = append(ids, *item.DomainName)
		}
		dataMap["result_list"] = resultList
	}
	_ = data.Set("data", []interface{}{dataMap})
	data.SetId(wangsuCommon.DataResourceIdsHash(ids))
	log.Printf("data_source.wangsu_cdn_domains.read success")
	return nil
}
