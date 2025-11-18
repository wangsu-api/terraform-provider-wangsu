package certificateapplication

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	wangsuCommon "github.com/wangsu-api/terraform-provider-wangsu/wangsu/common"
	"github.com/wangsu-api/wangsu-sdk-go/wangsu/certificateapplication"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceSslCertificateApplicationDetail() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSslCertificateApplicationDetailRead,
		Schema: map[string]*schema.Schema{
			"order_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The order ID.",
			},
			"purchase_record_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The purchase record ID.",
			},
			"code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The error code. 0 indicates success.",
			},
			"message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "message",
			},
			"data": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The order details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"order_id":          {Type: schema.TypeString, Computed: true, Description: "The order ID."},
						"order_status":      {Type: schema.TypeString, Computed: true, Description: "The status of the order."},
						"certificate_id":    {Type: schema.TypeInt, Computed: true, Description: "The certificate ID."},
						"certificate_name":  {Type: schema.TypeString, Computed: true, Description: "The name of the certificate."},
						"description":       {Type: schema.TypeString, Computed: true, Description: "A description of the order or certificate."},
						"auto_renew":        {Type: schema.TypeString, Computed: true, Description: "Indicates whether automatic renewal is enabled."},
						"certificate_brand": {Type: schema.TypeString, Computed: true, Description: "The brand of the certificate."},
						"certificate_spec":  {Type: schema.TypeString, Computed: true, Description: "The specification of the certificate."},
						"certificate_type":  {Type: schema.TypeString, Computed: true, Description: "The type of the certificate."},
						"algorithm":         {Type: schema.TypeString, Computed: true, Description: "The algorithm used for the certificate."},
						"auto_validate":     {Type: schema.TypeString, Computed: true, Description: "Indicates whether automatic validation is enabled."},
						"validate_method":   {Type: schema.TypeString, Computed: true, Description: "The validation method."},
						"auto_deploy":       {Type: schema.TypeString, Computed: true, Description: "Indicates whether automatic deployment is enabled."},
						"validity_days":     {Type: schema.TypeInt, Computed: true, Description: "The validity period of the certificate in days."},
						"identification_info": {Type: schema.TypeList, Computed: true, Description: "The identification information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"country":                   {Type: schema.TypeString, Computed: true, Description: "The country code."},
									"state":                     {Type: schema.TypeString, Computed: true, Description: "The state or province."},
									"city":                      {Type: schema.TypeString, Computed: true, Description: "The city."},
									"company":                   {Type: schema.TypeString, Computed: true, Description: "The company name."},
									"department":                {Type: schema.TypeString, Computed: true, Description: "The department name."},
									"common_name":               {Type: schema.TypeString, Computed: true, Description: "The common name of the certificate."},
									"email":                     {Type: schema.TypeString, Computed: true, Description: "The email address."},
									"street":                    {Type: schema.TypeString, Computed: true, Description: "The street address."},
									"subject_alternative_names": {Type: schema.TypeList, Elem: &schema.Schema{Type: schema.TypeString}, Computed: true, Description: "The subject alternative names."},
									"street1": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The second line of the street address.",
									},
									"phone": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Phone",
									},
									"postal_code": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Postal Code",
									},
								},
							},
						},
						"primary_certificate_brand": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Primary Certificate Brand",
						},
						"backup_certificate_brand": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Backup Certificate Brand",
						},
						"domain_type":   {Type: schema.TypeString, Computed: true, Description: "The domain type."},
						"create_time":   {Type: schema.TypeString, Computed: true, Description: "The creation time of the order."},
						"error_message": {Type: schema.TypeString, Computed: true, Description: "The error message, if any."},
						"remain_validity_days": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The remaining validity period in days.",
						},
						"dns_provider_infos": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Information about DNS providers.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The domain name associated with the certificate.",
									},
									"dns_api_access": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The DNS provider code. Supports CloudDNS (with optional certificate brands LE, TrustAsia, or GlobalSign) and CloudFlare (with optional certificate brands LE or ZeroSSL).",
									},

									"dns_provider_code": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The DNS provider code. Supports CloudDNS (with optional certificate brands LE, TrustAsia, or GlobalSign) and CloudFlare (with optional certificate brands LE or ZeroSSL).",
									},
									"enable_dns_alias_mode": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Indicates whether to use alias validation mode. Value range: true, false. Defaults to false.",
									},
									"validate_alias_domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The alias domain used for validation.",
									},
								},
							},
						},
						"org_validate_method": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Org Validate Method",
						},
						"admin": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Order admin contact.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"first_name": {Type: schema.TypeString, Computed: true, Description: "First name."},
									"last_name":  {Type: schema.TypeString, Computed: true, Description: "Last name."},
									"phone":      {Type: schema.TypeString, Computed: true, Description: "Phone."},
									"email":      {Type: schema.TypeString, Computed: true, Description: "Email."},
									"title":      {Type: schema.TypeString, Computed: true, Description: "Title."},
								},
							},
						},
						"tech": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Order tech contact.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"first_name": {Type: schema.TypeString, Computed: true, Description: "First name."},
									"last_name":  {Type: schema.TypeString, Computed: true, Description: "Last name."},
									"phone":      {Type: schema.TypeString, Computed: true, Description: "Phone."},
									"email":      {Type: schema.TypeString, Computed: true, Description: "Email."},
									"title":      {Type: schema.TypeString, Computed: true, Description: "Title."},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceSslCertificateApplicationDetailRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("data_source.wangsu_ssl_certificate_application_detail.read")
	orderId := d.Get("order_id").(string)
	purchaseRecordId := ""
	if v, ok := d.GetOk("purchase_record_id"); ok {
		purchaseRecordId = v.(string)
	}

	var diags diag.Diagnostics
	var response *certificateapplication.GetCertificateApplicationOrderForTerraformResponse
	var err error

	request := &certificateapplication.GetCertificateApplicationOrderForTerraformRequest{
		OrderId:          &orderId,
		PurchaseRecordId: nil,
	}
	if purchaseRecordId != "" {
		request.PurchaseRecordId = &purchaseRecordId
	}

	// SDK 查询，重试机制
	err = resource.RetryContext(ctx, time.Minute*2, func() *resource.RetryError {
		// 请根据实际 ProviderMeta 和 SDK 客户端替换此调用
		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseSslCertificateApplicationClient().GetCertificateApplicationDetail(request)
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
		d.SetId("")
		return nil
	}

	_ = d.Set("code", response.Code)
	_ = d.Set("message", response.Message)

	orderDetail := response.Data
	result := map[string]interface{}{
		"order_id":                  getStr(orderDetail.OrderId),
		"order_status":              getStr(orderDetail.OrderStatus),
		"certificate_id":            getInt(orderDetail.CertificateId),
		"certificate_name":          getStr(orderDetail.CertificateName),
		"certificate_brand":         getStr(orderDetail.CertificateBrand),
		"certificate_type":          getStr(orderDetail.CertificateType),
		"description":               getStr(orderDetail.Description),
		"algorithm":                 getStr(orderDetail.Algorithm),
		"auto_validate":             getStr(orderDetail.AutoValidate),
		"validate_method":           getStr(orderDetail.ValidateMethod),
		"auto_renew":                getStr(orderDetail.AutoRenew),
		"validity_days":             getInt(orderDetail.ValidityDays),
		"identification_info":       buildIdentificationInfo(orderDetail),
		"create_time":               getStr(orderDetail.CreateTime),
		"error_message":             getStr(orderDetail.ErrorMessage),
		"certificate_spec":          getStr(orderDetail.CertificateSpec),
		"domain_type":               getStr(orderDetail.DomainType),
		"auto_deploy":               getStr(orderDetail.AutoDeploy),
		"remain_validity_days":      getInt(orderDetail.RemainValidityDays),
		"dns_provider_infos":        flattenDnsProviderInfos(orderDetail.DnsProviderInfos),
		"primary_certificate_brand": getStr(orderDetail.PrimaryCertificateBrand),
		"backup_certificate_brand":  getStr(orderDetail.BackupCertificateBrand),
		"org_validate_method":       getStr(orderDetail.OrgValidateMethod),
		"admin":                     buildAdminLinkmanInfo(orderDetail.Admin),
		"tech":                      buildTechLinkmanInfo(orderDetail.Tech),
	}
	var resultList []interface{}
	resultList = append(resultList, result)
	_ = d.Set("data", resultList)

	d.SetId(orderId)
	log.Printf("data_source.wangsu_ssl_certificate_application_detail.read success")
	return nil
}

