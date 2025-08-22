package global

import (
	"github.com/davycun/eta/pkg/common/config"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/es/es_api"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"time"
)

var (
	globalApp *Application
)

func GetApplication() *Application {
	return globalApp
}
func GetRedis() *redis.Client {
	return globalApp.GetRedis()
}
func GetGin() *gin.Engine {
	return globalApp.GetGin()
}
func GetValidator() *validator.Validate {
	if globalApp == nil {
		return binding.Validator.Engine().(*validator.Validate)
	}
	return globalApp.GetValidator()
}
func GetConfig() *config.Configuration {
	return globalApp.GetConfig()
}
func GetES() *es_api.Api {
	return globalApp.esApi
}

func LoadGorm(database dorm.Database) (*gorm.DB, error) {
	return globalApp.LoadGorm(database)
}
func LoadGormByAppId(appId string) (*gorm.DB, error) {
	return globalApp.LoadGormByAppId(appId)
}
func LoadGormSetAppId(appId string, database dorm.Database) (*gorm.DB, error) {
	db, err := globalApp.LoadGorm(database)
	if err != nil {
		return db, err
	}
	dorm.SetAppId(db, appId)
	return db, err
}
func DeleteGorm(database dorm.Database) {
	globalApp.DeleteGorm(database)
}
func DeleteGormByAppId(appId string) {

}
func GetLocalGorm() *gorm.DB {
	return globalApp.GetLocalGorm()
}
func GetLocalDoris() *gorm.DB {
	return globalApp.GetLocalDoris()
}
func GetLocalDatabase() dorm.Database {
	return globalApp.GetLocalDatabase()
}

func LoadAppDoris(appSchema string, logLevel, slowThreshold int) *gorm.DB {
	if GetLocalDoris() == nil {
		return nil
	}
	cfg := GetConfig().Doris
	//cfg.DBName = appSchema
	cfg.Schema = appSchema
	cfg.LogLevel = logLevel
	cfg.SlowThreshold = slowThreshold
	appDoris, err1 := LoadGorm(cfg)
	if err1 != nil {
		logger.Errorf("load doris db err %s", err1)
	}
	return appDoris
}

func NewLogger(logLevel int, slowThreshold int) gormLogger.Interface {
	return gormLogger.New(logger.Logger,
		gormLogger.Config{
			SlowThreshold:             time.Duration(slowThreshold) * time.Millisecond,
			LogLevel:                  gormLogger.LogLevel(logLevel),
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
			ParameterizedQueries:      false,
		})
}
