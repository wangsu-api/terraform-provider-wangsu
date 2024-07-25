package appadomain

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wangsuCommon "github.com/wangsu-api/terraform-provider-wangsu/wangsu/common"
	appadomain "github.com/wangsu-api/wangsu-sdk-go/wangsu/appa/domain"
	"golang.org/x/net/context"
	"log"
	"time"
)

func DataSourceAppaDomainDetail() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppaDomainDetailRead,
		Schema: map[string]*schema.Schema{
			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Domain name you want to accelerate.A generic domain name is supported, starting with the symbol '.', such as .example.com.",
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
						"domain_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Domain ID",
						},
						"domain_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain name you want to accelerate.A generic domain name is supported, starting with the symbol '.', such as .example.com.",
						},
						"service_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The service type of the accelerated domain name. The value can be: appa: Application Acceleration",
						},
						"cname": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "cname",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "status",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "is enabled",
						},
						"origin_config": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Origin configuration. Example: 'originConfig':[{'level':1,'strategy':'robin','origin':[{'originIp':'1.1.1.1','weight':10},{'originIp':'2.2.2.2','weight':20}]},{'level':2,'strategy':'quick','origin':[{'originIp':'3.3.3.3','weight':10}]}]",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"level": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The level of the origin, which value can be an integer ranging from 1 to 5. Note:1. Must be configured level by level start from level 1. The same level cannot be configured repeatedly.2. The lower the value, the higher the priority.",
									},
									"strategy": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Origin selection strategy supports fast, robin and hash. The value can be: fast: Fast strategy, robin: Robin strategy,hash: Hash strategy",
									},
									"origin": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Origin information of a certain level. A level can be configured with multiple origin IP addresses or domain names.Example:'origin':[{'originIp':'1.1.1.1','weight':10},{'originIp':'2.2.2.2','weight':20}]",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"origin_ip": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Origin address, which can be an IP or domain name.",
												},
												"weight": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Weight, which is only useful for robin strategy. The value is an integer ranging from 1 to 10000. If this parameter is not specified, the default value is 10.",
												},
											},
										},
									},
								},
							},
						},
						"http_ports": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "HTTP port. The value is an integer ranging from 1 to 65535. Multiple ports are supported and can be configured in the following format:httpPorts:[1000,1001]Note: 1. Ports 2012, 2323, 2443, 4031, 12012, 20121, 57891, 62016, 65383, and 65529 do not support.2. The HTTP port and HTTPS port must be unique.3. At least one HTTP port or HTTPS port must be configured.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"https_ports": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "HTTPS port. The value is an integer ranging from 1 to 65535. Multiple ports are supported and can be configured in the following format:httpPorts:[1000,1001]Note: 1. Ports 2012, 2323, 2443, 4031, 12012, 20121, 57891, 62016, 65383, and 65529 do not support.2. The HTTP port and HTTPS port must be unique.3. At least one HTTP port or HTTPS port must be configured.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"carry_client_ip": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Carry client IP configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tcp_carry_client_ip": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Carry client IP configuration for TCP.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"carry_client_ip_enabled": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The function switch of TCP carry client IP. true: Enabled false: Disabled",
												},
												"protocol": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The protocol of TCP carry client IP. ws: The client IP is carried through TCP Option 78. The origin can get the client IP through the SDK provided by Wangsu or the F5 device. toa: Carry client IP through TOA. pp: Carry client IP through Proxy Protocol.",
												},
												"packet_num": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The number of packets carrying client IP other than SYN packet, when using ws or toa to carry client IP.",
												},
												"mode": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The mode of TCP carry client IP by ws or toa. 0: Only SYN packet carry client IP. 1: SYN packet and the first N packets carry client IP,N is specified by packetNum. 2: SYN packet, the first N packets, RST, PSH, URG and FIN/ACK packet carry client IP, N is specified by packetNum.",
												},
												"tcp_option_code": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The opcode to carry client IP by TOA.",
												},
											},
										},
									},
									"http_carry_client_ip": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Carry client IP configuration for HTTP.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"carry_client_ip_enabled": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The function switch of carry client IP through the HTTP header. true: Enabled. false: Disabled.",
												},
												"ports": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specify which acceleration ports to carry client IP, the function of carry client IP through the HTTP header is enabled.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"http_header_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The field name of HTTP header to carry client IP, such as X-Forwarded-For, Cdn-Src-Ip.",
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

func dataSourceAppaDomainDetailRead(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("data_source.wangsu_appa_domain_detail.read")

	domainName := data.Get("domain_name").(string)
	var response *appadomain.QueryAppaDomainForTerraformResponse
	var requestId string
	var diags diag.Diagnostics
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		requestId, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseAppaDomainClient().QueryAppaDomain(domainName)
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
	//读取response.Data的所有属性到map中
	var domainDetail = map[string]interface{}{
		"domain_id":       response.Data.DomainId,
		"domain_name":     response.Data.DomainName,
		"service_type":    response.Data.ServiceType,
		"cname":           response.Data.Cname,
		"status":          response.Data.Status,
		"enabled":         response.Data.Enabled,
		"origin_config":   flattenOriginConfig(response.Data.OriginConfig),
		"http_ports":      response.Data.HttpPorts,
		"https_ports":     response.Data.HttpsPorts,
		"carry_client_ip": flattenCarryClientIp(response.Data.CarryClientIp),
	}
	resultList = append(resultList, domainDetail)

	_ = data.Set("data", resultList)

	data.SetId(*response.Data.DomainName)
	log.Printf("data_source.wangsu_appa_domain_detail.read success")
	return nil
}

func flattenCarryClientIp(carryClientIp *appadomain.QueryAppaDomainForTerraformResponseDataCarryClientIp) interface{} {
	if carryClientIp == nil {
		return nil
	}
	var carryClientIpMap = map[string]interface{}{
		"tcp_carry_client_ip":  flattenTcpCarryClientIp(carryClientIp.TcpCarryClientIp),
		"http_carry_client_ip": flattenHttpCarryClientIp(carryClientIp.HttpCarryClientIp),
	}
	return []interface{}{carryClientIpMap}
}

func flattenTcpCarryClientIp(tcpCarryClientIp *appadomain.QueryAppaDomainForTerraformResponseDataCarryClientIpTcpCarryClientIp) interface{} {
	if tcpCarryClientIp == nil {
		return nil
	}
	var tcpCarryClientIpMap = map[string]interface{}{
		"carry_client_ip_enabled": *tcpCarryClientIp.CarryClientIpEnabled,
		"protocol":                *tcpCarryClientIp.Protocol,
		"packet_num":              *tcpCarryClientIp.PacketNum,
		"mode":                    *tcpCarryClientIp.Mode,
		"tcp_option_code":         *tcpCarryClientIp.TcpOptionCode,
	}
	return []interface{}{tcpCarryClientIpMap}
}

func flattenHttpCarryClientIp(httpCarryClientIp *appadomain.QueryAppaDomainForTerraformResponseDataCarryClientIpHttpCarryClientIp) interface{} {
	if httpCarryClientIp == nil {
		return nil
	}
	var httpCarryClientIpMap = map[string]interface{}{
		"carry_client_ip_enabled": *httpCarryClientIp.CarryClientIpEnabled,
		"ports":                   httpCarryClientIp.Ports,
		"http_header_name":        *httpCarryClientIp.HttpHeaderName,
	}
	return []interface{}{httpCarryClientIpMap}
}
