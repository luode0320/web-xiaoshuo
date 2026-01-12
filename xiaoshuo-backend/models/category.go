package models

import (
	"gorm.io/gorm"
)

// Category 分类模型
type Category struct {
	gorm.Model
	Name        string  `gorm:"uniqueIndex;size:255;not null" json:"name" validate:"required,min=1,max=50"`
	Description string  `json:"description" validate:"max=200"`
	ParentID    *uint   `json:"parent_id"` // 支持分类层级
	Parent      *Category `json:"parent"`
	Children    []Category `gorm:"foreignKey:ParentID" json:"children"`
	Novels      []Novel   `gorm:"many2many:novel_categories;" json:"novels"`
}

// TableName 指定表名
func (Category) TableName() string {
	return "categories"
}