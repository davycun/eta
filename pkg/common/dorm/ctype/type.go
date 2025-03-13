package ctype

type GeometryFormat interface {
	// GeomFormat
	//geoType 表示需要做个类型，百度、高德、84
	//geoFormat表示需要的是个，wkb、geojson、wkt
	GeomFormat(gcsType string, geoFormat string)
}
