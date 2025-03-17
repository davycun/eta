package security_test

import (
	"github.com/davycun/eta/pkg/common/crypt"
	_ "github.com/davycun/eta/pkg/common/http_tes"
	"github.com/davycun/eta/pkg/module/security"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestGetKey(t *testing.T) {
	sm2PubKey := make([]string, 0)
	sm2PriKey := make([]string, 0)
	rsaPubKey := make([]string, 0)
	rsaPriKey := make([]string, 0)
	size := 20
	w := &sync.WaitGroup{}
	for i := 0; i < size; i++ {
		w.Add(1)
		go func() {
			defer w.Done()
			rsaPubKey = append(rsaPubKey, security.GetPublicKey(crypt.AlgoASymRsaPKCS1v15))
			rsaPriKey = append(rsaPriKey, security.GetPrivateKey(crypt.AlgoASymRsaPKCS1v15))
			sm2PubKey = append(sm2PubKey, security.GetPublicKey(crypt.AlgoASymSm2Pkcs8C132))
			sm2PriKey = append(sm2PriKey, security.GetPrivateKey(crypt.AlgoASymSm2Pkcs8C132))
		}()
	}
	w.Wait()
	for i := 0; i < size; i++ {
		if i > 0 {
			assert.Equal(t, sm2PriKey[i], sm2PriKey[i-1])
			assert.Equal(t, sm2PubKey[i], sm2PubKey[i-1])
			assert.Equal(t, rsaPriKey[i], rsaPriKey[i-1])
			assert.Equal(t, rsaPriKey[i], rsaPriKey[i-1])
		}
	}
}
