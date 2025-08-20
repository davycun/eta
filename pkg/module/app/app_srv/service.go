package app_srv

import (
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/service"
	"github.com/davycun/eta/pkg/core/ws"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/eta/migrator"
	"github.com/davycun/eta/pkg/module/app"
)

type Service struct {
	service.DefaultService
}

func (s *Service) Migrate(migParam *migrator.MigrateAppParam, result *dto.Result) error {

	var (
		apps   []app.App
		err    error
		appIds = make([]string, 0, 2)
	)

	global.DeleteGorm(global.GetLocalDatabase())

	if len(migParam.AppIds) <= 0 {
		apps, err = app.LoadAllApp()
	} else {
		for _, aid := range migParam.AppIds {
			ap, err1 := app.LoadAppById(global.GetLocalGorm(), aid)
			if err1 != nil {
				return err1
			}
			apps = append(apps, ap)
		}
	}

	for _, ap := range apps {
		logger.Infof("migrate app: %s, 开始...", ap.ID)
		global.DeleteGorm(ap.Database)
		db, err1 := global.LoadGormSetAppId(ap.ID, ap.Database)
		if err1 != nil {
			return err1
		}

		err = migrator.MigrateApp(db, s.GetContext(), migParam)
		if err != nil {
			return err
		}

		if migParam.SendWsMessage {
			ws.SendMessage(constants.WsKeyMigrateApp, "")
		}
		logger.Infof("migrate app: %s, 结束!!!", ap.ID)
		appIds = append(appIds, ap.ID)
	}

	result.Data = appIds
	return nil
}
