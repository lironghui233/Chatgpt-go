package routers

import (
	"chatgpt-official/official"

	"github.com/gin-gonic/gin"
)

func InitRouters(r *gin.Engine) {
	api := r.Group("/")
	api.GET("/api/health", official.Health)
	initOfficial(api)
}
func initOfficial(group *gin.RouterGroup) {
	group.GET("/official/receive", official.CheckSignature)
	group.POST("/official/receive", official.ReceiveMessage)
}
