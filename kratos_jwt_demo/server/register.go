package server

import (
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

func RegistrarServer() (*etcd.Registry, func(), error) {
	// new etcd client
	client, err := clientv3.New(clientv3.Config{
		DialTimeout: 3 * time.Second,
		Endpoints:   []string{"127.0.0.1:2379"},
		Username:    "root",
		Password:    "3s4d5cvybnm3AS4gyi5F6G7H8J",
	})
	if err != nil {
		return nil, nil, err
	}
	// new reg with etcd client
	return etcd.New(client), func() {
		return
	}, nil
}
