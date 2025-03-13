package faker

import (
	"math/rand"
	"strings"
)

// Letter will generate a single random lower case letter
func Letter() string {
	return RandLetter()
}

// Lexify will replace ? will random generated letters
func Lexify(str string) string {
	return replaceWithLetters(str)
}

// ShuffleStrings will randomize a slice of strings
func ShuffleStrings(a []string) {
	for i := range a {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}
}

// RandString will take in a slice of string and return a randomly selected value
func RandString(a []string) string {
	return a[rand.Intn(len(a))]
}

// RandStrings will take in a slice of string and return randomly selected value
func RandStrings(a []string, l int) string {
	sb := strings.Builder{}
	for i := 0; i < l; i++ {
		sb.WriteString(RandString(a))
	}
	return sb.String()
}

// RandBool will take in a slice of bool and return a randomly selected value
func RandBool(a []bool) bool {
	return a[rand.Intn(len(a))]
}
