package main

import (
	"context"
	"log"
	"net"

	"github.com/microtwitch/chatedge/protos"
	"github.com/microtwitch/chatedge/receiver/edge"
	"github.com/microtwitch/chatedge/receiver/server"
	"google.golang.org/grpc"
)

const RECEIVER_TARGET string = "localhost:9090"
const EDGE_TARGET string = "localhost:8080"

func main() {
	lis, err := net.Listen("tcp", RECEIVER_TARGET)
	if err != nil {
		log.Fatalln(err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	server := server.NewServer()

	protos.RegisterEdgeReceiverServer(grpcServer, server)

	go grpcServer.Serve(lis)

	client, err := edge.NewChatEdgeClient(EDGE_TARGET)
	if err != nil {
		log.Fatalln(err)
	}

	err = client.JoinChat(context.Background(), "tmiloadtesting2", RECEIVER_TARGET)
	if err != nil {
		log.Fatalln(err)
	}

	for {
		select {}
	}

}
