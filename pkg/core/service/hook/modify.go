package hook

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/dorm/filter"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/iface"
	"gorm.io/gorm"
	"reflect"
	"strings"
	"time"
)

func (cfg *SrvConfig) modifyBefore(callbacks ...CallbackWrapper) error {

	var (
		err error
	)

	err = caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			return cfg.modifyCheck()
		}).
		Call(func(cl *caller.Caller) error {
			return cfg.modifyBeforeLoadValues(callbacks...)
		}).
		Call(func(cl *caller.Caller) error {
			cfg.modifyFillBefore()
			return callAuthCallback(cfg, CallbackBefore)
		}).
		Call(func(cl *caller.Caller) error {
			return callCallback(cfg, iface.CurdModify, cfg.Method, CallbackBefore, callbacks...)
		}).Err
	return err
}
func (cfg *SrvConfig) modifyAfter(callbacks ...CallbackWrapper) error {
	err := callCallback(cfg, iface.CurdModify, cfg.Method, CallbackAfter, callbacks...)
	cfg.modifyFillAfter()
	return err
}
func (cfg *SrvConfig) modifyCheck() error {
	if cfg.Param == nil || cfg.Param.Data == nil {
		logger.Errorf("cfg.Param.Data is null, cfg.Values will be null, tableName is %s, CurdType is %s", cfg.GetTableName(), cfg.CurdType)
	}
	if cfg.Result == nil {
		logger.Error("the ModifyConfig Result is nil")
	}
	if cfg.Ctx == nil {
		logger.Error("the ModifyConfig Context is nil")
	}
	return nil
}

func (cfg *SrvConfig) modifyBeforeLoadValues(modifyCallbacks ...CallbackWrapper) error {

	var (
		err    error
		dbType = dorm.GetDbType(cfg.OriginDB)
	)
	if cfg.CurdType != iface.CurdModify {
		return nil
	}

	//如果没有回调，就没必要加载数据
	if len(modifyCallbacks) < 1 && len(getModifyExtendCallback(cfg.GetTableName())) < 1 {
		return err
	}

	entitySlice := cfg.NewEntitySlicePointer()

	switch cfg.Method {
	case iface.MethodCreate:
		if cfg.Param != nil {
			cfg.NewValues = cfg.Param.Data
		}
	case iface.MethodUpdate:
		if cfg.Param != nil {
			cfg.NewValues = cfg.Param.Data
		}
		cfg.OldValues = entitySlice
		err = cfg.loadDataByIdAndUpdatedAt(cfg.OriginDB, cfg.Values, cfg.OldValues)
	case iface.MethodUpdateByFilters:
		if cfg.Param != nil {
			cfg.NewValues = cfg.Param.Data
		}
		wh := filter.ResolveWhereTable(cfg.GetTableName(), cfg.Param.Filters, dbType)
		if wh == "" {
			return errs.NewClientError("批量更新条件不能为空")
		}
		cfg.OldValues = entitySlice
		tx := dorm.Table(cfg.OriginDB, cfg.GetTableName())
		err = tx.Where(wh).Find(cfg.OldValues).Error
	case iface.MethodDelete:
		cfg.OldValues = entitySlice
		err = cfg.loadDataByIdAndUpdatedAt(cfg.OriginDB, cfg.Values, cfg.OldValues)
	case iface.MethodDeleteByFilters:
		cfg.OldValues = entitySlice
		wh := filter.ResolveWhereTable(cfg.GetTableName(), cfg.Param.Filters, dbType)
		if wh == "" {
			return errs.NewClientError("批量删除条件不能为空")
		}
		tx := dorm.Table(cfg.OriginDB, cfg.GetTableName())
		err = tx.Where(wh).Find(cfg.OldValues).Error
	default:

	}
	if err != nil {
		return err
	}
	if !cfg.checkContinue() {
		return errs.NoRecordAffected
	}

	return err
}

