package ioc

import (
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tencentSMS "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"webook/internal/service/sms"
	"webook/internal/service/sms/localsms"
	"webook/internal/service/sms/tencent"
)

func InitSMSService() sms.Service {
	return localsms.NewService()
	// 如果有需要，就可以用这个
	//return initTencentSMSService()
}

func initTencentSMSService() sms.Service {
	//secretId, ok := os.LookupEnv("SMS_SECRET_ID")
	//if !ok {
	//	panic("找不到腾讯 SMS 的 secret id")
	//}
	//secretKey, ok := os.LookupEnv("SMS_SECRET_KEY")
	//if !ok {
	//	panic("找不到腾讯 SMS 的 secret key")
	//}
	secretId := "1234"
	secretKey := "1234"
	c, err := tencentSMS.NewClient(
		common.NewCredential(secretId, secretKey),
		"ap-nanjing",
		profile.NewClientProfile(),
	)
	if err != nil {
		panic(err)
	}
	return tencent.NewService(c, "1400842696", "妙影科技")
}
