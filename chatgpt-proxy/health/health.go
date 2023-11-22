package health

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Health(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "健康检查成功",
	})
}
