package plugin_crypt

import (
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/module/app"
	"github.com/davycun/eta/pkg/module/setting"
	"strings"
)

// StringScanner
// 把解密注册在ctype.String或者ctype.Text的Scanner上
// 这样在单独查询sql进行反序列化的时候也能进行解密
func StringScanner(dt any) (any, error) {
	src := ctype.ToString(dt)
	if !strings.HasPrefix(src, metaFlag) {
		return src, nil
	}

	encryptPrefix, encryptStr := splitEncryptData(src)
	idTb := strings.Split(strings.ReplaceAll(encryptPrefix, metaFlag, ""), dataSplit)
	if len(idTb) != 3 {
		return src, nil
	}
	var (
		id     = idTb[0]
		tbName = idTb[1]
		field  = idTb[2]
	)

	ap, err := app.LoadAppById(global.GetLocalGorm(), id)
	if err != nil {
		logger.Errorf("splitEncryptData loadd app by id err %s", err)
		return src, nil
	}
	db, err := global.LoadGormSetAppId(ap.ID, ap.Database)
	if err != nil {
		logger.Errorf("splitEncryptData loadd app by id err %s", err)
		return src, nil
	}
	tb, b := setting.GetTableConfig(db, tbName)
	cryptFieldInfo := tb.GetCryptInfoByField(field)

	if !b || !cryptFieldInfo.Enable {
		return src, nil
	}

	plainText, b := decryptData(cryptFieldInfo, []rune(encryptStr))
	if b {
		return plainText, nil
	}
	return src, nil
}
