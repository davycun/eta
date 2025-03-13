package data_test

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/http_tes"
	"github.com/davycun/eta/pkg/common/mock/faker"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func templateCreateAllFieldTypeCode() string {
	return fmt.Sprintf("ut_all_field_type_%d%d", time.Now().UnixMilli(), faker.Number(100, 999))
}
func templateCreateSignCode() string {
	return fmt.Sprintf("ut_sign_%d%d", time.Now().UnixMilli(), faker.Number(100, 999))
}
func templateCreateCryptCode() string {
	return fmt.Sprintf("ut_crypt_%d%d", time.Now().UnixMilli(), faker.Number(100, 999))
}
func templateCreateFeatureCode() string {
	return fmt.Sprintf("ut_feature_%d%d", time.Now().UnixMilli(), faker.Number(100, 999))
}

func TestTemplate_Create(t *testing.T) {
	c1 := templateCreateAllFieldTypeCode()
	http_tes.Call(t, templateCreateAllFieldType(c1))
}

func TestTemplate_CreateFeature(t *testing.T) {
	c1 := templateCreateFeatureCode()
	http_tes.Call(t, templateCreateFeature(c1))
}

func templateCreateAllFieldType(code string) http_tes.HttpCase {
	return http_tes.HttpCase{
		Desc:    "表单配置-字段全类型",
		Method:  "POST",
		Path:    "/template/create",
		Headers: map[string]string{"Content-Type": "application/json"},
		Body: fmt.Sprintf(`{
    "single_transaction": true,
    "data": [{
        "code": "%s",
        "remark": "ut_全字段类型",
        "status": "ready",
        "detail": {
            "fields": [
                {
                    "name": "array_int",
                    "title": "整型数组",
                    "type": "array_int",
                    "comment": "整型数组"
                },
                {
                    "name": "array_string",
                    "title": "字符数组",
                    "type": "array_string",
                    "comment": "字符数组"
                },
                {
                    "name": "bool",
                    "title": "布尔",
                    "type": "bool",
                    "comment": "布尔"
                },
                {
                    "name": "numeric",
                    "title": "数字",
                    "type": "numeric(30,20)",
                    "comment": "数字"
                },
                {
                    "name": "geometry",
                    "title": "几何",
                    "type": "geometry",
                    "comment": "几何"
                },
                {
                    "name": "integer",
                    "title": "整数",
                    "type": "integer",
                    "comment": "整数"
                },
                {
                    "name": "bigint",
                    "title": "大整数",
                    "type": "bigint",
                    "comment": "大整数"
                },
                {
                    "name": "json",
                    "title": "json",
                    "type": "json",
                    "comment": "json"
                },
                {
                    "name": "string",
                    "title": "字符",
                    "type": "string",
                    "comment": "字符"
                },
                {
                    "name": "text",
                    "title": "文本",
                    "type": "text",
                    "comment": "文本"
                },
                {
                    "name": "time",
                    "title": "时间",
                    "type": "time",
                    "comment": "时间带时区"
                },
                {
                    "name": "file",
                    "title": "文件",
                    "type": "file",
                    "comment": "文件"
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
func templateCreateSign(code string) http_tes.HttpCase {
	return http_tes.HttpCase{
		Desc:    "表单配置-签名",
		Method:  "POST",
		Path:    "/template/create",
		Headers: map[string]string{"Content-Type": "application/json"},
		Body: fmt.Sprintf(`{
    "single_transaction": true,
    "data": [
        {
            "code": "%s",
            "remark": "ut-签名",
            "status": "ready",
            "detail": {
                "fields": [
                    {
                        "name": "array_int",
                        "title": "整型数组",
                        "type": "array_int",
                        "comment": "整型数组"
                    },
                    {
                        "name": "array_string",
                        "title": "字符数组",
                        "type": "array_string",
                        "comment": "字符数组"
                    },
                    {
                        "name": "bool",
                        "title": "布尔",
                        "type": "bool",
                        "comment": "布尔"
                    },
                    {
                        "name": "numeric",
                        "title": "数字",
                        "type": "numeric(30,20)",
                        "comment": "数字"
                    },
                    {
                        "name": "geometry",
                        "title": "几何",
                        "type": "geometry",
                        "comment": "几何"
                    },
                    {
                        "name": "integer",
                        "title": "整数",
                        "type": "integer",
                        "comment": "整数"
                    },
                    {
                        "name": "bigint",
                        "title": "大整数",
                        "type": "bigint",
                        "comment": "大整数"
                    },
                    {
                        "name": "json",
                        "title": "json",
                        "type": "json",
                        "comment": "json"
                    },
                    {
                        "name": "string",
                        "title": "字符",
                        "type": "string",
                        "comment": "字符"
                    },
                    {
                        "name": "text",
                        "title": "文本",
                        "type": "text",
                        "comment": "文本"
                    },
                    {
                        "name": "time",
                        "title": "时间",
                        "type": "time",
                        "comment": "时间带时区"
                    },
                    {
                        "name": "file",
                        "title": "文件",
                        "type": "file",
                        "comment": "文件"
                    },
                    {
                        "name": "sign1",
                        "title": "签名1",
                        "type": "string",
                        "comment": "签名1"
                    },
                    {
                        "name": "sign2",
                        "title": "签名2",
                        "type": "string",
                        "comment": "签名2"
                    }
                ],
                "index": []
            },
            "ft_sign": [
                {
                    "algo": "hmac_sm3",
                    "verify_field": "sign_matched1",
                    "salt": "sign1_salt",
                    "field": "sign1",
                    "fields": [
                        "array_int",
                        "bool",
                        "integer",
                        "string",
                        "file"
                    ]
                },{
                    "algo": "hmac_sm3",
                    "verify_field": "sign_matched2",
                    "salt": "sign2_salt",
                    "field": "sign2",
                    "fields": [
                        "array_string",
                        "time",
                        "bigint",
                        "text"
                    ]
                }
            ]
        }
    ]
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
func templateCreateCrypt(code string) http_tes.HttpCase {
	return http_tes.HttpCase{
		Desc:    "表单配置-加密",
		Method:  "POST",
		Path:    "/template/create",
		Headers: map[string]string{"Content-Type": "application/json"},
		Body: fmt.Sprintf(`{
    "single_transaction": true,
    "data": [
        {
            "code": "%s",
            "remark": "ut-加密",
            "status": "ready",
            "detail": {
                "fields": [
                    {
                        "name": "array_int",
                        "title": "整型数组",
                        "type": "array_int",
                        "comment": "整型数组"
                    },
                    {
                        "name": "array_string",
                        "title": "字符数组",
                        "type": "array_string",
                        "comment": "字符数组"
                    },
                    {
                        "name": "bool",
                        "title": "布尔",
                        "type": "bool",
                        "comment": "布尔"
                    },
                    {
                        "name": "numeric",
                        "title": "数字",
                        "type": "numeric(30,10)",
                        "comment": "数字"
                    },
                    {
                        "name": "geometry",
                        "title": "几何",
                        "type": "geometry",
                        "comment": "几何"
                    },
                    {
                        "name": "integer",
                        "title": "整数",
                        "type": "integer",
                        "comment": "整数"
                    },
                    {
                        "name": "bigint",
                        "title": "大整数",
                        "type": "bigint",
                        "comment": "大整数"
                    },
                    {
                        "name": "json",
                        "title": "json",
                        "type": "json",
                        "comment": "json"
                    },
                    {
                        "name": "string",
                        "title": "字符",
                        "type": "string",
                        "comment": "字符"
                    },
                    {
                        "name": "text",
                        "title": "文本",
                        "type": "text",
                        "comment": "文本"
                    },
                    {
                        "name": "time",
                        "title": "时间",
                        "type": "time",
                        "comment": "时间带时区"
                    },
                    {
                        "name": "file",
                        "title": "文件",
                        "type": "file",
                        "comment": "文件"
                    },
                    {
                        "name": "enc1",
                        "title": "加密字段1",
                        "type": "string",
                        "comment": "加密字段1"
                    },
                    {
                        "name": "enc2",
                        "title": "加密字段2",
                        "type": "string",
                        "comment": "加密字段2"
                    },
                    {
                        "name": "enc3",
                        "title": "加密字段3",
                        "type": "string",
                        "comment": "加密字段3"
                    }
                ],
                "index": []
            },
            "ft_encrypt": [
                {
                    "algo": "sm4_cbc_pkcs7padding",
                    "secret_key": ["8eadb267ecd6e860"],
                    "field": "enc1"
                },{
                    "algo": "sm4_cbc_pkcs7padding",
                    "secret_key": ["8eadb267ecd6e864"],
                    "field": "enc2",
                    "keep_txt_pre_cnt": 4
                },{
                    "algo": "sm4_cbc_pkcs7padding",
                    "secret_key": ["9eadb267ecd6e864"],
                    "field": "enc3",
                    "keep_txt_pre_cnt": 4,
                    "keep_txt_suf_cnt": 4
                }
            ]
        }
    ]
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
func templateCreateFeature(code string) http_tes.HttpCase {
	return http_tes.HttpCase{
		Desc:    "表单配置-feature",
		Method:  "POST",
		Path:    "/template/create",
		Headers: map[string]string{"Content-Type": "application/json"},
		Body: fmt.Sprintf(`{
    "single_transaction": true,
    "data": [
        {
            "code": "%s",
            "remark": "ut-feature",
            "status": "ready",
            "detail": {
                "fields": [
                    {
                        "name": "string",
                        "title": "字符",
                        "type": "string",
                        "comment": "字符"
                    }
                ],
                "index": []
            },
			"feature":{
				"history": true,
				"field_updater": true
			}
        }
    ]
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
