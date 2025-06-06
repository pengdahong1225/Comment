package internal

import (
	"Comment/app/common/serverBase"
	"Comment/app/gateway/internal/routers"
	"Comment/module/signalHandler"
	"fmt"
	"github.com/sirupsen/logrus"
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
	// http服务
	engine := routers.Router()
	dsn := fmt.Sprintf(":%d", receiver.Port)
	go func() {
		err := engine.Run(dsn)
		if err != nil {
			panic(err)
		}
	}()

	// 服务注册
	err := receiver.Register()
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
