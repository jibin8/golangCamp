package model

import (
	"context"

	"self/internal/schema"
)

type IAccountSystem interface {
	Query(ctx context.Context, params schema.AccountSystemQueryParam, opts ...schema.QueryOptions) (*schema.AccountSystemQueryResult, error)
}
