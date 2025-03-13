package weixin

type WxError struct {
	ErrCode int    `json:"errode"`
	ErrMsg  string `json:"errmsg"`
}
