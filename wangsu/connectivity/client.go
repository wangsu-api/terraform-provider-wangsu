package connectivity

import (
	appadomain "github.com/wangsu-api/wangsu-sdk-go/wangsu/appa/domain"
	cdn "github.com/wangsu-api/wangsu-sdk-go/wangsu/cdn/domain"
	"github.com/wangsu-api/wangsu-sdk-go/wangsu/common"
	edgeHostname "github.com/wangsu-api/wangsu-sdk-go/wangsu/edgehostname"
	monitorRule "github.com/wangsu-api/wangsu-sdk-go/wangsu/monitor/rule"
	"github.com/wangsu-api/wangsu-sdk-go/wangsu/policy"
	propertyConfig "github.com/wangsu-api/wangsu-sdk-go/wangsu/propertyconfig"
	"github.com/wangsu-api/wangsu-sdk-go/wangsu/ssl/certificate"
	userManage "github.com/wangsu-api/wangsu-sdk-go/wangsu/usermanage"
	waapCustomizerule "github.com/wangsu-api/wangsu-sdk-go/wangsu/waap/customizerule"
	waapDomain "github.com/wangsu-api/wangsu-sdk-go/wangsu/waap/domain"
	waapRatelimit "github.com/wangsu-api/wangsu-sdk-go/wangsu/waap/ratelimit"
	waapShareCustomizerule "github.com/wangsu-api/wangsu-sdk-go/wangsu/waap/share-customizerule"
	waapShareWhitelist "github.com/wangsu-api/wangsu-sdk-go/wangsu/waap/share-whitelist"
	waapWhitelist "github.com/wangsu-api/wangsu-sdk-go/wangsu/waap/whitelist"
)

type WangSuClient struct {
	Credential  *common.Credential
	HttpProfile *common.HttpProfile

	cdnConn                    *cdn.Client
	appaDomainConn             *appadomain.Client
	sslCertificateConn         *certificate.Client
	waapWhitelistConn          *waapWhitelist.Client
	waapCustomizeruleConn      *waapCustomizerule.Client
	waapRatelimitConn          *waapRatelimit.Client
	waapDomainConn             *waapDomain.Client
	waapShareWhitelistConn     *waapShareWhitelist.Client
	waapShareCustomizeruleConn *waapShareCustomizerule.Client
	monitorRuleConn            *monitorRule.Client
	policyConn                 *policy.Client
	userManageConn             *userManage.Client
	propertyConfigConn         *propertyConfig.Client
	edgeHostnameConn           *edgeHostname.Client
}

func (me *WangSuClient) UseCdnClient() *cdn.Client {
	if me.cdnConn != nil {
		return me.cdnConn
	}

	me.cdnConn, _ = cdn.NewClient(me.Credential, me.HttpProfile)

	return me.cdnConn
}

func (me *WangSuClient) UseAppaDomainClient() *appadomain.Client {
	if me.appaDomainConn != nil {
		return me.appaDomainConn
	}

	me.appaDomainConn, _ = appadomain.NewClient(me.Credential, me.HttpProfile)

	return me.appaDomainConn
}

func (me *WangSuClient) UseWaapWhitelistClient() *waapWhitelist.Client {
	if me.waapWhitelistConn != nil {
		return me.waapWhitelistConn
	}

	me.waapWhitelistConn, _ = waapWhitelist.NewClient(me.Credential, me.HttpProfile)

	return me.waapWhitelistConn
}

func (me *WangSuClient) UseWaapCustomizeruleClient() *waapCustomizerule.Client {
	if me.waapCustomizeruleConn != nil {
		return me.waapCustomizeruleConn
	}

	me.waapCustomizeruleConn, _ = waapCustomizerule.NewClient(me.Credential, me.HttpProfile)

	return me.waapCustomizeruleConn
}

func (me *WangSuClient) UseWaapRatelimitClient() *waapRatelimit.Client {
	if me.waapRatelimitConn != nil {
		return me.waapRatelimitConn
	}

	me.waapRatelimitConn, _ = waapRatelimit.NewClient(me.Credential, me.HttpProfile)

	return me.waapRatelimitConn
}

func (me *WangSuClient) UseWaapDomainClient() *waapDomain.Client {
	if me.waapDomainConn != nil {
		return me.waapDomainConn
	}

	me.waapDomainConn, _ = waapDomain.NewClient(me.Credential, me.HttpProfile)

	return me.waapDomainConn
}

func (me *WangSuClient) UseWaapShareWhitelistClient() *waapShareWhitelist.Client {
	if me.waapShareWhitelistConn != nil {
		return me.waapShareWhitelistConn
	}

	me.waapShareWhitelistConn, _ = waapShareWhitelist.NewClient(me.Credential, me.HttpProfile)

	return me.waapShareWhitelistConn
}

func (me *WangSuClient) UseWaapShareCustomizeruleClient() *waapShareCustomizerule.Client {
	if me.waapShareCustomizeruleConn != nil {
		return me.waapShareCustomizeruleConn
	}

	me.waapShareCustomizeruleConn, _ = waapShareCustomizerule.NewClient(me.Credential, me.HttpProfile)

	return me.waapShareCustomizeruleConn
}

func (me *WangSuClient) UseSslCertificateClient() *certificate.Client {
	if me.sslCertificateConn != nil {
		return me.sslCertificateConn
	}

	me.sslCertificateConn, _ = certificate.NewClient(me.Credential, me.HttpProfile)

	return me.sslCertificateConn
}

func (me *WangSuClient) UseMonitorRuleClient() *monitorRule.Client {
	if me.monitorRuleConn != nil {
		return me.monitorRuleConn
	}

	me.monitorRuleConn, _ = monitorRule.NewClient(me.Credential, me.HttpProfile)

	return me.monitorRuleConn
}
func (me *WangSuClient) UseUserManageClient() *userManage.Client {
	if me.userManageConn != nil {
		return me.userManageConn
	}

	me.userManageConn, _ = userManage.NewClient(me.Credential, me.HttpProfile)

	return me.userManageConn
}
func (me *WangSuClient) UsePolicyClient() *policy.Client {
	if me.policyConn != nil {
		return me.policyConn
	}

	me.policyConn, _ = policy.NewClient(me.Credential, me.HttpProfile)

	return me.policyConn
}

func (me *WangSuClient) UsePolicyAttachmentClient() *userManage.Client {
	if me.userManageConn != nil {
		return me.userManageConn
	}

	me.userManageConn, _ = userManage.NewClient(me.Credential, me.HttpProfile)

	return me.userManageConn
}

func (me *WangSuClient) UsePropertyConfigClient() *propertyConfig.Client {
	if me.propertyConfigConn != nil {
		return me.propertyConfigConn
	}

	me.propertyConfigConn, _ = propertyConfig.NewClient(me.Credential, me.HttpProfile)

	return me.propertyConfigConn
}

func (me *WangSuClient) UseEdgeHostnameClient() *edgeHostname.Client {
	if me.edgeHostnameConn != nil {
		return me.edgeHostnameConn
	}

	me.edgeHostnameConn, _ = edgeHostname.NewClient(me.Credential, me.HttpProfile)

	return me.edgeHostnameConn
}
