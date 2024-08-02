package server

import (
	"fmt"
	shared "plugin-demo/shard"
)

// Here is the RPC server that GreeterRPC talks to, conforming to
// the requirements of net/rpc
type GreeterRPCServer struct {
	// This is the real implementation
	Impl shared.Greeter
}

func (s *GreeterRPCServer) Greet(args interface{}, resp *string) error {
	fmt.Println("Greet GreeterRPCServer")
	*resp = s.Impl.Greet()
	return nil
}

func (s *GreeterRPCServer) SayHello(args interface{}, resp *string) error {
	fmt.Println("Greet GreeterRPCServer")
	*resp = s.Impl.SayHello()
	return nil
}
