package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"xiaoshuo-backend/config"
)

// CacheManager 缓存管理器
type CacheManager struct {
	client *redis.Client
}

// GlobalCache 全局缓存实例
var GlobalCache *CacheManager

// ctx Redis操作上下文
var ctx = context.Background()

// InitCache 初始化缓存
func InitCache() error {
	GlobalCache = &CacheManager{
		client: config.RDB, // 使用已有的Redis连接
	}
	return nil
}

// Set 将值存储到缓存
func (c *CacheManager) Set(key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("序列化缓存值失败: %v", err)
	}

	return c.client.Set(ctx, key, data, expiration).Err()
}

// Get 从缓存获取值
func (c *CacheManager) Get(key string, dest interface{}) error {
	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(data), dest)
}

// Exists 检查缓存键是否存在
func (c *CacheManager) Exists(key string) (bool, error) {
	count, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Delete 删除缓存键
func (c *CacheManager) Delete(key string) error {
	return c.client.Del(ctx, key).Err()
}

// SetNX 只设置缓存（仅当键不存在时）
func (c *CacheManager) SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, fmt.Errorf("序列化缓存值失败: %v", err)
	}

	return c.client.SetNX(ctx, key, data, expiration).Result()
}

// GetOrSet 先尝试从缓存获取，如果不存在则使用提供的函数生成值并存储到缓存
func (c *CacheManager) GetOrSet(key string, dest interface{}, expiration time.Duration, generator func() (interface{}, error)) error {
	// 尝试从缓存获取
	err := c.Get(key, dest)
	if err == nil {
		return nil // 缓存命中
	}

	// 缓存未命中，生成值
	value, err := generator()
	if err != nil {
		return fmt.Errorf("生成缓存值失败: %v", err)
	}

	// 存储到缓存
	if setErr := c.Set(key, value, expiration); setErr != nil {
		return fmt.Errorf("存储到缓存失败: %v", setErr)
	}

	// 返回生成的值
	return json.Unmarshal(json.RawMessage(fmt.Sprintf("%v", value)), dest)
}

// GetWithDefault 从缓存获取值，如果不存在则返回默认值
func (c *CacheManager) GetWithDefault(key string, defaultValue interface{}) interface{} {
	var value interface{}
	err := c.Get(key, &value)
	if err != nil {
		return defaultValue
	}
	return value
}