package ioc

import (
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	interactivev1 "webook/api/proto/gen/interactive/v1"
	"webook/internal/client"
	_ "webook/pkg/grpcx/balancer/wrr"
)

func InitEtcd() *clientv3.Client {
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

func InitIntrGRPCClient(etcdClient *clientv3.Client) interactivev1.InteractiveServiceClient {
	type Config struct {
		Addr   string `yaml:"addr"`
		Secure bool   `yaml:"secret"`
	}
	var config Config
	err := viper.UnmarshalKey("grpc.client.intr", &config)
	if err != nil {
		panic(err)
	}

	bd, err := resolver.NewBuilder(etcdClient)
	if err != nil {
		panic(err)
	}
	opts := []grpc.DialOption{
		grpc.WithResolvers(bd),
		// 负载均衡器(轮询)
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [ { "custom_wrr": {} } ]}`),
	}
	if !config.Secure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	cc, err := grpc.Dial(config.Addr, opts...)
	if err != nil {
		panic(err)
	}
	remote := interactivev1.NewInteractiveServiceClient(cc)
	return client.NewInteractiveRemoteClient(remote)
}

// InitIntrGRPCClientV1 流量控制
func InitIntrGRPCClientV1() interactivev1.InteractiveServiceClient {
	type Config struct {
		Addr   string `yaml:"addr"`
		Secure bool   `yaml:"secret"`
	}
	var config Config
	err := viper.UnmarshalKey("grpc.client.intr", &config)
	if err != nil {
		panic(err)
	}
	var opts []grpc.DialOption
	if config.Secure {
		// tls up
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	cc, err := grpc.NewClient(config.Addr, opts...)
	//local := client.NewLocalInteractiveServiceAdapter()
	remote := interactivev1.NewInteractiveServiceClient(cc)
	intrClient := client.NewInteractiveClient(remote)
	intrClient.UpdateThreshold(100)
	return intrClient
}
