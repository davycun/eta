package utils

import (
	"fmt"
	"sort"
)

// foreground background color
// ----------------------------
// 30  40  黑色
// 31  41  红色
// 32  42  绿色
// 33  43  黄色
// 34  44  蓝色
// 35  45  紫红色
// 36  46  青蓝色
// 37  47  白色
//
// style 说明
// -------------------------
//  0  终端默认设置
//  1  高亮显示
//  4  使用下划线
//  5  闪烁
//  7  反白显示
//  8  不可见

const (
	TextBlack = iota + 30
	TextRed
	TextGreen
	TextYellow
	TextBlue
	TextMagenta
	TextCyan
	TextWhite
)

func FmtUrl(url string) string {
	return FmtColor(url, 4, 0, 0)
}

func FmtTextBlack(msg string) string {
	return FmtColor(msg, 0, TextBlack, 0)
}
func FmtTextRed(msg string) string {
	return FmtColor(msg, 0, TextRed, 0)
}
func FmtTextGreen(msg string) string {
	return FmtColor(msg, 0, TextGreen, 0)
}
func FmtTextYellow(msg string) string {
	return FmtColor(msg, 0, TextYellow, 0)
}
func FmtTextBlue(msg string) string {
	return FmtColor(msg, 0, TextBlue, 0)
}
func FmtTextMagenta(msg string) string {
	return FmtColor(msg, 0, TextMagenta, 0)
}
func FmtTextCyan(msg string) string {
	return FmtColor(msg, 0, TextCyan, 0)
}
func FmtTextWhite(msg string) string {
	return FmtColor(msg, 0, TextWhite, 0)
}

func FmtTextRGB(msg string, style, r, g, b int) string {
	return fmt.Sprintf("\033[%d;38;2;%d;%d;%dm%s\033[0m", style, r, g, b, msg)
}

func FmtBackRGB(msg string, style, r, g, b int) string {
	return fmt.Sprintf("\033[%d;48;2;%d;%d;%dm%s\033[0m", style, r, g, b, msg)
}

func FmtColor(msg string, style, foreground, background int) string {
	cs := checkColor(style, foreground, background)
	return fmt.Sprintf("\033[%d;%d;%dm%s\033[0m", cs[0], cs[1], cs[2], msg)
}

func checkColor(style, foreground, background int) (cs []int) {
	cs = []int{checkStyle(style), checkForegroundColor(foreground), checkBackgroundColor(background)}
	sort.Ints(cs)
	return
}

func checkStyle(style int) (c int) {
	c = style
	if style > 7 || style < 0 {
		c = 0
	}
	return
}
func checkForegroundColor(color int) (c int) {
	c = color
	if color > 37 || color < 30 {
		c = 0
	}
	return
}
func checkBackgroundColor(color int) (c int) {
	c = color
	if color > 47 || color < 40 {
		c = 0
	}
	return
}

func checkRGB(r, g, b int) (r1, g1, b1 int) {
	if r < 0 || r > 255 {
		r1 = 0
	}
	if g < 0 || g > 255 {
		g1 = 0
	}
	if b < 0 || b > 255 {
		b1 = 0
	}
	return
}

func FmtBool(b bool, trueStr, falseStr string) string {
	if trueStr == "" {
		trueStr = "true"
	}
	if falseStr == "" {
		falseStr = "false"
	}
	if b {
		return trueStr
	}
	return falseStr
}
