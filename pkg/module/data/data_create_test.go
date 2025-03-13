package data_test

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/http_tes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestData_Create(t *testing.T) {
	c1, c2, c3, c4 := templateCreateAllFieldTypeCode(), templateCreateSignCode(), templateCreateCryptCode(), templateCreateFeatureCode()
	http_tes.Call(t, templateCreateAllFieldType(c1), templateCreateSign(c2), templateCreateCrypt(c3), templateCreateFeature(c4))
	http_tes.Call(t, dataCreateAllFieldType(c1), dataCreateSign(c2), dataCreateCrypt(c3), dataCreateFeature(c4))
}

func dataCreateAllFieldType(code string) http_tes.HttpCase {
	return http_tes.HttpCase{
		Desc:    "表单数据-字段全类型",
		Method:  "POST",
		Path:    fmt.Sprintf("/data/%s/create", code),
		Headers: map[string]string{"Content-Type": "application/json"},
		Body: `{
    "single_transaction": true,
    "data": [
        {
            "array_int": [1,2,3],
            "array_string": ["在","🤔","培"],
            "bool": true,
            "numeric": 322234.34434344,
            "integer": 9283,
            "bigint": 23783278327,
            "json": {"a":"s","轰轰烈烈":"敢梦的人"},
            "string": "凭什么倔强daskj232323",
            "text": "燃烧信仰",
            "time": "2023-12-15T17:40:39+08:00",
            "file": ["dir/file1.jpg","dir/xx/file.jpg"]
        }
    ]
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
func dataCreateSign(code string) http_tes.HttpCase {
	return http_tes.HttpCase{
		Desc:    "表单数据-签名",
		Method:  "POST",
		Path:    fmt.Sprintf("/data/%s/create", code),
		Headers: map[string]string{"Content-Type": "application/json"},
		Body: `{
    "single_transaction": true,
    "data": [
        {
            "array_int": [1,2,3],
            "array_string": ["在","🤔","培"],
            "bool": true,
            "numeric": 3222434.34434344,
            "integer": 9283,
            "bigint": 23783278327,
            "json": {"a":"s","轰轰烈烈":"敢梦的人"},
            "string": "凭什么倔强daskj232323",
            "text": "燃烧信仰",
            "time": "2023-12-15T17:40:39+08:00",
            "file": ["dir/file1.jpg","dir/xx/file.jpg"]
        }
    ]
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
func dataCreateCrypt(code string) http_tes.HttpCase {
	return http_tes.HttpCase{
		Desc:    "表单数据-加密",
		Method:  "POST",
		Path:    fmt.Sprintf("/data/%s/create", code),
		Headers: map[string]string{"Content-Type": "application/json"},
		Body: `{
    "single_transaction": true,
    "data": [
        {
            "array_int": [1,2,3],
            "array_string": ["在","🤔","培"],
            "bool": true,
            "numeric": 834.34434344,
            "integer": 9283,
            "bigint": 23783278327,
            "json": {"a":"s","轰轰烈烈":"敢梦的人"},
            "string": "凭什么倔强daskj232323",
            "text": "燃烧信仰",
            "time": "2023-12-15T17:40:39+08:00",
            "file": ["dir/file1.jpg","dir/xx/file.jpg"],
            "enc1":"在DM系统中，代理服务是运行在服务器端，调度并执行作业、监视警报的服务。通过它用户可以自动执行部分管理任务，如定期备份、出错通知等，减轻工作负担。必须启动代理服务后，作业与调度才能正常工作。代理服务加载系统定义的所有作业，并根据其调度信息安排其执行时间。当特定的时刻到来时，启动作业，并依次执行作业包含的每个步骤。代理服务不仅监控时间事件，同时也监控服务器内部的警报事件，当服务器在运行中产生某个特定事件时（如执行操作失败），代理服务会检测到这个事件的发生，并触发相应的警报。",
            "enc2":"在DM系统中，代理服务是运行在服务器端，调度并执行作业、监视警报的服务。通过它用户可以自动执行部分管理任务，如定期备份、出错通知等，减轻工作负担。必须启动代理服务后，作业与调度才能正常工作。代理服务加载系统定义的所有作业，并根据其调度信息安排其执行时间。当特定的时刻到来时，启动作业，并依次执行作业包含的每个步骤。代理服务不仅监控时间事件，同时也监控服务器内部的警报事件，当服务器在运行中产生某个特定事件时（如执行操作失败），代理服务会检测到这个事件的发生，并触发相应的警报。",
            "enc3":"在DM系统中，代理服务是运行在服务器端，调度并执行作业、监视警报的服务。通过它用户可以自动执行部分管理任务，如定期备份、出错通知等，减轻工作负担。必须启动代理服务后，作业与调度才能正常工作。代理服务加载系统定义的所有作业，并根据其调度信息安排其执行时间。当特定的时刻到来时，启动作业，并依次执行作业包含的每个步骤。代理服务不仅监控时间事件，同时也监控服务器内部的警报事件，当服务器在运行中产生某个特定事件时（如执行操作失败），代理服务会检测到这个事件的发生，并触发相应的警报。"
        }
    ]
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
func dataCreateFeature(code string) http_tes.HttpCase {
	return http_tes.HttpCase{
		Desc:    "表单数据-history",
		Method:  "POST",
		Path:    fmt.Sprintf("/data/%s/create", code),
		Headers: map[string]string{"Content-Type": "application/json"},
		Body: `{
				"single_transaction": true,
				"data": [
					{
						"string": "凭什么倔强daskj232323"
					}
				]
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
