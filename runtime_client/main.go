package main

import (
	"context"
	"fmt"
	"github.com/noovertime7/dailytest/informer_demo"
	v1 "k8s.io/api/apps/v1"
	runtimecluster "sigs.k8s.io/controller-runtime/pkg/cluster"
)

func main() {
	_, config := informer_demo.GetClientSet()
	c, err := runtimecluster.New(config)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()

	go c.Start(ctx)
	c.GetCache().WaitForCacheSync(ctx)

	out := &v1.DeploymentList{}
	if err := c.GetCache().List(ctx, out); err != nil {
		panic(err)
	}
	fmt.Println("list deployment success", out)

}
