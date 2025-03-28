package user_srv

import (
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/dorm/filter"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/core/service"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/core/ws"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/authorize/auth"
	"github.com/davycun/eta/pkg/module/dept"
	"github.com/davycun/eta/pkg/module/user"
	"github.com/davycun/eta/pkg/module/user/user2dept"
	"github.com/duke-git/lancet/v2/maputil"
	"github.com/duke-git/lancet/v2/slice"
)

func modifyCallbackUser2Dept(cfg *hook.SrvConfig, pos hook.CallbackPosition) error {

	err := caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			return hook.AfterCreate(cfg, pos, func(cfg *hook.SrvConfig, newValues []user2dept.User2Dept) error {
				cleanCache(newValues)
				return nil
			})
		}).
		Call(func(cl *caller.Caller) error {
			return hook.AfterModify(cfg, pos, func(cfg *hook.SrvConfig, oldValues []user2dept.User2Dept, newValues []user2dept.User2Dept) error {
				logout(oldValues, newValues)
				//因为是清除缓存，所以要清除更新之前的缓存，而是传入NewValues
				cleanCache(oldValues)
				return nil
			}, iface.MethodUpdate, iface.MethodUpdateByFilters, iface.MethodDelete, iface.MethodDeleteByFilters)
		}).
		Call(func(cl *caller.Caller) error {
			return hook.AfterModify(cfg, pos, notifyUser2DeptChanged)
		}).Err
	return err
}

func logout(oldValues, newValues []user2dept.User2Dept) {
	oldFromTo := valuesToMap(oldValues)
	newFromTo := valuesToMap(newValues)

	toLogout := make(map[string][]string) // userId:[]deptId
	maputil.ForEach(oldFromTo, func(userId string, olds []string) {
		if _, ok := newFromTo[userId]; !ok {
			toLogout[userId] = olds
		} else {
			news := newFromTo[userId]
			diff := slice.Difference(olds, news)
			if len(diff) > 0 {
				// 判断toLogout里是否有值，有值就合并
				if _, ok1 := toLogout[userId]; ok1 {
					toLogout[userId] = append(toLogout[userId], diff...)
				} else {
					// toLogout里没有值，则直接赋值
					toLogout[userId] = diff
				}
			}
		}
	})
	// 找出 toLogout 里已登录的用户，把登录到当前部门的用户踢下线
	maputil.ForEach(toLogout, func(userId string, deptIds []string) {
		err := user.DelUserTokenByIdAndDeptId(userId, deptIds...)
		if err != nil {
			logger.Errorf("user2dept callback logout error: %v", err)
		}
	})
}

func cleanCache(data []user2dept.User2Dept) {
	uIds := make([]string, 0, len(data))
	for _, v := range data {
		uIds = append(uIds, v.FromId)
	}
	auth.DelUserRoleCache(uIds...)
	dept.DelUser2DeptCache(uIds...)
}

// userId -> []deptIds
func valuesToMap(dt []user2dept.User2Dept) map[string][]string {
	fromIdToId := make(map[string][]string)
	if dt == nil {
		return fromIdToId
	}
	if len(dt) <= 0 {
		return fromIdToId
	}
	for _, v := range dt {
		if _, ok := fromIdToId[v.FromId]; !ok {
			fromIdToId[v.FromId] = make([]string, 0)
		}
		fromIdToId[v.FromId] = append(fromIdToId[v.FromId], v.ToId)
	}
	return fromIdToId
}

func modifyCallbackDept(cfg *hook.SrvConfig, pos hook.CallbackPosition) error {
	return hook.AfterDelete(cfg, pos, func(cfg *hook.SrvConfig, oldValues []dept.Department) error {
		if len(oldValues) < 1 {
			return nil
		}
		var (
			u2dArgs     dto.Param
			u2dRes      dto.Result
			batchSize   = 1000
			u2dSvc, err = service.NewService(constants.TableUser2Dept, cfg.Ctx, cfg.TxDB)
		)
		if err != nil {
			return err
		}

		for _, depts := range slice.Chunk(oldValues, batchSize) {
			u2dArgs.Filters = []filter.Filter{
				{
					LogicalOperator: filter.And,
					Column:          "to_id",
					Operator:        filter.IN,
					Value:           slice.Map(depts, func(i int, item dept.Department) string { return item.ID }),
				},
			}
			u2dArgs.Data = &user2dept.User2Dept{}
			err = u2dSvc.DeleteByFilters(&u2dArgs, &u2dRes)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// user2dept变更，给相关人发送 websocket 通知
func notifyUser2DeptChanged(cfg *hook.SrvConfig, oldValues []user2dept.User2Dept, newValues []user2dept.User2Dept) error {
	vs := slice.Concat(oldValues, newValues)
	uids := slice.Unique(slice.Map(vs, func(i int, v user2dept.User2Dept) string { return v.FromId }))
	slice.ForEach(uids, func(i int, v string) {
		ws.SendMessage(constants.WsKeyUser2DeptChanged, "", v)
	})
	return nil
}
