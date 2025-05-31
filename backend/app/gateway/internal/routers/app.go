package routers

import (
	"Comment/app/gateway/internal/controller"
	"Comment/app/gateway/internal/middlewares"
	"Comment/module/sse"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
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

	// sse服务
	sseServer := sse.NewServer()
	sseServer.Run()
	r.GET("/sse", sseServer.Handler)

	// 网关路由
	r.POST("/api/v1/comment", middlewares.RateLimitMiddleware(1*time.Second, 10), middlewares.AuthLogin(), controller.CommentHandler{}.HandleComment)

	return r
}
