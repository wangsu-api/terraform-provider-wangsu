package domain

import (
	"context"
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wangsuCommon "github.com/wangsu-api/terraform-provider-wangsu/wangsu/common"
	"github.com/wangsu-api/terraform-provider-wangsu/wangsu/services/waap"
	waapDomain "github.com/wangsu-api/wangsu-sdk-go/wangsu/waap/domain"
	"log"
	"time"
)

func ResourceWaapDomainCopy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWaapDomainCopyCreate,
		DeleteContext: resourceWaapDomainCopyDelete,
		ReadContext:   resourceWaapDomainCopyRead,
		UpdateContext: resourceWaapDomainCopyUpdate,

		Schema: map[string]*schema.Schema{
			"source_domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The reference hostname.",
			},
			"target_domains": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Hostname to be accessed.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceWaapDomainCopyCreate(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_waap_domain_copy.create")

	var diags diag.Diagnostics
	request := &waapDomain.UsingExistingHostnameToAddNewHostnameRequest{}
	if sourceDomain, ok := data.GetOk("source_domain"); ok {
		sourceDomainStr := sourceDomain.(string)
		request.SourceDomain = &sourceDomainStr
	}

	if targetDomains, ok := data.GetOk("target_domains"); ok {
		targetDomainsList := targetDomains.([]interface{})
		targetDomainsStr := make([]*string, len(targetDomainsList))
		for i, v := range targetDomainsList {
			str := v.(string)
			targetDomainsStr[i] = &str
		}
		request.TargetDomains = targetDomainsStr
	}

	var response *waapDomain.UsingExistingHostnameToAddNewHostnameResponse
	var err error
	err = resource.RetryContext(context, time.Duration(5)*time.Minute, func() *resource.RetryError {
		_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseWaapDomainClient().AddDomainByCopy(request)
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

	// 等待接入完成
	getRequest := &waapDomain.ListDomainInfoRequest{}
	getResponse := &waapDomain.ListDomainInfoResponse{}
	getRequest.SetDomainList(request.TargetDomains)
	for {
		err = resource.RetryContext(context, time.Duration(2)*time.Minute, func() *resource.RetryError {
			_, getResponse, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseWaapDomainClient().GetDomainList(getRequest)
			if err != nil {
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
			return diags
		}

		if getResponse == nil {
			return nil
		}
		if len(getResponse.Data) == 0 {
			time.Sleep(10 * time.Second)
			continue
		}
		successDomains := make([]string, 0, len(getResponse.Data))
		for _, v := range getResponse.Data {
			successDomains = append(successDomains, *v.Domain)
		}
		targetDomainsStr := make([]string, len(request.TargetDomains))
		for i, v := range request.TargetDomains {
			targetDomainsStr[i] = *v
		}
		waitDomains := waap.Difference(targetDomainsStr, successDomains)
		if len(waitDomains) == 0 {
			log.Printf("resource.wangsu_waap_domain_copy.create: all domains are ready")
			break
		}
		log.Printf("resource.wangsu_waap_domain_copy.create: waiting for domains %v", waitDomains)
		time.Sleep(10 * time.Second)
	}

	var ids = request.TargetDomains
	ids = append(ids, request.SourceDomain)
	idsStr := make([]string, len(ids))
	for i, v := range ids {
		idsStr[i] = *v
	}
	data.SetId(wangsuCommon.DataResourceIdsHash(idsStr))
	return diags
}

func resourceWaapDomainCopyDelete(context context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_waap_domain_copy.delete")

	var response *waapDomain.RemoveProtectedHostnameResponse
	var err error
	var diags diag.Diagnostics
	if targetDomains, ok := data.GetOk("target_domains"); ok {
		targetDomainsList := targetDomains.([]interface{})
		for _, v := range targetDomainsList {
			domain := v.(string)
			err = resource.RetryContext(context, time.Duration(5)*time.Minute, func() *resource.RetryError {
				request := &waapDomain.RemoveProtectedHostnameParameters{
					Domain: &domain,
				}
				_, response, err = meta.(wangsuCommon.ProviderMeta).GetAPIV3Conn().UseWaapDomainClient().DeleteDomain(request)
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
		}
	}
	return nil
}

func resourceWaapDomainCopyRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceWaapDomainCopyUpdate(_ context.Context, data *schema.ResourceData, _ interface{}) diag.Diagnostics {
	log.Printf("resource.wangsu_waap_domain_copy.update")
	var diags diag.Diagnostics
	err := errors.New("No changes allowed.")
	diags = append(diags, diag.FromErr(err)...)
	if data.HasChange("target_domains") {
		// 把domain强制刷回旧值，否则会有权限问题
		oldDomain, _ := data.GetChange("target_domains")
		_ = data.Set("target_domains", oldDomain)
	}
	if data.HasChange("source_domain") {
		// 把domain强制刷回旧值，否则会有权限问题
		oldDomain, _ := data.GetChange("source_domain")
		_ = data.Set("source_domain", oldDomain)
	}
	return diags
}
