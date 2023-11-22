package config

import (
	"chatgpt-proxy/pkg/cmd"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	Http struct {
		Host        string
		Port        int
		AccessToken string `mapstructure:"access_token"`
	}
	Chat struct {
		APIKeys          []string `mapstructure:"api_keys"`
		BaseURL          string   `mapstructure:"base_url"`
		Model            string   `mapstructure:"model"`
		MaxTokens        int      `mapstructure:"max_tokens"`
		Temperature      float32  `mapstructure:"temperature"`
		TopP             float32  `mapstructure:"top_p"`
		PresencePenalty  float32  `mapstructure:"presence_penalty"`
		FrequencyPenalty float32  `mapstructure:"frequency_penalty"`
	}
	Log struct {
		Level   string
		LogPath string `mapstructure:"log_path"`
	}
}

var cfg *Config

func init() {
	configPath := cmd.Args.Config
	if configPath == "" {
		panic("请指定应用程序配置文件")
	}
	_, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		panic("配置文件不存在")
	}
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigFile(configPath)
	v.ReadInConfig()
	cfg = &Config{}
	err = v.Unmarshal(cfg)
	if err != nil {
		panic(err.Error())
	}
}

func GetConf() *Config {
	return cfg
}
