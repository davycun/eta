package faker_test

import (
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/mock/faker"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJob(t *testing.T) {
	j := faker.Job()
	logger.Infof("job: %s", j)
	assert.NotNil(t, j)
}
