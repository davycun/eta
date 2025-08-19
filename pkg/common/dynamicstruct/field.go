package dynamicstruct

type FieldConfig interface {
	SetName(name string) FieldConfig
	SetPkg(pkg string) FieldConfig
	SetTypeVal(typ any) FieldConfig
	SetTag(tag string) FieldConfig
	SetAnonymous(anonymous bool) FieldConfig

	GetName() string
	GetPkg() string
	GetTypeVal() any
	GetTag() string
	GetAnonymous() bool
}

type fieldConfigImpl struct {
	name      string
	pkg       string
	typVal    interface{}
	tag       string
	anonymous bool
}

func (f *fieldConfigImpl) SetPkg(pkg string) FieldConfig {
	f.pkg = pkg
	return f
}

func (f *fieldConfigImpl) SetTypeVal(typ interface{}) FieldConfig {
	f.typVal = typ
	return f
}

func (f *fieldConfigImpl) SetTag(tag string) FieldConfig {
	f.tag = tag
	return f
}

func (f *fieldConfigImpl) SetAnonymous(anonymous bool) FieldConfig {
	f.anonymous = anonymous
	return f
}

func (f *fieldConfigImpl) SetName(name string) FieldConfig {
	f.name = name
	return f
}

func (f *fieldConfigImpl) GetName() string {
	return f.name
}

func (f *fieldConfigImpl) GetPkg() string {
	return f.pkg
}

func (f *fieldConfigImpl) GetTypeVal() any {
	return f.typVal
}

func (f *fieldConfigImpl) GetTag() string {
	return f.tag
}

func (f *fieldConfigImpl) GetAnonymous() bool {
	return f.anonymous
}
