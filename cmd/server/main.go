package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/alicerum/test-grpc/pkg/proto"
	"google.golang.org/grpc"
)

type greetingImpl struct {
	proto.UnimplementedGreetingServer
}

func (s greetingImpl) Hello(ctx context.Context, in *proto.UserInfo) (*proto.Response, error) {
	result := &proto.Response{}

	name := in.Name
	age := in.Age

	result.Result = fmt.Sprintf("Hello, user %q of age '%d'", name, age)
	return result, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Wrong number of arguments")
		os.Exit(1)
	}

	port, err := strconv.ParseInt(os.Args[1], 10, 32)
	if err != nil {
		fmt.Printf("Wrong port %q\n", os.Args[1])
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		fmt.Printf("Error while creating listener: %v", err)
		os.Exit(1)
	}

	var opts []grpc.ServerOption
	server := grpc.NewServer(opts...)
	proto.RegisterGreetingServer(server, greetingImpl{})
	if err := server.Serve(lis); err != nil {
		fmt.Printf("Error while serving the proto: %v\n", err)
		os.Exit(1)
	}
}
