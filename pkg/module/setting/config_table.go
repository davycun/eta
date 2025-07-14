package setting

import (
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/iface"
	"gorm.io/gorm"
)

//配置示例
// 下面这个配置表示user表的password和phone一起签名，并且对phone字段进行加密存储，只是没有启用，如果要启用，对应Enable设置为true即可
// st := Setting{}
//	st.Namespace = constants.NamespaceEta
//	st.Category = ConfigTableCategory
//	st.Name = ConfigTableName
//	tbConfig := TableConfig{
//		Tables: map[string]entity.Table{
//			constants.TableUser: {
//				TableName: constants.TableUser,
//				Feature: entity.Feature{
//					SignFields: []entity.SignFieldsInfo{
//						{
//							Enable: false,
//							Field:  "sign", //签名值存储的字段确保表中确实存在此字段
//							Fields: []string{
//								"password",
//								"phone",
//							},
//							VerifyField: "sign_matched",
//							Algo:        "hmac_sm3",
//							Key:         "citizen",
//						},
//					},
//					CryptFields: []entity.CryptFieldInfo{
//						{Enable: false, Field: "phone", Algo: crypt.AlgoSymSm4CbcPkcs7padding, SecretKey: []string{"isatest123456789"}, SliceSize: 1},
//					},
//				},
//			},
//		},
//	}
//	st.Content = ctype.NewJson(tbConfig)

type TableConfig struct {
	Tables map[string]entity.Table `json:"tables,omitempty"` //tableName-> entity.Table
}

// GetTableConfig
// appId可以为空，代表查询通用的配置，如果appId不为空，但是也查不到对应配置，那就返回默认配置
func GetTableConfig(db *gorm.DB, tableName string) (entity.Table, bool) {

	cfg, err := GetConfig[TableConfig](db, ConfigTableCategory, ConfigTableName)
	if err != nil {
		logger.Errorf("load table config err %s", err)
		return entity.Table{}, false
	}

	if cfg.Tables == nil {
		cfg.Tables = map[string]entity.Table{}
	}
	//如果在appDB中没有找到配置，又或者找到配置，但是配置里没有tableName的配置，那么再取localDb里面去找
	if _, ok := cfg.Tables[tableName]; !ok && isAppDb(db) {
		cfg, err = GetConfig[TableConfig](global.GetLocalGorm(), ConfigTableCategory, ConfigTableName)
		if cfg.Tables == nil {
			cfg.Tables = map[string]entity.Table{}
		}
	}
	if x, ok := cfg.Tables[tableName]; ok {
		return x, true
	}
	//如果localDb的全局配置里还是找不到就从实体的配置找（实体的特性是实现接口)
	ec, b := iface.GetEntityConfigByTableName(tableName)
	if b {
		return *ec.GetTable(), true
	}
	return entity.Table{}, false
}
