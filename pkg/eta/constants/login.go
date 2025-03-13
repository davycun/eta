package constants

const (
	LoginTypeAccount       = "account" //每种类型的登录方式，生成的token都会以下面定义的常量开头
	LoginTypeDingService   = "ding_service"
	LoginTypeDingQrcode    = "ding_qrcode"
	LoginTypeZzdService    = "zzd_service"
	LoginTypeZzdQrcode     = "zzd_qrcode"
	LoginTypeWechatService = "wechat_service"
	LoginTypeWechatQrcode  = "wechat_qrcode"
	LoginTypeWeComService  = "wecom_service"
	LoginTypeWeComQrcode   = "wecom_qrcode"
	LoginTypeSmsService    = "sms_service"
	LoginTypeAccessToken   = "access_token"
	LoginTypeFixToken      = "fixed_token" ///这种方式不需要登录，但是fix token 需要以这个常量作为前缀
)
