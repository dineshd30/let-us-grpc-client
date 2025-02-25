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

func (g *GreeterClient) SayHelloUnary() {
	log.Println("say hello unary")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	name := "Hulk"
	log.Println("sending uniary request:", name)
	res, err := g.client.SayHelloUnary(ctx, &proto.HelloRequest{
		Message: name,
	})
	if err != nil {
		log.Fatalf("failed to say hello to server: %v", err)
	}

	log.Println("got response:", res.Message)
}

func (g *GreeterClient) SayHelloServerStreaming() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*20))
	defer cancel()

	names := []string{
		"Thor", "Hulk", "Ultron", "Iron Man", "Thanos",
	}
	log.Println("sending server streaming request:", names)
	stream, err := g.client.SayHelloServerStreaming(ctx, &proto.NamesList{
		Name: names,
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

		log.Println("got response:", res.Message)
	}
}

func (g *GreeterClient) SayHelloClientStreaming() {
	names := []string{
		"Thor", "Hulk", "Black Widow", "Iron Man", "Captain Marvel",
	}
	stream, err := g.client.SayHelloClientStreaming(context.Background())
	if err != nil {
		log.Fatalf("failed to say hello to server: %v", err)
	}

	for _, name := range names {
		req := &proto.HelloRequest{
			Message: name,
		}

		log.Println("sending client streaming request:", name)
		if err := stream.Send(req); err != nil {
			log.Fatalf("failed to send message to server via client stream: %v", err)
		}
		time.Sleep(time.Second * 2)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("failed to get response from server : %v", err)
	}
	log.Println("got response:", res.Name)
}

func (g *GreeterClient) SayHelloBidirectionalStreaming() {
	names := []string{
		"Loki", "Hulk", "Black Widow", "Iron Man", "Captain Marvel",
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

			log.Println("got response:", res.Message)
		}
	}()

	for _, name := range names {
		req := &proto.HelloRequest{
			Message: name,
		}

		log.Println("sending bidirectional streaming request:", name)
		if err := stream.Send(req); err != nil {
			log.Fatalf("failed to send message to server via bidirectional stream: %v", err)
		}
		time.Sleep(time.Second * 2)
	}

	stream.CloseSend()
	log.Println("bidirectional stream completed successfully")
}
