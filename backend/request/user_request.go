package request

//UserRegisterRequest 用户注册请求
type UserRegisterRequest struct {
	Password       string `json:"password"`
	RepeatPassword string `json:"repeatPassword"`
	Account        string `json:"account"`
	Email          string `json:"email"`
	Code           string `json:"code"`
}

//UserLoginRequest 用户登录请求
type UserLoginRequest struct {
	Password   string `json:"password" binding:"required,min=6,max=100"`
	Account    string `json:"account" binding:"required,min=6,max=100"`
	VerifyCode string `json:"verifyCode" binding:"required,min=4,max=100"`
}

//UserLogoutRequest 用户退出登陆
type UserLogoutRequest struct {
	UserId uint64 `json:"userId" binding:"required"`
}

//UserProfileUpdateRequest 个人资料更新
type UserProfileUpdateRequest struct {
	UserId   uint64 `json:"userId" binding:""`
	NickName string `json:"nickName" binding:""`
	IconPath string `json:"iconPath" binding:""`
	Bio      string `json:"bio" binding:""`
	Mobile   string `json:"mobile" binding:""`
}
