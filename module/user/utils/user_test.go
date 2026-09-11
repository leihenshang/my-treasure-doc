package utils

import (
	"encoding/hex"
	"strings"
	"testing"
)

// TestGenerateLoginTokenLength 验证令牌为 128 位（16 字节 -> 32 个十六进制字符）。
func TestGenerateLoginTokenLength(t *testing.T) {
	token := GenerateLoginToken("9999999999")
	if len(token) != 32 {
		t.Fatalf("期望令牌长度为 32，实际为 %d（%q）", len(token), token)
	}
	if _, err := hex.DecodeString(token); err != nil {
		t.Fatalf("令牌应为合法十六进制字符串: %v", err)
	}
}

// TestGenerateLoginTokenNoPredictablePrefix 验证令牌不再拼接可预测前缀
// （旧实现为 md5("随机4字节"+"treasure-doc-"+userId)），避免离线猜测。
func TestGenerateLoginTokenNoPredictablePrefix(t *testing.T) {
	for _, uid := range []string{"", "9999999999", "treasuredocmgr"} {
		token := GenerateLoginToken(uid)
		if strings.Contains(token, "treasure-doc-") {
			t.Fatalf("令牌不应包含可预测前缀: %q", token)
		}
	}
}

// TestGenerateLoginTokenUniqueness 验证多次生成的令牌足够随机、无碰撞。
func TestGenerateLoginTokenUniqueness(t *testing.T) {
	const n = 10000
	seen := make(map[string]struct{}, n)
	for i := 0; i < n; i++ {
		token := GenerateLoginToken("user")
		if len(token) != 32 {
			t.Fatalf("第 %d 次生成长度异常: %q", i, token)
		}
		if _, ok := seen[token]; ok {
			t.Fatalf("令牌发生碰撞: %q", token)
		}
		seen[token] = struct{}{}
	}
}

// TestPasswordEncryptCompare 验证密码加密与比对的基本契约。
func TestPasswordEncryptCompare(t *testing.T) {
	const pwd = "S3cretP@ss-2026"
	hash, err := PasswordEncrypt(pwd)
	if err != nil {
		t.Fatalf("加密失败: %v", err)
	}
	if hash == pwd {
		t.Fatal("密文不应等于明文")
	}
	if !PasswordCompare(hash, pwd) {
		t.Fatal("正确密码应比对通过")
	}
	if PasswordCompare(hash, "wrong-password") {
		t.Fatal("错误密码应比对失败")
	}
}
