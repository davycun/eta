package service

import (
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/expr"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/core/builder"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/core/service/sqlbd"
	"strings"
)

func AggregateSql(cfg *hook.SrvConfig) (*sqlbd.SqlList, error) {

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

	listSql, countSql, err := aggBd.Build()

	return sqlbd.NewSqlList().
		SetNeedScan(true).
		AddEsBuilder(sqlbd.ListSql, BuilderEsApiForAggregate).
		AddSql(sqlbd.ListSql, listSql).
		AddSql(sqlbd.CountSql, countSql), err
}

// PartitionSql
// 实现的是postgresql的窗口函数查询。假设有一张表，记录了大区、年份、月份、营收，四个字段
// 我们需要查询 每个大区在2020~2023年中 每个月的营收 和 当月所有大区的月总营收和年总营收，那么可以如下查询
// select 大区,年份,月份,营收,
//
//		sum(营收) over (partition by 年份,月份) as 月总计,
//	 sum(营收) over (partition by 年份) as 年总计
//
// from 营收表
// where 年份 in ('2020','2021','2023')
// 注意如果传入了distinct ，并且如果需要order by，那么order by中必须出现distinct的字段，并且排在order by语句的最左边
func PartitionSql(cfg *hook.SrvConfig) (*sqlbd.SqlList, error) {

	var (
		dbType         = dorm.GetDbType(cfg.OriginDB)
		scm            = dorm.GetDbSchema(cfg.OriginDB)
		idsAlias       = "ids"
		defaultColumns = entity.GetDefaultColumns(cfg.NewEntityPointer())
		mustCols       = entity.GetMustColumns(cfg.NewEntityPointer())
		tbName         = cfg.GetTableName()
	)

	cte := builder.NewCteSqlBuilder(dbType, scm, tbName)

	if len(cfg.Param.Distinct) > 0 {
		switch dbType {
		case dorm.DaMeng:
			cte.AddExprColumn(expr.ExpColumn{
				Expression: expr.Expression{
					Expr: "distinct ?",
					Vars: []expr.ExpVar{
						{
							Type:  expr.VarTypeValue,
							Value: strings.Join(cfg.Param.Distinct, ","),
						},
					},
				},
			})
		case dorm.PostgreSQL:
			cte.AddExprColumn(expr.ExpColumn{
				Expression: expr.Expression{
					Expr: "distinct on (?)",
					Vars: []expr.ExpVar{
						{
							Type:  expr.VarTypeValue,
							Value: strings.Join(cfg.Param.Distinct, ","),
						},
					},
				},
			})
		}
	}

	if len(cfg.Param.Columns) > 0 {
		cte.AddColumn(utils.Merge(cfg.Param.Columns, mustCols...)...)
	} else if len(cfg.Param.ExtraColumns) < 1 {
		cte.AddColumn(defaultColumns...)
	}
	cte.AddExprColumn(cfg.Param.ExtraColumns...)
	cte.AddPartitionColumn(cfg.Param.PartitionColumns...)

	filterBd := buildIdListSqlBuilder(cfg)
	if filterBd != nil {
		cte.With(idsAlias, filterBd)
		cte.Join("", idsAlias, entity.IdDbName, cfg.GetTableName(), entity.IdDbName)
	}
	cte.AddOrderBy(cfg.Param.OrderBy...)
	if !cfg.Param.LoadAll {
		cte.Offset(cfg.Param.GetOffset()).Limit(cfg.Param.GetLimit())
	}
	listSql, countSql, err := cte.Build()

	return sqlbd.NewSqlList().SetNeedScan(true).AddSql(sqlbd.ListSql, listSql).AddSql(sqlbd.CountSql, countSql), err
}
