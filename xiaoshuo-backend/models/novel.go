package models

import (
	"gorm.io/gorm"
)

// Novel 小说模型
type Novel struct {
	gorm.Model
	Title         string    `gorm:"not null" json:"title" validate:"required,min=1,max=200"`
	Author        string    `gorm:"not null" json:"author" validate:"required,min=1,max=100"`
	Protagonist   string    `json:"protagonist" validate:"max=100"`
	Description   string    `json:"description" validate:"max=1000"`
	Filepath      string    `gorm:"not null" json:"file_path"`
	FileSize      int64     `json:"file_size"`
	WordCount     int       `json:"word_count"`
	ClickCount    int       `gorm:"default:0" json:"click_count"`
	TodayClicks   int       `gorm:"default:0" json:"today_clicks"`
	WeekClicks    int       `gorm:"default:0" json:"week_clicks"`
	MonthClicks   int       `gorm:"default:0" json:"month_clicks"`
	UploadTime    *gorm.DeletedAt `json:"upload_time"`
	LastReadTime  *gorm.DeletedAt `json:"last_read_time"`
	Status        string    `gorm:"default:'pending'" json:"status"` // pending, approved, rejected
	FileHash      string    `gorm:"uniqueIndex;size:255" json:"file_hash"`
	UploadUserID  uint      `json:"upload_user_id"`
	UploadUser    User      `json:"upload_user"`
	Categories    []Category `gorm:"many2many:novel_categories;" json:"categories"`
	Keywords      []Keyword `gorm:"many2many:novel_keywords;" json:"keywords"`
}

// TableName 指定表名
func (Novel) TableName() string {
	return "novels"
}