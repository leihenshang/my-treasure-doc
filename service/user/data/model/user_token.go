package model

import "time"

// UserToken 用户表
type UserToken struct {
	BaseModel
	UserId       int64      `gorm:"column:user_id;type:bigint(20) unsigned;default:0;NOT NULL" json:"userId"` // 用户id
	Token        string     `json:"token" gorm:"column:token;type:varchar(100);comment:登陆token"`
	TokenExpire  time.Time  `json:"tokenExpire" gorm:"column:token_expire;type:datetime;comment:token超时时间"`
	LoginIp      string     `json:"loginIp" gorm:"column:login_ip;type:varchar(100);comment:最后登陆ip地址"`
	LoginTime    time.Time  `json:"loginTime" gorm:"column:login_time;type:datetime;comment:最后登陆时间"`
	LoginOutTime *time.Time `json:"loginOutTime" gorm:"column:login_out_time;type:datetime;comment:退出登陆时间"`
}

type UserTokens []*UserToken

func (m *UserToken) TableName() string {
	return "td_user_token"
}
