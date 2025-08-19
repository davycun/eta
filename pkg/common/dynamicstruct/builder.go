package dynamicstruct

import (
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/duke-git/lancet/v2/slice"
	"reflect"
)

type Builder interface {
	AddField(name string, typ interface{}, tag string) Builder
	RemoveField(name string) Builder
	HasField(name string) bool
	Merge(builder Builder) Builder
	Build() DynamicStruct
	RangeFields(fn func(idx int, fd FieldConfig))
}
type builderImpl struct {
	fieldsName []string                    //为了保障顺序，所以存了一份key
	fields     map[string]*fieldConfigImpl //
}

func (b *builderImpl) init() *builderImpl {
	if b.fieldsName == nil {
		b.fieldsName = make([]string, 0)
	}
	if b.fields == nil {
		b.fields = make(map[string]*fieldConfigImpl)
	}
	return b
}

func (b *builderImpl) AddField(name string, typ interface{}, tag string) Builder {
	return b.addField(name, "", typ, tag, false)
}

func (b *builderImpl) addField(name string, pkg string, typ interface{}, tag string, anonymous bool) Builder {

	if _, ok := b.init().fields[name]; ok {
		logger.Warnf("dynamicstruct add field[%s] has exists, will be overrided", name)
	} else {
		b.fieldsName = append(b.fieldsName, name)
	}
	b.fields[name] = &fieldConfigImpl{
		name:      name,
		pkg:       pkg,
		typVal:    typ,
		tag:       tag,
		anonymous: anonymous,
	}
	return b
}

func (b *builderImpl) RemoveField(name string) Builder {
	if _, ok := b.init().fields[name]; !ok {
		return b
	}
	delete(b.fields, name)
	b.fieldsName = slice.Filter(b.fieldsName, func(index int, item string) bool {
		return item != name
	})
	return b
}

func (b *builderImpl) HasField(name string) bool {
	_, ok := b.init().fields[name]
	return ok
}
func (b *builderImpl) Merge(bd Builder) Builder {
	if bd == nil {
		return b
	}
	bd.RangeFields(func(idx int, fd FieldConfig) {
		b.addField(fd.GetName(), fd.GetPkg(), fd.GetTypeVal(), fd.GetTag(), fd.GetAnonymous())
	})
	return b
}

func (b *builderImpl) RangeFields(fn func(idx int, fd FieldConfig)) {
	for idx, fd := range b.fieldsName {
		fn(idx, b.fields[fd])
	}
}

func (b *builderImpl) Build() DynamicStruct {
	var structFields []reflect.StructField

	for _, field := range b.fields {
		structFields = append(structFields, reflect.StructField{
			Name:      field.name,
			PkgPath:   field.pkg,
			Type:      reflect.TypeOf(field.typVal),
			Tag:       reflect.StructTag(field.tag),
			Anonymous: field.anonymous,
		})
	}
	return &dynamicStructImpl{
		definition: reflect.StructOf(structFields),
	}
}

type DynamicStruct interface {
	New() interface{}
	NewSliceOfStructs() interface{}
	NewMapOfStructs(key interface{}) interface{}
}

type dynamicStructImpl struct {
	definition reflect.Type
}

func (ds *dynamicStructImpl) New() interface{} {
	return reflect.New(ds.definition).Interface()
}

func (ds *dynamicStructImpl) NewSliceOfStructs() interface{} {
	return reflect.New(reflect.SliceOf(ds.definition)).Interface()
}

func (ds *dynamicStructImpl) NewMapOfStructs(key interface{}) interface{} {
	return reflect.New(reflect.MapOf(reflect.Indirect(reflect.ValueOf(key)).Type(), ds.definition)).Interface()
}

func NewStruct() Builder {
	return &builderImpl{
		fields:     make(map[string]*fieldConfigImpl),
		fieldsName: make([]string, 0),
	}
}

func ExtendStruct(value interface{}) Builder {
	return MergeStructs(value)
}

func MergeStructs(values ...interface{}) Builder {
	builder := NewStruct()

	for _, value := range values {
		valueOf := reflect.Indirect(reflect.ValueOf(value))
		typeOf := valueOf.Type()

		for i := 0; i < valueOf.NumField(); i++ {
			fval := valueOf.Field(i)
			ftyp := typeOf.Field(i)
			builder.(*builderImpl).addField(ftyp.Name, ftyp.PkgPath, fval.Interface(), string(ftyp.Tag), ftyp.Anonymous)
		}
	}

	return builder
}
