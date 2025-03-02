package config

import (
	"github.com/spf13/viper"
	"log"
	"note_app_server_mq/model/appModel"
	"sync"
)

var AC *appModel.AppConfig

// InitAppConfig 读取config.yml文件
func InitAppConfig() {
	viper.SetConfigName("AppConfig")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config/yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("read config failed, err:%v\n", err)
	}
	AC = &appModel.AppConfig{}
	if err := viper.Unmarshal(AC); err != nil {
		log.Fatalf("unmarshal config failed, err:%v\n", err)
	}

	var once sync.Once
	var wg sync.WaitGroup
	once.Do(func() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			InitMysqlConfig()
		}()
	})
	wg.Wait()
}
