package ioc

import (
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	grpc2 "webook/article/grpc"
	"webook/pkg/grpcx"
	"webook/pkg/logger"
)

func NewGrpcxServer(articleSvc *grpc2.ArticleServiceServer, l logger.LoggerV1, etcdClient *clientv3.Client) *grpcx.Server {
	type Config struct {
		Port int `yaml:"port"`
	}
	var config Config
	err := viper.UnmarshalKey("grpc.server", &config)
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()
	articleSvc.Register(server)
	return &grpcx.Server{
		Server:     server,
		Port:       config.Port,
		Name:       "article",
		L:          l,
		EtcdClient: etcdClient,
	}

}
