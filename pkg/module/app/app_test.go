package app_test

import (
	"github.com/davycun/eta/pkg/common/http_tes"
	_ "github.com/davycun/eta/pkg/common/http_tes"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/module/app"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestName(t *testing.T) {

	rs, i := http_tes.Query[app.App](t, "/app/query", dto.RetrieveParam{AutoCount: true})
	assert.Greater(t, i, int64(0))
	assert.Greater(t, len(rs), 0)

}
