package faker

import "strconv"

// UserAgent will generate a random broswer user agent
func UserAgent() string {
	randNum := RandIntRange(0, 4)
	switch randNum {
	case 0:
		return ChromeUserAgent()
	case 1:
		return FirefoxUserAgent()
	case 2:
		return SafariUserAgent()
	case 3:
		return OperaUserAgent()
	default:
		return ChromeUserAgent()
	}
}

// ChromeUserAgent will generate a random chrome browser user agent string
func ChromeUserAgent() string {
	randNum1 := strconv.Itoa(RandIntRange(531, 536)) + strconv.Itoa(RandIntRange(0, 2))
	randNum2 := strconv.Itoa(RandIntRange(36, 40))
	randNum3 := strconv.Itoa(RandIntRange(800, 899))
	return "Mozilla/5.0 " + "(" + randomPlatform() + ") AppleWebKit/" + randNum1 + " (KHTML, like Gecko) Chrome/" + randNum2 + ".0." + randNum3 + ".0 Mobile Safari/" + randNum1
}

// FirefoxUserAgent will generate a random firefox broswer user agent string
func FirefoxUserAgent() string {
	ver := "Gecko/" + Date().Format("2006-02-01") + " Firefox/" + strconv.Itoa(RandIntRange(35, 37)) + ".0"
	platforms := []string{
		"(" + windowsPlatformToken() + "; " + "en-US" + "; rv:1.9." + strconv.Itoa(RandIntRange(0, 3)) + ".20) " + ver,
		"(" + linuxPlatformToken() + "; rv:" + strconv.Itoa(RandIntRange(5, 8)) + ".0) " + ver,
		"(" + macPlatformToken() + " rv:" + strconv.Itoa(RandIntRange(2, 7)) + ".0) " + ver,
	}

	return "Mozilla/5.0 " + RandString(platforms)
}

// SafariUserAgent will generate a random safari browser user agent string
func SafariUserAgent() string {
	randNum := strconv.Itoa(RandIntRange(531, 536)) + "." + strconv.Itoa(RandIntRange(1, 51)) + "." + strconv.Itoa(RandIntRange(1, 8))
	ver := strconv.Itoa(RandIntRange(4, 6)) + "." + strconv.Itoa(RandIntRange(0, 2))

	mobileDevices := []string{
		"iPhone; CPU iPhone OS",
		"iPad; CPU OS",
	}

	platforms := []string{
		"(Windows; U; " + windowsPlatformToken() + ") AppleWebKit/" + randNum + " (KHTML, like Gecko) Version/" + ver + " Safari/" + randNum,
		"(" + macPlatformToken() + " rv:" + strconv.Itoa(RandIntRange(4, 7)) + ".0; en-US) AppleWebKit/" + randNum + " (KHTML, like Gecko) Version/" + ver + " Safari/" + randNum,
		"(" + RandString(mobileDevices) + " " + strconv.Itoa(RandIntRange(7, 9)) + "_" + strconv.Itoa(RandIntRange(0, 3)) + "_" + strconv.Itoa(RandIntRange(1, 3)) + " like Mac OS X; " + "en-US" + ") AppleWebKit/" + randNum + " (KHTML, like Gecko) Version/" + strconv.Itoa(RandIntRange(3, 5)) + ".0.5 Mobile/8B" + strconv.Itoa(RandIntRange(111, 120)) + " Safari/6" + randNum,
	}

	return "Mozilla/5.0 " + RandString(platforms)
}

// OperaUserAgent will generate a random opera browser user agent string
func OperaUserAgent() string {
	platform := "(" + randomPlatform() + "; en-US) Presto/2." + strconv.Itoa(RandIntRange(8, 13)) + "." + strconv.Itoa(RandIntRange(160, 355)) + " Version/" + strconv.Itoa(RandIntRange(10, 13)) + ".00"

	return "Opera/" + strconv.Itoa(RandIntRange(8, 10)) + "." + strconv.Itoa(RandIntRange(10, 99)) + " " + platform
}

// linuxPlatformToken will generate a random linux platform
func linuxPlatformToken() string {
	return "X11; Linux " + getRandValue([]string{"computer", "linux_processor"})
}

// macPlatformToken will generate a random mac platform
func macPlatformToken() string {
	return "Macintosh; " + getRandValue([]string{"computer", "mac_processor"}) + " Mac OS X 10_" + strconv.Itoa(RandIntRange(5, 9)) + "_" + strconv.Itoa(RandIntRange(0, 10))
}

// windowsPlatformToken will generate a random windows platform
func windowsPlatformToken() string {
	return getRandValue([]string{"computer", "windows_platform"})
}

// randomPlatform will generate a random platform
func randomPlatform() string {
	platforms := []string{
		linuxPlatformToken(),
		macPlatformToken(),
		windowsPlatformToken(),
	}

	return RandString(platforms)
}
