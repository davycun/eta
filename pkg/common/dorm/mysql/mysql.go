package mysql

import (
	"database/sql"
	"fmt"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/go-sql-driver/mysql"
	ms "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/migrator"
	"gorm.io/gorm/schema"
)

type Migrator struct {
	ms.Migrator
	Dialector
}
type Dialector struct {
	ms.Dialector
}

func Open(dsn string) gorm.Dialector {
	dsnConf, _ := mysql.ParseDSN(dsn)
	dlt := ms.Dialector{Config: &ms.Config{DSN: dsn, DSNConfig: dsnConf}}
	return &Dialector{Dialector: dlt}
}

func (d Dialector) DataTypeOf(field *schema.Field) string {
	switch field.DataType {
	case schema.String, "varchar":
		return "TEXT"
	default:
		return d.Dialector.DataTypeOf(field)
	}
}

func (d Dialector) Migrator(db *gorm.DB) gorm.Migrator {
	mig := Migrator{
		Migrator: ms.Migrator{
			Migrator: migrator.Migrator{
				Config: migrator.Config{
					DB:        db,
					Dialector: d,
				},
			},
			Dialector: d.Dialector,
		},
		Dialector: d,
	}
	return mig
}

func EnsureDatabaseExists(database dorm.Database) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=True&loc=Local",
		database.User, database.Password, database.Host, database.Port)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logger.Warnf("errors: %v", err)
		return
	}
	defer db.Close()
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", database.DBName))
	if err != nil {
		logger.Warnf("errors: %v", err)
		return
	}
	db.Close()
	dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		database.User, database.Password, database.Host, database.Port, database.DBName)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		logger.Warnf("errors: %v", err)
		return
	}
	defer db.Close()
}
