package certificateapplication

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wangsuCommon "github.com/wangsu-api/terraform-provider-wangsu/wangsu/common"
	"github.com/wangsu-api/wangsu-sdk-go/wangsu/certificateapplication"
)

func DataSourceSslCertificateApplications() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSslCertificateApplicationsRead,
		Schema: map[string]*schema.Schema{
			// 请求参数
			"order_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The order ID.",
			},
			"order_status": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The order status. Valid values: ACCEPT_SUCCESS (Accepted), APPLYING (Applying), ISSUE_SUCCESS (Issued), ISSUE_FAILURE (Issue Failed), REVOKED (Revoked), CANCELED (Canceled).",
			},
			"certificate_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The certificate name.",
			},
			"domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The domain name.",
			},
			"start_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The order creation time, format: yyyy-MM-dd HH:mm:ss.",
			},
			"end_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The order end time, format: yyyy-MM-dd HH:mm:ss.",
			},
			"page_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     100,
				Description: "Page size, 1-100, default is 100.",
			},
			"page_number": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Page number, starting from 1, must be greater than or equal to 1, default is 1.",
			},

			// 响应参数
			"total_number": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of records.",
			},
			"total_page_number": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of pages.",
			},

			"orders": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of certificate application orders.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"order_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The order ID.",
						},
						"order_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The order status.",
						},
						"certificate_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The certificate ID.",
						},
						"certificate_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The certificate name.",
						},
						"certificate_brand": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The certificate brand.",
						},
						"certificate_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The certificate type.",
						},
						"auto_renew": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Whether auto-renewal is enabled.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A description of the order or certificate.",
						},
						"common_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The common name of the certificate.",
						},
						"subject_alternative_names": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "The subject alternative names.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The creation time of the order.",
						},
					},
				},
			},
		},
	}
}

func dataSourceSslCertificateApplicationsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("data_source.wangsu_ssl_certificate_applications.read")
	// 构建请求参数
	request := &certificateapplication.ListCertificateApplicationOrdersForTerraformRequest{}

	if v, ok := d.GetOk("order_id"); ok {
		orderId := v.(string)
		request.OrderId = &orderId
	}

	if v, ok := d.GetOk("order_status"); ok {
		rawOrderStatuses := v.([]interface{})
		orderStatuses := make([]*string, 0, len(rawOrderStatuses))
		for _, s := range rawOrderStatuses {
			status := s.(string)
			orderStatuses = append(orderStatuses, &status)
		}
		request.OrderStatus = orderStatuses
	}

	if v, ok := d.GetOk("certificate_name"); ok {
		certName := v.(string)
		request.CertificateName = &certName
	}

	if v, ok := d.GetOk("domain"); ok {
		domain := v.(string)
		request.Domain = &domain
	}

	if v, ok := d.GetOk("start_time"); ok {
		startTime := v.(string)
		request.StartTime = &startTime
	}

	if v, ok := d.GetOk("end_time"); ok {
		endTime := v.(string)
		request.EndTime = &endTime
	}

	pageParam := &certificateapplication.ListCertificateApplicationOrdersForTerraformRequestPageParam{}

	if v, ok := d.GetOk("page_size"); ok {
		pageSize := v.(int)
		pageParam.PageSize = &pageSize
	}

	if v, ok := d.GetOk("page_number"); ok {
		pageNumber := v.(int) - 1
		pageParam.PageNumber = &pageNumber
	}

	request.PageParam = pageParam

	var response *certificateapplication.ListCertificateApplicationOrdersForTerraformResponse
	var err error

	// 使用 SDK 调用 API
	err = resource.RetryContext(ctx, time.Minute*2, func() *resource.RetryError {
		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseSslCertificateApplicationClient().ListCertificateApplication(request)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(err)
	}

	if response == nil || response.Data == nil {
		d.SetId("")
		return nil
	}

	// 设置唯一ID
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	// 设置分页信息
	if response.Data.PageInfo != nil {
		pageInfo := response.Data.PageInfo
		if pageInfo.TotalNumber != nil {
			_ = d.Set("total_number", *pageInfo.TotalNumber) // ignore error
		}
		if pageInfo.TotalPageNumber != nil {
			_ = d.Set("total_page_number", *pageInfo.TotalPageNumber) // ignore error
		}
	}

	// 转换订单列表
	orders := make([]map[string]interface{}, 0, len(response.Data.Orders))
	for _, order := range response.Data.Orders {
		orderMap := map[string]interface{}{
			"order_id":                  getStrValue(order.OrderId),
			"order_status":              getStrValue(order.OrderStatus),
			"certificate_id":            getInt(order.CertificateId),
			"certificate_name":          getStrValue(order.CertificateName),
			"certificate_brand":         getStrValue(order.CertificateBrand),
			"certificate_type":          getStrValue(order.CertificateType),
			"auto_renew":                getStrValue(order.AutoRenew),
			"description":               getStrValue(order.Description),
			"common_name":               getStrValue(order.CommonName),
			"subject_alternative_names": flattenSubjectAlternativeNames(order.SubjectAlternativeNames),
			"create_time":               getStrValue(order.CreateTime),
		}
		orders = append(orders, orderMap)
	}

	if err := d.Set("orders", orders); err != nil {
		return diag.FromErr(err)
	}
	log.Printf("data_source.wangsu_ssl_certificate_applications.read success")
	return nil
}

// 工具函数，安全获取指针值
func getStrValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
