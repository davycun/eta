package doris

import "gorm.io/gorm"

func Update(db *gorm.DB, dst any) error {

	tx := db.Session(&gorm.Session{PrepareStmt: false, DryRun: true, SkipHooks: true}).Model(dst).Updates(dst)
	sq := db.Dialector.Explain(tx.Statement.SQL.String(), tx.Statement.Vars...)

	tx1 := db.Session(&gorm.Session{NewDB: true, PrepareStmt: false, AllowGlobalUpdate: true})
	return tx1.Raw(sq).Updates(dst).Error
}
