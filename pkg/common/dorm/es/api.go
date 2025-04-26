package es

import (
	"bytes"
	"context"
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
	"github.com/duke-git/lancet/v2/slice"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	jsoniter "github.com/json-iterator/go"
	"io"
	"strings"
	"time"
)

const (
	MaxResultWindow = 10000
)

type Api struct {
	Err           error
	Total         int64 //返回操作后的结果总数，比如withCount的结果
	esApi         *es_api.Api
	idx           string
	query         []filter.Filter
	orderBy       []dorm.OrderBy //排序
	columns       []string       //查询的列
	body          map[string]interface{}
	offset        int
	limit         int
	groupCol      []string               //需要被聚合的字段
	aggCol        []dorm.AggregateColumn //聚合后需要取更多的字段的聚合值
	aggAlias      map[string]string      //当取max、min、avg字段的时候，如果有别名需要对应到es响应中对应的字段agg_column的字段
	groupAggCol   []dorm.AggregateColumn
	groupAggAlias map[string]string //当取max、min、avg字段的时候，如果有别名需要对应到es响应中对应的字段agg_column的字段
	having        []filter.Having

	withCount  bool   //是否返回总数，最终值体现在Total字段上。针对agg：是否自动根据GroupCol进行count统计，不统计就是不取Count值，而只是取聚合后的字段值
	countAlias string //统计返回字段的别名，这个主要是针对Aggregate
}

// NewApi
// 暂时不支持多个idx，留待后续扩展
func NewApi(esApi *es_api.Api, idx string) *Api {
	sc := &Api{
		esApi:         esApi,
		aggAlias:      make(map[string]string),
		groupAggAlias: make(map[string]string),
	}
	sc.idx = idx
	sc.body = make(map[string]interface{})
	return sc
}

func (s *Api) Index(idx string) *Api {
	s.idx = idx
	return s
}
func (s *Api) AddColumn(col ...string) *Api {
	if len(col) < 1 {
		return s
	}
	s.columns = utils.Merge(s.columns, col...)
	return s
}
func (s *Api) AddFilters(flt ...filter.Filter) *Api {
	if len(flt) < 1 {
		return s
	}
	s.query = append(s.query, flt...)
	return s
}
func (s *Api) OrderBy(orderBy ...dorm.OrderBy) *Api {
	if len(orderBy) < 1 {
		return s
	}
	s.orderBy = append(s.orderBy, orderBy...)
	return s
}
func (s *Api) Offset(offset int) *Api {
	if offset < 0 {
		return s
	}
	s.offset = offset
	return s
}
func (s *Api) Limit(limit int) *Api {
	//也可能设置为0，表示不取结果，比如只进行count统计
	if limit < 0 {
		return s
	}
	s.limit = limit
	return s
}
func (s *Api) WithCount(flag bool) *Api {
	s.withCount = flag
	return s
}
func (s *Api) AddGroupCol(col ...string) *Api {
	s.groupCol = append(s.groupCol, col...)
	return s
}

// AddAggCol
// 可以对聚合之后的其余字段进行Aggregate操作，包括max、min、avg 和count(distinct col)
func (s *Api) AddAggCol(col ...dorm.AggregateColumn) *Api {
	s.aggCol = append(s.aggCol, col...)
	return s
}
func (s *Api) AddGroupAggCol(col ...dorm.AggregateColumn) *Api {
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
func (s *Api) AddHaving(hv ...filter.Having) *Api {
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
func (s *Api) Find(dest any) *Api {
	var (
		err  error
		resp *search.Response
	)
	if s.check().Err != nil {
		return s
	}

	s.Err = caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			resp, err = s.Search("")
			return err
		}).
		Call(func(cl *caller.Caller) error {
			if resp == nil { // 当 index 不存在时，resp 是 nil
				return errs.NewClientError("index 查询报错")
			}
			//处理相应结果
			if len(resp.Hits.Hits) > 0 {
				bf := bytes.Buffer{}
				bf.WriteByte('[')
				hasPre := false
				for _, v := range resp.Hits.Hits {
					if v.Source_ != nil {
						if hasPre {
							bf.WriteByte(',')
						}
						bf.Write(v.Source_)
						hasPre = true
					}
				}
				bf.WriteByte(']')
				if dest != nil {
					err = jsoniter.Unmarshal(bf.Bytes(), dest)
				}
			}
			if resp.Hits.Total != nil {
				s.Total = resp.Hits.Total.Value
			}
			return err
		}).Err

	return s
}

