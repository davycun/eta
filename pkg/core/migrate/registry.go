package migrate

import (
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/eta/constants"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"reflect"
	"strings"
	"sync"
)

//这个包在对表进行migrate的时候，提供了扩展功能，支持在表进行migrate的前后进行注册回调函数

const (
	CallbackForAll                  = "all_table" //如果通过这个注册migrate回调，那么表示任何表发生migrate行为都会触发当前回调
	CallbackBefore CallbackPosition = 1
	CallbackAfter  CallbackPosition = 1
)

var (
	callbackMap  = sync.Map{} //tableName ->  []Callback
	tbl          = reflect.TypeOf((*schema.Tabler)(nil)).Elem()
	tblWithNamer = reflect.TypeOf((*schema.TablerWithNamer)(nil)).Elem()
)

type CallbackPosition int
type Callback func(cfg *MigConfig, pos CallbackPosition) error

func AddCallback(tbName string, cb Callback) {
	if x, ok := callbackMap.Load(tbName); ok {
		cbList := x.([]Callback)
		cbList = append(cbList, cb)
		callbackMap.Store(tbName, cbList)
	} else {
		callbackMap.Store(tbName, []Callback{cb})
	}
}

type MigConfig struct {
	TxDB     *gorm.DB
	C        *ctx.Context
	TbOption entity.Table
}

func (mc *MigConfig) before() error {
	//TODO
	return nil
}
func (mc *MigConfig) after() error {
	return callAfterMigrate(mc, CallbackAfter)
}

func NewMigConfig(c *ctx.Context, db *gorm.DB, to entity.Table) *MigConfig {
	return &MigConfig{
		TxDB:     db,
		C:        c,
		TbOption: to,
	}
}

func callAfterMigrate(mc *MigConfig, pos CallbackPosition) (err error) {
	callbackList := getCallback(mc.TbOption.GetTableName())
	for _, fc := range callbackList {
		err = fc(mc, pos)
		if err != nil {
			return err
		}
	}
	return nil
}

func getCallback(tbName string) []Callback {
	mcs := make([]Callback, 0, 3)
	mc, ok := callbackMap.Load(CallbackForAll)
	if ok {
		tmp := mc.([]Callback)
		mcs = append(mcs, tmp...)
	}
	mc, ok = callbackMap.Load(tbName)
	if ok {
		tmp := mc.([]Callback)
		mcs = append(mcs, tmp...)
	}
	after, found := strings.CutPrefix(tbName, constants.TableTemplatePrefix)
	if found {
		mc, ok = callbackMap.Load(after)
		if ok {
			tmp := mc.([]Callback)
			mcs = append(mcs, tmp...)
		}
	}
	return mcs
}

func GetTableName(nm schema.Namer, dst interface{}) (schemaName, tableName string) {
	s, ok := dst.(string)
	if ok {
		return splitSchemaTable(s)
	}
	tpf := reflect.TypeOf(dst)
	if tpf.Implements(tbl) {
		table, ok1 := dst.(schema.Tabler)
		if ok1 {
			return splitSchemaTable(table.TableName())
		}
	}
	if tpf.Implements(tblWithNamer) {
		table, ok1 := dst.(schema.TablerWithNamer)
		if ok1 {
			return splitSchemaTable(table.TableName(nm))
		}
	}
	switch tpf.Kind() {
	case reflect.Slice, reflect.Map, reflect.Array, reflect.Chan, reflect.Pointer:
		return splitSchemaTable(strings.ToLower(tpf.Elem().Name()))
	default:
		return splitSchemaTable(strings.ToLower(tpf.Name()))
	}
}

func splitSchemaTable(s string) (schemaName, tableName string) {
	idx := strings.LastIndex(s, ".")
	if idx > -1 {
		return s[:idx], s[idx+1:]
	}
	return "", s
}
