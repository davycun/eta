package app

type MigrateAppParam struct {
	SendWsMessage bool               `json:"send_ws_message"`
	AppIds        []string           `json:"app_ids"`
	Dbs           []string           `json:"dbs"`
	Es            map[string]EsParam `json:"es"` // tableName:ES参数
}

type EsParam struct {
	Settings map[string]interface{} `json:"settings"`
}
