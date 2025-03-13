package dto

type BaseExportParam struct {
	ExportName   string        `json:"export_name,omitempty"`
	ExportFields []ExportField `json:"export_fields,omitempty"`
}

type ExportField struct {
	Field          string        `json:"field"`                   // 字段
	Name           string        `json:"name"`                    // 字段展示名
	Connector      string        `json:"connector,omitempty"`     // 连接符
	ConcatFields   []ExportField `json:"concat_fields,omitempty"` // [{field:b},{field:c.label.name}]
	File           bool          `json:"file"`                    // 是否文件
	FileStorageKey string        `json:"storage_key,omitempty"`   // [初始化时生成]存储里的文件名
	FileFsPath     string        `json:"file_fs_path,omitempty"`  // [初始化时生成]下载后存在文件系统的路径
	FileVal        string        `json:"file_val,omitempty"`      // [初始化时生成]因路径修改，所以值要相应变化
}

// DefaultQueryExportParamExtra 默认的RetrieveParam的Extra属性类型
func DefaultQueryExportParamExtra() any {
	return &BaseExportParam{}
}
