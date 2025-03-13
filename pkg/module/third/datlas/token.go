package datlas

import (
	"time"
)

const (
	TokenExpireIn int = 2 * 60 * 60 // 秒
)

type Token struct {
	Auth      string    `json:"auth"`
	MdtUser   string    `json:"mdt_user"`
	ExpiresIn int       `json:"expires_in"`
	ExpiredAt time.Time `json:"expired_at"`
}

func (t *Token) IsExpired() bool {
	return t.Auth == "" || t.ExpiredAt.Second()-time.Now().Second() < 300
}
