package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Config 应用配置结构体
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         int    `mapstructure:"port"`
	Mode         string `mapstructure:"mode"` // debug, release, test
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host            string `mapstructure:"host"`
	Port            string `mapstructure:"port"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	DBName          string `mapstructure:"dbname"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"` // 秒
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

var globalConfig *Config

// Load 加载配置
func Load(configPath string) (*Config, error) {
	viper.SetConfigType("yaml")
	viper.SetConfigFile(configPath)

	// 支持环境变量覆盖
	viper.AutomaticEnv()

	// 设置环境变量前缀（可选）
	viper.SetEnvPrefix("APP")

	// 设置环境变量替换函数，支持 ${VAR} 格式
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 替换配置结构体中的环境变量占位符
	expandConfigEnvVars(&config)

	globalConfig = &config
	log.Printf("配置文件加载成功: %s", configPath)
	return &config, nil
}

// expandConfigEnvVars 展开配置结构体中的环境变量
func expandConfigEnvVars(cfg *Config) {
	// 展开数据库配置中的环境变量
	cfg.Database.Port = os.ExpandEnv(cfg.Database.Port)
	cfg.Database.Host = os.ExpandEnv(cfg.Database.Host)
	cfg.Database.User = os.ExpandEnv(cfg.Database.User)
	cfg.Database.Password = os.ExpandEnv(cfg.Database.Password)
	cfg.Database.DBName = os.ExpandEnv(cfg.Database.DBName)

	// 展开服务器配置中的环境变量
	cfg.Server.Mode = os.ExpandEnv(cfg.Server.Mode)

	// 展开 Redis 配置中的环境变量
	cfg.Redis.Host = os.ExpandEnv(cfg.Redis.Host)
	cfg.Redis.Password = os.ExpandEnv(cfg.Redis.Password)
}

// Get 获取全局配置
func Get() *Config {
	return globalConfig
}

// GetDSN 获取数据库连接字符串
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.DBName,
	)
}

// GetAddr 获取 Redis 地址
func (c *RedisConfig) GetAddr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}
