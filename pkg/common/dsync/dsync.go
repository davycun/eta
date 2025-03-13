package dsync

import (
	"context"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/run"
	"sync"
	"time"
)

type SyncOption struct {
	ConsumerGoSize int  `json:"consumer_go_size,omitempty"`
	ChanSize       int  `json:"chan_size,omitempty"`
	Restore        bool `json:"clean,omitempty"` //是否对相关数据先清空
	Merge          bool `json:"merge,omitempty"` //同步的时候是否采用Merge的方式，比如同步宽表，如果以前表已经有数据就传入true
}

type defaultSyncService struct {
	Producer
	Consumer
	chanSize int
}

func NewDefaultSyncService(loader DataLoader, saver DataSaver, option SyncOption) SyncService {

	ds := &defaultSyncService{
		Producer: NewDefaultProducer(loader),
		Consumer: NewDefaultConsumer(saver, option.ConsumerGoSize),
		chanSize: option.ChanSize,
	}
	if ds.chanSize < 1 {
		ds.chanSize = 1
	}
	return ds
}

func (s *defaultSyncService) Sync(producerArgs, consumerArgs any) error {

	var (
		cn            = make(chan any, s.chanSize)
		wg            = &sync.WaitGroup{}
		c, cancelFunc = context.WithCancel(context.Background())
		produceErr    error
		consumeErr    error
	)
	defer func() {
		close(cn)
		cancelFunc()
	}()

	wg.Add(2)
	run.Go(func() {
		produceErr = s.Produce(c, producerArgs, cn)
		if produceErr != nil {
			cancelFunc()
			logger.Errorf("produce err %s", produceErr.Error())
		}
		wg.Done()
	})
	run.Go(func() {
		consumeErr = s.Consume(c, consumerArgs, cn)
		if consumeErr != nil {
			cancelFunc()
			cleanChan(cn)
			logger.Errorf("consumer err %s", consumeErr.Error())
		}
		wg.Done()
	})
	wg.Wait()
	if produceErr != nil {
		return produceErr
	}
	if consumeErr != nil {
		return consumeErr
	}
	return nil
}

func cleanChan(ch <-chan any) {

	t := time.NewTimer(time.Second * 10)
	for {
		select {
		case <-ch:
		case <-t.C:
			return
		}
	}
}
