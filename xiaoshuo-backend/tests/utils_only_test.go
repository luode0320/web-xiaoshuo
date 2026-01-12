// xiaoshuo-backend/tests\utils_only_test.go
// 仅测试工具函数，不依赖数据库

package tests

import (
	"testing"

	"xiaoshuo-backend/models"
	"xiaoshuo-backend/utils"

	"github.com/stretchr/testify/assert"
)

// TestGenerateTokenOnly 测试JWT token生成
func TestGenerateTokenOnly(t *testing.T) {
	// 测试生成token
	userId := uint(1)
	isAdmin := false
	token, err := utils.GenerateToken(userId, isAdmin)
	
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	
	// 测试解析token
	claims, err := utils.ParseToken(token)
	
	assert.NoError(t, err)
	assert.Equal(t, userId, claims.UserID)
	assert.Equal(t, isAdmin, claims.IsAdmin)
}

// TestParseTokenOnly 测试JWT token解析
func TestParseTokenOnly(t *testing.T) {
	// 测试有效token
	userId := uint(2)
	isAdmin := true
	token, err := utils.GenerateToken(userId, isAdmin)
	
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	
	claims, err := utils.ParseToken(token)
	
	assert.NoError(t, err)
	assert.Equal(t, userId, claims.UserID)
	assert.Equal(t, isAdmin, claims.IsAdmin)
	
	// 测试无效token
	invalidToken := "invalid.token.here"
	claims, err = utils.ParseToken(invalidToken)
	
	assert.Error(t, err)
	assert.Nil(t, claims)
}

// TestUserPasswordHashingOnly 测试用户密码哈希功能
func TestUserPasswordHashingOnly(t *testing.T) {
	user := &models.User{}
	password := "test_password_123"
	
	// 测试密码哈希生成
	err := user.HashPassword(password)
	
	assert.NoError(t, err)
	assert.NotEmpty(t, user.Password)
	assert.NotEqual(t, password, user.Password)
	
	// 测试密码验证
	err = user.CheckPassword(password)
	assert.NoError(t, err)
	
	// 测试错误密码验证
	err = user.CheckPassword("wrong_password")
	assert.Error(t, err)
}