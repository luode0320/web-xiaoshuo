package models

import (
	"encoding/json"
	"gorm.io/gorm"
)

// AdminLog 管理员操作日志模型
type AdminLog struct {
	gorm.Model
	AdminUserID uint   `json:"admin_user_id" validate:"required"`
	AdminUser   User   `json:"admin_user" validate:"required"`
	Action      string `json:"action" validate:"required,oneof=approve_novel reject_novel batch_approve_novel freeze_user unfreeze_user delete_pending_novels"` // 操作类型
	TargetType  string `json:"target_type" validate:"required,oneof=novel user comment rating"`                                                                 // 目标类型
	TargetID    uint   `json:"target_id" validate:"required"`                                                                                                // 目标ID
	Details     string `json:"details"`                                                                                                                      // 操作详情(json格式)
	Result      string `json:"result" validate:"required,oneof=success failed"`                                                                              // 操作结果
	IPAddress   string `json:"ip_address"`                                                                                                                   // 操作IP地址
}

// TableName 指定表名
func (AdminLog) TableName() string {
	return "admin_logs"
}

// SetDetails 设置操作详情
func (al *AdminLog) SetDetails(details map[string]interface{}) error {
	data, err := json.Marshal(details)
	if err != nil {
		return err
	}
	al.Details = string(data)
	return nil
}

// GetDetails 获取操作详情
func (al *AdminLog) GetDetails() (map[string]interface{}, error) {
	var details map[string]interface{}
	err := json.Unmarshal([]byte(al.Details), &details)
	return details, err
}