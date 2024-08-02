package v1

import (
	"context"
	"plugin-demo/common"
	"plugin-demo/proto"
)

type HelloPluginServer struct {
	mux *common.ServerMux
	proto.UnimplementedHelloServer
}

func (c *HelloPluginServer) SayHello(ctx context.Context, req *proto.HelloRequest) (*proto.HelloResponse, error) {
	return &proto.HelloResponse{Message: "hello" + req.Name}, nil
}
