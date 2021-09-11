package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/alicerum/test-grpc/pkg/proto"
	"google.golang.org/grpc"
)

type Client struct {
	Port int
	Name string
	Age  int
}

var (
	clientOpts Client
)

func init() {
	clientOpts.fill()
}

func (c *Client) fill() {
	flag.IntVar(&c.Port, "port", 5000, "gRPC port")
	flag.StringVar(&c.Name, "name", "", "User name for the message")
	flag.IntVar(&c.Age, "age", 0, "User age for the message")
}

func main() {
	flag.Parse()

	var callOpts []grpc.CallOption

	serverAddr := fmt.Sprintf("localhost:%d", clientOpts.Port)
	client, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		fmt.Printf("Error while dialing server: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	message := &proto.UserInfo{
		Name: clientOpts.Name,
		Age:  int32(clientOpts.Age),
	}

	cc := proto.NewGreetingClient(client)
	resp, err := cc.Hello(context.TODO(), message, callOpts...)
	if err != nil {
		fmt.Printf("Error while executing gRPC request: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Response from server!")
	fmt.Println(resp.Result)
}
