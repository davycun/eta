package broker

import (
	"context"
	"github.com/davycun/eta/pkg/common/id/nanoid"
	"github.com/davycun/eta/pkg/common/logger"
	jsoniter "github.com/json-iterator/go"
	"sync"
)

var (
	defaultBroker = NewBroker()
)

const (
	EventOptTypeInsert = "insert"
	EventOptTypeUpdate = "update"
	EventOptTypeDelete = "delete"
)

type Closer interface {
	Close(c context.Context) error
}

type Event struct {
	Id        string         `json:"id"`         //事件ID
	UserId    string         `json:"user_id"`    //用户的ID
	OptType   string         `json:"opt_type"`   //insert、update、delete
	Data      any            `json:"data"`       //实际时间的数据内容
	TableName string         `json:"table_name"` // 操作的表的ID
	Extra     map[string]any `json:"extra"`      //存储额外的值
}

func (e Event) String() string {
	s, _ := jsoniter.MarshalToString(e)
	return s
}

type EventFunc func(*Event)

func NewEvent(userId string, data any, ef ...EventFunc) *Event {
	e := &Event{
		Id:     nanoid.New(),
		UserId: userId,
		Data:   data,
	}
	for _, f := range ef {
		f(e)
	}
	return e
}

// Receiver 如果Receiver一直不返回，那么会影响Broker
// 函数允许返回error 目的是为了后续Broker对Receiver进行Retry
type Receiver func(event *Event) error

type Subscriber interface {
	Subscribe(rec Receiver) error
}

type Publisher interface {
	Publish(c context.Context, event *Event, autoCommit bool) error
	Commit(c context.Context, eventId ...string) error
	Rollback(c context.Context, eventId ...string) error
}

type Broker struct {
	BrokerOption
	cnMap *channelMap
	mutex *sync.Mutex
}

// BrokerOption
// Broker的选项配置
type BrokerOption struct {
}

type BrokeOptionFunc func(*BrokerOption)

func NewBroker(ofs ...BrokeOptionFunc) *Broker {
	bk := &Broker{}
	for _, f := range ofs {
		f(&bk.BrokerOption)
	}
	bk.cnMap = &channelMap{}
	bk.mutex = &sync.Mutex{}
	return bk
}

func (b *Broker) Load(name string) Channel {
	return b.cnMap.Load(name)
}
func (b *Broker) LoadAll() map[string]Channel {
	return b.cnMap.LoadAll()
}
func (b *Broker) LoadOrNew(name string, option ...ChannelOptionFunc) Channel {
	cn := b.cnMap.Load(name)
	if cn != nil {
		return cn
	}
	b.mutex.Lock()
	defer b.mutex.Unlock()

	cn = b.cnMap.Load(name)
	if cn != nil {
		return cn
	}
	cn = NewChannel(name, option...)
	b.cnMap.Store(name, cn)
	return cn
}

func Close(c context.Context, name string) {
	cn := defaultBroker.Load(name)
	if cn != nil {
		cn.Close(c)
	}
}
func CloseAll(c context.Context) {
	for _, v := range defaultBroker.LoadAll() {
		v.Close(c)
	}
}

func LoadChannel(name string, option ...ChannelOptionFunc) Channel {
	return defaultBroker.LoadOrNew(name, option...)
}

// Publish
// name 是channel的名称
func Publish(c context.Context, name string, event *Event, autoCommit bool, opts ...ChannelOptionFunc) error {
	return defaultBroker.LoadOrNew(name, opts...).Publish(c, event, autoCommit)
}
func Subscribe(name string, rec Receiver, opts ...ChannelOptionFunc) error {
	return defaultBroker.LoadOrNew(name, opts...).Subscribe(rec)
}
func Commit(c context.Context, name string, eventId ...string) error {
	cn := defaultBroker.Load(name)
	if cn == nil {
		logger.Infof("[%s] not found channel for commit event", name)
		return nil
	}
	return cn.Commit(c, eventId...)
}
func Rollback(c context.Context, name string, eventId ...string) error {
	cn := defaultBroker.Load(name)
	if cn == nil {
		logger.Infof("[%s] not found channel for rollback event", name)
		return nil
	}

	return cn.Rollback(c, eventId...)
}
