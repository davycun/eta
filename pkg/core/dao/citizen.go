package dao

import "github.com/davycun/eta/pkg/common/dorm/filter"

func GetDistinct(labelIds [][]string, filters ...[]filter.Filter) bool {
	tt := 0
	if len(labelIds) > 0 {
		tt += 1
		if len(labelIds) > 1 {
			tt += 1
		}
	}
	if filters != nil && len(filters) > 0 {
		for _, v := range filters {
			if len(v) > 0 {
				tt++
			}
		}
	}
	return tt < 2
}

func Distinct(filters ...[]filter.Filter) bool {
	tt := 0
	if filters != nil && len(filters) > 0 {
		for _, v := range filters {
			if len(v) > 0 {
				tt++
			}
		}
	}
	return tt < 2
}
