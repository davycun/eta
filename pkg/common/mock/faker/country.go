package faker

// Country will generate a random country string
func Country(langs ...string) string {
	// lang: zh_CN --> 中文
	// lang: en_US --> 英文
	lang := randLang(langs)
	return getRandValue([]string{"country", lang})
}
