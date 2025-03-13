package nanoid

import (
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	for i := 0; i < 100; i++ {
		nid := New()
		logger.Infof("nanoid:%s", nid)
		assert.NotEmpty(t, nid)
	}
}
