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

func DataSourceSslCertificates() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSslCertificatesRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the SSL certificate to be queried.",
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
				Description: "Response data array.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"certificate_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Certificate ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Certificate name, unique to customer granularity.",
						},
						"comment": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Remarks on certificate file.",
						},
						"share_ssl": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Shared, optional values are true and false, true represents shared certificates, false represents unshared certificates, default is false. This certificate allows cross-customer use when share-ssl is true. (The API does not support cross-customer use certificates. Contact customer service for manual configuration if required.)",
						},
						"certificate_validity_from": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Certificate effective start time (CST), such as 2016-08-01 07:00:00.",
						},
						"certificate_validity_to": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Certificate effective end time (CST), such as 2018-08-01 19:00:00.",
						},
						"related_domains": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of domain names using the current certificate.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domain_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Accelerated domain name ID.",
									},
									"domain_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name of accelerated domain name.",
									},
								},
							},
						},
						"dns_names": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "dns-names",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"certificate_serial": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRT certificate serial number.",
						},
						"certificate_issuer": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRT certificate issuer.",
						},
					},
				},
			},
		},
	}
}

func dataSourceSslCertificatesRead(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("data_source.wangsu_ssl_certificates.read")
	request := &certicate.QueryCertificateListForTerraformRequest{}
	if name, ok := data.Get("name").(string); ok && name != "" {
		request.Name = &name
	}
	var response *certicate.QueryCertificateListForTerraformResponse
	var requestId string
	var diags diag.Diagnostics
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		requestId, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseSslCertificateClient().QueryCertificateList(request)
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
	ids := make([]string, 0, len(response.Data))

	_ = data.Set("code", response.Code)
	_ = data.Set("x_cnc_request_id", requestId)
	_ = data.Set("message", response.Message)

	var certificates []interface{}
	for _, certificate := range response.Data {
		certificateData := map[string]interface{}{
			"certificate_id":            certificate.CertificateId,
			"name":                      certificate.Name,
			"comment":                   certificate.Comment,
			"share_ssl":                 certificate.ShareSsl,
			"certificate_validity_from": certificate.CertificateValidityFrom,
			"certificate_validity_to":   certificate.CertificateValidityTo,
			"related_domains":           flattenRelatedDomains(certificate.RelatedDomains),
			"dns_names":                 certificate.DnsNames,
			"certificate_serial":        certificate.CertificateSerial,
			"certificate_issuer":        certificate.CertificateIssuer,
		}
		ids = append(ids, wangsuCommon.Int64ToStr(*certificate.CertificateId))
		certificates = append(certificates, certificateData)
	}
	_ = data.Set("data", certificates)
	data.SetId(wangsuCommon.DataResourceIdsHash(ids))
	log.Printf("data_source.wangsu_ssl_certificates.success")
	return nil
}

func flattenRelatedDomains(domains []*certicate.QueryCertificateListForTerraformResponseDataRelatedDomains) []interface{} {
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
