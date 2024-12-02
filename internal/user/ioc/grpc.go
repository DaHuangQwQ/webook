package ioc

import (
	"github.com/DaHuangQwQ/webook/pkg/grpcx"
	"github.com/DaHuangQwQ/webook/pkg/logger"
	grpc2 "github.com/DaHuangQwQ/webook/user/grpc"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

func NewGrpcxServer(userSvc *grpc2.UserServiceServer, l logger.LoggerV1, etcdClient *clientv3.Client) *grpcx.Server {
	type Config struct {
		Port int `yaml:"port"`
	}
	var config Config
	err := viper.UnmarshalKey("grpc.server", &config)
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()
	userSvc.Register(server)
	return &grpcx.Server{
		Server:     server,
		Port:       config.Port,
		Name:       "user",
		L:          l,
		EtcdClient: etcdClient,
	}

}
