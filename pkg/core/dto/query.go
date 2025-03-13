package dto

//type QueryParam struct {
//	RetrieveParam
//}
//
//func (q QueryParam) Clone() QueryParam {
//	qp := QueryParam{}
//	qp.RetrieveParam = q.RetrieveParam.Clone()
//	return qp
//}
//
//type QueryResult struct {
//	Total    int64 `json:"total"`
//	PageSize int   `json:"page_size,omitempty"`
//	PageNum  int   `json:"page_num,omitempty"`
//	Data     any   `json:"data"`
//}
//
//// DefaultQueryParamExtra 默认的RetrieveParam的Extra属性类型
//func DefaultQueryParamExtra() any {
//	return &ctype.Map{}
//}

//func InitPage(args *QueryParam) {
//	if args.PageSize < 1 {
//		args.PageSize = 10
//	}
//	if args.PageNum < 1 {
//		args.PageNum = 1
//	}
//}
