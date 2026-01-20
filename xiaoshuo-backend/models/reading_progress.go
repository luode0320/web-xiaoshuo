package models

import (
	"gorm.io/gorm"
)

// ReadingProgress 阅读进度模型
type ReadingProgress struct {
	gorm.Model
	UserID      uint            `gorm:"comment:用户ID" json:"user_id"`             // 用户ID
	User        User            `json:"user"`                                    // 用户信息
	NovelID     uint            `gorm:"comment:小说ID" json:"novel_id"`            // 小说ID
	Novel       Novel           `json:"novel"`                                   // 小说信息
	ChapterID   uint            `gorm:"comment:当前阅读章节ID" json:"chapter_id"`      // 当前阅读章节ID
	ChapterName string          `gorm:"comment:当前阅读章节名称" json:"chapter_name"`    // 当前阅读章节名称
	Position    int             `gorm:"comment:当前阅读位置（字符位置或页数）" json:"position"` // 当前阅读位置（字符位置或页数）
	Progress    int             `gorm:"comment:阅读进度百分比" json:"progress"`         // 阅读进度百分比
	ReadingTime int             `gorm:"comment:阅读时长（秒）" json:"reading_time"`     // 阅读时长（秒）
	LastReadAt  *gorm.DeletedAt `json:"last_read_at"`                            // 最后阅读时间
}

// TableName 指定表名
func (ReadingProgress) TableName() string {
	return "reading_progress"
}
