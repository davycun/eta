package faker

import (
	"testing"
)

func TestName(t *testing.T) {
	for i := 0; i < 100; i++ {
		println(Name())
	}
}
