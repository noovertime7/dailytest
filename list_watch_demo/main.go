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
	if len(podList.Items) == 0 {
		fmt.Println("当前命名空间下POD资源为空")
	} else {
		for _, pod := range podList.Items {
			fmt.Printf("List获取到POD: %s\n", pod.Name)
		}
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
			fmt.Println("------------")
		}
	}
}

// 注意: 每次resourceVersion 都会增加
//新增 ADD
//修改 Modify  调度完成 绑定nodeName，更新status字段
//修改 Modify  status 增加更多字段，主要是修改status的message信息，例如 "containers with unready status: [nginx-app]"
//修改 Modify  从cni网络插件中获取到POD的ip地址并填充到status中，status的message置为空,修改pod运行状态为Running
