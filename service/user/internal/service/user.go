package service

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"sync"
	"time"
	"unicode/utf8"

	"fastduck/treasure-doc/service/user/data/model"
	userReq "fastduck/treasure-doc/service/user/data/request/user"
	"fastduck/treasure-doc/service/user/global"
	"fastduck/treasure-doc/service/user/utils"

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

func (user *UserService) RegisterRootUser() error {
	regRequest := &userReq.RegisterRequest{
		Password:   rootUser.Password,
		RePassword: rootUser.Password,
		Account:    rootUser.Account,
		Email:      rootUser.Email,
	}
	if checkAccountIsDuplicate(regRequest.Account) {
		log.Printf("root account [%v] already existes,cancel registration\n", regRequest.Account)
	} else {
		if _, err := userService.UserRegister(regRequest); err != nil {
			return err
		}
		log.Printf("root user is registered,account is [%v],password is [%v],"+
			"please update your password immediately\n", regRequest.Account, regRequest.Password)
	}
	return nil
}

// UserRegister 用户注册
func (user *UserService) UserRegister(r *userReq.RegisterRequest) (u *model.User, err error) {
	pwd, err := checkPasswordRule(r.Password, r.RePassword)
	if err != nil {
		return nil, err
	}

	encryptedPwd, err := utils.PasswordEncrypt(pwd)
	if err != nil {
		return nil, errors.New("加密密码失败")
	}

	if err := checkAccountRule(r.Account, 8); err != nil {
		return nil, err
	}

	if checkAccountIsDuplicate(r.Account) {
		return nil, errors.New("账号重复")
	}

	if checkEmailIsDuplicate(r.Email) {
		return nil, errors.New("邮箱重复")
	}

	u = &model.User{}
	u.Nickname = r.Account
	u.Account = r.Account
	u.Email = r.Email
	u.Password = encryptedPwd

	if r.Account == rootUser.Account && r.Password == rootUser.Password {
		u.UserStatus = model.UserStatusAvailable
		u.UserType = model.UserTypeRoot
	}

	trans := global.Db.Begin()
	err = trans.Create(&u).Error
	if err != nil {
		trans.Rollback()
		global.Log.Errorf("failed to create userReq:%v", err)
		return nil, errors.New("注册失败")
	}
	err = trans.Create(&model.Room{
		Name:   "个人空间",
		UserId: u.Id,
	}).Error
	if err != nil {
		trans.Rollback()
		global.Log.Errorf("failed to create room:%v", err)
		return nil, errors.New("创建空间失败")
	}
	trans.Commit()

	u.Password = ""
	return u, err
}

// checkAccountIsDuplicate 检查账号是否重复
func checkAccountIsDuplicate(account string) bool {
	var u *model.User
	err := global.Db.Where("account = ?", account).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	} else {
		global.Log.Errorf("failed to get userReq:%v", err)
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

// checkEmailIsDuplicate 检查邮箱是否重复
func checkEmailIsDuplicate(email string) bool {
	var u *model.User
	err := global.Db.Where("email = ?", email).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	} else {
		global.Log.Errorf("failed to get userReq from email:%v", err)
	}

	return true
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

	err = global.Db.Where("account = ? OR email = ?", r.Account, r.Account).First(&u).Error
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
		global.Log.Errorf("failed to get userReq token:%v", err)
		return nil, errors.New("获取用户token失败")
	}

	tx := global.Db.Begin()
	if len(userTokens) == 3 {
		if err = tx.Delete(&userTokens[0]).Error; err != nil {
			global.Log.Errorf("failed to delete userReq token:%v", err)
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
		global.Log.Errorf("failed to save userReq token:%v", err)
		tx.Rollback()
		return nil, errors.New("保存用户token失败")
	}

	tx.Commit()
	u.HiddenPwd().Token = userToken.Token
	return u, err
}

// UserLogout 用户退出登陆
func (user *UserService) UserLogout(userId string, token string) error {
	var userInfo model.User
	if err := global.Db.Where("id = ?", userId).First(&userInfo).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("用户不存在")
	} else if err != nil {
		global.Log.Errorf("userReq logout error:%v", err)
		return errors.New("查询用户信息失败")
	}

	tx := global.Db.Begin()
	userToken := &model.UserToken{}
	if err := tx.Model(&userToken).Where("user_id = ? AND token = ?", userId, token).Update("login_out_time", time.Now()).Error; err != nil {
		global.Log.Errorf("failed to update userReq token login out time:%v", err)
		tx.Rollback()
		return errors.New("更新用户token信息失败")
	}

	if err := tx.Model(&userToken).Where("user_id = ? AND token = ?", userId, token).Delete(&model.UserToken{}).Error; err != nil {
		global.Log.Errorf("failed to delete userReq token:%v", err)
		tx.Rollback()
		return errors.New("删除用户token信息失败")
	}
	tx.Commit()

	return nil
}

// UserProfileUpdate 更新用户个人资料
func (user *UserService) UserProfileUpdate(profile userReq.UpdateRequest, userId string) (u model.User, err error) {
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

func ResetPwd(account string, pwd string, rePwd string) error {
	if _, err := checkPasswordRule(pwd, rePwd); err != nil {
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

	tx := global.Db.Begin()
	userToken := &model.UserToken{}
	if err := tx.Model(&userToken).Where("user_id = ?", u.Id).Delete(&model.UserToken{}).Error; err != nil {
		global.Log.Errorf("failed to delete userReq token:%v", err)
		tx.Rollback()
		return errors.New("删除用户token信息失败")
	}
	tx.Commit()

	return nil
}
