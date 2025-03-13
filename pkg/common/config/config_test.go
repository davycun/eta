package config_test

import (
	"github.com/davycun/eta/pkg/common/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

var content = `
database:
  host: 127.0.0.1
  port: 15237
  user: citizen
  password: abc
  dbname: dameng
  schema: CITIZEN
  type: dm
  log_level: 4
  slow_threshold: 200
server:
  port: 8080
variables:
  MY_HOME: "/etc/eta"
  CONFIG: /etc/eta/defaultConf
`

func TestConfig(t *testing.T) {

	var c config.Configuration
	err := config.LoadFromContent(content, &c)

	assert.Nil(t, err)
	assert.Equal(t, 200, c.Database.SlowThreshold)
	assert.Equal(t, 8080, c.Server.Port)
	assert.Equal(t, 2, len(c.Variables))
}
