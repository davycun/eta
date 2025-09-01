package es

import (
	"bytes"
	"fmt"
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/dorm/es/es_api"
	"github.com/davycun/eta/pkg/common/dorm/expr"
	"github.com/davycun/eta/pkg/common/dorm/filter"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/utils"
	"io"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8/esapi"
	jsoniter "github.com/json-iterator/go"
)

const (
	groupName         = "group"
	aggName           = "agg"
	DefaultCountAlias = "count"
)

type Aggregate struct {
	esApi          *es_api.Api
	Err            error
	idx            string
	query          []filter.Filter
	groupCol       []string
	aggCol         []dorm.AggregateColumn
	aggAlias       map[string]string //当取max、min、avg字段的时候，如果有别名需要对应到es响应中对应的字段agg_column的字段
	groupAggCol    []dorm.AggregateColumn
	groupAggAlias  map[string]string //当取max、min、avg字段的时候，如果有别名需要对应到es响应中对应的字段agg_column的字段
	autoGroupCount bool              //是否自动根据GroupCol进行count统计，不统计就是不取Count值，而只是取聚合后的字段值
	having         []filter.Having
	orderBy        []dorm.OrderBy
	offset         int
	limit          int
	body           map[string]interface{}
	//GroupTotal     int64 //被统计的数据的条数，比如10条数据进行聚合，聚合后有三条聚合结果，那么AggTotal为3，QueryTotal为10
	//QueryTotal     int64 //统计结果条数

	countAlias string //统计返回字段的别名
}

func NewAggregate(esApi *es_api.Api) *Aggregate {
	sc := &Aggregate{
		esApi:         esApi,
		body:          make(map[string]interface{}),
		aggAlias:      make(map[string]string),
		groupAggAlias: make(map[string]string),
	}
	return sc
}

func (s *Aggregate) Index(idx string) *Aggregate {
	s.idx = idx
	return s
}
func (s *Aggregate) AddGroupCol(col ...string) *Aggregate {
	s.groupCol = append(s.groupCol, col...)
	return s
}

// AddAggCol
// 可以对聚合之后的其余字段进行Aggregate操作，包括max、min、avg 和count(distinct col)
func (s *Aggregate) AddAggCol(col ...dorm.AggregateColumn) *Aggregate {
	s.aggCol = append(s.aggCol, col...)
	return s
}
func (s *Aggregate) AddGroupAggCol(col ...dorm.AggregateColumn) *Aggregate {
	//只是为了取doc_count的别名
	for _, v := range col {
		fc := strings.TrimSpace(strings.ToLower(v.AggFunc))
		if fc == dorm.AggFuncCount {
			if v.Column == "*" || utils.ContainAny(s.groupCol, v.Column) {
				s.countAlias = v.Alias
				break
			}
		}
	}
	s.groupAggCol = append(s.groupAggCol, col...)
	return s
}

// AddHaving
// 暂时不支持子Having及OrHaving
func (s *Aggregate) AddHaving(hv ...filter.Having) *Aggregate {
	for _, v := range hv {
		if strings.TrimSpace(strings.ToLower(v.AggFunc)) == dorm.AggFuncCount {
			continue
		}
		ac := dorm.AggregateColumn{
			Column:  v.Column,
			AggFunc: v.AggFunc,
			Alias:   s.getGroupAggField(v.AggFunc, v.Column),
		}
		s.AddAggCol(ac)
	}
	s.having = append(s.having, hv...)
	return s
}
func (s *Aggregate) AddFilters(flt ...filter.Filter) *Aggregate {
	s.query = append(s.query, flt...)
	return s
}

// OrderBy
// orderBy 只是支持传入_count或者_key
func (s *Aggregate) OrderBy(orderBy ...dorm.OrderBy) *Aggregate {
	if len(orderBy) < 1 {
		return s
	}
	s.orderBy = append(s.orderBy, orderBy...)
	return s
}
func (s *Aggregate) Offset(offset int) *Aggregate {
	s.offset = offset
	return s
}
func (s *Aggregate) Limit(limit int) *Aggregate {
	s.limit = limit
	return s
}

func (s *Aggregate) WithGroupCount(flag bool) *Aggregate {
	s.autoGroupCount = flag
	return s
}
func (s *Aggregate) WithQueryCount(flag bool) *Aggregate {
	if flag {
		s.body["track_total_hits"] = true
	} else {
		s.body["track_total_hits"] = false
	}
	return s
}

