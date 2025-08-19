package migrator

import (
	"bytes"
	"context"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/es"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	jsoniter "github.com/json-iterator/go"
	"gorm.io/gorm"
	"net/http"
	"time"
)

/*
MigrateElasticsearch 迁移ES
param: tableName:ES参数
*/
func MigrateElasticsearch(db *gorm.DB, param map[string]EsParam) error {
	var (
		esApi = global.GetES()
		esIdx = iface.GetEsEntityConfig()
	)
	if esApi == nil {
		return nil
	}
	if param != nil && len(param) > 0 {
		slice.ForEach(esIdx, func(i int, v entity.Table) {
			tbName := v.GetTableName()
			if ep, ok := param[tbName]; ok {
				esIdx[i].EsSettings = ep.Settings
			}
			if ep, ok := param[tbName]; ok {
				esIdx[i].EsSettings = ep.Settings
			}
			if ep, ok := param[tbName]; ok {
				esIdx[i].EsSettings = ep.Settings
			}
		})
	}

	return CreateEsIndex(db, esIdx...)
}

func CreateEsIndex(db *gorm.DB, tbs ...entity.Table) error {
	for _, v := range tbs {
		var (
			err error
			et  = v.NewEsEntityPointer()
		)
		err = createSingleIndex(et, entity.GetEsIndexName(dorm.GetDbSchema(db), v.GetTableName()), v.EsSettings)
		if err != nil {
			return err
		}
	}
	return nil
}

func createSingleIndex(obj any, idxName string, sts map[string]interface{}) error {
	if indexExists(idxName) {
		//TODO 更新
		//props := map[string]interface{}{
		//	"properties": field2EsProps(getStructEsFields(reflect.TypeOf(obj))...),
		//}
		//bs, err := jsoniter.Marshal(props)
		//if err != nil {
		//	return err
		//}
		//
		//_, err = global.GetES().EsTypedApi.Indices.PutMapping(idxName).Raw(bytes.NewReader(bs)).Do(context.Background())
		//if err != nil {
		//	return err
		//}
		return nil
	}

	bs, err := ResolveEsIndex(obj, sts)
	if err != nil {
		return err
	}

	idxDataName := idxName + constants.EsIndexSubfix

	start := time.Now()
	resp, err := global.GetES().EsApi.Indices.Create(idxDataName, func(request *esapi.IndicesCreateRequest) {
		request.Body = bytes.NewReader(bs)
	})
	if resp != nil {
		es.LatencyLog(start, idxName, "create_index", bs, resp.StatusCode)
	}
	if err != nil {
		return err
	}
	if resp != nil {
		if resp.StatusCode != http.StatusOK {
			logger.Errorf("create es index error %s", resp.String())
		} else {
			defer resp.Body.Close()
		}
	}

	// 创建 alias
	alias, err := global.GetES().EsApi.Indices.PutAlias([]string{idxDataName}, idxName)
	if alias != nil {
		es.LatencyLog(start, idxName, "create_index_alias", bs, alias.StatusCode)
	}
	if err != nil {
		return err
	}
	if alias != nil {
		if alias.StatusCode != http.StatusOK {
			logger.Errorf("Put es index alias error %s", alias.String())
		} else {
			defer alias.Body.Close()
		}
	}

	return err
}

func ResolveEsIndex(obj any, sts map[string]interface{}) ([]byte, error) {

	var (
		numOfShards   = global.GetConfig().EsConfig.NumberOfShards
		numOfReplicas = global.GetConfig().EsConfig.NumberOfReplicas
	)
	if sts == nil {
		sts = make(map[string]interface{})
	}

	if numOfShards > 0 {
		sts["number_of_shards"] = numOfShards
	}
	if numOfReplicas >= 0 {
		sts["number_of_replicas"] = numOfReplicas
	}
	idx := map[string]interface{}{
		"settings": sts,
		"mappings": map[string]interface{}{
			"properties": es.GetEsMapping(obj),
		},
	}

	return jsoniter.Marshal(idx)
}

func indexExists(name string) bool {

	//resp, err := global.GetES().EsApi.Indices.Exists([]string{name})
	//if err != nil {
	//	logger.Errorf("indics exists err %s", err)
	//}
	//return resp.StatusCode == http.StatusNotFound

	exists, err := global.GetES().EsTypedApi.Indices.Exists(name).Do(context.Background())
	if err != nil {
		logger.Errorf("indics exists err %s", err)
	}
	if exists {
		return exists
	}
	exists, err = global.GetES().EsTypedApi.Indices.ExistsAlias(name).Do(context.Background())
	if err != nil {
		logger.Errorf("indics alias exists err %s", err)
	}
	return exists

}
