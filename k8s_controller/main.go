package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
)

// Controller POD控制器
type Controller struct {
	indexer  cache.Indexer
	queue    workqueue.RateLimitingInterface
	informer cache.Controller
}

func NewController(queue workqueue.RateLimitingInterface, indexer cache.Indexer, informer cache.Controller) *Controller {
	return &Controller{
		indexer:  indexer,
		queue:    queue,
		informer: informer,
	}
}

func (c *Controller) Run(threadiness int, stopCh chan struct{}) {
	defer runtime.HandleCrash()

	// 停止控制器后需要关闭queue
	defer c.queue.ShutDown()

	// 启动控制器
	klog.Infof("启动 pod controller")
	go c.informer.Run(stopCh)

	//等待所有相关的缓存同步完成，然后再开始处理队列中的数据
	if !cache.WaitForCacheSync(stopCh, c.informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("time out waiting for caches to sync"))
	}
	// 启动worker，处理元素
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}
	<-stopCh
	klog.Info("pod controller stopping")
}

func (c *Controller) runWorker() {
	for c.processNextItem() {

	}
}

func (c *Controller) processNextItem() bool {
	//从workQueue里面一个取出元素
	key, quit := c.queue.Get()
	//如果这个队列已经关闭了，返回一个false
	if quit {
		return false
	}
	//告诉队列已经处理了该key
	defer c.queue.Done(key)
	//根据key去处理我们的业务逻辑了
	err := c.syncToStdout(key.(string))
	c.handleErr(err, key)
	return true
}

// 业务逻辑处理
func (c *Controller) syncToStdout(key string) error {
	obj, exists, err := c.indexer.GetByKey(key)
	if err != nil {
		klog.Errorf("%s从缓存中获取对象失败%v", key, err)
		return err
	}
	if !exists {
		fmt.Printf("Pod %s 已经不存在了\n", key)
	} else {
		fmt.Printf("sync/add/update for pod %s\n", obj.(*v1.Pod).Name)
	}
	return nil
}

func (c *Controller) handleErr(err error, key interface{}) bool {
	if err == nil {
		c.queue.Forget(key)
		return false
	}
	//如果出现了问题，我们允许当前控制器重试5词
	if c.queue.NumRequeues(key) < 5 {
		//重新入队列
		c.queue.AddRateLimited(key)
		return false
	}
	c.queue.Forget(key)
	runtime.HandleError(err)
	//不允许继续重试了
	return true
}

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
	clientSet := getClientSet()
	//创建pod的listWatcher
	podlistWatcher := cache.NewListWatchFromClient(clientSet.CoreV1().RESTClient(), "pods", v1.NamespaceDefault, fields.Everything())
	//创建队列
	workQueue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	indexer, informer := cache.NewIndexerInformer(podlistWatcher, &v1.Pod{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			fmt.Println("新增", obj.(*v1.Pod).Name)
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err != nil {
				log.Fatalln(err)
			}
			workQueue.Add(key)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			fmt.Println("更新", oldObj.(*v1.Pod).Name)
			key, err := cache.MetaNamespaceKeyFunc(newObj)
			if err != nil {
				log.Fatalln(err)
			}
			workQueue.Add(key)
		},
		DeleteFunc: func(obj interface{}) {
			fmt.Println("删除", obj.(*v1.Pod).Name)
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err != nil {
				log.Fatalln(err)
			}
			workQueue.Add(key)
		},
	}, cache.Indexers{})
	//实例化控制器
	c := NewController(workQueue, indexer, informer)
	stopCh := make(chan struct{})
	go c.Run(10, stopCh)
	select {}
}
