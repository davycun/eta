package template

import (
	"errors"
	"fmt"
	"github.com/davycun/eta/pkg/common/crypt"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/duke-git/lancet/v2/slice"
	"strings"
)

// SignValidator 签名校验，校验模板中指定的签名字段是否有定义等
func SignValidator(dt []Template) error {
	supportedType := ctype.GetSupportType()
	for _, v := range dt {
		//基本类型都支持签名
		fieldList := append(supportedType, entity.DefaultVertexColumns...)
		for _, vv := range v.Table.GetFields() {
			//只允许字符串和文本类型的字段进行签名
			if utils.ContainAny(supportedType, vv.Type) {
				fieldList = append(fieldList, vv.Name)
			}
		}
		for _, vv := range v.Table.SignFields {
			// 签名算法
			if !crypt.ExistsAlgo(vv.Algo) {
				return errs.NewClientError(fmt.Sprintf("不支持指定的签名算法[%s],模板Code[%s]", vv.Algo, v.Code))
			}
			if !utils.ContainAll(fieldList, vv.Fields...) {
				return errs.NewClientError(fmt.Sprintf("需要签名的字段在表中不存在或者类型不是string或者text,模板Cocde[%s]", v.Code))
			}
			if !utils.ContainAll(fieldList, vv.Field) {
				return errs.NewClientError(fmt.Sprintf("存储签名值的字段在表中不存在或者类型不是string或者text,模板Cocde[%s]", v.Code))
			}
		}
	}
	return nil
}

// EncryptValidator 加密校验，校验模板中指定的加密字段是否有定义等
func EncryptValidator(dt []Template) error {
	for _, v := range dt {
		for _, cryptInfo := range v.Table.CryptFields {
			// 加密算法
			if crypt.ExistsAlgo(cryptInfo.Algo) {
				return errs.NewClientError(fmt.Sprintf("不支持指定的签名算法[%s],模板Code[%s]", cryptInfo.Algo, v.Code))
			}
			err := secretKeyValidator(cryptInfo)
			if err != nil {
				return err
			}
			// 加密字段是否已定义字段
			field := cryptInfo.Field
			foundField := false
			for _, f := range v.Table.GetFields() {
				if field == f.Name {
					if !slice.Contain([]string{ctype.TypeStringName, ctype.TypeTextName}, f.Type) {
						return errors.New(fmt.Sprintf("[%s]的加密配置[field]有误,字段类型必须是%s", v.Code, ctype.TypeStringName))
					}
					foundField = true
					break
				}
			}
			if !foundField {
				return errors.New(fmt.Sprintf("[%s]的加密配置[field]有误", v.Code))
			}
		}
	}
	return nil
}

// 校验加密信息中定义的算法及密钥信息是否合法等
func secretKeyValidator(cryptInfo entity.CryptFieldInfo) error {
	var (
		algo      = cryptInfo.Algo
		secretKey = cryptInfo.SecretKey
	)

	// 密钥数量
	if crypt.ExistsAlgo(cryptInfo.Algo) {
		if len(secretKey) != 1 {
			return errors.New(fmt.Sprintf("加密配置密钥数量不符合算法要求,%s", algo))
		}
		keyLen := len(secretKey[0])
		if strings.HasPrefix(algo, "sm4_") && keyLen != 16 {
			return errors.New(fmt.Sprintf("加密配置,密钥长度必须是16,%s", algo))
		} else if strings.HasPrefix(algo, "aes_") && (keyLen != 16 && keyLen != 24 && keyLen != 32) {
			return errors.New(fmt.Sprintf("加密配置,密钥长度只能是16/24/32,%s", algo))
		}
	} else {
		return errors.New(fmt.Sprintf("加密配置密钥配置不正确,%s", algo))
	}

	return nil
}

// RaDbFieldsValidator RA 校验
func RaDbFieldsValidator(dt []Template) error {
	for _, v := range dt {
		if len(v.Table.RaDbFields) <= 0 {
			continue
		}

		raDbFields := v.Table.RaDbFields
		fields := slice.Map(v.Table.GetFields(), func(_ int, v entity.TableField) string { return v.Name })
		if !slice.ContainSubSlice(fields, raDbFields) {
			return errors.New(fmt.Sprintf("[%s]的RA配置有误,ra字段需要全部取自于表字段", v.Code))
		}

	}
	return nil
}
