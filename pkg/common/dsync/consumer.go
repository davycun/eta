package dsync

import (
	"context"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/run"
	"sync"
	"sync/atomic"
)

type defaultConsumer struct {
	goroutineSize int
	saver         DataSaver
}

func NewDefaultConsumer(saver DataSaver, goSize int) Consumer {

	dc := &defaultConsumer{
		goroutineSize: goSize,
		saver:         saver,
	}
	return dc
}

func (s *defaultConsumer) Consume(c context.Context, args any, ch <-chan any) error {
	var (
		goroutineTotal = atomic.Int32{}
		routineSize    = s.goroutineSize
		err            error
		over           = false
		wg             = &sync.WaitGroup{}
	)
	goroutineTotal.Add(1)
	for {
		select {
		case <-c.Done():
			return err
		default:
			if over || err != nil {
				wg.Wait()
				return err
			}
			var src any
			select {
			case d := <-ch:
				switch x := d.(type) {
				case string:
					//代表结束了
					if x == overFlag {
						over = true
						goto saveData
					}
				default:
					src = d
				}
			}

		saveData:
			logger.Infof("当前的routineSize: %d", goroutineTotal.Load())
			if goroutineTotal.Load() >= int32(routineSize) {
				logger.Infof("main开始插入数据...")
				//这里需要用一个新的接受错误，否则可能会导致老的已经发生的错误被新nil覆盖，从而获取不到发生的错误
				err = errs.Cover(err, s.saver(args, src))
				if err != nil {
					logger.Errorf("main同步数据错误%s", err.Error())
				}
			} else {
				wg.Add(1)
				goroutineTotal.Add(1)
				run.Go(func() {
					defer wg.Done()
					logger.Infof("goroutine开始插入数据...")
					//这里需要用一个新的接受错误，否则可能会导致老的已经发生的错误被新nil覆盖，从而获取不到发生的错误
					err = errs.Cover(err, s.saver(args, src))
					if err != nil {
						logger.Errorf("goroutine同步数据错误%s", err.Error())
					}
					goroutineTotal.Add(-1)
				})
			}
		}
	}
}
