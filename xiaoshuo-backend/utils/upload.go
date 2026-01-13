package utils

import (
	"fmt"
	"time"
)

// GetUploadCount 获取用户今日上传次数
func GetUploadCount(userID uint) (int, error) {
	key := fmt.Sprintf("upload_count:%d:%s", userID, time.Now().Format("2006-01-02"))
	
	countVal := GlobalCache.GetWithDefault(key, 0)
	count, ok := countVal.(int)
	if !ok {
		return 0, fmt.Errorf("无法获取上传次数")
	}
	
	return count, nil
}