package loader

import (
	"errors"
	"fmt"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/es"
	"github.com/davycun/eta/pkg/common/dorm/filter"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/core/builder"
	"github.com/davycun/eta/pkg/core/entity"
	"reflect"
	"time"

	"github.com/duke-git/lancet/v2/slice"
	jsoniter "github.com/json-iterator/go"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

type KeyLoader struct {
	Err    error
	Schema string
	DB     *gorm.DB
	DbType dorm.DbType
	Config KeyLoaderConfig
}

type KeyLoaderConfig struct {
	Keys      []string //where条件字段的值
	KeyColumn string   //where条件字段
	TableName string   // entity对应哪张表
	IndexName string   // entity对应哪个索引
}

func NewKeyLoader(db *gorm.DB, config KeyLoaderConfig) *KeyLoader {
	l := &KeyLoader{
		DB:     db,
		DbType: dorm.GetDbType(db),
		Schema: dorm.GetDbSchema(db),
		Config: config,
	}
	return l
}

func (l *KeyLoader) SetTableName(tableName string) *KeyLoader {
	l.Config.TableName = tableName
	return l
}

func (l *KeyLoader) AddKey(id ...string) *KeyLoader {
	if l.Config.Keys == nil {
		l.Config.Keys = make([]string, 0, 10)
	}
	l.Config.Keys = append(l.Config.Keys, id...)
	return l
}
func (l *KeyLoader) check() *KeyLoader {
	if l.Config.KeyColumn == "" {
		l.Config.KeyColumn = entity.IdDbName
	}

	if l.Config.TableName == "" {
		l.Err = errors.New("tableName is empty")
	}

	if len(l.Config.Keys) < 1 {
		l.Err = NoNeedLoadError
	}
	return l
}

func (l *KeyLoader) Load(rs any) error {

	if l.check().Err != nil {
		if errors.Is(l.Err, NoNeedLoadError) {
			logger.Errorf("%s, for %s", NoNeedLoadError, l.Config.TableName)
			l.Err = nil
		}
		return l.Err
	}

	var (
		dbType    = dorm.GetDbType(l.DB)
		_, tbName = dorm.Quote(dbType, l.Schema), dorm.Quote(dbType, l.Config.TableName)
		scmTbName = fmt.Sprintf("%s.%s", l.Schema, l.Config.TableName)
		cols      = dorm.JoinColumns(dbType, l.Config.TableName, []string{l.Config.KeyColumn})
		idCol     = dorm.Quote(dbType, l.Config.KeyColumn)
	)

	if len(l.Config.Keys) == 1 {
		l.Err = dorm.Table(l.DB, l.Config.TableName).Select(cols).
			Where(fmt.Sprintf(`%s = ?`, idCol), l.Config.Keys[0]).Find(rs).Error
		return l.Err
	}
	if len(l.Config.Keys) < 6 {
		l.Err = dorm.Table(l.DB, l.Config.TableName).Select(cols).
			Where(fmt.Sprintf(`%s in ?`, idCol), l.Config.Keys).Find(rs).Error
		return l.Err
	}

	//raw sql 需要自己包
	scmTbName = fmt.Sprintf("%s.%s", dorm.Quote(dbType, l.Schema), dorm.Quote(dbType, l.Config.TableName))
	rSql := builder.BuildValueToTableSql(l.DbType, true, l.Config.Keys...)
	sq := fmt.Sprintf(`with r as (%s) select %s from r, %s where r.%s=%s.%s `,
		rSql, cols, scmTbName, dorm.Quote(dbType, "id"), tbName, idCol)

	l.Err = dorm.RawFetch(sq, l.DB.
		Session(&gorm.Session{
			NewDB:  true,
			Logger: gormLogger.New(logger.Logger, gormLogger.Config{LogLevel: gormLogger.Warn, SlowThreshold: 30 * time.Second})}),
		rs)
	return l.Err
}

func (l *KeyLoader) LoadFromEs(rs any) error {
	if l.check().Err != nil {
		if errors.Is(l.Err, NoNeedLoadError) {
			logger.Errorf("%s, for %s", NoNeedLoadError, l.Config.TableName)
			l.Err = nil
		}
		return l.Err
	}

	mapToSlice := func(tmp []map[string]interface{}, rsObj any) error {
		tmpRsMap := make([]interface{}, 0)
		for _, m := range tmp {
			tmpRsMap = append(tmpRsMap, m[l.Config.KeyColumn])
		}
		mars, err := jsoniter.Marshal(tmpRsMap)
		if err != nil {
			return err
		}
		return jsoniter.Unmarshal(mars, rsObj)
	}

	var (
		maxIdSize = 1000
		esApi     = es.NewApi(global.GetES(), l.Config.IndexName).AddColumn(l.Config.KeyColumn)
	)
	if len(l.Config.Keys) <= maxIdSize {
		rsMap := make([]map[string]interface{}, 0)
		_, err := esApi.
			AddFilters(filter.Filter{Column: l.Config.KeyColumn, Operator: filter.IN, Value: l.Config.Keys}).
			Limit(len(l.Config.Keys)).
			LoadAll(true).
			Find(&rsMap)
		if err != nil {
			return err
		}
		return mapToSlice(rsMap, rs)
	}

	val := reflect.ValueOf(rs)
	isPrt := false
	if val.Kind() == reflect.Pointer {
		isPrt = true
		val = val.Elem()
	}
	if val.Kind() != reflect.Slice {
		return errors.New("data must be a slice type")
	}

	newRs := reflect.MakeSlice(val.Type(), 0, val.Len())
	chunks := slice.Chunk(l.Config.Keys, maxIdSize)
	for _, chunk := range chunks {
		var (
			tmpRsVal = reflect.MakeSlice(val.Type(), 0, val.Len())
			tmpRsPrt = reflect.New(tmpRsVal.Type())
		)
		tmpRsPrt.Elem().Set(tmpRsVal)

		rsMap := make([]map[string]interface{}, 0)
		_, err := esApi.
			AddFilters(filter.Filter{Column: l.Config.KeyColumn, Operator: filter.IN, Value: chunk}).
			Limit(len(chunk)).
			LoadAll(true).
			Find(&rsMap)

		if err != nil {
			return err
		}
		err = mapToSlice(rsMap, tmpRsPrt.Interface())
		if err != nil {
			return err
		}

		for i := range tmpRsPrt.Elem().Len() {
			newRs = reflect.Append(newRs, tmpRsPrt.Elem().Index(i))
		}
	}

	if isPrt {
		reflect.ValueOf(rs).Elem().Set(newRs)
	} else {
		reflect.ValueOf(rs).Set(newRs)
	}

	return nil
}
