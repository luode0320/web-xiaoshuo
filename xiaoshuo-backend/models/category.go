package models

import (
	"gorm.io/gorm"
)

// Category 分类模型
type Category struct {
	gorm.Model
	Name        string  `gorm:"uniqueIndex;not null" json:"name" validate:"required,min=1,max=50"`
	Description string  `json:"description" validate:"max=200"`
	Novels      []Novel `gorm:"many2many:novel_categories;" json:"novels,omitempty"`
}

// TableName 指定表名
func (Category) TableName() string {
	return "categories"
}