func buildTechLinkmanInfo(tech *certificateapplication.GetCertificateApplicationOrderForTerraformResponseDataTech) []interface{} {
	if tech == nil {
		return nil
	}
	var techInfo = map[string]interface{}{
		"first_name": getStr(tech.FirstName),
		"last_name":  getStr(tech.LastName),
		"phone":      getStr(tech.Phone),
		"email":      getStr(tech.Email),
		"title":      getStr(tech.Title),
	}
	var resultList []interface{}
	resultList = append(resultList, techInfo)
	return resultList
}

func buildAdminLinkmanInfo(admin *certificateapplication.GetCertificateApplicationOrderForTerraformResponseDataAdmin) []interface{} {
	if admin == nil {
		return nil
	}
	var adminInfo = map[string]interface{}{
		"first_name": getStr(admin.FirstName),
		"last_name":  getStr(admin.LastName),
		"phone":      getStr(admin.Phone),
		"email":      getStr(admin.Email),
		"title":      getStr(admin.Title),
	}
	var resultList []interface{}
	resultList = append(resultList, adminInfo)
	return resultList
}

func buildIdentificationInfo(orderDetail *certificateapplication.GetCertificateApplicationOrderForTerraformResponseData) []interface{} {
	var identificationInfo = map[string]interface{}{
		"country":                   getStr(orderDetail.Country),
		"state":                     getStr(orderDetail.State),
		"city":                      getStr(orderDetail.City),
		"company":                   getStr(orderDetail.Company),
		"department":                getStr(orderDetail.Department),
		"common_name":               getStr(orderDetail.CommonName),
		"email":                     getStr(orderDetail.Email),
		"street":                    getStr(orderDetail.Street),
		"street1":                   getStr(orderDetail.Street1),
		"subject_alternative_names": flattenSubjectAlternativeNames(orderDetail.SubjectAlternativeNames),
		"phone":                     getStr(orderDetail.Phone),
		"postal_code":               getStr(orderDetail.PostalCode),
	}
	var resultList []interface{}
	resultList = append(resultList, identificationInfo)
	return resultList
}

