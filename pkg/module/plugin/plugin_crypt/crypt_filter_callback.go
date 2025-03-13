package plugin_crypt

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/expr"
	"github.com/davycun/eta/pkg/common/dorm/filter"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/duke-git/lancet/v2/slice"
	"gorm.io/gorm"
	"strings"
)

const (
	likeSymbol = "%"
)

// EncryptQueryParam 查询请求，对切片加密对字段进行加密
func EncryptQueryParam(cfg *hook.SrvConfig, pos hook.CallbackPosition) (err error) {
	if pos != hook.CallbackBefore {
		return
	}
	var (
		table = entity.GetContextTable(cfg.Ctx)
	)
	if len(table.CryptFields) < 1 {
		return
	}

	fieldCryptMap := make(map[string]entity.CryptFieldInfo)
	slice.ForEach(table.CryptFields, func(i int, v entity.CryptFieldInfo) {
		if v.Enable && v.Field != "" {
			fieldCryptMap[v.Field] = v
		}
	})
	if len(fieldCryptMap) <= 0 {
		return
	}

	logger.Debugf("EncryptQueryParam table: %s", table.GetTableName())

	cfg.Param.Filters = encryptFilters(cfg.OriginDB, fieldCryptMap, cfg.Param.Filters)
	cfg.Param.Auth2RoleFilters = encryptFilters(cfg.OriginDB, fieldCryptMap, cfg.Param.Auth2RoleFilters)
	cfg.Param.RecursiveFilters = encryptFilters(cfg.OriginDB, fieldCryptMap, cfg.Param.RecursiveFilters)
	cfg.Param.AuthRecursiveFilters = encryptFilters(cfg.OriginDB, fieldCryptMap, cfg.Param.AuthRecursiveFilters)

	return
}

func encryptFilters(curDb *gorm.DB, fieldCryptMap map[string]entity.CryptFieldInfo, filters []filter.Filter) []filter.Filter {
	for i, f := range filters {
		if ft, ok := fieldCryptMap[f.Column]; ok {
			switch x := f.Value.(type) {
			case string:
				var (
					prefix, suffix, ed = "", "", ""
					cryptOk            bool
				)
				switch f.Operator {
				case filter.Like:
					if strings.HasPrefix(x, likeSymbol) {
						prefix = likeSymbol
						x = strings.TrimPrefix(x, likeSymbol)
					}
					if strings.HasSuffix(x, likeSymbol) {
						suffix = likeSymbol
						x = strings.TrimSuffix(x, likeSymbol)
					}
					ed, cryptOk = encryptData(ft, []rune(x))
					ed = strings.TrimPrefix(ed, constants.CryptPrefix)
				default:
					ed, cryptOk = encryptData(ft, []rune(x))
				}
				if cryptOk {
					filters[i].Value = strings.Join([]string{prefix, ed, suffix}, "")
				}
			default:
				continue
			}
		}
		//解决达梦的 contains函数全文检索查询,contains(column, 'value')，暂时无法解决所有的表达式
		if strings.HasPrefix(f.Expr.Expr, "contains(") {
			var (
				dbType = dorm.GetDbType(curDb)
				ft     = entity.CryptFieldInfo{}
				ok     = false
			)
			if dbType == dorm.DaMeng {
				for _, v := range f.Expr.Vars {
					if v.Type == expr.VarTypeColumn {
						ft, ok = fieldCryptMap[fmt.Sprintf("%s", v.Value)]
					}
				}
				if ok {
					for x, v := range f.Expr.Vars {
						if v.Type == expr.VarTypeValue && ok {
							str, cryptOk := encryptData(ft, []rune(fmt.Sprintf("%s", v.Value)))
							if cryptOk {
								f.Expr.Vars[x].Value = str
							}
						}
					}
				}
			}
		}
		if len(f.Filters) > 0 {
			filters[i].Filters = encryptFilters(curDb, fieldCryptMap, f.Filters)
		}
	}
	return filters
}
