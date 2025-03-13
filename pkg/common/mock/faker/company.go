package faker

// Company will generate a random company name string
func Company(langs ...string) (company string) {
	lang := randLang(langs)
	if lang == "zh_CN" {
		return getRandValue([]string{"company", lang + "_prefix"}) + getRandValue([]string{"company", lang + "_suffix"})
	}

	switch randInt := RandIntRange(1, 3); randInt {
	case 1:
		company = LastName() + ", " + LastName() + " and " + LastName()
	case 2:
		company = LastName() + "-" + LastName()
	case 3:
		company = LastName() + " " + CompanySuffix()
	}

	return
}

// CompanySuffix will generate a random company suffix string
func CompanySuffix() string {
	return getRandValue([]string{"company", "suffix"})
}

// BuzzWord will generate a random company buzz word string
func BuzzWord() string {
	return getRandValue([]string{"company", "buzzwords"})
}

// BS will generate a random company bs string
func BS() string {
	return getRandValue([]string{"company", "bs"})
}
