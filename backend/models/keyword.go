package models

import (
	"gorm.io/gorm"
)

// Keyword 关键词模型
type Keyword struct {
	gorm.Model
	Keyword string  `gorm:"uniqueIndex;not null" json:"keyword" validate:"required,min=1,max=50"`
	Novels  []Novel `gorm:"many2many:novel_keywords;" json:"novels,omitempty"`
}

// TableName 指定表名
func (Keyword) TableName() string {
	return "keywords"
}