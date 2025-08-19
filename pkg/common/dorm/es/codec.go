package es

type Option func(esApi *Api)

type CodecConfig struct {
	DbSchema  string
	TableName string
}

func CodecOpt(dbSchema string, tableName string) Option {
	return func(esApi *Api) {
		esApi.codecConfig = CodecConfig{
			DbSchema:  dbSchema,
			TableName: tableName,
		}
	}
}

//func handleCodec(s *Api, t plugin.ProcessType, dest reflect.Value) {
//if t == plugin.ProcessBefore {
//	// es 不支持 before 操作
//	return
//}
//plugin.HandleCodec(s.codecConfig.DbSchema, s.codecConfig.TableName, t, dest)
//}
