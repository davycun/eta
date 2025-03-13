package dept_srv

import (
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/dorm/filter"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/common/run"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/core/builder"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/dept"
	"gorm.io/gorm"
	"sync"
)

type Result interface {
	dept.Department | ctype.Map
}

func fill(cfg *hook.SrvConfig, listRs []dept.Department) error {

	var (
		err error
		wg  = &sync.WaitGroup{}
	)
	wg.Add(1)
	run.Go(func() {
		defer wg.Done()
		err = errs.Cover(err, fillParent(cfg.Ctx.GetAppGorm(), cfg.Param, listRs))
	})

	wg.Wait()
	return err
}

func fillParent(db *gorm.DB, args *dto.Param, listRs []dept.Department) error {
	if !args.WithParent {
		return nil
	}

	var (
		err            error
		parentDeptList []dept.Department
	)
	pid := make([]string, 0, len(listRs))

	for _, v := range listRs {
		pid = append(pid, v.ParentId)
	}

	parentCte := builder.NewRecursiveSqlBuilder(dorm.GetDbType(db), dorm.GetDbSchema(db), constants.TableDept)
	cols := dept.DefaultColumns
	if len(args.Columns) > 0 {
		cols = utils.Merge(args.Columns, entity.IdDbName, "parent_id")
	}

	parentCte.AddColumn(cols...)
	parentListSql, _, err := parentCte.SetUp(true).AddRecursiveFilter(filter.Filter{Column: "id", Operator: filter.IN, Value: pid}).Build()
	if err != nil {
		return err
	}

	err = dorm.RawFetch(parentListSql, db, &parentDeptList)
	if err != nil {
		return err
	}

	parentDeptMap := make(map[string]dept.Department)
	for _, v := range parentDeptList {
		parentDeptMap[v.ID] = v
	}

	for i, _ := range listRs {
		fillParentRecursive(&listRs[i], parentDeptMap)
	}
	return err
}

func fillParentRecursive(dpt *dept.Department, parentDeptMap map[string]dept.Department) {
	if parentDept, ok := parentDeptMap[dpt.ParentId]; ok {
		dpt.Parent = &parentDept
		if parentDept.ParentId != "" {
			fillParentRecursive(&parentDept, parentDeptMap)
		}
	}
}
