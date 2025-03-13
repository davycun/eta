package constants

const (
	HeaderRespFetchCacheId string = "X-DELTA-Fetch-Cache-Id"
	HeaderRespOpFromEs     string = "X-DELTA-RA" // 来自es的操作

	HeaderCryptSymmetryAlgorithm   = "X-CRYPT-SYMMETRY-ALGO"   //对称加密算法
	HeaderCryptSymmetryKey         = "X-CRYPT-SYMMETRY-KEY"    //对称加密算法的key，通过非对称加密进行了加密，并且进行了base64编码
	HeaderCryptAsymmetricAlgorithm = "X-CRYPT-ASYMMETRIC-ALGO" //非对称加密算法

	HeaderAuthorization = "Authorization"
	HeaderRequestId     = "X-Request-ID"
	HeaderResponseFrom  = "X-Response-From" // db、es

	HeaderUserId = "X-User-Id"
	HeaderAppId  = "X-App-Id"

	HeaderOptClientTrigger = "X-Client-Trigger"
	HeaderOptClientType    = "X-Client-Type"
	HeaderOptType          = "X-Opt-Type"
	HeaderOptContent       = "X-Opt-Content"
	HeaderOptTarget        = "X-Opt-Target"
)
