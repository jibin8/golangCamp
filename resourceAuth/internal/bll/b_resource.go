package bll

import (
	"context"

	"self/internal/schema"
)

type IResource interface {
	// 查询数据
	Query(ctx context.Context, params schema.ResourceQueryParam, opts ...schema.ResourceQueryOptions) (*schema.ResourceQueryResult, error)
}
