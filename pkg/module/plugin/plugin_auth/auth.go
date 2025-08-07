package plugin_auth

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/cache"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/authorize/auth"
	"github.com/davycun/eta/pkg/module/setting"
	"strings"
)

//permission表定义了很多权限（filter）

// CleanUserAuthCache
// 当用户的权限变化的时候需要清空下缓存
func CleanUserAuthCache(userId string) {
	_, err := cache.Del(constants.RedisKey(constants.UserRoleIdsKey, userId))
	if err != nil {
		logger.Errorf("clean user role cache err %s", err)
	}
	_, err = cache.DelKeyPattern(fmt.Sprintf(`"auth:permission:%s:*"`, userId))
	if err != nil {
		logger.Errorf("clean user permission cache err %s", err)
	}
	_, err = cache.DelKeyPattern(fmt.Sprintf(`"auth:data_permission:%s:*"`, userId))
	if err != nil {
		logger.Errorf("clean user permission cache err %s", err)
	}
}

func AuthFilter(cfg *hook.SrvConfig, pos hook.CallbackPosition) error {
	if pos != hook.CallbackBefore {
		return nil
	}

	var (
		err            error
		authType       = auth.Read
		perms          []auth.Permission
		UseCurDeptAuth = cfg.Param.UseCurDeptAuth
		userId         = cfg.Ctx.GetContextUserId()
		curDeptId      = cfg.Ctx.GetContextCurrentDeptId()
		appDb          = cfg.Ctx.GetAppGorm()
		httpMethod     = utils.GetHttpMethod(cfg.GetContext().GetGinContext())
		uri            = utils.GetUrlPath(cfg.GetContext().GetGinContext())
	)

	if cfg.Ctx.GetContextIsManager() {
		return nil
	}

	//避免/oauth2/login 这类接口
	//if global.IsIgnoreUri(cfg.Ctx.GetGinContext().Request.RequestURI) {
	//	return cfg.CurDB, nil
	//}
	if setting.IsIgnoreTokenUri(cfg.GetDB(), httpMethod, uri) ||
		setting.IsIgnoreAuthUri(cfg.GetDB(), httpMethod, uri) {
		return nil
	}

	//设置当前的权限类型
	switch cfg.Method {
	case iface.MethodUpdate, iface.MethodUpdateByFilters:
		authType = auth.Edit
	case iface.MethodDelete, iface.MethodDeleteByFilters:
		authType = auth.Delete
	default:
		authType = auth.Read
	}

	//TODO 暂时的做法是，没有配置角色权限就不限制，也就是说没有配置的时候反而是admin
	if UseCurDeptAuth {
		perms, err = auth.FetchRolePermission(appDb, curDeptId, cfg.GetTableName(), authType)
	} else {
		perms, err = auth.FetchUserPermission(appDb, userId, cfg.GetTableName(), authType)
	}

	if err != nil {
		return err
	}

	if len(perms) > 0 {
		for _, v := range perms {
			//通用接口只支持单表的权限
			if !tbNameEquals(v.TbName, cfg.GetTableName()) && v.TbName != auth.PermissionForAll {
				continue
			}
			cfg.Param.AuthFilters = append(cfg.Param.AuthFilters, v.Filters...)
			cfg.Param.AuthRecursiveFilters = append(cfg.Param.AuthRecursiveFilters, v.RecursiveFilters...)
		}
	}
	//添加关联权限
	auth2roleFilters, err := auth.BuildJoinAuth2RoleFilter(appDb, cfg.GetTableName(), cfg.GetTableName(), userId, authType)
	cfg.Param.Auth2RoleFilters = append(cfg.Param.Auth2RoleFilters, auth2roleFilters...)
	return err
}

func authCreate(cfg *hook.SrvConfig, pos hook.CallbackPosition) error {

	return nil
}

func tbNameEquals(src, target string) bool {
	if src == target {
		return true
	}
	src = strings.TrimLeft(src, constants.TableTemplatePrefix)
	target = strings.TrimLeft(target, constants.TableTemplatePrefix)
	return src == target
}
