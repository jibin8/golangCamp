package model

import (
	"context"

	"self/internal/schema"
)

type IApp interface {
	Query(ctx context.Context, params schema.AppQueryParam, opts ...schema.AppQueryOptions) (*schema.AppQueryResult, error)
}
