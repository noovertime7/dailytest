package main

import (
	"github.com/hashicorp/go-plugin"
	"github.com/sirupsen/logrus"
	"plugin-demo/common"
	"plugin-demo/grpcplugin/plugin/v1"
)

func Handshake() plugin.HandshakeConfig {
	return plugin.HandshakeConfig{
		// The ProtocolVersion is the version that must match between Velero framework
		// and Velero client plugins. This should be bumped whenever a change happens in
		// one or the other that makes it so that they can't safely communicate.
		ProtocolVersion: 2,

		MagicCookieKey:   "VELERO_PLUGIN",
		MagicCookieValue: "hello",
	}
}

func main() {
	log := logrus.New()
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: Handshake(),
		Plugins: map[string]plugin.Plugin{
			string(common.PluginKV): v1.NewPlugin(common.ServerLogger(log)),
			string(common.HellO):    v1.NewHelloPlugin(common.ServerLogger(log)),
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
