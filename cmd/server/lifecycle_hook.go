package server

import (
	"github.com/davycun/eta/pkg/common/logger"
)

const (
	BeforeStage HookPosition = 1
	AfterStage  HookPosition = 2
)

var (
	lifeCycleHookMap = make(map[Stage]map[HookPosition][]LifeCycleHook)
)

type (
	HookPosition  int
	LifeCycleHook func() error
)

func AddLifeCycleHook(stage Stage, pos HookPosition, cb ...LifeCycleHook) {
	if len(cb) < 1 {
		return
	}
	if pos != BeforeStage && pos != AfterStage {
		logger.Errorf("add hook not support position %d,only support BeforeStage and AfterStage", pos)
		return
	}

	if !stageExists(stage) {
		logger.Errorf("add hook stage[%d] not exists, you need add LifeCycle first with stage[%d] first", stage, stage)
		return
	}

	if _, ok := lifeCycleHookMap[stage]; !ok {
		lifeCycleHookMap[stage] = make(map[HookPosition][]LifeCycleHook)
	}

	if _, ok := lifeCycleHookMap[stage][pos]; !ok {
		lifeCycleHookMap[stage][pos] = make([]LifeCycleHook, 0)
	}
	lifeCycleHookMap[stage][pos] = append(lifeCycleHookMap[stage][pos], cb...)
}

func callLifeCycleHook(stage Stage, pos HookPosition) error {
	if _, ok := lifeCycleHookMap[stage]; !ok {
		return nil
	}
	if cbList, ok := lifeCycleHookMap[stage][pos]; ok {
		for _, cb := range cbList {
			err := cb()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
