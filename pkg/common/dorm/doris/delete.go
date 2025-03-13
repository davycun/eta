package doris

import (
	"gorm.io/gorm"
)

// Delete
// RowsAffected 为0，返回没有意义
func Delete(db *gorm.DB, f func(tx *gorm.DB) *gorm.DB) error {
	tx := db.Session(&gorm.Session{PrepareStmt: false, DryRun: true})
	tx = f(tx)
	sq := db.Dialector.Explain(tx.Statement.SQL.String(), tx.Statement.Vars...)

	tx1 := db.Session(&gorm.Session{NewDB: true, PrepareStmt: false, AllowGlobalUpdate: true})
	return tx1.Raw(sq).Delete(tx.Statement.Dest).Error
}
