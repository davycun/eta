package service

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/core/service/sqlbd"
)

func (s *DefaultService) RetrieveFromEs(cfg *hook.SrvConfig, sqlList *sqlbd.SqlList) error {

	if sqlList == nil {
		return errs.NewServerError(fmt.Sprintf("[%s:%s]没有指定RetrieveFromEs函数", cfg.GetTableName(), cfg.Method))
	}

	var (
		esRetriever = sqlList.GetEsRetriever()
	)
	if esRetriever != nil {
		return esRetriever(cfg, sqlList)
	}
	return errs.NewServerError(fmt.Sprintf("[%s:%s]没有指定RetrieveFromEs函数", cfg.GetTableName(), cfg.Method))
}
