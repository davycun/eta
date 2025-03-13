package faker

import (
	"testing"
)

func TestIdCardOf(t *testing.T) {
	for i := 0; i < 10; i++ {
		birthDay := BirthDayFmt()
		gender := Gender()
		println(IdCardOf(birthDay, gender))
	}

	//rn := Number(0, 365*105) //最大年龄105岁
	//now := time.Now()
	//now = now.Truncate(time.Hour * 24)
	//rt := now.Add(-time.Hour * 24 * time.Duration(rn))
	//println(fmt.Sprintf("%v", rt))

}
