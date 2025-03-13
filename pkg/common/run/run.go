package run

func Go(fc func()) {
	//defer func() {
	//	if r := recover(); r != nil {
	//		logger.Errorf("goroutine panic %v", r)
	//	}
	//}()
	go fc()
}
