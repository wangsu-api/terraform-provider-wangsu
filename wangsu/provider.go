package wangsu

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wangsuCommon "github.com/wangsu/terraform-provider-wangsu/wangsu/common"
	"github.com/wangsu/terraform-provider-wangsu/wangsu/connectivity"
	"github.com/wangsu/terraform-provider-wangsu/wangsu/services/cdn/domain"
	waapCustomizerule "github.com/wangsu/terraform-provider-wangsu/wangsu/services/waap/customizerule"
	waapDomain "github.com/wangsu/terraform-provider-wangsu/wangsu/services/waap/domain"
	waapRatelimit "github.com/wangsu/terraform-provider-wangsu/wangsu/services/waap/ratelimit"
	waapWhitelist "github.com/wangsu/terraform-provider-wangsu/wangsu/services/waap/whitelist"
	sdkCommon "github.com/wangsu/wangsu-sdk-go/common"
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
				DefaultFunc:  schema.EnvDefaultFunc(PROVIDER_PROTOCOL, "HTTPS"),
				ValidateFunc: wangsuCommon.ValidateAllowedStringValue([]string{"HTTP", "HTTPS"}),
				Description:  "The protocol of the API request. Valid values: `HTTP` and `HTTPS`. Default is `HTTPS`.",
			},
			"domain": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(PROVIDER_DOMAIN, nil),
				Description: "The root domain of the API request, e.g. `api.example.com`. It must be provided.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"wangsu_cdn_domain":         domain.ResourceCdnDomain(),
			"wangsu_waap_whitelist":     waapWhitelist.ResourceWaapWhitelist(),
			"wangsu_waap_customizerule": waapCustomizerule.ResourceWaapCustomizeRule(),
			"wangsu_waap_ratelimit":     waapRatelimit.ResourceWaapRateLimit(),
			"wangsu_waap_domain_copy":   waapDomain.ResourceWaapDomainCopy(),
			"wangsu_waap_domain":        waapDomain.ResourceWaapDomain(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"wangsu_cdn_domains":        domain.DataSourceWangSuCdnDomains(),
			"wangsu_waap_whitelist":     waapWhitelist.DataSourceWaapWhitelist(),
			"wangsu_waap_customizerule": waapCustomizerule.DataSourceCustomizeRule(),
			"wangsu_waap_ratelimit":     waapRatelimit.DataSourceRateLimit(),
			"wangsu_waap_domain":        waapDomain.DataSourceWaapDomain(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	var (
		secretId  string
		secretKey string
		protocol  string
		domain    string
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

	var wangSuClient WangSuClient
	wangSuClient.apiV3Conn = &connectivity.WangSuClient{
		Credential: sdkCommon.NewCredential(secretId, secretKey),
		Protocol:   protocol,
		Domain:     domain,
	}

	return &wangSuClient, nil
}

// GetAPIV3Conn 返回访问云 API 的客户端连接对象
func (meta *WangSuClient) GetAPIV3Conn() *connectivity.WangSuClient {
	return meta.apiV3Conn
}
