package ioc

//func InitOssService() oss.Client {
//	return InitAliyunOss()
//}
//
//func InitAliyunOss() oss.Client {
//	type Config struct {
//		KEY    string `yaml:"KEY"`
//		SECRET string `yaml:"SECRET"`
//		ENDP   string `yaml:"ENDP"`
//	}
//	var config Config
//	err := viper.UnmarshalKey("OSS", &config)
//	if err != nil {
//		panic(err)
//	}
//	client, err := aliyunOss.New(config.ENDP, config.KEY, config.SECRET)
//	if err != nil {
//		panic(err)
//	}
//	return oss.NewAliyunOss(client)
//}
