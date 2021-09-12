package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"

	"github.com/alicerum/test-grpc/pkg/proto"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

type RootOptions struct {
	rootCmd *cobra.Command
	port    int
}

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

func (c *RootOptions) runServer(cmd *cobra.Command, args []string) error {
	if c.port == -1 {
		return errors.New("Must set up port")
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", c.port))
	if err != nil {
		return fmt.Errorf("Error while creating listener: %v", err)
	}

	var opts []grpc.ServerOption
	server := grpc.NewServer(opts...)
	proto.RegisterGreetingServer(server, greetingImpl{})
	if err := server.Serve(lis); err != nil {
		return fmt.Errorf("Error while serving the proto: %v\n", err)
	}

	return nil
}

func (c *RootOptions) setUpCobra() {
	c.rootCmd = &cobra.Command{
		Use:  os.Args[0],
		RunE: c.runServer,
	}

	c.rootCmd.Flags().IntVarP(&c.port, "port", "p", -1, "Server port to start")
}

func main() {
	rootOptions := &RootOptions{}
	rootOptions.setUpCobra()

	if err := rootOptions.rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
