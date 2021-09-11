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

func (c *Client) fill() {
	flag.IntVar(&c.Port, "port", 5000, "gRPC port")
	flag.StringVar(&c.Name, "name", "", "User name for the message")
	flag.IntVar(&c.Age, "age", 0, "User age for the message")

	flag.Parse()
}

func main() {
	c := Client{}
	c.fill()

	var callOpts []grpc.CallOption

	serverAddr := fmt.Sprintf("localhost:%d", c.Port)
	client, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		fmt.Printf("Error while dialing server: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	cc := proto.NewGreetingClient(client)
	resp, err := cc.Hello(context.TODO(), &proto.UserInfo{Name: c.Name, Age: int32(c.Age)}, callOpts...)
	if err != nil {
		fmt.Printf("Error while executing gRPC request: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Response from server!")
	fmt.Println(resp.Result)
}
