package response

const SUCCESS ErrorCode = 0
const ERROR ErrorCode = 1

// DocIsEdited doc error,start with 10000

// DocIsEdited the doc has been edited
const DocIsEdited ErrorCode = 10000

// CaptchaInvalid 图形验证码缺失、错误或已失效
const CaptchaInvalid ErrorCode = 40101
