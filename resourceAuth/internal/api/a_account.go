package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	bll2 "self/internal/bll"
	ginplus2 "self/internal/ginplus"
	"self/internal/schema"
)

// UserSet 注入User
var AccountSet = wire.NewSet(wire.Struct(new(Account), "*"))

// Account 账号管理
type Account struct {
	AccountBll bll2.IAccount
}

// Create 创建数据
func (a *Account) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.AccountCreateParam
	if err := ginplus2.ParseJSON(c, &item); err != nil {
		ginplus2.ResError(c, err)
		return
	}

	accountKey, err := a.AccountBll.Create(ctx, item)
	if err != nil {
		ginplus2.ResError(c, err)
		return
	}
	ginplus2.ResSuccess(c, accountKey)
}