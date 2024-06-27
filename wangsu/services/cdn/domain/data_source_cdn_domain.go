package domain

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wangsuCommon "github.com/wangsu/terraform-provider-wangsu/wangsu/common"
	cdn "github.com/wangsu/wangsu-sdk-go/wangsu/cdn/domain"
	"log"
	"time"
)

func DataSourceWangSuCdnDomains() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWangSuCdnDomainsRead,
		Schema: map[string]*schema.Schema{
			"domain_name": {
				Type:        schema.TypeList,
				Description: "Specifies the accelerated domain name for the query, allows multiple domains, commas delimited, and no default lookup of all domain names",
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"service_type": {
				Type:        schema.TypeString,
				Description: "Specifies the service type of the query, only one type per query, and no default lookup for all types",
				Optional:    true,
				ForceNew:    true,
			},
			"query_type": {
				Type:        schema.TypeString,
				Description: "Query to accelerated domain name, optional value is: fuzzy_match for fuzzy query; Full_match represents an exact query No fuzzy_match by default, for accelerated domain name only",
				Optional:    true,
				ForceNew:    true,
			},
			"start_time": {
				Type:        schema.TypeString,
				Description: "Query start time, support for years, months, days, hours, minutes, and seconds, for example: 20170101.09 million. Time equals",
				Optional:    true,
				ForceNew:    true,
			},
			"end_time": {
				Type:        schema.TypeString,
				Description: "Query end time, query time within the existence of the accelerated domain name, time is equal to, do not pass the default query all",
				Optional:    true,
				ForceNew:    true,
			},
			"domain_status": {
				Type:        schema.TypeString,
				Description: "Accelerate the status of the domain name, enabled indicates that it is in effect; Disabled indicates that it is Disabled; Deploying means in the process of deployment; Checking indicates that the audit is in progress; Disabling: Indicates disabled, no default lookup for all",
				Optional:    true,
				ForceNew:    true,
			},
			"page_number": {
				Type:        schema.TypeInt,
				Description: "Page number must be a positive integer greater than 0",
				Required:    true,
				ForceNew:    true,
			},
			"page_size": {
				Type:        schema.TypeInt,
				Description: "Number of domain name data items for paging, must be a positive integer greater than 0",
				Required:    true,
				ForceNew:    true,
			},
			"total_count": {
				Type:        schema.TypeInt,
				Description: "Responses the page number of the data",
				Required:    true,
				ForceNew:    true,
			},
			"total_page_number": {
				Type:        schema.TypeInt,
				Description: "total pages",
				Required:    true,
				ForceNew:    true,
			},
			"result_list": {
				Type:        schema.TypeList,
				Description: "Responses status information for the accelerated domain name",
				Optional:    true,
				ForceNew:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cname": {
							Type:        schema.TypeString,
							Description: "Accelerated domain CNAME corresponding to CNAME, for example: 7nt6mrh7sdkslj.cdn30.com",
							Computed:    true,
						},
						"config_form_name": {
							Type:        schema.TypeString,
							Description: "Configuration name",
							Computed:    true,
						},
						"create_time": {
							Type:        schema.TypeString,
							Description: "The time format is: 20160323112310",
							Computed:    true,
						},
						"domain_id": {
							Type:        schema.TypeString,
							Description: "Corresponding domain ID",
							Computed:    true,
						},
						"domain_name": {
							Type:        schema.TypeString,
							Description: "Accelerated domain name",
							Computed:    true,
						},
						"operator": {
							Type:        schema.TypeString,
							Description: "Operator of this query",
							Computed:    true,
						},
						"origin_ips": {
							Type:        schema.TypeString,
							Description: "Accelerate the origin IP of a domain name",
							Computed:    true,
						},
						"service_type": {
							Type:        schema.TypeString,
							Description: "Service type for accelerated domain name",
							Computed:    true,
						},
						"domain_status": {
							Type:        schema.TypeString,
							Description: "Status of accelerated domain name.",
							Computed:    true,
						},
						"deploy_version": {
							Type:        schema.TypeString,
							Description: "Deployment version code",
							Computed:    true,
						},
						"cdn_service_status": {
							Type:        schema.TypeString,
							Description: "Does the domain name enable CDN acceleration services, Y and N?",
							Computed:    true,
						},
						"is_enabled": {
							Type:        schema.TypeString,
							Description: "Whether the accelerated domain name is enabled, Y and N?",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceWangSuCdnDomainsRead(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("data_source.wangsu_cdn_domain.read")

	request := &cdn.GetFuzzyPagingDomainListRequest{}
	if pageNumber, ok := data.Get("page_number").(int); ok && pageNumber > 0 {
		request.PageNumber = &pageNumber
	}

	if pageSize, ok := data.Get("page_size").(int); ok && pageSize > 0 {
		request.PageSize = &pageSize
	}

	if serviceType, ok := data.Get("service_type").(string); ok && serviceType != "" {
		request.ServiceType = &serviceType
	}

	if domainNameList, ok := data.Get("domain_name").([]interface{}); ok && len(domainNameList) > 0 {
		domainNames := make([]*string, len(domainNameList))
		for i, v := range domainNameList {
			if domainName, ok := v.(string); ok {
				domainNames[i] = &domainName
			}
		}
		request.DomainName = domainNames
	}

	if queryType, ok := data.Get("query_type").(string); ok && queryType != "" {
		request.QueryType = &queryType
	}

	if startTime, ok := data.Get("start_time").(string); ok && startTime != "" {
		request.StartTime = &startTime
	}

	if endTime, ok := data.Get("end_time").(string); ok && endTime != "" {
		request.EndTime = &endTime
	}

	if domainStatus, ok := data.Get("domain_status").(string); ok && domainStatus != "" {
		request.DomainStatus = &domainStatus
	}

	var response *cdn.GetFuzzyPagingDomainListResponse
	var requestId string
	var diags diag.Diagnostics
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		requestId, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseCdnClient().GetCdnDomainList(request)
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

	ids := make([]string, 0, len(response.ResultList))

	_ = data.Set("code", response.Code)
	_ = data.Set("x_cnc_request_id", requestId)
	_ = data.Set("total_count", response.TotalCount)
	_ = data.Set("total_page_number", response.TotalPageNumber)
	_ = data.Set("page_number", response.PageNumber)
	_ = data.Set("page_size", response.PageSize)

	if response.ResultList != nil {
		resultList := make([]map[string]interface{}, len(response.ResultList))
		for i, item := range response.ResultList {
			resultList[i] = map[string]interface{}{
				"cname":              item.Cname,
				"config_form_name":   item.ConfigFormName,
				"create_time":        item.CreateTime,
				"domain_id":          item.DomainId,
				"domain_name":        item.DomainName,
				"operator":           item.Operator,
				"origin_ips":         item.OriginIps,
				"service_type":       item.ServiceType,
				"domain_status":      item.DomainStatus,
				"deploy_version":     item.DeployVersion,
				"cdn_service_status": item.CdnServiceStatus,
				"is_enabled":         item.IsEnabled,
			}
			ids = append(ids, *item.DomainName)
		}
		_ = data.Set("result_list", resultList)
	}

	data.SetId(wangsuCommon.DataResourceIdsHash(ids))

	return nil
}
