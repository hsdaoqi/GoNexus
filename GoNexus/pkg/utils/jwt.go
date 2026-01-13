package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"go-nexus/pkg/global"
	"time"
)

// MyClaims 自定义声明结构体并内嵌 jwt.RegisteredClaims
// jwt包自带的 RegisteredClaims 包含：iss(签发人), exp(过期时间), sub(主题), etc.
type MyClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenToken 生成 JWT
func GenToken(userID uint, username string) (string, error) {
	// 1. 创建声明
	claims := MyClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			//过期时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(global.Config.Jwt.Expire) * time.Hour)),
			Issuer:    "GoNexus", //签发人
		},
	}

	// 2. 使用指定的签名方法 (HS256) 创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 3. 使用 SecretKey 签名并获得完整的编码后的字符串 Token
	return token.SignedString([]byte(global.Config.Jwt.Secret))
}

// ParseToken 解析 JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	// 1. 解析 token
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.Config.Jwt.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	// 2. 校验 token 是否有效，并提取 Claims
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// GetUserIDFromToken 从token中提取用户ID (方便在中间件中使用)
func GetUserIDFromToken(tokenString string) (uint, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}
