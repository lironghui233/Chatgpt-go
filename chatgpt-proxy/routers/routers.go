package routers

import (
	"chatgpt-proxy/health"
	"chatgpt-proxy/middleware"
	"chatgpt-proxy/proxy"
	"github.com/gin-gonic/gin"
)

func InitRouters(r *gin.Engine) {
	r.GET("/health", health.Health)
	r.Use(middleware.Auth(), middleware.RateLimit(10, 10))
	initProxyRouter(r)
}
func initProxyRouter(r *gin.Engine) {
	p := proxy.NewProxy()
	v1 := r.Group("/v1")
	v1.Any("/*relativePath", p.ChatProxy)
}
