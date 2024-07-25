package domain

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wangsuCommon "github.com/wangsu-api/terraform-provider-wangsu/wangsu/common"
	cdn "github.com/wangsu-api/wangsu-sdk-go/wangsu/cdn/domain"
	"log"
	"time"
)

func DataSourceWangSuCdnDomainDetail() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWangSuCdnDomainDetailRead,
		Schema: map[string]*schema.Schema{
			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Accelerated domain name",
			},
			//cumputed
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
						"domain_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "domain ID",
						},
						"cname": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "cname",
						},
						"domain_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Need to access the domain name of the CDN. a generic domain name is supported, starting with the symbol '.'. Example: .example.com also contains a multilevel 'a.b.example.com'. If example.com is filed, the domain name xx.example.com does not need to be filed.",
						},
						"service_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The service type of the accelerated domain name (only one service type can be submitted at a time): web/web-https: Web page acceleration/Web page acceleration-https wsa/wsa-https: Full-station acceleration/full-station acceleration-https vodstream/vod-https: on-demand acceleration/on-demand acceleration-https download/dl-https: Download Acceleration/Download Acceleration-https livestream/live-https/cloudv-live: livestream acceleration v6sa/osv6: IPv6 Security&Acceleration Solution/IPv6 One-stop Solution Note: https in the code, such as web-https does not represent immediate support for https access, you need to upload the certificate to support https.",
						},
						"service_areas": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The acceleration area of the acceleration domain, if the resource coverage needs to be limited according to the area, the acceleration area needs to be specified. When no acceleration area is specified, we will provide acceleration services with optimal resource coverage according to the service area opened by the customer. Multiple regions are separated by semicolons, and the supported regions are as follows: cn (Mainland China), am (Americas), emea (Europe, Middle East, Africa), apac (Asia-Pacific region).",
						},
						"comment": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Remarks, up to 1000 characters",
						},
						"header_of_client_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Pass the response header of client IP. The optional values are Cdn-Src-Ip and X-Forwarded-For. The default value is Cdn-Src-Ip.",
						},
						"origin_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"origin_ips": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Origin address, which can be an IP or domain name. Multiple IPs are supported, separated by semicolons. Only one domain name is allowed. IP and domain name cannot exist at the same time. The length cannot exceed 500 characters. The number of IPs cannot exceed 15.",
									},
									"default_origin_host_header": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Back-to-origin HOST, used to change the HOST field in the back-to-origin HTTP request header. The supported formats are: ① domain name ② ip. Note: Must comply with the ip/domain name format specification. If it is a domain name, the length of the domain name must be less than or equal to 128 characters.",
									},
									"use_range": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "useRange",
									},
									"follow301": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "follow301",
									},
									"follow302": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "follow302",
									},
									"adv_src_setting": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"use_adv_src": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Use advance origin config, the optional values are true and false, true means to use advance origin config, false means not to use advance origin config",
												},
												"detect_url": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The advanced source monitors the url, and requests <master-ips> through the url. If the response is not 2**, 3** response, it is considered that the primary source ip is faulty, and <backup-ips> is used at this time.",
												},
												"detect_period": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Advanced source monitoring period, in seconds, optional as an integer greater than or equal to 0, 0 means no monitoring",
												},
												"master_ips": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Description: "The advanced source mainly returns the source IP. Multiple IPs are separated by a semicolon \";\", and the return source IP cannot be repeated.",
												},
												"backup_ips": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Description: "Advanced source backup source IP, multiple IPs are separated by semicolon \";\", and the return source IP cannot be duplicated.",
												},
											},
										},
										Description: "advance origin config",
									},
								},
							},
							Description: "Back to origin policy settings for setting source site information and return source policies for accelerated domain names",
						},
						"ssl": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"use_ssl": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Use a certificate, the optional values are true and false, true means to use the certificate, false means not to use the certificate",
									},
									"ssl_certificate_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Use sni certificate, the optional values are true and false, true means use sni certificate, false means use shared certificate (not supported)",
									},
									"backup_certificate_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Backup certificate ID",
									},
									"gm_certificate_ids": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Description: "SM2 certificate IDS",
									},
									"tls_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "TLS version. Optional values: SSLv3,TLSv1,TLSv1.1,TLSv1.2,TLSv1.3",
									},
									"enable_ocsp": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Enable OCSP(Online Certificate Status Protocol).",
									},
									"ssl_cipher_suite": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "This optional object is used to specify a colon separated list of cipher suites which are permitted when clients negotiate security settings to access your content. Cipher suites which you can specify are: LOW, ALL:!LOW, HIGH, !EXPORT, !aNULL, !RC4, !DH, !SHA, !MD5, @STRENGTH,  AES128-SHA, AES256-SHA, AES128-SHA256, AES256-SHA256, AES128-GCM-SHA256, AES256-GCM-SHA384, ECDHE-RSA-AES128-SHA, ECDHE-RSA-AES256-SHA, ECDHE-RSA-AES128-SHA256, ECDHE-RSA-AES256-SHA384, ECDHE-RSA-AES128-GCM-SHA256, and ECDHE-RSA-AES256-GCM-SHA384. These cipher suites are a subset of those supported by OpenSSL, https://www.openssl.org/docs/man1.0.2/man1/ciphers.html. Please note that !MD5 or !SHA must appear after HIGH..",
									},
								},
							},
							Description: "SSL settings, to bind a certificate with the accelerated domain. You can use the interface [AddCertificate] to upload your  certificates. If you want to modify a certificate, please use the interface: [UpdateCertificate]",
						},
						"cache_time_behaviors": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"data_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "dataId is to indicate a specific group configuration when the client has multiple groups of configurations.",
									},
									"path_pattern": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The url matching mode supports fuzzy regularization. If all matches, the input parameters can be configured as: *",
									},
									"except_path_pattern": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Exceptional url matching mode, except for some URLs: such as abc.jpg, do not do anti-theft chain function. E.g: ^https?://[^/]+/.*\\.m3u8",
									},
									"custom_pattern": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specify common types: Select the domain name that requires the cache to be all files or the home page.",
									},
									"file_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "File Type: Specify the file type for cache settings.",
									},
									"custom_file_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Custom file type: Fill in the appropriate identifiable file type according to your needs outside of the specified file type.",
									},
									"specify_url_pattern": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specify URL cache: Specify url according to requirements for cache. INS format does not support URI format with http(s)://",
									},
									"directory": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Directory: Specify the directory cache. Enter a legal directory format. Multiple separated by semicolons",
									},
									"cache_ttl": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cache time: set the time corresponding to the cache object. Input format: integer plus unit, such as 20s, 30m, 1h, 2d.",
									},
									"ignore_cache_control": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Ignore the source station does not cache the header. The optional values are true and false.",
									},
									"is_respect_server": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Respect the server: Accelerate whether to prioritize the source cache time. Optional values: true and false.",
									},
									"ignore_letter_case": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Ignore case, the optional value is true or false, true means to ignore case; false means not to ignore case.",
									},
									"reload_manage": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Reload processing rules, optional: ignore or if-modified-since. If-modified-since: indicates that you want to convert to if-modified-since. Ignore: means to ignore client refresh.",
									},
									"priority": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Indicates the priority execution order of multiple sets of redirected content by the customer. The higher the number, the higher the priority.",
									},
									"ignore_authentication_header": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "You can set it 'true' to cache ignoring the http header 'Authentication'.",
									},
								},
							},
						},
						"cache_key_rules": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"path_pattern": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The url matching mode supports fuzzy regularization. If all matches, the input parameters can be configured as: *",
									},
									"specify_url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specify a uri, such as /test/specifyurl",
									},
									"full_match4_specify_url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Whether to match specifyUrl exactly or not, you can select true and false. True: means match exactly. False: means fuzzy match, such as specifyUrl='/test/uri', and request for /test/uri?p=1 will be matched.",
									},
									"custom_pattern": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specify common types: Select the domain name that requires the cache to be all files or the home page.",
									},
									"file_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "File Type: Specify the file type for cache settings.",
									},
									"custom_file_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Custom file type: Fill in the appropriate identifiable file type according to your needs outside of the specified file type. Can be used with file-type. If the file-type is also configured, the actual file type is the sum of the two parameters.",
									},
									"directory": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Directory: Specify the directory cache. Enter a legal directory format. Multiple separated by semicolons.",
									},
									"ignore_case": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Ignore case, the optional value is true or false, true means to ignore case; false means not to ignore case.",
									},
									"header_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Header name. Example: If you specify a header as 'lang', Then, if the value of Lang is consistent, one copy will be cached.",
									},
									"parameter_of_header": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Parameter Of the specified Header. Example: Specifies the header as 'cookie', parameterOfHeader as 'name'. Then, if the value of name is consistent, one copy will be cached.",
									},
									"priority": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Indicates the priority execution order of multiple sets of redirected content by the customer. The higher the number, the higher the priority.",
									},
									"data_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "DataId is to indicate a specific group configuration when the client has multiple groups of configurations. dataId can be retrieved through a query interface.",
									},
								},
							},
							Description: "Custom Cachekey Configuration, parent node.",
						},
						"query_string_settings": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"path_pattern": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The url matching mode. If all matches, the input parameters can be configured as: .*",
									},
									"except_path_pattern": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Exceptional url matching mode, except for some URLs: such as abc.jpg, do not do anti-theft chain function. E.g: ^https?://[^/]+/.*\\.m3u8",
									},
									"file_types": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "File Type: Specify the file type for anti-theft chain settings.",
									},
									"custom_file_types": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Custom file type: Fill in the appropriate identifiable file type according to your needs outside of the specified file type.",
									},
									"custom_pattern": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specify common types: Select the domain name that requires the anti-theft chain to be all files or the home page.",
									},
									"specify_url_pattern": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specify URL cache: Specify url according to requirements for anti-theft chain setting.",
									},
									"directories": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Directory: Specify the directory for anti-theft chain settings.",
									},
									"priority": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Indicates the priority execution order of multiple sets of redirected content by the customer. The higher the number, the higher the priority.",
									},
									"ignore_letter_case": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Whether to ignore letter case.",
									},
									"ignore_query_string": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Define whether to cache with query strings, 'true' means ignore query strings, while 'false' means cache with all query strings.",
									},
									"query_string_kept": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cache with the specified query string parameters. If the kept parameter values are the same, one copy will be cached.",
									},
									"query_string_removed": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cache without the specified query string parameters. After deleting the specified parameter, if the other parameter values are the same, one copy will be cached.",
									},
									"source_with_query": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Whether to use the original URL back source, the allowable values are true and false.",
									},
									"source_key_kept": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Return to the source after specifying the reserved parameter value. Please separate them with semicolons, if no parameters reserved, please fill in:-.",
									},
									"source_key_removed": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Return to the source after specifying the deleted parameter value. Please separate them with semicolons, and if you do not delete any parameters, please fill in:-.",
									},
									"data_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Data-id is to indicate a specific group configuration when the client has multiple groups of configurations. Data-id can be retrieved through a query interface.",
									},
								},
							},
							Description: "Query String Settings Configuration, parent node.",
						},
						"cache_by_resp_headers": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Cache the file according to the response header content",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"data_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Add grid type identity, represents the customer multi - group configuration, a specific group of configuration",
									},
									"response_header": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the response header",
									},
									"path_pattern": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Url matching pattern, support fuzzy regular, if all match, the parameter can be configured as: *",
									},
									"except_path_pattern": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Exception to url matching pattern, except for some urls: abc.jpg, no content redirection",
									},
									"response_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Header values.",
									},
									"ignore_letter_case": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Ignore case, the optional value is true or false, true means ignore case;False means that case is not ignored.",
									},
									"priority": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Represents the priority execution order for multiple groups of customer redirected content.The larger the number, the higher the priority.",
									},
									"is_respheader": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Whether the response header content cache file is allowed. An optional value of true or false indicates that the response header content cache file is allowed.False indicates that the response header content cache file is not allowed.",
									},
								},
							},
						},
						"http_code_cache_rules": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Status Code Caching Rule Configuration, parent node 1. When you need to set status code caching rules, this must be filled in. 2. Configuration of Clear Status Code Caching Rules for <http-code-cache-rules/>.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"http_codes": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Description: "Configure HTTP status code, parent node",
									},
									"cache_ttl": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Define the caching time of the specified status code in units s, 0 to indicate no caching",
									},
									"data_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Data-id is to indicate a specific group configuration when the client has multiple groups of configurations. Data-id can be retrieved through a query interface. Note: A. If data-id is passed, it means that one group of configuration items is specified to be modified, and no other group configuration items need to be modified. B. If multiple groups of configurations are included, some of them are configured with data-id and others are not, then the expression of data-id is used to modify a specific group of configurations, and a new group of configurations is added on the original basis without the expression of data-id. C. If the data-id is not transmitted, it means that the original configuration will be fully covered by this configuration. D. If no configuration parameter is passed, only domain name and secondary label are passed, which means that all configuration of domain name secondary service corresponding to this interface is cleared. E. If there is no specific configuration item in a set of configurations, the data-id must be filled in, and the value is the actual data-id, which means clearing the value of the corresponding data-id configuration item; it is not allowed that there is no specific configuration item or data-id in a set of configurations.",
									},
								},
							},
						},
						"ignore_protocol_rules": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Ignore protocol caching and push configuration, parent tags 1. This must be filled when protocol cache and push configuration need to be ignored 2.<ignore-protocol-rules/>:Clear the configuration ignore about protocol cache and pushing",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"path_pattern": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Url matching pattern, support regular, if all matches, input parameters can be configured as:.*",
									},
									"except_path_pattern": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The exception url matches the pattern in the same format as the path-pattern",
									},
									"cache_ignore_protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Whether to ignore the protocol cache, with allowable values of true and false. True turns on the HTTP/HTTPS Shared cache. Not on by default.",
									},
									"purge_ignore_protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "It is recommended to use with cache-ignore protocol to avoid push failure. Note: 1. Once configured, the global effect is not applied to the matched path-pattern. 2. Directory push does not distinguish protocols, while url push can distinguish protocols",
									},
									"data_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "When configuring multiple groups of configurations, specify the id of a particular group of configurations",
									},
								},
							},
						},
						"http2_settings": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Http2.0 settings, used to enable or disable http2.0, parent node.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enable_http2": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Enable http2.0. The optional values are true and false. If it is empty, the default value is false. True means http2.0 is on; false means http2.0 is off.",
									},
									"back_to_origin_protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Back-to-origin protocol, the optional value is http1.1: Use the HTTP1.1 protocol version to back to source. if not filled, use it as default. follow-request: Same as client request protocol http2.0: Use the HTTP2.0 protocol. version to back to source.",
									},
								},
							},
						},
						"header_modify_rules": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Http header settings note: 1. When you need to cancel the http header setting, you can pass in the empty node <header-modify-rules></header-modify-rules>. 2. indicating that you need to set the http header, this field is required",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"data_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Add a grid type identifier to indicate a specific group configuration when the client has multiple groups of configurations.",
									},
									"path_pattern": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The url matching mode supports fuzzy regularization. If all matches, the input parameters can be configured as: *",
									},
									"except_path_pattern": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Exception url matching pattern, support regular. Example: ",
									},
									"custom_pattern": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Matching conditions: specify common types, optional values are all or homepage. 1. all: all files 2. homepage: home page",
									},
									"file_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Matching conditions: file type, please separate by semicolon, optional values: gif png bmp jpeg jpg html htm shtml mp3 wma flv mp4 wmv zip exe rar css txt ico js swf m3u8 xml f4m bootstarp ts.",
									},
									"custom_file_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Matching condition: Custom file type, separate by semicolon.",
									},
									"directory": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Directory",
									},
									"specify_url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Matching Condition: Specify URL. The input parameter does not support the URI format starting with http(s)://",
									},
									"request_method": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The matching request method, the optional values are: GET, POST, PUT, HEAD, DELETE, OPTIONS, separate by semicolons.",
									},
									"header_direction": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The control direction of the http header, the optional value is cache2visitor/cache2origin/visitor2cache/origin2cache, single-select. Cache2origin refers to the source direction---corresponding to the configuration item return source request; Cache2visitor refers to the direction of the client back - the corresponding configuration item returns to the client response; Visitor2cache refers to receiving client requests Origin2cache refers to the receiving source response",
									},
									"action": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The control type of the http header supports the addition and deletion of the http header value. The optional value is add|set|delete, which is single-selected. Corresponding to the header-name and header-value parameters. 1. Add: add a header 2. Set: modify the header value 3. Delete: delete the header Note: priority is delete > set > add",
									},
									"allow_regexp": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Http header regular match, optional value: true / false. True: indicates that the value of the header-name is handled as a regular match. False: indicates that the value of the header-name is processed according to the actual parameters, and no regular match is made. Do not pass the default is false",
									},
									"header_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Http header name, add or modify the http header, only one is allowed; delete the http header to allow multiple entries, separated by a semicolon ';'. Note: The operation of the special http header is limited, and the http header and operation type of the operation are allowed. This item is required and cannot be empty When the action is add: indicates that the header-name header is added. When the action is set: modify the header-name header When the action is delete: delete the header-name header",
									},
									"header_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The value corresponding to the HTTP header field, for example: mytest.example.com Note: 1. When the action is add or set, the input parameter must be passed a value 2. When the action is delete, the input parameter is not passed Support to get the value of specified variable by keyword, such as client IP, including: Key words: meaning #timestamp: current time, timestamp as 1559124945 #request-host: host in the request header #request-url: request url, which contains the full path of the protocol domain name, etc., such as http://aaa.aa.com/a.html #request-uri: request uri, relative path format, such as /index.html #origin- IP: return source IP #cache-ip: edge node IP #server-ip: external service IP #client-ip: client IP, or visitor IP #response-header{XXX} : get the value in the response header, such as #response-header{etag}, get the etag value in response-header #header{XXX} : to get the value in the HTTP header of the request, such as #header{user-agent}, is to get the user-agent value in the header #cookie{XXX} : get the value in the cookie, such as #cookie{account}, is to get the value of the account set in the cookie",
									},
									"request_header": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Match request header, header values support regular, header and header values separated by Spaces, e.g. : Range bytes=[0-9]{9,}",
									},
									"priority": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Indicates the priority of execution order for multiple sets of configurations. A higher number indicates higher priority. If no parameters are passed, the default value is 10 and cannot be cleared.",
									},
									"except_file_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Exception file type.",
									},
									"except_directory": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Exception directory.",
									},
									"except_request_method": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Exception request method.",
									},
									"except_request_header": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Exception request header.",
									},
								},
							},
						},
						"rewrite_rule_settings": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "redirection function note: 1. Define a set of internal redirected content. If there is internal redirected content, this field is required. 2. need to clear the content redirection content under the domain name, you can pass the empty node <rewrite-rule-settings></rewrite-rule-settings>",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"data_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Add a grid type identifier to indicate a specific group configuration when the client has multiple groups of configurations.",
									},
									"path_pattern": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The url matching mode supports fuzzy regularization. If all matches, the input parameters can be configured as: *",
									},
									"custom_pattern": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Matching conditions: specify common types, optional values are all or homepage 1. all: all files 2. homepage: home page",
									},
									"directory": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "directory",
									},
									"file_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "gif png bmp jpeg jpg html htm shtml mp3 wma flv mp4 wmv zip exe rar css txt ico js swf m3u8 xml f4m bootstarp ts",
									},
									"custom_file_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Matching condition: Custom file type, please separate them by semicolon.",
									},
									"except_path_pattern": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Exceptional url matching mode, except for certain URLs: such as abc.jpg, no content redirection Customer reference: ^https?://[^/]+/.*\\.m3u8",
									},
									"ignore_letter_case": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Ignore case, the optional value is true or false, true means to ignore case; false means not to ignore case; When adding a new configuration item, the default is not true. If the client passes a null value: such as <ignore-letter-case></ignore-letter-case>, the configuration is cleared.",
									},
									"publish_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Rewrite the location where the content is generated. The input value is: Cache indicates the node; Other input formats are not supported at this time",
									},
									"priority": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Indicates the priority execution order of multiple sets of redirected content by the customer. The higher the number, the higher the priority. When adding a new configuration item, the default is 10",
									},
									"before_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration item: old url Indicates the protocol mode before rewriting (that is, the object that needs to be rewritten), such as: ^https://([^/]+/.*)",
									},
									"after_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration item: new url Indicates the protocol method after rewriting, such as: http://$1",
									},
									"rewrite_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Redirection type; support for input: before: before the anti-theft chain after: after the anti-theft chain",
									},
									"request_header": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Matching condition: Request header",
									},
									"exception_request_header": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Matching condition: Exception request header",
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

