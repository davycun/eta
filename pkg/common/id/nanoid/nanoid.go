package nanoid

import "github.com/matoous/go-nanoid/v2"

var defaultAlphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

const (
	defaultSize = 21
)

func New() string {
	s, err := gonanoid.Generate(defaultAlphabet, defaultSize)
	if err != nil {
		return ""
	}
	return s
}
