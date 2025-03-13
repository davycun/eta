package broker_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/davycun/eta/pkg/module/broker"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewChannel(t *testing.T) {

	var (
		c  = context.Background()
		c1 = context.TODO()
		cn = make(chan int, 5)
		rs = ""
	)

	select {
	case <-time.After(500 * time.Millisecond):
		rs = "时间到"
	case <-c.Done():
		rs = "context done"
	case <-cn:
		rs = "获取成功"
	}
	assert.Equal(t, "时间到", rs)

	select {
	case <-time.After(500 * time.Millisecond):
		rs = "时间到"
	case <-c1.Done():
		rs = "context done"
	case <-cn:
		rs = "获取成功"
	}

	assert.Equal(t, "时间到", rs)
}

func TestConsumer(t *testing.T) {

	x := -5
	y := 0
	z := 5

	i := 0
	for range x {
		i++
	}
	assert.Equal(t, 0, i)

	i = 0
	for range y {
		i++
	}
	assert.Equal(t, 0, i)

	i = 0
	for range z {
		i++
	}
	assert.Equal(t, 5, i)
}

func TestBroker(t *testing.T) {
	var (
		channelId = "test"
	)
	err := broker.Publish(context.Background(), channelId, broker.NewEvent("my_user_id", "test"), true)
	assert.Nil(t, err)
	err = broker.Subscribe(channelId, func(event *broker.Event) error {
		assert.Equal(t, "test", fmt.Sprintf("%s", event.Data))
		assert.Equal(t, "my_user_id", event.UserId)
		return nil
	})
	assert.Nil(t, err)

	broker.CloseAll(context.Background())
	err = broker.Publish(context.Background(), channelId, broker.NewEvent("my_user_id", "test"), true)
	assert.True(t, errors.Is(err, broker.ChannelAlreadyClosed))
}
