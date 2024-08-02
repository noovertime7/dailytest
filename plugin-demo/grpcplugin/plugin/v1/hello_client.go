package v1

import (
	"context"
	"google.golang.org/grpc"
	"plugin-demo/common"
	"plugin-demo/proto"
)

type HelloPluginClient struct {
	Context context.Context
	*common.ClientBase
	grpcClient proto.HelloClient
}

func newHelloPluginClient(base *common.ClientBase, clientConn *grpc.ClientConn) interface{} {
	return &HelloPluginClient{
		Context:    context.Background(),
		ClientBase: base,
		grpcClient: proto.NewHelloClient(clientConn),
	}
}

func (c *HelloPluginClient) SayHello(ctx context.Context, req *proto.HelloRequest) (*proto.HelloResponse, error) {
	return c.grpcClient.SayHello(ctx, req)
}
