package entity

import (
	"gorm.io/gorm"
)

type TreeEntity interface {
	GetId() string
	GetParentId() string
	GetChildren() any
	SetChildren(any)
	GetParentIds(db *gorm.DB) []string //获取所有的parentId，按照切片顺序
}

func Tree[E any](db *gorm.DB, data []E) (rs []*E) {
	var (
		mp = make(map[string]*E)
	)
	for i, _ := range data {
		v := &data[i]
		id, _ := GetIdAndParentId(v)
		mp[id] = v
	}

	for i, _ := range data {
		v := &data[i]
		id, _ := GetIdAndParentId(v)
		curNode := mp[id]

		if pNode, ok := getParent(db, v, mp); ok {
			if x, ok1 := any(pNode).(TreeEntity); ok1 {
				if children, ok2 := x.GetChildren().([]*E); ok2 {
					children = append(children, any(curNode).(*E))
					x.SetChildren(children)
				}
			}
		} else {
			rs = append(rs, curNode)
		}
	}
	return rs
}
func getParent[E any](db *gorm.DB, e *E, data map[string]*E) (*E, bool) {
	_, parentId := GetIdAndParentId(e)
	if pNode, ok := data[parentId]; ok {
		return pNode, true
	}
	if x, ok := any(e).(TreeEntity); ok {
		pIds := x.GetParentIds(db)
		for _, v := range pIds {
			if pNode, ok1 := data[v]; ok1 {
				return pNode, true
			}
		}
	}
	return nil, false
}

// GetIdAndParentId
// data参数是针对实现了TreeEntity接口的结构体
func GetIdAndParentId(data any) (id, parentId string) {

	if x, ok := data.(TreeEntity); ok {
		return x.GetId(), x.GetParentId()
	}
	return GetString(data, IdFieldName), GetString(data, "ParentId")
}
