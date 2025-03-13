package sms_srv

import (
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/service"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/sms"
	"github.com/davycun/eta/pkg/module/sms/sms_sender"
	"time"
)

func init() {
	hook.AddModifyCallback(constants.TableSmsTask, modifyCallbacks)
}

func modifyCallbacks(cfg *hook.SrvConfig, pos hook.CallbackPosition) error {
	return caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			return hook.BeforeCreate(cfg, pos, func(cfg *hook.SrvConfig, newValues []sms.Task) error {
				for i, v := range newValues {
					if len(v.TargetList) < 1 {
						return errs.NewClientError("没有指定发送目标")
					}
					if ctype.IsValid(v.PlainTime) && v.PlainTime.Data.Sub(time.Now()) > 0 {
						newValues[i].Status = ctype.NewStringPrt("待发送")
					} else {
						newValues[i].Status = ctype.NewStringPrt("已发送")
					}
					newValues[i].TargetTotal = ctype.NewIntPrt(int64(len(v.TargetList)))
				}
				return nil
			})
		}).
		Call(func(cl *caller.Caller) error {
			return hook.AfterCreate(cfg, pos, func(cfg *hook.SrvConfig, newValues []sms.Task) error {
				targetResult, err := sms_sender.SendSms(cfg.TxDB, newValues...)
				if err != nil {
					return err
				}
				srv := service.NewService(constants.TableSmsTarget, cfg.Ctx, cfg.TxDB)
				return srv.Create(&dto.Param{
					ModifyParam: dto.ModifyParam{
						Data: &targetResult,
					},
				}, &dto.Result{})
			})
		}).Err
}
