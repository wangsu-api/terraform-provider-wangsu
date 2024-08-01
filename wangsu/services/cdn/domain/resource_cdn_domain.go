package domain

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wangsuCommon "github.com/wangsu-api/terraform-provider-wangsu/wangsu/common"
	cdn "github.com/wangsu-api/wangsu-sdk-go/wangsu/cdn/domain"
	"log"
	"time"
)

func ResourceCdnDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCdnDomainCreate,
		ReadContext:   resourceCdnDomainRead,
		UpdateContext: resourceCdnDomainUpdate,
		DeleteContext: resourceCdnDomainDelete,

		Schema: map[string]*schema.Schema{
			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Need to access the domain name of the CDN. a generic domain name is supported, starting with the symbol '.'. such as.example.com. which also contains a multilevel 'a.b.example.com'.If example.com is filed. the domain name xx.example.com does not need to be filed.",
			},
			"service_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The service type of the accelerated domain name (only one service type can be submitted at a time): web/web-https: Web page acceleration/Web page acceleration-https wsa/Wsa-https: Full-station acceleration/full-station acceleration-https vodstream/vod-https: on-demand acceleration/on-demand acceleration-https download/dl-https: Download Acceleration/Download Acceleration-https livestream/live-https/cloudv-live: livestream acceleration v6sa/osv6: IPv6 Security&Acceleration Solution/IPv6 One-stop Solution Note: 1. the https in the code. such as web-https does not represent immediate support for https access. you need to upload the certificate to support https.",
			},
			"service_areas": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The acceleration area of the acceleration domain. if the resource coverage needs to be limited according to the area. the acceleration area needs to be specified. When no acceleration area is specified. we will provide acceleration services with optimal resource coverage according to the service area opened by the customer. Multiple regions are separated by semicolons. and the supported regions are as follows: cn (Mainland China). am (Americas). emea (Europe. Middle East. Africa). apac (Asia-Pacific region).",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Remarks. up to 1000 characters",
			},
			"header_of_client_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Pass the response header of client IP. The optional values are Cdn-Src-Ip and X-Forwarded-For. The default value is Cdn-Src-Ip.",
			},
			"origin_config": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"origin_ips": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Origin address. which can be an IP or domain name. 1. Multiple IPs are supported. separated by semicolons. 2. Only one domain name is allowed. IP and domain name cannot exist at the same time. 3. The length cannot exceed 500 characters. 4. The number of IPs cannot exceed 15.",
						},
						"default_origin_host_header": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Back-to-origin HOST. used to change the HOST field in the back-to-origin HTTP request header. The supported formats are: ① domain name ③ ip Note: 1. Must comply with the ip/domain name format specification. If it is a domain name. the length of the domain name must be less than or equal to 128 characters.",
						},
						"use_range": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "useRange",
						},
						"follow301": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "follow301",
						},
						"follow302": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "follow302",
						},
						"adv_src_setting": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"use_adv_src": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Use advance origin config. the optional values are true and false. true means to use advance origin config. false means not to use advance origin config",
									},
									"detect_url": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The advanced source monitors the url. and requests <master-ips> through the url. If the response is not 2**. 3** response. it is considered that the primary source ip is faulty. and <backup-ips> is used at this time.",
									},
									"detect_period": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Advanced source monitoring period. in seconds. optional as an integer greater than or equal to 0. 0 means no monitoring",
									},
									"master_ips": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "The advanced source mainly returns the source IP. Multiple IPs are separated by a semicolon \";\". and the return source IP cannot be repeated.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"backup_ips": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Advanced source backup source IP. multiple IPs are separated by semicolon \";\". and the return source IP cannot be duplicated.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},
			"ssl": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "SSL settings, to bind a certificate with the accelerated domain. You can use the interface [AddCertificate] to upload your  certificates. If you want to modify a certificate, please use the interface: [UpdateCertificate]",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"use_ssl": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Use a certificate, the optional values are true and false, true means to use the certificate, false means not to use the certificate",
						},
						"ssl_certificate_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Use sni certificate, the optional values are true and false, true means use sni certificate, false means use shared certificate (not supported)",
						},
						"backup_certificate_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Backup certificate ID",
						},
						"gm_certificate_ids": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "SM2 certificate IDS",
						},
						"tls_version": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "TLS version. Optional values: SSLv3,TLSv1,TLSv1.1,TLSv1.2,TLSv1.3",
						},
						"enable_ocsp": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Enable OCSP(Online Certificate Status Protocol).",
						},
						"ssl_cipher_suite": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "This optional object is used to specify a colon separated list of cipher suites which are permitted when clients negotiate security settings to access your content. Cipher suites which you can specify are: LOW, ALL:!LOW, HIGH, !EXPORT, !aNULL, !RC4, !DH, !SHA, !MD5, @STRENGTH,  AES128-SHA, AES256-SHA, AES128-SHA256, AES256-SHA256, AES128-GCM-SHA256, AES256-GCM-SHA384, ECDHE-RSA-AES128-SHA, ECDHE-RSA-AES256-SHA, ECDHE-RSA-AES128-SHA256, ECDHE-RSA-AES256-SHA384, ECDHE-RSA-AES128-GCM-SHA256, and ECDHE-RSA-AES256-GCM-SHA384. These cipher suites are a subset of those supported by OpenSSL, https://www.openssl.org/docs/man1.0.2/man1/ciphers.html. Please note that !MD5 or !SHA must appear after HIGH..",
						},
					},
				},
			},
			"cache_time_behaviors": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Cache time configuration note: 1. When you need to cancel the cache time configuration setting, you can pass in the empty node <cache-time-behaviors></cache-time-behaviors>. 2. When it is required to set the cache time configuration, this item is required.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path_pattern": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The url matching mode supports fuzzy regularization. If all matches, the input parameters can be configured as: *",
						},
						"except_path_pattern": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Exceptional url matching mode, except for some URLs: such as abc.jpg, do not do anti-theft chain function E.g: ^https?://[^/]+/.*\\.m3u8",
						},
						"custom_pattern": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specify common types: Select the domain name that requires the cache  to be all files or the home page. : E.g: All: all files Homepage: homepage",
						},
						"file_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "File Type: Specify the file type for cache settings. File types include: gif png bmp jpeg jpg html htm shtml mp3 wma flv mp4 wmv zip exe rar css txt ico js swf If you need all types, pass all directly. Multiples are separated by semicolons, and all and specific file types cannot be configured at the same time.",
						},
						"custom_file_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Custom file type: Fill in the appropriate identifiable file type according to your needs outside of the specified file type. Can be used with file-type. If the file-type is also configured, the actual file type is the sum of the two parameters.",
						},
						"specify_url_pattern": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specify URL cache: Specify url according to requirements for cache INS format does not support URI format with http(s)://",
						},
						"directory": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Directory: Specify the directory cache. Enter a legal directory format. Multiple separated by semicolons",
						},
						"cache_ttl": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Cache time: set the time corresponding to the cache object Input format: integer plus unit, such as 20s, 30m, 1h, 2d, no cache is set to 0. Do not enter the unit default is seconds There is no upper limit on the cache time theory. This time is set according to the customer's own needs. If the customer feels that some of the files are not changed frequently, then the setting is longer. For example, the text class js, css, html, etc. can be set shorter, the picture, video and audio classes can be set longer (because the cache time will be replaced by the new file due to the file heat algorithm, the longest suggestion Do not exceed one month)",
						},
						"ignore_cache_control": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Ignore the source station does not cache the header. The optional values are true and false, which are used to ignore the two configurations of cache-control in the request header (private, no-cache) and the Authorization set by the client. The ture indicates that the source station's settings for the three are ignored. Enables resources to be cached on the service node in the form of cache-control: public, and then our nodes can cache this type of resource and provide acceleration services. False means that when the source station sets cache-control: private, cache-control: no-cache for a resource or specifies to cache according to authorization, our service node will not cache such files.",
						},
						"is_respect_server": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Respect the server: Accelerate whether to prioritize the source cache time. Optional values: true and false True: indicates that the server is time-first False: The cache time of the CDN configuration takes precedence.",
						},
						"ignore_letter_case": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Ignore case, the optional value is true or false, true means to ignore case; false means not to ignore case; When adding a new configuration item, the default is not true.",
						},
						"reload_manage": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Reload processing rules, optional: ignore or if-modified-since If-modified-since: indicates that you want to convert to if-modified-since Ignore: means to ignore client refresh",
						},
						"priority": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Indicates the priority execution order of multiple sets of redirected content by the customer. The higher the number, the higher the priority. When adding a new configuration item, the default is 10",
						},
						"ignore_authentication_header": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "You can set it 'true' to cache ignoring the http header 'Authentication'.  If it is empty, the header is not ignored by default.",
						},
					},
				},
			},
			"cache_key_rules": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Custom Cachekey Configuration, parent node 1. When you need to configure the cachekey rules,this must be filled in. 2. Configuration of clearing for <cacheKeyRules/>.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path_pattern": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The url matching mode supports fuzzy regularization. If all matches, the input parameters can be configured as: *",
						},
						"specify_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specify a uri, such as /test/specifyurl",
						},
						"full_match4_specify_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Whether to match specifyUrl exactly or not, you can select true and false. True:means match exactly. False: means fuzzy match, such as specifyUrl='/test/uri', and  request for /test/uri?p=1 will be matched.",
						},
						"custom_pattern": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specify common types: Select the domain name that requires the cache to be all files or the home page. : E.g: All: all files Homepage: homepage",
						},
						"file_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "File Type: Specify the file type for cache settings. File types include: gif png bmp jpeg jpg html htm shtml mp3 wma flv mp4 wmv zip exe rar css txt ico js swf If you need all types, pass all directly. Multiples are separated by semicolons, and all and specific file types cannot be configured at the same time.",
						},
						"custom_file_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Custom file type: Fill in the appropriate identifiable file type according to your needs outside of the specified file type. Can be used with file-type. If the file-type is also configured, the actual file type is the sum of the two parameters.",
						},
						"directory": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Directory: Specify the directory cache. Enter a legal directory format. Multiple separated by semicolons",
						},
						"ignore_case": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Ignore case, the optional value is true or false, true means to ignore case; false means not to ignore case; When adding a new configuration item, the default is true",
						},
						"header_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Header name. Example: If you specify a header as 'lang', Then, if the value of Lang is consistent, one copy will be cached",
						},
						"parameter_of_header": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Parameter Of the specified Header, Example: Specifies the header as 'cookie', parameterOfHeader as 'name'. Then, if the value of name is consistent, one copy will be cached.",
						},
						"priority": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Indicates the priority execution order of multiple sets of redirected content by the customer. The higher the number, the higher the priority. When adding a new configuration item, the default is 10",
						},
					},
				},
			},
			"query_string_settings": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Query String Settings Configuration, parent node\n1. When you need to configure the query string, this must be filled in.\n2. Configuration of clearing query string settings for <query-string-settings/>.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path_pattern": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The url matching mode. If all matches, the input parameters can be configured as: .*",
						},
						"except_path_pattern": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Exceptional url matching mode, except for some URLs: such as abc.jpg, do not do anti-theft chain function E.g: ^https?://[^/]+/.*\\.m3u8",
						},
						"file_types": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "File Type: Specify the file type for anti-theft chain settings.\nFile types include: gif png bmp jpeg jpg html htm shtml mp3 wma flv mp4 wmv zip exe rar css txt ico js swf\nIf you need all types, pass all directly. Multiples are separated by semicolons, and all and specific file types cannot be configured at the same time.",
						},
						"custom_file_types": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Custom file type: Fill in the appropriate identifiable file type according to your needs outside of the specified file type. Can be used with file-type. If the file-type is also configured, the actual file type is the sum of the two parameters.",
						},
						"custom_pattern": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specify common types: Select the domain name that requires the anti-theft chain to be all files or the home page. :\nE.g:\nAll: all files\nHomepage: homepage",
						},
						"specify_url_pattern": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specify URL cache: Specify url according to requirements for anti-theft chain setting\nINS format does not support URI format with http(s)://",
						},
						"directories": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Directory: Specify the directory for anti-theft chain settings\nEnter a legal directory format. Multiple separated by semicolons",
						},
						"priority": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Indicates the priority execution order of multiple sets of redirected content by the customer. The higher the number, the higher the priority.\nWhen adding a new configuration item, the default is 10",
						},
						"ignore_letter_case": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Whether to ignore letter case.",
						},
						"ignore_query_string": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Define whether to cache with query strings, 'true' means ignore query strings, while 'false' means cache with all query strings.",
						},
						"query_string_kept": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Cache with the specified query string parameters. If the kept parameter values are the same, one copy will be cached.\nNote:\n1. query-string-kept and query-string-removed are mutually exclusive, and only one of them has a value.\n2. query-string-kept and ignore-query-string are mutually exclusive, and only one has a value.",
						},
						"query_string_removed": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Cache without the specified query string parameters. After deleting the specified parameter, if the other parameter values are the same, one copy will be cached.\n1. query-string-kept and query string removed are mutually exclusive, and only one has a value.\n2. query-string-removed and ignore-query-string are mutually exclusive.",
						},
						"source_with_query": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Whether to use the original URL back source, the allowable values are true and false.\nWhen ignore-query-string is true or not set, source-with-query is true to indicate that the source is returned according to the original request, and false to indicate that the question mark is returned.\nWhen ignore-query-string is false, this default setting is empty (input is invalid)",
						},
						"source_key_kept": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Return to the source after specifying the reserved parameter value. Please separate them with semicolons, if no parameters reserved, please fill in:- . 1. Source-key-kept and ignore-query-string are mutually exclusive, and only one of them has a value. 2. Source-key-kept and source-key-removed are mutually exclusive, and only one of them has a value.",
						},
						"source_key_removed": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Return to the source after specifying the deleted parameter value. Please separate them with semicolons, and if you do not delete any parameters, please fill in:- . 1. Source-key-removed and ignore-query-string are mutually exclusive, and only one of them has a value. 2. Source-key-kept and source-key-removed are mutually exclusive, and only one of them has a value.",
						},
					},
				},
			},
			"cache_by_resp_headers": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Cache the file according to the response header content",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"response_header": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the response header",
						},
						"path_pattern": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Url matching pattern, support fuzzy regular, if all match, the parameter can be configured as: *",
						},
						"except_path_pattern": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Exception to url matching pattern, except for some urls: abc.jpg, no content redirection",
						},
						"response_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Header values.",
						},
						"ignore_letter_case": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Ignore case, the optional value is true or false, true means ignore case;False means that case is not ignored.",
						},
						"priority": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Represents the priority execution order for multiple groups of customer redirected content.The larger the number, the higher the priority.",
						},
						"is_respheader": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Whether the response header content cache file is allowed. An optional value of true or false indicates that the response header content cache file is allowed.False indicates that the response header content cache file is not allowed.",
						},
					},
				},
			},
			"http_code_cache_rules": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Status Code Caching Rule Configuration, parent node",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"http_codes": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Configure HTTP status code, parent node",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"cache_ttl": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Define the caching time of the specified status code in units s, 0 to indicate no caching",
						},
					},
				},
			},
			"ignore_protocol_rules": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Ignore protocol caching and push configuration, parent tags",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path_pattern": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Url matching pattern, support regular, if all matches, input parameters can be configured as:.*",
						},
						"except_path_pattern": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The exception url matches the pattern in the same format as the path-pattern",
						},
						"cache_ignore_protocol": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Whether to ignore the protocol cache, with allowable values of true and false. True turns on the HTTP/HTTPS Shared cache. Not on by default.",
						},
						"purge_ignore_protocol": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "It is recommended to use with cache-ignore protocol to avoid push failure.\n\nNote:\n\n1. Once configured, the global effect is not applied to the matched path-pattern.\n\n2. Directory push does not distinguish protocols, while url push can distinguish protocols",
						},
					},
				},
			},
			"http2_settings": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Http2.0 settings, used to enable or disable http2.0, parent node.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_http2": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Enable http2.0. The optional values are true and false. If it is empty, the default value is false. True means http2.0 is on; false means http2.0 is off.",
						},
						"back_to_origin_protocol": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Back-to-origin protocol, the optional value is http1.1: Use the HTTP1.1 protocol version to back to source. if not filled, use it as default. follow-request: Same as client request protocol http2.0: Use the HTTP2.0 protocol. version to back to source.",
						},
					},
				},
			},
			"header_modify_rules": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Http header settings note: 1. When you need to cancel the http header setting, you can pass in the empty node <header-modify-rules></header-modify-rules>. 2. indicating that you need to set the http header, this field is required",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path_pattern": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The url matching mode supports fuzzy regularization. If all matches, the input parameters can be configured as: *",
						},
						"except_path_pattern": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Exception url matching pattern, support regular. Example: ",
						},
						"custom_pattern": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Matching conditions: specify common types, optional values are all or homepage. 1. all: all files 2. homepage: home page",
						},
						"file_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Matching conditions: file type, please separate by semicolon, optional values: gif png bmp jpeg jpg html htm shtml mp3 wma flv mp4 wmv zip exe rar css txt ico js swf m3u8 xml f4m bootstarp ts.",
						},
						"custom_file_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Matching condition: Custom file type, separate by semicolon.",
						},
						"directory": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Directory",
						},
						"specify_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Matching Condition: Specify URL. The input parameter does not support the URI format starting with http(s)://",
						},
						"request_method": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The matching request method, the optional values are: GET, POST, PUT, HEAD, DELETE, OPTIONS, separate by semicolons.",
						},
						"header_direction": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The control direction of the http header, the optional value is cache2visitor/cache2origin/visitor2cache/origin2cache, single-select. Cache2origin refers to the source direction---corresponding to the configuration item return source request; Cache2visitor refers to the direction of the client back - the corresponding configuration item returns to the client response; Visitor2cache refers to receiving client requests Origin2cache refers to the receiving source response",
						},
						"action": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The control type of the http header supports the addition and deletion of the http header value. The optional value is add|set|delete, which is single-selected. Corresponding to the header-name and header-value parameters. 1. Add: add a header 2. Set: modify the header value 3. Delete: delete the header Note: priority is delete > set > add",
						},
						"allow_regexp": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Http header regular match, optional value: true / false. True: indicates that the value of the header-name is handled as a regular match. False: indicates that the value of the header-name is processed according to the actual parameters, and no regular match is made. Do not pass the default is false",
						},
						"header_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Http header name, add or modify the http header, only one is allowed; delete the http header to allow multiple entries, separated by a semicolon ';'. Note: The operation of the special http header is limited, and the http header and operation type of the operation are allowed. This item is required and cannot be empty When the action is add: indicates that the header-name header is added. When the action is set: modify the header-name header When the action is delete: delete the header-name header",
						},
						"header_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The value corresponding to the HTTP header field, for example: mytest.example.com Note: 1. When the action is add or set, the input parameter must be passed a value 2. When the action is delete, the input parameter is not passed Support to get the value of specified variable by keyword, such as client IP, including: Key words: meaning #timestamp: current time, timestamp as 1559124945 #request-host: host in the request header #request-url: request url, which contains the full path of the protocol domain name, etc., such as http://aaa.aa.com/a.html #request-uri: request uri, relative path format, such as /index.html #origin- IP: return source IP #cache-ip: edge node IP #server-ip: external service IP #client-ip: client IP, or visitor IP #response-header{XXX} : get the value in the response header, such as #response-header{etag}, get the etag value in response-header #header{XXX} : to get the value in the HTTP header of the request, such as #header{user-agent}, is to get the user-agent value in the header #cookie{XXX} : get the value in the cookie, such as #cookie{account}, is to get the value of the account set in the cookie",
						},
						"request_header": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Match request header, header values support regular, header and header values separated by Spaces, e.g. : Range bytes=[0-9]{9,}",
						},
						"priority": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Indicates the priority of execution order for multiple sets of configurations. A higher number indicates higher priority. If no parameters are passed, the default value is 10 and cannot be cleared.",
						},
						"except_file_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Exception file type.",
						},
						"except_directory": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Exception directory.",
						},
						"except_request_method": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Exception request method.",
						},
						"except_request_header": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Exception request header.",
						},
					},
				},
			},
			"rewrite_rule_settings": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "redirection function note: 1. Define a set of internal redirected content. If there is internal redirected content, this field is required. 2. need to clear the content redirection content under the domain name, you can pass the empty node <rewrite-rule-settings></rewrite-rule-settings>",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path_pattern": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The url matching mode supports fuzzy regularization. If all matches, the input parameters can be configured as: *",
						},
						"custom_pattern": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Matching conditions: specify common types, optional values are all or homepage 1. all: all files 2. homepage: home page",
						},
						"directory": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "directory",
						},
						"file_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "gif png bmp jpeg jpg html htm shtml mp3 wma flv mp4 wmv zip exe rar css txt ico js swf m3u8 xml f4m bootstarp ts",
						},
						"custom_file_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Matching condition: Custom file type, please separate them by semicolon.",
						},
						"except_path_pattern": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Exceptional url matching mode, except for certain URLs: such as abc.jpg, no content redirection Customer reference: ^https?://[^/]+/.*\\.m3u8",
						},
						"ignore_letter_case": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Ignore case, the optional value is true or false, true means to ignore case; false means not to ignore case; When adding a new configuration item, the default is not true. If the client passes a null value: such as <ignore-letter-case></ignore-letter-case>, the configuration is cleared.",
						},
						"publish_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Rewrite the location where the content is generated. The input value is: Cache indicates the node; Other input formats are not supported at this time",
						},
						"priority": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Indicates the priority execution order of multiple sets of redirected content by the customer. The higher the number, the higher the priority. When adding a new configuration item, the default is 10",
						},
						"before_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Configuration item: old url Indicates the protocol mode before rewriting (that is, the object that needs to be rewritten), such as: ^https://([^/]+/.*)",
						},
						"after_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Configuration item: new url Indicates the protocol method after rewriting, such as: http://$1",
						},
						"rewrite_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Redirection type; support for input: before: before the anti-theft chain after: after the anti-theft chain",
						},
						"request_header": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Matching condition: Request header",
						},
						"exception_request_header": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Matching condition: Exception request header",
						},
					},
				},
			},
		},
	}
}

