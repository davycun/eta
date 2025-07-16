package sqlbd

import (
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/filter"
	"github.com/davycun/eta/pkg/core/service/hook"
)

type (
	BuildSql      func(cfg *hook.SrvConfig) (*SqlList, error)
	BuildEsFilter func(cfg *hook.SrvConfig) ([]filter.Filter, error)
	BuildEsAggCol func(cfg *hook.SrvConfig) []dorm.AggregateColumn
	SqlListOption func(s *SqlList)
)

type SqlList struct {
	sqlMap   map[string]string //name -> sql
	EsFilter BuildEsFilter
	EsAggCol BuildEsAggCol //额外的统计字段
	IsAgg    bool          //是否是聚合相关的sql
	NeedScan bool          //是否需要通过scan，也就是查询了额外的字段，不能只是通过固定的结构体来获取数据，比如Group语句需要NeedScan为true
}

func NewSqlList(option ...SqlListOption) *SqlList {
	sl := &SqlList{
		sqlMap: map[string]string{},
	}
	for _, fc := range option {
		fc(sl)
	}
	return sl
}

// AddSql
// iface.Method -> sql
func (s *SqlList) AddSql(name, sql string) *SqlList {
	s.sqlMap[name] = sql
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
