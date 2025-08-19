package reload

import (
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/dsync"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/iface"
)

func GetSyncOption(srv iface.Service, args *dto.Param) *dsync.SyncOption {
	option := dsync.GetSyncOption(args)
	return option
}

func buildNewSyncOption(srv iface.Service, args *dto.Param) {
	sync2es := true
	err := checkEsIndex(srv)
	if err != nil {
		sync2es = false
	}

	so := *(args.Extra.(*dsync.SyncOption))
	so.SyncToEs = so.SyncToEs && sync2es
	so.UpdateDbRaContent = so.UpdateDbRaContent && needUpdateRaContent(srv, &so)
	so.UpdateDbEncrypt = so.UpdateDbEncrypt && needUpdateEncrypt(srv, &so)
	so.UpdateDbSign = so.UpdateDbSign && needUpdateSign(srv, &so)

	args.Extra = &so
}

func needSync(srv iface.Service) bool {
	// 自己实现了 ReloadService 接口，就需要同步
	if _, ok := srv.(dsync.ReloadService); ok {
		return true
	}
	// 有需要 reload 的事项，就需要同步
	return needUpdateDbRaContent(srv) || needUpdateDbSign(srv) || needUpdateDbEncrypt(srv) || needSync2Es(srv)
}
func needUpdateDbRaContent(srv iface.Service) bool {
	fields := entity.GetRaDbFields(srv.NewEntityPointer())
	if len(fields) > 0 {
		return true
	}
	return false
}
func needUpdateDbSign(srv iface.Service) bool {
	tb := srv.GetTable()
	return tb != nil && tb.NeedSign()
}
func needUpdateDbEncrypt(srv iface.Service) bool {
	tb := srv.GetTable()
	return tb != nil && tb.NeedCrypt()
}
func needSync2Es(srv iface.Service) bool {
	tb := srv.GetTable()
	return tb != nil && ctype.Bool(tb.EsEnable) && global.GetES() != nil
}

func needUpdateRaContent(srv iface.Service, opt *dsync.SyncOption) bool {
	if !opt.UpdateDbRaContent {
		return false
	}
	fields := entity.GetRaDbFields(srv.NewEntityPointer())
	if len(fields) > 0 {
		return true
	}
	return false
}
func needUpdateSign(srv iface.Service, opt *dsync.SyncOption) bool {
	if !opt.UpdateDbSign {
		return false
	}
	tb := srv.GetTable()
	return tb != nil && tb.NeedSign()
}
func needUpdateEncrypt(srv iface.Service, opt *dsync.SyncOption) bool {
	if !opt.UpdateDbEncrypt {
		return false
	}
	tb := srv.GetTable()

	return tb != nil && tb.NeedCrypt()
}
