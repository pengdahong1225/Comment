package middlewares

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func Cors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		method := ctx.Request.Method

		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, AccessToken, X-CSRF-Token, Authorization, Token, token")
		ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PATCH, PUT")
		ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access_Control-Allow-Origin, Access_Control-Headers, Content-Type")
		ctx.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
		}
	}
}

// 统计请求处理时耗
func statCost() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		ctx.Set("midName", "statCost")
		ctx.Next()
		cost := time.Since(start)
		log.Printf("----------- cost: %v -----------", cost)
	}
}
