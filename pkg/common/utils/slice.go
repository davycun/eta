package utils

func ConvertToSlice[T any](data any, dst *[]T) {

	var (
		dt []T
	)
	switch x := data.(type) {
	case *[]T:
		dt = *(x)
	case []T:
		dt = x
	case *T:
		dt = append(dt, *(x))
	case T:
		dt = append(dt, x)
	}
	*dst = dt
}

func AppendNoEmpty(src []string, dst ...string) []string {
	if len(dst) < 1 {
		return src
	}
	for _, v := range dst {
		if v != "" {
			src = append(src, v)
		}
	}
	return src
}
