package constants

const (
	LoginTypeAccount       = "eta_account" //每种类型的登录方式，生成的token都会以下面定义的常量开头
	LoginTypeDingService   = "eta_ding_service"
	LoginTypeDingQrcode    = "eta_ding_qrcode"
	LoginTypeZzdService    = "eta_zzd_service"
	LoginTypeZzdQrcode     = "eta_zzd_qrcode"
	LoginTypeWechatService = "eta_wechat_service"
	LoginTypeWechatQrcode  = "eta_wechat_qrcode"
	LoginTypeWeComService  = "eta_wecom_service"
	LoginTypeWeComQrcode   = "eta_wecom_qrcode"
	LoginTypeSmsService    = "eta_sms_service"
	LoginTypeAccessToken   = "eta_access_token"
	LoginTypeFixToken      = "eta_fixed_token" ///这种方式不需要登录，但是fix token 需要以这个常量作为前缀
)
