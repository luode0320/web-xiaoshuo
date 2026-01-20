package models

import (
	"gorm.io/gorm"
)

// Comment 评论模型
type Comment struct {
	gorm.Model
	Content    string    `gorm:"not null;comment:评论内容" json:"content" validate:"required,min=1,max=1000"` // 评论内容
	UserID     uint      `gorm:"comment:评论用户ID" json:"user_id"`                                           // 评论用户ID
	User       User      `json:"user"`                                                                    // 评论用户信息
	NovelID    uint      `gorm:"comment:小说ID" json:"novel_id"`                                            // 小说ID
	Novel      Novel     `json:"novel"`                                                                   // 关联的小说
	ChapterID  *uint     `gorm:"comment:章节ID，可选，用于章节评论" json:"chapter_id"`                                // 章节ID，可选，用于章节评论
	ParentID   *uint     `gorm:"comment:父评论ID，支持回复评论" json:"parent_id"`                                   // 父评论ID，支持回复评论
	Parent     *Comment  `json:"parent"`                                                                  // 父评论
	Replies    []Comment `gorm:"foreignKey:ParentID" json:"replies"`                                      // 子评论列表
	LikeCount  int       `gorm:"default:0;comment:点赞数" json:"like_count"`                                 // 点赞数
	IsApproved bool      `gorm:"default:true;comment:评论是否已审核通过" json:"is_approved"`                       // 评论是否已审核通过
}

// TableName 指定表名
func (Comment) TableName() string {
	return "comments"
}
