package utils

import (
	"time"
	"xiaoshuo-backend/config"

	"github.com/golang-jwt/jwt/v4"
)

// Claims JWT声明
type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT token
func GenerateToken(userID uint) (string, error) {
	// 设置过期时间
	expirationTime := time.Now().Add(time.Duration(config.GlobalConfig.JWT.Expires) * time.Second)
	
	// 创建声明
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "xiaoshuo-backend",
		},
	}

	// 创建token对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	// 生成token字符串
	tokenString, err := token.SignedString([]byte(config.GlobalConfig.JWT.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseToken 解析JWT token
func ParseToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GlobalConfig.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ValidationError{}
	}

	return claims, nil
}