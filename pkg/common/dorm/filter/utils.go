package filter

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/dorm/expr"
	"github.com/davycun/eta/pkg/common/utils"
	"slices"
)

// ConvertFilterColumnName
// 把指定列过滤名字换成目标列
func ConvertFilterColumnName(targetCol string, flt []Filter, originCol ...string) []Filter {
	rs := make([]Filter, 0, len(flt))
	for _, v := range flt {
		target := false
		if v.Expr.Expr != "" {
			v.Expr.Vars = slices.Clone(v.Expr.Vars)
			for x, y := range v.Expr.Vars {
				if y.Type == expr.VarTypeColumn {
					col := fmt.Sprintf("%s", y.Value)
					target = col != "" && len(originCol) < 1 || utils.ContainAny(originCol, col)
					if target {
						v.Expr.Vars[x].Value = targetCol
					}
				}
			}
			if target {
				rs = append(rs, v)
			}
			continue
		}

		target = v.Column != "" && (len(originCol) < 1 || utils.ContainAny(originCol, v.Column))
		if target {
			v.Column = targetCol
		}

		if len(v.Filters) > 0 {
			v.Filters = ConvertFilterColumnName(targetCol, v.Filters, originCol...)
			//如果当前的filter不是目标，但是子Filter有目标，那么去掉当前的filter内容
			if !target && len(v.Filters) > 0 {
				v.Column = ""
				v.Operator = ""
				v.Value = nil
			}
			//如果有子filter符合目标，那么也需要添加当前filter
			if len(v.Filters) > 0 {
				target = true
			}
		}

		if target {
			rs = append(rs, v)
		}
	}
	return rs
}
func AddFilterColumnPrefix(prefix string, flt []Filter, originCol ...string) []Filter {
	rs := make([]Filter, 0, len(flt))
	for _, v := range flt {
		target := false
		if v.Expr.Expr != "" {
			v.Expr.Vars = slices.Clone(v.Expr.Vars)
			for x, y := range v.Expr.Vars {
				if y.Type == expr.VarTypeColumn {
					col := fmt.Sprintf("%s", y.Value)
					target = col != "" && len(originCol) < 1 || utils.ContainAny(originCol, col)
					if target {
						v.Expr.Vars[x].Value = prefix + col
					}
				}
			}

			if target {
				rs = append(rs, v)
			}
			continue
		}

		target = v.Column != "" && (len(originCol) < 1 || utils.ContainAny(originCol, v.Column))
		if target {
			v.Column = prefix + v.Column
		}

		if len(v.Filters) > 0 {
			v.Filters = AddFilterColumnPrefix(prefix, v.Filters, originCol...)
			//如果当前的filter不是目标，但是子Filter有目标，那么去掉当前的filter内容
			if !target && len(v.Filters) > 0 {
				v.Column = ""
				v.Operator = ""
				v.Value = nil
			}
			//如果有子filter符合目标，那么也需要添加当前filter
			if len(v.Filters) > 0 {
				target = true
			}
		}
		if target {
			rs = append(rs, v)
		}
	}
	return rs
}
