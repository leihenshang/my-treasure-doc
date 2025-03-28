package model

// User 用户表
type User struct {
	BaseModel
	Nickname      string     `json:"nickname" gorm:"column:nickname;type:varchar(50);NOT NULL;comment:昵称;AUTO_INCREMENT"`
	Account       string     `json:"account" gorm:"column:account;type:varchar(100);default:'';comment:账号"`
	Email         string     `json:"email" gorm:"column:email;type:varchar(100);default:'';comment:邮箱"`
	Password      string     `json:"password" gorm:"column:password;type:varchar(100);default:'';comment:密码;"`
	UserType      UserType   `json:"userType" gorm:"column:user_type;type:tinyint(3) unsigned;default:1;NOT NULL;comment:1-普通用户,2管理员,100超级管理员"`
	UserStatus    UserStatus `json:"userStatus" gorm:"column:user_status;type:tinyint(3) unsigned;default:1;NOT NULL;comment:1-可用,2-不可用,3-未激活"`
	Mobile        string     `json:"mobile" gorm:"column:mobile;type:char(11);comment:手机号"`
	Avatar        string     `json:"avatar" gorm:"column:avatar;type:varchar(500);comment:头像地址"`
	Bio           string     `json:"bio" gorm:"column:bio;type:varchar(200);comment:个人说明"`
	CurrentRoomId string     `json:"currentRoomId" gorm:"column:current_room_id;type:varchar(100);NOT NULL;default:'';comment:当前所在房间"`
	Token         string     `json:"token" gorm:"-"`
}

type Users []*User

type UserType int

const (
	UserTypeUser  UserType = 1
	UserTypeAdmin UserType = 2
	UserTypeRoot  UserType = 100
)

type UserStatus int

const (
	// UserStatusAvailable 用户状态-可用
	UserStatusAvailable UserStatus = 1
	// UserStatusInvalid 用户状态-无效
	UserStatusInvalid UserStatus = 2
	// UserStatusNotActivated 用户状态-不活跃
	UserStatusNotActivated UserStatus = 3
)

func (u UserStatus) IsAvailable() bool {
	return u == 1
}

func (u UserType) IsUser() bool {
	return u == 1
}

func (u *User) TableName() string {
	return "td_user"
}

func (u *User) HiddenPwd() *User {
	u.Password = ""
	return u
}
