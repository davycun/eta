package data_test

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/http_tes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestData_Query(t *testing.T) {

	c1, c2, c3 := templateCreateAllFieldTypeCode(), templateCreateSignCode(), templateCreateCryptCode()

	http_tes.Call(t, templateCreateAllFieldType(c1), templateCreateSign(c2), templateCreateCrypt(c3))
	http_tes.Call(t, dataCreateAllFieldType(c1), dataCreateSign(c2), dataCreateCrypt(c3))
	http_tes.Call(t, dataQueryAllFieldType(c1), dataQuerySign(c2), dataQueryCrypt(c3))
}

func dataQueryAllFieldType(code string) http_tes.HttpCase {
	return http_tes.HttpCase{
		Desc:    "表单数据-字段全类型",
		Method:  "POST",
		Path:    fmt.Sprintf("/data/%s/query", code),
		Headers: map[string]string{"Content-Type": "application/json"},
		Body: `{
    "filters": [
        {
            "logical_operator": "and",
            "column": "id",
            "operator": "!=",
            "value": ""
        }
    ],
    "auto_count": true
}`,
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

func dataQuerySign(code string) http_tes.HttpCase {
	return http_tes.HttpCase{
		Desc:    "表单数据-签名",
		Method:  "POST",
		Path:    fmt.Sprintf("/data/%s/query", code),
		Headers: map[string]string{"Content-Type": "application/json"},
		Body: `{
    "filters": [
        {
            "logical_operator": "and",
            "column": "id",
            "operator": "!=",
            "value": ""
        }
    ],
    "auto_count": true
}`,
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

func dataQueryCrypt(code string) http_tes.HttpCase {
	return http_tes.HttpCase{
		Desc:    "表单数据-加密",
		Method:  "POST",
		Path:    fmt.Sprintf("/data/%s/query", code),
		Headers: map[string]string{"Content-Type": "application/json"},
		Body: `{
    "filters": [
        {
            "logical_operator": "and",
            "column": "id",
            "operator": "!=",
            "value": ""
        }
    ],
    "auto_count": true
}`,
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
