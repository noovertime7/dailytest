package plugin

import (
	"github.com/hashicorp/go-plugin"
	"net/rpc"
)

type MyPlugin interface {
	Run() string
}

var HandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

type MyPluginRPC struct {
	client *rpc.Client
}

func (g *MyPluginRPC) Run() string {
	var resp string
	err := g.client.Call("Plugin.Run", new(interface{}), &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}

	return resp
}

// Here is the RPC server that GreeterRPC talks to, conforming to
// the requirements of net/rpc
type MyPluginRPCRPCServer struct {
	// This is the real implementation
	Impl MyPlugin
}

func (s *MyPluginRPCRPCServer) Run(args interface{}, resp *string) error {
	*resp = s.Impl.Run()
	return nil
}

type MyPluginImpl struct {
	Impl MyPlugin
}

func (p *MyPluginImpl) Server(*plugin.MuxBroker) (interface{}, error) {
	return &MyPluginRPCRPCServer{Impl: p.Impl}, nil
}

func (MyPluginImpl) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &MyPluginRPC{client: c}, nil
}
