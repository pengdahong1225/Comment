package internal

import (
	"Comment/app/comment/internal/biz"
	"Comment/app/common/serverBase"
	"Comment/module/signalHandler"
	"Comment/proto/pb"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
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
	// grpc服务
	netAddr := fmt.Sprintf("%s:%d", receiver.Host, receiver.Port)
	listener, err := net.Listen("tcp", netAddr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	grpcServer := grpc.NewServer()

	// 健康检查
	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())
	// 服务
	commentSrv := &biz.CommentHandler{}
	pb.RegisterCommentServiceServer(grpcServer, commentSrv)
	go func() {
		err = grpcServer.Serve(listener)
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
