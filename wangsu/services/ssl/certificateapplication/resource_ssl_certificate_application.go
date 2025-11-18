package certificateapplication

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wangsuCommon "github.com/wangsu-api/terraform-provider-wangsu/wangsu/common"
	"github.com/wangsu-api/wangsu-sdk-go/wangsu/certificateapplication"
	"log"
	"time"
)

func ResourceSslCertificateApplication() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSslCertificateApplicationCreate,
		ReadContext:   resourceSslCertificateApplicationRead,
		DeleteContext: resourceSslCertificateApplicationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"certificate_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Name of the certificate. The certificate name cannot be the same as your existing certificate. Max 128 characters.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "A description of the new certificate. Max 256 characters.",
			},
			"auto_renew": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "false",
				Description: "Automatically renew your certificate. Allowed values: true, false (default).",
			},
			"certificate_brand": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Certificate Brand. Allowed values: LE, TrustAsia. If LE application fails, will switch to ZeroSSL.",
			},
			"certificate_spec": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Certificate Specification. See provider documentation for allowed values.",
			},
			"certificate_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Certificate Type. Allowed values: DV, OV. LE only supports DV.",
			},
			"algorithm": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Certificate Algorithm. Allowed values: RSA2048, ECDSA256.",
			},
			"auto_validate": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Automatically validate domain control. Allowed values: true, false.",
			},
			"validate_method": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Validate Method. Allowed values: HTTP, DNS. Wildcard domains require DNS.",
			},
			"auto_deploy": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Whether to deploy automatically. Allowed values: true, false.",
			},
			"validity_days": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Validity Days for the certificate. See certificateSpec for details.",
			},
			"identification_info": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Certificate Signing Request Information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"country":                   {Type: schema.TypeString, Optional: true, ForceNew: true, Description: "Country. ISO-3166 code, 2 characters."},
						"state":                     {Type: schema.TypeString, Optional: true, ForceNew: true, Description: "A state or province."},
						"city":                      {Type: schema.TypeString, Optional: true, ForceNew: true, Description: "A city."},
						"company":                   {Type: schema.TypeString, Optional: true, ForceNew: true, Description: "A company name."},
						"department":                {Type: schema.TypeString, Optional: true, ForceNew: true, Description: "The department associated with the certificate."},
						"common_name":               {Type: schema.TypeString, Required: true, ForceNew: true, Description: "A common name of the certificate."},
						"email":                     {Type: schema.TypeString, Optional: true, ForceNew: true, Description: "Email address."},
						"street":                    {Type: schema.TypeString, Optional: true, ForceNew: true, Description: "The street where the company is located."},
						"street1":                   {Type: schema.TypeString, Optional: true, ForceNew: true, Description: "street1."},
						"subject_alternative_names": {Type: schema.TypeList, Required: true, ForceNew: true, Elem: &schema.Schema{Type: schema.TypeString}, Description: "Hostnames that this certificate will serve."},
						"phone":                     {Type: schema.TypeString, Optional: true, ForceNew: true, Description: "phone."},
						"postal_code":               {Type: schema.TypeString, Optional: true, ForceNew: true, Description: "postalCode."},
					},
				},
			},
			"backup_certificate_brand": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Backup Certificate Brand. See provider documentation for allowed values.",
			},
			"org_validate_method": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Organization validate method. Allowed values: default, self_validate, none.",
			},
			"domain_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Domain type. Allowed values: single, multi.",
			},
			"admin": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Order admin contact.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"first_name": {Type: schema.TypeString, Optional: true, ForceNew: true, Description: "First name."},
						"last_name":  {Type: schema.TypeString, Optional: true, ForceNew: true, Description: "Last name."},
						"phone":      {Type: schema.TypeString, Optional: true, ForceNew: true, Description: "Phone."},
						"email":      {Type: schema.TypeString, Optional: true, ForceNew: true, Description: "Email."},
						"title":      {Type: schema.TypeString, Optional: true, ForceNew: true, Description: "Title."},
					},
				},
			},
			"tech": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Order tech contact.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"first_name": {Type: schema.TypeString, Optional: true, ForceNew: true, Description: "First name."},
						"last_name":  {Type: schema.TypeString, Optional: true, ForceNew: true, Description: "Last name."},
						"phone":      {Type: schema.TypeString, Optional: true, ForceNew: true, Description: "Phone."},
						"email":      {Type: schema.TypeString, Optional: true, ForceNew: true, Description: "Email."},
						"title":      {Type: schema.TypeString, Optional: true, ForceNew: true, Description: "Title."},
					},
				},
			},
			"dns_provider_infos": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "DNS Provider Information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain":                {Type: schema.TypeString, Required: true, ForceNew: true, Description: "Domain."},
						"dns_provider_code":     {Type: schema.TypeString, Optional: true, ForceNew: true, Description: "DNS Provider Code. Support CloudDNS."},
						"dns_api_access":        {Type: schema.TypeString, Optional: true, ForceNew: true, Description: "DNS API Access. JSON format, e.g. accessKey/secretKey for CloudDNS."},
						"enable_dns_alias_mode": {Type: schema.TypeString, Optional: true, ForceNew: true, Description: "Whether to use alias verification. Allowed values: true, false."},
						"validate_alias_domain": {Type: schema.TypeString, Optional: true, ForceNew: true, Description: "Validate alias"},
					},
				},
			},
		},
	}
}

