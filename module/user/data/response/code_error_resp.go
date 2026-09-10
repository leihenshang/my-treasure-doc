package response

const (
	SUCCESS ErrorCode = 0
	ERROR   ErrorCode = 1

	// CaptchaInvalid 图形验证码缺失、错误或已失效
	CaptchaInvalid ErrorCode = 40101
)
