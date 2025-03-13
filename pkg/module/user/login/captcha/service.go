package captcha

import (
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/module/setting"
	"github.com/davycun/eta/pkg/module/sms"
	"github.com/davycun/eta/pkg/module/sms/sms_sender"
	"github.com/duke-git/lancet/v2/random"
	"github.com/mojocn/base64Captcha"
	"time"
)

var (
	store = base64Captcha.NewMemoryStore(10240, 10*time.Second)
)

const (
	SmsForRegister  = "Register"  // 注册
	SmsForLogin     = "Login"     // 登录
	SmsForTwoFactor = "TwoFactor" // 双因子
	SmsForModify    = "Modify"    // 修改手机号
)

type Service struct {
	c *ctx.Context
}

func newService(c *ctx.Context) *Service {
	return &Service{c: c}
}

func (s *Service) GenerateImage(param *ImageCaptchaParam, result *dto.Result) error {

	var (
		cpt = Captcha{}
	)

	driverDigit := &base64Captcha.DriverDigit{
		Length:   4,   //数字个数
		Height:   50,  //高度
		Width:    120, //宽度
		MaxSkew:  0.7,
		DotCount: 80,
	}
	if param.Length > 0 {
		driverDigit.Length = param.Length
	}
	if param.Height > 0 {
		driverDigit.Height = param.Height
	}
	if param.Width > 0 {
		driverDigit.Width = param.Width
	}
	if param.MaxSkew > 0 {
		driverDigit.MaxSkew = param.MaxSkew
	}
	if param.DotCount > 0 {
		driverDigit.DotCount = param.DotCount
	}

	//生成验证码
	cp := base64Captcha.NewCaptcha(driverDigit, store)
	_, b64s, answer, err := cp.Generate()
	if err != nil {
		return err
	}
	result.Data = b64s
	cpt.Code = answer
	return StoreCaptcha(cpt.Code, cpt, time.Second*60)

}

func (s *Service) SendSmsCode(param *SmsCaptchaParam, result *dto.Result) error {

	var (
		cpt = Captcha{Phone: param.Phone, Code: random.RandNumeral(6)}
	)

	//发送验证码，不需要登录，所以
	cfg := setting.GetCommonConfig(global.GetLocalGorm())
	if cfg.SmsNeedImageCaptcha {
		if param.ImageCode == "" || param.Phone == "" {
			return errs.NewClientError("发送短信验证码时，图片验证码和手机号不能为空!")
		}
		if !Verify(Captcha{Code: param.ImageCode}) {
			return errs.NewClientError("图片验证码错误")
		}
	}

	target := sms.Target{}
	target.Mobile = param.Phone
	target.TemplateKey = "login_template"
	err := sms_sender.SendVerifyCode(s.c.GetAppGorm(), target, cpt.Code)
	if err != nil {
		return err
	}
	return StoreCaptcha(cpt.Code, cpt, time.Second*60)
}
