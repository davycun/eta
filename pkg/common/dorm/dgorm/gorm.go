package dgorm

import (
	"fmt"
	dameng "github.com/davycun/dm8-gorm"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/doris"
	"github.com/davycun/eta/pkg/common/dorm/mysql"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/utils"
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func CreateGorm(database dorm.Database) (*gorm.DB, error) {
	var (
		err error
	)
	err = binding.Validator.Engine().(*validator.Validate).Struct(database)
	if err != nil {
		return nil, err
	}
	slow := database.SlowThreshold
	logLv := database.LogLevel
	if logLv < 1 || logLv > 4 {
		logLv = 4
	}
	prepare := false // PrepareStmt=true开启的时候会内存泄露
	switch database.Type {
	case dorm.Doris.String(), dorm.Mysql.String():
		prepare = false
	}
	database.Host = utils.GetIP(database.Host)
	dl := GetDialect(database)
	conf := &gorm.Config{
		NamingStrategy: dorm.NewNamingStrategy(database),
		Logger: gormLogger.New(logger.Logger,
			gormLogger.Config{
				SlowThreshold:             time.Duration(slow) * time.Millisecond,
				LogLevel:                  gormLogger.LogLevel(logLv),
				IgnoreRecordNotFoundError: true,
				Colorful:                  true,
				ParameterizedQueries:      false,
			}),
		PrepareStmt: prepare,
	}
	dbe, err := gorm.Open(dl, conf)
	if err != nil {
		return dbe, err
	}
	//err = dbe.Use(plugin.NewCodecPlugin())
	for _, v := range gormPlugins {
		err = dbe.Use(v)
		if err != nil {
			return dbe, err
		}
	}
	db, err := dbe.DB()
	if err != nil {
		logger.Errorf("config db connection err %s", err)
		return nil, err
	}
	if database.MaxIdleCons > 0 {
		db.SetMaxIdleConns(database.MaxIdleCons)
	}
	if database.ConMaxLifetime > 0 {
		db.SetConnMaxLifetime(time.Second * time.Duration(database.ConMaxLifetime))
	}
	if database.ConMaxIdleTime > 0 {
		db.SetConnMaxIdleTime(time.Second * time.Duration(database.ConMaxIdleTime))
	}
	if database.MaxOpenCons > 0 {
		db.SetMaxOpenConns(database.MaxOpenCons)
	}
	return dbe, err
}

func CloseGorm(db *gorm.DB) {
	if db == nil {
		return
	}
	s, err := db.DB()
	if err != nil {
		logger.Errorf("close db err %s", err)
	}
	err = s.Close()
	if err != nil {
		logger.Errorf("close db err %s", err)
	}
	return
}

func GetDialect(database dorm.Database) gorm.Dialector {

	switch database.Type {
	case dorm.Doris.String():
		//interpolateParams=true -> Unsupported command(COM_STMT_SEND_LONG_DATA); Error 1047 (08S01): Unsupported command(COM_STMT_SEND_LONG_DATA)
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&interpolateParams=true",
			database.User, database.Password, database.Host, database.Port, database.DBName)
		dl := doris.Open(dsn)
		if x, ok := dl.(*doris.Dialector); ok {
			x.DefaultStringSize = 16382
		}
		return dl
	case dorm.Mysql.String():

		database.DBName = database.Schema
		mysql.EnsureDatabaseExists(database)
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			database.User, database.Password, database.Host, database.Port, database.Schema)
		dl := mysql.Open(dsn)
		if x, ok := dl.(*mysql.Dialector); ok {
			x.DefaultStringSize = 16382 // 65535(varchar最大字节数)/4(utf8mb4)=16383
		}
		return dl
	case dorm.PostgreSQL.String():
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=allow TimeZone=Asia/Shanghai", database.Host, database.User, database.Password, database.DBName, database.Port)
		return postgres.Open(dsn)
	case dorm.DaMeng.String():
		dsn := fmt.Sprintf("dm://%s:%s@%s:%d", database.User, database.Password, database.Host, database.Port)
		dl := dameng.Open(dsn)
		if x, ok := dl.(*dameng.Dialector); ok {
			x.DefaultStringSize = 8188
		}
		return dl
	default:
		logger.Errorf("unsupported dialect for %s", database.Type)
		panic(fmt.Sprintf("unsupported for dialect for %s", database.Type))
	}
}
