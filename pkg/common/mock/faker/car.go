package faker

// 获取所有车辆品牌的名称
func CarBrand(langs ...string) string {
	// lang: zh_CN --> 中文
	// lang: en_US --> 英文
	lang := randLang(langs)
	return getRandValue([]string{"carbrand", lang})
}
