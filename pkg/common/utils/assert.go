package utils

func AssertNil(err error) {
	if err != nil {
		panic(err.Error())
	}
}
