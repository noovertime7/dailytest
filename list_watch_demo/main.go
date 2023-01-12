package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	coreV1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func getClientSet() *kubernetes.Clientset {
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
	return clientSet
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func Must(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func main() {
	clientSet := getClientSet()
	ctx := context.Background()
	// kubectl run --image=nginx nginx-app --port=80
	podList, err := clientSet.CoreV1().Pods("default").List(ctx, metav1.ListOptions{})
	Must(err)
	for _, pod := range podList.Items {
		fmt.Printf("List获取到POD: %s\n", pod.Name)
	}

	fmt.Println("开始watch POD 变化...")
	w, err := clientSet.CoreV1().Pods("default").Watch(ctx, metav1.ListOptions{})
	Must(err)
	defer w.Stop()

	for {
		select {
		case event := <-w.ResultChan():
			pod, ok := event.Object.(*coreV1.Pod)
			if !ok {
				return
			}
			fmt.Printf("Watch 到 %s 变化,EventType: %s\n", pod.Name, event.Type)
		}
	}
}
