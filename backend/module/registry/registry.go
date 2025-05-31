package registry

import (
	"Comment/module/settings"
	"Comment/proto/pb"
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Registry struct {
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

// NewDBConnection db服务连接
// 注意需要后续手动close
func NewDBConnection() (*grpc.ClientConn, error) {
	register, err := NewRegistry()
	if err != nil {
		return nil, err
	}
	// registry.client.Health().Service返回的是对应服务的节点列表
	services, _, err := register.client.Health().Service("db-service", "db-service", true, nil)
	if err != nil {
		return nil, err
	}
	// 如果是集群的话，这里可以添加简单的负载均衡，访问压力均摊给集群中的每个服务
	dsn := fmt.Sprintf("%s:%d", services[0].Service.Address, services[0].Service.Port)
	return grpc.NewClient(dsn, grpc.WithTransportCredentials(insecure.NewCredentials())) // 不安全连接
}

func GetJudgeServerDsn(cfg *settings.RegistryConfig) (string, error) {
	register, err := NewRegistry()
	if err != nil {
		return "", err
	}
	// registry.client.Health().Service返回的是对应服务的节点列表
	services, _, err := register.client.Health().Service("judge-service", "judge-service", true, nil)
	if err != nil {
		return "", err
	}
	// 如果是集群的话，这里可以添加简单的负载均衡，访问压力均摊给集群中的每个服务
	dsn := fmt.Sprintf("%s:%d", services[0].Service.Address, services[0].Service.Port)
	return dsn, nil
}
