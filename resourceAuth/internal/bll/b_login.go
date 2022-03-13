package bll

import (
	"context"
	"self/internal/schema"
)

type ILogin interface {
	// 登录验证
	Verify(ctx context.Context, appKey, username, password string) (*schema.Account, error)
	// 生成令牌
	GenerateToken(ctx context.Context, accountKey string, info interface{}) (*schema.LoginTokenInfo, error)
	// 鉴权
	Authenticate(ctx context.Context, appKey, accountKey, resourceType, feature, method string) (bool, error)
}
