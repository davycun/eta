package dict

import (
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/duke-git/lancet/v2/slice"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

var (
	DefaultColumns = []string{"id", "updated_at", "parent_id", "namespace", "category", "parent_id", "name", "alias", "order", "remark", "order"}
)

type Dictionary struct {
	entity.BaseEntity
	Namespace string        `json:"namespace,omitempty" gorm:"column:namespace;comment:命名空间区别各个产品或定制项目" es:"type:keyword" binding:"required"`
	Category  *ctype.String `json:"category,omitempty" gorm:"column:category;comment:字典分类" es:"type:keyword"`
	ParentId  string        `json:"parent_id,omitempty" gorm:"column:parent_id;comment:字典父字典ID" es:"type:keyword"`
	Name      *ctype.String `json:"name,omitempty" gorm:"column:name;comment:字典名称" es:"type:keyword"`
	Alias     *ctype.String `json:"alias,omitempty" gorm:"column:alias;comment:字典别名" es:"type:keyword"`
	Order     int           `json:"order,omitempty" gorm:"column:order;comment:字典排序" es:"type:keyword"`
	Code      *ctype.String `json:"code,omitempty" gorm:"column:code;comment:字典编码" es:"type:keyword"`

	Children []*Dictionary `json:"children,omitempty" gorm:"-:all" binding:"ignore"`
	Parent   *Dictionary   `json:"parent,omitempty" gorm:"-:all" binding:"ignore"`
}

func (d Dictionary) TableName(namer schema.Namer) string {
	if namer == nil {
		return constants.TableDictionary
	}
	return namer.TableName(constants.TableDictionary)
}
func (d Dictionary) EsIndexName() string {
	return d.TableName(nil)
}
func (d Dictionary) RaDbFields() []string {
	return []string{"namespace", "category", "parent_id", "name", "alias", "order", "code"}
}
func (d Dictionary) MustColumns() []string {
	return []string{entity.IdDbName, entity.UpdatedAtDbName, "parent_id"}
}
func (d Dictionary) DefaultColumns() []string {
	return DefaultColumns
}

func (d *Dictionary) GetId() string {
	return d.ID
}
func (d *Dictionary) GetParentId() string {
	return d.ParentId
}
func (d *Dictionary) SetChildren(cd any) {
	if x, ok := cd.([]*Dictionary); ok {
		d.Children = x
	} else {
		logger.Error("Dictionary.SetChildren 参数不是[]*Dictionary类型")
	}
}
func (d *Dictionary) GetChildren() any {
	return d.Children
}
func (d *Dictionary) GetParentIds(db *gorm.DB) []string {
	return []string{d.ParentId}
}

func (d Dictionary) AfterMigrator(db *gorm.DB, c *ctx.Context) error {
	//如果没有数据就初始化一份数据
	return createCitizenDictionary(db, c)
}

func createCitizenDictionary(db *gorm.DB, c *ctx.Context) error {
	var batchSize = 100
	//TODO 理论上需要开启事务
	// 分批创建 修复一次性创建执行环境堆栈空间不足的问题
	for _, chunkDict := range slice.Chunk(defaultDictionary, batchSize) {
		cfl := clause.OnConflict{
			Columns: []clause.Column{
				{Name: entity.IdDbName},
			},
			DoNothing: true,
		}
		if err := dorm.TableWithContext(db, c, constants.TableDictionary).Clauses(cfl).Create(chunkDict).Error; err != nil {
			return err
		}
	}
	return nil
}
