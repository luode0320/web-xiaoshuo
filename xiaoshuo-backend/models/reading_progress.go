package models

import (
	"gorm.io/gorm"
)

// ReadingProgress 阅读进度模型
type ReadingProgress struct {
	gorm.Model
	UserID      uint   `json:"user_id"`
	User        User   `json:"user"`
	NovelID     uint   `json:"novel_id"`
	Novel       Novel  `json:"novel"`
	ChapterID   uint   `json:"chapter_id"`
	ChapterName string `json:"chapter_name"`
	Position    int    `json:"position"` // 当前阅读位置
	Progress    int    `json:"progress"` // 阅读进度百分比
	LastReadAt  *gorm.DeletedAt `json:"last_read_at"`
}

// TableName 指定表名
func (ReadingProgress) TableName() string {
	return "reading_progress"
}