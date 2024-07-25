package certificate

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wangsuCommon "github.com/wangsu-api/terraform-provider-wangsu/wangsu/common"
	certicate "github.com/wangsu-api/wangsu-sdk-go/wangsu/ssl/certificate"
	"golang.org/x/net/context"
	"log"
	"time"
)

func DataSourceSslCertificateDetail() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSslCertificateDetailRead,
		Schema: map[string]*schema.Schema{
			"certificate_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "ID of the SSL certificate to be queried.",
			},
			// computed
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
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"certificate_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "certificate Id",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "certificate name",
						},
						"cert": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Certificate, PEM certificate.",
						},
						"key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Private key of the certificate, PEM certificate.",
						},
						"domains": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domain_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Domain ID",
									},
									"domain_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Domain name",
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

func dataSourceSslCertificateDetailRead(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("data_source.wangsu_ssl_certificate_detail.read")

	certificateId := data.Get("certificate_id").(int)
	var response *certicate.QueryCertificateForTerraformResponse
	var requestId string
	var diags diag.Diagnostics
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		requestId, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseSslCertificateClient().QueryCertificate(int64(certificateId))
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

	_ = data.Set("code", response.Code)
	_ = data.Set("x_cnc_request_id", requestId)
	_ = data.Set("message", response.Message)

	var resultList []interface{}
	var certificateDetail = map[string]interface{}{
		"certificate_id": response.Data.CertificateId,
		"name":           response.Data.Name,
		"cert":           response.Data.Certificate,
		"key":            response.Data.PrivateKey,
		"domains":        flattenDomains(response.Data.Domains),
	}
	resultList = append(resultList, certificateDetail)

	_ = data.Set("data", resultList)

	data.SetId(wangsuCommon.Int64ToStr(*response.Data.CertificateId))
	log.Printf("data_source.wangsu_ssl_certificate_detail.read success")
	return nil

}

func flattenDomains(domains []*certicate.QueryCertificateForTerraformResponseDataDomains) []interface{} {
	var result []interface{}
	for _, d := range domains {
		domain := map[string]interface{}{
			"domain_id":   d.DomainId,
			"domain_name": d.DomainName,
		}
		result = append(result, domain)
	}
	return result
}
