package mock

import (
	"errors"
	"fmt"
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/eta/constants"
	"gorm.io/gorm"
)

const (
	MockParam = "MockParam"
)

type Mocker interface {
	iface.OrmService
	iface.ContextService
	Mock() (interface{}, error)
	CreateInBatches(value interface{}) error
	DeleteExistingData() error
	SetOrder(int)
	GetOrder() int
	SetConfig(*Config)
	GetConfig() *Config
	SetTableName(string)
	GetTableName() string
}

type DefaultMocker struct {
	Order     int
	Ctx       *ctx.Context
	DB        *gorm.DB
	Config    *Config
	TableName string
}

func (d *DefaultMocker) SetDB(orm *gorm.DB) {
	d.DB = orm
}
func (d *DefaultMocker) GetDB() *gorm.DB {
	return d.DB
}
func (d *DefaultMocker) GetContext() *ctx.Context {
	return d.Ctx
}
func (d *DefaultMocker) SetContext(c *ctx.Context) {
	d.Ctx = c
}
func (d *DefaultMocker) SetOrder(order int) {
	d.Order = order
}
func (d *DefaultMocker) GetOrder() int {
	return d.Order
}
func (d *DefaultMocker) SetConfig(config *Config) {
	d.Config = config
}
func (d *DefaultMocker) GetTableName() string {
	return d.TableName
}
func (d *DefaultMocker) SetTableName(tableName string) {
	d.TableName = tableName
}
func (d *DefaultMocker) GetConfig() *Config {
	return d.Config
}
func (d *DefaultMocker) Init(c *ctx.Context, db *gorm.DB, config *Config, order int, tableName string) {
	d.SetOrder(order)
	d.SetContext(c)
	d.SetDB(db)
	d.SetConfig(config)
	d.SetTableName(tableName)
	if config.CreateBatchSize <= 0 || config.CreateBatchSize > 2000 {
		config.CreateBatchSize = 500
	}
}
func (d *DefaultMocker) DeleteExistingData() error {
	dbType := dorm.GetDbType(d.GetDB())
	wh := ""
	switch dbType {
	case dorm.DaMeng:
		wh = fmt.Sprintf(`"%s".json_contains_any("extra",'mock','true',0)`, dorm.GetDbUser(d.DB))
	case dorm.PostgreSQL:
		wh = `"extra"->>'mock'='true'`
	case dorm.Mysql:
		wh = "`extra`->'$.mock'='true'"
	default:
		return errors.New("不支持的数据库类型")
	}
	return d.GetDB().Exec(fmt.Sprintf("delete from %s where %s", dorm.GetDbTable(d.GetDB(), d.GetTableName()), wh)).Error
}

// LoadRoom return roomId and buildingId
func (d *DefaultMocker) LoadRoom() ([]string, []string) {
	var (
		db         = d.GetDB()
		rs         []map[string]string
		roomId     = make([]string, 0, d.GetConfig().Size)
		buildingId = make([]string, 0, d.GetConfig().Size)
	)

	err := db.Exec(fmt.Sprintf("select %s,%s from %s limit %d",
		dorm.Quote(dorm.GetDbType(db), entity.IdDbName), dorm.Quote(dorm.GetDbType(db), "building_id"),
		dorm.GetDbTable(db, constants.TableRoom),
		d.GetConfig().Size)).Scan(&rs).Error
	if err != nil {
		logger.Errorf("loadRoom error: %v", err)
		return roomId, buildingId
	}
	for _, v := range rs {
		roomId = append(roomId, v[entity.IdDbName])
		buildingId = append(buildingId, v["building_id"])
	}
	return roomId, buildingId
}

// LoadFloor return floorId and buildingId
func (d *DefaultMocker) LoadFloor() ([]string, []string) {
	var (
		db         = d.GetDB()
		rs         []map[string]string
		roomId     = make([]string, 0, d.GetConfig().Size)
		buildingId = make([]string, 0, d.GetConfig().Size)
	)

	err := db.Exec(fmt.Sprintf("select %s,%s from %s limit %d",
		dorm.Quote(dorm.GetDbType(db), entity.IdDbName), dorm.Quote(dorm.GetDbType(db), "building_id"),
		dorm.GetDbTable(db, constants.TableFloor),
		d.GetConfig().Size)).Scan(&rs).Error
	if err != nil {
		logger.Errorf("loadRoom error: %v", err)
		return roomId, buildingId
	}
	for _, v := range rs {
		roomId = append(roomId, v[entity.IdDbName])
		buildingId = append(buildingId, v["building_id"])
	}
	return roomId, buildingId
}

func (d *DefaultMocker) CreateInBatches(value interface{}) error {
	return d.GetDB().CreateInBatches(value, d.GetConfig().CreateBatchSize).Error
}

func (d *DefaultMocker) CheckRelationData(table ...string) error {
	for _, t := range table {
		b := LoadContextData(d.GetContext(), "mock_"+t)
		if b == nil {
			return errors.New(fmt.Sprintf("关联表[%s]没有数据", t))
		}
	}

	return nil
}

type Service interface {
	iface.OrmService
	Mock(CitizenParam) error
	MockApp(CitizenAppParam) error
}

type DefaultService struct {
	Order int
	Param any
	DB    *gorm.DB
	Ctx   *ctx.Context
}

func (s *DefaultService) SetDB(orm *gorm.DB) {
	s.DB = orm
}
func (s *DefaultService) GetDB() *gorm.DB {
	return s.DB
}
func (s *DefaultService) GetContext() *ctx.Context {
	return s.Ctx
}
func (s *DefaultService) SetContext(c *ctx.Context) {
	s.Ctx = c
}
func (s *DefaultService) GetParam() any {
	return s.Param
}
func (s *DefaultService) SetParam(c any) {
	s.Param = c
}

func LoadContextData(c *ctx.Context, key string) any {
	d, ok := c.Get(key)
	if !ok {
		return nil
	}
	return d
}

func HasContextData(c *ctx.Context, key string) bool {
	_, ok := c.Get(key)
	return ok
}
