package user_srv

import (
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/run"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/loader"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/authorize/role"
	"github.com/davycun/eta/pkg/module/dept"
	"github.com/davycun/eta/pkg/module/user"
	"sync"
)

type Result interface {
	ctype.Map | user.ListResult | user.User
}

func fill[T Result](cfg *hook.SrvConfig, listRs []T) error {
	for i, _ := range listRs {
		switch val := any(&listRs[i]).(type) {
		case *user.ListResult:
			val.Password = ""
		case *user.User:
			val.Password = ""
		case *ctype.Map:
			(*val).Set("password", "")
		}
	}

	var (
		err error
		wg  = &sync.WaitGroup{}
	)
	wg.Add(2)

	run.Go(func() {
		defer wg.Done()
		err = errs.Cover(err, fillDept(cfg, listRs))
	})
	run.Go(func() {
		defer wg.Done()
		err = errs.Cover(err, fillRole(cfg, listRs))
	})

	wg.Wait()
	return err
}

func fillDept[T Result](cfg *hook.SrvConfig, listRs []T) error {

	var (
		err     error
		fromIds = make([]string, 0, len(listRs))
		rlMap   = make(map[string][]dept.RelationDept)
	)

	for _, v := range listRs {
		fromIds = append(fromIds, entity.GetString(v, entity.IdDbName))
	}

	ld := loader.NewRelationEntityLoader[dept.Department, dept.RelationDept](cfg.OriginDB, constants.TableUser2Dept, constants.TableDept)
	ld.AddRelationColumns(dept.DefaultRelationDeptColumns...)
	rlMap, err = ld.LoadToMap(fromIds...)

	if err != nil {
		logger.Errorf("load user2dept err %s", err)
		return err
	}
	for i, v := range listRs {
		if dps, ok := rlMap[entity.GetString(v, entity.IdDbName)]; ok {
			entity.Set(&listRs[i], "user2dept", dps)
		}
	}
	return nil
}

func fillRole[T Result](cfg *hook.SrvConfig, listRs []T) error {

	var (
		err     error
		fromIds = make([]string, 0, len(listRs))
		rlMap   = make(map[string][]role.RelationRole)
	)
	for _, v := range listRs {
		fromIds = append(fromIds, entity.GetString(v, entity.IdDbName))
	}
	ld := loader.NewRelationEntityLoader[role.Role, role.RelationRole](cfg.OriginDB, constants.TableUser2Role, constants.TableRole)
	ld.AddRelationColumns(dept.DefaultRelationDeptColumns...)
	rlMap, err = ld.LoadToMap(fromIds...)

	if err != nil {
		logger.Errorf("load user2role err %s", err)
		return err
	}
	for i, v := range listRs {
		if roles, ok := rlMap[entity.GetString(v, entity.IdDbName)]; ok {
			entity.Set(&listRs[i], "user2role", roles)
		}
	}
	return nil
}