// 如果是更新或者删除，那么OldValues需要匹配到数据，否则即使无需更新或者删除
func (cfg *SrvConfig) checkContinue() bool {
	switch cfg.Method {
	case iface.MethodUpdate, iface.MethodUpdateByFilters, iface.MethodDelete, iface.MethodDeleteByFilters:
		if cfg.OldValues == nil {
			return false
		}
		val := utils.GetRealValue(reflect.ValueOf(cfg.OldValues))
		if val.IsValid() && val.Kind() == reflect.Slice {
			if val.Len() > 0 {
				return true
			} else {
				return false
			}
		}
	default:
	}

	return true
}
func (cfg *SrvConfig) loadDataByIdAndUpdatedAt(db *gorm.DB, dataValues []reflect.Value, target any) error {

	var (
		dbType = dorm.GetDbType(db)
		sqs    = make([]string, 0, len(dataValues))
	)
	for _, v := range dataValues {
		id := v.FieldByName(entity.IdFieldName).String()
		updatedAt := v.FieldByName(entity.UpdatedAtFieldName).Int()
		sqs = append(sqs, fmt.Sprintf(`(%s='%s' and %s=%d)`, dorm.Quote(dbType, "id"), id, dorm.Quote(dbType, "updated_at"), updatedAt))
	}
	if len(sqs) < 1 {
		return nil
	}
	return dorm.Table(db, cfg.GetTableName()).Where(strings.Join(sqs, " or ")).Find(target).Error
}

func (cfg *SrvConfig) modifyFillBefore() {

	dbType := dorm.GetDbType(cfg.CurDB)
	for _, val := range cfg.Values {
		switch cfg.Method {
		case iface.MethodCreate:
			initId(dbType, val)
			initCreatedAt(val)
			initUpdatedAt(val)
			initCreatorId(val, cfg.Ctx)
			initUpdaterId(val, cfg.Ctx)
			initCreatorDeptId(val, cfg.Ctx)
			initUpdaterDeptId(val, cfg.Ctx)
			initFieldUpdaterIds(val, cfg.Ctx)
			initRaContent(val, cfg.Ctx)
		case iface.MethodUpdate, iface.MethodUpdateByFilters:
			initUpdatedAt(val)
			initUpdaterId(val, cfg.Ctx)
			initUpdaterDeptId(val, cfg.Ctx)
		default:
		}
	}
}
func (cfg *SrvConfig) modifyFillAfter() {

	switch cfg.Method {
	case iface.MethodCreate:
		rs := make([]ctype.Map, 0, len(cfg.Values))
		for _, v := range cfg.Values {
			vId := v.FieldByName(entity.IdFieldName)
			vUpdateAt := v.FieldByName(entity.UpdatedAtFieldName)
			dt := ctype.Map{
				entity.UpdatedAtDbName: vUpdateAt.Interface(),
				entity.IdDbName:        vId.Interface(),
			}
			rs = append(rs, dt)
		}
		cfg.Result.Data = rs
	default:

	}
}

func (cfg *SrvConfig) callExtendModifyCallback(pos CallbackPosition) error {
	var (
		err error
		mcs = getModifyExtendCallback(cfg.GetTableName())
	)
	for _, fc := range mcs {
		err = fc.Callback(cfg, pos)
		if err != nil {
			return err
		}
	}
	return err
}

