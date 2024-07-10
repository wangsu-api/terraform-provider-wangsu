package common

import (
	"github.com/wangsu-api/terraform-provider-wangsu/wangsu/connectivity"
)

// ProviderMeta Provider 元信息
type ProviderMeta interface {
	// GetAPIV3Conn 返回访问云 API 的客户端连接对象
	GetAPIV3Conn() *connectivity.WangSuClient
}
