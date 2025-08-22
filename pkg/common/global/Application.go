package global

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/davycun/eta/pkg/common/cache"
	"github.com/davycun/eta/pkg/common/config"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/dgorm"
	"github.com/davycun/eta/pkg/common/dorm/es/es_api"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/common/id"
	"github.com/davycun/eta/pkg/common/id/snow"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"strings"
	"sync"
)

type Application struct {
	dbEngine     sync.Map
	appIdDbKey   sync.Map //存储apId -> dorm.Database
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
	db, err := dgorm.CreateGorm(database)

	if err != nil || db == nil {
		return nil, err
	}
	app.dbEngine.Store(database.GetKey(), db)
	return db, nil
}
func (app *Application) LoadGormByAppId(appId string) (*gorm.DB, error) {
	if appId == "" {
		return nil, errs.NewServerError("appId为空")
	}

	if dbCfg, ok := app.appIdDbKey.Load(appId); ok {
		return app.LoadGorm(dbCfg.(dorm.Database))
	}
	app.mutex.Lock()
	defer app.mutex.Unlock()

	if dbCfg, ok := app.appIdDbKey.Load(appId); ok {
		return app.LoadGorm(dbCfg.(dorm.Database))
	}

	var (
		dbStr = ""
		dbCfg dorm.Database
		wh    = map[string]interface{}{"id": appId}
	)

	var (
		localGorm = app.GetLocalGorm()
		dbType    = dorm.GetDbType(localGorm)
		scm       = app.GetLocalDatabase().Schema
		column    = "database"
	)
	err := localGorm.Table(dorm.Quote(dbType, scm, constants.TableApp)).Select(dorm.Quote(dbType, column)).Where(wh).Find(&dbStr).Error
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(dbStr), &dbCfg)
	if err != nil {
		return nil, err
	}

	db, err := app.LoadGorm(dbCfg)
	if err != nil || db == nil {
		return db, err
	}
	app.appIdDbKey.Store(appId, dbCfg)
	return db, nil
}

func (app *Application) DeleteGorm(database dorm.Database) {
	app.dbEngine.Delete(database.GetKey())
}
func (app *Application) DeleteGormByAppId(appId string) {
	if dbCfg, ok := app.appIdDbKey.Load(appId); ok {
		app.DeleteGorm(dbCfg.(dorm.Database))
	}
	app.appIdDbKey.Delete(appId)
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

func (app *Application) closeDb() {
	app.dbEngine.Range(func(key, value any) bool {
		db, ok := value.(*gorm.DB)
		if ok {
			dgorm.CloseGorm(db)
		}
		return true
	})
	logger.Infof("sql.CurDB closed")
}
