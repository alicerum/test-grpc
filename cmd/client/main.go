package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/alicerum/test-grpc/pkg/proto"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

type Client struct {
	rootCmd *cobra.Command

	port int
	name string
	age  int
}

func (c *Client) runCommand(cmd *cobra.Command, args []string) error {
	if c.port == -1 {
		return errors.New("Must set port")
	}
	if len(c.name) == 0 {
		return errors.New("Must set name")
	}
	if c.age == -1 {
		return errors.New("Must set age")
	}

	var callOpts []grpc.CallOption

	serverAddr := fmt.Sprintf("localhost:%d", c.port)
	client, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("Error while dialing server: %v\n", err)
	}
	defer client.Close()

	message := &proto.UserInfo{
		Name: c.name,
		Age:  int32(c.age),
	}

	cc := proto.NewGreetingClient(client)
	resp, err := cc.Hello(context.TODO(), message, callOpts...)
	if err != nil {
		return fmt.Errorf("Error while executing gRPC request: %v\n", err)
	}

	fmt.Println("Response from server!")
	fmt.Println(resp.Result)

	return nil
}

func (c *Client) setUpCobra() {
	c.rootCmd = &cobra.Command{
		Use:  os.Args[0],
		RunE: c.runCommand,
	}

	c.rootCmd.Flags().IntVarP(&c.port, "port", "p", -1, "gRPC port")
	c.rootCmd.Flags().StringVar(&c.name, "name", "", "User name for the message")
	c.rootCmd.Flags().IntVar(&c.age, "age", -1, "User age for the message")
}

func main() {
	rootOptions := &Client{}
	rootOptions.setUpCobra()

	if err := rootOptions.rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
