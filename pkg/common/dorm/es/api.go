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
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	jsoniter "github.com/json-iterator/go"
)

const (
	MaxResultWindow   = 10000
	groupName         = "group"
	aggName           = "agg"
	DefaultCountAlias = "count"
)

type (
	Option          func(esApi *Api)
	AggregateResult struct {
		Group      []ctype.Map
		GroupTotal int64 //被统计的数据的条数，比如10条数据进行聚合，聚合后有三条聚合结果，那么GroupTotal为3，QueryTotal为10
		AfterKey   ctype.Map
		QueryTotal int64 //统计结果条数
	}

	Api struct {
		Err           error
		Total         int64 //返回操作后的结果总数，比如withCount的结果
		esApi         *es_api.Api
		dbType        dorm.DbType //当前部署的是什么数据库，不是dorm.ES，是需要具体的数据库
		idx           string
		body          map[string]interface{}
		columns       []string //查询的列
		queryFilters  []filter.Filter
		nestedFilters map[string][]filter.Filter //针对nested的查询, //path -> filters
		orderBy       []dorm.OrderBy             //排序
		offset        ctype.Integer
		limit         ctype.Integer
		groupCol      map[string][]string //针对聚合的字段是nested对象。index field name -> nested object field
		aggCol        map[string][]dorm.AggregateColumn
		aggColAlias   map[string]string //当取max、min、avg字段的时候，如果有别名需要对应到es响应中对应的字段agg_column的字段
		having        map[string][]filter.Having
		afterKey      map[string]ctype.Map

		loadAll    bool   //是否加载所有数据
		withCount  bool   //是否返回总数，最终值体现在Total字段上。针对agg：是否自动根据GroupCol进行count统计，不统计就是不取Count值，而只是取聚合后的字段值
		countAlias string //统计返回字段的别名，这个主要是针对Aggregate
	}
)

