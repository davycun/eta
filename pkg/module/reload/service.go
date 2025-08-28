package reload

import (
	"errors"
	"fmt"
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/run"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/service"
	"github.com/davycun/eta/pkg/eta/ecf"
	"sync"
)

type (
	rdFunc func(rds *RdService) error
)

var (
	rdFuncMap = map[RdType]rdFunc{
		RdTypeDb2Es: db2es,
	}
)

func processReload(rdType RdType, c *ctx.Context, param *RdParamList, result *RdResultList) error {

	var (
		appDb   = c.GetAppGorm()
		localDb = global.GetLocalGorm()
	)

	rdsList := make([]*RdService, 0, len(param.Items))
	for _, item := range param.Items {
		itemCt := c.Clone()
		ec, ok := ecf.GetEntityConfig(appDb, item.TableName)
		if !ok {
			return errors.New("table_name not found")
		}
		ecf.SetContextEntityConfig(itemCt, &ec)
		if ec.LocatedApp() {
			itemCt.SetContextGorm(appDb)
		} else {
			itemCt.SetContextGorm(localDb)
		}
		srv, err1 := service.NewService(item.TableName, itemCt, itemCt.GetContextGorm())
		if err1 != nil {
			return err1
		}
		item.Param.OrderBy = []dorm.OrderBy{{Column: "eid", Asc: true}}
		rdsList = append(rdsList, &RdService{srv: srv, param: item.Param, result: &dto.Result{}})
	}

	var (
		err     error
		rdf, ok = rdFuncMap[rdType]
		wg      = sync.WaitGroup{}
	)
	if !ok {
		return errs.NewServerError(fmt.Sprintf("指定的realod func [%s]不存在", rdType))
	}

	for _, rds := range rdsList {
		if err != nil {
			return err
		}
		if param.Concurrent {
			wg.Add(1)
			run.Go(func() {
				defer wg.Done()
				if err != nil {
					return
				}
				err = errs.Cover(err, rdf(rds))
			})
		} else {
			err = rdf(rds)
		}
	}
	wg.Wait()

	return err
}
