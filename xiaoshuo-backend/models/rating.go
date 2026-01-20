package models

import (
	"gorm.io/gorm"
)

// Rating 评分模型
type Rating struct {
	gorm.Model
	Score      float64 `gorm:"not null;comment:评分分数，0-10分制" json:"score" validate:"required,min=0,max=10"` // 评分分数，0-10分制
	Comment    string  `gorm:"comment:评分说明或评论" json:"comment" validate:"max=500"`                          // 评分说明或评论
	UserID     uint    `gorm:"comment:评分用户ID" json:"user_id"`                                              // 评分用户ID
	User       User    `json:"user"`                                                                       // 评分用户信息
	NovelID    uint    `gorm:"comment:被评分小说ID" json:"novel_id"`                                            // 被评分小说ID
	Novel      Novel   `json:"novel"`                                                                      // 关联的小说
	LikeCount  int     `gorm:"default:0;comment:点赞数" json:"like_count"`                                    // 点赞数
	IsApproved bool    `gorm:"default:true;comment:评分是否已审核通过" json:"is_approved"`                          // 评分是否已审核通过
}

// TableName 指定表名
func (Rating) TableName() string {
	return "ratings"
}
