package mock

type CitizenAppParam struct {
	Copy      Copy              `json:"copy,omitempty"`
	TblConfig map[string]Config `json:"table_config,omitempty"`
}

type CitizenParam struct {
	TblConfig map[string]Config `json:"table_config,omitempty"`
}

type Copy struct {
	DbSchema string   `json:"db_schema,omitempty"` // 数据库 schema
	Tables   []string `json:"tables,omitempty"`    // 要复制哪些数据库表
}

type Config struct {
	DeleteExistsData bool   `json:"delete_exists_data,omitempty"` // 删除已有数据
	Size             int    `json:"size,omitempty"`               // 生成多少条
	CreateBatchSize  int    `json:"create_batch_size,omitempty"`  // batch size of CreateInBatches
	ParentId         string `json:"parent_id,omitempty"`
}
