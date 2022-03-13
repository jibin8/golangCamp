package model

import (
	"context"

	"self/internal/schema"
)

type IRole interface {
	// 查询数据
	Query(ctx context.Context, params schema.RoleQueryParam, opts ...schema.RoleQueryOptions) (*schema.RoleQueryResult, error)
}
