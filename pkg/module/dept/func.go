package dept

import (
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/module/namer"
	"github.com/duke-git/lancet/v2/slice"
)

func LoadParents(ct *ctx.Context, deptId string, containCurrentDept bool) []namer.IdName {
	names := make([]namer.IdName, 0)
	if deptId == "" {
		return names
	}

	mnMap, err := namer.LoadAllIdName(ct)
	if err != nil {
		logger.Warnf("load all id name error: %v", err)
		return names
	}
	//mnMap = maputil.Filter(mnMap, func(k string, v namer.IdName) bool {
	//	return v.Tp == namer.TypeDept
	//})

	var dId = deptId
	for {
		if dId == "" {
			break
		}
		if n, ok := mnMap[dId]; ok {
			names = append(names, n)
			dId = n.ParentId
		} else {
			break
		}
	}

	if !containCurrentDept {
		names = slice.Filter(names, func(i int, v namer.IdName) bool {
			return v.ID != deptId
		})
	}
	return names
}
