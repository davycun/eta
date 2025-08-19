package ra

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/crypt"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/es"
	"github.com/davycun/eta/pkg/common/dorm/filter"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/setting"
	"strings"

	"github.com/duke-git/lancet/v2/slice"
	"gorm.io/gorm"
)

func CreateTrigger(db *gorm.DB, tableName string, raFields []string) error {
	if len(raFields) <= 0 {
		return nil
	}
	var (
		scm    = dorm.GetDbSchema(db)
		dbType = dorm.GetDbType(db)
	)
	return caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			switch dbType {
			case dorm.DaMeng:
				return createDmTrigger(db, tableName, raFields)
			case dorm.PostgreSQL:
				return createPostgresTrigger(db, scm, tableName, raFields)
			case dorm.Mysql:
				return createMysqlTrigger(db, scm, tableName, raFields)
			}
			return nil
		}).
		Call(func(cl *caller.Caller) error {
			switch dbType {
			case dorm.DaMeng:
				return dorm.CreateContextIndex(db, tableName, entity.RaContentDbName)
			}
			return nil
		}).Err
}

func DropTrigger(db *gorm.DB, tableName string) error {
	var (
		dbType      = dorm.GetDbType(db)
		scm         = dorm.GetDbSchema(db)
		triggerName = fmt.Sprintf(`trigger_%s_ra`, tableName)
	)
	switch dbType {
	case dorm.DaMeng, dorm.PostgreSQL:
		return dorm.TriggerDelete(db, scm+"."+triggerName, scm+"."+tableName)
	case dorm.Mysql:
		return dropMysqlTrigger(db, scm, tableName)
	}
	return nil
}

func KeywordToFilters(db *gorm.DB, tableName string, searchContent string, dbType dorm.DbType) []filter.Filter {
	if searchContent == "" {
		return nil
	}

	var (
		tb, b            = setting.GetTableConfig(db, tableName)
		plainTextFilters = make([]filter.Filter, 0)
		encTextFilters   = make([]filter.Filter, 0)
	)

	if dbType == dorm.ES {
		plainTextFilters = es.KeywordToFilters(entity.RaContentDbName, searchContent)
	} else {
		plainTextFilters = filter.KeywordToFilter(entity.RaContentDbName, searchContent)
	}

	if !b || len(tb.CryptFields) < 1 {
		return plainTextFilters
	}

	//只取一个？？
	cryptConf := tb.CryptFields[0]

	enc := func(c string) string {
		encStr, err := crypt.EncryptBase64(cryptConf.Algo, cryptConf.SecretKey[0], c)
		if err != nil {
			return ""
		}
		return encStr
	}
	encContent := make([]string, 0)

	scs := utils.Split(searchContent, ",", "，", " ")
	for _, sc := range scs {
		contentSlice := crypt.StringToSlice(sc, 1)
		encStrList := slice.Map(contentSlice, func(i int, v string) string { return enc(v) })
		encContent = append(encContent, strings.Join(encStrList, constants.CryptSliceSeparator))
	}

	if dbType == dorm.ES {
		encTextFilters = es.KeywordToFilters(entity.RaContentDbName, strings.Join(encContent, " "))
	} else {
		encTextFilters = filter.KeywordToFilter(entity.RaContentDbName, strings.Join(encContent, " "))
	}

	fls := make([]filter.Filter, 0, len(plainTextFilters))
	for i := range plainTextFilters {
		fls = append(fls, filter.Filter{
			LogicalOperator: filter.And,
			Filters: []filter.Filter{
				{LogicalOperator: filter.And, Filters: []filter.Filter{plainTextFilters[i]}},
				{LogicalOperator: filter.Or, Filters: []filter.Filter{encTextFilters[i]}},
			},
		})
	}
	return fls
}
