package ioc

import (
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	smsv1 "webook/api/proto/gen/sms/v1"
	_ "webook/pkg/grpcx/balancer/wrr"
)

func InitSmsGRPCClient(etcdClient *clientv3.Client) smsv1.SmsServiceClient {
	type Config struct {
		Addr   string `yaml:"addr"`
		Secure bool   `yaml:"secret"`
	}
	var config Config
	err := viper.UnmarshalKey("grpc.client.sms", &config)
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
	return smsv1.NewSmsServiceClient(cc)
}
