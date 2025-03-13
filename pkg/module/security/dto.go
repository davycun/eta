package security

type SignParam struct {
	Algo string `json:"algo" binding:"required"`
	Salt string `json:"salt" binding:"required"`
	Msg  string `json:"msg" binding:"required"`
}

type EncryptParam struct {
	Algo      string   `json:"algo" binding:"required"`
	SecretKey []string `json:"secret_key" binding:"required"`
	Msg       string   `json:"msg" binding:"required"`
}

type Result struct {
	Data string `json:"data"`
}

type KeyResult struct {
	Algo      string `json:"algo"`
	PublicKey string `json:"public_key"`
}
