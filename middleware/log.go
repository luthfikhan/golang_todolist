package middleware

import (
	"strconv"
	"time"
	"todolist/helper"

	"github.com/gin-gonic/gin"
)

func Log(ctx *gin.Context) {
	now := time.Now()

	ctx.Next()

	helper.Log.Info(gin.H{
		"method":    ctx.Request.Method,
		"path":      ctx.Request.URL.Path,
		"code":      ctx.Writer.Status(),
		"timetaken": strconv.Itoa(int(time.Now().UnixMilli()-now.UnixMilli())) + " ms",
	})
}
