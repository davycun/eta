package data_test

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/crypt"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/http_tes"
	"github.com/davycun/eta/pkg/common/mock/faker"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/module/template"
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

	tmp := template.Template{
		BaseEntity: entity.BaseEntity{
			Remark: "ut_全字段类型",
		},
		Code:   code,
		Status: template.Ready,
		Table: entity.Table{
			Feature: entity.Feature{},
			Fields: []entity.TableField{
				{
					Name:    "array_int",
					Title:   "整形数组",
					Type:    ctype.TypeArrayIntName,
					Comment: "整形数组",
				},
				{
					Name:    "array_string",
					Title:   "字符数组",
					Type:    ctype.TypeArrayStringName,
					Comment: "字符数组",
				},
				{
					Name:    "bool",
					Title:   "布尔",
					Type:    ctype.TypeBoolName,
					Comment: "布尔",
				},
				{
					Name:    "numeric",
					Title:   "数字",
					Type:    ctype.TypeNumericName,
					Comment: "数字",
				},
				{
					Name:    "geometry",
					Title:   "几何",
					Type:    ctype.TypeGeometryName,
					Comment: "几何",
				},
				{
					Name:    "integer",
					Title:   "整数",
					Type:    ctype.TypeIntegerName,
					Comment: "整数",
				},
				{
					Name:    "bigint",
					Title:   "大整数",
					Type:    ctype.TypeBigIntegerName,
					Comment: "大整数",
				},
				{
					Name:    "json",
					Title:   "json",
					Type:    ctype.TypeJsonName,
					Comment: "json",
				},
				{
					Name:    "string",
					Title:   "字符",
					Type:    ctype.TypeStringName,
					Comment: "字符",
				},
				{
					Name:    "text",
					Title:   "文本",
					Type:    ctype.TypeTextName,
					Comment: "文本",
				},
				{
					Name:    "time",
					Title:   "时间",
					Type:    ctype.TypeTimeName,
					Comment: "时间带时区",
				},
				{
					Name:    "file",
					Title:   "文件",
					Type:    ctype.TypeFileName,
					Comment: "文件",
				},
			},
		},
	}

	return http_tes.HttpCase{
		Desc:    "表单配置-字段全类型",
		Method:  "POST",
		Path:    "/template/create",
		Headers: map[string]string{"Content-Type": "application/json"},
		Body: &dto.Param{
			ModifyParam: dto.ModifyParam{
				SingleTransaction: true,
				Data:              &[]template.Template{tmp},
			},
		},
		ShowBody: true,
		Code:     "200",
		ValidateFunc: []http_tes.ValidateFunc{
			func(t *testing.T, resp *http_tes.Response) {
				res := resp.Resp.Result.(map[string]interface{})
				assert.NotNil(t, res["data"])
			},
		},
	}
}
func templateCreateSign(code string) http_tes.HttpCase {

	tmp := template.Template{
		BaseEntity: entity.BaseEntity{
			Remark: "ut_全字段类型",
		},
		Code:   code,
		Status: template.Ready,
		Table: entity.Table{
			Feature: entity.Feature{
				SignFields: []entity.SignFieldsInfo{
					{
						Enable:      true,
						Algo:        crypt.AlgoSignHmacSm3,
						Fields:      []string{"array_int", "bool", "integer", "string"},
						Field:       "sign1",
						VerifyField: "sign1_matched",
						Key:         "sign1_key",
					},
					{
						Enable:      true,
						Algo:        crypt.AlgoSignHmacSm3,
						Fields:      []string{"array_string", "time", "bigint", "text"},
						Field:       "sign2",
						VerifyField: "sign2_matched",
						Key:         "sign2_key",
					},
				},
			},
			Fields: []entity.TableField{
				{
					Name:    "array_int",
					Title:   "整形数组",
					Type:    ctype.TypeArrayIntName,
					Comment: "整形数组",
				},
				{
					Name:    "array_string",
					Title:   "字符数组",
					Type:    ctype.TypeArrayStringName,
					Comment: "字符数组",
				},
				{
					Name:    "bool",
					Title:   "布尔",
					Type:    ctype.TypeBoolName,
					Comment: "布尔",
				},
				{
					Name:    "numeric",
					Title:   "数字",
					Type:    ctype.TypeNumericName,
					Comment: "数字",
				},
				{
					Name:    "geometry",
					Title:   "几何",
					Type:    ctype.TypeGeometryName,
					Comment: "几何",
				},
				{
					Name:    "integer",
					Title:   "整数",
					Type:    ctype.TypeIntegerName,
					Comment: "整数",
				},
				{
					Name:    "bigint",
					Title:   "大整数",
					Type:    ctype.TypeBigIntegerName,
					Comment: "大整数",
				},
				{
					Name:    "json",
					Title:   "json",
					Type:    ctype.TypeJsonName,
					Comment: "json",
				},
				{
					Name:    "string",
					Title:   "字符",
					Type:    ctype.TypeStringName,
					Comment: "字符",
				},
				{
					Name:    "text",
					Title:   "文本",
					Type:    ctype.TypeTextName,
					Comment: "文本",
				},
				{
					Name:    "time",
					Title:   "时间",
					Type:    ctype.TypeTimeName,
					Comment: "时间带时区",
				},
				{
					Name:    "file",
					Title:   "文件",
					Type:    ctype.TypeFileName,
					Comment: "文件",
				},
				{
					Name:    "sign1",
					Title:   "签名1",
					Type:    ctype.TypeStringName,
					Comment: "签名1",
				},
				{
					Name:    "sign2",
					Title:   "签名2",
					Type:    ctype.TypeStringName,
					Comment: "签名2",
				},
			},
		},
	}

	return http_tes.HttpCase{
		Desc:    "表单配置-签名",
		Method:  "POST",
		Path:    "/template/create",
		Headers: map[string]string{"Content-Type": "application/json"},
		Body: &dto.Param{
			ModifyParam: dto.ModifyParam{
				SingleTransaction: true,
				Data:              &[]template.Template{tmp},
			},
		},
		ShowBody: true,
		Code:     "200",
		ValidateFunc: []http_tes.ValidateFunc{
			func(t *testing.T, resp *http_tes.Response) {
				res := resp.Resp.Result.(map[string]interface{})
				assert.NotNil(t, res["data"])
			},
		},
	}
}
func templateCreateCrypt(code string) http_tes.HttpCase {

	tmp := template.Template{
		BaseEntity: entity.BaseEntity{
			Remark: "ut-加密",
		},
		Code:   code,
		Status: template.Ready,
		Table: entity.Table{
			Feature: entity.Feature{
				CryptFields: []entity.CryptFieldInfo{
					{
						Enable:    true,
						Algo:      crypt.AlgoSymSm4CbcPkcs7padding,
						SecretKey: []string{"8eadb267ecd6e860"},
						Field:     "enc1",
						SliceSize: 1,
					},
					{
						Enable:        true,
						Algo:          crypt.AlgoSymSm4CbcPkcs7padding,
						SecretKey:     []string{"8eadb267ecd6e864"},
						Field:         "enc2",
						KeepTxtPreCnt: 4,
					},
					{
						Enable:        true,
						Algo:          crypt.AlgoSymSm4CbcPkcs7padding,
						SecretKey:     []string{"9eadb267ecd6e864"},
						Field:         "enc3",
						KeepTxtPreCnt: 4,
						KeepTxtSufCnt: 4,
					},
				},
			},
			Fields: []entity.TableField{
				{
					Name:    "array_int",
					Title:   "整形数组",
					Type:    ctype.TypeArrayIntName,
					Comment: "整形数组",
				},
				{
					Name:    "array_string",
					Title:   "字符数组",
					Type:    ctype.TypeArrayStringName,
					Comment: "字符数组",
				},
				{
					Name:    "bool",
					Title:   "布尔",
					Type:    ctype.TypeBoolName,
					Comment: "布尔",
				},
				{
					Name:    "numeric",
					Title:   "数字",
					Type:    ctype.TypeNumericName,
					Comment: "数字",
				},
				{
					Name:    "geometry",
					Title:   "几何",
					Type:    ctype.TypeGeometryName,
					Comment: "几何",
				},
				{
					Name:    "integer",
					Title:   "整数",
					Type:    ctype.TypeIntegerName,
					Comment: "整数",
				},
				{
					Name:    "bigint",
					Title:   "大整数",
					Type:    ctype.TypeBigIntegerName,
					Comment: "大整数",
				},
				{
					Name:    "json",
					Title:   "json",
					Type:    ctype.TypeJsonName,
					Comment: "json",
				},
				{
					Name:    "string",
					Title:   "字符",
					Type:    ctype.TypeStringName,
					Comment: "字符",
				},
				{
					Name:    "text",
					Title:   "文本",
					Type:    ctype.TypeTextName,
					Comment: "文本",
				},
				{
					Name:    "time",
					Title:   "时间",
					Type:    ctype.TypeTimeName,
					Comment: "时间带时区",
				},
				{
					Name:    "file",
					Title:   "文件",
					Type:    ctype.TypeFileName,
					Comment: "文件",
				},
				{
					Name:    "enc1",
					Title:   "加密1",
					Type:    ctype.TypeStringName,
					Comment: "加密1",
				},
				{
					Name:    "ecn2",
					Title:   "加密2",
					Type:    ctype.TypeStringName,
					Comment: "加密2",
				},
				{
					Name:    "ecn3",
					Title:   "加密3",
					Type:    ctype.TypeStringName,
					Comment: "加密3",
				},
			},
		},
	}

	return http_tes.HttpCase{
		Desc:    "表单配置-加密",
		Method:  "POST",
		Path:    "/template/create",
		Headers: map[string]string{"Content-Type": "application/json"},
		Body: &dto.Param{
			ModifyParam: dto.ModifyParam{
				SingleTransaction: true,
				Data:              &[]template.Template{tmp},
			},
		},
		ShowBody: true,
		Code:     "200",
		ValidateFunc: []http_tes.ValidateFunc{
			func(t *testing.T, resp *http_tes.Response) {
				res := resp.Resp.Result.(map[string]interface{})
				assert.NotNil(t, res["data"])
			},
		},
	}
}
func templateCreateFeature(code string) http_tes.HttpCase {
	tmp := template.Template{
		BaseEntity: entity.BaseEntity{
			Remark: "表单配置-feature",
		},
		Code:   code,
		Status: template.Ready,
		Table: entity.Table{
			Feature: entity.Feature{
				History:      ctype.NewBoolean(true, true),
				FieldUpdater: ctype.NewBoolean(true, true),
			},
			Fields: []entity.TableField{
				{
					Name:    "string",
					Title:   "字符",
					Type:    ctype.TypeStringName,
					Comment: "字符",
				},
			},
		},
	}
	return http_tes.HttpCase{
		Desc:    "表单配置-feature",
		Method:  "POST",
		Path:    "/template/create",
		Headers: map[string]string{"Content-Type": "application/json"},
		Body: &dto.Param{
			ModifyParam: dto.ModifyParam{
				SingleTransaction: true,
				Data:              &[]template.Template{tmp},
			},
		},
		ShowBody: true,
		Code:     "200",
		ValidateFunc: []http_tes.ValidateFunc{
			func(t *testing.T, resp *http_tes.Response) {
				res := resp.Resp.Result.(map[string]interface{})
				assert.NotNil(t, res["data"])
			},
		},
	}
}
