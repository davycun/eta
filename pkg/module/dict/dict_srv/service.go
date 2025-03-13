package dict_srv

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/core/builder"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/service"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/dict"
)

type Service struct {
	service.DefaultService
}

func (s *Service) TreeDelete(args *dto.Param, result *dto.Result) error {
	var (
		ids []string
	)
	bd := builder.NewRecursiveSqlBuilder(dorm.GetDbType(s.GetDB()), dorm.GetDbSchema(s.GetDB()), constants.TableDictionary)
	bd.AddRecursiveFilter(args.RecursiveFilters...).SetUp(args.IsUp)
	bd.AddColumn("id").AddFilter(args.Filters...)
	listSql, _, err := bd.Build()

	if err != nil {
		return err
	}
	// 获取出来关于要删除的字典以及子字典的ID
	err = dorm.RawFetch(listSql, s.GetDB(), &ids)

	tx := entity.SetTableName(s.GetDB(), &dict.Dictionary{})
	tx = tx.Where(fmt.Sprintf(`"%s"."id" in (%s)`, constants.TableDictionary, listSql)).Delete(&dict.Dictionary{})
	if tx.Error != nil {
		return tx.Error
	}
	result.RowsAffected = tx.RowsAffected
	//清除缓存
	dict.DataCache.DeleteAll(s.GetDB())
	return nil
}
