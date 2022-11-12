package conf

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

type ConfigModule struct {
}

func NewConfigModule() *ConfigModule {
	return &ConfigModule{}
}

func (this *ConfigModule) Run() error {
	viper.SetConfigName("application")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file, %s", err)
	}

	this.initDefaultConfig()

	if err := viper.Unmarshal(&Config); err != nil {
		return fmt.Errorf("unable to decode into struct, %v", err)
	}

	// 监控配置文件变化
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		if err := viper.Unmarshal(&Config); err != nil {
			log.Fatalln(fmt.Sprintf("unmarshal conf failed: %s", err.Error()))
		}
	})

	return nil
}

func (this *ConfigModule) initDefaultConfig() {
	viper.SetDefault("database.maxIdleConn", 10)
	viper.SetDefault("database.maxOpenConn", 100)
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.poolSize", 10)
	viper.SetDefault("redis.minIdleConn", 5)
	viper.SetDefault("jwt.exp", 3600*24*7)
	viper.SetDefault("jwt.key", "a403jk")
}
