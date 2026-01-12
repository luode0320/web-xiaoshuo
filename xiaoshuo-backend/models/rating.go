package models

import (
	"gorm.io/gorm"
)

// Rating 评分模型
type Rating struct {
	gorm.Model
	Score       float64     `gorm:"not null" json:"score" validate:"required,min=0,max=10"` // 0-10分制
	Comment     string      `json:"comment" validate:"max=500"` // 评分说明
	UserID      uint        `json:"user_id"`
	User        User        `json:"user"`
	NovelID     uint        `json:"novel_id"`
	Novel       Novel       `json:"novel"`
	LikeCount   int         `gorm:"default:0" json:"like_count"`
	IsApproved  bool        `gorm:"default:true" json:"is_approved"` // 评分审核
}

// TableName 指定表名
func (Rating) TableName() string {
	return "ratings"
}