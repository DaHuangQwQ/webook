package ioc

import (
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"webook/pkg/grpcx"
	"webook/pkg/logger"
	grpc2 "webook/tag/grpc"
)

func NewGrpcxServer(tagSvc *grpc2.TagServiceServer, l logger.LoggerV1, etcdClient *clientv3.Client) *grpcx.Server {
	type Config struct {
		Port int `yaml:"port"`
	}
	var config Config
	err := viper.UnmarshalKey("grpc.server", &config)
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()
	tagSvc.Register(server)
	return &grpcx.Server{
		Server:     server,
		Port:       config.Port,
		Name:       "interactive",
		L:          l,
		EtcdClient: etcdClient,
	}

}
