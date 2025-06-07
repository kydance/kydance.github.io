package main

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var DynamicConfig *viper.Viper

func init() {
	DynamicConfig = viper.New()

	DynamicConfig.AddConfigPath("./config")
	DynamicConfig.AddConfigPath(".")
	DynamicConfig.SetConfigName("dynamic")
	DynamicConfig.SetConfigType("yaml")

	if err := DynamicConfig.ReadInConfig(); err != nil {
		panic(err)
	}

	go func(dc *viper.Viper) {
		dc.WatchConfig()
		dc.OnConfigChange(func(e fsnotify.Event) {
			println("Config file changed:", e.Name)

			// Reload the configuration
			if err := dc.ReadInConfig(); err != nil {
				println("Error reloading config:", err.Error())
			}

			println("Reload config success")
		})
	}(DynamicConfig)
}
