package aliyun

import (
	"context"
	"encoding/json"
	"fmt"
	aliyunSms "github.com/alibabacloud-go/dysmsapi-20170525/v4/client"
	"webook/internal/service/sms"
)

type Service struct {
	appId    *string
	signName *string
	client   *aliyunSms.Client
}

func NewService(client *aliyunSms.Client, appId string, signName string) *Service {
	return &Service{
		appId:    &appId,
		signName: &signName,
		client:   client,
	}
}

func (s *Service) Send(ctx context.Context, tpl string, args []sms.NamedArg, numbers ...string) error {
	testStr := numbers[0]

	jsonMap := make(map[string]string)

	for _, arg := range args {
		jsonMap[arg.Name] = arg.Value
	}
	jsonData, err := json.Marshal(jsonMap)
	if err != nil {
		return fmt.Errorf("json marshal err: %v", err)
	}
	jsonStr := string(jsonData)
	println(jsonStr)

	sendSmsReq := &aliyunSms.SendSmsRequest{
		PhoneNumbers:  &testStr,
		SignName:      s.signName,
		TemplateCode:  &tpl,
		TemplateParam: &jsonStr,
	}
	res, err := s.client.SendSms(sendSmsReq)
	if err != nil {
		return err
	}
	println(res.String())
	return nil
}
