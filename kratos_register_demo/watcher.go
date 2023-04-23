package main

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
)

func WatchService(ctx context.Context, serviceName string, reg registry.Discovery) {
	watcher, err := reg.Watch(ctx, serviceName)
	if err != nil {
		log.Errorf("get watcher error %v", err)
		return
	}
	for {
		res, err := watcher.Next()
		if err != nil {
			return
		}
		fmt.Println(&res)
	}
}
