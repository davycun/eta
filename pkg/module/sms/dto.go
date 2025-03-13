package sms

import "github.com/davycun/eta/pkg/common/dorm/ctype"

// VerifyCodeParam
// 如果是非模版发送验证的就需要传入Content这个参数，示例："张先生在申请数据导出，验证码为：{code}"
// 如果采用短信模版方式发送验证码，那么就传入TemplateCode和TemplateParam
// 如果发送短信验证码需要通过图片验证码来验证，那么就传入ImageCode（暂未实现）
type VerifyCodeParam struct {
	Phone         string    `json:"phone,omitempty" binding:"required"`
	ImageCode     string    `json:"image_code,omitempty"`
	Content       string    `json:"content,omitempty" binding:"required"`
	SignName      string    `json:"sign_name,omitempty"`
	TemplateCode  string    `json:"template_code,omitempty"`
	TemplateParam ctype.Map `json:"template_param,omitempty"`
	TemplateKey   string    `json:"template_key,omitempty"`
}
