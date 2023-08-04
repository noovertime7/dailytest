package cluster

import (
	"context"
	"fmt"
	"k8s.io/client-go/tools/clientcmd"
	"strings"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	metricsvebeta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/cluster"
)

type Interface interface {
	cluster.Cluster
	Kubernetes() kubernetes.Interface
	Discovery() discovery.CachedDiscoveryInterface
	Watch(ctx context.Context, list client.ObjectList, callback func(watch.Event) error, opts ...client.ListOption) error
}

type Cluster struct {
	cluster.Cluster
	discovery  discovery.CachedDiscoveryInterface
	kubernetes kubernetes.Interface
}

func WithDisableCaches() func(o *cluster.Options) {
	disabled := []client.Object{
		&metricsvebeta1.NodeMetrics{},
		&metricsvebeta1.PodMetrics{},
	}
	return func(o *cluster.Options) { o.ClientDisableCacheFor = append(o.ClientDisableCacheFor, disabled...) }
}

type WatchableDelegatingClient struct {
	client.Client
	watchable client.WithWatch
}

func (c *WatchableDelegatingClient) Watch(ctx context.Context, obj client.ObjectList, opts ...client.ListOption) (watch.Interface, error) {
	return c.watchable.Watch(ctx, obj, opts...)
}

func WithInNamespace(ns string) func(o *cluster.Options) {
	return func(o *cluster.Options) {
		o.Namespace = ns
	}
}

func NewLocalAgentClusterAndStart(ctx context.Context, options ...cluster.Option) (*Cluster, error) {
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(), nil).ClientConfig()
	if err != nil {
		return nil, err
	}
	apply := func(c cluster.Cluster) error {
		return nil
	}
	return NewClusterAndStart(ctx, config, apply, options...)
}

func NewClusterAndStart(ctx context.Context, config *rest.Config, apply func(c cluster.Cluster) error, options ...cluster.Option) (*Cluster, error) {
	c, err := NewCluster(config, options...)
	if err != nil {
		return nil, err
	}
	if err := apply(c); err != nil {
		return nil, err
	}
	go c.Start(ctx)
	c.GetCache().WaitForCacheSync(ctx)
	return c, nil
}

func NewCluster(config *rest.Config, options ...cluster.Option) (*Cluster, error) {
	discovery, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create discovery client: %w", err)
	}

	options = append(options,
		WithDisableCaches())

	kubernetesClientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	c, err := cluster.New(config, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create cluster: %w", err)
	}
	return &Cluster{
		Cluster:    c,
		kubernetes: kubernetesClientSet,
		discovery:  memory.NewMemCacheClient(discovery),
	}, nil
}

func (c *Cluster) Kubernetes() kubernetes.Interface {
	return c.kubernetes
}

func (c *Cluster) Discovery() discovery.CachedDiscoveryInterface {
	return c.discovery
}

func (c *Cluster) Watch(ctx context.Context, list client.ObjectList, callback func(watch.Event) error, opts ...client.ListOption) error {
	gvk, err := apiutil.GVKForObject(list, c.GetScheme())
	if err != nil {
		return err
	}
	gvk.Kind = strings.TrimSuffix(gvk.Kind, "List")

	mapping, err := c.Cluster.GetRESTMapper().RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return err
	}

	if callback == nil {
		return errors.NewBadRequest("no callback provided")
	}

	listOpts := client.ListOptions{}
	listOpts.ApplyOptions(opts)

	config := c.GetConfig()
	nclient, err := dynamic.NewForConfig(config)
	if err != nil {
		return err
	}

	watcher, err := nclient.
		Resource(mapping.Resource).
		Namespace(listOpts.Namespace).
		Watch(ctx, *listOpts.AsListOptions())
	if err != nil {
		return err
	}
	defer watcher.Stop()

	for {
		select {
		case event := <-watcher.ResultChan():
			if err := callback(event); err != nil {
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}
}
