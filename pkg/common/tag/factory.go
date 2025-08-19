package tag

import "strings"

const (
	GormTagName    = "gorm"
	EsTagName      = "es"
	JsonTagName    = "json"
	BindingTagName = "binding"
	DorisTagName   = "doris"
)

func NewTag(tagName, keySplit, keyValueSplit string) *Tag {
	return &Tag{
		name:          tagName,
		keySplit:      keySplit,
		keyValueSplit: keyValueSplit,
		keys:          make([]string, 0),
		props:         make(map[string]string),
	}
}

func ParseTag(tagName string, text string, keySplit string, keyValueSplit string) *Tag {
	var (
		t = &Tag{
			name:          tagName,
			keySplit:      keySplit,
			keyValueSplit: keyValueSplit,
			props:         make(map[string]string),
		}
	)
	if t.keySplit == "" {
		t.keySplit = SplitSemicolon
	}
	if t.keyValueSplit == "" {
		t.keyValueSplit = SplitColon
	}
	if text == "" {
		return t
	}
	ts := strings.Split(text, t.keySplit)
	for _, v := range ts {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		i2 := strings.Split(v, t.keyValueSplit)
		if len(i2) < 1 {
			continue
		}
		switch len(i2) {
		case 1:
			tmp := strings.TrimSpace(i2[0])
			t.props[tmp] = ""
			t.keys = append(t.keys, tmp)
		default:
			tmp := strings.TrimSpace(i2[0])
			t.keys = append(t.keys, tmp)
			t.props[tmp] = strings.TrimSpace(strings.Join(i2[1:], ""))
		}
	}
	return t
}

func NewBindingTag() *Tag {
	return NewTag(BindingTagName, SplitComma, SplitEq)
}
func NewGormTag() *Tag {
	return NewTag(GormTagName, SplitSemicolon, SplitColon)
}
func NewEsTag() *Tag {
	return NewTag(EsTagName, SplitSemicolon, SplitColon)
}
func NewDorisTag() *Tag {
	return NewTag(EsTagName, SplitSemicolon, SplitColon)
}
func NewJsonTag() *Tag {
	return NewTag(EsTagName, SplitComma, "")
}

func ParseBindingTag(text string) *Tag {
	return ParseTag(BindingTagName, text, SplitComma, SplitEq)
}
func ParseGormTag(text string) *Tag {
	return ParseTag(GormTagName, text, SplitSemicolon, SplitColon)
}
func ParseEsTag(text string) *Tag {
	return ParseTag(EsTagName, text, SplitSemicolon, SplitColon)
}
func ParseDorisTag(text string) *Tag {
	return ParseTag(DorisTagName, text, SplitSemicolon, SplitColon)
}
func ParseJsonTag(text string) *Tag {
	return ParseTag(JsonTagName, text, SplitComma, "")
}
