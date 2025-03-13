package bean_test

import (
	"encoding/json"
	"fmt"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/utils/bean"
	"reflect"
	"testing"
	"time"
)

func TestCp(t *testing.T) {

	type B struct {
		BA string
	}

	type Dst struct {
		Name string
		Age  int
		//High bool
		Tm  *ctype.LocalTime
		Son struct {
			A string
		}
	}

	type Src struct {
		Name string
		Age  int
		High bool
		Tm   ctype.LocalTime
		Son  struct {
			A string
		}
	}

	var dt Dst

	sc := Src{
		Name: "Davy",
		Age:  23,
		High: true,
		Tm:   ctype.LocalTime{Data: time.Now(), Valid: true},
		Son:  struct{ A string }{A: "ps"},
	}

	err := bean.Copy(&dt, &sc)
	println(err)

	//cpi := reflect.TypeOf((*CopyInterface)(nil)).Elem()
	//tf := reflect.TypeOf(ctype.LocalTime{})
	//tf1 := reflect.TypeOf(&ctype.LocalTime{})
	//
	//println(reflect.PointerTo(tf).Implements(cpi))
	//println(reflect.PointerTo(tf).AssignableTo(cpi))
	//println(reflect.PointerTo(tf1).Implements(cpi))
	//println(reflect.PointerTo(tf1).AssignableTo(cpi))

	//i := reflect.Copy(, reflect.ValueOf(sc))
	//reflect.ValueOf(&dt).Elem().Set(reflect.ValueOf(sc))
	//println("")

}

func TestImpl(t *testing.T) {

	lt := ctype.LocalTime{Data: time.Now(), Valid: true}

	cpi := reflect.TypeOf((*bean.CopyInterface)(nil)).Elem()
	dt := reflect.TypeOf(lt)
	vo := reflect.ValueOf(lt)

	fmt.Printf("cpi_name: %s, cpi_kind: %d\n", cpi.Name(), cpi.Kind())
	fmt.Printf("dt_name: %s, dt_kind: %d\n", dt.Name(), dt.Kind())
	fmt.Printf("vot_name: %s, vo_kind: %d, vot_kind: %d\n", vo.Type().Name(), vo.Kind(), vo.Type().Kind())

	if vo.Type().Kind() == reflect.Pointer {
		println("==========")
		println(vo.Type().Implements(cpi))
		println(vo.Type().Kind())
		println(vo.Type().Name())
		println(vo.Interface())
		localTime := vo.Interface().(bean.CopyInterface)
		println(localTime)
	} else {
		println("----------")
		println(vo.Type().Kind())
		println(reflect.PointerTo(vo.Type()).Implements(cpi))
		println(vo.CanAddr())
		println(vo.CanInterface())
		println(vo.CanConvert(cpi))
		//println(vo.UnsafeAddr())
		//v := vo.Convert(cpi)
		//vo.UnsafePointer()

		//println(vo.Interface())
		//println(vo.CanConvert(cpi))
		//i := vo.Addr().Interface().(CopyInterface)
		//i := v.Interface().(CopyInterface)
		//println(i)
	}

	//
	//
	//println(dt.Implements(cpi))
}

func TestJson(t *testing.T) {

	j := JS{Name: "PS"}

	marshal, err := json.Marshal(&j)
	if err != nil {
		println(err.Error())
	}
	println(string(marshal))
}

type JS struct {
	Name string
}

func (j *JS) MarshalJSON() ([]byte, error) {
	return []byte("{\"Name\":\"PS\"}"), nil
}
