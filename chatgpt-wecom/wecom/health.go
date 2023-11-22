package wecom

import "github.com/gin-gonic/gin"

func Health(ctx *gin.Context) {
	ctx.JSON(200, nil)
	return
}
