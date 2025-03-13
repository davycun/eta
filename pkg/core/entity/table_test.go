package entity

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestTable(t *testing.T) {
	tb1 := Table{}
	tb2 := &Table{
		TableName: "test",
	}
	tb1.Merge(tb2)
	assert.Equal(t, "test", tb1.TableName)

	assert.True(t, strings.EqualFold("btree", "BTREE"))
	assert.True(t, strings.EqualFold("BTREE", "BTREE"))
	assert.True(t, strings.EqualFold("hash", "HASH"))
	assert.True(t, strings.EqualFold("HASH", "HASH"))
}
