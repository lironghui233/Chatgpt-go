package main

import (
	"chatgpt-official/pkg/cmd"
	"chatgpt-official/pkg/config"
	"chatgpt-official/pkg/log"
	"chatgpt-official/routers"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	loadDependOn()

	gin.SetMode(gin.ReleaseMode)
	gin.New()
	r := gin.Default()
	routers.InitRouters(r)

	cnf := config.GetConf()
	addr := fmt.Sprintf("%s:%d", cnf.Http.Host, cnf.Http.Port)
	err := r.Run(addr)
	if err != nil {
		log.Error(err)
	}
}

func loadDependOn() {
	config.InitConf(cmd.Args.Config)
	cnf := config.GetConf()
	log.SetOutput(log.GetRotateWriter(cnf.Log.LogPath))
	log.SetLevel(cnf.Log.Level)
}
