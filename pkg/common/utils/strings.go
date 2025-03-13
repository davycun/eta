package utils

import "strings"

func HumpToUnderline(src string) string {
	bd := strings.Builder{}
	pre := 0
	mp := make(map[int]int32)
	for j, v := range src {
		mp[j] = v
		if v >= 'A' && v <= 'Z' && j > 0 {
			p := mp[pre]
			if (p <= 'A' || p >= 'Z') && ((p >= '0' && p <= '9') || (p >= 'a' && p <= 'z')) {
				bd.WriteByte('_')
			}
		}
		bd.WriteString(strings.ToLower(string(v)))
		pre = j
	}
	return bd.String()
}
func UnderlineToHump(src string) string {
	//特殊处理下ID
	if strings.ToLower(src) == "id" {
		return "ID"
	}

	bd := strings.Builder{}
	pre := false
	for j, v := range src {
		if v == '_' {
			pre = true
			continue
		}

		//首字母大写
		//TODO 可以增加一个参数考虑是否首字母大写
		if pre || j == 0 {
			bd.WriteString(strings.ToUpper(string(v)))
		} else {
			bd.WriteRune(v)
		}
		pre = false
	}
	return bd.String()
}

func Split(src string, sep ...string) []string {
	if src == "" {
		return []string{}
	}
	if len(sep) < 1 {
		return []string{src}
	}
	var (
		rs = make([]string, 0, 2)
		bd = strings.Builder{}
	)
	for _, v := range src {
		if ContainAny(sep, string(v)) {
			tmp := bd.String()
			bd.Reset()
			if strings.TrimSpace(tmp) != "" {
				rs = append(rs, tmp)
			}
			continue
		}
		bd.WriteRune(v)
	}
	if bd.Len() > 0 {
		rs = append(rs, bd.String())
	}
	return rs
}

func AddPrefix(prefix string, src ...string) []string {
	rs := make([]string, 0, len(src))
	for _, v := range src {
		rs = append(rs, prefix+v)
	}
	return rs
}
