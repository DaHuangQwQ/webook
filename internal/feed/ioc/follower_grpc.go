package ioc

import (
	followv1 "github.com/DaHuangQwQ/webook/api/proto/gen/follow/v1"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitFollowClient() followv1.FollowServiceClient {
	type config struct {
		Target string `yaml:"target"`
	}
	var cfg config
	err := viper.UnmarshalKey("grpc.client.follow", &cfg)
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
