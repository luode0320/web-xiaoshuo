package config

import (
	"log"
	"sync"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	redis "github.com/go-redis/redis/v8"
)

var (
	DB   *gorm.DB
	RDB  *redis.Client
	once sync.Once
)

// Config 应用配置
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port     string `mapstructure:"port"`
	Mode     string `mapstructure:"mode"`
	BasePath string `mapstructure:"base_path"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	Charset  string `mapstructure:"charset"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret  string `mapstructure:"secret"`
	Expires int64  `mapstructure:"expires"`
}

// GlobalConfig 全局配置变量
var GlobalConfig *Config

// InitConfig 初始化配置
func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	// 设置默认值
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("server.base_path", "/api/v1")
	viper.SetDefault("database.host", "192.168.3.3")
	viper.SetDefault("database.port", "3306")
	viper.SetDefault("database.user", "root")
	viper.SetDefault("database.password", "Ld588588")
	viper.SetDefault("database.dbname", "xiaoshuo")
	viper.SetDefault("database.charset", "utf8mb4")
	viper.SetDefault("redis.addr", "192.168.3.3:6379")
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("jwt.secret", "xiaoshuo_secret_key")
	viper.SetDefault("jwt.expires", 3600)

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("配置文件读取失败: %v", err)
	}

	GlobalConfig = &Config{}
	if err := viper.Unmarshal(GlobalConfig); err != nil {
		log.Fatalf("配置文件解析失败: %v", err)
	}

	log.Println("配置初始化完成")
}

// InitDB 初始化数据库连接
func InitDB() {
	once.Do(func() {
		dsn := GlobalConfig.Database.User + ":" + GlobalConfig.Database.Password + 
			"@tcp(" + GlobalConfig.Database.Host + ":" + GlobalConfig.Database.Port + 
			")/" + GlobalConfig.Database.DBName + "?charset=" + GlobalConfig.Database.Charset + "&parseTime=True&loc=Local"

		var err error
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			log.Fatalf("数据库连接失败: %v", err)
		}

		log.Println("数据库连接成功")
	})
}

// InitRedis 初始化Redis连接
func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     GlobalConfig.Redis.Addr,
		Password: GlobalConfig.Redis.Password,
		DB:       GlobalConfig.Redis.DB,
	})

	log.Println("Redis连接成功")
}