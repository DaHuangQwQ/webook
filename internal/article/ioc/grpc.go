package ioc

import (
	grpc2 "github.com/DaHuangQwQ/webook/article/grpc"
	"github.com/DaHuangQwQ/webook/pkg/grpcx"
	"github.com/DaHuangQwQ/webook/pkg/logger"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
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
