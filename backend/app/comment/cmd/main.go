package main

import "Comment/app/comment/internal"

func main() {
	server := internal.Server{}
	server.Name = "comment-service"
	server.SrvType = "grpc"

	err := server.Init()
	if err != nil {
		panic(err)
	}
	server.Start()
}
