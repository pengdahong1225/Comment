package serverBase

import (
	"Comment/module/logger"
	"Comment/module/registry"
	"Comment/module/services"
	"Comment/module/settings"
	"Comment/proto/pb"
	"errors"
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
)

type Server struct {
	NodeType int
	NodeId   int
	Host     string
	Port     int
	Name     string
	SrvType  string // grpc | http
}

func (receiver *Server) Initialize() error {
	err := logger.Init("./log", "debug")
	if err != nil {
		return err
	}

	flag.IntVar(&receiver.NodeType, "node_type", 0, "node type")
	flag.IntVar(&receiver.NodeId, "node_id", 0, "node id")
	flag.StringVar(&receiver.Host, "host", "", "host")
	flag.IntVar(&receiver.Port, "port", 0, "port")
	flag.StringVar(&receiver.Name, "name", "", "name")
	flag.Parse()

	logrus.Debugf("--------------- name:%v, node_type:%v, node_id:%v, host:%v, port:%v ---------------", receiver.Name, receiver.NodeType, receiver.NodeId, receiver.Host, receiver.Port)

	// 读取配置
	err = settings.Instance().LoadConfig()
	if err != nil {
		return err
	}

	// 初始化services manager
	registryConfig := settings.Instance().RegistryConfig
	dsn := fmt.Sprintf("%s:%d", registryConfig.Host, registryConfig.Port)
	services.Init(dsn)

	return nil
}

func (receiver *Server) Register() error {
	register, err := registry.NewRegistry()
	if err != nil {
		return err
	}

	if receiver.SrvType == "grpc" {
		err = register.RegisterServiceWithGrpc(&pb.PBNodeInfo{
			NodeType: int32(receiver.NodeType),
			NodeId:   int32(receiver.NodeId),
			Ip:       receiver.Host,
			Port:     int32(receiver.Port),
			State:    int32(pb.ENNodeState_EN_NODE_STATE_ONLINE),
			Name:     receiver.Name,
		})
		if err != nil {
			return err
		}
		return nil
	} else if receiver.SrvType == "http" {
		err = register.RegisterServiceWithHttp(&pb.PBNodeInfo{
			NodeType: int32(receiver.NodeType),
			NodeId:   int32(receiver.NodeId),
			Ip:       receiver.Host,
			Port:     int32(receiver.Port),
			State:    int32(pb.ENNodeState_EN_NODE_STATE_ONLINE),
			Name:     receiver.Name,
		})
		if err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("unknown srv type")
	}
}

func (receiver *Server) UnRegister() error {
	register, err := registry.NewRegistry()
	if err != nil {
		return err
	}
	return register.UnRegisterService(&pb.PBNodeInfo{
		NodeType: int32(receiver.NodeType),
		NodeId:   int32(receiver.NodeId),
		Ip:       receiver.Host,
		Port:     int32(receiver.Port),
		State:    int32(pb.ENNodeState_EN_NODE_STATE_OFFLINE),
		Name:     receiver.Name,
	})
}
