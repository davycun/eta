package dorm

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/davycun/dm8-gorm"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/utils"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"strings"
	"time"
)

type ExecutedSql struct {
	sql []string
}

func (b *ExecutedSql) Printf(s string, args ...interface{}) {
	if len(args) > 0 {
		b.sql = append(b.sql, strings.TrimSpace(fmt.Sprintf("%s", args[len(args)-1])))
	}
	logger.Logger.Printf(s, args...)
}

func (b *ExecutedSql) Exists(rawSql string) bool {
	return utils.ContainAny(b.sql, rawSql)
}
func (b *ExecutedSql) Reset() {
	b.sql = b.sql[0:0]
}

func NewTestDB(dbType DbType, schema string) (*gorm.DB, *ExecutedSql, sqlmock.Sqlmock, error) {
	var (
		rawSq    = &ExecutedSql{}
		database = Database{
			Host:          "127.0.0.1",
			Port:          1234,
			DBName:        "eta",
			Schema:        schema,
			User:          "test",
			Password:      "test",
			Type:          dbType.String(),
			LogLevel:      4,
			SlowThreshold: 200,
		}
	)

	lg := gormLogger.New(rawSq,
		gormLogger.Config{
			SlowThreshold:             time.Duration(database.SlowThreshold) * time.Millisecond,
			LogLevel:                  gormLogger.LogLevel(database.LogLevel),
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
			ParameterizedQueries:      false,
		})

	conf := &gorm.Config{
		NamingStrategy:       NewNamingStrategy(database),
		Logger:               lg,
		PrepareStmt:          false,
		DisableAutomaticPing: true,
	}

	db, mock, err := sqlmock.New()

	if err != nil {
		return nil, rawSq, nil, err
	}

	var dialect gorm.Dialector
	switch dbType {
	case DaMeng:
		dialect = dmgorm.New(dmgorm.Config{
			Conn: db,
		})
	case PostgreSQL:
		dialect = postgres.New(postgres.Config{
			Conn: db,
		})
	case Mysql:
		dialect = mysql.New(mysql.Config{
			Conn:                      db,
			SkipInitializeWithVersion: true,
		})
	}

	dbe, err := gorm.Open(dialect, conf)
	return dbe.Session(&gorm.Session{DryRun: true, PrepareStmt: false}), rawSq, mock, err
}
