package internal

import (
	"Comment/app/common/serverBase"
	"Comment/app/gateway/internal/routers"
	"fmt"
)

type Server struct {
	serverBase.Server
}

func (receiver *Server) Init() error {
	err := receiver.Initialize()
	if err != nil {
		return err
	}

	return nil
}

func (receiver *Server) Start() {
	engine := routers.Router()
	dsn := fmt.Sprintf(":%d", receiver.Port)
	go func() {
		err := engine.Run(dsn)
		if err != nil {
			panic(err)
		}
	}()

	err := receiver.Register()
	if err != nil {
		panic(err)
	}

	select {}
}