func dataSourceWangSuCdnDomainDetailRead(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("data_source.wangsu_cdn_domain_detail.read")

	domainName := data.Get("domain_name").(string)
	var response *cdn.QueryDomainForTerraformResponse
	var diags diag.Diagnostics
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseCdnClient().QueryCdnDomain(domainName)
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
	_ = data.Set("message", response.Message)
	var resultList []interface{}
	//读取response.Data的所有属性到map中
	var domainDetail = map[string]interface{}{
		"domain_id":             response.Data.DomainId,
		"cname":                 response.Data.Cname,
		"domain_name":           response.Data.DomainName,
		"service_type":          response.Data.ServiceType,
		"service_areas":         response.Data.ServiceAreas,
		"comment":               response.Data.Comment,
		"header_of_client_ip":   response.Data.HeaderOfClientIp,
		"origin_config":         buildOriginConfig(response.Data.OriginConfig),
		"ssl":                   buildSsl(response.Data.Ssl),
		"cache_time_behaviors":  buildCacheTimeBehaviors(response.Data.CacheTimeBehaviors),
		"cache_key_rules":       buildCacheKeyRules(response.Data.CacheKeyRules),
		"query_string_settings": buildQueryStringSettings(response.Data.QueryStringSettings),
		"cache_by_resp_headers": buildCacheByRespHeaders(response.Data.CacheByRespHeaders),
		"http_code_cache_rules": buildHttpCodeCacheRules(response.Data.HttpCodeCacheRules),
		"ignore_protocol_rules": buildIgnoreProtocolRules(response.Data.IgnoreProtocolRules),
		"http2_settings":        buildHttp2Settings(response.Data.Http2Settings),
		"header_modify_rules":   buildHeaderModifyRules(response.Data.HeaderModifyRules),
		"rewrite_rule_settings": buildRewriteRuleSettings(response.Data.RewriteRuleSettings),
	}
	resultList = append(resultList, domainDetail)

	_ = data.Set("data", resultList)

	data.SetId(*response.Data.DomainName)
	log.Printf("data_source.wangsu_cdn_domain_detail.read success")
	return nil
}

