package dto

import "github.com/davycun/eta/pkg/common/dsync"

type SyncArgs struct {
	CurId string
}

func GetSyncOption(args *Param) *dsync.SyncOption {
	qe := &dsync.SyncOption{
		ConsumerGoSize: 10,
		ChanSize:       10,
	}
	if args.Extra == nil {
		args.Extra = qe
	} else {
		if x, ok := args.Extra.(*dsync.SyncOption); ok {
			if x.ConsumerGoSize < 1 {
				x.ConsumerGoSize = 10
			}
			if x.ChanSize < 1 {
				x.ChanSize = 10
			}
			return x
		}
	}
	return args.Extra.(*dsync.SyncOption)
}