// NewApi
// 暂时不支持多个idx，留待后续扩展
func NewApi(esApi *es_api.Api, idx string, opts ...Option) *Api {
	sc := &Api{
		esApi:         esApi,
		nestedFilters: make(map[string][]filter.Filter),
		groupCol:      make(map[string][]string),
		aggCol:        make(map[string][]dorm.AggregateColumn),
		aggColAlias:   make(map[string]string),
		having:        make(map[string][]filter.Having),
		afterKey:      make(map[string]ctype.Map),
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
	s.queryFilters = append(s.queryFilters, flt...)
	return s
}
func (s *Api) AddNestedFilters(path string, flt ...filter.Filter) *Api {
	if len(flt) < 1 || path == "" {
		return s
	}
	if s.nestedFilters == nil {
		s.nestedFilters = make(map[string][]filter.Filter)
	}
	s.nestedFilters[path] = append(s.nestedFilters[path], flt...)
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
	s.offset = ctype.Integer{Valid: true, Data: int64(offset)}
	return s
}
func (s *Api) Limit(limit int) *Api {
	//也可能设置为0，表示不取结果，比如只进行count统计
	if limit < 0 {
		return s
	}

	s.limit = ctype.Integer{Valid: true, Data: int64(limit)}
	return s
}
func (s *Api) LoadAll(flag bool) *Api {
	s.loadAll = flag
	return s
}
func (s *Api) WithCount(flag bool) *Api {
	s.withCount = flag
	return s
}

// AddgroupCol
// path针对是nested类型字段的名称，如果path为空，则表示普通字段
func (s *Api) AddGroupCol(path string, col ...string) *Api {
	if s.groupCol == nil {
		s.groupCol = make(map[string][]string)
	}
	s.groupCol[path] = append(s.groupCol[path], col...)
	return s
}
func (s *Api) AddAggCol(path string, col ...dorm.AggregateColumn) *Api {
	if s.aggCol == nil {
		s.aggCol = make(map[string][]dorm.AggregateColumn)
	}
	s.aggCol[path] = append(s.aggCol[path], col...)
	return s
}
func (s *Api) AddHaving(path string, hv ...filter.Having) *Api {
	if s.having == nil {
		s.having = make(map[string][]filter.Having)
	}
	s.having[path] = append(s.having[path], hv...)
	return s
}

func (s *Api) Count() int64 {
	defer s.cleanBody()
	total, err := s.Limit(0).Offset(0).WithCount(true).LoadAll(false).Find(nil)
	if err != nil {
		s.Err = err
		logger.Errorf("es count err %s", err)
	}
	return total
}

func (s *Api) Find(dest any) (total int64, err error) {

	defer s.cleanBody()
	if s.check().Err != nil {
		return 0, s.Err
	}
	if s.loadAll {
		return s.findAll(dest)
	}
	return s.find(dest)
}

func (s *Api) cleanBody() {
	clear(s.body)
}

func (s *Api) findAll(dest any) (total int64, err error) {
	if !s.loadAll {
		return
	}
	var (
		dfSize      = 5000
		searchAfter []types.FieldValue
		orderBy     = s.resolveOrderBy()
		resp        *search.Response
		hits        = bytes.Buffer{}
	)

	if len(orderBy) < 1 {
		orderBy = []map[string]interface{}{
			{
				"eid": "asc",
			},
		}
	}

	//pit的作用是防止获取全量数据的过程中，有动态新增数据，导致获取不准
	if err = s.openPit(); err == nil {
		orderBy = append(orderBy, map[string]interface{}{"_shard_doc": "desc"})
	} else {
		return
	}

	hits.WriteByte('[')
	hasPre := false
	for {
		resp, err = s.reqSearch(dfSize, searchAfter, orderBy, true)
		if err != nil {
			return
		}

		var tmpHits []byte
		tmpHits, searchAfter, total, err = s.readSearchRespBytes(resp, false)
		if err != nil {
			return
		}
		if hasPre && len(tmpHits) > 0 {
			hits.WriteByte(',')
		}
		hits.Write(tmpHits)
		if len(tmpHits) > 0 {
			hasPre = true
		}
		//当前返回的数据已经小于size，表示已经load完毕
		if resp != nil && len(resp.Hits.Hits) < dfSize {
			break
		}
	}
	hits.WriteByte(']')
	err = jsoniter.Unmarshal(hits.Bytes(), dest)
	return
}
func (s *Api) openPit() error {
	var (
		start = time.Now()
	)
	//pit的作用是防止获取全量数据的过程中，有动态新增数据，导致获取不准
	pitFunc := s.esApi.EsTypedClient.OpenPointInTime(s.idx)
	pit, err := pitFunc.KeepAlive("1m").Do(context.Background())
	if err != nil {
		return err
	}
	s.body["pit"] = map[string]interface{}{
		"id":         pit.Id,
		"keep_alive": "1m",
	}
	LatencyLog(start, s.idx, optSearch, []byte(fmt.Sprintf("/%s/_pit?keep_alive=1m", s.idx)), GetSearchResultCode(err))
	return nil
}

func (s *Api) find(dest any) (total int64, err error) {
	var (
		ofs         = s.getOffset()
		lmt         = s.getLimit()
		searchAfter []types.FieldValue
		orderBy     = s.resolveOrderBy()
	)

	//有一个默认排序，避免第一次请求没有排序，然后点击第二页，第二页是经过排序的导致问题
	//默认排序不放在resolveOrderBy()函数的原因是，resolveOrderBy()可能还会被其他函数使用
	if len(orderBy) < 1 {
		orderBy = []map[string]interface{}{{"eid": "asc"}}
	}

	tmpOfs := ofs
	for tmpOfs > 5000 {
		_, searchAfter, err = s.reqSearchAndReadResponse(5000, searchAfter, orderBy, dest, false)
		if err != nil {
			return
		}
		tmpOfs = tmpOfs - 5000
	}
	if tmpOfs > 0 {
		_, searchAfter, err = s.reqSearchAndReadResponse(tmpOfs, searchAfter, orderBy, dest, false)
		if err != nil {
			return
		}
	}
	total, searchAfter, err = s.reqSearchAndReadResponse(lmt, searchAfter, orderBy, dest, true)
	return
}

func (s *Api) reqSearchAndReadResponse(size int, searchAfter []types.FieldValue, sort []map[string]interface{}, dest any, isResult bool) (int64, []types.FieldValue, error) {

	resp, err := s.reqSearch(size, searchAfter, sort, isResult)
	if err != nil {
		return 0, nil, err
	}
	total, after, err := s.readSearchResponse(resp, dest, isResult)
	return total, after, err
}

// isResult 表示是否是获取最终结果，还是只是中间过程为了获取searchAfterKey而已
func (s *Api) reqSearch(size int, searchAfter []types.FieldValue, sort []map[string]interface{}, isResult bool) (*search.Response, error) {
	var (
		err        error
		searchBody []byte
		resp       *search.Response
	)
	if s.check().Err != nil {
		return resp, s.Err
	}

	s.Err = caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			if isResult {
				if len(s.columns) > 0 {
					s.body["_source"] = s.columns
				}
				s.body["track_total_hits"] = s.withCount
			} else {
				s.body["_source"] = []string{"id", "eid"}
				s.body["track_total_hits"] = false
			}
			s.body["size"] = size
			if len(searchAfter) > 0 {
				s.body["search_after"] = searchAfter
			}
			if len(sort) > 0 {
				s.body["sort"] = sort
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
				start    = time.Now()
				esSearch = s.esApi.EsTypedApi.Search()
			)
			//如果有pit就无需携带索引信息
			if _, ok := s.body["pit"]; !ok {
				esSearch = esSearch.Index(s.idx)
			}
			resp, err = esSearch.Raw(bytes.NewReader(searchBody)).Do(context.Background())
			LatencyLog(start, s.idx, optSearch, searchBody, GetSearchResultCode(err))
			return err
		}).Err

	return resp, err
}
func (s *Api) readSearchResponse(resp *search.Response, dest any, isResult bool) (total int64, searchAfter []types.FieldValue, err error) {
	if resp == nil { // 当 index 不存在时，resp 是 nil
		return
	}
	if resp.Hits.Total != nil {
		total = resp.Hits.Total.Value
	}
	if len(resp.Hits.Hits) < 1 {
		return
	}
	if isResult && dest != nil {
		bs := make([]byte, 0)
		bs, searchAfter, total, err = s.readSearchRespBytes(resp, true)
		err = jsoniter.Unmarshal(bs, dest)
		return
	} else {
		for i, v := range resp.Hits.Hits {
			if i == len(resp.Hits.Hits)-1 {
				searchAfter = v.Sort
			}
		}
	}
	return
}
func (s *Api) readSearchRespBytes(resp *search.Response, withSquareBrackets bool) (hits []byte, searchAfter []types.FieldValue, total int64, err error) {
	if resp == nil { // 当 index 不存在时，resp 是 nil
		return
	}

	bf := bytes.Buffer{}
	if withSquareBrackets {
		bf.WriteByte('[')
	}
	hasPre := false
	for i, v := range resp.Hits.Hits {
		//只有取结果的时候才处理结果
		if v.Source_ != nil {
			if hasPre {
				bf.WriteByte(',')
			}
			bf.Write(v.Source_)
			hasPre = true
		}
		if i == len(resp.Hits.Hits)-1 {
			searchAfter = v.Sort
		}
	}
	if withSquareBrackets {
		bf.WriteByte(']')
	}
	hits = bf.Bytes()
	if resp.Hits.Total != nil {
		total = resp.Hits.Total.Value
	}
	return
}

