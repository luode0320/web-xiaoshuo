package models

import (
	"gorm.io/gorm"
)

// ReviewCriteria 审核标准配置模型
type ReviewCriteria struct {
	gorm.Model
	Name        string `gorm:"not null;size:255" json:"name" validate:"required,max=255"` // 标准名称
	Description string `json:"description" validate:"max=1000"`                          // 标准描述
	Type        string `json:"type" validate:"oneof=novel content"`                     // 标准类型 (小说审核/内容审核)
	Content     string `gorm:"type:text" json:"content"`                                // 审核标准内容
	IsActive    bool   `gorm:"default:true" json:"is_active"`                           // 是否启用
	Weight      int    `gorm:"default:1" json:"weight"`                                 // 重要程度权重
	CreatedBy   uint   `json:"created_by"`                                              // 创建者ID
	UpdatedBy   uint   `json:"updated_by"`                                              // 更新者ID
}

// TableName 指定表名
func (ReviewCriteria) TableName() string {
	return "review_criteria"
}