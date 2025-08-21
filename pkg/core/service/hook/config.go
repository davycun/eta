package hook

import (
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/es"
	"github.com/davycun/eta/pkg/common/dorm/xa"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/iface"
	"gorm.io/gorm"
	"reflect"
)

type SrvConfigFunc func(o *SrvConfig)
type SrvConfig struct {
	iface.SrvOptions
	Method    iface.Method   //当前调用服务的方法
	CurdType  iface.CurdType //表示当前方法是对数据进行修改还是删除
	DbType    dorm.DbType
	CurDB     *gorm.DB //将会在方法中去执行数据的DB
	TxDB      *gorm.DB //如果是修改数据，这个就是事务DB
	EsApi     *es.Api
	Param     *dto.Param      //请求参数的指针
	Result    *dto.Result     //响应参数的指针
	Values    []reflect.Value //新增或者修改数据的反射
	OldValues any             //修改之前的旧数据切片
	NewValues any             //修改之后的新数据切片
}

func (cfg *SrvConfig) Before() error {
	switch cfg.CurdType {
	case iface.CurdModify:
		return cfg.modifyBefore()
	case iface.CurdRetrieve:
		return cfg.retrieveBefore()
	}
	return nil
}
func (cfg *SrvConfig) After() error {
	switch cfg.CurdType {
	case iface.CurdModify:
		return cfg.modifyAfter()
	case iface.CurdRetrieve:
		return cfg.retrieveAfter()
	}
	return nil
}
func (cfg *SrvConfig) RetrieveEnableEs() bool {
	return global.GetES() != nil && cfg.GetTable().EsRetrieveEnabled()
}

// CommitOrRollback
// 对cfg中的TxDB进行事务提交或者回滚
func (cfg *SrvConfig) CommitOrRollback(err error) error {
	if cfg.TxDB == nil {
		return nil
	}
	if !dorm.InTransaction(cfg.OriginDB) {
		xa.CommitOrRollback(cfg.TxDB, err)
	}
	return nil
}

func NewSrvConfig(curdType iface.CurdType, method iface.Method, opt iface.SrvOptions,
	args *dto.Param, result *dto.Result, mcf ...SrvConfigFunc) *SrvConfig {

	cfg := &SrvConfig{
		Method:   method,
		CurdType: curdType,
		SrvOptions: iface.SrvOptions{
			EC:       opt.EC,
			Ctx:      opt.Ctx,
			OriginDB: opt.OriginDB,
		},
		Param:  args,
		Result: result,
		DbType: dorm.GetDbType(opt.OriginDB),
	}
	for _, ff := range mcf {
		ff(cfg)
	}

	switch curdType {
	case iface.CurdModify:
		if args.Data != nil {
			cfg.NewValues = args.Data
			cfg.Values = utils.ConvertToValueArray(args.Data)
		}
		cfg.TxDB = dorm.Transaction(cfg.OriginDB)
		cfg.CurDB = cfg.TxDB
	case iface.CurdRetrieve:
		cfg.CurDB = cfg.OriginDB
		cfg.TxDB = cfg.OriginDB

		if global.GetES() != nil {
			cfg.EsApi = es.NewApi(global.GetES(), cfg.GetEsIndexName())
		}
	}

	//这里可能已经对AuthFunc进行了设置

	if cfg.Result == nil {
		cfg.Result = &dto.Result{}
	}
	if cfg.Param == nil {
		cfg.Param = &dto.Param{}
	}

	if cfg.EC == nil {
		cfg.EC = iface.GetContextEntityConfig(cfg.GetContext())
	}
	//收尾
	//if cfg.Table == nil {
	//	cfg.Table = entity.GetContextTable(cfg.GetContext())
	//}
	return cfg
}
