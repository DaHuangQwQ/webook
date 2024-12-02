package ioc

import (
	accountv1 "github.com/DaHuangQwQ/webook/api/proto/gen/account/v1"
	_ "github.com/DaHuangQwQ/webook/pkg/grpcx/balancer/wrr"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitAccountGrpcClient(etcdClient *clientv3.Client) accountv1.AccountServiceClient {
	type Config struct {
		Addr   string `yaml:"addr"`
		Secure bool   `yaml:"secret"`
	}
	var config Config
	err := viper.UnmarshalKey("grpc.client.account", &config)
	if err != nil {
		panic(err)
	}
	bd, err := resolver.NewBuilder(etcdClient)
	if err != nil {
		panic(err)
	}
	opts := []grpc.DialOption{
		grpc.WithResolvers(bd),
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [ { "custom_wrr": {} } ]}`),
	}
	if !config.Secure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	cc, err := grpc.Dial(config.Addr, opts...)
	if err != nil {
		panic(err)
	}
	return accountv1.NewAccountServiceClient(cc)
}
