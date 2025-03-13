package storage_test

import (
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/storage"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAppStorageFolder(t *testing.T) {
	str := storage.AppStorageFolder("appId")
	logger.Infof(str)
	assert.Equal(t, str, "/data/eta_storage/appId")
}
