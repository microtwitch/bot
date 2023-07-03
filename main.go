package main

import (
	"context"
	"log"
	"net"
	"strings"

	"github.com/microtwitch/bot/config"
	"github.com/microtwitch/chatedge/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	config.Init()

	lis, err := net.Listen("tcp", config.ReceiverTarget)
	if err != nil {
		log.Fatalln(err)
	}

	bot := Bot{}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	server := newServer(bot.handleMessage)

	protos.RegisterEdgeReceiverServer(grpcServer, server)

	go grpcServer.Serve(lis)

	client, err := newChatEdgeClient(config.EdgeTarget)
	if err != nil {
		log.Fatalln(err)
	}

	bot.client = client

	err = bot.client.JoinChat(context.Background(), config.ReceiverTarget)
	if err != nil {
		log.Fatalln(err)
	}

	err = bot.client.Send(context.Background(), "OpieOP bot ready!")
	if err != nil {
		log.Fatalln(err)
	}

	for {
		select {}
	}

}

type Bot struct {
	client *ChatEdgeClient
}

func (bot *Bot) handleMessage(msg *protos.ChatMessage) {
	if msg.User != config.Admin {
		return
	}

	parts := strings.Split(msg.Message, " ")

	switch parts[0] {
	case ";ping":
		bot.client.Send(context.Background(), "PONG!")
	}
}

type receiverServer struct {
	protos.UnimplementedEdgeReceiverServer

	handleMsg func(*protos.ChatMessage)
}

func newServer(handleMsg func(*protos.ChatMessage)) *receiverServer {
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

func newChatEdgeClient(target string) (*ChatEdgeClient, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial(target, opts...)
	if err != nil {
		return nil, err
	}

	client := protos.NewChatEdgeClient(conn)

	return &ChatEdgeClient{client}, nil
}

func (c *ChatEdgeClient) JoinChat(ctx context.Context, callback string) error {
	joinRequest := protos.JoinRequest{Channel: config.Channel, Callback: callback}
	_, err := c.client.JoinChat(ctx, &joinRequest)
	if err != nil {
		return err
	}

	// TODO: do something with the id in resp

	return nil
}

func (c *ChatEdgeClient) Send(ctx context.Context, msg string) error {
	sendRequest := protos.SendRequest{Token: config.Token, User: config.BotUser, Channel: config.Channel, Msg: msg}
	_, err := c.client.Send(ctx, &sendRequest)
	return err
}