func ResolveValue(val reflect.Value) []reflect.Value {
	var vs []reflect.Value
	if !val.IsValid() {
		return vs
	}
	switch val.Kind() {
	case reflect.Pointer:
		return utils.ConvertToValueArray(val.Elem())
	case reflect.Slice:
		vs = make([]reflect.Value, 0, val.Len())
		for i := 0; i < val.Len(); i++ {
			v := val.Index(i)
			vs = append(vs, v)
		}
		return vs
	case reflect.Struct:
		vs = make([]reflect.Value, 1)
		vs[0] = val
		return vs
	default:
	}
	return vs
}
func GenerateRaContent(val reflect.Value) {
	if val.IsValid() && val.CanInterface() {
		obj := val.Interface()
		fields := entity.GetRaDbFields(val.Interface())
		if len(fields) <= 0 {
			return
		}
		strs := make([]string, 0, len(fields))
		for _, v := range fields {
			strs = append(strs, entity.GetString(obj, v))
		}

		raVal := val.FieldByName(entity.RaContentFieldName)
		if raVal.IsValid() && raVal.IsZero() && raVal.CanSet() {
			raVal.Set(reflect.ValueOf(ctype.NewTextPrt(strings.Join(strs, " "))))
		}
	}
}
func initId(dbType dorm.DbType, val reflect.Value) {
	id := val.FieldByName(entity.IdFieldName)
	if id.IsValid() && id.IsZero() {
		id.SetString(global.GenerateIDStr())
	}
}
func initCreatedAt(val reflect.Value) {
	createdAt := val.FieldByName(entity.CreatedAtFieldName)
	if !createdAt.IsValid() {
		return
	}
	if createdAt.IsZero() {
		i := createdAt.Interface()
		switch i.(type) {
		case *ctype.LocalTime:
			now := ctype.LocalTime{Data: time.Now(), Valid: true}
			createdAt.Set(reflect.ValueOf(&now))
		case ctype.LocalTime:
			now := ctype.LocalTime{Data: time.Now(), Valid: true}
			createdAt.Set(reflect.ValueOf(now))
		case time.Time:
			createdAt.Set(reflect.ValueOf(time.Now()))
		case *time.Time:
			now := time.Now()
			createdAt.Set(reflect.ValueOf(&now))
		}
	}
}
func initUpdatedAt(val reflect.Value) {
	updatedAt := val.FieldByName(entity.UpdatedAtFieldName)
	if updatedAt.IsValid() && updatedAt.IsZero() {
		updatedAt.SetInt(time.Now().UnixMilli())
	}
}
func initCreatorId(val reflect.Value, c *ctx.Context) {
	uid := c.GetContextUserId()
	createUser := val.FieldByName(entity.CreatorIdFieldName)
	if createUser.IsValid() && createUser.IsZero() {
		createUser.SetString(uid)
	}
}
func initUpdaterId(val reflect.Value, c *ctx.Context) {
	uid := c.GetContextUserId()
	updateUser := val.FieldByName(entity.UpdaterIdFieldName)
	if updateUser.IsValid() && updateUser.IsZero() {
		updateUser.SetString(uid)
	}
}
func initCreatorDeptId(val reflect.Value, c *ctx.Context) {
	deptId := c.GetContextCurrentDeptId()
	if deptId == "" {
		deptId = c.GetContextUserId()
	}
	createDeptId := val.FieldByName(entity.CreatorDeptIdFieldName)
	if createDeptId.IsValid() && createDeptId.IsZero() {
		createDeptId.SetString(deptId)
	}
}
func initUpdaterDeptId(val reflect.Value, c *ctx.Context) {
	deptId := c.GetContextCurrentDeptId()
	if deptId == "" {
		deptId = c.GetContextUserId()
	}
	updaterDeptId := val.FieldByName(entity.UpdaterDeptIdFieldName)
	if updaterDeptId.IsValid() && updaterDeptId.IsZero() {
		updaterDeptId.SetString(deptId)
	}
}
func initFieldUpdaterIds(val reflect.Value, c *ctx.Context) {

	updaterId := c.GetContextUserId()
	ids := val.FieldByName(entity.FieldUpdaterIdsFieldName)
	if ids.IsValid() && ids.IsZero() {
		ids.Set(reflect.ValueOf(ctype.NewStringArrayPrt(updaterId)))
	}
}

func initRaContent(val reflect.Value, c *ctx.Context) {
	GenerateRaContent(val)
}
