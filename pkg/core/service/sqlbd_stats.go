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

func PartitionSql(cfg *hook.SrvConfig) (*sqlbd.SqlList, error) {
	listSql, countSql, err := buildPartitionSql(cfg)
	return sqlbd.NewSqlList().SetNeedScan(true).AddSql(sqlbd.ListSql, listSql).AddSql(sqlbd.CountSql, countSql).SetEsFilter(buildListFilter), err
}

func buildPartitionSql(cfg *hook.SrvConfig) (listSql, countSql string, err error) {

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
	} else {
		cte.AddColumn(defaultColumns...)
	}
	if len(cfg.Param.ExtraColumns) > 0 {
		cte.AddExprColumn(cfg.Param.ExtraColumns...)
	}
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
	listSql, countSql, err = cte.Build()

	return
}
