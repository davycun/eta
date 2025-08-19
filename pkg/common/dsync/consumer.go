package dsync

import (
	"context"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/run"
	"sync"
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
		routineSize      = s.goroutineSize
		err              error
		over             = false
		wg               = &sync.WaitGroup{}
		goroutineChannel = make(chan int, routineSize)
	)
	for i := 0; i < routineSize; i++ {
		goroutineChannel <- i
	}
	defer func() {
		close(goroutineChannel)
	}()
	for {
		select {
		case <-c.Done():
			wg.Wait()
			return err
		default:
			if over || err != nil {
				wg.Wait()
				return err
			}
			var src any
			select {
			//被close的时候d可能是nil
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
			g := <-goroutineChannel
			wg.Add(1)
			run.Go(func() {
				defer func() {
					goroutineChannel <- g
					wg.Done()
				}()
				logger.Infof("当前协程编号: %d", g)
				//这里需要用一个新的接受错误，否则可能会导致老的已经发生的错误被新nil覆盖，从而获取不到发生的错误
				if src != nil { ///可能会是close事件
					err = errs.Cover(err, s.saver(args, src))
					if err != nil {
						logger.Errorf("同步数据错误%s", err.Error())
					}
				}
			})
		}
	}
}
