package dict_srv

import (
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/common/run"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/module/dict"
	"gorm.io/gorm"
	"sync"
)

type Result interface {
	dict.Dictionary | ctype.Map
}

func processResult(cfg *hook.SrvConfig, listRs any) error {
	switch x := listRs.(type) {
	case []dict.Dictionary:
		return fill(cfg.OriginDB, cfg.Param, x)
	case []ctype.Map:
		//暂不支持
		//return fill(db, args, x)
	}
	return nil
}

func fill(db *gorm.DB, args *dto.Param, listRs []dict.Dictionary) error {

	var (
		err error
		wg  = &sync.WaitGroup{}
	)
	wg.Add(1)
	run.Go(func() {
		defer wg.Done()
		err = errs.Cover(err, fillParent(db, args, listRs))
	})

	wg.Wait()
	return err
}

func fillParent(db *gorm.DB, args *dto.Param, listRs []dict.Dictionary) error {
	if !args.WithParent {
		return nil
	}

	for i, _ := range listRs {
		lbMap, err := dict.LoadAllDictionary(db)
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

func fillParentRecursive(db *gorm.DB, dic *dict.Dictionary, allDic map[string]dict.Dictionary) error {
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
