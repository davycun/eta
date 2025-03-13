package data

//func modifyCallback(cfg *service.ModifyConfig, pos service.CallbackPosition) error {
//
//	switch pos {
//	case service.CallbackBefore:
//	case service.CallbackAfter:
//		switch cfg.Curd {
//		case service.CurdCreate:
//			DelMenuCache()
//		case service.CurdUpdate, service.CurdUpdateByFilters:
//			DelMenuCache()
//		case service.CurdDelete, service.CurdDeleteByFilters:
//			DelMenuCache()
//		}
//
//	}
//	return nil
//}
//
//func retrieveCallbacks(cfg *hook.SrvConfig, pos hook.CallbackPosition) error {
//
//
//	switch pos {
//	case service.CallbackBefore:
//	case service.CallbackAfter:
//		return hook.AfterRetrieve(cfg, service.ProcessGeometryType, service.CurdQuery, service.CurdPartition, service.CurdDetail)
//	}
//	return nil
//}
