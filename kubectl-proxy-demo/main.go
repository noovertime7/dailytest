package main

import (
	"flag"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"kubectl-proxy-demo/server"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {

	_, cfg := GetClientSet()

	s, err := server.NewProxyHandler("/", nil, cfg, 0, false)
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", s)

	log.Fatal(http.ListenAndServe(":8082", nil))

}

func GetClientSet() (*kubernetes.Clientset, *rest.Config) {
	var err error
	var config *rest.Config
	var kubeConfig *string
	if home := homeDir(); home != "" {
		kubeConfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeConfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	// 使用 ServiceAccount 创建集群配置（InCluster模式）
	if config, err = rest.InClusterConfig(); err != nil {
		// 使用 KubeConfig 文件创建集群配置
		if config, err = clientcmd.BuildConfigFromFlags("", *kubeConfig); err != nil {
			panic(err.Error())
		}
	}
	// 创建 clientSet
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientSet, config
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
