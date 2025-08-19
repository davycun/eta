package dsync

import (
	"errors"
	"fmt"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/panjf2000/ants/v2"
	"sync"
)

type AntsSyncService struct {
	AntsProducer
	AntsConsumer
	Pool *ants.Pool
}

func NewAntsSyncService(loader DataLoader, saver DataSaver, option SyncOption) SyncService {
	var (
		poolSize = 1
	)
	if option.ConsumerGoSize > 1 {
		poolSize = option.ConsumerGoSize
	}
	as := &AntsSyncService{
		AntsProducer: NewAntsProducer(loader),
		AntsConsumer: NewAntsConsumer(saver),
		Pool:         errs.TryUnwrap(ants.NewPool(poolSize, ants.WithNonblocking(false))),
	}
	return as
}

func (a *AntsSyncService) Release() {
	if a.Pool != nil {
		a.Pool.Release()
		logger.Infof("pool released!")
	}
}

func (a *AntsSyncService) Sync(producerArgs, consumerArgs any) error {
	logger.Infof("sync start...")
	var (
		wg      sync.WaitGroup
		taskErr error // task 执行错误
	)
	defer a.Release()

	for {
		data, err := a.Produce(producerArgs)
		if err != nil {
			return err
		}
		if utils.IsEmptySlice(data) {
			break
		}

		pool := a.Pool
		wg.Add(1)
		err1 := pool.Submit(func() {
			defer func() {
				wg.Done()
				if r := recover(); r != nil {
					logger.Errorf("AntsSyncService consume panic. %v", r)
					taskErr = errors.New(fmt.Sprintf("AntsSyncService consume panic. %v", r))
				}
			}()
			logger.Debugf("pool cap: %d, running:%d, free:%d, waiting:%d", pool.Cap(), pool.Running(), pool.Free(), pool.Waiting())
			taskErr = a.Consume(consumerArgs, data)
			if taskErr != nil {
				logger.Error(taskErr)
			}
		})
		if err1 != nil {
			logger.Errorf("AntsSyncService submit consume task err:%v\n", err1)
			wg.Done()
			return err1
		}
		if taskErr != nil {
			// 每次循环时，都判断 taskErr 是否为空，如果不为空，就直接返回
			return taskErr
		}
	}
	wg.Wait()
	if taskErr != nil {
		// 任务全部执行结束后，如果 taskErr 不为空，就直接返回
		return taskErr
	}
	logger.Infof("sync end!!!")
	return nil
}

type AntsProducer struct {
	loader DataLoader
}

func NewAntsProducer(loader DataLoader) AntsProducer {
	pd := AntsProducer{
		loader: loader,
	}
	return pd
}

// Produce 返回的 data 是 nil 时，同步结束
func (p *AntsProducer) Produce(args any) (data any, err error) {
	data, _, err = p.loader(args)
	return
}

type AntsConsumer struct {
	saver DataSaver
}

func NewAntsConsumer(saver DataSaver) AntsConsumer {
	dc := AntsConsumer{
		saver: saver,
	}
	return dc
}

func (c *AntsConsumer) Consume(args any, data any) error {
	return c.saver(args, data)
}
