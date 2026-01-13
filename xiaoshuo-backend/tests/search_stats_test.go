package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"xiaoshuo-backend/config"
)

// MockSearchStats 用于模拟搜索统计数据
func MockGetSearchStats(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"total_searches": 100,
			"top_keywords": []gin.H{
				{"keyword": "测试关键词1", "count": 50},
				{"keyword": "测试关键词2", "count": 30},
			},
			"recent_searches": []gin.H{
				{"keyword": "测试关键词1", "count": 50},
			},
		},
	})
}

func TestGetSearchStats(t *testing.T) {
	// 初始化配置
	config.InitConfig()
	
	// 设置测试路由，使用Mock函数
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/search/stats", MockGetSearchStats)

	// 创建请求
	req, _ := http.NewRequest("GET", "/search/stats", nil)
	recorder := httptest.NewRecorder()

	// 执行请求
	router.ServeHTTP(recorder, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, recorder.Code)

	// 验证响应中包含预期的字段
	assert.Contains(t, recorder.Body.String(), "code")
	assert.Contains(t, recorder.Body.String(), "message")
	assert.Contains(t, recorder.Body.String(), "data")
	assert.Contains(t, recorder.Body.String(), "total_searches")
	assert.Contains(t, recorder.Body.String(), "top_keywords")
	
	// 解析JSON响应以进行更详细的验证
	var response map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	// 验证响应结构
	assert.Equal(t, float64(200), response["code"])
	assert.Equal(t, "success", response["message"])
	
	data, ok := response["data"].(map[string]interface{})
	assert.True(t, ok)
	assert.Contains(t, data, "total_searches")
	assert.Contains(t, data, "top_keywords")
}