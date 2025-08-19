package dsync

import (
	"context"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/run"
	"github.com/davycun/eta/pkg/core/dto"
	"sync"
	"time"
)

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

func GetSyncOption(args *dto.Param) *SyncOption {
	qe := &SyncOption{
		ConsumerGoSize: 10,
		ChanSize:       10,
	}
	if args.Extra == nil {
		args.Extra = qe
	} else {
		if x, ok := args.Extra.(*SyncOption); ok {
			if x.ConsumerGoSize < 1 {
				x.ConsumerGoSize = 10
			}
			if x.ChanSize < 1 {
				x.ChanSize = 10
			}
			return x
		}
	}
	return args.Extra.(*SyncOption)
}
