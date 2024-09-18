package ioc

import (
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	grpc2 "webook/interactive/grpc"
	"webook/pkg/grpcx"
)

func NewGrpcxServer(intrSvc *grpc2.InteractiveServiceServer) *grpcx.Server {
	type Config struct {
		Addr string `yaml:"addr"`
	}
	var config Config
	err := viper.UnmarshalKey("grpc.server", &config)
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()
	intrSvc.Register(server)
	return &grpcx.Server{
		Server: server,
		Addr:   config.Addr,
	}

}