func resourceSslCertificateApplicationDelete(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_ssl_certificate_application.delete")
	response, diagnostics, hasError := getCertificateApplicationDetail(context, data, meta)
	if hasError {
		return diagnostics
	}
	if getStr(response.Data.OrderStatus) != "ACCEPT_SUCCESS" && getStr(response.Data.OrderStatus) != "APPLYING" {
		return nil
	}

	client := meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseSslCertificateApplicationClient()
	orderId := data.Id()
	if orderId == "" {
		data.SetId("")
		return nil
	}
	// Wangsu API may use a cancel method for certificate application deletion
	request := &certificateapplication.CancelCertificateApplicationOrderForTerraformRequest{
		OrderId: &orderId,
	}
	var err error
	_, _, err = client.CancelCertificateApplication(request)
	if err != nil {
		return diag.FromErr(err)
	}
	data.SetId("")
	return nil
}

func resourceSslCertificateApplicationCreate(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_ssl_certificate_application.create")
	var diags diag.Diagnostics
	request := &certificateapplication.CreateCertificateApplicationOrderForTerraformRequest{}
	if v, ok := data.GetOk("certificate_name"); ok {
		request.CertificateName = strPtr(v)
	}
	if v, ok := data.GetOk("description"); ok {
		request.Description = strPtr(v)
	}
	if v, ok := data.GetOk("auto_renew"); ok {
		request.AutoRenew = strPtr(v)
	}
	if v, ok := data.GetOk("certificate_brand"); ok {
		request.CertificateBrand = strPtr(v)
	}
	if v, ok := data.GetOk("certificate_spec"); ok {
		request.CertificateSpec = strPtr(v)
	}
	if v, ok := data.GetOk("certificate_type"); ok {
		request.CertificateType = strPtr(v)
	}
	if v, ok := data.GetOk("algorithm"); ok {
		request.Algorithm = strPtr(v)
	}
	if v, ok := data.GetOk("auto_validate"); ok {
		request.AutoValidate = strPtr(v)
	}
	if v, ok := data.GetOk("validate_method"); ok {
		request.ValidateMethod = strPtr(v)
	}
	if v, ok := data.GetOk("auto_deploy"); ok {
		request.AutoDeploy = strPtr(v)
	}
	if v, ok := data.GetOk("validity_days"); ok {
		val := v.(int)
		request.ValidityDays = &val
	}
	if v, ok := data.GetOk("backup_certificate_brand"); ok {
		request.BackupCertificateBrand = strPtr(v)
	}
	if v, ok := data.GetOk("org_validate_method"); ok {
		request.OrgValidateMethod = strPtr(v)
	}
	if v, ok := data.GetOk("domain_type"); ok {
		request.DomainType = strPtr(v)
	}
	if v, ok := data.GetOk("batch_apply"); ok {
		request.BatchApply = strPtr(v)
	}
	// identification_info
	if v, ok := data.GetOk("identification_info"); ok {
		list := v.([]interface{})
		if len(list) > 0 && list[0] != nil {
			info := list[0].(map[string]interface{})
			request.IdentificationInfo = expandIdentificationInfo(info)
		}
	}
	// admin
	if v, ok := data.GetOk("admin"); ok {
		list := v.([]interface{})
		if len(list) > 0 && list[0] != nil {
			info := list[0].(map[string]interface{})
			request.Admin = expandAdminLinkman(info)
		}
	}
	// tech
	if v, ok := data.GetOk("tech"); ok {
		list := v.([]interface{})
		if len(list) > 0 && list[0] != nil {
			info := list[0].(map[string]interface{})
			request.Tech = expandTechLinkman(info)
		}
	}
	// dns_provider_infos
	if v, ok := data.GetOk("dns_provider_infos"); ok {
		list := v.([]interface{})
		request.DnsProviderInfos = expandDnsProviderInfos(list)
	}

	//start to create a domain in 2 minutes
	var response *certificateapplication.CreateCertificateApplicationOrderForTerraformResponse
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseSslCertificateApplicationClient().CreateCertificateApplication(request)
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
	data.SetId(*response.Data.OrderId)
	return nil
}

