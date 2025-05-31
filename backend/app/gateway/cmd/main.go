package main

import "Comment/app/gateway/internal"

func main() {
	server := internal.Server{}
	server.Name = "gateway-service"
	server.SrvType = "http"

	err := server.Init()
	if err != nil {
		panic(err)
	}
	server.Start()
}
