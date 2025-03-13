package faker

// Color will generate a random color string
func Gender(langs ...string) string {
	// lang: zh_CN --> 中文
	// lang: en_US --> 英文
	lang := randLang(langs)
	return getRandValue([]string{"gender", lang})
}
