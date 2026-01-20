package models

import (
	"gorm.io/gorm"
)

// Keyword 关键词模型
type Keyword struct {
	gorm.Model
	Word   string  `gorm:"uniqueIndex;size:255;not null;comment:关键词内容，唯一索引" json:"word" validate:"required,min=1,max=50"` // 关键词内容，唯一索引
	Weight float64 `gorm:"default:1.0;comment:权重，用于搜索排序" json:"weight"`                                                   // 权重，用于搜索排序
	Novels []Novel `gorm:"many2many:novel_keywords;" json:"novels"`                                                       // 关联的小说列表
}

// TableName 指定表名
func (Keyword) TableName() string {
	return "keywords"
}