func (s *Api) GroupBy() ([]AggregateResult, error) {

	defer s.cleanBody()
	if s.check().Err != nil {
		return nil, s.Err
	}

	if s.loadAll {
		return s.groupByAll()
	}
	return s.groupBy()
}

func (s *Api) groupBy() ([]AggregateResult, error) {
	var (
		err    error
		ofs    = s.getOffset()
		lmt    = s.getLimit()
		arList = make([]AggregateResult, 0, 1)
	)
	//groupCol 是map[string][]string，key是nested类型对象的字段名，如果需要聚合的字段不是嵌套对象，则key是空字符串
	for k, _ := range s.groupCol {
		ar := AggregateResult{}

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
func (s *Api) groupByAll() ([]AggregateResult, error) {

	var (
		dfSize   = 5000
		afterKey = ctype.Map{}
		arList   = make([]AggregateResult, 0)
	)
	//groupCol 是map[string][]string，key是nested类型对象的字段名，如果需要聚合的字段不是嵌套对象，则key是空字符串
	for k, _ := range s.groupCol {

		curAggRs := AggregateResult{}
		for {
			ar, err := s.reqGroup(k, dfSize, afterKey)
			if err != nil {
				return arList, err
			}

			curAggRs.Group = append(curAggRs.Group, ar.Group...)
			curAggRs.GroupTotal = ar.GroupTotal
			afterKey = ar.AfterKey
			if len(ar.Group) < dfSize {
				break
			}
		}
		arList = append(arList, curAggRs)
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

// 返回的第一个参数buckets列表，第二个参数是after_key
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
	for k, v := range s.aggColAlias {
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
	if !ctype.IsValid(s.offset) {
		return 0
	}

	return int(s.offset.Data)
}
func (s *Api) getLimit() int {
	if !ctype.IsValid(s.limit) {
		return 10
	}
	return int(s.limit.Data)
}
func (s *Api) getCountAlias() string {
	if s.countAlias == "" {
		return DefaultCountAlias
	}
	return s.countAlias
}
func (s *Api) resolveQuery() (map[string]interface{}, error) {
	rsQr := make(map[string]interface{})
	//解析查询条件
	qr, err1 := ResolveEsQuery(s.getDbType(), s.queryFilters...)
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
	for k, v := range s.nestedFilters {
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
	}
	return map[string]interface{}{}, nil
}
func (s *Api) resolveOrderBy() []map[string]interface{} {
	return dorm.ResolveEsOrderBy(s.orderBy...)
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

// 当groupBy多个字段的时候，我们需要通过runtime_mapping生成一个新的字段，这样cardinality才会准
func (s *Api) formatGroupCountRuntimeFieldKey(path string) string {
	if s.groupCol == nil {
		return ""
	}
	if x, ok := s.groupCol[path]; ok {
		rf := slice.JoinFunc(x, "_", func(v string) string {
			return strings.ReplaceAll(v, ".", "_")
		})
		if path != "" {
			return fmt.Sprintf("%s_%s", path, rf)
		}
		return rf
	}
	return ""
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

	if flt, ok := s.nestedFilters[path]; ok {
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
	rm := s.resolveGroupCountRuntimeMapping(path)
	if len(rm) > 0 {
		s.body["runtime_mappings"] = rm
	}
	return map[string]interface{}{
		s.formatGroupCountKey(path): map[string]interface{}{
			"cardinality": map[string]interface{}{
				"field":               s.resolveGroupCountFieldName(path),
				"precision_threshold": 10000,
			},
		},
	}
}
func (s *Api) resolveGroupCountFieldName(path string) string {

	if s.groupCol == nil {
		return ""
	}
	groupCols := s.groupCol[path]
	switch len(groupCols) {
	case 0:
		return ""
	case 1:
		//表示只是聚合一个字段，直接返回当前字段名称
		return groupCols[0] + ".keyword"
	default:
		//表示聚合多个字段，需要把多个字段通过runtime_mapping生成一个运行时字段，然后再对这个运行时字段进行cardinality
		return s.formatGroupCountRuntimeFieldKey(path)
	}
}

// 是否需要runtimeMapping，是根据聚合字段的多少来判断的，大于等于2个就需要
func (s *Api) resolveGroupCountRuntimeMapping(path string) map[string]interface{} {

	//当groupBy多个字段的时候，我们需要通过runtime_mapping生成一个新的字段，这样cardinality才会准
	if s.groupCol == nil {
		return nil
	}

	//如果只是聚合一个字段，那么就不需要runtime_mapping
	groupCols, ok := s.groupCol[path]
	if !ok || len(groupCols) < 2 {
		return nil
	}

	runtimeFieldName := s.formatGroupCountRuntimeFieldKey(path)
	return map[string]interface{}{
		runtimeFieldName: map[string]interface{}{
			"type":   "keyword",
			"script": s.resolveGroupCountRuntimeFieldScript(path),
		},
	}
}
func (s *Api) resolveGroupCountRuntimeFieldScript(path string) string {

	if s.groupCol == nil {
		return ""
	}
	var (
		cols, colExists = s.groupCol[path]
		sc              = ""
	)
	if !colExists {
		return ""
	}

	//this is nested type
	if path != "" {
		sc = `
			def result = [];
			def src = params['_source'];
			if (src.containsKey('%s')) {
			  def al = src['%s'];
			  def keys = new String[]{%s};
			  for (addr in al) {
				def rs = "";
				for (k in keys) {
				  if (addr.containsKey(k) && !addr[k].isEmpty()) {
					if (rs!="") {
					  rs += "_" + addr[k];
					}else{
					  rs += addr[k];
					}
				  }
				}
				if (rs != ""){
				  result.add(rs)
				}
			  }
			}
			for (rs in result) {
			  emit(rs); // 逐个 emit 每个 name
			}`
		keys := "'" + strings.Join(cols, "','") + "'" //def keys = new String[]{'addr_type','land_type'}; 此处不能加keyword，因为是从_source取出来的对象
		sc = fmt.Sprintf(sc, path, path, keys)
	} else {
		sc = `def keys = new String[]{%s};
				def rs = "";
				for (k in keys) {
				  if (doc.containsKey(k) && !doc[k].isEmpty()){
					def arr = doc[k];
					StringBuilder sb = new StringBuilder();
					for (int i = 0; i < arr.length; i++) {
						sb.append(String.valueOf(arr[i]));
						if (i < arr.length - 1) {
							sb.append('_');
						}
					}
					if (rs!=""){
					  rs += "_" + sb.toString();
					}else{
					  rs += sb.toString();
					}
				  }
				}
             emit(rs);
		`
		keys := "'" + strings.Join(cols, ".keyword','") + ".keyword'" //def keys = new String[]{'addr_ids_list.building_id.keyword','addr_ids_list.level2_id.keyword'};
		sc = fmt.Sprintf(sc, keys)
	}

	return sc
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

	if flt, ok := s.nestedFilters[path]; ok {
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
		aggCol   = s.resolveAggCol(path)
		hv       = s.resolveHaving(path)
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

	if s.groupCol == nil {
		return nil
	}
	rs := make(map[string]interface{})
	rs["size"] = size
	if after != nil && len(after) > 0 {
		rs["after"] = after
	}
	sources := make([]map[string]interface{}, 0)
	if cols, ok := s.groupCol[path]; ok {
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
func (s *Api) resolveAggCol(path string) map[string]interface{} {

	if s.aggCol == nil {
		return nil
	}
	aggColList, ok := s.aggCol[path]
	if !ok {
		return nil
	}

	if s.aggColAlias == nil {
		s.aggColAlias = make(map[string]string)
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
			v.AggFunc: s.resolveAggColFunc(v.AggFunc, v.Column),
		}
		s.aggColAlias[aggField] = aggAliasName
	}
	return aggColParam
}
func (s *Api) resolveAggColFunc(aggFunc, field string) map[string]interface{} {
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
func (s *Api) resolveHaving(path string) map[string]interface{} {

	if s.having == nil {
		return nil
	}
	hv, ok := s.having[path]
	if !ok {
		return nil
	}

	subAgg := make(map[string]interface{})
	for i, v := range hv {
		havingName := fmt.Sprintf("%s_%s_having_%d", v.AggFunc, v.Column, i)
		aggColPath := s.formatAggCol(path, v.Column, v.AggFunc)
		switch v.AggFunc {
		case dorm.AggFuncCount:
			havingName = fmt.Sprintf("%s_having_%d", v.AggFunc, i)
			aggColPath = "count_agg"
			subAgg[havingName] = map[string]interface{}{
				"bucket_selector": map[string]interface{}{
					"buckets_path": map[string]interface{}{
						aggColPath: "_count",
					},
					"script": s.resolveHavingScript(v.Operator, aggColPath, v.Value),
				},
			}
		case dorm.AggFuncMax, dorm.AggFuncMin, dorm.AggFuncAvg:
			subAgg[havingName] = map[string]interface{}{
				"bucket_selector": map[string]interface{}{
					"buckets_path": map[string]interface{}{
						aggColPath: aggColPath,
					},
					"script": s.resolveHavingScript(v.Operator, aggColPath, v.Value),
				},
			}
		}
	}
	return subAgg
}
func (s *Api) resolveHavingScript(opt, path string, val any) map[string]interface{} {

	script := make(map[string]interface{})
	switch opt {
	case filter.Eq, filter.Neq, filter.GT, filter.LT, filter.GTE, filter.LTE:
		script["source"] = fmt.Sprintf("params.%s %s %s", path, opt, expr.ExplainExprValue(dorm.DaMeng, val))
	case filter.IN:

		bd := strings.Builder{}
		bd.WriteString(fmt.Sprintf(`
							   def %s = params.%s;
                               if (%s instanceof Long || %s instanceof Integer || %s instanceof Float || %s instanceof Double)  {
									return params.pm1.contains((int)%s);
								} `, path, path, path, path, path, path, path))
		bd.WriteString(fmt.Sprintf(" return params.pm1.contains(%s);", path))
		script["source"] = bd.String()
		script["params"] = map[string]interface{}{
			"pm1": val,
		}

	case filter.NotIn:
		bd := strings.Builder{}
		bd.WriteString(fmt.Sprintf(`
							   def %s = params.%s;
                               if (%s instanceof Long || %s instanceof Integer || %s instanceof Float || %s instanceof Double)  {
									return !params.pm1.contains((int)%s);
								} `, path, path, path, path, path, path, path))
		bd.WriteString(fmt.Sprintf(" return !params.pm1.contains(%s);", path))
		script["source"] = bd.String()
		script["params"] = map[string]interface{}{
			"pm1": val,
		}

	}
	return script
}
