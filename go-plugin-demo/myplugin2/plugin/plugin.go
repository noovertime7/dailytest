package main

import (
	"github.com/hashicorp/go-plugin"
	"net/rpc"
	"plugin/myplugin2/types"
)

type GreeterPlugin struct{}

func (GreeterPlugin) Greet() string {
	return "Hello!"
}

type MyPluginImpl struct {
	Impl types.OS
}

func (p *MyPluginImpl) Server(*plugin.MuxBroker) (interface{}, error) {
	return &MyPluginImpl{Impl: p.Impl}, nil
}

func (MyPluginImpl) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return nil, nil
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  1,
			MagicCookieKey:   "GREETER_PLUGIN",
			MagicCookieValue: "hello",
		},
		Plugins: map[string]plugin.Plugin{
			"os": &MyPluginImpl{
				Impl: &GreeterPlugin{},
			},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
