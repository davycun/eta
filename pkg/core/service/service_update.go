package service

import (
	"errors"
	"fmt"
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/dorm/filter"
	"github.com/davycun/eta/pkg/common/dorm/xa"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/core/service/hook"
	"gorm.io/gorm"
	"reflect"
)

func (s *DefaultService) UpdateByFilters(args *dto.Param, result *dto.Result) error {

	var (
		err error
		cfg = hook.NewSrvConfig(iface.CurdModify, iface.MethodUpdateByFilters, s.GetContext(), s.GetDB(), args, result, func(o *hook.SrvConfig) {
			//互相拷贝同步，以Service的配置优先
			o.SrvOptions.Merge(s.SrvOptions)
			s.SrvOptions.Merge(o.SrvOptions)
			o.EC = s.EC
		})
	)
	defer func() {
		if !dorm.InTransaction(s.GetDB()) {
			xa.CommitOrRollback(cfg.TxDB, err)
		}
	}()

	//TODO 是否允许全量更新，应该通过比对带filters的统计结果与全量总量是否一致来确定是否是全量更新
	//不严格的做法是判断filter的size，否则统计可能会影响性能
	if len(args.Filters) < 1 {
		return errors.New("不允许更新全量数据")
	}
	cfg.CurDB = entity.SetTableName(cfg.CurDB, s.NewEntityPointer())
	err = caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			return cfg.Before()
		}).
		Call(func(cl *caller.Caller) error {
			if len(args.AuthFilters) > 0 {
				cfg.CurDB = cfg.CurDB.Where(filter.ResolveWhereTable(cfg.GetTableName(), args.AuthFilters, s.GetDbType()))
			}
			if len(args.Filters) > 0 {
				cfg.CurDB = cfg.CurDB.Where(filter.ResolveWhereTable(cfg.GetTableName(), args.Filters, s.GetDbType()))
			}
			return nil
		}).
		Call(func(cl *caller.Caller) error {
			if len(args.Columns) > 0 {
				cols := args.Columns
				cols = utils.Merge(cols, entity.UpdaterIdDbName, entity.UpdaterDeptIdDbName)
				cols = GetUpdateColumn(args.Data, cols...)
				cfg.CurDB = cfg.CurDB.Select(cols)
			}
			return cfg.CurDB.Updates(args.Data).Error
		}).
		Call(func(cl *caller.Caller) error {
			//TODO 再查一遍新的值，这里也有问题，如果在Hooks中有人修改CurDB的Where，那么Old值和NewValues就对应不上了
			//这种情况可能发生在通过回调函数来过滤不允许更新的数据（权限设置）
			cfg.NewValues = s.NewEntitySlicePointer()
			return cfg.TxDB.Model(cfg.NewValues).Where(filter.ResolveWhereTable(cfg.GetTableName(), args.Filters, s.GetDbType())).Find(cfg.NewValues).Error
		}).
		Call(func(cl *caller.Caller) error {
			result.RowsAffected = cfg.CurDB.RowsAffected
			return cfg.After()
		}).Err

	return err
}

// Update
// 注意这个方式的Update因为会执行多次Update，所以不是基于CurDB操作。如果是事务的用的是TxDB，如果不是事务的用的OriginDB
// TODO 这个是否会影响回调？
func (s *DefaultService) Update(args *dto.Param, result *dto.Result) error {
	//TODO 可以采用如下的批量更新语句，不过问题是可能每个对象更新的字段数量不一致（所以，只有保障更新所有的字段一致才可以使用）
	//update test set info=tmp.info from (values (1,'new1'),(2,'new2'),(6,'new6')) as tmp (id,info) where test.id=tmp.id;
	var (
		err error
		cfg = hook.NewSrvConfig(iface.CurdModify, iface.MethodUpdate, s.GetContext(), s.GetDB(), args, result)
	)
	defer func() {
		if !dorm.InTransaction(s.GetDB()) {
			xa.CommitOrRollback(cfg.TxDB, err)
		}
	}()

	err = caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			return cfg.Before()
		}).
		Call(func(cl *caller.Caller) error {
			return s.update(cfg)
		}).
		Call(func(cl *caller.Caller) error {
			//在查询一遍是为了保障只是获取更新成功的数据
			tmpDb := cfg.OriginDB
			if args.SingleTransaction {
				tmpDb = cfg.TxDB
			}
			cfg.NewValues = s.NewEntitySlicePointer()
			rs, ok := cfg.Result.Data.(ctype.Map)
			if !ok {
				return nil
			}
			sd, ok := rs["success_data"]
			if !ok {
				return nil
			}

			rsMap, ok := sd.([]ctype.Map)
			if !ok || len(rsMap) < 1 {
				return nil
			}
			return s.LoadUpdatedDataBySuccessId(tmpDb, rsMap, cfg.NewValues)
		}).
		Call(func(cl *caller.Caller) error {
			return cfg.After()
		}).Err

	return err
}

