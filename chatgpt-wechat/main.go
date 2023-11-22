package main

import (
	msg_handler "chatgpt-wechat/msg-handler"
	"chatgpt-wechat/pkg/cmd"
	"chatgpt-wechat/pkg/config"
	"chatgpt-wechat/pkg/log"

	"github.com/eatmoreapple/openwechat"
)

func main() {
	loadDependOn()
	bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式

	dispatcher := openwechat.NewMessageMatchDispatcher()
	dispatcher.OnText(msg_handler.NewMsgHandler().TextHandler)

	// 注册消息处理函数
	bot.MessageHandler = dispatcher.AsMessageHandler()
	// 注册登陆二维码回调
	bot.UUIDCallback = openwechat.PrintlnQrcodeUrl

	// 创建热存储容器对象
	reloadStorage := openwechat.NewFileHotReloadStorage("storage.json")
	defer reloadStorage.Close()
	// 免扫码登录
	err := bot.PushLogin(reloadStorage, openwechat.NewRetryLoginOption())
	if err != nil {
		panic(err)
	}

	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	bot.Block()
}
func loadDependOn() {
	config.InitConf(cmd.Args.Config)
	cnf := config.GetConf()
	log.SetOutput(log.GetRotateWriter(cnf.Log.LogPath))
	log.SetLevel(cnf.Log.Level)
}
