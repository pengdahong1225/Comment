package services

import "google.golang.org/grpc"

var (
	Instance *ConnectionManager
)

func Init(consulAddress string) {
	Instance = &ConnectionManager{
		consulAddress: consulAddress,
		conns:         make(map[string]*grpc.ClientConn),
	}
}
