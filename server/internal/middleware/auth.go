package middleware

import (
	"github.com/gin-gonic/gin"
)

func Auth(authorization string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if authorization != token {
			ctx.JSON(400, gin.H{
				"error": "Authorization token error",
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
