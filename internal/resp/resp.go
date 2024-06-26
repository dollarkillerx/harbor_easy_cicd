package resp

import "github.com/gin-gonic/gin"

func Resp(ctx *gin.Context, ok bool, message string, data interface{}) {
	ctx.JSON(200, gin.H{
		"success": ok,
		"message": message,
		"data":    data,
	})
}