func resourceCdnDomainDelete(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_cdn_domain.delete")

	var response *cdn.DeleteDomainForTerraformResponse
	var requestId string
	var err error
	var diags diag.Diagnostics
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		requestId, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseCdnClient().DeleteCdnDomain(data.Id())
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
		return nil
	}

	time.Sleep(3 * time.Second)
	//query domain deployment status
	var deploymentResponse *cdn.QueryDeployResultForTerraformResponse
	err = resource.RetryContext(context, time.Duration(5)*time.Minute, func() *resource.RetryError {
		deploymentResponse, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseCdnClient().QueryDomainDeployStatus(requestId)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if deploymentResponse != nil && deploymentResponse.Data != nil && *deploymentResponse.Data.DeployResult != "SUCCESS" {
			return resource.RetryableError(fmt.Errorf("domain deployment status is in progress, retrying"))
		}
		return nil
	})
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	return nil
}

func resourceCdnDomainRead(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_cdn_domain.read")
	//query domain information
	var response *cdn.QueryDomainForTerraformResponse
	var diags diag.Diagnostics
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseCdnClient().QueryCdnDomain(data.Id())
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
	responseData := response.Data

	_ = data.Set("domain_id", responseData.DomainId)
	_ = data.Set("domain_name", responseData.DomainName)
	_ = data.Set("cname", responseData.Cname)
	_ = data.Set("service_type", responseData.ServiceType)
	_ = data.Set("service_areas", responseData.ServiceAreas)
	_ = data.Set("comment", responseData.Comment)
	_ = data.Set("header_of_client_ip", responseData.HeaderOfClientIp)
	if responseData.OriginConfig != nil {
		originConfig := map[string]interface{}{}
		originConfig["origin_ips"] = responseData.OriginConfig.OriginIps
		originConfig["default_origin_host_header"] = responseData.OriginConfig.DefaultOriginHostHeader
		originConfig["use_range"] = responseData.OriginConfig.UseRange
		originConfig["follow301"] = responseData.OriginConfig.Follow301
		originConfig["follow302"] = responseData.OriginConfig.Follow302
		advSrcSetting := responseData.OriginConfig.AdvSrcSetting
		if advSrcSetting != nil {
			advSrcConfig := map[string]interface{}{}
			advSrcConfig["use_adv_src"] = advSrcSetting.UseAdvSrc
			advSrcConfig["detect_url"] = advSrcSetting.DetectUrl
			advSrcConfig["detect_period"] = advSrcSetting.DetectPeriod
			advSrcConfig["master_ips"] = advSrcSetting.MasterIps
			advSrcConfig["backup_ips"] = advSrcSetting.BackupIps
			originConfig["adv_src_setting"] = []interface{}{advSrcConfig}
		}
		_ = data.Set("origin_config", []interface{}{originConfig})
	}

	ssl := make([]interface{}, 0)
	if responseData.Ssl != nil {
		ssl = append(ssl, map[string]interface{}{
			"use_ssl":               responseData.Ssl.UseSsl,
			"ssl_certificate_id":    responseData.Ssl.SslCertificateId,
			"backup_certificate_id": responseData.Ssl.BackupCertificateId,
			"gm_certificate_ids":    responseData.Ssl.GmCertificateIds,
			"tls_version":           responseData.Ssl.TlsVersion,
			"enable_ocsp":           responseData.Ssl.EnableOcsp,
			"ssl_cipher_suite":      responseData.Ssl.SslCipherSuite,
		})
		_ = data.Set("ssl", ssl)
	}
	if responseData.CacheTimeBehaviors != nil && len(responseData.CacheTimeBehaviors) > 0 {
		cacheTimeBehaviors := make([]interface{}, 0)
		for _, cacheTimeBehavior := range responseData.CacheTimeBehaviors {
			cacheTimeBehaviors = append(cacheTimeBehaviors, map[string]interface{}{
				"path_pattern":                 cacheTimeBehavior.PathPattern,
				"except_path_pattern":          cacheTimeBehavior.ExceptPathPattern,
				"custom_pattern":               cacheTimeBehavior.CustomPattern,
				"file_type":                    cacheTimeBehavior.FileType,
				"custom_file_type":             cacheTimeBehavior.CustomFileType,
				"specify_url_pattern":          cacheTimeBehavior.SpecifyUrlPattern,
				"directory":                    cacheTimeBehavior.Directory,
				"cache_ttl":                    cacheTimeBehavior.CacheTtl,
				"ignore_cache_control":         cacheTimeBehavior.IgnoreCacheControl,
				"is_respect_server":            cacheTimeBehavior.IsRespectServer,
				"ignore_letter_case":           cacheTimeBehavior.IgnoreLetterCase,
				"reload_manage":                cacheTimeBehavior.ReloadManage,
				"priority":                     cacheTimeBehavior.Priority,
				"ignore_authentication_header": cacheTimeBehavior.IgnoreAuthenticationHeader,
			})
		}
		_ = data.Set("cache_time_behaviors", cacheTimeBehaviors)
	}

	if responseData.CacheKeyRules != nil && len(responseData.CacheKeyRules) > 0 {
		cacheKeyRules := make([]interface{}, 0)
		for _, cacheKeyRule := range responseData.CacheKeyRules {
			cacheKeyRules = append(cacheKeyRules, map[string]interface{}{
				"path_pattern":            cacheKeyRule.PathPattern,
				"specify_url":             cacheKeyRule.SpecifyUrl,
				"full_match4_specify_url": cacheKeyRule.FullMatch4SpecifyUrl,
				"custom_pattern":          cacheKeyRule.CustomPattern,
				"file_type":               cacheKeyRule.FileType,
				"custom_file_type":        cacheKeyRule.CustomFileType,
				"directory":               cacheKeyRule.Directory,
				"ignore_case":             cacheKeyRule.IgnoreCase,
				"header_name":             cacheKeyRule.HeaderName,
				"parameter_of_header":     cacheKeyRule.ParameterOfHeader,
				"priority":                cacheKeyRule.Priority,
			})
		}
		_ = data.Set("cache_key_rules", cacheKeyRules)

	}

	if responseData.QueryStringSettings != nil && len(responseData.QueryStringSettings) > 0 {
		queryStringSettings := make([]interface{}, 0)
		for _, queryStringSetting := range responseData.QueryStringSettings {
			queryStringSettings = append(queryStringSettings, map[string]interface{}{
				"path_pattern":         queryStringSetting.PathPattern,
				"except_path_pattern":  queryStringSetting.ExceptPathPattern,
				"file_types":           queryStringSetting.FileTypes,
				"custom_file_types":    queryStringSetting.CustomFileTypes,
				"custom_pattern":       queryStringSetting.CustomPattern,
				"specify_url_pattern":  queryStringSetting.SpecifyUrlPattern,
				"directories":          queryStringSetting.Directories,
				"priority":             queryStringSetting.Priority,
				"ignore_letter_case":   queryStringSetting.IgnoreLetterCase,
				"ignore_query_string":  queryStringSetting.IgnoreQueryString,
				"query_string_kept":    queryStringSetting.QueryStringKept,
				"query_string_removed": queryStringSetting.QueryStringRemoved,
				"source_with_query":    queryStringSetting.SourceWithQuery,
				"source_key_kept":      queryStringSetting.SourceKeyKept,
				"source_key_removed":   queryStringSetting.SourceKeyRemoved,
			})
		}
		_ = data.Set("query_string_settings", queryStringSettings)
	}

	if responseData.CacheByRespHeaders != nil && len(responseData.CacheByRespHeaders) > 0 {
		cacheByRespHeaders := make([]interface{}, 0)
		for _, cacheByRespHeader := range responseData.CacheByRespHeaders {
			cacheByRespHeaders = append(cacheByRespHeaders, map[string]interface{}{
				"response_header":     cacheByRespHeader.ResponseHeader,
				"path_pattern":        cacheByRespHeader.PathPattern,
				"except_path_pattern": cacheByRespHeader.ExceptPathPattern,
				"response_value":      cacheByRespHeader.ResponseValue,
				"ignore_letter_case":  cacheByRespHeader.IgnoreLetterCase,
				"priority":            cacheByRespHeader.Priority,
				"is_respheader":       cacheByRespHeader.IsRespheader,
			})
		}
		_ = data.Set("cache_by_resp_headers", cacheByRespHeaders)
	}

	if responseData.HttpCodeCacheRules != nil && len(responseData.HttpCodeCacheRules) > 0 {
		httpCodeCacheRules := make([]interface{}, 0)
		for _, httpCodeCacheRule := range responseData.HttpCodeCacheRules {
			httpCodes := make([]*string, 0)
			for _, httpCode := range httpCodeCacheRule.HttpCodes {
				httpCodes = append(httpCodes, httpCode)
			}
			httpCodeCacheRules = append(httpCodeCacheRules, map[string]interface{}{
				"http_codes": httpCodes,
				"cache_ttl":  httpCodeCacheRule.CacheTtl,
			})
		}
		_ = data.Set("http_code_cache_rules", httpCodeCacheRules)
	}

	if responseData.IgnoreProtocolRules != nil && len(responseData.IgnoreProtocolRules) > 0 {
		ignoreProtocolRules := make([]interface{}, 0)
		for _, ignoreProtocolRule := range responseData.IgnoreProtocolRules {
			ignoreProtocolRules = append(ignoreProtocolRules, map[string]interface{}{
				"path_pattern":          ignoreProtocolRule.PathPattern,
				"except_path_pattern":   ignoreProtocolRule.ExceptPathPattern,
				"cache_ignore_protocol": ignoreProtocolRule.CacheIgnoreProtocol,
				"purge_ignore_protocol": ignoreProtocolRule.PurgeIgnoreProtocol,
			})
		}
		_ = data.Set("ignore_protocol_rules", ignoreProtocolRules)
	}

	if responseData.Http2Settings != nil {
		http2Settings := map[string]interface{}{
			"enable_http2":            responseData.Http2Settings.EnableHttp2,
			"back_to_origin_protocol": responseData.Http2Settings.BackToOriginProtocol,
		}
		_ = data.Set("http2_settings", []interface{}{http2Settings})
	}

	if responseData.HeaderModifyRules != nil && len(responseData.HeaderModifyRules) > 0 {
		headerModifyRules := make([]interface{}, 0)
		for _, headerModifyRule := range responseData.HeaderModifyRules {
			headerModifyRules = append(headerModifyRules, map[string]interface{}{
				"path_pattern":          headerModifyRule.PathPattern,
				"except_path_pattern":   headerModifyRule.ExceptPathPattern,
				"custom_pattern":        headerModifyRule.CustomPattern,
				"file_type":             headerModifyRule.FileType,
				"custom_file_type":      headerModifyRule.CustomFileType,
				"directory":             headerModifyRule.Directory,
				"specify_url":           headerModifyRule.SpecifyUrl,
				"request_method":        headerModifyRule.RequestMethod,
				"header_direction":      headerModifyRule.HeaderDirection,
				"action":                headerModifyRule.Action,
				"allow_regexp":          headerModifyRule.AllowRegexp,
				"header_name":           headerModifyRule.HeaderName,
				"header_value":          headerModifyRule.HeaderValue,
				"request_header":        headerModifyRule.RequestHeader,
				"priority":              headerModifyRule.Priority,
				"except_file_type":      headerModifyRule.ExceptFileType,
				"except_directory":      headerModifyRule.ExceptDirectory,
				"except_request_method": headerModifyRule.ExceptRequestMethod,
				"except_request_header": headerModifyRule.ExceptRequestHeader,
			})
		}
		_ = data.Set("header_modify_rules", headerModifyRules)
	}

	if responseData.RewriteRuleSettings != nil && len(responseData.RewriteRuleSettings) > 0 {
		rewriteRuleSettings := make([]interface{}, 0)
		for _, rewriteRuleSetting := range responseData.RewriteRuleSettings {
			rewriteRuleSettings = append(rewriteRuleSettings, map[string]interface{}{
				"path_pattern":             rewriteRuleSetting.PathPattern,
				"custom_pattern":           rewriteRuleSetting.CustomPattern,
				"directory":                rewriteRuleSetting.Directory,
				"file_type":                rewriteRuleSetting.FileType,
				"custom_file_type":         rewriteRuleSetting.CustomFileType,
				"except_path_pattern":      rewriteRuleSetting.ExceptPathPattern,
				"ignore_letter_case":       rewriteRuleSetting.IgnoreLetterCase,
				"publish_type":             rewriteRuleSetting.PublishType,
				"priority":                 rewriteRuleSetting.Priority,
				"before_value":             rewriteRuleSetting.BeforeValue,
				"after_value":              rewriteRuleSetting.AfterValue,
				"rewrite_type":             rewriteRuleSetting.RewriteType,
				"request_header":           rewriteRuleSetting.RequestHeader,
				"exception_request_header": rewriteRuleSetting.ExceptionRequestHeader,
			})
		}
		_ = data.Set("rewrite_rule_settings", rewriteRuleSettings)
	}
	log.Printf("resource.wangsu_cdn_domain.read success")
	return nil
}

