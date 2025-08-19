package history

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/dorm/db_table"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/eta/constants"
	"reflect"
	"sync"
	"time"

	dynamicstruct "github.com/ompluscator/dynamic-struct"
	"gorm.io/gorm"
)

var (
	historyTableCache = sync.Map{} //存储的是schemaTableName -> bool
)

type WrapperHistory struct {
	entity.History
	Entity any `json:"entity,omitempty" gorm:"embedded;embeddedPrefix:h_"`
}

func DeleteHistoryCallback(cfg *hook.SrvConfig, pos hook.CallbackPosition) error {

	return hook.AfterDelete(cfg, pos, func(cfg *hook.SrvConfig, oldValues []reflect.Value) error {
		var (
			userId = cfg.Ctx.GetContextUserId()
			deptId = cfg.Ctx.GetContextCurrentDeptId()
		)
		return CreateDeleteHistory(cfg.TxDB, cfg.GetTableName(), NewWrapperHistory(userId, deptId, oldValues))
	})
}

// NewWrapperHistory
// 返回一个History 包含Entity字段的动态结构体切片指针
func NewWrapperHistory[T any](userId, deptId string, values []T) any {
	var (
		bd = dynamicstruct.ExtendStruct(entity.History{})
	)

	if len(values) < 1 {
		return nil
	}
	switch x := any(values[0]).(type) {
	case reflect.Value:
		if x.CanInterface() {
			bd.AddField("Entity", reflect.New(x.Type()).Elem().Interface(), `json:"entity,omitempty" gorm:"embedded;embeddedPrefix:h_"`)
		}
	default:
		bd.AddField("Entity", x, `json:"entity,omitempty" gorm:"embedded;embeddedPrefix:h_"`)
	}

	st := bd.Build()
	dataList := st.NewSliceOfStructs()
	val := reflect.ValueOf(dataList)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}

	for _, v := range values {
		data := st.New()
		tmpVal := reflect.ValueOf(data)
		if tmpVal.Kind() == reflect.Pointer {
			tmpVal = tmpVal.Elem()
		}
		tmpVal.FieldByName("ID").Set(reflect.ValueOf(global.GenerateIDStr()))
		tmpVal.FieldByName("CreatedAt").Set(reflect.ValueOf(ctype.NewLocalTimePrt(time.Now())))
		tmpVal.FieldByName("OpType").Set(reflect.ValueOf(entity.HistoryDelete))
		tmpVal.FieldByName("OptUserId").Set(reflect.ValueOf(userId))
		tmpVal.FieldByName("OptDeptId").Set(reflect.ValueOf(deptId))

		switch x := any(v).(type) {
		case reflect.Value:
			if x.CanInterface() {
				tmpVal.FieldByName("Entity").Set(x)
			}
		default:
			tmpVal.FieldByName("Entity").Set(reflect.ValueOf(x))
		}

		val = reflect.Append(val, tmpVal)
	}
	return val.Interface()
}

// CreateDeleteHistory
// hisList 需要是一个切片
func CreateDeleteHistory(txDb *gorm.DB, tableName string, hisList any) error {

	if hisList == nil {
		return nil
	}

	var (
		scm       = dorm.GetDbSchema(txDb)
		scmTbName = fmt.Sprintf("%s.%s", scm, tableName)
	)

	//判断是否存在历史表，存在则插入删除记录，不存在则直接返回
	//这里针对Template 动态删除了历史记录表的情况是会存在问题的。
	val, ok := historyTableCache.Load(scmTbName)
	if !ok {
		exists := db_table.TableExists(txDb, scm, tableName+constants.TableHistorySubFix)
		historyTableCache.Store(scmTbName, exists)
		if !exists {
			logger.Infof("不需要记录删除历史，因为表[%s]没有开启历史记录", scmTbName)
			return nil
		}
	} else {
		if !ctype.Bool(val) {
			logger.Infof("不需要记录删除历史，因为表[%s]没有开启历史记录", scmTbName)
			return nil
		}
	}
	return dorm.Table(txDb, tableName+constants.TableHistorySubFix).Create(hisList).Error
}