func buildOriginConfig(config *cdn.QueryDomainForTerraformResponseDataOriginConfig) interface{} {
	if config == nil {
		return nil
	}
	var originConfig = map[string]interface{}{
		"origin_ips":                 config.OriginIps,
		"default_origin_host_header": config.DefaultOriginHostHeader,
		"use_range":                  config.UseRange,
		"follow301":                  config.Follow301,
		"follow302":                  config.Follow302,
		"adv_src_setting":            buildAdvSrcSetting(config.AdvSrcSetting),
	}
	return []interface{}{originConfig}
}

func buildAdvSrcSetting(setting *cdn.QueryDomainForTerraformResponseDataOriginConfigAdvSrcSetting) interface{} {
	if setting == nil {
		return nil
	}
	var advSrcSetting = map[string]interface{}{
		"use_adv_src":   setting.UseAdvSrc,
		"detect_url":    setting.DetectUrl,
		"detect_period": setting.DetectPeriod,
		"master_ips":    setting.MasterIps,
		"backup_ips":    setting.BackupIps,
	}
	return []interface{}{advSrcSetting}
}

func buildSsl(ssl *cdn.QueryDomainForTerraformResponseDataSsl) interface{} {
	if ssl == nil {
		return nil
	}
	var sslConfig = map[string]interface{}{
		"use_ssl":               ssl.UseSsl,
		"ssl_certificate_id":    ssl.SslCertificateId,
		"backup_certificate_id": ssl.BackupCertificateId,
		"gm_certificate_ids":    ssl.GmCertificateIds,
		"tls_version":           ssl.TlsVersion,
		"enable_ocsp":           ssl.EnableOcsp,
		"ssl_cipher_suite":      ssl.SslCipherSuite,
	}
	return []interface{}{sslConfig}
}
func buildCacheTimeBehaviors(behaviors []*cdn.QueryDomainForTerraformResponseDataCacheTimeBehaviors) interface{} {
	if behaviors == nil {
		return nil
	}
	var cacheTimeBehaviors []interface{}
	for _, behavior := range behaviors {
		var cacheTimeBehavior = map[string]interface{}{
			"data_id":                      behavior.DataId,
			"path_pattern":                 behavior.PathPattern,
			"except_path_pattern":          behavior.ExceptPathPattern,
			"custom_pattern":               behavior.CustomPattern,
			"file_type":                    behavior.FileType,
			"custom_file_type":             behavior.CustomFileType,
			"specify_url_pattern":          behavior.SpecifyUrlPattern,
			"directory":                    behavior.Directory,
			"cache_ttl":                    behavior.CacheTtl,
			"ignore_cache_control":         behavior.IgnoreCacheControl,
			"is_respect_server":            behavior.IsRespectServer,
			"ignore_letter_case":           behavior.IgnoreLetterCase,
			"reload_manage":                behavior.ReloadManage,
			"priority":                     behavior.Priority,
			"ignore_authentication_header": behavior.IgnoreAuthenticationHeader,
		}
		cacheTimeBehaviors = append(cacheTimeBehaviors, cacheTimeBehavior)
	}
	return cacheTimeBehaviors
}

