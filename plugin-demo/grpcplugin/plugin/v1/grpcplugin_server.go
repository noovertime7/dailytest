package v1

import (
	"context"
	"fmt"
	"plugin-demo/common"
	"plugin-demo/proto"
)

type PluginServer struct {
	mux *common.ServerMux
}

func (PluginServer) Get(context.Context, *proto.GetRequest) (*proto.GetResponse, error) {
	return &proto.GetResponse{
		Value: []byte("hello new  new"),
	}, nil
}
func (PluginServer) Put(context.Context, *proto.PutRequest) (*proto.Empty, error) {
	fmt.Println("Put")
	return nil, nil
}
