package bll

import (
	"context"

	"self/internal/schema"
)

type IAccount interface {
	// 查询数据
	Query(ctx context.Context, params schema.AccountQueryParam, opts ...schema.AccountQueryOptions) (*schema.AccountQueryResult, error)
	// 创建数据
	Create(ctx context.Context, item schema.AccountCreateParam) (string, error)
}
