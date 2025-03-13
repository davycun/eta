package oauth2

type LoginParam struct {
	LoginType string `json:"login_type"`
	Param     any    `json:"args"`
}

type LoginByUsernameParam struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	AppId    string `json:"app_id"`
}

type LoginByCodeParam struct {
	LoginType string `json:"login_type" binding:"required"`
	Code      string `json:"code" binding:"required"`
	AppId     string `json:"app_id"`
	Phone     string `json:"phone" binding:"mobile"`
}

type AccessKeyParam struct {
	AccessKey string `json:"access_key" form:"access_key" binding:"required"` //UserKey.AccessKey
	Algo      string `json:"algo" form:"algo" binding:"required"`             //加密的算法，默认是hmac_sha256
	Nonce     string `json:"nonce" form:"nonce" binding:"required"`           //随机数
	Ts        int64  `json:"ts" form:"ts" binding:"required"`                 // unix时间戳, 秒，UTC
	Sign      string `json:"sign" form:"sign" binding:"required"`             //根据accessKey对应的AccessSecure 对Ts+Nonce进行Algo签名
}

type LoginResult struct {
	Authorization string `json:"authorization"` //token
	ExpiresIn     int64  `json:"expires_in"`    //多少秒之后过期
	Data          any    `json:"data"`          //携带登录成功的用户和APP信息
}
