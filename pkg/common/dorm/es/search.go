package es

import (
	"errors"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/utils"
	"reflect"
)

func ScanSearch[T any](esApi *Api, after SearchAfter, fc func(dest []T) (bool, error)) error {

	var (
		err error
	)

	for {
		var (
			rs   []T
			stop bool
		)
		after, _, err = esApi.FindByAfter(&rs, after)
		if err != nil {
			break
		}
		if len(rs) < 1 {
			break
		}

		stop, err = fc(rs)
		if err != nil || stop {
			break
		}
	}
	return err
}
func ScanSearchWithType(esApi *Api, after SearchAfter, tp reflect.Type, fc func(dest any) (bool, error)) error {

	var (
		err  error
		stop bool
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
		stop, err = fc(rs)
		if err != nil || stop {
			break
		}
	}
	return err
}

// ScanGroupBy
// path 是nested类型的列名
func ScanGroupBy(esApi *Api, after ctype.Map, fc func(path string, dest AggregateResult) (bool, error)) error {

	var (
		size = esApi.getLimit()
		stop bool
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

			stop, err = fc(k, ar)
			if err != nil || stop {
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
