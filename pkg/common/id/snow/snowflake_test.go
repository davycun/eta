package snow_test

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/id/snow"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	mp        = make(map[int64]int64, 10000000)
	snowflake = snow.NewSnowflake(1, "2021-02-02")
	//mp = sync.Map{}
)

func BenchmarkGenerate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		j := snowflake.Generate()
		_, ok := mp[j]
		if ok {
			b.Error("same id error")
			fmt.Printf("map length %d\n", len(mp))
			fmt.Printf("conflict id %d\n", j)
		}
		mp[j] = j
	}
}

func TestNodeId(t *testing.T) {
	assert.Equal(t, 1023, snow.NodeMax)
	assert.Equal(t, 4095, snow.SequenceMax)
}

func TestGen(t *testing.T) {
	for i := 0; i < 100; i++ {
		fmt.Printf("%d,", snowflake.Generate())
	}
}
func TestDefault(t *testing.T) {
	assert.True(t, snow.DefaultSnowflake().Generate() > 0)
}
