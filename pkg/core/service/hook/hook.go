package hook

import (
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/eta/constants"
	"slices"
	"strings"
	"sync"
)

var (
	//存储的内容是tableName -> []Callback
	extendCallbacks = sync.Map{}
)

const (
	CallbackForAll = "all_table_callback" //这是为了回调所有表（所有的service）
)

type Callback func(cfg *SrvConfig, pos CallbackPosition) error
type CallbackOption struct {
	TableName string
	Order     int
	CurdType  iface.CurdType //是修改还是查询
	Methods   []iface.Method
	IsAuth    bool
	Name      string //回调的名称
}
type CallbackOptionFunc func(*CallbackOption)

type CallbackWrapper struct {
	CallbackOption //顺序
	Callback       Callback
}

func AddModifyCallback(tableName string, cb Callback, callbackOption ...CallbackOptionFunc) {
	addCallback(tableName, cb, iface.CurdModify, false, callbackOption...)
}
func AddRetrieveCallback(tableName string, cb Callback, callbackOption ...CallbackOptionFunc) {
	addCallback(tableName, cb, iface.CurdRetrieve, false, callbackOption...)
}
func AddAuthCallback(tableName string, cb Callback, callbackOption ...CallbackOptionFunc) {
	addCallback(tableName, cb, iface.CurdAll, true, callbackOption...)
}

func addCallback(tableName string, cb Callback, curdType iface.CurdType, isAuth bool, callbackOption ...CallbackOptionFunc) {
	var (
		rs      = make([]CallbackWrapper, 0, 2)
		wrapper = CallbackWrapper{
			Callback: cb,
			CallbackOption: CallbackOption{
				TableName: tableName,
			},
		}
		val, ok = extendCallbacks.Load(tableName)
	)
	if ok {
		rs = val.([]CallbackWrapper)
	}
	for _, fc := range callbackOption {
		fc(&wrapper.CallbackOption)
	}
	//默认不添加排序就是100，拍得比较靠后
	if wrapper.Order == 0 {
		wrapper.Order = 100
	}
	wrapper.IsAuth = isAuth
	//放在这里是避免callbackOption 配置错误
	wrapper.CurdType = curdType

	if len(wrapper.Methods) < 1 {
		wrapper.Methods = []iface.Method{iface.MethodAll}
	}

	//如果是权限回调，那么如果重复就需要被覆盖。如果是常规的回调，是允许重复的
	if isAuth {
		rs = slices.DeleteFunc(rs, func(v CallbackWrapper) bool {
			if utils.ContainAll(v.Methods, wrapper.Methods...) && v.CurdType == wrapper.CurdType && v.IsAuth {
				logger.Errorf("在添加权限回调[name:%s,tableName:%s,method:%s,curdType:%s]的时候，有重复会被覆盖", wrapper.Name, tableName, wrapper.Methods, wrapper.CurdType)
				return true
			}
			return false
		})
		rs = append(rs, wrapper)
	} else {
		rs = append(rs, wrapper)
	}

	//排序
	slices.SortFunc(rs, func(a, b CallbackWrapper) int {
		return a.Order - b.Order
	})
	extendCallbacks.Store(tableName, rs)
}

func callCallback(cfg *SrvConfig, curdType iface.CurdType, method iface.Method, pos CallbackPosition, callbacks ...CallbackWrapper) error {
	var (
		err    error
		cbList []CallbackWrapper
	)
	switch curdType {
	case iface.CurdModify:
		cbList = getModifyExtendCallback(cfg.GetTableName())
	case iface.CurdRetrieve:
		cbList = getRetrieveExtendCallback(cfg.GetTableName())
	}
	cbList = append(cbList, callbacks...)
	slices.SortFunc(cbList, func(a, b CallbackWrapper) int {
		return a.Order - b.Order
	})
	for _, fc := range cbList {
		if len(fc.Methods) < 1 || utils.ContainAny(fc.Methods, method, iface.MethodAll) {
			err = fc.Callback(cfg, pos)
			if err != nil {
				return err
			}
		}
	}
	return err

}
func callAuthCallback(cfg *SrvConfig, pos CallbackPosition) error {
	if cfg.UseParamAuth && cfg.Param.DisablePermFilter {
		return nil
	}
	var (
		first  []CallbackWrapper
		second []CallbackWrapper
		third  []CallbackWrapper
		fourth []CallbackWrapper
	)
	cbList := getAuthExtendCallback(cfg.GetTableName())
	for _, fc := range cbList {
		if cfg.GetTableName() == fc.TableName {
			if utils.ContainAll(fc.Methods, cfg.Method) {
				first = append(first, fc)
			}
			if utils.ContainAll(fc.Methods, iface.MethodAll) {
				second = append(second, fc)
			}
		}
		if fc.TableName == CallbackForAll {
			if utils.ContainAll(fc.Methods, cfg.Method) {
				third = append(third, fc)
			}
			if utils.ContainAll(fc.Methods, iface.MethodAll) {
				fourth = append(fourth, fc)
			}
		}
	}

	if len(first) > 0 {
		for _, cb := range first {
			err := cb.Callback(cfg, pos)
			if err != nil {
				return err
			}
		}
	} else if len(second) > 0 {
		for _, cb := range second {
			err := cb.Callback(cfg, pos)
			if err != nil {
				return err
			}
		}
	} else if len(third) > 0 {
		for _, cb := range third {
			err := cb.Callback(cfg, pos)
			if err != nil {
				return err
			}
		}
	} else {
		for _, cb := range fourth {
			err := cb.Callback(cfg, pos)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func getModifyExtendCallback(tbName string) []CallbackWrapper {
	mcs := getExtendCallback(tbName)
	mcs = slices.DeleteFunc(mcs, func(wrapper CallbackWrapper) bool {
		return wrapper.CurdType != iface.CurdModify || wrapper.IsAuth
	})
	return mcs
}
func getRetrieveExtendCallback(tbName string) []CallbackWrapper {
	mcs := getExtendCallback(tbName)
	mcs = slices.DeleteFunc(mcs, func(wrapper CallbackWrapper) bool {
		return wrapper.CurdType != iface.CurdRetrieve || wrapper.IsAuth
	})
	return mcs
}
func getAuthExtendCallback(tbName string) []CallbackWrapper {
	mcs := getExtendCallback(tbName)
	mcs = slices.DeleteFunc(mcs, func(wrapper CallbackWrapper) bool {
		return !wrapper.IsAuth
	})
	return mcs
}
func getExtendCallback(tbName string) []CallbackWrapper {
	mcs := make([]CallbackWrapper, 0, 3)
	mc, ok := extendCallbacks.Load(CallbackForAll)
	if ok {
		tmp := mc.([]CallbackWrapper)
		mcs = append(mcs, tmp...)
	}
	mc, ok = extendCallbacks.Load(tbName)
	if ok {
		tmp := mc.([]CallbackWrapper)
		mcs = append(mcs, tmp...)
	}
	after, found := strings.CutPrefix(tbName, constants.TableTemplatePrefix)
	if found {
		mc, ok = extendCallbacks.Load(after)
		if ok {
			tmp := mc.([]CallbackWrapper)
			mcs = append(mcs, tmp...)
		}
	}
	return mcs
}
