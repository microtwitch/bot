package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/microtwitch/bot/config"
	"github.com/microtwitch/chatedge/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const RECEIVER_TARGET string = "127.0.0.1:9090"
const EDGE_TARGET string = "127.0.0.1:8080"

func main() {
	config.Init()

	lis, err := net.Listen("tcp", RECEIVER_TARGET)
	if err != nil {
		log.Fatalln(err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	server := NewServer(HandleMessage)

	protos.RegisterEdgeReceiverServer(grpcServer, server)

	go grpcServer.Serve(lis)

	client, err := NewChatEdgeClient(EDGE_TARGET)
	if err != nil {
		log.Fatalln(err)
	}

	err = client.JoinChat(context.Background(), config.Channel, RECEIVER_TARGET)
	if err != nil {
		log.Fatalln(err)
	}

	for {
		select {}
	}

}

func HandleMessage(msg *protos.ChatMessage) {
	log.Println(fmt.Sprintf("#%s %s: %s", msg.Channel, msg.User, msg.Message))
}

type receiverServer struct {
	protos.UnimplementedEdgeReceiverServer

	handleMsg func(*protos.ChatMessage)
}

func NewServer(handleMsg func(*protos.ChatMessage)) *receiverServer {
	s := &receiverServer{handleMsg: handleMsg}
	return s
}

func (s *receiverServer) Send(ctx context.Context, chatMessage *protos.ChatMessage) (*protos.Empty, error) {
	s.handleMsg(chatMessage)
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
