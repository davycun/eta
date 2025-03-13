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
	"github.com/duke-git/lancet/v2/slice"
	"gorm.io/gorm"
	"strings"
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
		triggerName = fmt.Sprintf(`"trigger_%s_ra"`, tableName)
	)
	switch dbType {
	case dorm.DaMeng:
		return db.Exec(fmt.Sprintf(`drop trigger if exists "%s".%s`, scm, triggerName)).Error
	case dorm.PostgreSQL:
		return db.Exec(fmt.Sprintf(`drop trigger if exists %s on "%s"."%s" cascade`, triggerName, scm, tableName)).Error
	case dorm.Mysql:
		return dropMysqlTrigger(db, scm, tableName)
	}
	return nil
}

func KeywordToFilters(db *gorm.DB, tableName string, searchContent string) []filter.Filter {
	tb, b := setting.GetTableConfig(db, tableName)
	plainTextFilters := es.KeywordToFilters(entity.RaContentDbName, searchContent)

	if !b {
		return plainTextFilters
	}

	if len(tb.CryptFields) <= 0 {
		return plainTextFilters
	}

	//TODO 只是取一个，这里理论上也是不对的
	//TODO 应该把所有的检索内容根据所有的加密字段进行加密，然后取或的关系
	cryptConf := tb.CryptFields[0]

	enc := func(c string) string {
		encStr, err := crypt.NewEncrypt(cryptConf.Algo, cryptConf.SecretKey[0]).FromRawString(c).ToBase64String()
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

	encTextFilters := es.KeywordToFilters(entity.RaContentDbName, strings.Join(encContent, " "))

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
