package snow

import (
	"github.com/davycun/eta/pkg/common/id"
	"github.com/davycun/eta/pkg/common/logger"
	"math/rand"
	"sync"
	"time"
)

const (
	SequenceBits = int64(12)
	SequenceMax  = 1<<SequenceBits - 1 //4095

	NodeBits  = int64(10)
	NodeMax   = 1<<NodeBits - 1 //1024
	NodeShift = SequenceBits

	TimeBits  = int64(41)
	TimeShift = NodeBits + SequenceBits

	DefaultEpoch = "2023-01-01"
)

var (
	defaultSnowFlake = NewSnowflake(rand.Int63n(1000)+1, DefaultEpoch) //采用随机数，减少并发测试中id重复的可能性
)

type Snowflake struct {
	seq   int64
	node  int64
	epoch time.Time
	mu    sync.Mutex
	time  int64
}

func (s *Snowflake) Generate() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Since(s.epoch).Milliseconds()

	if now == s.time {
		s.seq = (s.seq + 1) & SequenceMax
		if s.seq == 0 {
			for now <= s.time {
				now = time.Since(s.epoch).Milliseconds()
			}
		}
	} else {
		s.seq = 0
	}
	s.time = now

	return now<<TimeShift | s.node | s.seq
}

// NewSnowflake epoch is the start time layout is YYYY-MM-DD
func NewSnowflake(nodeId int64, epoch string) id.Generator {

	if nodeId < 1 {
		nodeId = rand.Int63n(1000) + 1
	}

	if nodeId > NodeMax {
		logger.Errorf("the nodeId must less than %d", NodeMax)
		nodeId = nodeId & NodeMax
	}

	if epoch == "" || len(epoch) != 10 {
		epoch = DefaultEpoch
	}
	parse, err := time.Parse("2006-01-02", epoch)
	if err != nil {
		//1672531200000 is the timestamp of 2023-01-01
		parse = time.UnixMilli(1672531200000)
	}
	sf := Snowflake{}
	sf.node = nodeId << NodeShift
	sf.epoch = parse
	return &sf
}

func DefaultSnowflake() id.Generator {
	return defaultSnowFlake
}
