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
	"io"
	"strings"
	"time"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/scroll"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	jsoniter "github.com/json-iterator/go"
)

const (
	MaxResultWindow = 10000
)

type AggregateResult struct {
	Agg        ctype.Map
	Group      []ctype.Map
	GroupTotal int64 //被统计的数据的条数，比如10条数据进行聚合，聚合后有三条聚合结果，那么GroupTotal为3，QueryTotal为10
	AfterKey   ctype.Map
	QueryTotal int64 //统计结果条数
}

type Api struct {
	Err             error
	Total           int64 //返回操作后的结果总数，比如withCount的结果
	esApi           *es_api.Api
	dbType          dorm.DbType //当前部署的是什么数据库，不是dorm.ES，是需要具体的数据库
	idx             string
	query           []filter.Filter
	nestedQuery     map[string][]filter.Filter //针对nested的查询, //path -> filters
	orderBy         []dorm.OrderBy             //排序
	columns         []string                   //查询的列
	body            map[string]interface{}
	offset          int
	limit           int
	forceLimit      bool
	groupCol        []string               //需要被聚合的字段
	aggCol          []dorm.AggregateColumn //聚合后需要取更多的字段的聚合值
	groupCol2       map[string][]string    //针对聚合的字段是nested对象。index field name -> nested object field
	aggCol2         map[string][]dorm.AggregateColumn
	aggCol2Alias    map[string]string //当取max、min、avg字段的时候，如果有别名需要对应到es响应中对应的字段agg_column的字段
	groupCountField string
	having2         map[string][]filter.Having
	afterKey        map[string]ctype.Map

	groupAggAlias map[string]string //当取max、min、avg字段的时候，如果有别名需要对应到es响应中对应的字段agg_column的字段
	aggAlias      map[string]string //当取max、min、avg字段的时候，如果有别名需要对应到es响应中对应的字段agg_column的字段
	groupAggCol   []dorm.AggregateColumn
	having        []filter.Having

	withCount  bool   //是否返回总数，最终值体现在Total字段上。针对agg：是否自动根据GroupCol进行count统计，不统计就是不取Count值，而只是取聚合后的字段值
	countAlias string //统计返回字段的别名，这个主要是针对Aggregate

	codecConfig CodecConfig
}

// NewApi
// 暂时不支持多个idx，留待后续扩展
func NewApi(esApi *es_api.Api, idx string, opts ...Option) *Api {
	sc := &Api{
		esApi:         esApi,
		aggAlias:      make(map[string]string),
		groupAggAlias: make(map[string]string),
		nestedQuery:   make(map[string][]filter.Filter),
	}
	sc.idx = idx
	sc.body = make(map[string]interface{})
	for _, opt := range opts {
		opt(sc)
	}
	return sc
}

