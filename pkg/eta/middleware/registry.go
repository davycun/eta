package middleware

import (
	"github.com/davycun/eta/pkg/common/logger"
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
