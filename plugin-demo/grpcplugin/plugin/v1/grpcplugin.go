package v1

import (
	"context"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	"plugin-demo/common"
	"plugin-demo/proto"
)

type Plugin struct {
	plugin.NetRPCUnsupportedPlugin
	*common.PluginBase
}

func NewPlugin(options ...common.PluginOption) *Plugin {
	return &Plugin{
		PluginBase: common.NewPluginBase(options...),
	}
}

// GRPCClient returns a clientDispenser for BackupItemAction gRPC clients.
func (p *Plugin) GRPCClient(_ context.Context, _ *plugin.GRPCBroker, clientConn *grpc.ClientConn) (interface{}, error) {
	return common.NewClientDispenser(p.ClientLogger, clientConn, newBackupItemActionGRPCClient), nil
}

// GRPCServer registers a BackupItemAction gRPC server.
func (p *Plugin) GRPCServer(_ *plugin.GRPCBroker, server *grpc.Server) error {
	proto.RegisterKVServer(server, &PluginServer{mux: p.ServerMux})
	return nil
}