func (s *Api) SetDbType(dbType dorm.DbType) *Api {
	s.dbType = dbType
	return s
}
func (s *Api) getDbType() dorm.DbType {
	if s.dbType != "" {
		return s.dbType
	}
	return dorm.PostgreSQL
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
func (s *Api) AddNestedFilters(path string, flt ...filter.Filter) *Api {
	if len(flt) < 1 || path == "" {
		return s
	}
	if s.nestedQuery == nil {
		s.nestedQuery = make(map[string][]filter.Filter)
	}
	s.nestedQuery[path] = append(s.nestedQuery[path], flt...)
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
	s.forceLimit = true
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

// AddGroupCol2
// path针对是nested类型字段的名称，如果path为空，则表示普通字段
func (s *Api) AddGroupCol2(path string, col ...string) *Api {
	if s.groupCol2 == nil {
		s.groupCol2 = make(map[string][]string)
	}
	s.groupCol2[path] = append(s.groupCol2[path], col...)
	return s
}
func (s *Api) AddAggCol2(path string, col ...dorm.AggregateColumn) *Api {
	s.aggCol = append(s.aggCol, col...)
	if s.aggCol2 == nil {
		s.aggCol2 = make(map[string][]dorm.AggregateColumn)
	}
	s.aggCol2[path] = append(s.aggCol2[path], col...)
	return s
}
func (s *Api) AddHaving2(path string, hv ...filter.Having) *Api {
	if s.having2 == nil {
		s.having2 = make(map[string][]filter.Having)
	}
	s.having2[path] = append(s.having2[path], hv...)
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
			if resp.Hits.Total != nil && s.withCount {
				s.Total = resp.Hits.Total.Value
			}
			return err
		}).
		Call(func(cl *caller.Caller) error {
			if dest == nil {
				return nil
			}
			//handleCodec(s, plugin.ProcessAfter, reflect.ValueOf(dest))
			return nil
		}).Err

	return s
}
func (s *Api) Scroll(dest any, scrollId string, scrollDuration string) (newScrollId string) {
	var (
		err        error
		resp       *search.Response
		scrollResp *scroll.Response
	)

	newScrollId = scrollId
	if s.check().Err != nil {
		return newScrollId
	}

	s.Err = caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			if scrollId == "" {
				resp, err = s.Search(scrollDuration)
			} else {
				scrollResp, err = s.esApi.EsTypedApi.Scroll().
					ScrollId(scrollId).
					Scroll(scrollDuration).
					Do(context.Background())
			}
			return err
		}).
		Call(func(cl *caller.Caller) error {
			if resp == nil && scrollResp == nil { // 当 index 不存在时，resp 是 nil
				return errs.NewClientError("index 查询报错")
			}
			var hits [][]byte
			if scrollId == "" {
				newScrollId = *resp.ScrollId_
				if len(resp.Hits.Hits) > 0 {
					hits = slice.Map(resp.Hits.Hits, func(_ int, v types.Hit) []byte { return v.Source_ })
				}
			} else {
				newScrollId = *scrollResp.ScrollId_
				if len(scrollResp.Hits.Hits) > 0 {
					hits = slice.Map(scrollResp.Hits.Hits, func(_ int, v types.Hit) []byte { return v.Source_ })
				}
			}
			//处理相应结果
			if len(hits) > 0 {
				bf := bytes.Buffer{}
				bf.WriteByte('[')
				bf.Write(bytes.Join(hits, utils.StringToBytes(",")))
				bf.WriteByte(']')
				if dest != nil {
					err = jsoniter.Unmarshal(bf.Bytes(), dest)
				}
			}
			return err
		}).
		Call(func(cl *caller.Caller) error {
			if dest == nil {
				return nil
			}
			//handleCodec(s, plugin.ProcessAfter, reflect.ValueOf(dest))
			return nil
		}).Err

	return newScrollId
}
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
	if resp == nil { // 当 index 不存在时，resp 是 nil
		s.Err = errs.NewClientError("index 查询报错")
		return s
	}
	scrollId := resp.ScrollId_
	defer func() {
		s.ClearScroll(*scrollId)
	}()

	if len(resp.Hits.Hits) <= 0 {
		return s
	}
	hits = append(hits, slice.Map(resp.Hits.Hits, func(_ int, v types.Hit) []byte { return v.Source_ })...)

	for {
		scrollResp, err1 := s.esApi.EsTypedApi.Scroll().ScrollId(*scrollId).Scroll(scrollDuration).Do(ct)
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
			if scrollDuration == "" { // scroll 时，不支持设置 track_total_hits
				s.body["track_total_hits"] = s.withCount
			}
			s.body["from"] = s.offset
			if s.forceLimit || s.limit > 0 {
				s.body["size"] = s.limit
			}
			return nil
		}).
		Call(func(cl *caller.Caller) error {
			s.body["query"], s.Err = s.resolveQuery()
			return s.Err
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
			LatencyLog(start, s.idx, optSearch, searchBody, GetSearchResultCode(err))
			return err
		}).Err

	return resp, s.Err
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
			s.body["query"], s.Err = s.resolveQuery()
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

func (s *Api) GroupBy() ([]AggregateResult, error) {

	var (
		err    error
		arList = make([]AggregateResult, 0, 1)
	)

	if s.check().Err != nil {
		return arList, s.Err
	}

	var (
		ofs = s.getOffset()
		lmt = s.getLimit()
	)

	for k, _ := range s.groupCol2 {
		ar := AggregateResult{}
		if len(s.afterKey) > 0 {
			if x, ok := s.afterKey[k]; ok && len(x) > 0 {
				ar, err = s.reqGroup(k, lmt, x)
				arList = append(arList, ar)
				if err != nil {
					return arList, err
				}
				continue
			}
		}

		tmpOfs := ofs
		for tmpOfs > 5000 {
			ar, _ = s.reqGroup(k, 5000, ar.AfterKey) //丢弃掉的请求，为了分页
			tmpOfs = tmpOfs - 5000
		}
		if tmpOfs > 0 {
			ar, _ = s.reqGroup(k, tmpOfs, ar.AfterKey) //丢弃掉的请求，为了分页
		}
		ar, err = s.reqGroup(k, lmt, ar.AfterKey) //正式请求
		if err != nil {
			return arList, err
		}
		arList = append(arList, ar)
	}

	return arList, s.Err
}

