package routers

import (
	"Comment/app/gateway/internal/biz"
	"Comment/module/middlewares"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"time"
)

func Router() *gin.Engine {
	path := fmt.Sprintf("%s/web.log", "./log")
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Errorf("web日志文件打开失败：%s", err.Error())
	}
	gin.DefaultWriter = io.MultiWriter(os.Stdout, file)
	gin.SetMode(os.Getenv("GIN_MODE"))

	r := gin.Default()
	r.Use(middlewares.Cors()) // 跨域处理

	// 初始化路由
	healthCheckRouters(r)

	// 路由
	group := r.Group("/api/v1")
	group.POST("/api/v1/comment", middlewares.RateLimitMiddleware(1*time.Second, 10), middlewares.AuthLogin(), biz.CommentHandler{}.HandleAddComment)

	return r
}

// healthCheckRouters 健康检查路由
func healthCheckRouters(engine *gin.Engine) {
	engine.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    "0",
			"message": "health",
		})
	})
}
