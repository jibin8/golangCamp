package bll

import (
	"context"

	"self/internal/schema"
)

type IApp interface {
	// 查询数据
	Query(ctx context.Context, params schema.AppQueryParam, opts ...schema.AppQueryOptions) (*schema.AppQueryResult, error)
}
