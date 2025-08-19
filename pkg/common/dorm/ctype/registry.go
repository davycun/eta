package ctype

var (
	scannerMap = map[string]Scanner{}
)

type Scanner func(src any) (any, error)

func RegistryScanner(tp string, fn Scanner) {
	scannerMap[tp] = fn
}

func RemoveScanner(tp string) {
	delete(scannerMap, tp)
}
func GetScanner(tp string) (Scanner, bool) {
	s, ok := scannerMap[tp]
	return s, ok
}
