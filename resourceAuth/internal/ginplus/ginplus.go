package ginplus

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"self/internal/schema"
	"self/pkg/errors"
	"self/pkg/logger"
	"self/pkg/util"
)

// 定义上下文中的键
const (
	prefix           = "rj-auth"
	AccountKeyKey    = prefix + "/account-key"
	ReqBodyKey       = prefix + "/req-body"
	ResBodyKey       = prefix + "/res-body"
	LoggerReqBodyKey = prefix + "/logger-req-body"
)

// GetToken 获取用户令牌
func GetToken(c *gin.Context) string {
	var token string
	auth := c.GetHeader("RJ-Authorization")
	prefix := "Bearer "
	if strings.HasPrefix(auth, prefix) {
		token = auth[len(prefix):]
	} else {
		token = auth
	}
	return token
}

// GetAccountKey 获取账号标识
func GetAccountKey(c *gin.Context) string {
	return c.GetString(AccountKeyKey)
}

// SetAccountKey 设定账号标识
func SetAccountKey(c *gin.Context, accountKey string) {
	c.Set(AccountKeyKey, accountKey)
}

// ParseJSON 解析请求JSON
func ParseJSON(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		return errors.Wrap400Response(err, fmt.Sprintf("解析请求参数发生错误 - %s", err.Error()))
	}
	return nil
}

// ResSuccess 响应成功
func ResSuccess(c *gin.Context, v interface{}) {
	ResJSON(c, http.StatusOK, v)
}

// ResJSON 响应JSON数据
func ResJSON(c *gin.Context, status int, v interface{}) {
	buf, err := util.JSONMarshal(v)
	if err != nil {
		panic(err)
	}
	c.Set(ResBodyKey, buf)
	c.Data(status, "application/json; charset=utf-8", buf)
	c.Abort()
}

// ResError 响应错误
func ResError(c *gin.Context, err error, status ...int) {
	ctx := c.Request.Context()
	var res *errors.ResponseError
	if err != nil {
		if e, ok := err.(*errors.ResponseError); ok {
			res = e
		} else {
			res = errors.UnWrapResponse(errors.Wrap500Response(err, "服务器错误"))
		}
	} else {
		res = errors.UnWrapResponse(errors.ErrInternalServer)
	}

	if len(status) > 0 {
		res.StatusCode = status[0]
	}

	if err := res.ERR; err != nil {
		if status := res.StatusCode; status >= 400 && status < 500 {
			logger.StartSpan(ctx).Warnf(err.Error())
		} else if status >= 500 {
			logger.ErrorStack(ctx, err)
		}
	}

	eitem := schema.ErrorItem{
		Code:    res.Code,
		Message: res.Message,
	}
	ResJSON(c, res.StatusCode, eitem)
}