func (s *Aggregate) getOffset() int {
	if s.offset < 0 {
		return 0
	}
	return s.offset
}
func (s *Aggregate) getLimit() int {
	if s.limit < 1 {
		return 10
	}
	return s.limit
}
func (s *Aggregate) getCountAlias() string {
	if s.countAlias == "" {
		return DefaultCountAlias
	}
	return s.countAlias
}
func (s *Aggregate) getAggField(col string) string {
	col = strings.ReplaceAll(col, ".", "_")
	return fmt.Sprintf("%s_%s", aggName, col)
}
func (s *Aggregate) getGroupAggField(aggFunc, column string) string {
	column = strings.ReplaceAll(column, ".", "_")
	return fmt.Sprintf("%s_%s", aggFunc, column)
}
func (s *Aggregate) getGroupRealCol(col string) string {
	if strings.Contains(col, ".") {
		col = col[strings.LastIndex(col, ".")+1:]
	}
	return col
}

func (s *Aggregate) resolveQuery() map[string]interface{} {
	query, err := ResolveEsQuery(dorm.DaMeng, s.query...)
	if err != nil {
		s.Err = err
	}
	return query
}

// 返回的string 是 multi_terms 或 terms
func (s *Aggregate) resolveAggCol() map[string]interface{} {
	//这里不用考虑重复的问题，重复了就行覆盖就行了，因为对于取某个字段的Max、Min、Avg，就算AggregateColumn重复了也没问题，map覆盖即可
	//也就是aggAliasName
	bodyAgg := make(map[string]interface{})

	for _, v := range s.aggCol {
		//aggCount 的聚合是默认的，path是_count，不需要额外聚合函数计算
		af := strings.TrimSpace(strings.ToLower(v.AggFunc))
		if af == dorm.AggFuncCount {
			if v.Column == "*" {
				v.AggFunc = dorm.AggFuncValueCount
				continue
			} else {
				v.AggFunc = dorm.AggFuncCardinality
			}
		}
		aggAliasName := v.Alias
		aggField := s.getAggField(v.Column) ///重复了就进行覆盖
		if aggAliasName == "" {
			aggAliasName = aggField
		}
		s.aggAlias[aggAliasName] = aggField

		bodyAgg[aggField] = map[string]interface{}{
			v.AggFunc: s.resolveAggFunc(v.AggFunc, v.Column),
		}
	}

	if s.autoGroupCount {
		for _, v := range s.groupCol {
			bodyAgg[s.getAggField(v)] = map[string]interface{}{
				dorm.AggFuncCardinality: s.resolveAggFunc(dorm.AggFuncCardinality, v),
			}
		}
	}

	return bodyAgg
}
func (s *Aggregate) resolveGroupTerms() map[string]interface{} {

	var (
		termsName = "terms"
		aggTerms  = make(map[string]interface{})
	)
	if len(s.groupCol) == 1 {
		aggTerms["field"] = s.groupCol[0]
	} else if len(s.groupCol) > 1 {
		termsName = "multi_terms"
		tms := make([]map[string]interface{}, 0, len(s.groupCol))
		for _, v := range s.groupCol {
			tms = append(tms, map[string]interface{}{
				"field": v,
			})
		}
		aggTerms["terms"] = tms
	}

	//size的处理，为了达到分页的目的，这里的size是计算的最大size
	aggTerms["size"] = s.getOffset() + s.getLimit()

	orderBy := "_count"
	order := "desc"
	if len(s.orderBy) > 0 {
		//暂时只支持一个orderBy
		ob := s.orderBy[0]
		order = dorm.ResolveOrderDesc(ob.Asc)
		if utils.ContainAny(s.groupCol, ob.Column) {
			orderBy = "_key"
		}
	}
	aggTerms["order"] = map[string]interface{}{
		orderBy: order,
	}

	return map[string]interface{}{
		termsName: aggTerms,
	}
}
func (s *Aggregate) resolveGroupSubAgg() map[string]interface{} {
	subAgg := make(map[string]interface{})
	for k, v := range s.resolveGroupAggCol() {
		subAgg[k] = v
	}
	for k, v := range s.resolveGroupHaving() {
		subAgg[k] = v
	}
	for k, v := range s.resolveGroupPagination() {
		subAgg[k] = v
	}

	return subAgg
}
func (s *Aggregate) resolveGroupAggCol() map[string]interface{} {
	subAgg := make(map[string]interface{})
	//这里不用考虑重复的问题，重复了就行覆盖就行了，因为对于取某个字段的Max、Min、Avg，就算AggregateColumn重复了也没问题，map覆盖即可
	//也就是aggAliasName
	for _, v := range s.groupAggCol {
		//aggCount 的聚合是默认的，path是_count，不需要额外聚合函数计算
		if strings.TrimSpace(strings.ToLower(v.AggFunc)) == dorm.AggFuncCount {
			if v.Column == "*" || utils.ContainAny(s.groupCol, v.Column) {
				s.countAlias = v.Alias
				continue
			}
			v.AggFunc = dorm.AggFuncCardinality
		}
		aggAliasName := v.Alias
		aggField := s.getGroupAggField(v.AggFunc, v.Column) ///重复了就进行覆盖
		if aggAliasName == "" {
			aggAliasName = aggField
		}
		s.groupAggAlias[aggAliasName] = aggField

		subAgg[aggField] = map[string]interface{}{
			v.AggFunc: s.resolveAggFunc(v.AggFunc, v.Column),
		}
	}
	return subAgg
}
func (s *Aggregate) resolveGroupHaving() map[string]interface{} {
	subAgg := make(map[string]interface{})
	for i, v := range s.having {
		havingName := fmt.Sprintf("%s_%s_having_%d", v.AggFunc, v.Column, i)
		havingPath := s.getGroupAggField(v.AggFunc, v.Column)
		switch v.AggFunc {
		case dorm.AggFuncCount:
			havingName = fmt.Sprintf("%s_having_%d", v.AggFunc, i)
			havingPath = "count_agg"
			subAgg[havingName] = map[string]interface{}{
				"bucket_selector": map[string]interface{}{
					"buckets_path": map[string]interface{}{
						havingPath: "_count",
					},
				},
				"script": s.resolveGroupHavingScript(v.Operator, havingPath, v.Value),
			}
		case dorm.AggFuncMax, dorm.AggFuncMin, dorm.AggFuncAvg:
			subAgg[havingName] = map[string]interface{}{
				"bucket_selector": map[string]interface{}{
					"buckets_path": map[string]interface{}{
						havingPath: havingPath,
					},
				},
				"script": s.resolveGroupHavingScript(v.Operator, havingPath, v.Value),
			}
		}
	}
	return subAgg
}
func (s *Aggregate) resolveGroupHavingScript(opt, path string, val any) map[string]interface{} {

	script := make(map[string]interface{})
	switch opt {
	case filter.Eq, filter.Neq, filter.GT, filter.LT, filter.GTE, filter.LTE:
		script["source"] = fmt.Sprintf("params.%s %s %s", path, opt, expr.ExplainExprValue(dorm.DaMeng, val))
	case filter.IN:

		bd := strings.Builder{}
		bd.WriteString(fmt.Sprintf(`
							   def %s = params.%s;
                               if (%s instanceof Long || %s instanceof Integer)  {
                               		%s = (int)%s;
								}\n`, path, path, path, path, path, path))
		bd.WriteString(fmt.Sprintf("return params.pm1.contains(%s);", path))
		script["source"] = bd.String()
		script["params"] = map[string]interface{}{
			"pm1": val,
		}

	case filter.NotIn:
		bd := strings.Builder{}
		bd.WriteString(fmt.Sprintf(`
							   def %s = params.%s;
                               if (%s instanceof Long || %s instanceof Integer)  {
                               		%s = (int)%s;
								}\n`, path, path, path, path, path, path))
		bd.WriteString(fmt.Sprintf("return !params.pm1.contains(%s);", path))
		script["source"] = bd.String()
		script["params"] = map[string]interface{}{
			"pm1": val,
		}

	}
	return script
}
func (s *Aggregate) resolveGroupPagination() map[string]interface{} {
	subAgg := map[string]interface{}{
		"agg_page": map[string]interface{}{
			"bucket_sort": map[string]interface{}{
				"from": s.getOffset(),
				"size": s.getLimit(),
			},
		},
	}
	return subAgg
}
func (s *Aggregate) resolveAggFunc(aggFunc, field string) map[string]interface{} {
	aggF := map[string]interface{}{
		"field": field,
	}
	if aggFunc == dorm.AggFuncCardinality {
		aggF["precision_threshold"] = 10000
	}
	return aggF
}

