package model

import (
	"context"

	"self/internal/schema"
)

type IAccount interface {
	Query(ctx context.Context, params schema.AccountQueryParam, opts ...schema.AccountQueryOptions) (*schema.AccountQueryResult, error)
	Create(ctx context.Context, item schema.Account) error
}
