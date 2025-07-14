package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig
	MySQL  MySQLConfig
	Redis  RedisConfig
}

type ServerConfig struct {
	Port int
}

type MySQLConfig struct {
	Host         string
	Port         int
	User         string
	Password     string
	DBName       string
	MaxIdleConns int
	MaxOpenConns int
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
	UserName string
}

var Conf *Config

func InitConfig() {
	viper.SetConfigName("config")   // 配置文件名 (不含扩展名)
	viper.SetConfigType("yaml")     // 配置文件类型
	viper.AddConfigPath("./config") // 配置文件路径

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("读取配置文件失败: %w", err))
	}

	Conf = &Config{}

	if err := viper.Unmarshal(Conf); err != nil {
		panic(fmt.Errorf("解析配置文件失败: %w", err))
	}

	fmt.Println("配置文件加载成功")
}
