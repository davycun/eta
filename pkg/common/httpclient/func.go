package httpclient

func Get(url string, dst any) error {
	err := DefaultHttpClient.Get(url, dst)

	return err
}
func Post(url, contentType string, reqBody, respBody any) error {
	err := DefaultHttpClient.Post(url, contentType, reqBody, respBody)
	return err
}
func Put(url, contentType string, reqBody, respBody any) error {
	err := DefaultHttpClient.Put(url, contentType, reqBody, respBody)
	return err
}
