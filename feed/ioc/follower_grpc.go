package ioc

import (
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	followv1 "webook/api/proto/gen/follow/v1"
)

func InitFollowClient() followv1.FollowServiceClient {
	type config struct {
		Target string `yaml:"target"`
	}
	var cfg config
	err := viper.UnmarshalKey("grpc.client.sms", &cfg)
	if err != nil {
		panic(err)
	}
	conn, err := grpc.Dial(
		cfg.Target,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	client := followv1.NewFollowServiceClient(conn)
	return client
}
