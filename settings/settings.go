package settings

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

func Init() (err error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")
	if err = viper.ReadInConfig(); err != nil {
		log.Printf("viper.ReadInConfig failed, err :%s", err)
		return err
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		log.Println("配置文件修改了!")
	})
	return nil
}
