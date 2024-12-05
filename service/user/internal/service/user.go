package service

import (
	"errors"
	"fmt"
	"regexp"
	"time"
	"unicode/utf8"

	"fastduck/treasure-doc/service/user/data/model"
	"fastduck/treasure-doc/service/user/data/request/user"
	"fastduck/treasure-doc/service/user/global"
	"fastduck/treasure-doc/service/user/utils"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// UserRegister 用户注册
func UserRegister(r user.UserRegisterRequest) (u model.User, err error) {

	pwd, err := checkPasswordRule(r.Password, r.RePassword)
	if err != nil {
		return u, err
	}

	//对密码进行加密
	encryptedPwd, err := utils.PasswordEncrypt(pwd)
	if err != nil {
		return u, errors.New("加密密码失败")
	}

	//检查账号格式
	if err := checkAccountRule(r.Account, 8); err != nil {
		return u, err
	}

	//检查账号是否重复
	if checkAccountIsDuplicate(r.Account) {
		return u, errors.New("账号重复")
	}

	//检查邮箱是否重复
	if checkEmailIsDuplicate(r.Email) {
		return u, errors.New("邮箱重复")
	}

	if u.Nickname == "" {
		u.Nickname = r.Account
	}

	u.Account = r.Account
	u.Email = r.Email
	u.Password = encryptedPwd
	err = global.Db.Create(&u).Error
	//返回数据不显示密码
	u.Password = ""
	return u, err
}

// checkAccountIsDuplicate 检查账号是否重复
func checkAccountIsDuplicate(account string) bool {
	var user *model.User
	result := global.Db.Where("account = ?", account).First(&user)
	if result.RowsAffected > 0 {
		return true
	}

	return false
}

// checkAccountRule 检查账号规则
func checkAccountRule(account string, accountLen int) (err error) {
	if accountLen == 0 {
		global.Log.Fatal("accountLen is zero.")
		return errors.New("accountLen 设置错误")
	}

	if len(account) < accountLen {
		return errors.New(fmt.Sprintf("账号长度不能小于%d", accountLen))
	}

	//需要检查一下账号只能使用英文和数字
	reg := regexp.MustCompile(`^[a-zA-Z\d]*$`)
	if isAccord := reg.MatchString(account); !isAccord {
		return errors.New("账号必须为数字或英文")
	}

	return
}

// checkEmailIsDuplicate 检查邮箱是否重复
func checkEmailIsDuplicate(email string) bool {
	var user *model.User
	result := global.Db.Where("email = ?", email).First(&user)
	if result.RowsAffected > 0 {
		return true
	}

	return false
}

// checkPasswordRule 检查密码规则是否符合规则
func checkPasswordRule(password string, repeatPassword string) (string, error) {
	if utf8.RuneCountInString(password) < 8 {
		return "", errors.New("密码长度不能低于8位")
	}

	if password != repeatPassword {
		return "", errors.New("两次输入的密码不一致")
	}

	return password, nil
}

// UserLogin 用户登录
func UserLogin(r user.UserLoginRequest, clientIp string) (u model.User, err error) {
	if len(r.Password) == 0 || len(r.Account) == 0 {
		return u, errors.New("密码或账号(邮箱)不能为空")
	}

	result := global.Db.Where("account = ? OR email = ?", r.Account, r.Account).First(&u)
	if result.RowsAffected <= 0 {
		return u, errors.New(fmt.Sprintf("账号 %s 没有找到", r.Account))
	}

	//检查账号状态
	if u.UserStatus != model.UserStatusAvailable {
		return u, errors.New("账号不可用或未激活")
	}

	if utils.PasswordCompare(u.Password, r.Password) == false {
		return u, errors.New("账号或密码错误")
	}

	//下发token以及设置token过期时间
	u.Token = utils.GenerateLoginToken(u.Id)
	u.TokenExpire = time.Now().Add(time.Hour * 24 * 7)

	//设置登录时间记录用户登录ip
	u.LastLoginIp = clientIp
	u.LastLoginTime = time.Now()
	if err := global.Db.Select("LastLoginIp", "LastLoginTime", "Token", "TokenExpire").Save(&u).Error; err != nil {
		return u, errors.New("登录失败: 更新登录状态发生错误")
	}

	//密码置空
	u.Password = ""
	return u, err
}

// UserLogout 用户退出登陆
func UserLogout(userId int64) error {
	var userInfo model.User
	if errors.Is(global.Db.Where("id = ?", userId).First(&userInfo).Error, gorm.ErrRecordNotFound) {
		return errors.New("用户不存在")
	}

	userInfo.Token = ""
	userInfo.TokenExpire = utils.GetZeroDateTime()

	if err := global.Db.Save(&userInfo).Error; err != nil {
		global.Zap.Error("退出登陆，更新信息失败", zap.Any("dbErr", err))
		return errors.New("更新信息失败")
	}

	return nil
}

// UserProfileUpdate 更新用户个人资料
func UserProfileUpdate(profile user.UserProfileUpdateRequest, userId int64) (u model.User, err error) {
	if errors.Is(global.Db.Where("id = ?", userId).First(&u).Error, gorm.ErrRecordNotFound) {
		return u, errors.New("用户没有找到")
	}

	if err := global.Db.Model(&u).
		Select("NickName", "IconPath", "Bio", "Mobile").
		Updates(model.User{
			Nickname: profile.NickName,
			Avatar:   profile.NickName,
			Bio:      profile.Bio,
			Mobile:   profile.Mobile,
		}).
		Error; err != nil {
		return u, errors.New("更新个人资料失败")
	}

	return u, nil
}

// GetUserByToken 通过token获取用户
func GetUserByToken(token string) (u *model.User, err error) {

	now := time.Now().Format("2006-01-02 15:04:05")
	err = global.Db.Model(&model.User{}).Select(
		"user_type",
		"user_status",
		"token_expire",
		"token",
		"nickname",
		"mobile",
		"id",
		"email",
		"bio",
		"avatar",
		"account",
	).
		Where("token = ? AND token_expire > ?", token, now).
		First(&u).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		global.Log.Errorf("用户 token : %s ,expire_time: %s \n", token, now)
		return nil, errors.New("用户信息没有找到")
	}

	return
}

func ResetPwd(account string, pwd string) error {
	if _, err := checkPasswordRule(pwd, pwd); err != nil {
		return err
	}

	var u *model.User
	result := global.Db.Where("account = ?", account).First(&u)
	if result.RowsAffected <= 0 {
		return errors.New(fmt.Sprintf("账号 %s 没有找到", account))
	}

	//对密码进行加密
	encryptedPwd, err := utils.PasswordEncrypt(pwd)
	if err != nil {
		return errors.New("加密密码失败")
	}
	u.Password = encryptedPwd

	if err := global.Db.Select("Password").Save(&u).Error; err != nil {
		return errors.New("更新密码失败")
	}

	return nil
}
