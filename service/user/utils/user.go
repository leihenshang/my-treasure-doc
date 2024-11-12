package utils

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

// PasswordEncrypt 对密码使用bcrypt进行加密，返回一个加密字符串
func PasswordEncrypt(password string) (encrypted string, err error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), err
}

// PasswordCompare 对比加密后的hash 和 用户输入的密码是否匹配
// true 为正确，false 为密码不正确
func PasswordCompare(encryptedPassword string, inputPassword string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(inputPassword)); err == nil {
		return true
	}

	return false
}

func GenerateLoginToken(userId int64) string {
	str := "apiDocGo" + strconv.FormatInt(userId, 2)

	// 生成随机数（生成4字节的随机数）
	randomBytes := make([]byte, 4)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return ""
	}
	randomInt := int(binary.BigEndian.Uint32(randomBytes))

	data := []byte(strconv.Itoa(randomInt) + str)
	has := md5.Sum(data)

	return fmt.Sprintf("%x", has)
}
