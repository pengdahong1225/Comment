package captcha

import (
	"Comment/module/settings"
	"Comment/module/utils"
	"encoding/json"
	"errors"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/sirupsen/logrus"
)

func SendSmsCode(mobile string) (string, error) {
	// 生成随机数
	c := utils.GenerateSmsCode(6)
	param := map[string]string{
		"code": c,
	}
	data, _ := json.Marshal(param)

	// 调用第三方服务发送
	if err := send(data, mobile); err != nil {
		return "", err
	}

	return c, nil
}

func send(param []byte, phone string) error {
	config := &openapi.Config{
		// 您的AccessKey ID
		AccessKeyId: &settings.Instance().SmsConfig.AccessKeyId,
		// 您的AccessKey Secret
		AccessKeySecret: &settings.Instance().SmsConfig.AccessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String(settings.Instance().SmsConfig.Endpoint)
	client, _ := dysmsapi.NewClient(config)

	request := &dysmsapi.SendSmsRequest{}
	request.SetSignName(settings.Instance().SmsConfig.SignName)
	request.SetTemplateCode(settings.Instance().SmsConfig.TemplateCode)
	request.SetPhoneNumbers(phone)
	request.SetTemplateParam(string(param))

	response, err := client.SendSms(request)
	if err != nil {
		return err
	}
	if *response.Body.Code != "OK" {
		return errors.New(tea.StringValue(response.Body.Message))
	}
	logrus.Debugln(tea.StringValue(response.Body.RequestId))

	return nil
}