func resourceCdnDomainCreate(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_cdn_domain.create")

	var diags diag.Diagnostics
	request := &cdn.AddDomainForTerraformRequest{}
	if domainName, ok := data.Get("domain_name").(string); ok && domainName != "" {
		request.DomainName = &domainName
	}
	if serviceType, ok := data.Get("service_type").(string); ok && serviceType != "" {
		request.ServiceType = &serviceType
	}
	if serviceAreas, ok := data.Get("service_areas").(string); ok && serviceAreas != "" {
		request.ServiceAreas = &serviceAreas
	}
	if comment, ok := data.Get("comment").(string); ok && comment != "" {
		request.Comment = &comment
	}
	if headerOfClientIp, ok := data.Get("header_of_client_ip").(string); ok && headerOfClientIp != "" {
		request.HeaderOfClientIp = &headerOfClientIp
	}
	if originConfig, ok := data.Get("origin_config").([]interface{}); ok && len(originConfig) > 0 {
		for _, v := range originConfig {
			originConfigMap := v.(map[string]interface{})
			originIps := originConfigMap["origin_ips"].(string)
			defaultOriginHostHeader := originConfigMap["default_origin_host_header"].(string)
			useRange := originConfigMap["use_range"].(string)
			follow301 := originConfigMap["follow301"].(string)
			follow302 := originConfigMap["follow302"].(string)
			config := &cdn.AddDomainForTerraformRequestOriginConfig{
				OriginIps:               &originIps,
				DefaultOriginHostHeader: &defaultOriginHostHeader,
				UseRange:                &useRange,
				Follow301:               &follow301,
				Follow302:               &follow302,
			}
			advSrcSettings := originConfigMap["adv_src_setting"].([]interface{})
			if advSrcSettings != nil && len(advSrcSettings) > 0 {
				for _, item := range advSrcSettings {
					advSrcSetting := item.(map[string]interface{})
					useAdvSrc := advSrcSetting["use_adv_src"].(string)
					detectUrl := advSrcSetting["detect_url"].(string)
					detectPeriod := advSrcSetting["detect_period"].(string)
					masterIps := advSrcSetting["master_ips"].([]interface{})
					backupIps := advSrcSetting["backup_ips"].([]interface{})
					originConfigAdvSrcSetting := &cdn.AddDomainForTerraformRequestOriginConfigAdvSrcSetting{
						UseAdvSrc:    &useAdvSrc,
						DetectUrl:    &detectUrl,
						DetectPeriod: &detectPeriod,
					}
					if masterIps != nil && len(masterIps) > 0 {
						originConfigAdvSrcSetting.MasterIps = make([]*string, 0, len(masterIps))
						for _, ip := range masterIps {
							if ip == nil {
								diags = append(diags, diag.FromErr(errors.New("The master ip could not be empty."))...)
								return diags
							}
							masterIp := ip.(string)
							originConfigAdvSrcSetting.MasterIps = append(originConfigAdvSrcSetting.MasterIps, &masterIp)
						}
					}
					if backupIps != nil && len(backupIps) > 0 {
						originConfigAdvSrcSetting.BackupIps = make([]*string, 0, len(backupIps))
						for _, ip := range backupIps {
							if ip == nil {
								diags = append(diags, diag.FromErr(errors.New("The backup ip could not be empty."))...)
								return diags
							}
							backupIp := ip.(string)
							originConfigAdvSrcSetting.BackupIps = append(originConfigAdvSrcSetting.BackupIps, &backupIp)
						}
					}
					config.AdvSrcSetting = originConfigAdvSrcSetting
				}
			}
			request.OriginConfig = config
		}
	}

	if ssl, ok := data.Get("ssl").([]interface{}); ok && len(ssl) > 0 {
		for _, v := range ssl {
			sslMap := v.(map[string]interface{})
			useSsl := sslMap["use_ssl"].(string)
			sslCertificateId := sslMap["ssl_certificate_id"].(string)
			backupCertificateId := sslMap["backup_certificate_id"].(string)
			gmCertificateIds := sslMap["gm_certificate_ids"].([]interface{})
			var gmCertIds []*string = nil
			if gmCertificateIds != nil {
				gmCertIds = make([]*string, 0)
				for _, certId := range gmCertificateIds {
					if certId == nil {
						diags = append(diags, diag.FromErr(errors.New("The gm certificate id could not be empty."))...)
						return diags
					}
					certificateId := certId.(string)
					gmCertIds = append(gmCertIds, &certificateId)
				}
			}
			tlsVersion := sslMap["tls_version"].(string)
			enableOcsp := sslMap["enable_ocsp"].(string)
			sslCipherSuite := sslMap["ssl_cipher_suite"].(string)
			request.Ssl = &cdn.AddDomainForTerraformRequestSsl{
				UseSsl:              &useSsl,
				SslCertificateId:    &sslCertificateId,
				BackupCertificateId: &backupCertificateId,
				GmCertificateIds:    gmCertIds,
				TlsVersion:          &tlsVersion,
				EnableOcsp:          &enableOcsp,
				SslCipherSuite:      &sslCipherSuite,
			}
		}
	}

	if cacheTimeBehaviors, ok := data.Get("cache_time_behaviors").([]interface{}); ok && len(cacheTimeBehaviors) > 0 {
		for _, v := range cacheTimeBehaviors {
			cacheTimeBehaviorMap := v.(map[string]interface{})
			pathPattern := cacheTimeBehaviorMap["path_pattern"].(string)
			exceptPathPattern := cacheTimeBehaviorMap["except_path_pattern"].(string)
			customPattern := cacheTimeBehaviorMap["custom_pattern"].(string)
			fileType := cacheTimeBehaviorMap["file_type"].(string)
			customFileType := cacheTimeBehaviorMap["custom_file_type"].(string)
			specifyUrlPattern := cacheTimeBehaviorMap["specify_url_pattern"].(string)
			directory := cacheTimeBehaviorMap["directory"].(string)
			cacheTtl := cacheTimeBehaviorMap["cache_ttl"].(string)
			ignoreCacheControl := cacheTimeBehaviorMap["ignore_cache_control"].(string)
			isRespectServer := cacheTimeBehaviorMap["is_respect_server"].(string)
			ignoreLetterCase := cacheTimeBehaviorMap["ignore_letter_case"].(string)
			reloadManage := cacheTimeBehaviorMap["reload_manage"].(string)
			priority := cacheTimeBehaviorMap["priority"].(string)
			ignoreAuthenticationHeader := cacheTimeBehaviorMap["ignore_authentication_header"].(string)
			request.CacheTimeBehaviors = append(request.CacheTimeBehaviors, &cdn.AddDomainForTerraformRequestCacheTimeBehaviors{
				PathPattern:                &pathPattern,
				ExceptPathPattern:          &exceptPathPattern,
				CustomPattern:              &customPattern,
				FileType:                   &fileType,
				CustomFileType:             &customFileType,
				SpecifyUrlPattern:          &specifyUrlPattern,
				Directory:                  &directory,
				CacheTtl:                   &cacheTtl,
				IgnoreCacheControl:         &ignoreCacheControl,
				IsRespectServer:            &isRespectServer,
				IgnoreLetterCase:           &ignoreLetterCase,
				ReloadManage:               &reloadManage,
				Priority:                   &priority,
				IgnoreAuthenticationHeader: &ignoreAuthenticationHeader,
			})
		}
	}

	if cacheKeyRules, ok := data.Get("cache_key_rules").([]interface{}); ok && len(cacheKeyRules) > 0 {
		for _, v := range cacheKeyRules {
			cacheKeyRuleMap := v.(map[string]interface{})
			pathPattern := cacheKeyRuleMap["path_pattern"].(string)
			specifyUrl := cacheKeyRuleMap["specify_url"].(string)
			fullMatch4SpecifyUrl := cacheKeyRuleMap["full_match4_specify_url"].(string)
			customPattern := cacheKeyRuleMap["custom_pattern"].(string)
			fileType := cacheKeyRuleMap["file_type"].(string)
			customFileType := cacheKeyRuleMap["custom_file_type"].(string)
			directory := cacheKeyRuleMap["directory"].(string)
			ignoreCase := cacheKeyRuleMap["ignore_case"].(string)
			headerName := cacheKeyRuleMap["header_name"].(string)
			parameterOfHeader := cacheKeyRuleMap["parameter_of_header"].(string)
			priority := cacheKeyRuleMap["priority"].(string)
			request.CacheKeyRules = append(request.CacheKeyRules, &cdn.AddDomainForTerraformRequestCacheKeyRules{
				PathPattern:          &pathPattern,
				SpecifyUrl:           &specifyUrl,
				FullMatch4SpecifyUrl: &fullMatch4SpecifyUrl,
				CustomPattern:        &customPattern,
				FileType:             &fileType,
				CustomFileType:       &customFileType,
				Directory:            &directory,
				IgnoreCase:           &ignoreCase,
				HeaderName:           &headerName,
				ParameterOfHeader:    &parameterOfHeader,
				Priority:             &priority,
			})
		}

	}

	if queryStringSettings, ok := data.Get("query_string_settings").([]interface{}); ok && len(queryStringSettings) > 0 {
		for _, v := range queryStringSettings {
			queryStringSettingMap := v.(map[string]interface{})
			pathPattern := queryStringSettingMap["path_pattern"].(string)
			exceptPathPattern := queryStringSettingMap["except_path_pattern"].(string)
			fileTypes := queryStringSettingMap["file_types"].(string)
			customFileTypes := queryStringSettingMap["custom_file_types"].(string)
			customPattern := queryStringSettingMap["custom_pattern"].(string)
			specifyUrlPattern := queryStringSettingMap["specify_url_pattern"].(string)
			directories := queryStringSettingMap["directories"].(string)
			priority := queryStringSettingMap["priority"].(string)
			ignoreLetterCase := queryStringSettingMap["ignore_letter_case"].(string)
			ignoreQueryString := queryStringSettingMap["ignore_query_string"].(string)
			queryStringKept := queryStringSettingMap["query_string_kept"].(string)
			queryStringRemoved := queryStringSettingMap["query_string_removed"].(string)
			sourceWithQuery := queryStringSettingMap["source_with_query"].(string)
			sourceKeyKept := queryStringSettingMap["source_key_kept"].(string)
			sourceKeyRemoved := queryStringSettingMap["source_key_removed"].(string)
			request.QueryStringSettings = append(request.QueryStringSettings, &cdn.AddDomainForTerraformRequestQueryStringSettings{
				PathPattern:        &pathPattern,
				ExceptPathPattern:  &exceptPathPattern,
				FileTypes:          &fileTypes,
				CustomFileTypes:    &customFileTypes,
				CustomPattern:      &customPattern,
				SpecifyUrlPattern:  &specifyUrlPattern,
				Directories:        &directories,
				Priority:           &priority,
				IgnoreLetterCase:   &ignoreLetterCase,
				IgnoreQueryString:  &ignoreQueryString,
				QueryStringKept:    &queryStringKept,
				QueryStringRemoved: &queryStringRemoved,
				SourceWithQuery:    &sourceWithQuery,
				SourceKeyKept:      &sourceKeyKept,
				SourceKeyRemoved:   &sourceKeyRemoved,
			})
		}

	}

	if cacheByRespHeaders, ok := data.Get("cache_by_resp_headers").([]interface{}); ok && len(cacheByRespHeaders) > 0 {
		for _, v := range cacheByRespHeaders {
			cacheByRespHeaderMap := v.(map[string]interface{})
			responseHeader := cacheByRespHeaderMap["response_header"].(string)
			pathPattern := cacheByRespHeaderMap["path_pattern"].(string)
			exceptPathPattern := cacheByRespHeaderMap["except_path_pattern"].(string)
			responseValue := cacheByRespHeaderMap["response_value"].(string)
			ignoreLetterCase := cacheByRespHeaderMap["ignore_letter_case"].(string)
			priority := cacheByRespHeaderMap["priority"].(string)
			isRespheader := cacheByRespHeaderMap["is_respheader"].(string)
			request.CacheByRespHeaders = append(request.CacheByRespHeaders, &cdn.AddDomainForTerraformRequestCacheByRespHeaders{
				ResponseHeader:    &responseHeader,
				PathPattern:       &pathPattern,
				ExceptPathPattern: &exceptPathPattern,
				ResponseValue:     &responseValue,
				IgnoreLetterCase:  &ignoreLetterCase,
				Priority:          &priority,
				IsRespheader:      &isRespheader,
			})
		}
	}

	if httpCodeCacheRules, ok := data.Get("http_code_cache_rules").([]interface{}); ok && len(httpCodeCacheRules) > 0 {
		for _, v := range httpCodeCacheRules {
			httpCodeCacheRuleMap := v.(map[string]interface{})
			httpCodes := httpCodeCacheRuleMap["http_codes"].([]interface{})
			var codes []*string
			if httpCodes != nil {
				codes = make([]*string, 0)
				for _, code := range httpCodes {
					if code == nil {
						diags = append(diags, diag.FromErr(errors.New("The http code could not be empty."))...)
						return diags
					}
					httpCode := code.(string)
					codes = append(codes, &httpCode)
				}
			}
			cacheTtl := httpCodeCacheRuleMap["cache_ttl"].(string)
			request.HttpCodeCacheRules = append(request.HttpCodeCacheRules, &cdn.AddDomainForTerraformRequestHttpCodeCacheRules{
				HttpCodes: codes,
				CacheTtl:  &cacheTtl,
			})
		}
	}

	if ignoreProtocolRules, ok := data.Get("ignore_protocol_rules").([]interface{}); ok && len(ignoreProtocolRules) > 0 {
		for _, v := range ignoreProtocolRules {
			ignoreProtocolRuleMap := v.(map[string]interface{})
			pathPattern := ignoreProtocolRuleMap["path_pattern"].(string)
			exceptPathPattern := ignoreProtocolRuleMap["except_path_pattern"].(string)
			cacheIgnoreProtocol := ignoreProtocolRuleMap["cache_ignore_protocol"].(string)
			purgeIgnoreProtocol := ignoreProtocolRuleMap["purge_ignore_protocol"].(string)
			request.IgnoreProtocolRules = append(request.IgnoreProtocolRules, &cdn.AddDomainForTerraformRequestIgnoreProtocolRules{
				PathPattern:         &pathPattern,
				ExceptPathPattern:   &exceptPathPattern,
				CacheIgnoreProtocol: &cacheIgnoreProtocol,
				PurgeIgnoreProtocol: &purgeIgnoreProtocol,
			})
		}
	}

	if http2Settings, ok := data.Get("http2_settings").([]interface{}); ok && len(http2Settings) > 0 {
		for _, v := range http2Settings {
			http2SettingMap := v.(map[string]interface{})
			enableHttp2 := http2SettingMap["enable_http2"].(string)
			backToOriginProtocol := http2SettingMap["back_to_origin_protocol"].(string)
			request.Http2Settings = &cdn.AddDomainForTerraformRequestHttp2Settings{
				EnableHttp2:          &enableHttp2,
				BackToOriginProtocol: &backToOriginProtocol,
			}
		}
	}

	if headerModifyRules, ok := data.Get("header_modify_rules").([]interface{}); ok && len(headerModifyRules) > 0 {
		for _, v := range headerModifyRules {
			headerModifyRuleMap := v.(map[string]interface{})
			pathPattern := headerModifyRuleMap["path_pattern"].(string)
			exceptPathPattern := headerModifyRuleMap["except_path_pattern"].(string)
			customPattern := headerModifyRuleMap["custom_pattern"].(string)
			fileType := headerModifyRuleMap["file_type"].(string)
			customFileType := headerModifyRuleMap["custom_file_type"].(string)
			directory := headerModifyRuleMap["directory"].(string)
			specifyUrl := headerModifyRuleMap["specify_url"].(string)
			requestMethod := headerModifyRuleMap["request_method"].(string)
			headerDirection := headerModifyRuleMap["header_direction"].(string)
			action := headerModifyRuleMap["action"].(string)
			allowRegexp := headerModifyRuleMap["allow_regexp"].(string)
			headerName := headerModifyRuleMap["header_name"].(string)
			headerValue := headerModifyRuleMap["header_value"].(string)
			requestHeader := headerModifyRuleMap["request_header"].(string)
			priority := headerModifyRuleMap["priority"].(string)
			exceptFileType := headerModifyRuleMap["except_file_type"].(string)
			exceptDirectory := headerModifyRuleMap["except_directory"].(string)
			exceptRequestMethod := headerModifyRuleMap["except_request_method"].(string)
			exceptRequestHeader := headerModifyRuleMap["except_request_header"].(string)
			request.HeaderModifyRules = append(request.HeaderModifyRules, &cdn.AddDomainForTerraformRequestHeaderModifyRules{
				PathPattern:         &pathPattern,
				ExceptPathPattern:   &exceptPathPattern,
				CustomPattern:       &customPattern,
				FileType:            &fileType,
				CustomFileType:      &customFileType,
				Directory:           &directory,
				SpecifyUrl:          &specifyUrl,
				RequestMethod:       &requestMethod,
				HeaderDirection:     &headerDirection,
				Action:              &action,
				AllowRegexp:         &allowRegexp,
				HeaderName:          &headerName,
				HeaderValue:         &headerValue,
				RequestHeader:       &requestHeader,
				Priority:            &priority,
				ExceptFileType:      &exceptFileType,
				ExceptDirectory:     &exceptDirectory,
				ExceptRequestMethod: &exceptRequestMethod,
				ExceptRequestHeader: &exceptRequestHeader,
			})
		}
	}

	if rewriteRuleSettings, ok := data.Get("rewrite_rule_settings").([]interface{}); ok && len(rewriteRuleSettings) > 0 {
		for _, v := range rewriteRuleSettings {
			rewriteRuleSettingMap := v.(map[string]interface{})
			pathPattern := rewriteRuleSettingMap["path_pattern"].(string)
			customPattern := rewriteRuleSettingMap["custom_pattern"].(string)
			directory := rewriteRuleSettingMap["directory"].(string)
			fileType := rewriteRuleSettingMap["file_type"].(string)
			customFileType := rewriteRuleSettingMap["custom_file_type"].(string)
			exceptPathPattern := rewriteRuleSettingMap["except_path_pattern"].(string)
			ignoreLetterCase := rewriteRuleSettingMap["ignore_letter_case"].(string)
			publishType := rewriteRuleSettingMap["publish_type"].(string)
			priority := rewriteRuleSettingMap["priority"].(string)
			beforeValue := rewriteRuleSettingMap["before_value"].(string)
			afterValue := rewriteRuleSettingMap["after_value"].(string)
			rewriteType := rewriteRuleSettingMap["rewrite_type"].(string)
			requestHeader := rewriteRuleSettingMap["request_header"].(string)
			exceptionRequestHeader := rewriteRuleSettingMap["exception_request_header"].(string)
			request.RewriteRuleSettings = append(request.RewriteRuleSettings, &cdn.AddDomainForTerraformRequestRewriteRuleSettings{
				PathPattern:            &pathPattern,
				CustomPattern:          &customPattern,
				Directory:              &directory,
				FileType:               &fileType,
				CustomFileType:         &customFileType,
				ExceptPathPattern:      &exceptPathPattern,
				IgnoreLetterCase:       &ignoreLetterCase,
				PublishType:            &publishType,
				Priority:               &priority,
				BeforeValue:            &beforeValue,
				AfterValue:             &afterValue,
				RewriteType:            &rewriteType,
				RequestHeader:          &requestHeader,
				ExceptionRequestHeader: &exceptionRequestHeader,
			})
		}
	}
	//start to create a domain in 2 minutes
	var createDomainResponse *cdn.AddDomainForTerraformResponse
	var requestId string
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		requestId, createDomainResponse, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseCdnClient().AddCdnDomain(request)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}
	if createDomainResponse == nil {
		data.SetId("")
		return nil
	}

	data.SetId(*request.DomainName)

	time.Sleep(3 * time.Second)
	//query domain deployment status
	var response *cdn.QueryDeployResultForTerraformResponse
	err = resource.RetryContext(context, time.Duration(5)*time.Minute, func() *resource.RetryError {
		response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseCdnClient().QueryDomainDeployStatus(requestId)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if response != nil && response.Data != nil && *response.Data.DeployResult != "SUCCESS" {
			return resource.RetryableError(fmt.Errorf("domain deployment status is in progress, retrying"))
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

	log.Printf("resource.wangsu_cdn_domain.create success")
	_ = data.Set("xCncRequestId", response.Data.RequestId)
	//set status
	return resourceCdnDomainRead(context, data, meta)
}

func resourceCdnDomainUpdate(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_cdn_domain.update")
	request := &cdn.UpdateDomainForTerraformRequest{}
	var diags diag.Diagnostics
	if data.HasChanges("service_areas") {
		serviceAreas := data.Get("service_areas").(string)
		request.ServiceAreas = &serviceAreas
	}
	if data.HasChanges("comment") {
		comment := data.Get("comment").(string)
		request.Comment = &comment
	}
	if data.HasChanges("header_of_client_ip") {
		headerOfClientIp := data.Get("header_of_client_ip").(string)
		request.HeaderOfClientIp = &headerOfClientIp
	}
	if data.HasChanges("origin_config") {
		if originConfig, ok := data.Get("origin_config").([]interface{}); ok && len(originConfig) > 0 {
			for _, v := range originConfig {
				originConfigMap := v.(map[string]interface{})
				originIps := originConfigMap["origin_ips"].(string)
				defaultOriginHostHeader := originConfigMap["default_origin_host_header"].(string)
				useRange := originConfigMap["use_range"].(string)
				follow301 := originConfigMap["follow301"].(string)
				follow302 := originConfigMap["follow302"].(string)
				config := &cdn.UpdateDomainForTerraformRequestOriginConfig{
					OriginIps:               &originIps,
					DefaultOriginHostHeader: &defaultOriginHostHeader,
					UseRange:                &useRange,
					Follow301:               &follow301,
					Follow302:               &follow302,
				}
				advSrcSettings := originConfigMap["adv_src_setting"].([]interface{})
				if advSrcSettings != nil && len(advSrcSettings) > 0 {
					for _, item := range advSrcSettings {
						advSrcSetting := item.(map[string]interface{})
						useAdvSrc := advSrcSetting["use_adv_src"].(string)
						detectUrl := advSrcSetting["detect_url"].(string)
						detectPeriod := advSrcSetting["detect_period"].(string)
						masterIps := advSrcSetting["master_ips"].([]interface{})
						backupIps := advSrcSetting["backup_ips"].([]interface{})
						originConfigAdvSrcSetting := &cdn.UpdateDomainForTerraformRequestOriginConfigAdvSrcSetting{
							UseAdvSrc:    &useAdvSrc,
							DetectUrl:    &detectUrl,
							DetectPeriod: &detectPeriod,
						}
						if masterIps != nil && len(masterIps) > 0 {
							originConfigAdvSrcSetting.MasterIps = make([]*string, 0, len(masterIps))
							for _, ip := range masterIps {
								if ip == nil {
									diags = append(diags, diag.FromErr(errors.New("The master ip could not be empty."))...)
									return diags
								}
								masterIp := ip.(string)
								originConfigAdvSrcSetting.MasterIps = append(originConfigAdvSrcSetting.MasterIps, &masterIp)
							}
						}
						if backupIps != nil && len(backupIps) > 0 {
							originConfigAdvSrcSetting.BackupIps = make([]*string, 0, len(backupIps))
							for _, ip := range backupIps {
								if ip == nil {
									diags = append(diags, diag.FromErr(errors.New("The backup ip could not be empty."))...)
									return diags
								}
								backupIp := ip.(string)
								originConfigAdvSrcSetting.BackupIps = append(originConfigAdvSrcSetting.BackupIps, &backupIp)
							}
						}
						config.AdvSrcSetting = originConfigAdvSrcSetting
					}
				} else {
					useAdvSrc := "false"
					config.AdvSrcSetting = &cdn.UpdateDomainForTerraformRequestOriginConfigAdvSrcSetting{UseAdvSrc: &useAdvSrc}
				}
				request.OriginConfig = config
			}
		}
	}

	if data.HasChanges("ssl") {
		if ssl, ok := data.Get("ssl").([]interface{}); ok && len(ssl) > 0 {
			for _, v := range ssl {
				sslMap := v.(map[string]interface{})
				useSsl := sslMap["use_ssl"].(string)
				sslCertificateId := sslMap["ssl_certificate_id"].(string)
				backupCertificateId := sslMap["backup_certificate_id"].(string)
				gmCertificateIds := sslMap["gm_certificate_ids"].([]interface{})
				var gmCertIds = make([]*string, 0)
				if gmCertificateIds != nil {
					for _, certId := range gmCertificateIds {
						if certId == nil {
							diags = append(diags, diag.FromErr(errors.New("The gm certificate id could not be empty."))...)
							return diags
						}
						certificateId := certId.(string)
						gmCertIds = append(gmCertIds, &certificateId)
					}
				}
				tlsVersion := sslMap["tls_version"].(string)
				enableOcsp := sslMap["enable_ocsp"].(string)
				sslCipherSuite := sslMap["ssl_cipher_suite"].(string)
				request.Ssl = &cdn.UpdateDomainForTerraformRequestSsl{
					UseSsl:              &useSsl,
					SslCertificateId:    &sslCertificateId,
					BackupCertificateId: &backupCertificateId,
					GmCertificateIds:    gmCertIds,
					TlsVersion:          &tlsVersion,
					EnableOcsp:          &enableOcsp,
					SslCipherSuite:      &sslCipherSuite,
				}
			}
		} else {
			useSsl := "false"
			request.Ssl = &cdn.UpdateDomainForTerraformRequestSsl{}
			request.Ssl.UseSsl = &useSsl
		}
	}

	if data.HasChanges("cache_time_behaviors") {
		if cacheTimeBehaviors, ok := data.Get("cache_time_behaviors").([]interface{}); ok && len(cacheTimeBehaviors) > 0 {
			for _, v := range cacheTimeBehaviors {
				cacheTimeBehaviorMap := v.(map[string]interface{})
				pathPattern := cacheTimeBehaviorMap["path_pattern"].(string)
				exceptPathPattern := cacheTimeBehaviorMap["except_path_pattern"].(string)
				customPattern := cacheTimeBehaviorMap["custom_pattern"].(string)
				fileType := cacheTimeBehaviorMap["file_type"].(string)
				customFileType := cacheTimeBehaviorMap["custom_file_type"].(string)
				specifyUrlPattern := cacheTimeBehaviorMap["specify_url_pattern"].(string)
				directory := cacheTimeBehaviorMap["directory"].(string)
				cacheTtl := cacheTimeBehaviorMap["cache_ttl"].(string)
				ignoreCacheControl := cacheTimeBehaviorMap["ignore_cache_control"].(string)
				isRespectServer := cacheTimeBehaviorMap["is_respect_server"].(string)
				ignoreLetterCase := cacheTimeBehaviorMap["ignore_letter_case"].(string)
				reloadManage := cacheTimeBehaviorMap["reload_manage"].(string)
				priority := cacheTimeBehaviorMap["priority"].(string)
				ignoreAuthenticationHeader := cacheTimeBehaviorMap["ignore_authentication_header"].(string)
				request.CacheTimeBehaviors = append(request.CacheTimeBehaviors, &cdn.UpdateDomainForTerraformRequestCacheTimeBehaviors{
					PathPattern:                &pathPattern,
					ExceptPathPattern:          &exceptPathPattern,
					CustomPattern:              &customPattern,
					FileType:                   &fileType,
					CustomFileType:             &customFileType,
					SpecifyUrlPattern:          &specifyUrlPattern,
					Directory:                  &directory,
					CacheTtl:                   &cacheTtl,
					IgnoreCacheControl:         &ignoreCacheControl,
					IsRespectServer:            &isRespectServer,
					IgnoreLetterCase:           &ignoreLetterCase,
					ReloadManage:               &reloadManage,
					Priority:                   &priority,
					IgnoreAuthenticationHeader: &ignoreAuthenticationHeader,
				})
			}
		} else {
			request.CacheTimeBehaviors = make([]*cdn.UpdateDomainForTerraformRequestCacheTimeBehaviors, 0)
		}
	}

	if data.HasChanges("cache_key_rules") {
		if cacheKeyRules, ok := data.Get("cache_key_rules").([]interface{}); ok && len(cacheKeyRules) > 0 {
			for _, v := range cacheKeyRules {
				cacheKeyRuleMap := v.(map[string]interface{})
				pathPattern := cacheKeyRuleMap["path_pattern"].(string)
				specifyUrl := cacheKeyRuleMap["specify_url"].(string)
				fullMatch4SpecifyUrl := cacheKeyRuleMap["full_match4_specify_url"].(string)
				customPattern := cacheKeyRuleMap["custom_pattern"].(string)
				fileType := cacheKeyRuleMap["file_type"].(string)
				customFileType := cacheKeyRuleMap["custom_file_type"].(string)
				directory := cacheKeyRuleMap["directory"].(string)
				ignoreCase := cacheKeyRuleMap["ignore_case"].(string)
				headerName := cacheKeyRuleMap["header_name"].(string)
				parameterOfHeader := cacheKeyRuleMap["parameter_of_header"].(string)
				priority := cacheKeyRuleMap["priority"].(string)
				request.CacheKeyRules = append(request.CacheKeyRules, &cdn.UpdateDomainForTerraformRequestCacheKeyRules{
					PathPattern:          &pathPattern,
					SpecifyUrl:           &specifyUrl,
					FullMatch4SpecifyUrl: &fullMatch4SpecifyUrl,
					CustomPattern:        &customPattern,
					FileType:             &fileType,
					CustomFileType:       &customFileType,
					Directory:            &directory,
					IgnoreCase:           &ignoreCase,
					HeaderName:           &headerName,
					ParameterOfHeader:    &parameterOfHeader,
					Priority:             &priority,
				})
			}
		} else {
			request.CacheKeyRules = make([]*cdn.UpdateDomainForTerraformRequestCacheKeyRules, 0)
		}
	}

	if data.HasChanges("query_string_settings") {
		if queryStringSettings, ok := data.Get("query_string_settings").([]interface{}); ok && len(queryStringSettings) > 0 {
			for _, v := range queryStringSettings {
				queryStringSettingMap := v.(map[string]interface{})
				pathPattern := queryStringSettingMap["path_pattern"].(string)
				exceptPathPattern := queryStringSettingMap["except_path_pattern"].(string)
				fileTypes := queryStringSettingMap["file_types"].(string)
				customFileTypes := queryStringSettingMap["custom_file_types"].(string)
				customPattern := queryStringSettingMap["custom_pattern"].(string)
				specifyUrlPattern := queryStringSettingMap["specify_url_pattern"].(string)
				directories := queryStringSettingMap["directories"].(string)
				priority := queryStringSettingMap["priority"].(string)
				ignoreLetterCase := queryStringSettingMap["ignore_letter_case"].(string)
				ignoreQueryString := queryStringSettingMap["ignore_query_string"].(string)
				queryStringKept := queryStringSettingMap["query_string_kept"].(string)
				queryStringRemoved := queryStringSettingMap["query_string_removed"].(string)
				sourceWithQuery := queryStringSettingMap["source_with_query"].(string)
				sourceKeyKept := queryStringSettingMap["source_key_kept"].(string)
				sourceKeyRemoved := queryStringSettingMap["source_key_removed"].(string)
				request.QueryStringSettings = append(request.QueryStringSettings, &cdn.UpdateDomainForTerraformRequestQueryStringSettings{
					PathPattern:        &pathPattern,
					ExceptPathPattern:  &exceptPathPattern,
					FileTypes:          &fileTypes,
					CustomFileTypes:    &customFileTypes,
					CustomPattern:      &customPattern,
					SpecifyUrlPattern:  &specifyUrlPattern,
					Directories:        &directories,
					Priority:           &priority,
					IgnoreLetterCase:   &ignoreLetterCase,
					IgnoreQueryString:  &ignoreQueryString,
					QueryStringKept:    &queryStringKept,
					QueryStringRemoved: &queryStringRemoved,
					SourceWithQuery:    &sourceWithQuery,
					SourceKeyKept:      &sourceKeyKept,
					SourceKeyRemoved:   &sourceKeyRemoved,
				})
			}
		} else {
			request.QueryStringSettings = make([]*cdn.UpdateDomainForTerraformRequestQueryStringSettings, 0)
		}
	}

	if data.HasChanges("cache_by_resp_headers") {
		if cacheByRespHeaders, ok := data.Get("cache_by_resp_headers").([]interface{}); ok && len(cacheByRespHeaders) > 0 {
			for _, v := range cacheByRespHeaders {
				cacheByRespHeaderMap := v.(map[string]interface{})
				responseHeader := cacheByRespHeaderMap["response_header"].(string)
				pathPattern := cacheByRespHeaderMap["path_pattern"].(string)
				exceptPathPattern := cacheByRespHeaderMap["except_path_pattern"].(string)
				responseValue := cacheByRespHeaderMap["response_value"].(string)
				ignoreLetterCase := cacheByRespHeaderMap["ignore_letter_case"].(string)
				priority := cacheByRespHeaderMap["priority"].(string)
				isRespheader := cacheByRespHeaderMap["is_respheader"].(string)
				request.CacheByRespHeaders = append(request.CacheByRespHeaders, &cdn.UpdateDomainForTerraformRequestCacheByRespHeaders{
					ResponseHeader:    &responseHeader,
					PathPattern:       &pathPattern,
					ExceptPathPattern: &exceptPathPattern,
					ResponseValue:     &responseValue,
					IgnoreLetterCase:  &ignoreLetterCase,
					Priority:          &priority,
					IsRespheader:      &isRespheader,
				})
			}
		} else {
			request.CacheByRespHeaders = make([]*cdn.UpdateDomainForTerraformRequestCacheByRespHeaders, 0)
		}
	}

	if data.HasChanges("http_code_cache_rules") {
		if httpCodeCacheRules, ok := data.Get("http_code_cache_rules").([]interface{}); ok && len(httpCodeCacheRules) > 0 {
			for _, v := range httpCodeCacheRules {
				httpCodeCacheRuleMap := v.(map[string]interface{})
				httpCodes := httpCodeCacheRuleMap["http_codes"].([]interface{})
				var codes []*string
				if httpCodes != nil {
					codes = make([]*string, 0)
					for _, code := range httpCodes {
						if code == nil {
							diags = append(diags, diag.FromErr(errors.New("The http code could not be empty."))...)
							return diags
						}
						httpCode := code.(string)
						codes = append(codes, &httpCode)
					}
				}
				cacheTtl := httpCodeCacheRuleMap["cache_ttl"].(string)
				request.HttpCodeCacheRules = append(request.HttpCodeCacheRules, &cdn.UpdateDomainForTerraformRequestHttpCodeCacheRules{
					HttpCodes: codes,
					CacheTtl:  &cacheTtl,
				})
			}
		} else {
			request.HttpCodeCacheRules = make([]*cdn.UpdateDomainForTerraformRequestHttpCodeCacheRules, 0)
		}

	}

	if data.HasChanges("ignore_protocol_rules") {
		if ignoreProtocolRules, ok := data.Get("ignore_protocol_rules").([]interface{}); ok && len(ignoreProtocolRules) > 0 {
			for _, v := range ignoreProtocolRules {
				ignoreProtocolRuleMap := v.(map[string]interface{})
				pathPattern := ignoreProtocolRuleMap["path_pattern"].(string)
				exceptPathPattern := ignoreProtocolRuleMap["except_path_pattern"].(string)
				cacheIgnoreProtocol := ignoreProtocolRuleMap["cache_ignore_protocol"].(string)
				purgeIgnoreProtocol := ignoreProtocolRuleMap["purge_ignore_protocol"].(string)
				request.IgnoreProtocolRules = append(request.IgnoreProtocolRules, &cdn.UpdateDomainForTerraformRequestIgnoreProtocolRules{
					PathPattern:         &pathPattern,
					ExceptPathPattern:   &exceptPathPattern,
					CacheIgnoreProtocol: &cacheIgnoreProtocol,
					PurgeIgnoreProtocol: &purgeIgnoreProtocol,
				})
			}
		} else {
			request.IgnoreProtocolRules = make([]*cdn.UpdateDomainForTerraformRequestIgnoreProtocolRules, 0)
		}
	}

	if data.HasChanges("http2_settings") {
		if http2Settings, ok := data.Get("http2_settings").([]interface{}); ok && len(http2Settings) > 0 {
			for _, v := range http2Settings {
				http2SettingMap := v.(map[string]interface{})
				enableHttp2 := http2SettingMap["enable_http2"].(string)
				backToOriginProtocol := http2SettingMap["back_to_origin_protocol"].(string)
				request.Http2Settings = &cdn.UpdateDomainForTerraformRequestHttp2Settings{
					EnableHttp2:          &enableHttp2,
					BackToOriginProtocol: &backToOriginProtocol,
				}
			}
		} else {
			request.Http2Settings = &cdn.UpdateDomainForTerraformRequestHttp2Settings{}
		}
	}

	if data.HasChanges("header_modify_rules") {
		if headerModifyRules, ok := data.Get("header_modify_rules").([]interface{}); ok && len(headerModifyRules) > 0 {
			for _, v := range headerModifyRules {
				headerModifyRuleMap := v.(map[string]interface{})
				pathPattern := headerModifyRuleMap["path_pattern"].(string)
				exceptPathPattern := headerModifyRuleMap["except_path_pattern"].(string)
				customPattern := headerModifyRuleMap["custom_pattern"].(string)
				fileType := headerModifyRuleMap["file_type"].(string)
				customFileType := headerModifyRuleMap["custom_file_type"].(string)
				directory := headerModifyRuleMap["directory"].(string)
				specifyUrl := headerModifyRuleMap["specify_url"].(string)
				requestMethod := headerModifyRuleMap["request_method"].(string)
				headerDirection := headerModifyRuleMap["header_direction"].(string)
				action := headerModifyRuleMap["action"].(string)
				allowRegexp := headerModifyRuleMap["allow_regexp"].(string)
				headerName := headerModifyRuleMap["header_name"].(string)
				headerValue := headerModifyRuleMap["header_value"].(string)
				requestHeader := headerModifyRuleMap["request_header"].(string)
				priority := headerModifyRuleMap["priority"].(string)
				exceptFileType := headerModifyRuleMap["except_file_type"].(string)
				exceptDirectory := headerModifyRuleMap["except_directory"].(string)
				exceptRequestMethod := headerModifyRuleMap["except_request_method"].(string)
				exceptRequestHeader := headerModifyRuleMap["except_request_header"].(string)
				request.HeaderModifyRules = append(request.HeaderModifyRules, &cdn.UpdateDomainForTerraformRequestHeaderModifyRules{
					PathPattern:         &pathPattern,
					ExceptPathPattern:   &exceptPathPattern,
					CustomPattern:       &customPattern,
					FileType:            &fileType,
					CustomFileType:      &customFileType,
					Directory:           &directory,
					SpecifyUrl:          &specifyUrl,
					RequestMethod:       &requestMethod,
					HeaderDirection:     &headerDirection,
					Action:              &action,
					AllowRegexp:         &allowRegexp,
					HeaderName:          &headerName,
					HeaderValue:         &headerValue,
					RequestHeader:       &requestHeader,
					Priority:            &priority,
					ExceptFileType:      &exceptFileType,
					ExceptDirectory:     &exceptDirectory,
					ExceptRequestMethod: &exceptRequestMethod,
					ExceptRequestHeader: &exceptRequestHeader,
				})
			}
		} else {
			request.HeaderModifyRules = make([]*cdn.UpdateDomainForTerraformRequestHeaderModifyRules, 0)
		}
	}

	if data.HasChanges("rewrite_rule_settings") {
		if rewriteRuleSettings, ok := data.Get("rewrite_rule_settings").([]interface{}); ok && len(rewriteRuleSettings) > 0 {
			for _, v := range rewriteRuleSettings {
				rewriteRuleSettingMap := v.(map[string]interface{})
				pathPattern := rewriteRuleSettingMap["path_pattern"].(string)
				customPattern := rewriteRuleSettingMap["custom_pattern"].(string)
				directory := rewriteRuleSettingMap["directory"].(string)
				fileType := rewriteRuleSettingMap["file_type"].(string)
				customFileType := rewriteRuleSettingMap["custom_file_type"].(string)
				exceptPathPattern := rewriteRuleSettingMap["except_path_pattern"].(string)
				ignoreLetterCase := rewriteRuleSettingMap["ignore_letter_case"].(string)
				publishType := rewriteRuleSettingMap["publish_type"].(string)
				priority := rewriteRuleSettingMap["priority"].(string)
				beforeValue := rewriteRuleSettingMap["before_value"].(string)
				afterValue := rewriteRuleSettingMap["after_value"].(string)
				rewriteType := rewriteRuleSettingMap["rewrite_type"].(string)
				requestHeader := rewriteRuleSettingMap["request_header"].(string)
				exceptionRequestHeader := rewriteRuleSettingMap["exception_request_header"].(string)
				request.RewriteRuleSettings = append(request.RewriteRuleSettings, &cdn.UpdateDomainForTerraformRequestRewriteRuleSettings{
					PathPattern:            &pathPattern,
					CustomPattern:          &customPattern,
					Directory:              &directory,
					FileType:               &fileType,
					CustomFileType:         &customFileType,
					ExceptPathPattern:      &exceptPathPattern,
					IgnoreLetterCase:       &ignoreLetterCase,
					PublishType:            &publishType,
					Priority:               &priority,
					BeforeValue:            &beforeValue,
					AfterValue:             &afterValue,
					RewriteType:            &rewriteType,
					RequestHeader:          &requestHeader,
					ExceptionRequestHeader: &exceptionRequestHeader,
				})
			}
		} else {
			request.RewriteRuleSettings = make([]*cdn.UpdateDomainForTerraformRequestRewriteRuleSettings, 0)
		}
	}

	var editResponse *cdn.UpdateDomainForTerraformResponse
	var requestId string
	var err error
	err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
		requestId, editResponse, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseCdnClient().UpdateCdnDomain(request, data.Id())
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	if editResponse == nil {
		data.SetId("")
		return nil
	}

	time.Sleep(3 * time.Second)

	//query domain deployment status
	var response *cdn.QueryDeployResultForTerraformResponse
	err = resource.RetryContext(context, time.Duration(5)*time.Minute, func() *resource.RetryError {
		response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseCdnClient().QueryDomainDeployStatus(requestId)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if response != nil && response.Data != nil && *response.Data.DeployResult != "SUCCESS" {
			return resource.RetryableError(fmt.Errorf("domain deployment status is in progress, retrying"))
		}
		return nil
	})
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	log.Printf("resource.wangsu_cdn_domain.update success")
	return resourceCdnDomainRead(context, data, meta)
}
