package models

import (
	"gorm.io/gorm"
)

// AdminLog 管理日志模型
type AdminLog struct {
	gorm.Model
	AdminUserID uint   `gorm:"comment:管理员用户ID" json:"admin_user_id"`                                                                    // 管理员用户ID
	AdminUser   User   `json:"admin_user"`                                                                                              // 管理员用户信息
	Action      string `gorm:"not null;comment:操作类型，如 approve_novel, delete_comment 等" json:"action" validate:"required,min=1,max=100"` // 操作类型，如 approve_novel, delete_comment 等
	TargetType  string `gorm:"comment:目标类型，如 novel, user, comment, rating" json:"target_type" validate:"max=50"`                        // 目标类型，如 "novel", "user", "comment", "rating"
	TargetID    uint   `gorm:"comment:目标ID，对应目标类型的记录ID" json:"target_id"`                                                               // 目标ID，对应目标类型的记录ID
	Details     string `gorm:"comment:操作详情，描述具体操作内容" json:"details"`                                                                    // 操作详情，描述具体操作内容
	IPAddress   string `gorm:"comment:操作时的IP地址" json:"ip_address"`                                                                      // 操作时的IP地址
	UserAgent   string `gorm:"comment:操作时的用户代理信息" json:"user_agent"`                                                                    // 操作时的用户代理信息
}

// TableName 指定表名
func (AdminLog) TableName() string {
	return "admin_logs"
}
