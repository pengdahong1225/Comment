package services

import (
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sync"
)

type ConnectionManager struct {
	// 注册中心地址
	consulAddress string

	// 服务连接
	conns map[string]*grpc.ClientConn
	mu    sync.RWMutex
}

func (m *ConnectionManager) GetConnection(serviceName string) (*grpc.ClientConn, error) {
	m.mu.RLock()
	if conn, ok := m.conns[serviceName]; ok {
		m.mu.RUnlock()
		return conn, nil
	}
	m.mu.RUnlock()

	m.mu.Lock()
	defer m.mu.Unlock()

	// 双重检查
	if conn, ok := m.conns[serviceName]; ok {
		return conn, nil
	}

	// 创建服务连接
	target := fmt.Sprintf("consul://%s/%s?wait=10s&healthy=true", m.consulAddress, serviceName)
	conn, err := grpc.NewClient(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),                // 无认证连接
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`), // 负载均衡，轮训策略
	)
	if err != nil {
		logrus.Errorf("Failed to create connection with %s: %v", serviceName, err)
		return nil, err
	}
	m.conns[serviceName] = conn

	return conn, err
}

func (m *ConnectionManager) GetConnectionByNodeId(serviceName string, nodeId string) (*grpc.ClientConn, error) {
	return nil, nil
}
