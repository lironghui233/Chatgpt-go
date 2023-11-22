package main

import (
	"chatgpt-proxy/pkg/config"
	"chatgpt-proxy/routers"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	cnf := config.GetConf()
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	routers.InitRouters(r)
	r.Run(fmt.Sprintf("%s:%d", cnf.Http.Host, cnf.Http.Port))
}
