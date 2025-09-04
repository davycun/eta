package es

import (
	"errors"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/utils"
	"reflect"
)

func ScanSearch[T any](esApi *Api, after SearchAfter, fc func(dest []T) error) error {

	var (
		err error
	)

	for {
		var (
			rs []T
		)
		after, _, err = esApi.FindByAfter(&rs, after)
		if err != nil {
			break
		}
		if len(rs) < 1 {
			break
		}

		err = fc(rs)
		if err != nil {
			break
		}
	}
	return err
}
func ScanSearchWithType(esApi *Api, after SearchAfter, tp reflect.Type, fc func(dest any) error) error {

	var (
		err error
	)

	if tp == nil {
		return errors.New("tp is nil")
	}

	for {
		rs := reflect.New(reflect.SliceOf(utils.GetRealType(tp))).Interface()
		after, _, err = esApi.FindByAfter(rs, after)
		if err != nil {
			break
		}
		val := reflect.ValueOf(rs).Elem()
		if val.Len() < 1 {
			break
		}
		err = fc(rs)
		if err != nil {
			break
		}
	}
	return err
}

// ScanGroupBy
// path 是nested类型的列名
func ScanGroupBy(esApi *Api, after ctype.Map, fc func(path string, dest AggregateResult) error) error {

	var (
		size = esApi.getLimit()
	)

	for k, _ := range esApi.groupCol {
		for {
			ar, err := esApi.reqGroup(k, size, after)
			if err != nil {
				return err
			}
			if len(ar.Group) < 1 {
				break
			}

			err = fc(k, ar)
			if err != nil {
				return err
			}

			after = ar.AfterKey
			if len(ar.Group) < size {
				break
			}
		}
	}
	return nil
}
