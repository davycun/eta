package dsync

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type syncArgs struct {
	size int
}

func TestDsync(t *testing.T) {
	var (
		args = &syncArgs{}
	)
	sc := NewDefaultSyncService(MyDataLoader, MyDataSaver, SyncOption{
		ConsumerGoSize: 2,
		ChanSize:       100,
	})
	err := sc.Sync(args, args)
	assert.Nil(t, err)
}

func MyDataLoader(args any) (data any, over bool, err error) {
	var (
		sa = args.(*syncArgs)
	)
	sa.size++
	return sa.size, sa.size > 30, nil
}

func MyDataSaver(args any, data any) error {
	var (
		dt = data.(int)
	)
	if dt == 10 {
		time.Sleep(30 * time.Second)
	} else {
		time.Sleep(1 * time.Second)
	}
	return nil
}
