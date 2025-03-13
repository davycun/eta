package httpclient

import (
	"bytes"
	"errors"
	"github.com/davycun/eta/pkg/common/logger"
	jsoniter "github.com/json-iterator/go"
	"io"
	"net/http"
)

var (
	JSON = JsonBinding{}
)

type Binding interface {
	Bind(dt []byte, dst any) error
	UnBind(dst any) (io.ReadCloser, error)
}

type JsonBinding struct {
}

func (j JsonBinding) Bind(dt []byte, dst any) error {
	if len(dt) > 0 {
		return jsoniter.Unmarshal(dt, dst)
	}
	return errors.New("binding target dst is nil")
}
func (j JsonBinding) UnBind(dst any) (io.ReadCloser, error) {

	dt, err := jsoniter.Marshal(dst)
	if err != nil {
		return nil, err
	}
	return io.NopCloser(bytes.NewReader(dt)), nil
}

func GetBinding(contentType string) Binding {

	switch contentType {
	case MIMEJSON:
		return JSON
	default:
		return JSON
	}
}

// 读了多少返回多少
func readValue(res *http.Response) ([]byte, error) {
	rs := make([]byte, res.ContentLength)
	for {
		tmp := make([]byte, res.ContentLength)
		n, err := res.Body.Read(rs)
		if n > 0 {
			rs = append(rs, tmp[:n]...)
		}
		if err == io.EOF || n < 1 {
			break
		}
		if err != nil {
			logger.Errorf("readValur from http.Response error %s ", err.Error())
			return rs, err
		}
	}
	return rs, nil
}
