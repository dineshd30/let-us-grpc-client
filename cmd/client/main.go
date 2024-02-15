package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"

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

	for {
		fmt.Println("---------------------------------------")
		fmt.Println("Enter 1 for SayHelloUnary()")
		fmt.Println("Enter 2 for SayHelloServerStreaming()")
		fmt.Println("Enter 3 for SayHelloClientStreaming()")
		fmt.Println("Enter 4 for SayHelloBidirectionalStreaming()")
		fmt.Println("Enter 5 to exit")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		option, _ := strconv.Atoi(scanner.Text())
		fmt.Println("---------------------------------------")
		switch option {
		case 1:
			client.SayHelloUnary()
		case 2:
			client.SayHelloServerStreaming()
		case 3:
			client.SayHelloClientStreaming()
		case 4:
			client.SayHelloBidirectionalStreaming()
		case 5:
			os.Exit(0)
		default:
			fmt.Println("Incorrect option")
		}
	}

}