func (s *Api) reqGroup(path string, size int, after map[string]interface{}) (AggregateResult, error) {

	var (
		err        error
		searchBody []byte
		resp       *esapi.Response
		group      = make([]ctype.Map, 0, 1)
		afterKey   = ctype.Map{}
		ar         = AggregateResult{}
		respMap    ctype.Map
	)

	if err = s.check().Err; err != nil {
		return ar, err
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
			s.body["query"], s.Err = s.resolveQuery()
			return s.Err
		}).
		Call(func(cl *caller.Caller) error {

			//TODO 应该添加onlyCount的支持
			aggs := make(map[string]interface{})
			groupBucket := s.resolveGroupNested(path, size, after)
			groupCount := s.resolveGroupCountNested(path)
			for k, v := range groupBucket {
				aggs[k] = v
			}
			for k, v := range groupCount {
				aggs[k] = v
			}
			s.body["aggs"] = aggs
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

			err = jsoniter.Unmarshal(bs, &respMap)
			return err
		}).
		Call(func(cl *caller.Caller) error {
			group, afterKey, err = s.readGroupBuckets(path, respMap)
			ar.Group = group
			ar.AfterKey = afterKey
			ar.GroupTotal = s.readGroupCount(path, respMap)
			return err
		}).Err

	return ar, s.Err
}

func (s *Api) readGroupBuckets(path string, respMap ctype.Map) ([]ctype.Map, ctype.Map, error) {

	var (
		group = make([]ctype.Map, 0, 10)
		gb    = s.formatGroupBucketsKey(path)
		gf    = s.formatGroupFilterKey(path)
		gn    = s.formatGroupNestedKey(path)
	)

	//获取聚合后的列表
	buckets := ctype.GetMapValue(respMap, fmt.Sprintf("aggregations.%s.buckets", gb))
	if buckets == nil {
		buckets = ctype.GetMapValue(respMap, fmt.Sprintf("aggregations.%s.%s.buckets", gf, gb))
	}
	if buckets == nil {
		buckets = ctype.GetMapValue(respMap, fmt.Sprintf("aggregations.%s.%s.%s.buckets", gn, gf, gb))
	}
	if buckets == nil {
		return nil, nil, nil
	}

	afterKey := ctype.GetMapValue(respMap, fmt.Sprintf("aggregations.%s.after_key", gb))
	if afterKey == nil {
		afterKey = ctype.GetMapValue(respMap, fmt.Sprintf("aggregations.%s.%s.after_key", gf, gb))
	}
	if afterKey == nil {
		afterKey = ctype.GetMapValue(respMap, fmt.Sprintf("aggregations.%s.%s.%s.after_key", gn, gf, gb))
	}

	for _, x := range ctype.ToSlice(buckets) {
		if bucket := s.readBucket(ctype.ToMap(x)); len(bucket) > 0 {
			group = append(group, bucket)
		}
	}

	return group, ctype.ToMap(afterKey), nil

}

func (s *Api) readBucket(bucket ctype.Map) ctype.Map {
	rs := ctype.ToMap(ctype.GetMapValue(bucket, "key"))
	rs[s.getCountAlias()] = ctype.GetMapValue(bucket, "doc_count")
	for k, v := range s.aggCol2Alias {
		rs[v] = ctype.GetMapValue(bucket, fmt.Sprintf("%s.value", k))
	}
	return rs
}