func buildCacheKeyRules(rules []*cdn.QueryDomainForTerraformResponseDataCacheKeyRules) interface{} {
	if rules == nil {
		return nil
	}
	var cacheKeyRules []interface{}
	for _, rule := range rules {
		var cacheKeyRule = map[string]interface{}{
			"path_pattern":            rule.PathPattern,
			"specify_url":             rule.SpecifyUrl,
			"full_match4_specify_url": rule.FullMatch4SpecifyUrl,
			"custom_pattern":          rule.CustomPattern,
			"file_type":               rule.FileType,
			"custom_file_type":        rule.CustomFileType,
			"directory":               rule.Directory,
			"ignore_case":             rule.IgnoreCase,
			"header_name":             rule.HeaderName,
			"parameter_of_header":     rule.ParameterOfHeader,
			"priority":                rule.Priority,
			"data_id":                 rule.DataId,
		}
		cacheKeyRules = append(cacheKeyRules, cacheKeyRule)
	}
	return cacheKeyRules
}

func buildQueryStringSettings(settings []*cdn.QueryDomainForTerraformResponseDataQueryStringSettings) interface{} {
	if settings == nil {
		return nil
	}
	var queryStringSettings []interface{}
	for _, setting := range settings {
		var queryStringSetting = map[string]interface{}{
			"path_pattern":         setting.PathPattern,
			"except_path_pattern":  setting.ExceptPathPattern,
			"file_types":           setting.FileTypes,
			"custom_file_types":    setting.CustomFileTypes,
			"custom_pattern":       setting.CustomPattern,
			"specify_url_pattern":  setting.SpecifyUrlPattern,
			"directories":          setting.Directories,
			"priority":             setting.Priority,
			"ignore_letter_case":   setting.IgnoreLetterCase,
			"ignore_query_string":  setting.IgnoreQueryString,
			"query_string_kept":    setting.QueryStringKept,
			"query_string_removed": setting.QueryStringRemoved,
			"source_with_query":    setting.SourceWithQuery,
			"source_key_kept":      setting.SourceKeyKept,
			"source_key_removed":   setting.SourceKeyRemoved,
			"data_id":              setting.DataId,
		}
		queryStringSettings = append(queryStringSettings, queryStringSetting)
	}
	return queryStringSettings
}

