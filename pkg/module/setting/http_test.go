package setting_test

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/http_tes"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCRUD(t *testing.T) {
	var (
		id              string
		updatedAt       float64
		namespace       = "unit_test_ns"
		name            = "unit_test_name"
		category        = "unit_test_category"
		namespaceUpdate = "namespace_xxx"
		nameUpdate      = "name_xxx"
		categoryUpdate  = "category_xxx"
	)
	http_tes.Call(t, createCase(namespace, name, category, func(result map[string]interface{}) {
		id = result["data"].([]interface{})[0].(map[string]any)["id"].(string)
		updatedAt = result["data"].([]interface{})[0].(map[string]any)["updated_at"].(float64)
	}))
	http_tes.Call(t, updateByIdCase(id, updatedAt, namespaceUpdate, nameUpdate, categoryUpdate, func(result map[string]interface{}) {
		updatedAt = result["data"].(map[string]interface{})["success_data"].([]interface{})[0].(map[string]any)["updated_at"].(float64)
	}))
	http_tes.Call(t, query(id, func(result map[string]interface{}) {
		d := result["data"].([]interface{})[0].(map[string]any)
		// namespace、category、name未被修改
		assert.Equal(t, namespace, d["namespace"])
		assert.Equal(t, name, d["name"])
		assert.Equal(t, category, d["category"])
	}))
	http_tes.Call(t, del(id, updatedAt, func(result map[string]interface{}) {
		logger.Info("success")
	}))
}

func createCase(namespace, name, category string, f func(resp map[string]interface{})) http_tes.HttpCase {
	return http_tes.HttpCase{
		Desc:    "配置管理-新增",
		Method:  "POST",
		Path:    "/setting/create",
		Headers: map[string]string{"Content-Type": "application/json"},
		Body: fmt.Sprintf(`{
				  "data": [
					{
					  "namespace": "%s",
					  "category": "%s",
					  "name": "%s",
					  "content": {
						"field1": "aaaaa",
						"field2": [1,2,3,4,5],
						"field3": ["a","b","c"]
					  }
					}
				  ]
				}`, namespace, category, name),
		ShowBody: true,
		Code:     "200",
		ValidateFunc: []http_tes.ValidateFunc{
			func(t *testing.T, resp *http_tes.Response) {
				res := resp.Resp.Result.(map[string]interface{})
				assert.NotNil(t, res["data"])
				f(res)
			},
		},
	}
}

func updateByIdCase(id string, updatedAt float64, namespace, name, category string, f func(resp map[string]interface{})) http_tes.HttpCase {
	return http_tes.HttpCase{
		Desc:    "配置管理-根据ID更新",
		Method:  "POST",
		Path:    "/setting/update",
		Headers: map[string]string{"Content-Type": "application/json"},
		Body: fmt.Sprintf(`{
				  "data": [
					{
					  "id": "%s",
                      "updated_at": %d,
					  "category": "%s",
					  "namespace": "%s",
					  "name": "%s",
					  "content": {
						"field1": "aaaaa",
						"field2": [1,2,3,4,5],
						"field3": ["a","b","c"]
					  }
					}
				  ]
				}`, id, int(updatedAt), category, namespace, name),
		ShowBody: true,
		Code:     "200",
		ValidateFunc: []http_tes.ValidateFunc{
			func(t *testing.T, resp *http_tes.Response) {
				res := resp.Resp.Result.(map[string]interface{})
				assert.NotNil(t, res["data"])
				f(res)
			},
		},
	}
}

func query(id string, f func(resp map[string]interface{})) http_tes.HttpCase {
	return http_tes.HttpCase{
		Desc:    "配置管理-查询校验",
		Method:  "POST",
		Path:    "/setting/query",
		Headers: map[string]string{"Content-Type": "application/json"},
		Body: fmt.Sprintf(`{
								   "filters": [
										{
											"logical_operator": "and",
											"column": "id",
											"operator": "=",
											"value": "%s"
										}
									]
								}`, id),
		ShowBody: true,
		Code:     "200",
		ValidateFunc: []http_tes.ValidateFunc{
			func(t *testing.T, resp *http_tes.Response) {
				res := resp.Resp.Result.(map[string]interface{})
				assert.NotNil(t, res["data"])
				f(res)
			},
		},
	}
}

func del(id string, updatedAt float64, f func(resp map[string]interface{})) http_tes.HttpCase {
	return http_tes.HttpCase{
		Desc:    "配置管理-根据ID删除",
		Method:  "POST",
		Path:    "/setting/delete",
		Headers: map[string]string{"Content-Type": "application/json"},
		Body: fmt.Sprintf(`{
				  "data": [
					{
					  "id": "%s",
                      "updated_at": %d
					}
				  ]
				}`, id, int(updatedAt)),
		ShowBody: true,
		Code:     "200",
		ValidateFunc: []http_tes.ValidateFunc{
			func(t *testing.T, resp *http_tes.Response) {
				res := resp.Resp.Result.(map[string]interface{})
				assert.NotNil(t, res["rows_affected"])
				f(res)
			},
		},
	}
}
