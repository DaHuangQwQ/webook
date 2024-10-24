package logger

import (
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

//func InitViper() {
//	file := pflag.String("config", "config/dev.yaml", "指定文件路径")
//	pflag.Parse()
//	viper.SetConfigFile(*file)
//	err := viper.ReadInConfig()
//	if err != nil {
//		panic(fmt.Errorf("Fatal error config file: %s \n", err))
//	}
//}

func InitViper() {
	err := viper.AddRemoteProvider("etcd3", "http://etcd:2379", "/config/config.yaml")
	if err != nil {
		panic(err)
	}
	viper.SetConfigType("yaml")
	err = viper.ReadRemoteConfig()
	if err != nil {
		panic(err)
	}
}
