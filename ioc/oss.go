package ioc

import (
	aliyunOss "github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/spf13/viper"
	"webook/internal/repository/dao/oss"
)

func InitOssService() oss.Client {
	return InitAliyunOss()
}

func InitAliyunOss() oss.Client {
	type Config struct {
		KEY    string `yaml:"KEY"`
		SECRET string `yaml:"SECRET"`
		ENDP   string `yaml:"ENDP"`
	}
	var config Config
	err := viper.UnmarshalKey("OSS", &config)
	if err != nil {
		panic(err)
	}
	client, err := aliyunOss.New(config.ENDP, config.KEY, config.SECRET)
	if err != nil {
		panic(err)
	}
	return oss.NewAliyunOss(client)
}
