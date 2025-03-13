package middleware

import (
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/logger"
	"slices"
)

var (
	// key: middleware name, value: middle
	middlewareMap = map[string]MidOption{}
)

func Registry(midList ...MidOption) {
	for _, v := range midList {
		if _, ok := middlewareMap[v.Name]; ok {
			logger.Errorf("the middleware[%s] has exists,can not be overwrite", v.Name)
		} else {
			middlewareMap[v.Name] = v
		}
	}
}
func Remove(name string) {
	delete(middlewareMap, name)
}

func InitMiddleware() {
	mds := sortMiddleware()
	for _, v := range mds {
		global.GetGin().Use(v.HandlerFunc)
	}
}

func sortMiddleware() []MidOption {
	mds := make([]MidOption, 0, len(middlewareMap))
	for _, v := range middlewareMap {
		mds = append(mds, v)
	}
	slices.SortFunc(mds, func(a, b MidOption) int {
		return a.Order - b.Order
	})
	return mds
}
