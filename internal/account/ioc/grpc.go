package ioc

import (
	grpc2 "github.com/DaHuangQwQ/webook/account/grpc"
	"github.com/DaHuangQwQ/webook/pkg/grpcx"
	"github.com/DaHuangQwQ/webook/pkg/logger"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
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
