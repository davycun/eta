package plugin_crypt

import (
	"github.com/davycun/eta/pkg/common/crypt"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/core/service/hook"
	"reflect"
	"slices"
	"strings"
)

// StoreSign 在数据插入和修改前，填充签名
func StoreSign(cfg *hook.SrvConfig, pos hook.CallbackPosition) (err error) {

	//如果是新增，那么就是所有字段进行拼接然后签名即可
	//如果是update_by_filters，需要更新签名的情况
	//1）param.Data中签名相关字段有值
	//2）columns中有指定强制更新签名相关字段
	//如果是update_by_filters，需要更新签名的情况
	//1）param.Data中签名相关字段有值
	//2）columns中有指定强制更新签名相关字段

	//如果涉及到签名字段更新，那么就重新计算签名，重新计算签名就可能涉及到需要从oldValues取没有值的签名字段

	if pos != hook.CallbackBefore {
		return
	}

	table := entity.GetContextTable(cfg.Ctx)
	valMap := getOldValueMap(cfg.OldValues)

	if len(table.SignFields) < 1 {
		return
	}

	switch cfg.Method {
	case iface.MethodCreate, iface.MethodUpdate:
		for _, v := range cfg.Values {
			for _, signInfo := range table.SignFields {
				if !signInfo.Enable || len(table.SignFields) < 1 {
					continue
				}
				//如果是更新操作，columns强制指定了需要签名的字段，那么就强制更新签名字段
				if cfg.Method == iface.MethodUpdate || cfg.Method == iface.MethodUpdateByFilters {
					if len(cfg.Param.Columns) > 0 && utils.ContainAny(cfg.Param.Columns, signInfo.Fields...) {
						tmp := entity.GetValue(v, signInfo.Field)
						if tmp.IsValid() { //存在签名字段，那么就强制更新
							cfg.Param.Columns = append(cfg.Param.Columns, signInfo.Field)
						}
					}
				}
				err = signValue(v, valMap[entity.GetString(v, entity.IdDbName)], signInfo)
				if err != nil {
					return err
				}
			}
		}
	}
	logger.Debugf("%s 已加密%d条数据", table.GetTableName(), len(cfg.Values))
	return
}

func getOldValueMap(oldValues any) map[string]reflect.Value {

	var (
		valMap = make(map[string]reflect.Value)
	)
	if oldValues == nil {
		return valMap
	}
	valList := utils.ConvertToValueArray(oldValues)
	for _, v := range valList {
		str := entity.GetString(v, entity.IdDbName)
		if str != "" {
			valMap[str] = v
		}
	}
	return valMap
}

func signValue(val reflect.Value, oldVal reflect.Value, signInfo entity.SignFieldsInfo) error {
	str, err := signData(val, oldVal, signInfo)
	if err != nil || str == "" {
		return err
	}

	return entity.SetValue(val, signInfo.Field, str)
}

// signData 计算数据签名
func signData(val reflect.Value, oldVal reflect.Value, signInfo entity.SignFieldsInfo) (sign string, err error) {

	if !val.IsValid() || val.IsZero() {
		return
	}

	//如果是模版表，属性会自动添加F前缀，GetValueField能够通过json或者gorm来查找，能解决
	signFieldVal := entity.GetValue(val, signInfo.Field)
	if !signFieldVal.IsValid() {
		return
	}
	var (
		bd = strings.Builder{}
	)
	slices.Sort(signInfo.Fields)
	for _, v := range signInfo.Fields {
		valStr := entity.GetString(val, v)
		if valStr == "" {
			valStr = entity.GetString(oldVal, v)
		}
		bd.WriteString(valStr)
	}

	if bd.String() == "" {
		return
	}
	return crypt.NewEncrypt(signInfo.Algo, signInfo.Key).FromRawString(bd.String()).ToBase64String()
}

// VerifySign 数据查询后验签
func VerifySign(cfg *hook.SrvConfig, pos hook.CallbackPosition) (err error) {

	if cfg.CurdType != iface.CurdRetrieve {
		return
	}

	var (
		table     = entity.GetContextTable(cfg.Ctx)
		valueList = utils.ConvertToValueArray(cfg.Result.Data)
	)

	if len(table.SignFields) < 1 {
		return
	}

	switch pos {
	case hook.CallbackBefore:
		for _, signInfo := range table.SignFields {
			if !signInfo.Enable || len(table.SignFields) < 1 {
				continue
			}
			//TODO 可能配置错误Field字段不存在
			cfg.Param.MustColumns = append(cfg.Param.MustColumns, signInfo.Field)
		}
	case hook.CallbackAfter:
		for _, v := range valueList {
			for _, signInfo := range table.SignFields {
				if !signInfo.Enable || len(table.SignFields) < 1 {
					continue
				}
				//如果获取的字段不包括所有字段，那么就不校验
				if len(cfg.Param.Columns) > 0 && !utils.ContainAll(cfg.Param.Columns, signInfo.Fields...) {
					continue
				}
				err = verifyValue(v, signInfo)
				if err != nil {
					return err
				}
			}
		}
	}
	return
}
func verifyValue(val reflect.Value, signInfo entity.SignFieldsInfo) error {

	if !val.IsValid() || !signInfo.Enable {
		return nil
	}
	fieldStr := entity.GetString(val, signInfo.Field)
	if fieldStr == "" {
		//如果之前的加签字段为空（有可能是select没有取这个字段）那么就跳过
		return nil
	}

	str, err := signData(val, reflect.Value{}, signInfo)
	if err != nil || str == "" {
		return err
	}

	return entity.SetValue(val, signInfo.VerifyField, str)
}
