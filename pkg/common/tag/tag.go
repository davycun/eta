package tag

import (
	"github.com/duke-git/lancet/v2/slice"
	"strings"
)

const (
	SplitSemicolon = ";"
	SplitColon     = ":"
	SplitComma     = ","
	SplitEq        = "="
)

type Tag struct {
	name          string
	keySplit      string   //多个key和value 对之间的分隔符，一般是分号";"
	keyValueSplit string   //key和value之间分隔符，一般是冒号":"，也可能等号"="
	keys          []string // 为了保障顺序，这里存的是props的key，比如json这样的tag，第一个必须是名字
	props         map[string]string
}

func (t *Tag) Exists(key string) bool {
	if t.props == nil {
		t.props = make(map[string]string)
	}
	_, ok := t.props[key]
	return ok
}

func (t *Tag) Get(key string) string {
	if t.props == nil {
		t.props = make(map[string]string)
	}
	v, ok := t.props[key]
	if !ok {
		return ""
	}
	if v == "" {
		v = key
	}
	return v
}

// GetFirstKey
// 主要是为了json这个tag的，返回名字用
func (t *Tag) GetFirstKey() string {
	if t.props == nil {
		t.props = make(map[string]string)
	}
	if len(t.keys) < 1 {
		return ""
	}
	return t.keys[0]
}
func (t *Tag) GetTagName() string {
	return t.name
}
func (t *Tag) GetAll() map[string]string {
	mp := make(map[string]string)
	for k, v := range t.props {
		mp[k] = v
	}
	return mp
}
func (t *Tag) GetArray(key string) []string {
	return t.GetArrayWithSplit(key, SplitComma)
}
func (t *Tag) GetArrayWithSplit(key string, split string) []string {
	if t.props == nil {
		t.props = make(map[string]string)
	}
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

func (t *Tag) SetKeySplit(split string) *Tag {
	t.keySplit = split
	return t
}
func (t *Tag) SetKeyValueSplit(split string) *Tag {
	t.keyValueSplit = split
	return t
}

func (t *Tag) Add(key string, value string) *Tag {
	if key == "" {
		key = value
	}
	if key == "" {
		return t
	}
	if t.props == nil {
		t.props = make(map[string]string)
	}
	t.keys = append(t.keys, key)
	t.props[key] = value
	return t
}
func (t *Tag) Remove(key string) *Tag {
	delete(t.props, key)
	t.keys = slice.Filter(t.keys, func(index int, item string) bool {
		return item != key
	})
	return t
}
func (t *Tag) String() string {
	bd := strings.Builder{}
	if t.name == "" {
		return ""
	}
	if len(t.props) < 1 {
		return t.name + `:""`
	}
	for _, k := range t.keys {
		if k == "" {
			continue
		}
		v := t.props[k]
		if bd.Len() > 0 {
			bd.WriteString(t.keySplit)
		}
		bd.WriteString(k)
		if t.keyValueSplit != "" && v != "" {
			bd.WriteString(t.keyValueSplit)
			bd.WriteString(v)
		}
	}

	return t.name + `:"` + bd.String() + `"`
}