func (s *Aggregate) check() *Aggregate {
	if s.idx == "" {
		s.Err = errs.NewClientError("es查询索引不能为空")
	}
	return s
}
func (s *Aggregate) Find() (AggregateResult, error) {
	var (
		err        error
		searchBody []byte
		resp       *esapi.Response
		ar         = AggregateResult{Group: make([]ctype.Map, 0, 1), Agg: ctype.Map{}}
	)

	if s.check().Err != nil {
		return ar, s.Err
	}

	defer func() {
		if resp != nil {
			err1 := resp.Body.Close()
			if err1 != nil {
				logger.Infof("close es response body err %s", err1)
			}
		}
	}()

	s.Err = caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			s.body["size"] = 0
			return s.Err
		}).
		Call(func(cl *caller.Caller) error {
			//组装search请求中请求体的的query部分的内容
			query := s.resolveQuery()
			if len(query) > 0 {
				s.body["query"] = query
			}
			return s.Err
		}).
		Call(func(cl *caller.Caller) error {
			//组装search请求的aggs（聚合）部分的内容
			bodyAgg := s.resolveAggCol()
			//暂时只支持对一个字段进行统计，group 的count统计值是不准确的

			aggs := make(map[string]interface{})
			for k, v := range s.resolveGroupTerms() {
				aggs[k] = v
			}

			subAgg := s.resolveGroupSubAgg()
			if len(subAgg) > 0 {
				aggs["aggs"] = subAgg
			}

			if len(aggs) > 0 {
				bodyAgg[groupName] = aggs
			}

			if len(bodyAgg) > 0 {
				s.body["aggs"] = bodyAgg
			}

			return s.Err
		}).
		Call(func(cl *caller.Caller) error {
			//序列化请求体为json
			searchBody, err = jsoniter.Marshal(s.body)
			return err
		}).
		Call(func(cl *caller.Caller) error {
			//发起请求
			var (
				start    = time.Now()
				esClient = s.esApi.EsApi
			)
			resp, err = esClient.Search(
				esClient.Search.WithIndex(s.idx),
				esClient.Search.WithBody(bytes.NewReader(searchBody)),
			)
			LatencyLog(start, s.idx, optSearch, searchBody, GetSearchResultCode(err))
			return err
		}).
		Call(func(cl *caller.Caller) error {

			bs, err2 := io.ReadAll(resp.Body)
			if err2 != nil {
				return err2
			}
			if resp.IsError() {
				return errs.NewServerError(string(bs))
			}

			var respMap ctype.Map
			err = jsoniter.Unmarshal(bs, &respMap)
			if err != nil {
				return err
			}

			if s.autoGroupCount {
				for _, v := range s.groupCol {
					tt := ctype.GetMapValue(respMap, fmt.Sprintf("aggregations.%s.value", s.getAggField(v)))
					if ar.GroupTotal == 0 {
						ar.GroupTotal += ctype.ToInt64(tt)
					} else {
						ar.GroupTotal *= ctype.ToInt64(tt)
					}
				}
			}

			//一些非group内的独立统计的结果
			for key, val := range s.aggAlias {
				ar.Agg[key] = ctype.GetMapValue(respMap, fmt.Sprintf("aggregations.%s.value", val))
			}

			//获取总数
			ar.QueryTotal = ctype.ToInt64(ctype.GetMapValue(respMap, "hits.total.value"))

			//获取聚合后的列表
			buckets := ctype.GetMapValue(respMap, fmt.Sprintf("aggregations.%s.buckets", groupName))
			for _, x := range ctype.ToSlice(buckets) {
				bucket := ctype.ToMap(x)
				obj := ctype.Map{}

				if len(s.groupCol) == 1 {
					obj[s.getGroupRealCol(s.groupCol[0])] = ctype.GetMapValue(bucket, "key")
				} else if len(s.groupCol) > 1 {
					keys := ctype.ToSlice(ctype.GetMapValue(bucket, "key"))
					for i, col := range s.groupCol {
						obj[s.getGroupRealCol(col)] = keys[i]
					}
				}
				obj[s.getCountAlias()] = ctype.GetMapValue(bucket, "doc_count")

				for key, val := range s.groupAggAlias {
					obj[key] = ctype.GetMapValue(bucket, fmt.Sprintf("%s.value", val))
				}

				ar.Group = append(ar.Group, obj)
			}
			return err
		}).Err

	return ar, s.Err
}
