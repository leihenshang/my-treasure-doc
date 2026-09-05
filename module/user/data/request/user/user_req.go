package user

// LoginRequest 用户登录请求
// CaptchaId 与 VerifyCode 由 GET /api/user/captcha 获取，captcha.enabled 为 true 时必填。
// 这里不使用 required，改由服务层按配置校验，便于关闭验证码时无需改动请求体。
type LoginRequest struct {
	Password   string `json:"password" binding:"required,min=6,max=100"`
	Account    string `json:"account" binding:"required,min=6,max=100"`
	CaptchaId  string `json:"captchaId" binding:"omitempty,max=100"`
	VerifyCode string `json:"verifyCode" binding:"omitempty,max=100"`
}
