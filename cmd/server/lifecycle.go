package server

import (
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/duke-git/lancet/v2/maputil"
	"math"
	"slices"
)

// 下面定义的生命周期阶段，按顺序从小到大执行
const (
	initConfig             Stage = 0
	InitPlugin             Stage = 10          //注册一些回调函数，包括gorm的插件等，这个需要再InitApplication之前
	InitApplication        Stage = 20          //创建Application,初始化gin、redis、db、es等
	InitMiddleware         Stage = 30          //初始化gin的中间件，需要放在InitRouter及InitModules之前
	InitValidator          Stage = 40          //添加binding自定义校验器
	InitData               Stage = 50          //主要是注册一些EntityConfig及一些默认的初始化数据（在Migrate阶段之后会插入到数据库）
	InitEntityConfigRouter Stage = 60          // 主要是注册路由，注册是方式是通过InitData阶段注册的EntityConfig
	InitModules            Stage = 70          //初始化每个模块，每个模块中可以注册各个模块自己服务的回调，添加模块自定义的gin路由等
	InitMigrator           Stage = 80          //必须在Migrate之前，主要是注册migrate的回调
	Migrate                Stage = 90          //启动执行Migrate动作
	StartServer            Stage = math.MaxInt //必须最后执行，启动http服务器
	Shutdown               Stage = 1000000     //这个是特殊节点在执行生命周期函数的时候不会去执行
)

type (
	Stage         int
	LifeCycleFunc func() error
)

var (
	lifeCycles = map[Stage]LifeCycleFunc{}
)

func AddLifeCycle(stage Stage, fc LifeCycleFunc) {
	if _, ok := lifeCycles[stage]; ok || stage == Shutdown || stage < 0 {
		logger.Errorf("Add new life cycle err,because stage[%d] had exists or stage less then zero!", stage)
		return
	}
	lifeCycles[stage] = fc
}
func stageExists(stage Stage) bool {
	_, ok := lifeCycles[stage]
	return ok
}

// 根据Stage的顺序调用生命周期函数
func callLifeCycle() error {
	var (
		err       error
		stageList = maputil.Keys(lifeCycles)
	)
	slices.Sort(stageList)
	for _, stage := range stageList {
		//shutdown是特殊阶段执行的 请看watchSignal函数
		if stage == Shutdown {
			continue
		}
		err = callLifeCycleHook(stage, BeforeStage)
		if err != nil {
			return err
		}
		err = lifeCycles[stage]()
		if err != nil {
			return err
		}
		err = callLifeCycleHook(stage, AfterStage)
		if err != nil {
			return err
		}
	}
	return nil
}
