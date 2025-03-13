package faker

import (
	"testing"
)

func TestTelPhone(t *testing.T) {
	for i := 0; i < 10; i++ {
		println(TelPhoneShort())
	}
}