func resourceSslCertificateApplicationRead(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_ssl_certificate_application.read")
	response, diagnostics, hasError := getCertificateApplicationDetail(context, data, meta)
	if hasError {
		return diagnostics
	}

	orderDetail := response.Data
	_ = data.Set("certificate_name", getStr(orderDetail.CertificateName))
	_ = data.Set("description", getStr(orderDetail.Description))
	_ = data.Set("auto_renew", getStr(orderDetail.AutoRenew))
	_ = data.Set("certificate_brand", getStr(orderDetail.CertificateBrand))
	_ = data.Set("certificate_spec", getStr(orderDetail.CertificateSpec))
	_ = data.Set("certificate_type", getStr(orderDetail.CertificateType))
	_ = data.Set("algorithm", getStr(orderDetail.Algorithm))
	_ = data.Set("auto_validate", getStr(orderDetail.AutoValidate))
	_ = data.Set("validate_method", getStr(orderDetail.ValidateMethod))
	_ = data.Set("auto_deploy", getStr(orderDetail.AutoDeploy))
	_ = data.Set("validity_days", getInt(orderDetail.ValidityDays))
	// identification_info
	identInfo := flattenIdentificationInfo(orderDetail)
	_ = data.Set("identification_info", identInfo)
	_ = data.Set("backup_certificate_brand", getStr(orderDetail.BackupCertificateBrand))
	_ = data.Set("org_validate_method", getStr(orderDetail.OrgValidateMethod))
	_ = data.Set("domain_type", getStr(orderDetail.DomainType))
	// admin, tech
	_ = data.Set("admin", flattenAdminLinkman(orderDetail.Admin))
	_ = data.Set("tech", flattenTechLinkman(orderDetail.Tech))
	//dns_provider_infos
	_ = data.Set("dns_provider_infos", flattenDnsProviderInfos(orderDetail.DnsProviderInfos))
	return nil
}

