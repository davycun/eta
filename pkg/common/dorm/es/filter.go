package es

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/expr"
	"github.com/davycun/eta/pkg/common/dorm/filter"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/common/utils"
	"strings"
)

// ResolveEsQuery 解析es查询
// dbType : 数据库类型。请求参数的filter语法源自什么数据库
func ResolveEsQuery(dbType dorm.DbType, cds ...filter.Filter) (map[string]interface{}, error) {

	var (
		err         error
		filterQuery = make([]map[string]interface{}, 0, len(cds))
		shouldQuery = make([]map[string]interface{}, 0, len(cds))
	)
	for _, v := range cds {
		keywordCol := v.Column + ".keyword"
		curQuery := make(map[string]interface{})
		if v.Expr.Expr != "" {
			curQuery, err = explainExprEs(v)
			if err != nil {
				return nil, err
			}
		} else {
			switch strings.TrimSpace(strings.ToLower(v.Operator)) {
			case filter.Eq:
				curQuery = map[string]interface{}{
					"term": map[string]interface{}{
						keywordCol: v.Value,
					},
				}
			case filter.Neq:
				curQuery = map[string]interface{}{
					"bool": map[string]interface{}{
						"must_not": []map[string]interface{}{
							{
								"term": map[string]interface{}{
									keywordCol: v.Value,
								},
							},
						},
					},
				}
			case filter.GT:
				curQuery = map[string]interface{}{
					"range": map[string]interface{}{
						keywordCol: map[string]interface{}{
							"gt": v.Value,
						},
					},
				}
			case filter.LT:
				curQuery = map[string]interface{}{
					"range": map[string]interface{}{
						keywordCol: map[string]interface{}{
							"lt": v.Value,
						},
					},
				}
			case filter.GTE:
				curQuery = map[string]interface{}{
					"range": map[string]interface{}{
						keywordCol: map[string]interface{}{
							"gte": v.Value,
						},
					},
				}

			case filter.LTE:
				curQuery = map[string]interface{}{
					"range": map[string]interface{}{
						keywordCol: map[string]interface{}{
							"lte": v.Value,
						},
					},
				}
			case filter.IN:
				curQuery = map[string]interface{}{
					"terms": map[string]interface{}{
						keywordCol: v.Value,
					},
				}
			case filter.NotIn:
				curQuery = map[string]interface{}{
					"bool": map[string]interface{}{
						"must_not": []map[string]interface{}{
							{
								"terms": map[string]interface{}{
									keywordCol: v.Value,
								},
							},
						},
					},
				}
			case filter.Like:
				curQuery = map[string]interface{}{
					"match_phrase": map[string]interface{}{
						v.Column: strings.ReplaceAll(fmt.Sprintf("%s", v.Value), "%", ""),
					},
				}
			case filter.IS:
				curQuery = map[string]interface{}{
					"bool": map[string]interface{}{
						"must_not": []map[string]interface{}{
							{
								"exists": map[string]interface{}{
									"field": v.Column,
								},
							},
						},
					},
				}
			case filter.IsNot:
				curQuery = map[string]interface{}{
					"exists": map[string]interface{}{
						"field": v.Column,
					},
				}
			default:
				//should never happen,if it happens,please check the operator or return error
				if len(v.Filters) < 1 {
					continue
				}
			}
		}

		if len(curQuery) > 0 {
			switch v.LogicalOperator {
			case filter.Or:
				shouldQuery = append(shouldQuery, curQuery)
			default:
				filterQuery = append(filterQuery, curQuery)
			}
		}

		if len(v.Filters) > 0 {
			childQuery, err1 := ResolveEsQuery(dbType, v.Filters...)
			if err1 != nil {
				return nil, err1
			}
			if len(childQuery) > 0 {
				switch v.LogicalOperator {
				case filter.Or:
					shouldQuery = append(shouldQuery, childQuery)
				default:
					filterQuery = append(filterQuery, childQuery)
				}
			}
		}
	}

	//返回的是一个bool -> must, should, must_not
	boolQuery := make(map[string]interface{})

	if len(shouldQuery) > 0 {
		boolQuery["should"] = shouldQuery
	}
	if len(filterQuery) > 0 {
		boolQuery["filter"] = filterQuery
	}
	if len(boolQuery) < 1 {
		return map[string]interface{}{}, err
	}

	return map[string]interface{}{"bool": boolQuery}, err
}

