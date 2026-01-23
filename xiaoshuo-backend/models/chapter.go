package models

import (
	"gorm.io/gorm"
)

// Chapter 章节模型
type Chapter struct {
	gorm.Model
	Title     string `gorm:"not null;size:255;comment:章节标题" json:"title" validate:"required,min=1,max=200"` // 章节标题
	Content   string `gorm:"type:text;comment:章节内容" json:"content"`                                         // 章节内容
	WordCount int    `gorm:"comment:章节字数" json:"word_count"`                                                // 章节字数
	Position  int    `gorm:"index:idx_novel_position;comment:章节在小说中的位置" json:"position"`                     // 章节在小说中的位置
	NovelID   uint   `gorm:"index:idx_novel_position;comment:所属小说ID" json:"novel_id"`                        // 所属小说ID
	Novel     Novel  `json:"novel"`                                                                         // 关联的小说
	FilePath  string `gorm:"comment:章节内容文件路径（对于大章节）" json:"file_path"`                                      // 章节内容文件路径（对于大章节）
	FileSize  int64  `gorm:"comment:章节文件大小（字节）" json:"file_size"`                                           // 章节文件大小（字节）
}

// TableName 指定表名
func (Chapter) TableName() string {
	return "chapters"
}