func (s *DefaultService) update(cfg *hook.SrvConfig) error {

	var (
		errData     = make([]ctype.Map, 0, 5)
		successData = make([]ctype.Map, 0, len(cfg.Values))
		err         error
	)
	for _, v := range cfg.Values {
		var (
			rs   = ctype.Map{}
			err1 error
		)

		tmpDb := cfg.OriginDB
		if cfg.Param.SingleTransaction {
			tmpDb = cfg.TxDB
		}
		err1 = s.updateSingle(tmpDb, cfg, &rs, v)
		if err1 != nil {
			errData = append(errData, ctype.Map{
				"id":  v.FieldByName(entity.IdFieldName).String(),
				"err": err1.Error(),
			})
			//统一一个事务，发生一个错误就直接回滚返回了
			if cfg.Param.SingleTransaction {
				err = err1
				break
			}
		} else {
			successData = append(successData, rs)
		}
	}

	if err != nil {
		cfg.Result.RowsAffected = 0
		cfg.Result.Data = ctype.Map{
			"err_data": errData,
		}
	} else {
		cfg.Result.RowsAffected = int64(len(successData))
		cfg.Result.Data = ctype.Map{
			"success_data": successData,
			"err_data":     errData,
		}
	}
	return err
}
func (s *DefaultService) updateSingle(db *gorm.DB, cfg *hook.SrvConfig, result *ctype.Map, v reflect.Value) error {
	var (
		dbType   = dorm.GetDbType(db)
		id       = v.FieldByName(entity.IdFieldName)
		updateAt = v.FieldByName(entity.UpdatedAtFieldName)
		data     = v.Addr().Interface()
		tbName   = cfg.GetTableName()
	)
	tx := dorm.TableWithContext(db, s.GetContext(), tbName)

	err := caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			if len(cfg.Param.Columns) > 0 {
				cols := cfg.Param.Columns
				cols = utils.Merge(cols, entity.UpdaterIdDbName, entity.UpdaterDeptIdDbName)
				cols = GetUpdateColumn(data, cols...)
				tx = tx.Select(cols)
			}
			return nil
		}).
		Call(func(cl *caller.Caller) error {
			if len(cfg.Param.AuthFilters) > 0 {
				tx = tx.Where(filter.ResolveWhereTable(tbName, cfg.Param.AuthFilters, dbType))
			}
			tx = tx.Where(fmt.Sprintf(`%s=? and %s=?`, dorm.Quote(dbType, "id"), dorm.Quote(dbType, "updated_at")), id.String(), updateAt.Int())
			return nil
		}).
		Call(func(cl *caller.Caller) error {
			tx = tx.Updates(data)
			return tx.Error
		}).
		Call(func(cl *caller.Caller) error {
			(*result)[entity.IdDbName] = id.String()
			(*result)[entity.UpdatedAtDbName] = updateAt.Int()
			if tx.RowsAffected < 1 {
				return errors.New("没有更新成功，没有权限或者指定的ID和updated_at数据不存在")
			}
			return nil
		}).Err

	return err
}

// GetUpdateColumn
// 获取非nil及非零值的字段，当指定更新字段的时候，同时还需要更新非零值字段
func GetUpdateColumn(obj any, mustCols ...string) []string {
	if len(mustCols) < 1 {
		return []string{}
	}
	var (
		gormFieldMap = entity.GetGormFieldName(obj)
		cols         = append([]string{}, mustCols...)
		val          = reflect.ValueOf(obj)
	)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}
	for k, v := range gormFieldMap {
		fieldVal := val.FieldByName(k)
		if fieldVal.IsValid() && !fieldVal.IsZero() {
			cols = append(cols, v)
		}
	}
	return cols
}
