package ioc

import (
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// etcdctl --endpoints=localhost:12379 get service/interactive --prefix

func InitEtcdClient() *clientv3.Client {
	var config clientv3.Config
	err := viper.UnmarshalKey("etcd", &config)
	if err != nil {
		panic(err)
	}
	etcdClient, err := clientv3.New(config)
	if err != nil {
		panic(err)
	}
	return etcdClient
}