// 工具函数，安全转换指针类型
func getStr(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}
func getInt(v interface{}) int {
	switch val := v.(type) {
	case *int:
		if val != nil {
			return *val
		}
	case *int64:
		if val != nil {
			return int(*val)
		}
	case int:
		return val
	case int64:
		return int(val)
	}
	return 0
}

// 将 []*string 转为 []interface{}
func flattenSubjectAlternativeNames(names []*string) []interface{} {
	if names == nil {
		return []interface{}{}
	}
	result := make([]interface{}, 0, len(names))
	for _, n := range names {
		if n != nil {
			result = append(result, *n)
		} else {
			result = append(result, "")
		}
	}
	return result
}

func flattenDnsProviderInfos(infos []*certificateapplication.GetCertificateApplicationOrderForTerraformResponseDataDnsProviderInfos) []interface{} {
	if infos == nil {
		return []interface{}{}
	}
	result := make([]interface{}, 0, len(infos))
	for _, info := range infos {
		if info == nil {
			continue
		}
		item := map[string]interface{}{
			"domain":                getStr(info.Domain),
			"dns_provider_code":     getStr(info.DnsProviderCode),
			"dns_api_access":        getStr(info.DnsApiAccess),
			"enable_dns_alias_mode": getStr(info.EnableDnsAliasMode),
			"validate_alias_domain": getStr(info.ValidateAliasDomain),
		}
		result = append(result, item)
	}
	return result
}
