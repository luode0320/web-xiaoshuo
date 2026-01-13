package models

import (
	"gorm.io/gorm"
)

// Chapter 章节模型
type Chapter struct {
	gorm.Model
	Title       string `gorm:"not null;size:255" json:"title" validate:"required,min=1,max=200"`
	Content     string `gorm:"type:text" json:"content"`
	WordCount   int    `json:"word_count"`
	Position    int    `json:"position"`        // 章节在小说中的位置
	NovelID     uint   `json:"novel_id"`        // 所属小说ID
	Novel       Novel  `json:"novel"`           // 关联的小说
	FilePath    string `json:"file_path"`       // 章节内容文件路径（对于大章节）
	FileSize    int64  `json:"file_size"`       // 章节文件大小
}

// TableName 指定表名
func (Chapter) TableName() string {
	return "chapters"
}