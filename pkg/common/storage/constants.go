package storage

const (
	StorePlatformMinio     = "minio"
	StorePlatformAliYunOss = "aliyun_oss"
	StorePlatformLocal     = "local"

	ParamKeyAlgorithm  = "X-Eta-Algorithm"
	ParamKeyCredential = "X-Eta-Credential"
	ParamKeyDate       = "X-Eta-Date"
	ParamKeyExpires    = "X-Eta-Expires"
	ParamKeySignature  = "X-Eta-Signature"
	DefaultAlgorithm   = "ETA-HMAC-SHA256"
)
