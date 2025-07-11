package entity

import (
	"gorm.io/gorm"
)

type TreeEntity interface {
	GetId() string
	GetParentId() string
	GetChildren() any
	SetChildren(any)
	// GetParentIds
	//获取当前节点的父节点、爷爷节点等直到顶点，返回的上层节点必须按照顺序返回。确保在缺失父节点的情况下，能够把孙子节点挂接到爷爷节点上
	//比如有a、b、c、d、e四个节点（他们ID分别为1、2、3、4、5），如果当前节点是e，GetParentIds返回的结果应该是[4、3、2、1]
	//假设我们现在需要组装树结构的切片节点是[a、c、e]，如果不考虑跨级挂靠，那么a、c、e是平级返回的，如果考虑跨级挂靠那么树形算法后，应该返回 a->c->e
	//这个是要考虑e能挂靠到c上，c能挂靠到a上，那么就需要知道每个节点的所有父级节点，所以GetParentIds就是这个目的
	GetParentIds(db *gorm.DB) []string
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
