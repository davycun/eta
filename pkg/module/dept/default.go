package dept

import (
	"github.com/davycun/eta/pkg/common/global"
	"sync"
)

var (
	//userId -> RelationDept
	cacheRd = sync.Map{}
	//
	lok = sync.Mutex{}
)

func GetDefaultUser2Dept(userId string, name string) RelationDept {

	dp, ok := cacheRd.Load(userId)
	if ok {
		return dp.(RelationDept)
	}
	lok.Lock()
	defer lok.Unlock()

	dp, ok = cacheRd.Load(userId)
	if ok {
		return dp.(RelationDept)
	}

	rd := RelationDept{}
	rd.ID = global.GenerateIDStr()
	rd.FromId = userId
	rd.ToId = userId
	rd.Dept.ID = userId
	rd.Dept.Name = name
	rd.IsManager = true
	rd.IsMain = true

	cacheRd.Store(userId, rd)
	return rd
}
