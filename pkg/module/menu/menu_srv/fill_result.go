package menu_srv

import (
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/common/run"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/module/menu"
	"gorm.io/gorm"
	"sync"
)

type Result interface {
	menu.Menu | ctype.Map
}

func fill(cfg *hook.SrvConfig, listRs []menu.Menu) error {

	var (
		err error
		wg  = &sync.WaitGroup{}
	)
	wg.Add(1)
	run.Go(func() {
		defer wg.Done()
		err = errs.Cover(err, fillParent(cfg.OriginDB, cfg.Param, listRs))
	})

	wg.Wait()
	return err
}

func fillParent(db *gorm.DB, args *dto.Param, listRs []menu.Menu) error {
	if !args.WithParent {
		return nil
	}

	for i, _ := range listRs {
		lbMap, err := menu.LoadAllMenu(db)
		if err != nil {
			return err
		}
		err = fillParentRecursive(db, &listRs[i], lbMap)
		if err != nil {
			return err
		}
	}
	return nil
}

func fillParentRecursive(db *gorm.DB, dic *menu.Menu, allDic map[string]menu.Menu) error {
	pDic, ok := allDic[dic.ParentId]
	if !ok {
		return nil
	}
	dic.Parent = &pDic
	if pDic.ParentId != "" {
		return fillParentRecursive(db, &pDic, allDic)
	}
	return nil
}
