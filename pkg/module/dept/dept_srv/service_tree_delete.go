package dept_srv

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/cache"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/core/builder"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/authorize/auth"
	"github.com/davycun/eta/pkg/module/dept"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

func (s *Service) TreeDelete(args *dto.Param, result *dto.Result) error {
	var (
		ids []string
		db  = dorm.Transaction(s.GetDB())
		err error
	)
	defer func() {
		if !dorm.InTransaction(s.GetDB()) {
			dorm.CommitOrRollback(db, err)
		}
	}()

	bd := builder.NewRecursiveSqlBuilder(dorm.GetDbType(s.GetDB()), dorm.GetDbSchema(s.GetDB()), constants.TableDept)
	bd.AddRecursiveFilter(args.Filters...).SetUp(args.IsUp)
	bd.AddColumn("id").AddFilter(args.Filters...)
	listSql, _, err := bd.Build()
	if err != nil {
		return err
	}
	// 获取出来关于要删除的部门以及子部门的ID 用于删除r_user2dept 表里面用户与部门的关联关系
	err = dorm.RawFetch(listSql, db, &ids)
	if len(ids) > 0 {
		err = deleteUser2Dept(db, ids)
		if err != nil {
			return err
		}
		err = deleteAuth2Role(db, ids)
		if err != nil {
			return err
		}
	}

	tx := entity.SetTableName(db, &dept.Department{})
	tx = tx.Where(fmt.Sprintf(`"%s"."id" in (%s)`, constants.TableDept, listSql)).Delete(&dept.Department{})

	if tx.Error != nil {
		return tx.Error
	}
	result.RowsAffected = tx.RowsAffected
	return tx.Error
}

func deleteUser2Dept(db *gorm.DB, ids []string) error {
	var (
		scm     = dorm.GetDbSchema(db)
		dbType  = dorm.GetDbType(db)
		listBd  = strings.Builder{}
		col     = dorm.JoinColumns(dbType, "", []string{"from_id"})
		user2Rs []string
		delBd   = strings.Builder{}
	)
	ids = ListStrToIds(ids)
	listBd.WriteString(fmt.Sprintf(`select %s from "%s"."r_user2dept" where "to_id" in (%s)`, col, scm, strings.Join(ids, ",")))
	err := dorm.RawFetch(listBd.String(), db, &user2Rs)
	if err != nil {
		return err
	}
	// 清理 UserDeptCache
	for _, id := range user2Rs {
		// 清理用户的token 防止删除完部门以后部门中的用户还能在页面去访问其他页面报500的错误
		err, _ = cache.Del(constants.RedisKey(constants.UserTokenKey, id))
		if err != nil {
			logger.Errorf("del user token cache err %s", err)
		}
		dept.DelUser2DeptCache(id)
	}
	if len(user2Rs) > 0 {
		delBd.WriteString(fmt.Sprintf(`DELETE FROM "%s"."r_user2dept" WHERE "to_id" IN (%s)`, scm, strings.Join(ids, ",")))
		err = dorm.RawFetch(delBd.String(), db, &ids)
		if err != nil {
			return err
		}
	}
	return err
}

func deleteAuth2Role(db *gorm.DB, ids []string) error {
	var deletedIds []string

	if len(ids) > 0 {
		err := db.Model(&auth.Auth2Role{}).
			Clauses(clause.Returning{Columns: []clause.Column{{Name: "to_id"}}}).
			Where(`"to_id" in ? and "to_table" = 't_department'`, ids).
			Delete(&deletedIds).Error
		if err != nil {
			return err
		}
	}
	// 清理 auth2Role cache
	if len(ids) > 0 {
		for _, id := range ids {
			auth.DelAuth2RoleCache(dorm.GetDbSchema(db), id)
		}
	}
	return nil
}

func ListStrToIds(list []string) []string {
	var ids []string
	for _, id := range list {
		ids = append(ids, fmt.Sprintf("'%s'", id))
	}
	return ids
}
