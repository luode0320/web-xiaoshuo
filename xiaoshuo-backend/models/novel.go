package models

import (
	"gorm.io/gorm"
)

// Novel 小说模型
type Novel struct {
	gorm.Model
	Title         string          `gorm:"not null;comment:小说标题" json:"title" validate:"required,min=1,max=200"`                    // 小说标题
	Author        string          `gorm:"not null;comment:小说作者" json:"author" validate:"required,min=1,max=100"`                   // 小说作者
	Protagonist   string          `gorm:"comment:小说主角" json:"protagonist" validate:"max=100"`                                      // 小说主角
	Description   string          `gorm:"comment:小说描述" json:"description" validate:"max=1000"`                                     // 小说描述
	Filepath      string          `gorm:"not null;comment:小说文件存储路径" json:"file_path"`                                              // 小说文件存储路径
	FileSize      int64           `gorm:"comment:小说文件大小（字节）" json:"file_size"`                                                     // 小说文件大小（字节）
	WordCount     int             `gorm:"comment:小说总字数" json:"word_count"`                                                         // 小说总字数
	ClickCount    int             `gorm:"default:0;comment:总点击量" json:"click_count"`                                               // 总点击量
	TodayClicks   int             `gorm:"default:0;comment:今日点击量" json:"today_clicks"`                                             // 今日点击量
	WeekClicks    int             `gorm:"default:0;comment:本周点击量" json:"week_clicks"`                                              // 本周点击量
	MonthClicks   int             `gorm:"default:0;comment:本月点击量" json:"month_clicks"`                                             // 本月点击量
	UploadTime    *gorm.DeletedAt `json:"upload_time"`                                                                             // 上传时间
	LastReadTime  *gorm.DeletedAt `json:"last_read_time"`                                                                          // 最后阅读时间
	Status        string          `gorm:"default:'pending';comment:小说状态：pending(待审核), approved(已通过), rejected(已拒绝)" json:"status"` // 小说状态：pending(待审核), approved(已通过), rejected(已拒绝)
	ChapterStatus string          `gorm:"default:'pending';comment:章节解析状态：pending(待解析), processing(解析中), completed(已完成), failed(解析失败)" json:"chapter_status"` // 章节解析状态：pending(待解析), processing(解析中), completed(已完成), failed(解析失败)
	FileHash      string          `gorm:"uniqueIndex;size:255;comment:小说文件哈希值，用于去重" json:"file_hash"`                              // 小说文件哈希值，用于去重
	UploadUserID  uint            `gorm:"comment:上传用户ID" json:"upload_user_id"`                                                    // 上传用户ID
	UploadUser    User            `json:"upload_user"`                                                                             // 上传用户信息
	Categories    []Category      `gorm:"many2many:novel_categories;" json:"categories"`                                           // 小说分类
	Keywords      []Keyword       `gorm:"many2many:novel_keywords;" json:"keywords"`                                               // 小说关键词
	AverageRating float64         `gorm:"default:0;comment:平均评分" json:"average_rating"`                                            // 平均评分
	RatingCount   int             `gorm:"default:0;comment:评分数量" json:"rating_count"`                                              // 评分数量
	Chapters      []Chapter       `json:"chapters"`                                                                                // 小说章节
}

// TableName 指定表名
func (Novel) TableName() string {
	return "novels"
}