func flattenTechLinkman(tech *certificateapplication.GetCertificateApplicationOrderForTerraformResponseDataTech) []interface{} {
	if tech == nil {
		return nil
	}
	techInfo := map[string]interface{}{
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

func flattenAdminLinkman(admin *certificateapplication.GetCertificateApplicationOrderForTerraformResponseDataAdmin) []interface{} {
	if admin == nil {
		return nil
	}
	adminInfo := map[string]interface{}{
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

func getCertificateApplicationDetail(context context.Context, data *schema.ResourceData, meta interface{}) (*certificateapplication.GetCertificateApplicationOrderForTerraformResponse, diag.Diagnostics, bool) {
	var diags diag.Diagnostics
	orderId := data.Id()

	request := &certificateapplication.GetCertificateApplicationOrderForTerraformRequest{}
	request.OrderId = &orderId

	var response *certificateapplication.GetCertificateApplicationOrderForTerraformResponse

	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseSslCertificateApplicationClient().GetCertificateApplicationDetail(request)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return nil, diags, true
	}
	if response == nil || response.Data == nil {
		data.SetId("")
		return nil, nil, true
	}
	return response, nil, false
}

func strPtr(v interface{}) *string {
	if v == nil {
		return nil
	}
	s := v.(string)
	return &s
}

func expandIdentificationInfo(info map[string]interface{}) *certificateapplication.CreateCertificateApplicationOrderForTerraformRequestIdentificationInfo {
	obj := &certificateapplication.CreateCertificateApplicationOrderForTerraformRequestIdentificationInfo{}
	if v, ok := info["country"]; ok {
		obj.Country = strPtr(v)
	}
	if v, ok := info["state"]; ok {
		obj.State = strPtr(v)
	}
	if v, ok := info["city"]; ok {
		obj.City = strPtr(v)
	}
	if v, ok := info["company"]; ok {
		obj.Company = strPtr(v)
	}
	if v, ok := info["department"]; ok {
		obj.Department = strPtr(v)
	}
	if v, ok := info["common_name"]; ok {
		obj.CommonName = strPtr(v)
	}
	if v, ok := info["email"]; ok {
		obj.Email = strPtr(v)
	}
	if v, ok := info["street"]; ok {
		obj.Street = strPtr(v)
	}
	if v, ok := info["subject_alternative_names"]; ok {
		arr := v.([]interface{})
		var names []*string
		for _, n := range arr {
			if n != nil {
				str := n.(string)
				names = append(names, &str)
			}
		}
		obj.SubjectAlternativeNames = names
	}
	if v, ok := info["street1"]; ok {
		obj.Street1 = strPtr(v)
	}
	if v, ok := info["phone"]; ok {
		obj.Phone = strPtr(v)
	}
	if v, ok := info["postal_code"]; ok {
		obj.PostalCode = strPtr(v)
	}
	return obj
}

func expandAdminLinkman(info map[string]interface{}) *certificateapplication.CreateCertificateApplicationOrderForTerraformRequestAdmin {
	obj := &certificateapplication.CreateCertificateApplicationOrderForTerraformRequestAdmin{}
	if v, ok := info["first_name"]; ok {
		obj.FirstName = strPtr(v)
	}
	if v, ok := info["last_name"]; ok {
		obj.LastName = strPtr(v)
	}
	if v, ok := info["phone"]; ok {
		obj.Phone = strPtr(v)
	}
	if v, ok := info["email"]; ok {
		obj.Email = strPtr(v)
	}
	if v, ok := info["title"]; ok {
		obj.Title = strPtr(v)
	}
	return obj
}
func expandTechLinkman(info map[string]interface{}) *certificateapplication.CreateCertificateApplicationOrderForTerraformRequestTech {
	obj := &certificateapplication.CreateCertificateApplicationOrderForTerraformRequestTech{}
	if v, ok := info["first_name"]; ok {
		obj.FirstName = strPtr(v)
	}
	if v, ok := info["last_name"]; ok {
		obj.LastName = strPtr(v)
	}
	if v, ok := info["phone"]; ok {
		obj.Phone = strPtr(v)
	}
	if v, ok := info["email"]; ok {
		obj.Email = strPtr(v)
	}
	if v, ok := info["title"]; ok {
		obj.Title = strPtr(v)
	}
	return obj
}

func expandDnsProviderInfos(list []interface{}) []*certificateapplication.CreateCertificateApplicationOrderForTerraformRequestDnsProviderInfos {
	var result []*certificateapplication.CreateCertificateApplicationOrderForTerraformRequestDnsProviderInfos
	for _, item := range list {
		if item == nil {
			continue
		}
		info := item.(map[string]interface{})
		obj := &certificateapplication.CreateCertificateApplicationOrderForTerraformRequestDnsProviderInfos{}
		if v, ok := info["domain"]; ok {
			obj.Domain = strPtr(v)
		}
		if v, ok := info["dns_provider_code"]; ok {
			obj.DnsProviderCode = strPtr(v)
		}
		if v, ok := info["dns_api_access"]; ok {
			obj.DnsApiAccess = strPtr(v)
		}
		if v, ok := info["enable_dns_alias_mode"]; ok {
			obj.EnableDnsAliasMode = strPtr(v)
		}
		if v, ok := info["validate_alias_domain"]; ok {
			obj.ValidateAliasDomain = strPtr(v)
		}
		result = append(result, obj)
	}
	return result
}

func flattenIdentificationInfo(orderDetail *certificateapplication.GetCertificateApplicationOrderForTerraformResponseData) []interface{} {
	m := map[string]interface{}{
		"country":                   getStr(orderDetail.Country),
		"state":                     getStr(orderDetail.State),
		"city":                      getStr(orderDetail.City),
		"company":                   getStr(orderDetail.Company),
		"department":                getStr(orderDetail.Department),
		"common_name":               getStr(orderDetail.CommonName),
		"email":                     getStr(orderDetail.Email),
		"street":                    getStr(orderDetail.Street),
		"subject_alternative_names": flattenSubjectAlternativeNames(orderDetail.SubjectAlternativeNames),
		"street1":                   getStr(orderDetail.Street1),
		"phone":                     getStr(orderDetail.Phone),
		"postal_code":               getStr(orderDetail.PostalCode),
	}
	return []interface{}{m}
}
