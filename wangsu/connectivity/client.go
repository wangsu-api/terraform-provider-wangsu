package connectivity

import (
	appadomain "github.com/wangsu-api/wangsu-sdk-go/wangsu/appa/domain"
	cdn "github.com/wangsu-api/wangsu-sdk-go/wangsu/cdn/domain"
	"github.com/wangsu-api/wangsu-sdk-go/wangsu/common"
	edgeHostname "github.com/wangsu-api/wangsu-sdk-go/wangsu/edgehostname"
	monitorRule "github.com/wangsu-api/wangsu-sdk-go/wangsu/monitor/rule"
	"github.com/wangsu-api/wangsu-sdk-go/wangsu/policy"
	propertyConfig "github.com/wangsu-api/wangsu-sdk-go/wangsu/propertyconfig"
	"github.com/wangsu-api/wangsu-sdk-go/wangsu/securitypolicy"
	"github.com/wangsu-api/wangsu-sdk-go/wangsu/ssl/certificate"
	userManage "github.com/wangsu-api/wangsu-sdk-go/wangsu/usermanage"
	waapBot "github.com/wangsu-api/wangsu-sdk-go/wangsu/waap/bot"
	waapBotSceneWhitelist "github.com/wangsu-api/wangsu-sdk-go/wangsu/waap/bot-scene-whitelist"
	waapCustomizerule "github.com/wangsu-api/wangsu-sdk-go/wangsu/waap/customizerule"
	waapDDoSProtection "github.com/wangsu-api/wangsu-sdk-go/wangsu/waap/ddosprotection"
	waapDomain "github.com/wangsu-api/wangsu-sdk-go/wangsu/waap/domain"
	waapPreDeploy "github.com/wangsu-api/wangsu-sdk-go/wangsu/waap/predeploy"
	waapRatelimit "github.com/wangsu-api/wangsu-sdk-go/wangsu/waap/ratelimit"
	waapShareCustomizeBot "github.com/wangsu-api/wangsu-sdk-go/wangsu/waap/share-customizebot"
	waapShareCustomizerule "github.com/wangsu-api/wangsu-sdk-go/wangsu/waap/share-customizerule"
	waapShareWhitelist "github.com/wangsu-api/wangsu-sdk-go/wangsu/waap/share-whitelist"
	waapWAF "github.com/wangsu-api/wangsu-sdk-go/wangsu/waap/waf"
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
	waapWAFConn                *waapWAF.Client
	waapBotConn                *waapBot.Client
	waapDDoSProtectionConn     *waapDDoSProtection.Client
	waapPreDeployConn          *waapPreDeploy.Client
	monitorRuleConn            *monitorRule.Client
	policyConn                 *policy.Client
	userManageConn             *userManage.Client
	propertyConfigConn         *propertyConfig.Client
	edgeHostnameConn           *edgeHostname.Client
	securityPolicyConn         *securitypolicy.Client
	waapBotSceneWhitelistConn  *waapBotSceneWhitelist.Client
	waapShareCustomizeBotConn  *waapShareCustomizeBot.Client
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

func (me *WangSuClient) UseWaapBotSceneWhiteListClient() *waapBotSceneWhitelist.Client {
	if me.waapBotConn != nil {
		return me.waapBotSceneWhitelistConn
	}

	me.waapBotSceneWhitelistConn, _ = waapBotSceneWhitelist.NewClient(me.Credential, me.HttpProfile)

	return me.waapBotSceneWhitelistConn
}

func (me *WangSuClient) UseWaapShareCustomizeBotClient() *waapShareCustomizeBot.Client {
	if me.waapShareCustomizeBotConn != nil {
		return me.waapShareCustomizeBotConn
	}

	me.waapShareCustomizeBotConn, _ = waapShareCustomizeBot.NewClient(me.Credential, me.HttpProfile)

	return me.waapShareCustomizeBotConn
}

func (me *WangSuClient) UseWaapPreDeployClient() *waapPreDeploy.Client {
	if me.waapPreDeployConn != nil {
		return me.waapPreDeployConn
	}

	me.waapPreDeployConn, _ = waapPreDeploy.NewClient(me.Credential, me.HttpProfile)

	return me.waapPreDeployConn
}

func (me *WangSuClient) UseWaapWAFClient() *waapWAF.Client {
	if me.waapWAFConn != nil {
		return me.waapWAFConn
	}

	me.waapWAFConn, _ = waapWAF.NewClient(me.Credential, me.HttpProfile)

	return me.waapWAFConn
}

func (me *WangSuClient) UseWaapBotClient() *waapBot.Client {
	if me.waapBotConn != nil {
		return me.waapBotConn
	}

	me.waapBotConn, _ = waapBot.NewClient(me.Credential, me.HttpProfile)

	return me.waapBotConn
}

func (me *WangSuClient) UseWaapDDoSProtectionClient() *waapDDoSProtection.Client {
	if me.waapDDoSProtectionConn != nil {
		return me.waapDDoSProtectionConn
	}

	me.waapDDoSProtectionConn, _ = waapDDoSProtection.NewClient(me.Credential, me.HttpProfile)

	return me.waapDDoSProtectionConn
}

func (me *WangSuClient) UseSecurityPolicyClient() *securitypolicy.Client {
	if me.securityPolicyConn != nil {
		return me.securityPolicyConn
	}

	me.securityPolicyConn, _ = securitypolicy.NewClient(me.Credential, me.HttpProfile)

	return me.securityPolicyConn
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
