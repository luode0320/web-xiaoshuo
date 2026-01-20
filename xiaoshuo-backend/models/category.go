package models

import (
	"gorm.io/gorm"
)

// Category 分类模型
type Category struct {
	gorm.Model
	Name        string     `gorm:"uniqueIndex;size:255;not null;comment:分类名称，唯一索引" json:"name" validate:"required,min=1,max=50"` // 分类名称，唯一索引
	Description string     `gorm:"comment:分类描述" json:"description" validate:"max=200"`                                           // 分类描述
	ParentID    *uint      `gorm:"comment:父分类ID，支持分类层级" json:"parent_id"`                                                        // 父分类ID，支持分类层级
	Parent      *Category  `json:"parent"`                                                                                       // 父分类
	Children    []Category `gorm:"foreignKey:ParentID" json:"children"`                                                          // 子分类列表
	Novels      []Novel    `gorm:"many2many:novel_categories;" json:"novels"`                                                    // 该分类下的小说列表
}

// TableName 指定表名
func (Category) TableName() string {
	return "categories"
}
