package main

import "Comment/app/sse/internal"

func main() {
	server := internal.Server{}
	server.Name = "sse-service"
	server.SrvType = "http"

	err := server.Init()
	if err != nil {
		panic(err)
	}
	server.Start()
}
