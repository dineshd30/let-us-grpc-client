package grpcclient

import (
	"context"
	"errors"
	"io"
	"log"
	"time"

	"github.com/dineshd30/let-us-grpc-proto/proto"
	"google.golang.org/grpc"
)

type GreeterClient struct {
	client proto.GreeterServiceClient
}

func New(cc grpc.ClientConnInterface) *GreeterClient {
	return &GreeterClient{
		client: proto.NewGreeterServiceClient(cc),
	}
}

func (g *GreeterClient) SayHelloUniary() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := g.client.SayHelloUniary(ctx, &proto.HelloRequest{
		Message: "Dinesh",
	})
	if err != nil {
		log.Fatalf("failed to say hello to server: %v", err)
	}

	log.Printf("%s\n", res.Message)
}

func (g *GreeterClient) SayHelloServerStreaming() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*10))
	defer cancel()

	stream, err := g.client.SayHelloServerStreaming(ctx, &proto.NamesList{
		Name: []string{
			"Dinesh",
			"Bob",
			"Alice",
		},
	})
	if err != nil {
		log.Fatalf("failed to say hello to server: %v", err)
	}

	for {
		res, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			log.Println("reached end of stream")
			return
		}
		if err != nil {
			log.Fatalf("failed to get message from server steam: %v", err)
		}

		log.Println(res.Message)
	}
}

func (g *GreeterClient) SayHelloClientStreaming() {
	names := []string{
		"Dinesh", "Bob", "Alice",
	}
	stream, err := g.client.SayHelloClientStreaming(context.Background())
	if err != nil {
		log.Fatalf("failed to say hello to server: %v", err)
	}

	for _, name := range names {
		req := &proto.HelloRequest{
			Message: name,
		}

		if err := stream.Send(req); err != nil {
			log.Fatalf("failed to send message to server via client stream: %v", err)
		}
		log.Printf("Sent the request message to server - %s", name)
		time.Sleep(time.Second * 2)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("failed to get response from server : %v", err)
	}
	log.Println(res.Name)
}

func (g *GreeterClient) SayHelloBidirectionalStreaming() {
	names := []string{
		"Dinesh", "Bob", "Alice",
	}

	stream, err := g.client.SayHelloBidirectionalStreaming(context.Background())
	if err != nil {
		log.Fatalf("failed to say hello to server: %v", err)
	}

	go func() {
		for {
			res, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				log.Fatalf("failed to receive message from server via bidirectional stream: %v", err)
			}

			log.Println(res.Message)
		}
	}()

	for _, name := range names {
		req := &proto.HelloRequest{
			Message: name,
		}

		if err := stream.Send(req); err != nil {
			log.Fatalf("failed to send message to server via bidirectional stream: %v", err)
		}
		log.Printf("sent the request message to server - %s", name)
		time.Sleep(time.Second * 2)
	}

	stream.CloseSend()
	log.Println("bidirectional stream completed successfully")
}
