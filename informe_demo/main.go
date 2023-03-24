package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
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

func main() {
	// 初始化clientSet与factory
	clientSet := getClientSet()
	factory := informers.NewSharedInformerFactoryWithOptions(clientSet, 0)

	deploymentInformer := factory.Apps().V1().Deployments()
	_ = deploymentInformer.Informer()
	deploymentLister := deploymentInformer.Lister()

	podInformer := factory.Core().V1().Pods()
	// 向factory注册podInformer
	sharedIndexInformer := podInformer.Informer()
	podLister := podInformer.Lister()

	sharedIndexInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			fmt.Println("add", obj.(*v1.Pod).Name)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			fmt.Println("update", oldObj.(*v1.Pod).Name)
		},
		DeleteFunc: func(obj interface{}) {
			fmt.Println("delete", obj.(*v1.Pod).Name)
		},
	})

	// 可以添加多个EventHandler
	sharedIndexInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			fmt.Println("add2", obj.(*v1.Pod).Name)
		},
	})

	stopCh := make(chan struct{})
	// 启动factory
	factory.Start(stopCh)
	defer close(stopCh)
	//等待所有的informer同步完成
	factory.WaitForCacheSync(stopCh)

	fmt.Println("从indexer缓存中获取POD数据")
	pods, err := podLister.Pods(v1.NamespaceAll).List(labels.Everything())
	if err != nil {
		log.Fatalln(err)
	}
	for index, pod := range pods {
		fmt.Println(index, "->", pod.Name)
	}

	fmt.Println("从indexer缓存中获取deployment数据")
	deployments, err := deploymentLister.List(labels.Everything())
	if err != nil {
		log.Fatalln(err)
	}
	for index, deployment := range deployments {
		fmt.Println(index, "->", deployment.Name)
	}
}
