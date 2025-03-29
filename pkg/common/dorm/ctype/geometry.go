package ctype

import (
	"context"
	"fmt"
	"github.com/davycun/dm8-go-driver"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/utils"
	jsoniter "github.com/json-iterator/go"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"github.com/twpayne/go-geom/encoding/ewkbhex"
	"github.com/twpayne/go-geom/encoding/geojson"
	"github.com/twpayne/go-geom/encoding/wkbhex"
	"github.com/twpayne/go-geom/encoding/wkt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"strings"
)

const (
	GeomMarshalWkb     = "wkb"
	GeomMarshalEWkb    = "ewkb"
	GeomMarshalWkt     = "wkt"
	GeomMarshalGeoJson = "geojson"
)

type Geometry struct {
	Data       geom.T
	Srid       int
	GeoTypeId  int
	Valid      bool
	FormatType string //wkb,wkt,geojson
	GcsType    string //坐标系类型，gcj02,wgs84,bd09,cgs2000,.....
}

func (g Geometry) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	if !g.Valid || g.Data == nil {
		return clause.Expr{SQL: "null"}
	}
	//encode, err := ewkb.Marshal(g.Data, wkbhex.XDR)
	var (
		dbType      = dorm.GetDbType(db)
		encode, err = ewkbhex.Encode(g.Data, wkbhex.XDR)
	)
	if err != nil {
		logger.Errorf("Geometry.FormValue err on ekbhex.Encode %s", err)
	}
	switch dbType {
	case dorm.DaMeng:
		srid := g.Data.SRID()
		if srid == 0 {
			srid = g.Srid
		}
		return clause.Expr{
			SQL:  "dmgeo.ST_GeomFromWKB(?,?)",
			Vars: []interface{}{encode, srid},
		}
	case dorm.PostgreSQL:
		return clause.Expr{
			SQL:  "?",
			Vars: []interface{}{encode},
		}
	case dorm.Mysql:
		return clause.Expr{
			SQL:  "?",
			Vars: []interface{}{encode},
		}
	case dorm.Doris:
		return clause.Expr{SQL: "'" + encode + "'"}
	}
	return clause.Expr{SQL: "null"}

}

func (g *Geometry) Scan(value interface{}) error {
	switch x := value.(type) {
	case nil:
		g.Valid = false
	case *[]byte:
		return g.scanByte(*x)
	case []byte:
		return g.scanByte(x)
	case string:
		//dameng or postgis ST_AsText(geometry) return string
		//postgis and mysql default return wkbhex is string
		// dmgeo.ST_AsGEOJSON(geometry) 返回的数据是错误的，在数据库范文的字符串就是错误的，因为格式不是json的
		return g.scanText(x)
	case *dm.DmStruct:
		return g.scanDmStruct(x)
	}
	return nil
}
func (g *Geometry) scanByte(bs []byte) error {
	var err error
	g.Data, err = ewkb.Unmarshal(bs)
	if err != nil {
		g.Data, err = ewkbhex.Decode(utils.BytesToString(bs))
	}
	g.Valid = err == nil
	return nil
}
func (g *Geometry) scanText(s string) error {
	var err error
	if strings.Contains(s, "(") {
		//WKT
		g.Data, err = wkt.Unmarshal(s)
	} else if strings.Contains(s, "{") {
		//GEOJSON
		err = geojson.Unmarshal(utils.StringToBytes(s), &g.Data)
	} else {
		//WKB
		g.Data, err = ewkbhex.Decode(s)
	}
	if err == nil {
		g.Valid = true
	}
	return err
}
func (g *Geometry) scanDmStruct(ds *dm.DmStruct) error {
	var (
		geoWkb *dm.DmBlob
		length int64
		err    error
	)
	//如果返回的是dmgeo2的ST_Geometry信息，数组只有一个元素
	//如果返回dmgeo的ST_Geometry对象，数组有三个元素
	attributes, err := ds.GetAttributes()
	switch len(attributes) {
	case 0:
		return err
	case 1:
		geoWkb = attributes[0].(*dm.DmBlob)
		//g.Srid = int(attributes[0].(int32))
		//g.GeoTypeId = int(attributes[2].(int32))
	case 2:
		//暂时没有这种情况
		geoWkb = attributes[0].(*dm.DmBlob)
		g.Srid = int(attributes[1].(int32))
	case 3:
		g.Srid = int(attributes[0].(int32))
		geoWkb = attributes[1].(*dm.DmBlob)
		g.GeoTypeId = int(attributes[2].(int32))
	default:
		geoWkb = attributes[0].(*dm.DmBlob)
	}

	length, err = geoWkb.GetLength()
	bs := make([]byte, length)
	_, err = geoWkb.Read(bs)

	g.Data, err = ewkb.Unmarshal(bs)
	if err == nil {
		g.Valid = true
	}
	return err
}

func (g Geometry) MarshalJSON() ([]byte, error) {
	if !g.Valid || g.Data == nil {
		return nullValue, nil
	}
	if g.Data == nil {
		return nullValue, nil
	}

	switch g.FormatType {
	case GeomMarshalWkb:
		rs, err := wkbhex.Encode(g.Data, wkbhex.XDR)
		return []byte(`"` + rs + `"`), err
	case GeomMarshalEWkb:
		rs, err := ewkbhex.Encode(g.Data, wkbhex.XDR)
		return []byte(`"` + rs + `"`), err
	case GeomMarshalWkt:
		rs, err := wkt.Marshal(g.Data)
		return []byte(`"` + rs + `"`), err
	default:
		gJson, err := geojson.Encode(g.Data)
		srid := g.Srid
		if srid == 0 {
			srid = 4326
		}
		crsProp := map[string]interface{}{"name": fmt.Sprintf(`EPSG:%d`, srid)}
		gJson.CRS = &geojson.CRS{Type: "name", Properties: crsProp}
		if err != nil {
			return nil, err
		}
		return jsoniter.Marshal(gJson)
	}
}
func (g *Geometry) UnmarshalJSON(bytes []byte) error {
	//maybe null
	if bytes == nil || len(bytes) < 5 {
		return nil
	}
	var (
		err error
	)

	if (bytes[0] == '"' && bytes[1] == '0' && bytes[2] == '0') || (bytes[0] == '"' && bytes[1] == '0' && bytes[2] == '1') {
		//WKB
		g.Data, err = ewkbhex.Decode(utils.BytesToString(bytes[1 : len(bytes)-1]))
	} else if bytes[0] == '"' && bytes[1] == '(' {
		//WKT
		g.Data, err = wkt.Unmarshal(utils.BytesToString(bytes[1 : len(bytes)-1]))
	} else if bytes[0] == '{' {
		//GEOJSON
		err = geojson.Unmarshal(bytes, &g.Data)
	} else {
		g.Data, err = ewkb.Unmarshal(bytes)
	}

	if err != nil {
		return err
	}
	g.Valid = true
	return nil
}

func (g Geometry) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	tp, err := GetDbTypeName(db, TypeGeometryName)
	if err != nil {
		logger.Error(err)
	}
	return tp
}

// GormDataType 在建表期间实现这个才不会报错，否则需要再gorm tag中显示指定type
func (g Geometry) GormDataType() string {
	return TypeGeometryName
}

func (g *Geometry) GeomFormat(gcsType string, geoFormat string) {
	g.GcsType = gcsType
	g.FormatType = geoFormat
}
