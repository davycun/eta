package service

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/core/service/hook"
	"gorm.io/gorm"
)

type CommonService struct {
	iface.SrvOptions
	Callbacks []hook.Callback
}

func (s *CommonService) Options(fs ...iface.SrvOptionsFunc) {
	for _, f := range fs {
		f(&s.SrvOptions)
	}
}
func (s *CommonService) EsRetrieveEnabled() bool {
	return global.GetES() != nil && s.GetTable().EsRetrieveEnabled()
}

type CommonDbService struct {
	CommonService
}

type DefaultService struct {
	CommonDbService
}

func (s *DefaultService) Init(c *ctx.Context, db *gorm.DB, ec *iface.EntityConfig) error {
	s.SetContext(c)
	s.SetDB(db)
	s.SetEntityConfig(ec)
	return nil
}

// NewModifyConfig
// TODO 这里其实可能会有一个问题，就是加载的OldValues是接口期望的，但是经过Before的Callbacks 增加了一些Scopes然后就导致OldValues与实际的发生不一致的情况
// 但是有时候可能又需要再Before的回调中用到OldValues数据，所以提前获取

func (s *DefaultService) LoadUpdatedDataBySuccessId(db *gorm.DB, rs []ctype.Map, target any) error {

	ids := make([]string, 0, len(rs))
	for _, v := range rs {
		sid, ok := v[entity.IdDbName]
		if ok {
			ids = append(ids, sid.(string))
		}
	}
	if len(ids) < 1 {
		return nil
	}
	return dorm.TableWithContext(db, s.GetContext(), s.GetTableName()).
		Where(fmt.Sprintf(`%s in ?`, dorm.Quote(dorm.GetDbType(db), "id")), ids).Find(target).Error
}