// LoadAll
// TODO 关于Scroll需要重构及适应新版本
func (s *Api) LoadAll(dest any) *Api {
	var (
		hits           = make([][]byte, 0) // 查询出来的结果
		ct             = context.Background()
		pageSize       = MaxResultWindow
		scrollDuration = "1m"
	)
	if s.check().Err != nil {
		return s
	}
	s.Offset(0).Limit(pageSize)

	resp, err := s.Search(scrollDuration)
	if err != nil {
		s.Err = err
		return s
	}
	scrollId := resp.ScrollId_
	defer func() {
		_, _ = s.esApi.EsTypedApi.ClearScroll().ScrollId(*scrollId).Do(ct)
	}()

	if len(resp.Hits.Hits) <= 0 {
		return s
	}
	hits = append(hits, slice.Map(resp.Hits.Hits, func(_ int, v types.Hit) []byte { return v.Source_ })...)

	for {
		scrollResp, err1 := s.esApi.EsTypedApi.Scroll().ScrollId(*scrollId).Do(ct)
		if err1 != nil {
			s.Err = err1
			return s
		}
		if len(scrollResp.Hits.Hits) <= 0 {
			break
		}
		hits = append(hits, slice.Map(scrollResp.Hits.Hits, func(_ int, v types.Hit) []byte { return v.Source_ })...)
	}

	if dest != nil {
		bf := bytes.Buffer{}
		bf.WriteByte('[')
		bf.Write(bytes.Join(hits, utils.StringToBytes(",")))
		bf.WriteByte(']')
		s.Err = jsoniter.Unmarshal(bf.Bytes(), dest)
	}
	s.Total = int64(len(hits))
	return s
}
func (s *Api) Search(scrollDuration string) (*search.Response, error) {
	var (
		err        error
		searchBody []byte
		resp       *search.Response
	)
	if s.check().Err != nil {
		return nil, s.Err
	}

	s.Err = caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			if len(s.columns) > 0 {
				s.body["_source"] = s.columns
			}
			st := s.resolveOrderBy()
			if len(st) > 0 {
				s.body["sort"] = st
			}
			if s.withCount {
				s.body["track_total_hits"] = true
			}
			if s.limit > 0 {
				s.body["size"] = s.limit
			}
			if s.offset > 0 {
				s.body["from"] = s.offset
			}
			return nil
		}).
		Call(func(cl *caller.Caller) error {
			//解析查询条件
			qr, err1 := ResolveEsQuery(dorm.DaMeng, s.query...)
			if len(qr) > 0 {
				s.body["query"] = qr
			}
			return err1
		}).
		Call(func(cl *caller.Caller) error {
			//序列化请求体
			searchBody, err = jsoniter.Marshal(s.body)
			return err
		}).
		Call(func(cl *caller.Caller) error {
			//发起请求查询
			var (
				start = time.Now()
			)
			esSearch := s.esApi.EsTypedApi.Search().Index(s.idx)
			if scrollDuration != "" {
				esSearch.Scroll(scrollDuration)
			}
			resp, err = esSearch.Raw(bytes.NewReader(searchBody)).Do(context.Background())
			LatencyLog(start, s.idx, optSearch, searchBody, getSearchResultCode(err))
			return err
		}).Err

	return resp, nil
}
func (s *Api) Count() int64 {
	return s.Limit(0).WithCount(true).Find(nil).Total
}
func (s *Api) Aggregate() (AggregateResult, error) {

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
			//聚合统计的时候，就不取列表结果
			s.body["size"] = 0
			return s.Err
		}).
		Call(func(cl *caller.Caller) error {
			//组装search请求中请求体的的query部分的内容
			//解析查询条件
			qr, err1 := ResolveEsQuery(dorm.DaMeng, s.query...)
			if len(qr) > 0 {
				s.body["query"] = qr
			}
			return err1
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
			LatencyLog(start, s.idx, optSearch, searchBody, getSearchResultCode(err))
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

			if s.withCount {
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

// Delete
// dest可以是id的切片，也可以是实体对象的切片，但是实体对象要包含ID字段，也可以是map的切片，但是map里也需要id字段
// 最终会根据ID进行删除
func (s *Api) Delete(dest any) error {
	return Delete(s.esApi, s.idx, dest)
}

// Upsert
// dest 是具体要操作的文档的切片，文档内容需要包含ID字段
// 如果对应的id数据已经存在则更新，不存在则插入
func (s *Api) Upsert(dest any) error {
	return Upsert(s.esApi, s.idx, dest)
}

func (s *Api) check() *Api {
	if s.idx == "" {
		s.Err = errs.NewClientError("es查询索引不能为空")
	}
	return s
}
func (s *Api) getOffset() int {
	if s.offset < 0 {
		return 0
	}
	return s.offset
}
func (s *Api) getLimit() int {
	if s.limit < 1 {
		return 10
	}
	return s.limit
}
func (s *Api) getCountAlias() string {
	if s.countAlias == "" {
		return DefaultCountAlias
	}
	return s.countAlias
}
func (s *Api) getAggField(col string) string {
	col = strings.ReplaceAll(col, ".", "_")
	return fmt.Sprintf("%s_%s", aggName, col)
}
func (s *Api) getGroupAggField(aggFunc, column string) string {
	column = strings.ReplaceAll(column, ".", "_")
	return fmt.Sprintf("%s_%s", aggFunc, column)
}
func (s *Api) getGroupRealCol(col string) string {
	if strings.Contains(col, ".") {
		col = col[strings.LastIndex(col, ".")+1:]
	}
	return col
}

func (s *Api) resolveOrderBy() []map[string]interface{} {
	return dorm.ResolveEsOrderBy(s.orderBy...)
}

// 返回的string 是 multi_terms 或 terms
func (s *Api) resolveAggCol() map[string]interface{} {
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

	if s.withCount {
		for _, v := range s.groupCol {
			bodyAgg[s.getAggField(v)] = map[string]interface{}{
				dorm.AggFuncCardinality: s.resolveAggFunc(dorm.AggFuncCardinality, v),
			}
		}
	}

	return bodyAgg
}
func (s *Api) resolveGroupTerms() map[string]interface{} {

	var (
		termsName = "terms"
		aggTerms  = make(map[string]interface{})
	)
	if len(s.groupCol) == 1 {
		aggTerms["field"] = s.groupCol[0] + ".keyword"
	} else if len(s.groupCol) > 1 {
		termsName = "multi_terms"
		tms := make([]map[string]interface{}, 0, len(s.groupCol))
		for _, v := range s.groupCol {
			tms = append(tms, map[string]interface{}{
				"field": v + ".keyword",
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
func (s *Api) resolveGroupSubAgg() map[string]interface{} {
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
func (s *Api) resolveGroupAggCol() map[string]interface{} {
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
func (s *Api) resolveGroupHaving() map[string]interface{} {
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
func (s *Api) resolveGroupHavingScript(opt, path string, val any) map[string]interface{} {

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
func (s *Api) resolveGroupPagination() map[string]interface{} {
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
func (s *Api) resolveAggFunc(aggFunc, field string) map[string]interface{} {
	aggF := map[string]interface{}{
		"field": field + ".keyword",
	}

	switch aggFunc {
	case dorm.AggFuncMax, dorm.AggFuncMin, dorm.AggFuncAvg:
		aggF = map[string]interface{}{
			"field": field,
		}
	}
	if aggFunc == dorm.AggFuncCardinality {
		aggF["precision_threshold"] = 10000 //cardinality的精度
	}
	return aggF
}
