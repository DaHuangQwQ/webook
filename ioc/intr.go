package ioc

import (
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	interactivev1 "webook/api/proto/gen/interactive/v1"
	"webook/internal/client"
)

func InitIntrGRPCClient() interactivev1.InteractiveServiceClient {
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
