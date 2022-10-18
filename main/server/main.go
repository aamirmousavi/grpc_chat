package main

import (
	"grpc_chat/internal/chat"
	"grpc_chat/internal/hub"
	"net"

	"google.golang.org/grpc"
)

func main() {
	srv := grpc.NewServer()
	chat.RegisterChatServer(srv, hub.Hub)

	lis, err := net.Listen("tcp", ":6565")
	if err != nil {
		panic(err)
	}

	if err := srv.Serve(lis); err != nil {
		panic(err)
	}
}
