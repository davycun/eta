package bean

import (
	"errors"
	"fmt"
	"reflect"
)

// Copy
// TODO 应该要有缓存，否则会影响性能，类型JSON序列化一样
func Copy(dst, src any) error {

	dt := reflect.TypeOf(dst)
	st := reflect.TypeOf(src)
	if dt == nil || st == nil || dt.Kind() != reflect.Pointer || st.Kind() != reflect.Pointer {
		return errors.New("dst和src需要是两个结构体的指针")
	}

	if dt.Elem().Kind() != reflect.Struct || st.Elem().Kind() != reflect.Struct {
		return errors.New("dst和src需要是两个结构体的指针")
	}

	err := copyStruct(reflect.ValueOf(dst), reflect.ValueOf(src))
	if err != nil {
		return err
	}
	return nil
}

var cpi = reflect.TypeOf((*CopyInterface)(nil)).Elem()

func copyStruct(dstValue, srcValue reflect.Value) error {

	dst := dstValue
	src := srcValue
	if dst.Kind() == reflect.Pointer {
		dst = dst.Elem()
	}
	if src.Kind() == reflect.Pointer {
		src = src.Elem()
	}

	if dst.Kind() != reflect.Struct || src.Kind() != reflect.Struct {
		return errors.New("copyStruct argus LiveType must be struct or struct point")
	}

	if reflect.PointerTo(dst.Type()).Implements(cpi) {
		va := dst
		if dst.CanAddr() {
			va = dst.Addr()
			if va.IsNil() {
				return errors.New("dst addr is nil")
			}
		}

		cp, ok := va.Interface().(CopyInterface)
		if ok {
			err := cp.Copy(src)
			if err != nil {
				return err
			}
			return nil
		}
	}

	fieldCount := dst.Type().NumField()

	for i := 0; i < fieldCount; i++ {
		dstTypField := dst.Type().Field(i)

		srcValField := src.FieldByName(dstTypField.Name)
		if !srcValField.IsValid() || srcValField.IsZero() {
			continue
		}

		//TODO 这里可以考虑下如果int8 、int16等不一样之间的的copy
		if !typeEquals(dstTypField.Type, srcValField.Type()) {
			continue
		}

		dstValField := dst.FieldByName(dstTypField.Name)

		if dstValField.IsValid() && dstValField.CanSet() {
			err := copyField(dstValField, srcValField)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func copyField(dstValue, srcValue reflect.Value) error {
	if srcValue.IsZero() {
		return nil
	}
	dst := dstValue
	if dstValue.Kind() == reflect.Pointer {
		dst = dstValue.Elem()
	}
	src := srcValue
	if dst.Kind() != src.Kind() {
		return errors.New(fmt.Sprintf("copyBase %s not equals %s", dst.Kind().String(), src.Kind().String()))
	}
	switch dst.Kind() {
	case reflect.String:
		dst.SetString(src.String())
	case reflect.Int8, reflect.Int16, reflect.Int, reflect.Int32, reflect.Int64,
		reflect.Uint8, reflect.Uint16, reflect.Uint, reflect.Uint32, reflect.Uint64:
		dst.SetInt(src.Int())
	case reflect.Float64, reflect.Float32:
		dst.SetFloat(src.Float())
	case reflect.Bool:
		dst.SetBool(src.Bool())
	case reflect.Interface:
		dst.Set(src)
	case reflect.Struct:
		err := copyStruct(dst, src)
		if err != nil {
			return err
		}
	case reflect.Slice:
		s := reflect.AppendSlice(dst, src)
		dst.Set(s)
	case reflect.Map:
		for _, k := range src.MapKeys() {
			v := src.MapIndex(k)
			if !v.IsValid() {
				continue
			}
			if !dst.IsNil() {
				// 如果dst可以寻址，则直接使用地址设置
				if dst.CanAddr() {
					dst.SetMapIndex(k, v)
				} else {
					// 如果dst不能寻址，则先创建一个新映射再拷贝
					if dst.IsNil() {
						// 如果dst是nil，则应该是赋值操作失败，应该返回错误
						return errors.New("dst is nil")
					}
					// 先创建一个与dst相同大小的新映射
					newMap := reflect.MakeMapWithSize(dst.Type(), src.Len())
					for _, innerK := range src.MapKeys() {
						newMap.SetMapIndex(innerK, src.MapIndex(innerK))
					}
					// 再将新映射赋值给dst
					dst.Set(newMap)
				}
			} else {
				// 如果dst是nil，则创建一个新映射并开始填充
				dst.Set(reflect.MakeMapWithSize(dst.Type(), src.Len()))
				for _, innerK := range src.MapKeys() {
					dst.SetMapIndex(innerK, src.MapIndex(innerK))
				}
			}
		}
	default:
		return errors.New(fmt.Sprintf("copyBase not support the type %s", dst.Kind().String()))
	}
	return nil
}

func getElem(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Pointer {
		return v.Elem()
	} else {
		return v
	}
}

func typeEquals(dst, src reflect.Type) bool {

	d := dst
	s := src
	if dst.Kind() == reflect.Pointer {
		d = dst.Elem()
	}
	if src.Kind() == reflect.Pointer {
		s = src.Elem()
	}
	return d.Kind() == s.Kind()
}

type CopyInterface interface {
	Copy(src reflect.Value) error
}