// TODO 不支持数据库本身带函数的表达式
// 支持类似
// 达梦全文检索：contains(?,?)
// 数组索引：contains_arr(?,?)
// 与符号：? & ? 等等
func explainExprEs(flt filter.Filter) (map[string]interface{}, error) {

	var (
		ifCond       = make([]string, 0, 2)
		varScript    = strings.Builder{}
		returnScript = strings.Builder{}
		sourceScript = strings.Builder{}
		exp          = flt.Expr
		expScript    = strings.Builder{}
	)
	if strings.Count(exp.Expr, "?") != len(exp.Vars) {
		return nil, errs.NewClientError("表达式参数与实际参数数量不匹配")
	}

	//可能是达梦的权限检索contains，或者数据检索contains_arr
	if strings.HasPrefix(strings.TrimSpace(strings.ToLower(exp.Expr)), "contains") {

		var (
			col string
			val any
		)

		for _, v := range exp.Vars {
			switch v.Type {
			case expr.VarTypeColumn:
				col = fmt.Sprintf("%s", v.Value)
			case expr.VarTypeValue:
				val = v.Value
			}
		}

		return map[string]interface{}{
			"match_phrase": map[string]interface{}{
				col: val,
			},
		}, nil

	}

	idx := 0
	for _, v := range []byte(exp.Expr) {
		switch v {
		case '?':
			ev := exp.Vars[idx]
			switch ev.Type {
			case expr.VarTypeColumn:
				colName := fmt.Sprintf("%s", ev.Value)
				varCol := colName
				if strings.Contains(varCol, ".") {
					varCol = varCol[strings.LastIndex(varCol, ".")+1:]
				}
				ifCond = append(ifCond, fmt.Sprintf(`doc.containsKey('%s') && !doc['%s'].empty && doc['%s'].size() > 0`, colName, colName, colName))

				varScript.WriteString(fmt.Sprintf(`
							   def %s = doc['%s'].value;
                               if (%s instanceof Long || %s instanceof Integer)  {
                               		%s = (int)%s;
								}`, varCol, colName, varCol, varCol, varCol, varCol))
				expScript.WriteString(varCol)
			case expr.VarTypeValue:
				expScript.WriteString(expr.ExplainExprValue(dorm.DaMeng, ev.Value))
			}
			idx++
		default:
			expScript.WriteByte(v)
		}
	}

	opt := strings.ToLower(strings.TrimSpace(flt.Operator))
	switch opt {
	case filter.Eq, filter.Neq, filter.GT, filter.LT, filter.GTE, filter.LTE:
		if opt == filter.Eq {
			opt = "=="
		}
		returnScript.WriteString("return ")
		returnScript.WriteString("(" + expScript.String() + ")")
		returnScript.WriteString(opt)
		returnScript.WriteString(expr.ExplainExprValue(dorm.DaMeng, flt.Value))
		returnScript.WriteString(";")
	case filter.IN:
		returnScript.WriteString("return params.pm1.contains(")
		returnScript.WriteString(expScript.String())
		returnScript.WriteString(");")
	case filter.NotIn:
		returnScript.WriteString("return !params.pm1.contains(")
		returnScript.WriteString(expScript.String())
		returnScript.WriteString(");")
	case filter.Like:
		//注意表达式返回的必须是string类型，不能是text，在es中对应java的String类型才有contains函数
		//如果表达式返回的是索引中某个字段的值，那么这个值必须是keyword类型
		returnScript.WriteString(fmt.Sprintf("return (%s).contains(%s);", expScript.String(), strings.ReplaceAll(expr.ExplainExprValue(dorm.DaMeng, flt.Value), "%", "")))
	case filter.IS:
		//Not Supported
		return nil, errs.NewClientError("暂不支持")
	case filter.IsNot:
		//Not Supported
		return nil, errs.NewClientError("暂不支持")
	default:
		returnScript.WriteString("return ")
		returnScript.WriteString(expScript.String())
		returnScript.WriteString(";")
	}

	if len(ifCond) > 0 {
		sourceScript.WriteString("if (")
		sourceScript.WriteString(strings.Join(ifCond, " && "))
		sourceScript.WriteString(") {")
		sourceScript.WriteString(varScript.String())
		sourceScript.WriteString(returnScript.String())
		sourceScript.WriteString("}")
		sourceScript.WriteString("return false;")
	} else {
		sourceScript.WriteString(varScript.String())
		sourceScript.WriteString(returnScript.String())
	}

	internalScript := map[string]interface{}{
		"source": sourceScript.String(),
	}
	if flt.Value != nil {
		internalScript["params"] = map[string]interface{}{
			"pm1": flt.Value,
		}
	}

	script := map[string]interface{}{
		"script": map[string]interface{}{
			"script": internalScript,
		},
	}
	return script, nil
}

func KeywordToFilters(column string, keyword string) []filter.Filter {
	rs := make([]filter.Filter, 0, 1)
	if keyword != "" {
		scs := utils.Split(keyword, ",", "，", " ")
		for _, v := range scs {
			rs = append(rs, filter.Filter{
				LogicalOperator: filter.And,
				Column:          column,
				Operator:        filter.Like,
				Value:           v,
			})
		}
	}

	return rs
}
