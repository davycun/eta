package broker

import (
	"context"
	"errors"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/utils"
	"sync"
	"time"
)

var (
	ChannelAlreadyClosed = errors.New("channel is closed")
)

type Channel interface {
	Subscriber
	Publisher
	Close(c context.Context)
	IsClosed() bool
}

type ChannelOption struct {
	name         string
	channelSize  int
	consumerSize int
}

func (c *ChannelOption) SetChannelSize(sz int) {
	c.channelSize = sz
}
func (c *ChannelOption) SetConsumerSize(sz int) {
	c.consumerSize = sz
}

type ChannelOptionFunc func(*ChannelOption)

func setChannelSize(co *ChannelOption) {
	co.SetChannelSize(20)
}
func setConsumerSize(co *ChannelOption) {
	co.SetConsumerSize(1)
}

func NewChannel(name string, opts ...ChannelOptionFunc) Channel {
	if name == "" {
		name = "default"
		logger.Warnf("you will get a default channel, because the name is empty")
	}
	cn := &defaultChannel{
		ChannelOption: ChannelOption{
			name: name,
		},
	}
	for _, opt := range opts {
		opt(&cn.ChannelOption)
	}
	if cn.channelSize < 1 {
		setChannelSize(&cn.ChannelOption)
	}
	if cn.consumerSize < 1 {
		setConsumerSize(&cn.ChannelOption)
	}
	cn.channel = make(chan *Event, cn.channelSize)
	cn.close = make(chan struct{}, 3)
	cn.start()

	return cn
}

type defaultChannel struct {
	ChannelOption
	recList []Receiver
	dirty   sync.Map //eventId -> *event
	channel chan *Event
	close   chan struct{}
	closed  bool
}

func (d *defaultChannel) Subscribe(rec Receiver) error {
	if d.IsClosed() {
		return ChannelAlreadyClosed
	}
	d.recList = append(d.recList, rec)
	return nil
}

// Publish
// 默认两秒钟后超时
func (d *defaultChannel) Publish(c context.Context, event *Event, autoCommit bool) error {

	if autoCommit {
		return d.publish(c, event)
	} else {
		logger.Infof("[%s] publish event waiting commit or rollback: %s ", utils.FmtTextBlue(d.name), event.String())
		d.dirty.Store(event.Id, event)
	}
	return nil
}

func (d *defaultChannel) Commit(c context.Context, eventIds ...string) error {

	var (
		err error
	)
	for _, v := range eventIds {
		value, loaded := d.dirty.LoadAndDelete(v)
		if !loaded {
			logger.Infof("[%s] not found event[%s] for commit ", d.name, v)
			return nil
		}
		event := value.(*Event)
		err = d.publish(c, event)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *defaultChannel) publish(c context.Context, event *Event) error {
	if event == nil {
		return nil
	}
	if event.Id == "" {
		return errors.New("event id can not be empty")
	}
	if d.IsClosed() {
		return ChannelAlreadyClosed
	}
	if c == nil {
		c = context.Background()
	}

	select {
	case <-time.After(2 * time.Second):
		logger.Errorf("[%s] the channel is full, can not publish event: %s", utils.FmtTextBlue(d.name), event.String())
	case <-c.Done():
		logger.Errorf("[%s] context had done, can not publish event: %s", utils.FmtTextBlue(d.name), event.String())
		return c.Err()
	case d.channel <- event:
		logger.Infof("[%s] commit publish event: %s ", utils.FmtTextBlue(d.name), event.String())
	case <-d.close:
		logger.Infof("[%s] the channel has closed ,can not publish event: %s", utils.FmtTextBlue(d.name), event.String())
	}
	return nil
}

func (d *defaultChannel) Rollback(c context.Context, eventIds ...string) error {
	for _, v := range eventIds {
		logger.Infof("[%s] rollback event %s", utils.FmtTextBlue(d.name), v)
		d.dirty.Delete(v)
	}
	return nil
}

func (d *defaultChannel) Close(c context.Context) {
	if d.IsClosed() {
		return
	}
	if c == nil {
		c = context.Background()
	}
	select {
	case <-c.Done():
		logger.Infof("[%s] channel has closed err %s", utils.FmtTextBlue(d.name), c.Err())
	case d.close <- struct{}{}:
		d.closed = true
		close(d.channel)
		close(d.close)
		logger.Infof("[%s] The Channel is closed using the Close method", utils.FmtTextBlue(d.name))
	}
}
func (d *defaultChannel) IsClosed() bool {
	return d.closed
}
func (d *defaultChannel) start() {
	logger.Infof("[%s] start channel, consumerSize is %d", utils.FmtTextBlue(d.name), d.consumerSize)
	for range d.consumerSize {
		go func() {
			for {
				select {
				case <-d.close:
					logger.Infof("[%s] receive channel closed event, we will exit", utils.FmtTextBlue(d.name))
					return
				case event := <-d.channel:
					if event == nil {
						//这里直接return也行，毕竟在publish event 的时候是不允许传入nil进来的，所以为nil的情况只能是channel被close了
						logger.Infof("[%s] channel maybe closed,because event is nil", utils.FmtTextBlue(d.name))
						continue
					}
					logger.Infof("[%s] start process event: %s", utils.FmtTextBlue(d.name), event.String())
					for _, rec := range d.recList {
						if err := rec(event); err != nil {
							//TODO 这里应该支持重试
							logger.Infof("receive err %s ", err)
						}
					}
				}
			}
		}()
	}
	return
}

type channelMap struct {
	channels sync.Map
}

func (c *channelMap) Load(name string) Channel {
	value, ok := c.channels.Load(name)
	if ok {
		return value.(Channel)
	}
	return nil
}
func (c *channelMap) LoadAll() map[string]Channel {

	mp := make(map[string]Channel)
	c.channels.Range(func(key, value any) bool {
		k := key.(string)
		if val, ok := value.(Channel); ok {
			mp[k] = val
		}
		return true
	})
	return mp
}

func (c *channelMap) Store(name string, channel Channel) {
	c.channels.Store(name, channel)
}
