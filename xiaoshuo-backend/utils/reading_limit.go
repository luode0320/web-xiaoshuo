package utils

import (
	"fmt"
	"xiaoshuo-backend/models"
	"gorm.io/gorm"
)

// CheckReadingRestrictions 检查用户的阅读限制
func CheckReadingRestrictions(userID uint) error {
	// 检查用户是否被冻结
	var user models.User
	if err := models.DB.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("用户不存在")
		}
		return err
	}

	// 检查用户是否已激活
	if !user.IsActivated {
		return fmt.Errorf("用户尚未激活，请先激活账户")
	}

	// 检查用户是否被冻结
	if !user.IsActive {
		return fmt.Errorf("用户已被冻结，无法阅读")
	}

	// 可以在这里添加更多限制检查，如：
	// - 每日最大阅读时长
	// - 阅读频率限制
	// - 特定内容访问限制等

	return nil
}