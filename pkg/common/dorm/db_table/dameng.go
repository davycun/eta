package db_table

import (
	"github.com/davycun/eta/pkg/common/logger"
	"gorm.io/gorm"
)

var (
	objCache = make(map[int64]objects)
)

type objects struct {
	Id      int64  `json:"id,omitempty" gorm:"column:ID"`
	Name    string `json:"name,omitempty" gorm:"column:NAME"`
	ObjType string `json:"obj_type,omitempty" gorm:"column:TYPE$"`
	SubType string `json:"sub_type,omitempty" gorm:"column:SUBTYPE$"`
	SchId   int64  `json:"sch_id,omitempty" gorm:"column:SCHID"`
}

func fetchObjectById(db *gorm.DB, id int64) objects {

	obj, ok := objCache[id]
	if ok {
		return obj
	}
	var objs objects
	err := db.Table("SYS.SYSOBJECTS").Select("ID", "NAME", "TYPE$", "SUBTYPE$", "SCHID").Where(`ID = ?`, id).First(&objs).Error
	if err != nil {
		logger.Errorf("can not find objects by id[%d] in dameng because %s", id, err)
	}
	return objs
}
