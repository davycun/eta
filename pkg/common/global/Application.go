package global

import (
	"context"
	"fmt"
	"github.com/davycun/dm8-gorm"
	"github.com/davycun/eta/pkg/common/cache"
	"github.com/davycun/eta/pkg/common/config"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/doris"
	"github.com/davycun/eta/pkg/common/dorm/es/es_api"
	"github.com/davycun/eta/pkg/common/dorm/mysql"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/common/id"
	"github.com/davycun/eta/pkg/common/id/snow"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"strings"
	"sync"
	"time"
)

type Application struct {
	dbEngine     sync.Map
	nebulaEngine sync.Map
	gin          *gin.Engine
	validator    *validator.Validate
	config       *config.Configuration
	redis        *redis.Client
	idGenerator  id.Generator
	dorisEngin   sync.Map
	mutex        sync.Mutex
	esApi        *es_api.Api
}

func NewApplication(cfg *config.Configuration) *Application {
	logger.Info("Application init...")
	ap := &Application{config: cfg}
	ap.validator = binding.Validator.Engine().(*validator.Validate)
	ap.initGorm()
	ap.initRedis()
	ap.initServer()
	ap.initIdGenerator()
	ap.initDoris()
	ap.initElasticsearch()
	return ap
}

// LoadGorm 当前应用的使用的数据
func (app *Application) LoadGorm(database dorm.Database) (*gorm.DB, error) {
	if database.Host == "" || database.Port == 0 || database.DBName == "" {
		return nil, errs.NewServerError("数据库信息不全")
	}
	value, ok := app.dbEngine.Load(database.GetKey())
	if ok {
		return value.(*gorm.DB), nil
	}
	app.mutex.Lock()
	defer app.mutex.Unlock()
	value, ok = app.dbEngine.Load(database.GetKey())
	if ok {
		return value.(*gorm.DB), nil
	}
	db, err := app.createGorm(database)

	if err != nil || db == nil {
		return nil, err
	}
	app.dbEngine.Store(database.GetKey(), db)
	return db, nil
}

func (app *Application) DeleteGorm(database dorm.Database) {
	app.dbEngine.Delete(database.GetKey())
}
func (app *Application) GetLocalGorm() *gorm.DB {
	orm, err := app.LoadGorm(app.config.Database)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return orm
}
func (app *Application) GetLocalDatabase() dorm.Database {
	return app.config.Database
}

func (app *Application) GetLocalNebulaConfig() config.Nebula {
	return app.config.Nebula
}

func (app *Application) GetRedis() *redis.Client {
	return app.redis
}
func (app *Application) GetLocalDoris() *gorm.DB {
	orm, err := app.LoadGorm(app.config.Doris)
	if err != nil {
		return nil
	}
	return orm
}
func (app *Application) GetGin() *gin.Engine {
	return app.gin
}
func (app *Application) GetValidator() *validator.Validate {
	return app.validator
}
func (app *Application) GetConfig() *config.Configuration {
	return app.config
}
func (app *Application) GetIdGenerator() id.Generator {
	return app.idGenerator
}
func (app *Application) Shutdown(ctx context.Context) error {
	var (
		err error
	)
	err = app.redis.Close()
	logger.Infof("redis closed")
	app.closeDb()
	return err
}

