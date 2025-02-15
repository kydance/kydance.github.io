package main

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		Name string
		Port int

		Info struct {
			Author string
			Email  string
		}
	}

	Database struct {
		Host   string
		Port   int
		User   string
		Passwd string
	}

	Redis struct {
		Host string
		Port int
	}
}

func init() {
	// 默认值
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("app.info", map[string]string{
		"author": "kyden",
		"email":  "kytedance@gmail.com",
	})

	// VIPER_DEMO_xxx
	viper.SetEnvPrefix("VIPER_DEMO")
	fmt.Println(viper.GetEnvPrefix())

	// 绑定环境变量
	if err := viper.BindEnv("app.port", "APP_PORT"); err != nil {
		log.Fatalf("Error binding env: %v", err)
	}
	viper.AutomaticEnv() // 自动读取环境变量
}

func main() {
	var cfg Config

	// 添加搜索路径（可以多个）
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	// 设置配置文件名(无扩展名)
	viper.SetConfigName("cfg")
	// 设置配置文件类型
	viper.SetConfigType("yaml")

	// Read
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %v", err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}

	fmt.Printf("App Name: %s\n", cfg.App.Name)
	fmt.Printf("App Port: %d\n", cfg.App.Port)
	fmt.Printf("App Info Author: %s\n", cfg.App.Info.Author)
	fmt.Printf("App Info Email: %s\n", cfg.App.Info.Email)

	fmt.Printf("Database Host: %s\n", cfg.Database.Host)
	fmt.Printf("Database Port: %d\n", cfg.Database.Port)
	fmt.Printf("Database User: %s\n", cfg.Database.User)
	fmt.Printf("Database Passwd: %s\n", cfg.Database.Passwd)

	fmt.Printf("Redis Host: %s\n", cfg.Redis.Host)
	fmt.Printf("Redis Port: %d\n", cfg.Redis.Port)
}
