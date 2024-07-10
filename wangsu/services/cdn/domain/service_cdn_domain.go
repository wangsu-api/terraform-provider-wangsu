package domain

import "github.com/wangsu-api/terraform-provider-wangsu/wangsu/connectivity"

func NewCdnService(client *connectivity.WangSuClient) CdnService {
	return CdnService{client: client}
}

type CdnService struct {
	client *connectivity.WangSuClient
}
