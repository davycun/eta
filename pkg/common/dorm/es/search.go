package es

import (
	"bytes"
	"context"
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/es/es_api"
	"github.com/davycun/eta/pkg/common/dorm/filter"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/common/utils"
	"time"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	jsoniter "github.com/json-iterator/go"
)

type Search struct {
	esApi   *es_api.Api
	Err     error
	idx     string
	query   []filter.Filter
	orderBy []dorm.OrderBy
	body    map[string]interface{}
	Total   int64
}

func NewSearch(esApi *es_api.Api) *Search {
	sc := &Search{
		esApi: esApi,
	}
	sc.body = make(map[string]interface{})
	return sc
}

func (s *Search) Index(idx string) *Search {
	s.idx = idx
	return s
}
func (s *Search) AddColumn(col ...string) *Search {
	if len(col) < 1 {
		return s
	}
	if cs, ok := s.body["_source"]; ok {
		s.body["_source"] = utils.Merge(cs.([]string), col...)
	} else {
		s.body["_source"] = col
	}
	return s
}
func (s *Search) AddFilters(flt ...filter.Filter) *Search {
	if len(flt) < 1 {
		return s
	}
	s.query = append(s.query, flt...)
	return s
}
func (s *Search) OrderBy(orderBy ...dorm.OrderBy) *Search {
	if len(orderBy) < 1 {
		return s
	}
	s.body["sort"] = dorm.ResolveEsOrderBy(orderBy...)
	return s
}
func (s *Search) Offset(offset int) *Search {
	s.body["from"] = offset
	return s
}
func (s *Search) Limit(limit int) *Search {
	s.body["size"] = limit
	return s
}
func (s *Search) WithCount(flag bool) *Search {
	if flag {
		s.body["track_total_hits"] = true
	} else {
		s.body["track_total_hits"] = false
	}
	return s
}

func (s *Search) check() *Search {
	if s.idx == "" {
		s.Err = errs.NewClientError("es查询索引不能为空")
	}
	return s
}

func (s *Search) Find(dest any) *Search {
	var (
		err        error
		searchBody []byte
		resp       *search.Response
	)
	if s.check().Err != nil {
		return s
	}

	s.Err = caller.NewCaller().
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
			resp, err = s.esApi.EsTypedApi.Search().Index(s.idx).Raw(bytes.NewReader(searchBody)).Do(context.Background())
			LatencyLog(start, s.idx, optSearch, searchBody, GetSearchResultCode(err))
			return err
		}).
		Call(func(cl *caller.Caller) error {
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
func (s *Search) Count() int64 {
	return s.Limit(0).WithCount(true).Find(nil).Total
}
