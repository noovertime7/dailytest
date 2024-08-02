package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"plugin-demo/grpcplugin/manager"
	"plugin-demo/grpcplugin/process"
)

func main() {
	log := logrus.New()

	reg := process.NewRegistry("E:\\code\\dailytest\\plugin-demo\\grpcplugin\\target\\", log, logrus.DebugLevel)
	err := reg.DiscoverPlugins()
	if err != nil {
		log.Fatal(err)
	}

	m := manager.NewManager(log, logrus.DebugLevel, reg)

	res, err := m.GetKV("1.0.3")
	if err != nil {
		log.Fatal(err)
	}
	defer m.CleanupClients()
	resp, err := res.Get("abc")
	fmt.Println(string(resp), err)

	//client := plugin.NewClient(&plugin.ClientConfig{
	//	HandshakeConfig:  Handshake(),
	//	AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
	//	Plugins: map[string]plugin.Plugin{
	//		string(common.PluginKV): v1.NewPlugin(common.ServerLogger(log)),
	//		string(common.HellO):    v1.NewHelloPlugin(common.ServerLogger(log)),
	//	},
	//	Cmd: exec.Command("E:\\code\\dailytest\\plugin-demo\\grpcplugin\\target\\KV-1.0.2.exe"),
	//})
	//defer client.Kill()
	//
	//// Connect via RPC
	//rpcClient, err := client.Client()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//// Request the plugin
	//raw, err := rpcClient.Dispense("KV")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(raw)
	//dis := raw.(common.ClientDispenser)
	//
	//kvServer := dis.ClientFor("KV").(KV)
	//res, err := kvServer.Get("abc")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(string(res))

}

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
