package faker

import (
	"fmt"
	"testing"
)

func ExampleLetter() {
	Seed(11)
	fmt.Println(Letter())
	// g
}

func BenchmarkLetter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Letter()
	}
}

func ExampleLexify() {
	Seed(11)
	fmt.Println(Lexify("?????"))
	// gbrma
}

func BenchmarkLexify(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Lexify("??????")
	}
}

func ExampleShuffleStrings() {
	Seed(11)
	strings := []string{"happy", "times", "for", "everyone", "have", "a", "good", "day"}
	ShuffleStrings(strings)
	fmt.Println(strings)
	// [everyone a for good have times happy day]
}

func BenchmarkShuffleStrings(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ShuffleStrings([]string{"happy", "times", "for", "everyone", "have", "a", "good", "day"})
	}
}
