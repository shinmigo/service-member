package logic

import (
	"golang.org/x/crypto/bcrypt"
)

// 生成密码
func GeneratePassword(passwordUser string) (string, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(passwordUser), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// 获取
func GetMemberPassword() {

}

// 验证密码
func VerifyPassword(passwordUser string, passwordDb string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(passwordDb), []byte(passwordUser)); err != nil {
		return false
	}

	return true
}
