package template_srv

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/core/history"
	"github.com/davycun/eta/pkg/core/ra"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/core/updater"
	"github.com/davycun/eta/pkg/module/data/template"
	"time"
)

func afterDelete(cfg *hook.SrvConfig, oldValues []template.Template) error {
	var (
		err    error
		dbType = dorm.GetDbType(cfg.TxDB)
	)
	for i, _ := range oldValues {
		var (
			execSql    = make([]string, 0, 20)
			p          = &oldValues[i]
			tbName     = p.GetTableName()
			targetName = fmt.Sprintf("%s_%s_deleted", tbName, time.Now().Format("20060102150405"))
		)
		execSql = append(execSql, fmt.Sprintf(`alter table %s rename to %s`, dorm.GetScmTableName(cfg.TxDB, tbName), dorm.Quote(dbType, targetName)))
		if p.Table.History.Data {
			//execSql = append(execSql, fmt.Sprintf(`drop table if exists "%s"."%s" cascade`, scm, p.HistoryTableName()))
			historyTargetName := fmt.Sprintf("%s_%s_deleted", p.HistoryTableName(), time.Now().Format("20060102150405"))
			execSql = append(execSql, fmt.Sprintf(`alter table %s rename to %s`,
				dorm.GetScmTableName(cfg.TxDB, p.HistoryTableName()), dorm.Quote(dbType, historyTargetName)))
			err = history.DropTrigger(cfg.TxDB, tbName)
			if err != nil {
				return err
			}
		}
		if ctype.Bool(p.Table.FieldUpdater) {
			err = updater.DropUpdaterTrigger(cfg.TxDB, tbName)
			if err != nil {
				return err
			}
		}
		if len(p.Table.RaDbFields) > 0 {
			err = ra.DropTrigger(cfg.TxDB, tbName)
			if err != nil {
				return err
			}
		}

		for _, v := range execSql {
			err = cfg.TxDB.Exec(v).Error
			if err != nil {
				return err
			}
		}
	}

	return err
}
