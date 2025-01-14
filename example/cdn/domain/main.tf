terraform {
  required_providers {
    wangsu = {
      source = "registry.terraform.io/wangsu-api/wangsu"
    }
  }
}

provider "wangsu" {
  secret_id  = "my-secret-id"
  secret_key = "my-secret-key"
}

resource "wangsu_cdn_domain" "domain20240712001" {
  domain_name         = "20240712001.conftest.com"
  service_type        = "download"
  comment             = "test2"
  header_of_client_ip = "Cdn-Src-Ip"

  origin_config {
    origin_ips                 = "1.1.1.1"
    default_origin_host_header = "test.qq.com"
    use_range                  = "false"
    follow301                  = "false"
    follow302                  = "false"
    adv_src_setting {
      use_adv_src   = "true"
      detect_url    = "http://test.com/test2"
      detect_period = "60000"
      master_ips    = ["1.1.1.6", "1.1.1.2"]
      backup_ips    = ["2.2.2.6", "2.2.2.3"]
    }
  }
  ssl {
    use_ssl               = "true"
    ssl_certificate_id    = "1464893"
    backup_certificate_id = "1465017"
    #     gm_certificate_ids = ["1465027","1465028"]
    enable_ocsp      = "true"
    tls_version      = "TLSv1.1;TLSv1.2;GMTLS"
    ssl_cipher_suite = "AES128-SHA"
  }
  cache_time_behaviors {
    path_pattern                 = "*"
    cache_ttl                    = "2m"
    ignore_cache_control         = "true"
    is_respect_server            = "true"
    ignore_letter_case           = "true"
    reload_manage                = "ignore"
    priority                     = "11"
    ignore_authentication_header = "true"
  }
  cache_time_behaviors {
    path_pattern                 = "*"
    except_path_pattern          = "*.jpg"
    cache_ttl                    = "0"
    ignore_cache_control         = "true"
    is_respect_server            = "true"
    ignore_letter_case           = "true"
    reload_manage                = "ignore"
    priority                     = "12"
    ignore_authentication_header = "true"
  }
  cache_time_behaviors {
    custom_pattern               = "all"
    cache_ttl                    = "2m"
    ignore_cache_control         = "true"
    is_respect_server            = "true"
    ignore_letter_case           = "true"
    reload_manage                = "ignore"
    priority                     = "11"
    ignore_authentication_header = "true"
  }
  cache_time_behaviors {
    file_type                    = "png"
    cache_ttl                    = "2m"
    ignore_cache_control         = "true"
    is_respect_server            = "true"
    ignore_letter_case           = "true"
    reload_manage                = "ignore"
    priority                     = "11"
    ignore_authentication_header = "true"
  }
  cache_time_behaviors {
    custom_file_type             = "txt"
    cache_ttl                    = "2m"
    ignore_cache_control         = "true"
    is_respect_server            = "true"
    ignore_letter_case           = "true"
    reload_manage                = "ignore"
    priority                     = "11"
    ignore_authentication_header = "true"
  }
  cache_time_behaviors {
    specify_url_pattern          = "test.cachetimebehaviors.com"
    cache_ttl                    = "2m"
    ignore_cache_control         = "true"
    is_respect_server            = "true"
    ignore_letter_case           = "true"
    reload_manage                = "ignore"
    priority                     = "11"
    ignore_authentication_header = "true"
  }
  cache_time_behaviors {
    directory                    = "/test/cache/time/behaviors/"
    cache_ttl                    = "2s"
    ignore_cache_control         = "false"
    is_respect_server            = "false"
    ignore_letter_case           = "false"
    reload_manage                = "ignore"
    priority                     = "12"
    ignore_authentication_header = "false"
  }
  cache_key_rules {
    path_pattern        = "*"
    ignore_case         = "true"
    header_name         = "test-header1"
    parameter_of_header = "t1"
    priority            = "1"
  }
  cache_key_rules {
    specify_url             = "/test/specifyurl"
    full_match4_specify_url = "true"
    ignore_case             = "true"
    header_name             = "test-header2"
    parameter_of_header     = "t2"
    priority                = "2"
  }
  cache_key_rules {
    custom_pattern      = "all"
    ignore_case         = "false"
    header_name         = "test-header3"
    parameter_of_header = "t3"
    priority            = "3"
  }
  cache_key_rules {
    file_type           = "png"
    ignore_case         = "false"
    header_name         = "test-header4"
    parameter_of_header = "t4"
    priority            = "4"
  }
  cache_key_rules {
    custom_file_type    = "exe"
    ignore_case         = "false"
    header_name         = "test-header5"
    parameter_of_header = "t5"
    priority            = "5"
  }
  cache_key_rules {
    directory           = "/test6/"
    ignore_case         = "false"
    header_name         = "test-header6"
    parameter_of_header = "t6"
    priority            = "6"
  }
  query_string_settings {
    path_pattern        = "*"
    except_path_pattern = "abc.jpg"
    priority            = "1"
    ignore_letter_case  = "true"
    query_string_kept   = "a"
    source_with_query   = "false"
    source_key_kept     = "c"
  }
  query_string_settings {
    file_types           = "png"
    priority             = "2"
    ignore_letter_case   = "true"
    query_string_removed = "b"
    source_with_query    = "false"
    source_key_removed   = "d"
  }
  query_string_settings {
    custom_file_types    = "exe"
    priority             = "3"
    ignore_letter_case   = "true"
    query_string_removed = "b"
    source_with_query    = "false"
    source_key_removed   = "d"
  }
  query_string_settings {
    custom_pattern       = "all"
    priority             = "4"
    ignore_letter_case   = "true"
    query_string_removed = "b"
    source_with_query    = "false"
    source_key_removed   = "d"
  }
  query_string_settings {
    specify_url_pattern  = "/test/specify/url/pattern"
    priority             = "5"
    ignore_letter_case   = "true"
    query_string_removed = "b"
    source_with_query    = "false"
    source_key_removed   = "d"
  }
  query_string_settings {
    directories          = "/test/querystring/"
    priority             = "6"
    ignore_letter_case   = "false"
    query_string_removed = "b"
    source_with_query    = "true"
    source_key_removed   = "d"
  }
  query_string_settings {
    directories         = "/test/querystring2/"
    priority            = "7"
    ignore_letter_case  = "false"
    ignore_query_string = "true"
    source_with_query   = "true"
  }
  cache_by_resp_headers {
    response_header     = "test-header1"
    path_pattern        = ".*"
    except_path_pattern = "abc.png"
    response_value      = "resp1"
    ignore_letter_case  = "true"
    priority            = "1"
    is_respheader       = "true"
  }
  cache_by_resp_headers {
    response_header     = "test-header2"
    path_pattern        = ".*"
    except_path_pattern = "abc.jpg"
    response_value      = "resp2"
    ignore_letter_case  = "false"
    priority            = "2"
    is_respheader       = "false"
  }
  http_code_cache_rules {
    http_codes = ["200", "201"]
    cache_ttl  = "3600"
  }
  http_code_cache_rules {
    http_codes = ["500", "501"]
    cache_ttl  = "7200"
  }
  ignore_protocol_rules {
    path_pattern          = "*.jpg"
    except_path_pattern   = "abc.jpg"
    cache_ignore_protocol = "true"
    purge_ignore_protocol = "true"
  }
  ignore_protocol_rules {
    path_pattern          = "*.png"
    except_path_pattern   = "abc.png"
    cache_ignore_protocol = "false"
    purge_ignore_protocol = "false"
  }
  http2_settings {
    enable_http2            = "true"
    back_to_origin_protocol = "http2.0"
  }
  header_modify_rules {
    path_pattern          = "*.jpg"
    except_path_pattern   = "abc.jpg"
    request_method        = "GET"
    header_direction      = "visitor2cache"
    action                = "add"
    allow_regexp          = "false"
    header_name           = "test-header1"
    header_value          = "value1"
    request_header        = "r1"
    except_request_method = "POST"
    except_request_header = "r2"
    priority              = "1"
  }
  header_modify_rules {
    custom_pattern        = "all"
    request_method        = "GET"
    header_direction      = "visitor2cache"
    action                = "add"
    allow_regexp          = "false"
    header_name           = "test-header1"
    header_value          = "value1"
    request_header        = "r1"
    except_request_method = "POST"
    except_request_header = "r2"
    priority              = "2"
  }
  header_modify_rules {
    file_type             = "png"
    request_method        = "GET"
    header_direction      = "visitor2cache"
    action                = "add"
    allow_regexp          = "false"
    header_name           = "test-header1"
    header_value          = "value1"
    request_header        = "r1"
    except_request_method = "POST"
    except_request_header = "r2"
    except_file_type      = "jpg"
    priority              = "3"
  }
  header_modify_rules {
    custom_file_type      = "exe"
    request_method        = "GET"
    header_direction      = "visitor2cache"
    action                = "add"
    allow_regexp          = "false"
    header_name           = "test-header1"
    header_value          = "value1"
    request_header        = "r1"
    except_request_method = "POST"
    except_request_header = "r2"
    priority              = "4"
  }
  header_modify_rules {
    directory             = "/header/modify/rules/"
    request_method        = "GET"
    header_direction      = "visitor2cache"
    action                = "add"
    allow_regexp          = "false"
    header_name           = "test-header1"
    header_value          = "value1"
    request_header        = "r1"
    except_request_method = "POST"
    except_request_header = "r2"
    except_directory      = "/header/modify/rules2/"
    priority              = "5"
  }
  header_modify_rules {
    specify_url           = "/test/specify/url"
    request_method        = "GET"
    header_direction      = "visitor2cache"
    action                = "add"
    allow_regexp          = "false"
    header_name           = "test-header1"
    header_value          = "value1"
    request_header        = "r1"
    except_request_method = "POST"
    except_request_header = "r2"
    priority              = "6"
  }
  rewrite_rule_settings {
    path_pattern             = "*.jpg"
    except_path_pattern      = "abc.jpg"
    ignore_letter_case       = "true"
    publish_type             = "Cache"
    before_value             = "abc"
    after_value              = "def"
    rewrite_type             = "before"
    request_header           = "request_header"
    exception_request_header = "exception_request_header"
    priority                 = "1"
  }
  rewrite_rule_settings {
    custom_pattern           = "all"
    ignore_letter_case       = "true"
    publish_type             = "Cache"
    before_value             = "abc"
    after_value              = "def"
    rewrite_type             = "before"
    request_header           = "request_header"
    exception_request_header = "exception_request_header"
    priority                 = "2"
  }
  rewrite_rule_settings {
    directory                = "/rewrite_rule_settings/"
    ignore_letter_case       = "true"
    publish_type             = "Cache"
    before_value             = "abc"
    after_value              = "def"
    rewrite_type             = "before"
    request_header           = "request_header"
    exception_request_header = "exception_request_header"
    priority                 = "3"
  }
  rewrite_rule_settings {
    file_type                = "png"
    ignore_letter_case       = "true"
    publish_type             = "Cache"
    before_value             = "abc"
    after_value              = "def"
    rewrite_type             = "before"
    request_header           = "request_header"
    exception_request_header = "exception_request_header"
    priority                 = "4"
  }
  rewrite_rule_settings {
    custom_file_type         = "exe"
    ignore_letter_case       = "true"
    publish_type             = "Cache"
    before_value             = "abc"
    after_value              = "def"
    rewrite_type             = "before"
    request_header           = "request_header"
    exception_request_header = "exception_request_header"
    priority                 = "5"
  }
  back_to_origin_rewrite_rule {
    protocol                 = "https"
    port                     = "8443"
  }
}