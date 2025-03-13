package template_test

import (
	"github.com/davycun/eta/pkg/module/data/template"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTable(t *testing.T) {

	tmpl := template.Template{}
	tmpl.Code = "test"
	tb := tmpl.GetTable()
	assert.Equal(t, tb.TableName, tmpl.Table.TableName)
}
