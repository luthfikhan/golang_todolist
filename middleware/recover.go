package middleware

import (
	"net/http"
	"todolist/helper"

	"github.com/gin-gonic/gin"
)

func Recover() gin.HandlerFunc {
	return gin.CustomRecovery(func(ctx *gin.Context, err interface{}) {
		helper.Log.Error(gin.H{
			"error": err,
		}, "Recover from panic")
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})

		gin.RecoveryWithWriter(gin.DefaultErrorWriter)
	})
}