func buildCacheByRespHeaders(headers []*cdn.QueryDomainForTerraformResponseDataCacheByRespHeaders) interface{} {
	if headers == nil {
		return nil
	}
	var cacheByRespHeaders []interface{}
	for _, header := range headers {
		var cacheByRespHeader = map[string]interface{}{
			"data_id":             header.DataId,
			"response_header":     header.ResponseHeader,
			"path_pattern":        header.PathPattern,
			"except_path_pattern": header.ExceptPathPattern,
			"response_value":      header.ResponseValue,
			"ignore_letter_case":  header.IgnoreLetterCase,
			"priority":            header.Priority,
			"is_respheader":       header.IsRespheader,
		}
		cacheByRespHeaders = append(cacheByRespHeaders, cacheByRespHeader)
	}
	return cacheByRespHeaders
}

func buildHttpCodeCacheRules(rules []*cdn.QueryDomainForTerraformResponseDataHttpCodeCacheRules) interface{} {
	if rules == nil {
		return nil
	}
	var httpCodeCacheRules []interface{}
	for _, rule := range rules {
		var httpCodeCacheRule = map[string]interface{}{
			"http_codes": rule.HttpCodes,
			"cache_ttl":  rule.CacheTtl,
			"data_id":    rule.DataId,
		}
		httpCodeCacheRules = append(httpCodeCacheRules, httpCodeCacheRule)
	}
	return httpCodeCacheRules
}

