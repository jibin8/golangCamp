package schema

// LoginParam 登录参数
type LoginParam struct {
	Username    string 	`json:"username" binding:"required"`    	// 用户名
	Password    string 	`json:"password" binding:"required"`     	// 密码(md5加密)
	AppKey    	string 	`json:"appKey" binding:"required"`     		// 应用标识
	//CaptchaID   string `json:"captcha_id" binding:"required"`   // 验证码ID
	//CaptchaCode string `json:"captcha_code" binding:"required"` // 验证码
}

// LoginTokenInfo 登录令牌信息
type LoginTokenInfo struct {
	Token 		string 	`json:"token"` 		// 访问令牌
	ExpiresAt   int64  	`json:"expiresAt"`  // 令牌过期时间戳
}

// AuthenticateParam 鉴权信息
type AuthenticateParam struct {
	Feature   		string 	`json:"feature" binding:"required"`   	// 资源特征值
	ResourceType   	string 	`json:"resourceType,default=api"`   	// 资源类型（api|...）
	Method 			string 	`json:"method,default="` 				// 请求类型
}

// AuthenticateInfo 授权信息
type AuthenticateInfo struct {
	Grant   		int64 	`json:"grant"`   		// 鉴权结果：0不通过 1通过
}
