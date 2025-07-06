package subscribe_test

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/http_tes"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/module/setting"
	"github.com/davycun/eta/pkg/module/subscribe"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type People struct {
	Name string
	IdNo string
}

func TestRegistrySub(t *testing.T) {

	var (
		//pepId = "123123123"
		//subId = "321321321"
		ch = make(chan struct{}, 1)
		st = setting.Setting{Namespace: "my", Name: "my1", Category: "my2", Content: ctype.NewJson(People{Name: "name:" + global.GenerateIDStr(), IdNo: "idNo:" + global.GenerateIDStr()})}
	)

	server := httptest.NewServer(http.DefaultServeMux)
	//注册回调，
	http_tes.Call(t, http_tes.HttpCase{
		Method: http.MethodPost,
		Path:   "/subscribe/create",
		Body: dto.ModifyParam{
			Data: []subscribe.Subscriber{
				{
					//BaseEntity: entity.BaseEntity{ID: subId},
					Method: http.MethodPost,
					Url:    fmt.Sprintf("%s/people/receive", server.URL),
					Target: "t_people",
				},
			},
		},
		ShowBody: true,
	})

	http.HandleFunc("/setting/receive", func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			ch <- struct{}{}
		}()
		bs, err1 := io.ReadAll(request.Body)
		assert.Nil(t, err1)

		args := dto.ModifyParam{}
		dt := make([]People, 0, 5)
		args.Data = &dt
		err := jsoniter.Unmarshal(bs, &args)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(dt))
		assert.Equal(t, st.Name, dt[0].Name)
		_, err = writer.Write([]byte("ok"))
		assert.Nil(t, err)
	})

	http_tes.Call(t, http_tes.HttpCase{
		Method: http.MethodPost,
		Path:   "/setting/create",
		Body: dto.ModifyParam{
			Data: []setting.Setting{st},
		},
		ShowBody: true,
	})

	ticker := time.NewTicker(time.Second * 3)
	select {
	case <-ch:
		logger.Info("正确返回了")
	case <-ticker.C:
		logger.Info("超时了")
	}
}
