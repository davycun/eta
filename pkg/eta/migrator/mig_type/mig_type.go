package mig_type

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/dorm"
	"gorm.io/gorm"
)

func MigrateTypeAndFunction(db *gorm.DB) (err error) {

	return caller.NewCaller().
		Call(func(cl *caller.Caller) error { return migrateDb(db) }).
		Call(func(cl *caller.Caller) error { return migrateType(db) }).
		Call(func(cl *caller.Caller) error { return migrateFunction(db) }).
		Err
}

func migrateDb(db *gorm.DB) error {
	var (
		err    error
		dbType = dorm.GetDbType(db)
		scm    = dorm.GetDbSchema(db)
	)

	switch dbType {
	case dorm.PostgreSQL:
		//给一些trigger插入数据时用的id生成
		err = db.Exec(fmt.Sprintf(`CREATE SEQUENCE IF NOT EXISTS "%s"."%s"`, scm, dorm.SequenceIdName)).Error
	case dorm.DaMeng:
		//给一些trigger插入数据时用的id生成
		err = db.Exec(fmt.Sprintf(`CREATE SEQUENCE IF NOT EXISTS "%s"."%s"`, scm, dorm.SequenceIdName)).Error
		if err != nil {
			return err
		}
		err = initDm(db)
	case dorm.Mysql:
		tableName := dorm.GetScmTableName(db, "sequence")
		return caller.NewCaller().
			Call(func(cl *caller.Caller) error {
				sql := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s ( 
						name VARCHAR(255) NOT NULL, 
						current_value INT NOT NULL, 
						increment INT NOT NULL DEFAULT 1, 
						PRIMARY KEY (name) 
					)`, tableName)
				return db.Exec(sql).Error
			}).
			Call(func(cl *caller.Caller) error {
				return db.Exec(fmt.Sprintf("set global log_bin_trust_function_creators=1")).Error
			}).
			Call(func(cl *caller.Caller) error {
				tblSql := fmt.Sprintf(`
					CREATE FUNCTION IF NOT EXISTS %s (seq_name VARCHAR(255)) RETURNS INTEGER
					BEGIN
						DECLARE value INTEGER;
						SET value = 0; 
						SELECT current_value INTO value 
						FROM %s
						WHERE name = seq_name; 
						RETURN value; 
					END`, dorm.GetScmTableName(db, "currval"), tableName)
				return db.Exec(fmt.Sprintf(tblSql)).Error
			}).
			Call(func(cl *caller.Caller) error {
				tblSql := fmt.Sprintf(`
					CREATE FUNCTION IF NOT EXISTS %s (seq_name VARCHAR(255)) RETURNS INTEGER
					BEGIN
						UPDATE %s
						SET current_value = current_value + increment 
						WHERE name = seq_name;
						if ROW_COUNT()<=0 then
						  INSERT INTO %s VALUES (seq_name, 0, 1);
						end if;
						RETURN currval(seq_name); 
					END`, dorm.GetScmTableName(db, "nextval"), tableName, tableName)
				return db.Exec(fmt.Sprintf(tblSql)).Error
			}).
			Call(func(cl *caller.Caller) error {
				tblSql := fmt.Sprintf(`
					CREATE FUNCTION IF NOT EXISTS %s (seq_name VARCHAR(255), value INTEGER) RETURNS INTEGER
					BEGIN
						UPDATE %s
						SET current_value = value 
						WHERE name = seq_name; 
						if ROW_COUNT()<=0 then
						  INSERT INTO %s VALUES (seq_name, 0, 1);
						end if;
						RETURN currval(seq_name); 
					END`, dorm.GetScmTableName(db, "setval"), tableName, tableName)
				return db.Exec(fmt.Sprintf(tblSql)).Error
			}).Err
	}

	return err
}

func migrateType(db *gorm.DB) error {
	var (
		err    error
		dbType = dorm.GetDbType(db)
		dbUser = dorm.GetDbUser(db)
	)

	switch dbType {
	case dorm.DaMeng:
		var total int
		err = db.Raw(fmt.Sprintf(`select count(*) from DBA_SOURCE where NAME='%s' and OWNER='%s'`, "ARR_INT", dbUser)).Find(&total).Error
		if total > 0 {
			return nil
		}
		err = caller.NewCaller().
			Call(func(cl *caller.Caller) error {
				return db.Exec(fmt.Sprintf(`CREATE OR REPLACE TYPE %s.ARR_STR AS VARRAY(100000) OF VARCHAR2`, dbUser)).Error
			}).
			Call(func(cl *caller.Caller) error {
				return db.Exec(fmt.Sprintf(`CREATE OR REPLACE TYPE %s.ARR_INT AS VARRAY(100000) OF BIGINT`, dbUser)).Error
			}).
			Call(func(cl *caller.Caller) error {
				return db.Exec(fmt.Sprintf(`CREATE OR REPLACE CLASS %s.ARR_STR_CLS AS V %s.ARR_STR; end;`, dbUser, dbUser)).Error
			}).
			Call(func(cl *caller.Caller) error {
				return db.Exec(fmt.Sprintf(`CREATE OR REPLACE CLASS %s.ARR_INT_CLS AS V %s.ARR_INT; end;`, dbUser, dbUser)).Error
			}).
			Err

	case dorm.PostgreSQL:

	}
	return err
}

func migrateFunction(db *gorm.DB) error {
	var (
		err    error
		dbType = dorm.GetDbType(db)
	)

	switch dbType {
	case dorm.PostgreSQL:
		err = catArray(db)
	case dorm.DaMeng:
		err = caller.NewCaller().
			Call(func(cl *caller.Caller) error { return arrIntAnyBetween(db) }).
			Call(func(cl *caller.Caller) error { return arrContainsStr(db) }).
			Call(func(cl *caller.Caller) error { return arrContainsInt(db) }).
			Call(func(cl *caller.Caller) error { return jsonSize(db) }).
			Call(func(cl *caller.Caller) error { return jsonScalarContainsAny(db) }).
			Call(func(cl *caller.Caller) error { return jsonObjectContainsAny(db) }).
			Call(func(cl *caller.Caller) error { return jsonArrayContainsAny(db) }).
			Call(func(cl *caller.Caller) error { return jsonContainsAny(db) }).
			Err
	}
	return err
}
