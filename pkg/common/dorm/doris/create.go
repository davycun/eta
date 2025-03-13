package doris

import "gorm.io/gorm"

func Upsert(db *gorm.DB, values interface{}) error {
	if db == nil || values == nil {
		return nil
	}
	//tx := db.Session(&gorm.Session{NewDB: true, PrepareStmt: false})
	//tx = tx.Exec("set enable_insert_strict=false").Exec("set enable_unique_key_partial_update=true")
	return db.Model(values).Create(values).Error
}
