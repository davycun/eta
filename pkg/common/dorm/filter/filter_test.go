package filter_test

import (
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/filter"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResolveWhere(t *testing.T) {
	js := `[
    {
      "logical_operator": "and",
      "filters": [
        {
          "logical_operator": "and",
          "column": "shop_type",
          "operator": "like",
          "value": "%经营性住宿场所%"
        },
        {
          "logical_operator": "or",
          "column": "shop_type",
          "operator": "like",
          "value": "%医疗服务场所%"
        }
      ]
    }
  ]`
	var fs []filter.Filter
	err := jsoniter.Unmarshal([]byte(js), &fs)
	wh := filter.ResolveWhere(fs, dorm.DaMeng)
	assert.Equal(t, ` ( ( "shop_type"  like '%经营性住宿场所%' or  "shop_type"  like '%医疗服务场所%') ) `, wh)
	assert.Nil(t, err)

}

//
//import (
//	"encoding/json"
//	"github.com/stretchr/testify/assert"
//	"github.com/davycun/eta/pkg/common/global"
//	"github.com/davycun/eta/pkg/common/dorm/filter"
//	"github.com/davycun/eta/pkg/core/dto"
//	"strings"
//	"testing"
//)
//
//var (
//	AND = "and"
//	OR  = "or"
//)
//
//// (k1='v1' or k2='v2') or (k3=true and k4=5)
//func TestFilter(t *testing.T) {
//
//	c1 := filter.Filter{
//		LogicalOperator: AND,
//		Filters:         make([]filter.Filter, 2),
//	}
//
//	c11 := filter.Filter{
//		Operator:        filter.Eq,
//		LogicalOperator: AND,
//		Column:          "k1",
//		Value:           "v1",
//	}
//	c12 := filter.Filter{
//		Operator:        filter.Eq,
//		LogicalOperator: OR,
//		Column:          "k2",
//		Value:           "v2",
//	}
//	c1.Filters[0] = c11
//	c1.Filters[1] = c12
//
//	c2 := filter.Filter{
//		LogicalOperator: OR,
//		Filters:         make([]filter.Filter, 2),
//	}
//
//	c21 := filter.Filter{
//		Operator:        filter.Eq,
//		LogicalOperator: OR,
//		Column:          "k3",
//		Value:           true,
//	}
//	c22 := filter.Filter{
//		Operator:        filter.Eq,
//		LogicalOperator: AND,
//		Column:          "k4",
//		Value:           5,
//	}
//	c2.Filters[0] = c21
//	c2.Filters[1] = c22
//
//	c3 := filter.Filter{
//		Operator:        filter.IN,
//		Column:          "k0",
//		LogicalOperator: AND,
//		Value:           []string{"a", "b", "c"},
//	}
//
//	cds := make([]filter.Filter, 3)
//	cds[0] = c3
//	cds[1] = c1
//	cds[2] = c2
//
//	r := getJson(t, cds...)
//
//	rs := filter.ResolveWhere(r.Filters, global.ResolveDbType("dm"))
//	rs = strings.Trim(rs, " ")
//	assert.Equal(t, `( k0 in ('a','b','c') and ( k1 = 'v1' or k2 = 'v2' ) or ( k3 = true and k4 = 5 ) )`, rs)
//}
//
//func TestIn(t *testing.T) {
//
//	cds := make([]filter.Filter, 3)
//	rs := []string{`( k in (3,4,5,6,7,8) )`, `( k not in (3.5,4,1235.9875,6.7) )`, `( k in ('dvay','buhao','什么') )`}
//
//	cds[0] = filter.Filter{
//		Operator:        filter.IN,
//		LogicalOperator: AND,
//		Column:          "k",
//		Value:           []int64{3, 4, 5, 6, 7, 8},
//	}
//	cds[1] = filter.Filter{
//		Operator:        filter.NotIn,
//		LogicalOperator: AND,
//		Column:          "k",
//		Value:           []float64{3.5, 4, 1235.9875, 6.7},
//	}
//	cds[2] = filter.Filter{
//		Operator:        filter.IN,
//		LogicalOperator: AND,
//		Column:          "k",
//		Value:           []string{"dvay", "buhao", "什么"},
//	}
//	for i, v := range cds {
//		r := getJson(t, v)
//		w := filter.ResolveWhere(r.Filters, global.ResolveDbType("dm"))
//		assert.Equal(t, rs[i], strings.Trim(w, " "))
//	}
//
//}
//
//func getJson(t *testing.T, cds ...filter.Filter) dto.RetrieveParam {
//
//	c := make([]filter.Filter, 0, len(cds))
//	c = append(c, cds...)
//	p := dto.RetrieveParam{}
//	p.Filters = c
//
//	marshal, err := json.Marshal(p)
//	assert.Nil(t, err)
//
//	r := dto.RetrieveParam{}
//	err = json.Unmarshal(marshal, &r)
//	assert.Nil(t, err)
//	return r
//}
