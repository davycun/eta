package dto

type DataImportParam struct {
	FilePath string `json:"file_path,omitempty" binding:"required"` // 导入文件路径
}
