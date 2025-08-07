package sqlbd

import (
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/filter"
	"github.com/davycun/eta/pkg/core/service/hook"
	"maps"
	"reflect"
)

type (
	BuildSql      func(cfg *hook.SrvConfig) (*SqlList, error)
	BuildEsFilter func(cfg *hook.SrvConfig) ([]filter.Filter, error)
	BuildEsAggCol func(cfg *hook.SrvConfig) []dorm.AggregateColumn
	SqlListOption func(s *SqlList)
)

// SqlList
// sqlMap中的key和rsMap中的key要保持对应，也就是可以针对每个命名的sql指定对应的sql结果接收类型
type SqlList struct {
	sqlMap   map[string]string //name -> sql
	EsFilter BuildEsFilter
	EsAggCol BuildEsAggCol           //额外的统计字段
	IsAgg    bool                    //是否是聚合相关的sql
	NeedScan bool                    //是否需要通过scan，也就是查询了额外的字段，不能只是通过固定的结构体来获取数据，比如Group语句需要NeedScan为true
	rsMap    map[string]reflect.Type //接收对应sql结果的类型
	onlyOne  map[string]bool         //name-> bool，表示对应的名字的sql是否只返回一条数据（可能只有一个字段比如统计，也可能有多个字段，但只有一条数据，通常统计结果类的设置为true）
}

func NewSqlList(option ...SqlListOption) *SqlList {
	sl := &SqlList{
		sqlMap: map[string]string{},
		rsMap:  make(map[string]reflect.Type),
	}
	for _, fc := range option {
		fc(sl)
	}
	return sl
}

// AddSql
// iface.Method -> sql
func (s *SqlList) AddSql(name, sql string) *SqlList {
	if s.sqlMap == nil {
		s.sqlMap = make(map[string]string)
	}
	s.sqlMap[name] = sql
	return s
}
func (s *SqlList) AddResultType(name string, rsType reflect.Type) *SqlList {
	if s.rsMap == nil {
		s.rsMap = make(map[string]reflect.Type)
	}
	s.rsMap[name] = rsType
	return s
}
func (s *SqlList) SetOnlyOne(name string, onlyOne bool) *SqlList {
	if s.onlyOne == nil {
		s.onlyOne = make(map[string]bool)
	}
	s.onlyOne[name] = onlyOne
	return s
}
func (s *SqlList) SetEsFilter(esFilter BuildEsFilter) *SqlList {
	s.EsFilter = esFilter
	return s
}
func (s *SqlList) SetEsAggCol(esAggCol BuildEsAggCol) *SqlList {
	s.EsAggCol = esAggCol
	return s
}
func (s *SqlList) SetNeedScan(needScan bool) *SqlList {
	s.NeedScan = needScan
	return s
}
func (s *SqlList) SetIsAgg(isAgg bool) *SqlList {
	s.IsAgg = isAgg
	return s
}

func (s *SqlList) ListSql() string {
	return s.sqlMap[ListSql]
}
func (s *SqlList) CountSql() string {
	return s.sqlMap[CountSql]
}
func (s *SqlList) TotalSql() string {
	return s.sqlMap[TotalSql]
}
func (s *SqlList) Sql(name string) string {
	return s.sqlMap[name]
}
func (s *SqlList) AllSql() map[string]string {
	return maps.Clone(s.sqlMap)
}
func (s *SqlList) OnlyOne(name string) bool {
	if s.onlyOne == nil {
		return false
	}
	return s.onlyOne[name]
}

func (s *SqlList) NewResultPointer(name string) any {
	if s.rsMap == nil {
		return nil
	}
	if x, ok := s.rsMap[name]; ok {
		return reflect.New(x).Interface()
	}
	return nil
}
func (s *SqlList) NewResultSlicePointer(name string) any {
	if s.rsMap == nil {
		return nil
	}
	if x, ok := s.rsMap[name]; ok {
		reflect.New(reflect.SliceOf(x)).Interface()
	}
	return nil
}
func (s *SqlList) NewResultOrSlicePointer(name string) any {

	if s.OnlyOne(name) {
		return s.NewResultPointer(name)
	}
	return s.NewResultSlicePointer(name)
}
