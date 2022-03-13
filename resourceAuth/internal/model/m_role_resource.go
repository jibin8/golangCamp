package model

import (
	"context"

	"self/internal/schema"
)

type IRoleResource interface {
	// 查询数据
	Query(ctx context.Context, params schema.RoleResourceQueryParam, opts ...schema.RoleResourceQueryOptions) (*schema.RoleResourceQueryResult, error)

}
