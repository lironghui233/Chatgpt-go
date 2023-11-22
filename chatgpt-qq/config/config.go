package config

import (
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Enterprise struct {
		Id string `mapstructure:"id"`
	} `mapstructure:"enterprise"`
	CqHttp struct {
		WebSocket    string `mapstructure:"websocket"`
		WsServerHost string `mapstructure:"ws_server_host"`
		WsServerPort int    `mapstructure:"ws_server_port"`
		AccessToken  string `mapstructure:"access_token"`
		AtOnly       bool   `mapstructure:"at_only"`
		UseKeyword   bool   `mapstructure:"use_keyword"`
		KeywordType  string `mapstructure:"keyword_type"`
		Keyword      string `mapstructure:"keyword"`
		TimeOut      int    `mapstructure:"timeout"`
	}
	Context struct {
		PrivateContext bool `mapstructure:"private_context"`
		GroupContext   bool `mapstructure:"group_context"`
	}
	ChatGPTService struct {
		Address     string `mapstructure:"address"`
		AccessToken string `mapstructure:"access_token"`
	} `mapstructure:"chatgpt-service"`
	Log struct {
		Level   string
		LogPath string `mapstructure:"log_path"`
	}
	Chat struct {
		Model             string  `mapstructure:"model"`
		MaxTokens         int     `mapstructure:"max_tokens"`
		Temperature       float32 `mapstructure:"temperature"`
		TopP              float32 `mapstructure:"top_p"`
		PresencePenalty   float32 `mapstructure:"presence_penalty"`
		FrequencyPenalty  float32 `mapstructure:"frequency_penalty"`
		BotDesc           string  `mapstructure:"bot_desc"`
		ContextTTL        int     `mapstructure:"context_ttl"`
		ContextLen        int     `mapstructure:"context_len"`
		MinResponseTokens int     `mapstructure:"min_response_tokens"`
	}
}

var Cfg Config

func init() {
	if _, err := os.Stat("config.cfg"); os.IsNotExist(err) {
		f, err := os.Create("config.cfg")
		if err != nil {
			log.Println(err)
		}
		// 自动生成配置文件
		_, err = f.Write([]byte("# config.toml 配置文件\n\n" +
			"# 企业ID\n[enterprise]\n" +
			"id = \"0voice\"\n" +
			"# cqhttp机器人配置\n[cqhttp]\n" +
			"# go-cqhttp的正向WebSocket地址\n" +
			"websocket = \"ws://127.0.0.1:8080\"\n" +
			"# WebSocket服务监听的主机地址\n" +
			"ws_server_host = \"0.0.0.0\"\n" +
			"# WebSocket服务监听的端口\n" +
			"ws_server_port = 8080\n" +
			"access_token = \"ksSdkFdDngKie8nth0yi29405nr9jey84prEhw5u43780yjr3h7s\"\n" +
			"# 群聊是否需要@机器人才能触发\n" +
			"at_only = true\n" +
			"# 是否开启触发关键词\n" +
			"use_keyword = false\n" +
			"# 触发关键词场合 可选值: all, group, private, 开启群聊关键词建议关闭at_only\n" +
			"keyword_type = \"group\"\n" +
			"# 触发关键词\n" +
			"keyword = \"对话\"\n" +
			"# 生成中提醒时间秒数\n" +
			"timeout = 30\n\n" +
			"# 连续对话相关（实际使用中，连续对话似乎会导致更多的token使用，在这里可以设置是否启用这个功能。默认关闭。另注：预设角色不支持连续对话。）\n[context]\n" +
			"# 是否在私聊中启用连续对话\n" +
			"private_context = true\n" +
			"# 是否在群聊中启用连续对话\n" +
			"group_context = true\n" +
			"# chatgpt-qq依赖的chatgpt-service 服务 \n" +
			"[chatgpt-service]\n" +
			"address = \"192.168.239.154:50051\"\n" +
			"access_token = \"05nr9jey84prEhw5u43780yjr3h7sksSdkFdDngKie8nth0yi294\"\n" +
			"# 日志配置\n" +
			"[log]\n" +
			"# panic,fatal,error,warn,warning,info,debug,trace\n" +
			"level = \"info\"\n" +
			"log_path = \"runtime/app.log\"\n" +
			"[chat]\n" +
			"# 使用的训练模型\n" +
			"model = \"gpt-3.5-turbo-0301\"\n" +
			"# 单次请求的上下文总长度，包括 请求消息+completion.maxToken 两者总计不能超过4097\n" +
			"max_tokens = 1024\n" +
			"# 表示语言模型输出的随机性和创造性\n" +
			"# 取值范围 0 ~ 1，值越大随机性和创造性越高\n" +
			"temperature = 0.8\n" +
			"# 用于生成文本时控制选词的随机程度\n" +
			"# 即下一个预测单词考虑的概率范围\n" +
			"# 取值范围 0 ~ 1，例如：0.5 表示考虑选择的单词累计概率大于等于0.5\n" +
			"top_p = 0.9\n" +
			"# 存在惩罚，用于生成文本时控制重复使用单词的程度\n" +
			"# 取值 0 ~ 1,0表示不惩罚，1表示完全禁止重复单词\n" +
			"# 完全禁止重复单词会影响生成文本的流畅性和连贯性\n" +
			"presence_penalty = 0.8\n" +
			"# 用于控制模型生成回复时重复单词出现的频率\n" +
			"# 取值 0 ~ 1，值越大生成的回复会更注重避免使用已经出现的单词\n" +
			"frequency_penalty = 0.5\n" +
			"# AI助手特征描述\n" +
			"bot_desc = \"你是一个AI助手，我需要你模拟一名资深软件工程师来回答我的问题\"\n" +
			"# 上下文缓存时长，单位s\n" +
			"context_ttl = 1800\n" +
			"# 上下文消息条数\n" +
			"context_len = 4\n" +
			"# 单次请求，保留的响应tokens数量\n" +
			"min_response_tokens = 512"))
		if err != nil {
			panic(err)
		}
		log.Println("配置文件不存在, 已自动生成配置文件, 请修改配置文件后再次运行程序, 5秒后退出程序...")
		time.Sleep(5 * time.Second)
		os.Exit(0)
	}
	viper.SetConfigName("config.cfg")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".") // 指定查找配置文件的路径
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&Cfg)
	if err != nil {
		panic(err)
	}
}

func GetConf() Config {
	return Cfg
}
