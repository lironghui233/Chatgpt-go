package routers

import (
	"chatgpt-wecom/wecom"

	"github.com/gin-gonic/gin"
)

func InitRouters(r *gin.Engine) {
	api := r.Group("/")
	api.GET("/api/health", wecom.Health)
	initOfficial(api)
}
func initOfficial(group *gin.RouterGroup) {
	group.GET("/wecom/receive", wecom.CheckSignature)
	group.POST("/wecom/receive", wecom.ReceiveMessage)
}