func (app *Application) initRedis() {
	//TODO need support ClusterClient decide by app.config
	app.redis = redis.NewClient(&redis.Options{
		Addr:     app.config.Redis.Addr(),
		Username: app.config.Redis.Username,
		Password: app.config.Redis.Password,
		DB:       app.config.Redis.DB,
	})

	//如果想去除global对cache包的依赖，可以把下面的调用迁移到的server包
	cache.InitCacheWithRedis(app.redis)

}
func (app *Application) initServer() {
	err := app.validator.Struct(app.config.Server)
	utils.AssertNil(err)
	if app.config.Server.GinMode != "" {
		gin.SetMode(app.config.Server.GinMode)
	}
	app.gin = gin.New()

}
func (app *Application) initGorm() {

	if app.config.Database.IsEmpty() {
		logger.Errorf("database config is empty")
		return
	}
	_, err := app.LoadGorm(app.config.Database)
	if err != nil {
		logger.Error(err)
	}
}
func (app *Application) initIdGenerator() {
	idg := app.config.ID
	if strings.ToLower(idg.Type) == config.IdGeneratorTypeRedis {
		ir := app.redis.IncrBy(context.Background(), config.IdGeneratorRedisKey, 1)
		nodeId := ir.Val()
		if fmt.Sprintf("%s", ir.Err()) == "ERR increment or decrement would overflow" {
			ir1 := app.redis.Set(context.Background(), config.IdGeneratorRedisKey, 1, 0)
			if ir1.Err() != nil {
				logger.Errorf("set %s error: %s", config.IdGeneratorRedisKey, ir1.Err())
			}
			nodeId = 1
		}
		app.idGenerator = snow.NewSnowflake(nodeId, idg.Epoch)
	} else {
		app.idGenerator = snow.NewSnowflake(idg.NodeId, idg.Epoch)
	}
}
func (app *Application) initDoris() {
	if app.config.Doris.IsEmpty() {
		logger.Errorf("doris config is empty")
		return
	}
	if app.config.Doris.IsEmpty() {
		logger.Errorf("database config is empty")
		return
	}
	_, err := app.LoadGorm(app.config.Doris)
	if err != nil {
		logger.Error(err)
	}
}

func (app *Application) initElasticsearch() {
	app.esApi = es_api.New(app.config.EsConfig)
}

func (app *Application) createGorm(database dorm.Database) (*gorm.DB, error) {
	var (
		err error
	)
	err = app.validator.Struct(database)
	if err != nil {
		return nil, err
	}
	slow := IfLtDefault(database.SlowThreshold, 200, 500)
	logLv := database.LogLevel
	if logLv < 1 || logLv > 4 {
		logLv = 4
	}
	prepare := true
	switch database.Type {
	case dorm.Doris.String(), dorm.Mysql.String():
		prepare = false
	}
	database.Host = utils.GetIP(database.Host)
	dl := app.getDialect(database)
	conf := &gorm.Config{
		NamingStrategy: dorm.NamingStrategy{Config: database},
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
	configConnPool(dbe, database)
	return dbe, err
}

func configConnPool(dbe *gorm.DB, dbCfg dorm.Database) {
	db, err := dbe.DB()
	if err != nil {
		logger.Errorf("config db connection err %s", err)
		return
	}
	if dbCfg.MaxIdleCons > 0 {
		db.SetMaxIdleConns(dbCfg.MaxIdleCons)
	}
	if dbCfg.ConMaxLifetime > 0 {
		db.SetConnMaxLifetime(time.Second * time.Duration(dbCfg.ConMaxLifetime))
	}
	if dbCfg.ConMaxIdleTime > 0 {
		db.SetConnMaxIdleTime(time.Second * time.Duration(dbCfg.ConMaxIdleTime))
	}
	if dbCfg.MaxOpenCons > 0 {
		db.SetMaxOpenConns(dbCfg.MaxOpenCons)
	}
}

func IfLtDefault(src, target, dft int) int {
	if src < target {
		return dft
	}
	return src
}
func IfGtDefault(src, target, dft int) int {
	if src > target {
		return dft
	}
	return src
}
func (app *Application) getDialect(database dorm.Database) gorm.Dialector {

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
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai", database.Host, database.User, database.Password, database.DBName, database.Port)
		return postgres.Open(dsn)
	case dorm.DaMeng.String():
		dsn := fmt.Sprintf("dm://%s:%s@%s:%d", database.User, database.Password, database.Host, database.Port)
		dl := dmgorm.Open(dsn)
		if x, ok := dl.(*dmgorm.Dialector); ok {
			x.DefaultStringSize = 8188
		}
		return dl
	default:
		logger.Errorf("unsupported dialect for %s", database.Type)
		panic(fmt.Sprintf("unsupported for dialect for %s", database.Type))
	}
}

func (app *Application) closeDb() {
	app.dbEngine.Range(func(key, value any) bool {
		db, ok := value.(*gorm.DB)
		if ok {
			s, err := db.DB()
			err = s.Close()
			if err != nil {
				logger.Errorf("close db err %s", err)
			}
		}
		return true
	})
	logger.Infof("sql.CurDB closed")
}
