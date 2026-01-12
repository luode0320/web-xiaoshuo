package models

import (
	"gorm.io/gorm"
)

// Keyword 关键词模型
type Keyword struct {
	gorm.Model
	Word        string  `gorm:"uniqueIndex;not null" json:"word" validate:"required,min=1,max=50"`
	Weight      float64 `gorm:"default:1.0" json:"weight"` // 权重，用于搜索排序
	Novels      []Novel `gorm:"many2many:novel_keywords;" json:"novels"`
}

// TableName 指定表名
func (Keyword) TableName() string {
	return "keywords"
}