package settings

import (
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"sync"
)

var (
	conf *AppConfig
	once sync.Once
)

func Instance() *AppConfig {
	once.Do(func() {
		conf = new(AppConfig)
	})
	return conf
}

func (receiver *AppConfig) LoadConfig() error {
	// 读取配置
	v := viper.New()
	v.AddConfigPath("config")
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	if e := v.ReadInConfig(); e != nil {
		return e
	}
	if e := v.Unmarshal(conf); e != nil {
		return e
	}
	// 配置热重载
	viper.OnConfigChange(func(event fsnotify.Event) {
		logrus.Errorf("config file changed: %s", event.Name)
		if e := v.Unmarshal(conf); e != nil {
			logrus.Errorf("config file update failed: %s", event.Name)
		}
	})
	// 监听配置文件
	viper.WatchConfig()
	return nil
}
