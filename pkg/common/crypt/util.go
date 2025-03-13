package crypt

func StringToSlice(text string, sliceSize int) []string {
	if sliceSize < 1 {
		return []string{text}
	}
	var (
		msg        = []rune(text)
		length     = len(msg)
		encStrList = make([]string, 0, length/sliceSize+1)
	)
	if sliceSize > length {
		encStrList = append(encStrList, string(msg))
	} else {
		for i, _ := range msg {
			if i+sliceSize > length {
				break
			}
			encStrList = append(encStrList, string(msg[i:i+sliceSize]))
		}
	}
	return encStrList
}

func SliceToString(slc []string) string {
	length := len(slc)
	if length == 0 {
		return ""
	}
	var res []rune
	for i, m := range slc {
		if i == length-1 {
			res = append(res, []rune(m)...)
		} else if len([]rune(m)) > 0 {
			res = append(res, []rune(m)[0])
		}
	}
	return string(res)
}
