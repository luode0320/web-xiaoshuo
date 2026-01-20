package models

import (
	"gorm.io/gorm"
)

// ReviewCriteria 审核标准配置模型
type ReviewCriteria struct {
	gorm.Model
	Name        string `gorm:"not null;size:255;comment:审核标准名称" json:"name" validate:"required,max=255"`             // 审核标准名称
	Description string `gorm:"comment:审核标准描述" json:"description" validate:"max=1000"`                                // 审核标准描述
	Type        string `gorm:"comment:审核标准类型：novel(小说审核), content(内容审核)" json:"type" validate:"oneof=novel content"` // 审核标准类型：novel(小说审核), content(内容审核)
	Content     string `gorm:"type:text;comment:审核标准具体内容" json:"content"`                                            // 审核标准具体内容
	IsActive    bool   `gorm:"default:true;comment:是否启用该审核标准" json:"is_active"`                                      // 是否启用该审核标准
	Weight      int    `gorm:"default:1;comment:重要程度权重，数值越大越重要" json:"weight"`                                       // 重要程度权重，数值越大越重要
	CreatedBy   uint   `gorm:"comment:创建者ID" json:"created_by"`                                                      // 创建者ID
	UpdatedBy   uint   `gorm:"comment:更新者ID" json:"updated_by"`                                                      // 更新者ID
}

// TableName 指定表名
func (ReviewCriteria) TableName() string {
	return "review_criteria"
}
