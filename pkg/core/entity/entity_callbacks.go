package entity

import (
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/global"
	"gorm.io/gorm"
	"time"
)

func (b *BaseEntity) BeforeCreate(db *gorm.DB) error {
	return BeforeCreate(b, dorm.GetDbContext(db))
}

func (b *BaseEntity) BeforeUpdate(db *gorm.DB) error {
	return beforeUpdate(db, b)
}
func (b *BaseEdgeEntity) BeforeCreate(db *gorm.DB) error {
	return BeforeCreate(&b.BaseEntity, dorm.GetDbContext(db))
}

func (b *BaseEdgeEntity) BeforeUpdate(db *gorm.DB) error {
	return beforeUpdate(db, &(b.BaseEntity))
}

func BeforeCreate(b *BaseEntity, c *ctx.Context) error {

	if b.ID == "" {
		b.ID = global.GenerateIDStr()
	}
	if b.CreatedAt == nil || !b.CreatedAt.IsValid() {
		b.CreatedAt = ctype.NewLocalTimePrt(time.Now())
	}
	if b.UpdatedAt < 1 {
		b.UpdatedAt = time.Now().UnixMilli()
	}

	if c != nil {
		userId := c.GetContextUserId()
		curDeptId := c.GetContextCurrentDeptId()
		if curDeptId == "" {
			curDeptId = userId
		}
		if b.CreatorId == "" {
			b.CreatorId = userId
		}
		if b.UpdaterId == "" {
			b.UpdaterId = userId
		}
		if b.CreatorDeptId == "" {
			b.CreatorDeptId = curDeptId
		}
		if b.UpdaterDeptId == "" {
			b.UpdaterDeptId = curDeptId
		}
		if b.FieldUpdaterIds == nil || !b.FieldUpdaterIds.Valid {
			b.FieldUpdaterIds = ctype.NewStringArrayPrt(b.UpdaterId)
		}
	}

	return nil
}
func beforeUpdate(db *gorm.DB, b *BaseEntity) error {
	c, ok := ctx.GetCurrentContext()
	if ok {
		if b.UpdatedAt < 1 {
			b.UpdatedAt = time.Now().UnixMilli()
		}
		if b.UpdaterId == "" {
			b.UpdaterId = c.GetContextUserId()
		}
		if b.UpdaterDeptId == "" {
			b.UpdaterDeptId = c.GetContextCurrentDeptId()
		}
	}
	return nil
}
