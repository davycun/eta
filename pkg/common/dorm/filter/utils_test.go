package filter_test

import (
	"github.com/davycun/eta/pkg/common/dorm/expr"
	"github.com/davycun/eta/pkg/common/dorm/filter"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConvertFilter(t *testing.T) {

	fs := []filter.Filter{
		{
			Column:   "name",
			Operator: filter.Like,
			Value:    "%大%",
			Filters: []filter.Filter{
				{
					Column:   "age",
					Operator: filter.Eq,
					Value:    23,
				},
			},
		},
		{
			Expr: expr.Expression{
				Expr: "? & ?",
				Vars: []expr.ExpVar{
					{
						Type:  expr.VarTypeColumn,
						Value: "usage",
					},
					{
						Type:  expr.VarTypeValue,
						Value: 1,
					},
				},
			},
			Operator: filter.Eq,
			Value:    1,
		},
		{
			Column:   "other",
			Operator: filter.Eq,
			Value:    "test",
		},
	}

	fs2 := filter.ConvertFilterColumnName("pep", fs)
	assert.Equal(t, "pep", fs2[0].Column)
	assert.Equal(t, "pep", fs2[0].Filters[0].Column)
	assert.Equal(t, "pep", fs2[1].Expr.Vars[0].Value)
	assert.Equal(t, "pep", fs2[2].Column)

	fs3 := filter.ConvertFilterColumnName("pep", fs, "age", "usage")
	assert.Equal(t, 2, len(fs3))
	assert.Equal(t, "", fs3[0].Column)
	assert.Nil(t, fs3[0].Value)
	assert.Equal(t, "pep", fs3[0].Filters[0].Column)
	assert.Equal(t, "pep", fs3[1].Expr.Vars[0].Value)

	fs4 := filter.AddFilterColumnPrefix("pep.", fs)
	assert.Equal(t, 3, len(fs4))
	assert.Equal(t, "pep.name", fs4[0].Column)
	assert.Equal(t, "pep.age", fs4[0].Filters[0].Column)
	assert.Equal(t, "pep.usage", fs4[1].Expr.Vars[0].Value)
	assert.Equal(t, "pep.other", fs4[2].Column)

	fs5 := filter.AddFilterColumnPrefix("pep.", fs, "age", "usage")
	assert.Equal(t, 2, len(fs5))
	assert.Equal(t, "", fs5[0].Column)
	assert.Equal(t, "pep.age", fs5[0].Filters[0].Column)
	assert.Equal(t, "pep.usage", fs5[1].Expr.Vars[0].Value)

	//幂等
	fs3 = filter.ConvertFilterColumnName("pep", fs, "age", "usage")
	assert.Equal(t, 2, len(fs3))
	assert.Equal(t, "", fs3[0].Column)
	assert.Equal(t, "pep", fs3[0].Filters[0].Column)
	assert.Equal(t, "pep", fs3[1].Expr.Vars[0].Value)

	fs6 := filter.ConvertFilterColumnName("pep", fs, "other")
	assert.Equal(t, 1, len(fs6))
	assert.Equal(t, "pep", fs6[0].Column)

	fs7 := filter.AddFilterColumnPrefix("pep.", fs, "other")
	assert.Equal(t, 1, len(fs7))
	assert.Equal(t, "pep.other", fs7[0].Column)

}
