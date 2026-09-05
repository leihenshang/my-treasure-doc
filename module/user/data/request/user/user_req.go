package user

// RegisterRequest 用户注册请求
type RegisterRequest struct {
	Password   string `json:"password"`
	RePassword string `json:"rePassword"`
	Account    string `json:"account"`
	Email      string `json:"email"`
	Code       string `json:"code"`
}

// LoginRequest 用户登录请求
// CaptchaId 与 VerifyCode 由 GET /api/user/captcha 获取，captcha.enabled 为 true 时必填。
// 这里不使用 required，改由服务层按配置校验，便于关闭验证码时无需改动请求体。
type LoginRequest struct {
	Password   string `json:"password" binding:"required,min=6,max=100"`
	Account    string `json:"account" binding:"required,min=6,max=100"`
	CaptchaId  string `json:"captchaId" binding:"omitempty,max=100"`
	VerifyCode string `json:"verifyCode" binding:"omitempty,max=100"`
}

// UpdateRequest 个人资料更新
type UpdateRequest struct {
	NickName string `json:"nickName" binding:"max=40"`
	IconPath string `json:"iconPath" binding:"max=100"`
	Bio      string `json:"bio" binding:"max=200"`
	Mobile   string `json:"mobile" binding:"omitempty,len=11"`
}
