package sms

type Sender interface {
	SendCustom(task Task) ([]Target, error)
	// SendTemplate
	//如果是模版发送
	// 1）如果短信平台不支持模版发送，那么需要自己填写Content，Content示例：您的登录验证码为:{code} 其中code的值从TemplateParam获取
	// 2）如果短信平台支持模版发送，那么就不需要填写Content（就算填写了也不会采用），需要填写短信平台的TemplateCode和TemplateParam
	SendTemplate(target Target) (Target, error)
}
