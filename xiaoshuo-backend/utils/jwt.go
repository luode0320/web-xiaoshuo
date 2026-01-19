package utils

import (
	"net/http"
	"time"
	"xiaoshuo-backend/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// JwtCustomClaims JWT自定义声明
type JwtCustomClaims struct {
	UserID uint `json:"user_id"`
	IsAdmin bool `json:"is_admin"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT token
func GenerateToken(userID uint, isAdmin bool) (string, error) {
	// 设置过期时间
	expirationTime := time.Now().Add(time.Duration(config.GlobalConfig.JWT.Expires) * time.Second)
	
	// 创建声明
	claims := &JwtCustomClaims{
		UserID: userID,
		IsAdmin: isAdmin,
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
func ParseToken(tokenString string) (*JwtCustomClaims, error) {
	claims := &JwtCustomClaims{}

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

func GetClaims(c *gin.Context) *JwtCustomClaims {
		// 从JWT token获取用户信息
	claimsGet, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取用户信息失败"})
		return nil
	}

	claims, ok := claimsGet.(*JwtCustomClaims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取用户信息失败"})
		return nil
	}

	return claims
}