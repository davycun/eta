package user_srv

import (
	"errors"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/filter"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/core/builder"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/core/service/sqlbd"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/authorize/auth"
	"github.com/davycun/eta/pkg/module/dept"
	"github.com/davycun/eta/pkg/module/user"
	"github.com/duke-git/lancet/v2/slice"
)

var (
	NoDeptUser = errors.New("NoDeptUser")
	NoRoleUser = errors.New("NoRoleUser")
)

// SqlBuilder
// 该数据权限只在后台管理中启用:
// 1. 超级管理员和系统管理员角色可以在后台管理中看到所有的部门和所有的用户（包括虚拟用户）
// 2. 部门管理员可以在后台管理中看到他有部门管理员角色的所在部门及其子部门和部门及子部门的用户。（不包括虚拟用户）
// 3. 一般人如果被分配了后台管理的功能权限，那么在后台管理中只能看到他自己，看不到部门列表。
// disable_perm_filter 控制是否开启权限过滤

func buildListSql(cfg *hook.SrvConfig) (*sqlbd.SqlList, error) {
	var (
		userId = cfg.Ctx.GetContextUserId()
		col    = user.DefaultUserColumns
		extra  = cfg.Param.Extra.(*user.ListParam)
		dbType = dorm.GetDbType(cfg.OriginDB)
		scm    = dorm.GetDbSchema(cfg.OriginDB)
		args   = cfg.Param

		appDb     = cfg.Ctx.GetAppGorm()
		appScm    = dorm.GetDbSchema(appDb)
		appDbType = dorm.GetDbType(appDb)
		allSb     = make([]*builder.SqlBuilder, 0, 4)
		sqlList   = sqlbd.NewSqlList()
	)

	if !args.DisablePermFilter && !cfg.Ctx.GetContextIsManager() && !auth.IsSystemAdmin(cfg.Ctx.GetAppGorm(), userId) {
		rd, err1 := dept.LoadUser2DeptByUserId(cfg.Ctx, userId)
		if err1 != nil {
			return sqlList, err1
		}
		managerDept := slice.FilterMap(rd, func(i int, v dept.RelationDept) (string, bool) {
			return v.ToId, v.IsManager
		})
		if len(managerDept) < 1 {
			//只能看到他自己
			args.Filters = append(args.Filters, filter.Filter{
				LogicalOperator: filter.And,
				Column:          entity.IdDbName,
				Operator:        filter.Eq,
				Value:           userId,
			})
		} else {
			// 能看到他管理的部门及子部门
			extra.User2DeptFilters = append(extra.User2DeptFilters, filter.Filter{
				LogicalOperator: filter.And,
				Column:          entity.ToIdDbName,
				Operator:        filter.IN,
				Value:           managerDept,
			})
			args.Filters = append(args.Filters, filter.Filter{
				LogicalOperator: filter.And,
				Column:          "category",
				Operator:        filter.IN,
				Value:           user.NotVirtualUser,
			})
		}
	}
	//因为有可能r_user2dept和 t_user表不在同一个db内，所以不能做关联查询
	if len(extra.User2DeptFilters) > 0 || len(extra.DeptFilters) > 0 {
		tmpBd := builder.NewSqlBuilder(appDbType, appScm, constants.TableUser2Dept)
		tmpBd.AddFilter(extra.User2DeptFilters...)
		tmpBd.AddColumn(entity.FromIdDbName)

		if len(extra.DeptFilters) > 0 {
			tmpBd.AddTableFilter(constants.TableDept, extra.DeptFilters...)
			tmpBd.Join(appScm, constants.TableDept, entity.IdDbName, constants.TableUser2Dept, entity.ToIdDbName)
		}
		allSb = append(allSb, tmpBd)
	}

	if len(extra.User2RoleFilters) > 0 {
		tmpBd := builder.NewSqlBuilder(appDbType, appScm, constants.TableUser2Role)
		tmpBd.AddColumn(entity.FromIdDbName).AddFilter(extra.User2RoleFilters...)
		allSb = append(allSb, tmpBd)
	}

	cte := builder.NewCteSqlBuilder(dbType, scm, constants.TableUser)

	if len(allSb) > 0 {
		var first *builder.SqlBuilder
		first = allSb[0]
		for i := 1; i < len(allSb); i++ {
			first.UnionIntersect(allSb[i])
		}
		if len(allSb) < 2 {
			first.SetDistinct(true)
		}

		//如果用户表和用户部门表在同一个db上，就直接join，否则先查询userId再添加filter去查User
		dbCfg := dorm.GetDbConfig(cfg.OriginDB)
		appDbCfg := dorm.GetDbConfig(appDb)
		if dbCfg.Host == appDbCfg.Host && dbCfg.Port == appDbCfg.Port {
			//在同一个DB上
			cte.With("userId", first)
			cte.Join("", "userId", entity.FromIdDbName, constants.TableUser, entity.IdDbName)
		} else {

			//如果不在同一个DB上，只能是先把关联的APP的ID查询出来再作为filter条件去查询user表
			listSql, _, err := first.Build()
			if err != nil {
				return sqlList, err
			}
			var userIds []string
			err = dorm.RawFetch(listSql, appDb, &userId)
			if err != nil {
				return sqlList, err
			}
			if len(userIds) < 1 {
				return sqlList, nil
			}
			if len(userIds) < 6 {
				args.Filters = append(args.Filters, filter.Filter{
					LogicalOperator: filter.And,
					Column:          entity.IdDbName,
					Operator:        filter.IN,
					Value:           userIds,
				})
			} else {
				vb := builder.NewValueBuilder(dbType, entity.IdDbName, userIds...)
				cte.With("userId", vb)
				cte.Join("", "userId", entity.IdDbName, constants.TableUser, entity.IdDbName)
			}
		}
	}

	if len(args.Columns) > 0 {
		col = utils.Merge(args.Columns, "id", "name")
	}

	//只能查询当前APP的
	flt := filter.Filter{
		LogicalOperator: filter.And,
		Column:          "app_id",
		Operator:        filter.Eq,
		Value:           cfg.Ctx.GetContextAppId(),
		Filters: []filter.Filter{{
			LogicalOperator: filter.And,
			Filters:         args.Filters,
		}},
	}

	args.Filters = []filter.Filter{flt}
	cte.AddColumn(col...).AddFilter(args.Filters...)
	if !args.LoadAll {
		cte.Offset(args.GetOffset()).Limit(args.GetLimit())
	}

	listSql, countSql, err := cte.Build()
	sqlList.AddSql(sqlbd.ListSql, listSql).AddSql(sqlbd.CountSql, countSql)

	return sqlList, err
}
