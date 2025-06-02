package internal

import (
	"Comment/app/common/serverBase"
	sse "Comment/app/sse/internal/sse_server"
	"Comment/module/middlewares"
	"Comment/module/signalHandler"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
)

type Server struct {
	serverBase.Server
	SignalHandler *signalHandler.SignalHandler
}

func (receiver *Server) Init() error {
	receiver.SignalHandler = signalHandler.NewSignalHandler(receiver.Reload, receiver.Online, receiver.Offline)
	go receiver.SignalHandler.Listen()

	err := receiver.Initialize()
	if err != nil {
		return err
	}

	return nil
}

func (receiver *Server) Start() {
	// sse服务
	path := fmt.Sprintf("%s/web.log", "./log")
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Errorf("web日志文件打开失败：%s", err.Error())
	}
	gin.DefaultWriter = io.MultiWriter(os.Stdout, file)
	gin.SetMode(os.Getenv("GIN_MODE"))

	engine := gin.Default()
	engine.Use(middlewares.Cors())
	engine.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    "0",
			"message": "health",
		})
	})
	engine.GET("/api/v1/sse", middlewares.AuthLogin(), sse.Instance().Handler)

	dsn := fmt.Sprintf(":%d", receiver.Port)
	go func() {
		err := engine.Run(dsn)
		if err != nil {
			panic(err)
		}
	}()

	// 服务注册
	err = receiver.Register()
	if err != nil {
		panic(err)
	}

	select {}
}

func (receiver *Server) Reload() {
	logrus.Infof("服务重载")
}
func (receiver *Server) Online() {
	// 注册
	err := receiver.Register()
	if err != nil {
		logrus.Errorf("服务上线失败,err=%s", err.Error())
	}
}
func (receiver *Server) Offline() {
	err := receiver.UnRegister()
	if err != nil {
		logrus.Errorf("服务下线失败,err=%s", err.Error())
	}
}
