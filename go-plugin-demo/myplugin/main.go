package main

import (
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"os"
	plugin2 "plugin/myplugin/plugin"
)

type MyPluginBasic struct {
}

func (m *MyPluginBasic) Run() string {
	return "MyPluginBasic Run "
}

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})
	runner := &MyPluginBasic{}
	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"myplugin": &plugin2.MyPluginImpl{Impl: runner},
	}
	logger.Debug("message from plugin", "foo", "bar")
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugin2.HandshakeConfig,
		Plugins:         pluginMap,
	})
}
