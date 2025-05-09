package faker

import (
	"fmt"
	"testing"
)

func ExampleNumber() {
	Seed(11)
	fmt.Println(Number(50, 23456))
	// 14866
}

func BenchmarkNumber(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Number(10, 999999)
	}
}

func ExampleUint8() {
	Seed(11)
	fmt.Println(Uint8())
	// 152
}

func BenchmarkUint8(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Uint8()
	}
}

func ExampleUint16() {
	Seed(11)
	fmt.Println(Uint16())
	// 34968
}

func BenchmarkUint16(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Uint16()
	}
}

func ExampleUint32() {
	Seed(11)
	fmt.Println(Uint32())
	// 1075055705
}

func BenchmarkUint32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Uint32()
	}
}

func ExampleUint64() {
	Seed(11)
	fmt.Println(Uint64())
	// 843730692693298265
}

func BenchmarkUint64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Uint64()
	}
}

func ExampleInt8() {
	Seed(11)
	fmt.Println(Int8())
	// 24
}

func BenchmarkInt8(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Int8()
	}
}

func ExampleInt16() {
	Seed(11)
	fmt.Println(Int16())
	// 2200
}

func BenchmarkInt16(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Int16()
	}
}

func ExampleInt32() {
	Seed(11)
	fmt.Println(Int32())
	// -1072427943
}

func BenchmarkInt32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Int32()
	}
}

func ExampleInt64() {
	Seed(11)
	fmt.Println(Int64())
	// -8379641344161477543
}

func BenchmarkInt64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Int64()
	}
}

func ExampleFloat32() {
	Seed(11)
	fmt.Println(Float32())
	// 3.1128167e+37
}

func BenchmarkFloat32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Float32()
	}
}

func ExampleFloat64() {
	Seed(11)
	fmt.Println(Float64())
	// 1.644484108270445e+307
}

func BenchmarkFloat64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Float64()
	}
}

func ExampleNumerify() {
	Seed(11)
	fmt.Println(Numerify("###-###-####"))
	// 328-727-1570
}

func BenchmarkNumerify(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Numerify("###-###-####")
	}
}

func ExampleShuffleInts() {
	Seed(11)

	ints := []int{52, 854, 941, 74125, 8413, 777, 89416, 841657}
	ShuffleInts(ints)
	fmt.Println(ints)
	// [74125 777 941 89416 8413 854 52 841657]
}

func BenchmarkShuffleInts(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ShuffleInts([]int{52, 854, 941, 74125, 8413, 777, 89416, 841657})
	}
}
