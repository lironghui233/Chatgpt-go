package config

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Host        string
		Port        int
		AccessToken string `mapstructure:"access_token"`
	}
	Redis struct {
		Host string
		Port int
		Pwd  string
	}
	Log struct {
		Level   string
		LogPath string `mapstructure:"log_path"`
	}
	WxOfficials []struct {
		AppId  string `mapstructure:"appid"`
		Secret string `mapstructure:"secret"`
	} `mapstructure:"wx_officials"`
	WeComs []struct {
		CorpId     string `mapstructure:"corp_id"`
		CorpSecret string `mapstructure:"corp_secret"`
		App        string `mapstructure:"app"`
	} `mapstructure:"we_coms"`
}

var cfg *Config

func InitConf(configPath string) {
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
