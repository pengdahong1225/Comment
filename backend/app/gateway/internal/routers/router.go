package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// healthCheckRouters 健康检查路由
func healthCheckRouters(engine *gin.Engine) {
	engine.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    "0",
			"message": "health",
		})
	})
}
