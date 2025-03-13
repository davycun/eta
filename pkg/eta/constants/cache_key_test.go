package constants

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/id/nanoid"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/oklog/ulid/v2"
	"math/rand"
	"testing"
	"time"
)

func TestRedisKey(t *testing.T) {
	println(RedisKey(TokenKey, "x1"))
}

func TestUlid(t *testing.T) {
	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	ms := ulid.Timestamp(time.Now())
	u, _ := ulid.New(ms, entropy)
	println(u.String())
	ul := ulid.Make().String()
	println(ul)
	parse, err := ulid.Parse(ul)
	if err != nil {
		logger.Errorf("ulid parse error: %v", err)
	}
	println(fmt.Sprintf("%v", parse.Time()))
	//println(uuid.NewString())
	//println(uuid.New().String())
	println(nanoid.New())
}
