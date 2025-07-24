package sqlbd

import (
	"fmt"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/core/service/hook"
	"sync"
)

const (
	CountSql = "count_sql"
	ListSql  = "list_sql"
	TotalSql = "total_sql"
)

var (
	sqlBuilderStore = sync.Map{} //存储的内容是tableName -> map[Methods]SqlList
)

const (
	BuildForAllTable = "all_table_sql_builder" //这是为了回调所有表（所有的service）
)

type BuildOption struct {
	TableName string
	Method    iface.Method
	Name      string //回调的名称
}
type CallbackOptionFunc func(*BuildOption)

type BuildSqlWrapper struct {
	BuildSql BuildSql
	BuildOption
}

func AddSqlBuilder(tableName string, sqlBd BuildSql, method iface.Method, buildOptions ...CallbackOptionFunc) {
	var (
		wrapper = BuildSqlWrapper{
			BuildSql: sqlBd,
		}
	)

	for _, f := range buildOptions {
		f(&wrapper.BuildOption)
	}
	wrapper.TableName = tableName
	wrapper.Method = method
	sqlBuilderStore.Store(getKey(tableName, method), wrapper)
}

func getKey(tableName string, method iface.Method) string {

	return fmt.Sprintf("%s_%s", tableName, method)
}

func getSqlBuilder(tbName string, method iface.Method) BuildSqlWrapper {

	bs, ok := sqlBuilderStore.Load(getKey(tbName, method))
	if ok {
		return bs.(BuildSqlWrapper)
	}
	return BuildSqlWrapper{}
}

func Build(cfg *hook.SrvConfig, tableName string, method iface.Method) (*SqlList, error) {
	bs := getSqlBuilder(tableName, method)
	if bs.BuildSql == nil {
		bs = getSqlBuilder(BuildForAllTable, method)
	}

	if bs.BuildSql == nil {
		return nil, fmt.Errorf("not found sql builder for %s.%s", tableName, method)
	}
	return bs.BuildSql(cfg)
}
