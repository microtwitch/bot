package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/microtwitch/chatedge/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const RECEIVER_TARGET string = "127.0.0.1:9090"
const EDGE_TARGET string = "127.0.0.1:8080"

func main() {
	lis, err := net.Listen("tcp", RECEIVER_TARGET)
	if err != nil {
		log.Fatalln(err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	server := NewServer()

	protos.RegisterEdgeReceiverServer(grpcServer, server)

	go grpcServer.Serve(lis)

	client, err := NewChatEdgeClient(EDGE_TARGET)
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

type receiverServer struct {
	protos.UnimplementedEdgeReceiverServer
}

func NewServer() *receiverServer {
	s := &receiverServer{}
	return s
}

func (s *receiverServer) Send(ctx context.Context, chatMessage *protos.ChatMessage) (*protos.Empty, error) {
	log.Println(fmt.Sprintf("#%s %s: %s", chatMessage.Channel, chatMessage.User, chatMessage.Message))
	return &protos.Empty{}, nil
}

type ChatEdgeClient struct {
	client protos.ChatEdgeClient
}

func NewChatEdgeClient(target string) (*ChatEdgeClient, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial(target, opts...)
	if err != nil {
		return nil, err
	}

	client := protos.NewChatEdgeClient(conn)

	return &ChatEdgeClient{client}, nil
}

func (c *ChatEdgeClient) JoinChat(ctx context.Context, channel string, callback string) error {
	joinRequest := protos.JoinRequest{Channel: channel, Callback: callback}
	_, err := c.client.JoinChat(ctx, &joinRequest)
	if err != nil {
		return err
	}

	// TODO: do something with the id in resp

	return nil
}
