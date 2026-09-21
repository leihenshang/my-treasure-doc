package service

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"fastduck/treasure-doc/module/user/data/model"
	"fastduck/treasure-doc/module/user/global"
	"fastduck/treasure-doc/module/user/utils"

	"github.com/glebarez/sqlite"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// newUserServiceTestDB 用真实 SQLite 搭建用户表，并把 global.Log 的输出捕获到 buffer。
// 捕获日志的目的不是断言「写了什么」，而是断言「不该写什么」——例如账号已存在时
// 曾被错误打印成 failed to get user by account:<nil>。
func newUserServiceTestDB(t *testing.T) *bytes.Buffer {
	t.Helper()

	dsn := "file:" + filepath.Join(t.TempDir(), "user_service.db") + "?_busy_timeout=5000&_journal_mode=WAL&_foreign_keys=on"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("打开 SQLite 失败: %v", err)
	}
	if err := db.AutoMigrate(&model.User{}, &model.UserToken{}); err != nil {
		t.Fatalf("AutoMigrate 失败: %v", err)
	}

	logs := &bytes.Buffer{}
	prevDb, prevLog := global.Db, global.Log
	global.Db = db
	global.Log = zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(logs),
		zapcore.DebugLevel,
	)).Sugar()
	t.Cleanup(func() {
		global.Db, global.Log = prevDb, prevLog
		if sqlDB, err := db.DB(); err == nil {
			_ = sqlDB.Close()
		}
	})

	return logs
}

func mustCreateUser(t *testing.T, account, password string, userType model.UserType) *model.User {
	t.Helper()

	encrypted, err := utils.PasswordEncrypt(password)
	if err != nil {
		t.Fatalf("加密密码失败: %v", err)
	}
	u := &model.User{
		Nickname:   account,
		Account:    account,
		Email:      account,
		Password:   encrypted,
		UserStatus: model.UserStatusAvailable,
		UserType:   userType,
	}
	if err := global.Db.Create(u).Error; err != nil {
		t.Fatalf("创建用户失败: %v", err)
	}
	return u
}

// 账号已存在是查询成功的结果（err == nil），不能落到错误分支打印 <nil> 错误。
func TestCheckAccountIsDuplicateDoesNotLogForExistingAccount(t *testing.T) {
	logs := newUserServiceTestDB(t)

	duplicated, err := checkAccountIsDuplicate(DefaultAdminAccount)
	if err != nil {
		t.Fatalf("账号不存在时不应返回错误: %v", err)
	}
	if duplicated {
		t.Fatal("账号不存在时 duplicated 应为 false")
	}

	mustCreateUser(t, DefaultAdminAccount, "InitPass123", model.UserTypeRoot)

	// 大小写不敏感，与登录查询保持一致
	duplicated, err = checkAccountIsDuplicate(strings.ToUpper(DefaultAdminAccount))
	if err != nil {
		t.Fatalf("账号已存在不是错误: %v", err)
	}
	if !duplicated {
		t.Fatal("账号存在时 duplicated 应为 true")
	}

	if output := logs.String(); strings.Contains(output, "failed to get user by account") {
		t.Fatalf("账号存在不应记录查询失败日志: %s", output)
	}
}

// 默认管理员已存在时，启动注册应静默跳过：既不重复创建，也不打印查询失败。
func TestRegisterRootUserSkipsExistingAccountQuietly(t *testing.T) {
	logs := newUserServiceTestDB(t)
	mustCreateUser(t, DefaultAdminAccount, "InitPass123", model.UserTypeRoot)

	if err := (&UserService{}).RegisterRootUser(); err != nil {
		t.Fatalf("账号已存在时注册应返回 nil: %v", err)
	}

	var count int64
	if err := global.Db.Model(&model.User{}).
		Where("LOWER(account) = LOWER(?)", DefaultAdminAccount).Count(&count).Error; err != nil {
		t.Fatalf("统计用户失败: %v", err)
	}
	if count != 1 {
		t.Fatalf("不应重复创建默认账号，当前数量 %d", count)
	}
	if output := logs.String(); strings.Contains(output, "failed to get user by account") {
		t.Fatalf("账号已存在不应记录查询失败日志: %s", output)
	}
}

// 重置默认管理员密码：新密码生效、旧密码失效，清除强制改密标记并清空其它登录态。
func TestResetDefaultAdminPasswordReplacesPasswordAndRevokesTokens(t *testing.T) {
	newUserServiceTestDB(t)

	u := mustCreateUser(t, DefaultAdminAccount, "OldPass123", model.UserTypeRoot)
	// 默认管理员首次启动会被置位强制改密；重置密码后必须清掉，否则登录后会被
	// RequireAdmin 以 40301 拦在管理接口外，表现成「重置密码无效」。
	if err := global.Db.Model(&model.User{}).Where("id = ?", u.Id).
		Update("require_pwd_reset", true).Error; err != nil {
		t.Fatalf("设置强制改密标记失败: %v", err)
	}
	if err := global.Db.Create(&model.UserToken{
		UserId:      u.Id,
		Token:       "old-token",
		TokenExpire: time.Now().Add(time.Hour),
	}).Error; err != nil {
		t.Fatalf("创建登录态失败: %v", err)
	}

	if err := (&UserService{}).ResetDefaultAdminPassword("NewPass123"); err != nil {
		t.Fatalf("重置密码失败: %v", err)
	}

	var got model.User
	if err := global.Db.Where("id = ?", u.Id).First(&got).Error; err != nil {
		t.Fatalf("读取用户失败: %v", err)
	}
	if !utils.PasswordCompare(got.Password, "NewPass123") {
		t.Fatal("新密码未生效")
	}
	if utils.PasswordCompare(got.Password, "OldPass123") {
		t.Fatal("旧密码在重置后仍然可用")
	}
	if got.RequirePwdReset {
		t.Fatal("重置密码后应清除 require_pwd_reset")
	}

	var tokens int64
	if err := global.Db.Model(&model.UserToken{}).Where("user_id = ?", u.Id).Count(&tokens).Error; err != nil {
		t.Fatalf("统计登录态失败: %v", err)
	}
	if tokens != 0 {
		t.Fatalf("重置后其它登录态应失效，剩余 %d", tokens)
	}
}

// 重置不存在的账号必须返回错误，而不是静默成功。
func TestResetPasswordRejectsUnknownAccount(t *testing.T) {
	newUserServiceTestDB(t)

	if err := (&UserService{}).ResetDefaultAdminPassword("NewPass123"); err == nil {
		t.Fatal("账号不存在时应返回错误")
	}
}

// 新密码仍走统一规则校验，过短的密码不允许重置。
func TestResetPasswordRejectsShortPassword(t *testing.T) {
	newUserServiceTestDB(t)
	mustCreateUser(t, DefaultAdminAccount, "OldPass123", model.UserTypeRoot)

	if err := (&UserService{}).ResetDefaultAdminPassword("short"); err == nil {
		t.Fatal("少于 8 位的密码应被拒绝")
	}
}
