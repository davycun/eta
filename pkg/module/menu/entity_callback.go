package menu

import (
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/duke-git/lancet/v2/slice"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (m Menu) AfterMigrator(db *gorm.DB, c *ctx.Context) error {
	//如果没有数据就初始化一份数据
	return initMenu(c, db)
}

func initMenu(c *ctx.Context, db *gorm.DB) error {

	ids := make([]string, 0, len(defaultMenu))
	for _, v := range defaultMenu {
		ids = append(ids, v.ID)
	}
	var batchSize = 100
	//TODO 理论上需要开启事务
	// 分批创建 修复一次性创建执行环境堆栈空间不足的问题
	for _, chunkDict := range slice.Chunk(defaultMenu, batchSize) {
		cfl := clause.OnConflict{
			Columns: []clause.Column{
				{Name: entity.IdDbName},
			},
			DoNothing: true,
		}
		if err := dorm.TableWithContext(db, c, constants.TableMenu).Clauses(cfl).Create(chunkDict).Error; err != nil {
			return err
		}
	}
	return nil
}
