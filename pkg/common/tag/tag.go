package tag

import "strings"

const (
	SplitSemicolon = ";"
	SplitColon     = ":"
	SplitComma     = ","
)

type Tag struct {
	Text  string
	split string
	props map[string]string
}

// New 这个方式创建的tag的格式是key:value;key:value;key:value
func New(text string) Tag {
	return NewWithSplit(text, SplitSemicolon)
}
func NewWithSplit(text string, split string) Tag {

	if text == "" {
		return Tag{}
	}
	if split == "" {
		split = SplitSemicolon
	}
	tg := Tag{Text: text, split: split, props: make(map[string]string)}

	ts := strings.Split(text, split)

	for _, v := range ts {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		i2 := strings.Split(v, SplitColon)
		if len(i2) < 1 {
			continue
		}
		switch len(i2) {
		case 1:
			tmp := strings.TrimSpace(i2[0])
			tg.props[tmp] = tmp
		default:
			tmp1 := strings.TrimSpace(i2[0])
			tg.props[tmp1] = strings.TrimSpace(strings.Join(i2[1:], ""))
		}
	}
	return tg
}

func (t Tag) Get(key string) string {
	return t.props[key]
}
func (t Tag) GetAll() map[string]string {
	mp := make(map[string]string)
	for k, v := range t.props {
		mp[k] = v
	}
	return mp
}
func (t Tag) GetArray(key string) []string {
	return t.GetArrayWithSplit(key, SplitComma)
}
func (t Tag) GetArrayWithSplit(key string, split string) []string {
	val := t.props[key]
	if val == "" {
		return make([]string, 0)
	}
	rs := strings.Split(val, split)
	for i, v := range rs {
		rs[i] = strings.TrimSpace(v)
	}
	return rs
}

type JsonTag struct {
	Tag
}

func NewJsonTag(text string) JsonTag {
	jt := JsonTag{}
	jt.Tag = NewWithSplit(text, SplitComma)
	return jt
}
func (t JsonTag) GetName() string {
	for k, v := range t.GetAll() {
		if k == "omitempty" {
			continue
		}
		return v
	}
	return ""
}
