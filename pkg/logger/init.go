package logger

import (
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func InitViper() {
	file := pflag.String("config", "config/dev.yaml", "指定文件路径")
	pflag.Parse()
	viper.SetConfigFile(*file)
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
