package config

import (
	"github.com/dapr-platform/common"
	"github.com/spf13/viper"
	"os"
	"strings"
)

// MiniProgramConfig 小程序配置结构体
type MiniProgramConfig struct {
	AppID     string `mapstructure:"app_id"`
	AppSecret string `mapstructure:"app_secret"`
}

// WeChatConfig 微信服务的整体配置
type WeChatConfig struct {
	MiniPrograms map[string]MiniProgramConfig `mapstructure:"mini_programs"`
}

var CLIENT_ID = "default" //需要和authz-service配置一致
var CLIENT_SECRET = "secret"

func init() {
	if val := os.Getenv("CLIENT_ID"); val != "" {
		CLIENT_ID = val
	}
	if val := os.Getenv("CLIENT_SECRET"); val != "" {
		CLIENT_SECRET = val
	}
}

func GetWeChatConfig() (*WeChatConfig, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")

	// 允许使用环境变量，环境变量名自动转换为大写，使用下划线连接
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.ReadInConfig()
	if err != nil {
		common.Logger.Fatalf("Error reading config file: %s \n", err)
		return nil, err
	}

	// 解析到结构体
	var config WeChatConfig
	err = viper.Unmarshal(&config)
	if err != nil {
		common.Logger.Fatalf("Unable to decode into struct: %s \n", err)
		return nil, err
	}

	return &config, nil
}