func buildIgnoreProtocolRules(rules []*cdn.QueryDomainForTerraformResponseDataIgnoreProtocolRules) interface{} {
	if rules == nil {
		return nil
	}
	var ignoreProtocolRules []interface{}
	for _, rule := range rules {
		var ignoreProtocolRule = map[string]interface{}{
			"path_pattern":          rule.PathPattern,
			"except_path_pattern":   rule.ExceptPathPattern,
			"cache_ignore_protocol": rule.CacheIgnoreProtocol,
			"purge_ignore_protocol": rule.PurgeIgnoreProtocol,
			"data_id":               rule.DataId,
		}
		ignoreProtocolRules = append(ignoreProtocolRules, ignoreProtocolRule)
	}
	return ignoreProtocolRules
}

func buildHttp2Settings(settings *cdn.QueryDomainForTerraformResponseDataHttp2Settings) interface{} {
	if settings == nil {
		return nil
	}
	var http2Settings = map[string]interface{}{
		"enable_http2":            settings.EnableHttp2,
		"back_to_origin_protocol": settings.BackToOriginProtocol,
	}
	return []interface{}{http2Settings}
}

func buildHeaderModifyRules(rules []*cdn.QueryDomainForTerraformResponseDataHeaderModifyRules) interface{} {
	if rules == nil {
		return nil
	}
	var headerModifyRules []interface{}
	for _, rule := range rules {
		var headerModifyRule = map[string]interface{}{
			"data_id":               rule.DataId,
			"path_pattern":          rule.PathPattern,
			"except_path_pattern":   rule.ExceptPathPattern,
			"custom_pattern":        rule.CustomPattern,
			"file_type":             rule.FileType,
			"custom_file_type":      rule.CustomFileType,
			"directory":             rule.Directory,
			"specify_url":           rule.SpecifyUrl,
			"request_method":        rule.RequestMethod,
			"header_direction":      rule.HeaderDirection,
			"action":                rule.Action,
			"allow_regexp":          rule.AllowRegexp,
			"header_name":           rule.HeaderName,
			"header_value":          rule.HeaderValue,
			"request_header":        rule.RequestHeader,
			"priority":              rule.Priority,
			"except_file_type":      rule.ExceptFileType,
			"except_directory":      rule.ExceptDirectory,
			"except_request_method": rule.ExceptRequestMethod,
			"except_request_header": rule.ExceptRequestHeader,
		}
		headerModifyRules = append(headerModifyRules, headerModifyRule)
	}
	return headerModifyRules
}

func buildRewriteRuleSettings(settings []*cdn.QueryDomainForTerraformResponseDataRewriteRuleSettings) interface{} {
	if settings == nil {
		return nil
	}
	var rewriteRuleSettings []interface{}
	for _, setting := range settings {
		var rewriteRuleSetting = map[string]interface{}{
			"data_id":                  setting.DataId,
			"path_pattern":             setting.PathPattern,
			"custom_pattern":           setting.CustomPattern,
			"directory":                setting.Directory,
			"file_type":                setting.FileType,
			"custom_file_type":         setting.CustomFileType,
			"except_path_pattern":      setting.ExceptPathPattern,
			"ignore_letter_case":       setting.IgnoreLetterCase,
			"publish_type":             setting.PublishType,
			"priority":                 setting.Priority,
			"before_value":             setting.BeforeValue,
			"after_value":              setting.AfterValue,
			"rewrite_type":             setting.RewriteType,
			"request_header":           setting.RequestHeader,
			"exception_request_header": setting.ExceptionRequestHeader,
		}
		rewriteRuleSettings = append(rewriteRuleSettings, rewriteRuleSetting)
	}
	return rewriteRuleSettings
}
