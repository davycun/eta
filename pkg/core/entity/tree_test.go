package entity_test

import (
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"reflect"
	"testing"
)

type Addr struct {
	ID       string  `json:"id"`
	ParentId string  `json:"parent_id"`
	Children []*Addr `json:"children"`
	Parent   *Addr   `json:"parent"`
}

func (a *Addr) GetId() string {
	return a.ID
}
func (a *Addr) GetParentId() string {
	return a.ParentId
}
func (a *Addr) GetParentIds(db *gorm.DB) []string {
	return []string{a.ParentId}
}
func (a *Addr) SetChildren(cd any) {
	if x, ok := cd.([]*Addr); ok {
		a.Children = x
	} else {
		logger.Error("Address.SetChildren 参数不是[]*Address类型")
	}
}
func (a *Addr) GetChildren() any {
	return a.Children
}

func TestTree(t *testing.T) {

	db, _, _, _ := dorm.NewTestDB(dorm.DaMeng, "delta_dev_backend")
	ds := make([]Addr, 0, 5)

	ds = append(ds, Addr{ID: "1"})
	ds = append(ds, Addr{ID: "2", ParentId: "1"})
	ds = append(ds, Addr{ID: "3", ParentId: "2"})
	ds = append(ds, Addr{ID: "4"})
	ds = append(ds, Addr{ID: "5", ParentId: "4"})

	rs := entity.Tree(db, ds)

	assert.Equal(t, 2, len(rs))
	assert.Equal(t, "4", rs[1].ID)
	assert.Equal(t, "3", rs[0].Children[0].Children[0].ID)
}

func TestTypeName(t *testing.T) {
	//assert.Equal(t, "Addr", reflect.TypeOf(Addr{}).Name())
	//assert.Equal(t, "Addr", reflect.TypeOf(&Addr{}).String())

	assert.Equal(t, "entity_test.Addr", reflect.TypeOf(&Addr{}).Elem().String())
	assert.Equal(t, "Addr", reflect.TypeOf(&Addr{}).Elem().Name())
	assert.Equal(t, "github.com/davycun/eta/pkg/core/entity_test", reflect.TypeOf(&Addr{}).Elem().PkgPath())
	assert.Equal(t, "entity_test.Addr", reflect.TypeOf(Addr{}).String())
	assert.Equal(t, "Addr", reflect.TypeOf(Addr{}).Name())
	assert.Equal(t, "github.com/davycun/eta/pkg/core/entity_test", reflect.TypeOf(Addr{}).PkgPath())
}
