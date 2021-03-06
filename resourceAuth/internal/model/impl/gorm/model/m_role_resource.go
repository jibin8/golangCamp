package model

import (
	"context"
	model2 "self/internal/model"
	entity2 "self/internal/model/impl/gorm/entity"

	"self/internal/schema"
	"self/pkg/errors"
	"github.com/google/wire"
	"github.com/jinzhu/gorm"
)

var _ model2.IRoleResource = (*RoleResource)(nil)

// RoleResourceSet 注入RoleResource
var RoleResourceSet = wire.NewSet(wire.Struct(new(RoleResource), "*"), wire.Bind(new(model2.IRoleResource), new(*RoleResource)))

// RoleResource 角色资源存储
type RoleResource struct {
	DB *gorm.DB
}

func (a *RoleResource) getQueryOption(opts ...schema.RoleResourceQueryOptions) schema.RoleResourceQueryOptions {
	var opt schema.RoleResourceQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// Query 查询数据
func (a *RoleResource) Query(ctx context.Context, params schema.RoleResourceQueryParam, opts ...schema.RoleResourceQueryOptions) (*schema.RoleResourceQueryResult, error) {
	opt := a.getQueryOption(opts...)

	db := entity2.GetRoleResourceDB(ctx, a.DB)
	if v := params.RoleID; v != "" {
		db = db.Where("roleID = ?", v)
	}
	if v := params.RoleIDs; len(v) > 0 {
		db = db.Where("roleID IN (?)", v)
	}

	opt.OrderFields = append(opt.OrderFields, schema.NewOrderField("id", schema.OrderByDESC))
	db = db.Order(ParseOrder(opt.OrderFields))

	var list entity2.RoleResources
	pr, err := WrapPageQuery(ctx, db, params.PaginationParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.RoleResourceQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaRoleResources(),
	}

	return qr, nil
}

//// Get 查询指定数据
//func (a *RoleMenu) Get(ctx context.Context, id string, opts ...schema.RoleMenuQueryOptions) (*schema.RoleMenu, error) {
//	db := entity.GetRoleMenuDB(ctx, a.DB).Where("id=?", id)
//	var item entity.RoleMenu
//	ok, err := FindOne(ctx, db, &item)
//	if err != nil {
//		return nil, errors.WithStack(err)
//	} else if !ok {
//		return nil, nil
//	}
//
//	return item.ToSchemaRoleMenu(), nil
//}
//
//// Create 创建数据
//func (a *RoleMenu) Create(ctx context.Context, item schema.RoleMenu) error {
//	eitem := entity.SchemaRoleMenu(item).ToRoleMenu()
//	result := entity.GetRoleMenuDB(ctx, a.DB).Create(eitem)
//	if err := result.Error; err != nil {
//		return errors.WithStack(err)
//	}
//	return nil
//}
//
//// Update 更新数据
//func (a *RoleMenu) Update(ctx context.Context, id string, item schema.RoleMenu) error {
//	eitem := entity.SchemaRoleMenu(item).ToRoleMenu()
//	result := entity.GetRoleMenuDB(ctx, a.DB).Where("id=?", id).Updates(eitem)
//	if err := result.Error; err != nil {
//		return errors.WithStack(err)
//	}
//	return nil
//}
//
//// Delete 删除数据
//func (a *RoleMenu) Delete(ctx context.Context, id string) error {
//	result := entity.GetRoleMenuDB(ctx, a.DB).Where("id=?", id).Delete(entity.RoleMenu{})
//	if err := result.Error; err != nil {
//		return errors.WithStack(err)
//	}
//	return nil
//}
//
//// DeleteByRoleID 根据角色ID删除数据
//func (a *RoleMenu) DeleteByRoleID(ctx context.Context, roleID string) error {
//	result := entity.GetRoleMenuDB(ctx, a.DB).Where("role_id=?", roleID).Delete(entity.RoleMenu{})
//	if err := result.Error; err != nil {
//		return errors.WithStack(err)
//	}
//	return nil
//}
