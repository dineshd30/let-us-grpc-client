package main

import (
	"fmt"
	"log"

	"github.com/dineshd30/let-us-grpc-client/internal/grpcclient"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	port = "8080"
)

func main() {
	cc, err := grpc.Dial(fmt.Sprintf("localhost:%s", port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to create connection to server: %v", err)
	}

	defer cc.Close()

	client := grpcclient.New(cc)
	// client.SayHelloUniary()
	// client.SayHelloServerStreaming()
	// client.SayHelloClientStreaming()
	client.SayHelloBidirectionalStreaming()
}
