package registry

import (
	"Comment/module/settings"
	"Comment/proto/pb"
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
)

type Registry struct {
	// consul客户端
	client *consulapi.Client
}

func NewRegistry() (*Registry, error) {
	// 配置中心地址
	cfg := settings.Instance().RegistryConfig
	dsn := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	consulConf := consulapi.DefaultConfig()
	consulConf.Address = dsn
	// client
	c, err := consulapi.NewClient(consulConf)
	if err != nil {
		return nil, err
	}
	return &Registry{client: c}, nil
}

func (receiver *Registry) RegisterServiceWithGrpc(info *pb.PBNodeInfo) error {
	id := fmt.Sprintf("%d:%d", info.NodeType, info.NodeId)
	srv := &consulapi.AgentServiceRegistration{
		ID:      id,        // 服务唯一ID
		Name:    info.Name, // 服务名称
		Tags:    []string{info.Name},
		Port:    int(info.Port),
		Address: info.Ip,
		Check: &consulapi.AgentServiceCheck{
			CheckID:                        id,
			GRPC:                           fmt.Sprintf("%s:%d", info.Ip, info.Port),
			Interval:                       "5s",  // 每5秒检测一次
			Timeout:                        "5s",  // 5秒超时
			DeregisterCriticalServiceAfter: "10s", // 超时10秒注销节点
		},
	}
	err := receiver.client.Agent().ServiceRegister(srv)
	if err != nil {
		return err
	}
	return nil
}
func (receiver *Registry) RegisterServiceWithHttp(info *pb.PBNodeInfo) error {
	id := fmt.Sprintf("%d:%d", info.NodeType, info.NodeId)
	srv := &consulapi.AgentServiceRegistration{
		ID:      id,        // 服务唯一ID
		Name:    info.Name, // 服务名称
		Tags:    []string{info.Name},
		Port:    int(info.Port),
		Address: info.Ip,
		Check: &consulapi.AgentServiceCheck{
			CheckID:                        id,
			HTTP:                           fmt.Sprintf("http://%s:%d/%s", info.Ip, info.Port, "health"),
			Interval:                       "5s",  // 每5秒检测一次
			Timeout:                        "5s",  // 5秒超时
			DeregisterCriticalServiceAfter: "10s", // 超时10秒注销节点
		},
	}
	err := receiver.client.Agent().ServiceRegister(srv)
	if err != nil {
		return err
	}
	return nil
}
func (receiver *Registry) UnRegisterService(info *pb.PBNodeInfo) error {
	id := fmt.Sprintf("%d:%d", info.NodeType, info.NodeId)
	err := receiver.client.Agent().ServiceDeregister(id)
	if err != nil {
		return err
	}
	return nil
}
