package main

import (
	"chatgpt-qq/config"
	"chatgpt-qq/log"
	"context"
	"os"
	"os/signal"

	"chatgpt-qq/cmd/cqhttp"
)

func main() {
	cnf := config.GetConf()
	log.SetLevel(cnf.Log.Level)
	log.SetOutput(log.GetRotateWriter(cnf.Log.LogPath))
	if cnf.CqHttp.WsServerPort != 0 {
		// 启动websocket server
		go cqhttp.RunWsServer()
	} else {
		go cqhttp.Run()
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer stop()
	<-ctx.Done()
}
