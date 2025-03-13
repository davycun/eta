package doris

import (
	"github.com/davycun/eta/pkg/common/dorm"
	"gorm.io/gorm"
)

func Select(db *gorm.DB, dst any, f func(tx *gorm.DB) *gorm.DB) error {
	tx := db.Session(&gorm.Session{PrepareStmt: false, DryRun: true})
	tx = f(tx)
	sq := db.Dialector.Explain(tx.Statement.SQL.String(), tx.Statement.Vars...)
	tx = db.Session(&gorm.Session{NewDB: true, PrepareStmt: false})
	return dorm.RawFetch(sq, tx, dst)
}
