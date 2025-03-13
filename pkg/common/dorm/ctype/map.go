package ctype

import (
	"encoding/xml"
	"github.com/davycun/eta/pkg/common/utils"
	"strings"
)

type Map map[string]interface{}

// MarshalXML allows type H to be used with xml.Marshal.
func (m Map) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{
		Space: "",
		Local: "map",
	}
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	for key, value := range m {
		elem := xml.StartElement{
			Name: xml.Name{Space: "", Local: key},
			Attr: []xml.Attr{},
		}
		if err := e.EncodeElement(value, elem); err != nil {
			return err
		}
	}

	return e.EncodeToken(xml.EndElement{Name: start.Name})
}

func (m *Map) GeomFormat(gcsType string, geoFormat string) {
	for fk, fv := range *m {
		switch g := fv.(type) {
		case *Geometry:
			g.GeomFormat(gcsType, geoFormat)
		case Geometry:
			g.GeomFormat(gcsType, geoFormat)
			(*m)[fk] = g
		}
	}
}

// GetKey
// 实现KeyInterface接口
// 可以通过GetKey来获取指定字段内容
func (m Map) GetKey(field string) string {
	if v, ok := m[strings.ToLower(field)]; ok {
		return ToString(v)
	}
	return ""
}

// GetId
// 下面四个方法是实现了TreeEntity接口
func (m Map) GetId() string {
	return m.GetKey("id")
}

func (m Map) GetParentId() string {
	return m.GetKey("parent_id")
}

func (m Map) GetChildren() any {
	return m["children"]
}

func (m Map) SetChildren(a any) {
	m["children"] = a
}

func (m Map) Get(key string) (interface{}, bool) {
	x, ok := m[key]
	if !ok {
		key = utils.HumpToUnderline(key)
		x, ok = m[key]
	}
	return x, ok
}
func (m *Map) Set(key string, val interface{}) {
	(*m)[key] = val
}
