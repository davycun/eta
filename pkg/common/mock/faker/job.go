package faker

func Job(langs ...string) string {
	// lang: zh_CN --> 中文
	// lang: en_US --> 英文
	lang := randLang(langs)
	return getRandValue([]string{"job", lang})
}
