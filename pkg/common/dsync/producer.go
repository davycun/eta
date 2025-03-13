package dsync

import (
	"context"
	"github.com/davycun/eta/pkg/common/logger"
)

const (
	overFlag = "over"
)

type defaultProducer struct {
	loader DataLoader
}

func NewDefaultProducer(loader DataLoader) Producer {

	pd := &defaultProducer{
		loader: loader,
	}
	return pd
}

func (s *defaultProducer) Produce(c context.Context, args any, cn chan<- any) error {
	defer func() {
		cn <- overFlag
	}()
	for {
		select {
		case <-c.Done():
			logger.Infof("producer 被强制退出了")
			return nil
		default:
			//查询中必须要满足id 升序
			data, over, err := s.loader(args)
			if err != nil {
				return err
			}
			if data != nil {
				cn <- data
			}
			logger.Infof("当前管道大小:%d", len(cn))
			if over {
				return nil
			}
		}
	}
}
