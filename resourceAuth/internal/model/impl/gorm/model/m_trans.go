package model

import (
	"context"
	icontext2 "self/internal/icontext"
	model2 "self/internal/model"

	"self/pkg/errors"
	"github.com/google/wire"
	"github.com/jinzhu/gorm"
)

var _ model2.ITrans = new(Trans)

// TransSet 注入Trans
var TransSet = wire.NewSet(wire.Struct(new(Trans), "*"), wire.Bind(new(model2.ITrans), new(*Trans)))

// Trans 事务管理
type Trans struct {
	DB *gorm.DB
}

// Exec 执行事务
func (a *Trans) Exec(ctx context.Context, fn func(context.Context) error) error {
	if _, ok := icontext2.FromTrans(ctx); ok {
		return fn(ctx)
	}

	err := a.DB.Transaction(func(db *gorm.DB) error {
		return fn(icontext2.NewTrans(ctx, db))
	})
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
