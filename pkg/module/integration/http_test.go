package integration_test

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/http_tes"
	"github.com/davycun/eta/pkg/common/mock/faker"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	templateCode     = fmt.Sprintf("ut_multi_table_operate_%d%d", time.Now().UnixMilli(), faker.Number(100, 999))
	deltaSettingName = fmt.Sprintf("name_%d%d", time.Now().UnixMilli(), faker.Number(100, 999))
)

func TestTransaction(t *testing.T) {
	http_tes.Call(t, templateCreate(templateCode))
	http_tes.Call(t, dataTransaction())
}

func dataTransaction() http_tes.HttpCase {
	return http_tes.HttpCase{
		Desc:    "多表写操作",
		Method:  "POST",
		Path:    fmt.Sprintf("/integration/Transaction"),
		Headers: map[string]string{"Content-Type": "application/json"},
		Body: fmt.Sprintf(`{
    "items": [
        {
            "command": "create",
            "entity_code": "eta_setting",
            "param": {
                "data": [
                    {
						"namespace": "ut",
						"category": "xxx",
						"name": "%s",
						"content": "content"
                    }
                ]
            }
        },
        {
            "command": "create",
            "entity_code": "%s",
            "param": {
                "data": [
                    {
                        "text": "燃烧信仰"
                    }
                ]
            }
        }
    ]
}`, deltaSettingName, templateCode),
		ShowBody: true,
		Code:     "200",
		ValidateFunc: []func(t *testing.T, resp *http_tes.Resp){
			func(t *testing.T, resp *http_tes.Resp) {
				res := resp.Result.(map[string]interface{})
				assert.NotNil(t, res["items"])
				assert.NotEmpty(t, res["items"])
			},
		},
	}
}

func templateCreate(code string) http_tes.HttpCase {
	return http_tes.HttpCase{
		Desc:    "表单配置-多表操作",
		Method:  "POST",
		Path:    "/template/create",
		Headers: map[string]string{"Content-Type": "application/json"},
		Body: fmt.Sprintf(`{
    "single_transaction": true,
    "data": [{
        "code": "%s",
        "remark": "ut_多表操作",
        "status": "ready",
        "detail": {
            "fields": [
                {
                    "name": "text",
                    "title": "文本",
                    "type": "text",
                    "comment": "文本"
                }
            ],
            "index": []
        }
    }]
}`, code),
		ShowBody: true,
		Code:     "200",
		ValidateFunc: []func(t *testing.T, resp *http_tes.Resp){
			func(t *testing.T, resp *http_tes.Resp) {
				res := resp.Result.(map[string]interface{})
				assert.NotNil(t, res["data"])
			},
		},
	}
}