func (s *Api) readGroupCount(path string, respMap ctype.Map) int64 {
	if !s.withCount {
		return 0
	}
	var (
		nk = s.formatGroupCountNestedKey(path)
		fk = s.formatGroupCountFilterKey(path)
		ck = s.formatGroupCountKey(path)
	)
	tt := ctype.GetMapValue(respMap, fmt.Sprintf("aggregations.%s.value", ck))
	if tt == nil {
		tt = ctype.GetMapValue(respMap, fmt.Sprintf("aggregations.%s.%s.value", fk, ck))
	}
	if tt == nil {
		tt = ctype.GetMapValue(respMap, fmt.Sprintf("aggregations.%s.%s.%s.value", nk, fk, ck))
	}
	return ctype.ToInt64(tt)
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

func (s *Api) ClearScroll(scrollId string) {
	_, _ = s.esApi.EsTypedApi.ClearScroll().ScrollId(scrollId).Do(context.Background())
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
func (s *Api) resolveQuery() (map[string]interface{}, error) {
	rsQr := make(map[string]interface{})
	//解析查询条件
	qr, err1 := ResolveEsQuery(s.getDbType(), s.query...)
	if err1 != nil {
		s.Err = err1
		return qr, err1
	}
	if bf, ok := qr["bool"].(map[string]interface{}); ok {
		for k, v := range bf {
			rsQr[k] = v
		}
	}

	boolFilter := make([]map[string]interface{}, 0)
	if bf, ok := rsQr["filter"].([]map[string]interface{}); ok {
		boolFilter = bf
	}

	//解析nested，//TODO 还没有解决两个path之间如果是或的问题
	for k, v := range s.nestedQuery {
		nqr, err2 := ResolveEsQuery(s.getDbType(), v...)
		if err2 != nil {
			s.Err = err2
			return nil, err2
		}
		nestQr := map[string]interface{}{
			"nested": map[string]interface{}{
				"path":  k,
				"query": nqr,
			},
		}
		boolFilter = append(boolFilter, nestQr)
	}
	rsQr["filter"] = boolFilter

	if len(rsQr) > 0 {
		return map[string]interface{}{
			"bool": rsQr,
		}, nil
		//s.body["query"] = map[string]interface{}{
		//	"bool": rsQr,
		//}
	}
	return map[string]interface{}{}, nil
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

func (s *Api) formatGroupCol(path string, col string) string {
	col = strings.ReplaceAll(col, ".", "_")
	if path == "" {
		return col
	}
	return fmt.Sprintf("%s_%s", path, col)
}
func (s *Api) formatAggCol(path string, col string, aggFunc string) string {
	col = strings.ReplaceAll(col, ".", "_")
	if path != "" {
		col = fmt.Sprintf("%s_%s", path, col)
	}
	if aggFunc != "" {
		col = fmt.Sprintf("%s_%s", col, aggFunc)
	}
	return col
}
func (s *Api) formatGroupBucketsKey(path string) string {
	if path != "" {
		return fmt.Sprintf("%s_%s", path, "group_buckets")
	}
	return "group_buckets"
}
func (s *Api) formatGroupFilterKey(path string) string {

	if path != "" {
		return fmt.Sprintf("%s_%s", path, "filter_buckets")
	}
	return "filter_buckets"
}
func (s *Api) formatGroupNestedKey(path string) string {

	if path != "" {
		return fmt.Sprintf("%s_%s", path, "nested_buckets")
	}
	return "nested_buckets"
}

func (s *Api) formatGroupCountKey(path string) string {
	if path != "" {
		return fmt.Sprintf("%s_%s", path, "group_count")
	}
	return "group_count"
}
func (s *Api) formatGroupCountFilterKey(path string) string {

	if path != "" {
		return fmt.Sprintf("%s_%s", path, "filter_group_count")
	}
	return "filter_group_count"
}
func (s *Api) formatGroupCountNestedKey(path string) string {

	if path != "" {
		return fmt.Sprintf("%s_%s", path, "nested_group_count")
	}
	return "nested_group_count"
}

func (s *Api) resolveGroupCountNested(path string) map[string]interface{} {
	var (
		gf = s.resolveGroupCountFilter(path)
	)
	if path == "" {
		return gf
	}

	return map[string]interface{}{
		s.formatGroupCountNestedKey(path): map[string]interface{}{
			"nested": map[string]interface{}{
				"path": path,
			},
			"aggs": gf,
		},
	}
}
func (s *Api) resolveGroupCountFilter(path string) map[string]interface{} {
	var (
		rs           = make(map[string]interface{})
		nestedFilter = make(map[string]interface{})
		gc           = s.resolveGroupCount(path)
	)

	if flt, ok := s.nestedQuery[path]; ok {
		nestedFilter, s.Err = ResolveEsQuery(s.getDbType(), flt...)
		if len(nestedFilter) > 0 {
			rs["filter"] = nestedFilter
		}
	}
	if len(rs) < 1 {
		return gc
	}
	rs["aggs"] = gc
	return map[string]interface{}{
		s.formatGroupCountFilterKey(path): rs,
	}

}
func (s *Api) resolveGroupCount(path string) map[string]interface{} {

	groupCountField := ""
	//TODO ES如何支持多个字段的统计
	if s.groupCol2 != nil {
		groupCountField = s.groupCol2[path][0]
	}

	if groupCountField == "" {
		groupCountField = "id"
	}

	return map[string]interface{}{
		s.formatGroupCountKey(path): map[string]interface{}{
			"cardinality": map[string]interface{}{
				"field":               groupCountField + ".keyword",
				"precision_threshold": 10000,
			},
		},
	}
}

func (s *Api) resolveGroupNested(path string, size int, after map[string]interface{}) map[string]interface{} {

	var (
		gf = s.resolveGroupFilter(path, size, after)
	)

	if path == "" {
		return gf
	}
	return map[string]interface{}{
		s.formatGroupNestedKey(path): map[string]interface{}{
			"nested": map[string]interface{}{
				"path": path,
			},
			"aggs": gf,
		},
	}
}
func (s *Api) resolveGroupFilter(path string, size int, after map[string]interface{}) map[string]interface{} {

	var (
		nestedFilter = make(map[string]interface{})
		rs           = make(map[string]interface{})
		gb           = s.resolveGroupBuckets(path, size, after)
	)

	if flt, ok := s.nestedQuery[path]; ok {
		nestedFilter, s.Err = ResolveEsQuery(s.getDbType(), flt...)
		if len(nestedFilter) > 0 {
			rs["filter"] = nestedFilter
		}
	}
	if len(rs) < 1 {
		return gb
	}
	rs["aggs"] = gb
	return map[string]interface{}{
		s.formatGroupFilterKey(path): rs,
	}
}
func (s *Api) resolveGroupBuckets(path string, size int, after map[string]interface{}) map[string]interface{} {
	var (
		cs       = s.resolveComposite(path, size, after)
		aggCol   = s.resolveAggCol2(path)
		hv       = s.resolveHaving2(path)
		rs       = make(map[string]interface{})
		colAndHv = make(map[string]interface{})
	)
	if len(cs) > 0 {
		rs["composite"] = cs
	}
	if len(hv) > 0 {
		for k, v := range hv {
			colAndHv[k] = v
		}
	}
	if len(aggCol) > 0 {
		for k, v := range aggCol {
			colAndHv[k] = v
		}
	}
	if len(colAndHv) > 0 {
		rs["aggs"] = colAndHv
	}
	if len(rs) < 1 {
		return map[string]interface{}{}
	}

	return map[string]interface{}{
		s.formatGroupBucketsKey(path): rs,
	}
}

func (s *Api) resolveComposite(path string, size int, after map[string]interface{}) map[string]interface{} {

	if s.groupCol2 == nil {
		return nil
	}
	rs := make(map[string]interface{})
	rs["size"] = size
	if after != nil && len(after) > 0 {
		rs["after"] = after
	}
	sources := make([]map[string]interface{}, 0)
	if cols, ok := s.groupCol2[path]; ok {
		for _, col := range cols {
			gCol := map[string]interface{}{
				s.formatGroupCol("", col): map[string]interface{}{
					"terms": map[string]interface{}{
						"field": col,
						//"order" //不支持排序，composite的排序只是支持针对group column的排序
					},
				},
			}
			sources = append(sources, gCol)
		}
	}
	rs["sources"] = sources
	return rs
}
func (s *Api) resolveHaving2(path string) map[string]interface{} {

	if s.having2 == nil {
		return nil
	}
	hv, ok := s.having2[path]
	if !ok {
		return nil
	}

	subAgg := make(map[string]interface{})
	for i, v := range hv {
		havingName := fmt.Sprintf("%s_%s_having_%d", v.AggFunc, v.Column, i)
		havingPath := s.formatAggCol(path, v.Column, v.AggFunc)
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
func (s *Api) resolveAggCol2(path string) map[string]interface{} {

	if s.aggCol2 == nil {
		return nil
	}
	aggColList, ok := s.aggCol2[path]
	if !ok {
		return nil
	}

	if s.aggCol2Alias != nil {
		s.aggCol2Alias = make(map[string]string)
	}

	aggColParam := make(map[string]interface{})
	for _, v := range aggColList {
		//aggCount 的聚合是默认的，path是doc_count，不需要额外聚合函数计算
		if strings.TrimSpace(strings.ToLower(v.AggFunc)) == dorm.AggFuncCount {
			if v.Alias != "" {
				s.countAlias = v.Alias
			}
			continue
		}
		aggAliasName := v.Alias
		if aggAliasName == "" {
			aggAliasName = s.formatAggCol("", v.Column, v.AggFunc)
		}

		aggField := s.formatAggCol(path, v.Column, v.AggFunc) ///重复了就进行覆盖
		aggColParam[aggField] = map[string]interface{}{
			v.AggFunc: s.resolveAggFunc(v.AggFunc, v.Column),
		}
		s.aggCol2Alias[aggField] = aggAliasName
	}
	return aggColParam
}
