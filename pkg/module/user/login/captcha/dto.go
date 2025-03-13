package captcha

type ImageCaptchaParam struct {
	Length   int     `json:"length" form:"length"`
	Height   int     `json:"height" form:"height"`
	Width    int     `json:"width" form:"width"`
	MaxSkew  float64 `json:"max_skew" form:"max_skew"`
	DotCount int     `json:"dot_count" form:"dot_count"`
}

type SmsCaptchaParam struct {
	SmsType   string `json:"sms_type,omitempty" binding:"required,oneof=Register Login TwoFactor Modify ''"`
	Phone     string `json:"phone,omitempty" binding:"required,mobile"` // 手机号
	ImageCode string `json:"image_code,omitempty" binding:"required"`   // 图形验证码
}

// Captcha
// 存储生成短信二维码或者图片验证码后对应的数据
type Captcha struct {
	Code  string `json:"code" redis:"code"`   // 短信验证码或者图片验证码
	Phone string `json:"phone" redis:"phone"` // 手机号
}

func (c Captcha) Verify(target Captcha) bool {

	if c.Code != target.Code {
		return false
	}
	if c.Phone != target.Phone {
		return false
	}
	return true
}
