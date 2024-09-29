package ioc

import (
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	grpc2 "webook/account/grpc"
	"webook/pkg/grpcx"
	"webook/pkg/logger"
)

func NewGrpcServer(etcdClient *clientv3.Client, l logger.LoggerV1, accountSvc *grpc2.AccountServiceServer) *grpcx.Server {
	type Config struct {
		port int `yaml:"port"`
	}
	var config Config
	err := viper.UnmarshalKey("grpc.server", &config)
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()
	accountSvc.Register(server)
	return &grpcx.Server{
		Server:     server,
		Port:       config.port,
		EtcdClient: etcdClient,
		Name:       "payment",
		L:          l,
	}
}
