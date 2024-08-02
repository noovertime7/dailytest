package v1

import (
	"context"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	"plugin-demo/common"
	"plugin-demo/proto"
)

type HelloPlugin struct {
	plugin.NetRPCUnsupportedPlugin
	*common.PluginBase
}

func NewHelloPlugin(options ...common.PluginOption) *HelloPlugin {
	return &HelloPlugin{
		PluginBase: common.NewPluginBase(options...),
	}
}

// GRPCClient returns a clientDispenser for BackupItemAction gRPC clients.
func (p *HelloPlugin) GRPCClient(_ context.Context, _ *plugin.GRPCBroker, clientConn *grpc.ClientConn) (interface{}, error) {
	return common.NewClientDispenser(p.ClientLogger, clientConn, newHelloPluginClient), nil
}

// GRPCServer registers a BackupItemAction gRPC server.
func (p *HelloPlugin) GRPCServer(_ *plugin.GRPCBroker, server *grpc.Server) error {
	proto.RegisterHelloServer(server, &HelloPluginServer{mux: p.ServerMux})
	return nil
}
