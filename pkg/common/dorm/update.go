package dorm

import (
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
	"gorm.io/gorm/clause"
	gormLogger "gorm.io/gorm/logger"
)

var (
	silentLogger = gormLogger.New(logger.Logger, gormLogger.Config{LogLevel: gormLogger.Silent})
)

//TODO
//采用clause.OnConflict 重新处理下关联的字段，以结构体重新设计功能

func BatchUpdate(db *gorm.DB, dest any, cols ...string) error {
	var (
		dbType = GetDbType(db)
	)
	tx := db.Session(&gorm.Session{NewDB: true, SkipHooks: true, DryRun: true, Logger: silentLogger}).Model(dest)
	if len(cols) > 0 {
		cols = utils.Merge(cols, "id")
		tx = tx.Select(cols)
	}
	tx = tx.Callback().Create().Execute(tx)
	values := callbacks.ConvertToCreateValues(tx.Statement)

	cs := cols
	if len(cs) < 1 {
		cs = make([]string, 0, len(values.Columns))
		for _, v := range values.Columns {
			cs = append(cs, v.Name)
		}
	}

	switch dbType {
	case DaMeng:
		dmBatchUpdate(tx, values, cs...)
		sq := db.Dialector.Explain(tx.Statement.SQL.String(), tx.Statement.Vars...)
		return db.Exec(sq).Error
	case PostgreSQL:
		pgBatchUpdate(tx, values, cs...)
		sq := db.Dialector.Explain(tx.Statement.SQL.String(), tx.Statement.Vars...)
		return db.Exec(sq).Error
	case Mysql:
		cfl := clause.OnConflict{
			Columns: []clause.Column{
				{Name: "id"},
			},
			DoUpdates: clause.AssignmentColumns(cols),
		}
		return db.Model(dest).Clauses(cfl).Create(dest).Error
	}
	return nil
}
