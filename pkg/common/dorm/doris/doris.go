package doris

import (
	mysql2 "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/migrator"
	"gorm.io/gorm/schema"
)

type Config struct {
	mysql.Config
}

type Dialector struct {
	mysql.Dialector
}

func Open(dsn string) gorm.Dialector {
	dsnConf, _ := mysql2.ParseDSN(dsn)
	dlt := mysql.Dialector{Config: &mysql.Config{DSN: dsn, DSNConfig: dsnConf}}
	return &Dialector{Dialector: dlt}
}

func New(config Config) gorm.Dialector {

	switch {
	case config.DSN == "" && config.DSNConfig != nil:
		config.DSN = config.DSNConfig.FormatDSN()
	case config.DSN != "" && config.DSNConfig == nil:
		config.DSNConfig, _ = mysql2.ParseDSN(config.DSN)
	}
	dlt := mysql.Dialector{Config: &config.Config}

	return &Dialector{Dialector: dlt}
}

func (dialector Dialector) DataTypeOf(field *schema.Field) string {
	switch field.DataType {
	case schema.Bool, schema.Int, schema.Uint, schema.Float, schema.Time, schema.Bytes:
		return dialector.Dialector.DataTypeOf(field)
	case schema.String:
		return "VARCHAR"
	default:
		return dialector.Dialector.DataTypeOf(field)
	}
}

func (dialector Dialector) Migrator(db *gorm.DB) gorm.Migrator {
	mig := Migrator{
		Migrator: mysql.Migrator{
			Migrator: migrator.Migrator{
				Config: migrator.Config{
					DB:        db,
					Dialector: dialector,
				},
			},
			Dialector: dialector.Dialector,
		},
		Dialector: dialector,
	}
	return mig
}
