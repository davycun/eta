package migrator_test

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/dorm/es/es_api"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/eta/migrator"
	"github.com/duke-git/lancet/v2/slice"
	"reflect"
	"testing"
)

func TestGenCrtIdxReq(t *testing.T) {
	schema := "Eta_dev_frontend"

	esIdx := slice.Filter(iface.GetMigrateEntityConfig(), func(_ int, v entity.Table) bool {
		//return slice.Contain(v.EnableDbType, dorm.ES)
		return ctype.Bool(v.EsEnable)
	})
	slice.ForEach(esIdx, func(index int, item entity.Table) {
		if ctype.Bool(item.EsEnable) {
			et := reflect.New(item.EntityType)
			idxName := fmt.Sprintf("%s_%s", schema, entity.GetTableName(et))
			esIndex, _ := migrator.ResolveEsIndex(et, es_api.DefaultSetting())
			fmt.Printf("PUT /%s\n", idxName)
			fmt.Printf("%s\n\n", string(esIndex))
		}
	})
}
