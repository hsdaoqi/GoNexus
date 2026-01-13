package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword 密码加密
func HashPassword(pwd string) (string, error) {
	// GenerateFromPassword 会自动生成盐值并哈希
	// 第二个参数 cost 是计算强度，取 12 比较均衡
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPassword 校验密码
func CheckPassword(pwd, hash string) bool {
	// CompareHashAndPassword 比较密文和明文是否匹配
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
	return err == nil
}
