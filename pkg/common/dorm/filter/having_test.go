package filter_test

import (
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/filter"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResolveHaving(t *testing.T) {
	js := `[
    {
      "logical_operator": "and",
      "having": [
        {
          "logical_operator": "and",
          "agg_func":"count",
          "column": "shop_type",
          "operator": ">",
          "value": 21
        },
        {
          "logical_operator": "and",
          "agg_func":"count",
          "column": "*",
          "operator": ">",
          "value": 21
        },
        {
          "logical_operator": "or",
		  "agg_func":"max",
          "column": "shop_type",
          "operator": "=",
          "value": "abcd"
        }
      ]
    }
  ]`
	var fs []filter.Having
	err := jsoniter.Unmarshal([]byte(js), &fs)
	wh := filter.ResolveHaving(dorm.DaMeng, fs...)
	wh2 := filter.ResolveHavingTable(dorm.DaMeng, "my_table", fs...)
	assert.Equal(t, ` ( ( count("shop_type")  > 21 and  count(*)  > 21 or  max("shop_type")  = 'abcd') ) `, wh)
	assert.Equal(t, ` ( ( count("my_table"."shop_type")  > 21 and  count(*)  > 21 or  max("my_table"."shop_type")  = 'abcd') ) `, wh2)
	assert.Nil(t, err)

}
