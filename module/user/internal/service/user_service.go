package service

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"sync"
	"time"
	"unicode/utf8"

	"fastduck/treasure-doc/module/user/data/model"
	userReq "fastduck/treasure-doc/module/user/data/request/user"
	"fastduck/treasure-doc/module/user/global"
	"fastduck/treasure-doc/module/user/utils"

	"gorm.io/gorm"
)

type UserService struct{}

var userService *UserService

var userOnce = sync.Once{}

func NewUserService() *UserService {
	userOnce.Do(func() {
		userService = &UserService{}
	})

	if err := userService.RegisterRootUser(); err != nil {
		log.Fatalf("register root user failed: %v", err)
	}

	return userService
}

var rootUser = &model.User{
	Account:  "treasure-root",
	Email:    "treasure-root",
	Password: "treasure-root",
}

// RegisterRootUser 确保博客管理员默认账号存在，仅用于服务启动时初始化。
func (user *UserService) RegisterRootUser() error {
	if checkAccountIsDuplicate(rootUser.Account) {
		log.Printf("root account [%v] already exists, cancel registration\n", rootUser.Account)
		return nil
	}

	pwd, err := checkPasswordRule(rootUser.Password, rootUser.Password)
	if err != nil {
		return err
	}
	encryptedPwd, err := utils.PasswordEncrypt(pwd)
	if err != nil {
		return errors.New("加密密码失败")
	}
	if err := checkAccountRule(rootUser.Account, 8); err != nil {
		return err
	}

	u := &model.User{
		Nickname:   rootUser.Account,
		Account:    rootUser.Account,
		Email:      rootUser.Email,
		Password:   encryptedPwd,
		UserStatus: model.UserStatusAvailable,
		UserType:   model.UserTypeRoot,
	}
	if err := global.Db.Create(&u).Error; err != nil {
		global.Log.Errorf("failed to create root user: %v", err)
		return errors.New("创建默认管理员失败")
	}
	log.Printf("root user is registered,account is [%v],password is [%v], please update your password immediately\n", u.Account, rootUser.Password)
	return nil
}

// checkAccountIsDuplicate 检查账号是否重复
func checkAccountIsDuplicate(account string) bool {
	var u *model.User
	err := global.Db.Where("LOWER(account) = LOWER(?)", account).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	} else {
		global.Log.Errorf("failed to get user by account:%v", err)
	}
	return true
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
	reg := regexp.MustCompile(`^[a-zA-Z-_\d]*$`)
	if isAccord := reg.MatchString(account); !isAccord {
		return errors.New("账号必须为数字或英文")
	}

	return
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
func (user *UserService) UserLogin(r userReq.LoginRequest, clientIp string) (u *model.User, err error) {
	if len(r.Password) == 0 || len(r.Account) == 0 {
		return nil, errors.New("密码或账号(邮箱)不能为空")
	}

	// 先校验验证码，避免被自动化脚本用于撞库；验证码一次性消费。
	if global.GetConf() != nil && global.GetConf().Captcha.Enabled {
		if err = verifyCaptcha(r.CaptchaId, r.VerifyCode, true); err != nil {
			return nil, err
		}
	}

	err = global.Db.Where("LOWER(account) = LOWER(?) OR LOWER(email) = LOWER(?)", r.Account, r.Account).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New(fmt.Sprintf("账号 %s 没有找到", r.Account))
	}

	if u.UserStatus != model.UserStatusAvailable {
		return nil, errors.New("账号不可用或未激活")
	}

	if !utils.PasswordCompare(u.Password, r.Password) {
		return nil, errors.New("账号或密码错误")
	}

	var userTokens model.UserTokens
	if err = global.Db.Where("user_id = ?", u.Id).Order("created_at ASC").Find(&userTokens).Error; err != nil {
		global.Log.Errorf("failed to get user token:%v", err)
		return nil, errors.New("获取用户token失败")
	}

	tx := global.Db.Begin()
	if len(userTokens) == 3 {
		if err = tx.Delete(&userTokens[0]).Error; err != nil {
			global.Log.Errorf("failed to delete user token:%v", err)
			tx.Rollback()
			return nil, errors.New("删除用户token失败")
		}
	}

	userToken := &model.UserToken{
		Token:       utils.GenerateLoginToken(u.Id),
		TokenExpire: time.Now().Add(time.Hour * 24 * 7),
		LoginIp:     clientIp,
		LoginTime:   time.Now(),
		UserId:      u.Id,
	}

	if err = tx.Save(&userToken).Error; err != nil {
		global.Log.Errorf("failed to save user token:%v", err)
		tx.Rollback()
		return nil, errors.New("保存用户token失败")
	}

	tx.Commit()
	u.HiddenPwd().Token = userToken.Token
	return u, err
}

// UserLogout 用户退出登陆
// 幂等操作：按 user_id + token 使 token 失效即可，无需校验用户记录是否存在
// （开发模式下的 mock 超级管理员在 td_user 中并不存在，也能正常退出）。
func (user *UserService) UserLogout(userId string, token string) error {
	tx := global.Db.Begin()
	userToken := &model.UserToken{}
	if err := tx.Model(&userToken).Where("user_id = ? AND token = ?", userId, token).Update("login_out_time", time.Now()).Error; err != nil {
		global.Log.Errorf("failed to update user token login out time:%v", err)
		tx.Rollback()
		return errors.New("更新用户token信息失败")
	}

	if err := tx.Model(&userToken).Where("user_id = ? AND token = ?", userId, token).Delete(&model.UserToken{}).Error; err != nil {
		global.Log.Errorf("failed to delete user token:%v", err)
		tx.Rollback()
		return errors.New("删除用户token信息失败")
	}
	tx.Commit()

	return nil
}

// GetUserByToken 通过token获取用户
func GetUserByToken(token string) (u *model.User, err error) {
	now := time.Now()
	err = global.Db.Select("td_user.*").Joins("inner join td_user_token "+
		"on td_user_token.user_id = td_user.id AND td_user_token.token = ? "+
		"AND td_user_token.token_expire > ? AND td_user_token.deleted_at IS NULL", token, now).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		global.Log.Errorf("token : %s ,expire_time: %s  not found\n", token, now)
		return nil, errors.New("用户信息没有找到")
	}

	if !u.UserStatus.IsAvailable() {
		return nil, errors.New("用户不可用，请联系管理员")
	}

	u.HiddenPwd().Token = token
	return
}
