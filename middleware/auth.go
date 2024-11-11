package middleware

import (
	"net/http"
	"todolist/helper"

	"github.com/gin-gonic/gin"
)

func Auth(ctx *gin.Context) {
	authorization := ctx.GetHeader("authorization")
	_, err := helper.VerifyToken(authorization)

	if err != nil {
		helper.Log.Info(gin.H{"error": err.Error()}, "Unauthorized")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
	}

	ctx.Next()
}
