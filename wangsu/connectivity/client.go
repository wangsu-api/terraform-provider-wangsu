package connectivity

import (
	"github.com/wangsu-api/wangsu-sdk-go/common"
	cdn "github.com/wangsu-api/wangsu-sdk-go/wangsu/cdn/domain"
	waapCustomizerule "github.com/wangsu-api/wangsu-sdk-go/wangsu/waap/customizerule"
	waapDomain "github.com/wangsu-api/wangsu-sdk-go/wangsu/waap/domain"
	waapRatelimit "github.com/wangsu-api/wangsu-sdk-go/wangsu/waap/ratelimit"
	waapWhitelist "github.com/wangsu-api/wangsu-sdk-go/wangsu/waap/whitelist"
)

type WangSuClient struct {
	Credential *common.Credential
	Protocol   string
	Domain     string

	cdnConn               *cdn.Client
	waapWhitelistConn     *waapWhitelist.Client
	waapCustomizeruleConn *waapCustomizerule.Client
	waapRatelimitConn     *waapRatelimit.Client
	waapDomainConn        *waapDomain.Client
}

func (me *WangSuClient) UseCdnClient() *cdn.Client {
	if me.cdnConn != nil {
		return me.cdnConn
	}

	me.cdnConn, _ = cdn.NewClient(me.Credential)

	return me.cdnConn
}

func (me *WangSuClient) UseWaapWhitelistClient() *waapWhitelist.Client {
	if me.waapWhitelistConn != nil {
		return me.waapWhitelistConn
	}

	me.waapWhitelistConn, _ = waapWhitelist.NewClient(me.Credential)

	return me.waapWhitelistConn
}

func (me *WangSuClient) UseWaapCustomizeruleClient() *waapCustomizerule.Client {
	if me.waapCustomizeruleConn != nil {
		return me.waapCustomizeruleConn
	}

	me.waapCustomizeruleConn, _ = waapCustomizerule.NewClient(me.Credential)

	return me.waapCustomizeruleConn
}

func (me *WangSuClient) UseWaapRatelimitClient() *waapRatelimit.Client {
	if me.waapRatelimitConn != nil {
		return me.waapRatelimitConn
	}

	me.waapRatelimitConn, _ = waapRatelimit.NewClient(me.Credential)

	return me.waapRatelimitConn
}

func (me *WangSuClient) UseWaapDomainClient() *waapDomain.Client {
	if me.waapDomainConn != nil {
		return me.waapDomainConn
	}

	me.waapDomainConn, _ = waapDomain.NewClient(me.Credential)

	return me.waapDomainConn
}
