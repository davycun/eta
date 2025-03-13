package faker

import (
	"fmt"
	"testing"
)

func ExampleOneColor() {
	Seed(11)
	fmt.Println(Color("zh_CN"))
	// 黑色
}

func ExampleRandomColor() {
	Seed(11)
	fmt.Println(Color("zh_CN", "en_US"))
	// 黑色
}

func BenchmarkOneColor(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Color("zh_CN")
	}
}
