package bll

import (
	"context"
	"github.com/casbin/casbin"
	"github.com/google/wire"
	bll2 "self/internal/bll"
	model2 "self/internal/model"
	"self/internal/schema"
	"self/pkg/auth"
	"self/pkg/errors"
	"self/pkg/util"
)

var _ bll2.ILogin = (*Login)(nil)

// LoginSet 注入Login
var LoginSet = wire.NewSet(wire.Struct(new(Login), "*"), wire.Bind(new(bll2.ILogin), new(*Login)))

// Login 登录管理
type Login struct {
	Enforcer          *casbin.SyncedEnforcer
	Auth              auth.Auther
	AppModel          model2.IApp
	AccountModel      model2.IAccount
	AccountRoleModel  model2.IAccountRole
	RoleModel         model2.IRole
	RoleResourceModel model2.IRoleResource
	ResourceModel     model2.IResource
}

// Verify 登录验证
func (a *Login) Verify(ctx context.Context, appKey, username, password string) (*schema.Account, error) {
	// 检查是否是超级用户
	root := schema.GetRootAccount()
	if root.Username == username && root.Password == password {
		return root, nil
	}

	appResult, err := a.AppModel.Query(ctx, schema.AppQueryParam{
		PaginationParam: schema.PaginationParam{OnlyCount: true},
		AppKey:          appKey,
	})
	if err != nil {
		return nil, err
	} else if appResult.PageResult.Total == 0 {
		return nil, errors.New403Response(errors.ErrCodeAppNotExist, "该应用尚未接入")
	}

	result, err := a.AccountModel.Query(ctx, schema.AccountQueryParam{
		AppKey:   appKey,
		Username: username,
		Password: util.MD5HashString(password),
	})

	if err != nil {
		return nil, err
	} else if len(result.Data) == 0 {
		return nil, errors.New403Response(errors.ErrCodeLoginFailed, "登录失败，账号或密码错误")
	}

	item := result.Data[0]

	return item, nil
}

// GenerateToken 生成令牌
func (a *Login) GenerateToken(ctx context.Context, accountKey string, info interface{}) (*schema.LoginTokenInfo, error) {
	tokenInfo, err := a.Auth.GenerateToken(ctx, accountKey, info)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	item := &schema.LoginTokenInfo{
		Token:     tokenInfo.GetAccessToken(),
		ExpiresAt: tokenInfo.GetExpiresAt(),
	}
	return item, nil
}

// Authenticate 鉴权
func (a *Login) Authenticate(ctx context.Context, appKey, accountKey, resourceType, feature, method string) (bool, error) {
	b := a.Enforcer.Enforce(accountKey, appKey, resourceType, feature, method)
	return b, nil
}