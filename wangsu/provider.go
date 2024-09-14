package wangsu

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wangsuCommon "github.com/wangsu-api/terraform-provider-wangsu/wangsu/common"
	"github.com/wangsu-api/terraform-provider-wangsu/wangsu/connectivity"
	appadomain "github.com/wangsu-api/terraform-provider-wangsu/wangsu/services/appa/domain"
	"github.com/wangsu-api/terraform-provider-wangsu/wangsu/services/cdn/domain"
	"github.com/wangsu-api/terraform-provider-wangsu/wangsu/services/ssl/certificate"
	waapCustomizerule "github.com/wangsu-api/terraform-provider-wangsu/wangsu/services/waap/customizerule"
	waapDomain "github.com/wangsu-api/terraform-provider-wangsu/wangsu/services/waap/domain"
	waapRatelimit "github.com/wangsu-api/terraform-provider-wangsu/wangsu/services/waap/ratelimit"
	waapWhitelist "github.com/wangsu-api/terraform-provider-wangsu/wangsu/services/waap/whitelist"
	sdkCommon "github.com/wangsu-api/wangsu-sdk-go/common"
)

const (
	PROVIDER_SECRET_ID  = "WANGSU_SECRET_ID"
	PROVIDER_SECRET_KEY = "WANGSU_SECRET_KEY"
	PROVIDER_PROTOCOL   = "WANGSU_PROTOCOL"
	PROVIDER_DOMAIN     = "WANGSU_DOMAIN"
)

type WangSuClient struct {
	apiV3Conn *connectivity.WangSuClient
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"secret_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(PROVIDER_SECRET_ID, nil),
				Description: "This is the wangsu access key. It must be provided, but it can also be sourced from the `WANGSU_SECRET_ID` environment variable.",
			},
			"secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(PROVIDER_SECRET_KEY, nil),
				Description: "This is the wangsu secret key. It must be provided, but it can also be sourced from the `WANGSU_SECRET_KEY` environment variable.",
				Sensitive:   true,
			},
			"protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc(PROVIDER_PROTOCOL, "https"),
				ValidateFunc: wangsuCommon.ValidateAllowedStringValue([]string{"http", "https"}),
				Description:  "(Optional)The protocol of the API request. Valid values: `http` and `https`. Default is `https`.",
			},
			"domain": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(PROVIDER_DOMAIN, nil),
				Description: "(Optional)The root domain of the API request.Default is `open.chinanetcenter.com`. It is optional",
			},
			"service_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "(Optional)Security service type. Please enter a specific service type, if you purchase multiple security services.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"wangsu_cdn_domain":         domain.ResourceCdnDomain(),
			"wangsu_ssl_certificate":    certificate.ResourceSslCertificate(),
			"wangsu_appa_domain":        appadomain.ResourceAppaDomain(),
			"wangsu_waap_whitelist":     waapWhitelist.ResourceWaapWhitelist(),
			"wangsu_waap_customizerule": waapCustomizerule.ResourceWaapCustomizeRule(),
			"wangsu_waap_ratelimit":     waapRatelimit.ResourceWaapRateLimit(),
			"wangsu_waap_domain_copy":   waapDomain.ResourceWaapDomainCopy(),
			"wangsu_waap_domain":        waapDomain.ResourceWaapDomain(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"wangsu_cdn_domains":            domain.DataSourceWangSuCdnDomains(),
			"wangsu_cdn_domain_detail":      domain.DataSourceWangSuCdnDomainDetail(),
			"wangsu_ssl_certificate_detail": certificate.DataSourceSslCertificateDetail(),
			"wangsu_appa_domain_detail":     appadomain.DataSourceAppaDomainDetail(),
			"wangsu_ssl_certificates":       certificate.DataSourceSslCertificates(),
			"wangsu_waap_whitelist":         waapWhitelist.DataSourceWaapWhitelist(),
			"wangsu_waap_whitelists":        waapWhitelist.DataSourceWaapWhitelists(),
			"wangsu_waap_customizerule":     waapCustomizerule.DataSourceCustomizeRule(),
			"wangsu_waap_customizerules":    waapCustomizerule.DataSourceCustomizeRules(),
			"wangsu_waap_ratelimit":         waapRatelimit.DataSourceRateLimit(),
			"wangsu_waap_ratelimits":        waapRatelimit.DataSourceRateLimits(),
			"wangsu_waap_domain":            waapDomain.DataSourceWaapDomain(),
			"wangsu_waap_domains":           waapDomain.DataSourceWaapDomains(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var (
		secretId    string
		secretKey   string
		protocol    string
		domain      string
		serviceType string
	)
	if v, ok := d.GetOk("secret_id"); ok {
		secretId = v.(string)
	}

	if v, ok := d.GetOk("secret_key"); ok {
		secretKey = v.(string)
	}
	if v, ok := d.GetOk("protocol"); ok {
		protocol = v.(string)
	}

	if v, ok := d.GetOk("domain"); ok {
		domain = v.(string)
	}

	if v, ok := d.GetOk("service_type"); ok {
		serviceType = v.(string)
	}

	var wangSuClient WangSuClient
	wangSuClient.apiV3Conn = &connectivity.WangSuClient{
		Credential:  sdkCommon.NewCredential(secretId, secretKey),
		HttpProfile: sdkCommon.NewHttpProfile(domain, protocol, serviceType),
	}

	return &wangSuClient, nil
}

// GetAPIV3Conn 返回访问云 API 的客户端连接对象
func (client *WangSuClient) GetAPIV3Conn() *connectivity.WangSuClient {
	return client.apiV3Conn
}
