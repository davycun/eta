package service

import (
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/core/builder"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/core/service/sqlbd"
)

func AggregateSql(cfg *hook.SrvConfig) (*sqlbd.SqlList, error) {
	var (
		sqlList = sqlbd.NewSqlList(iface.MethodAggregate, true)
	)
	listSql, countSql, err := buildListAggregateSql(cfg)
	sqlList.AddSql(sqlbd.ListSql, listSql).AddSql(sqlbd.CountSql, countSql).SetEsFilter(buildListFilter)
	return sqlList, err
}

func buildListAggregateSql(cfg *hook.SrvConfig) (listSql, countSql string, err error) {

	var (
		dbType = dorm.GetDbType(cfg.OriginDB)
		scm    = dorm.GetDbSchema(cfg.OriginDB)
		tbName = cfg.GetTableName()
		aggBd  = builder.NewAggregateSqlBuilder(dbType, scm, tbName).SetCteName("r")
	)

	filterBd := buildIdListSqlBuilder(cfg)
	if filterBd != nil {
		aggBd.With("ids", filterBd)
		aggBd.Join("", "ids", entity.IdDbName, tbName, entity.IdDbName)
	}

	aggBd.AddGroupColumn(cfg.Param.GroupColumns...).
		AddHavingFilter(cfg.Param.Having...).
		AddAggregateColumn(cfg.Param.AggregateColumns...).
		AddOrderBy(cfg.Param.OrderBy...).
		Offset(cfg.Param.GetOffset()).
		Limit(cfg.Param.GetPageSize())

	listSql, countSql, err = aggBd.Build()
	return
}
