package main

import (
	"fmt"
	"github.com/hashicorp/go-hclog"
	"os"
	"os/exec"
	"path/filepath"
	plugin2 "plugin/myplugin/plugin"

	"github.com/hashicorp/go-plugin"
)

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Debug,
	})
	// 设置插件路径
	pluginPath, _ := filepath.Abs("./myplugin")
	os.Setenv("PLUGIN_PATH", pluginPath)

	// 创建并启动插件进程
	cmd := exec.Command("./myplugin/KV-1.0.2.exe")
	pluginClient := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: plugin2.HandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			"myplugin": &plugin2.MyPluginImpl{},
		},
		Logger: logger,
		Cmd:    cmd,
	})
	rpcClient, err := pluginClient.Client()
	if err != nil {
		fmt.Println("Failed to start plugin:", err)
		return
	}
	defer pluginClient.Kill()

	// 连接插件进程并调用方法
	raw, err := rpcClient.Dispense("myplugin")
	if err != nil {
		fmt.Println("Failed to dispense plugin:", err)
		return
	}
	myPlugin := raw.(plugin2.MyPlugin)
	res := myPlugin.Run()
	fmt.Println(res)
}
