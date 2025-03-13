package global

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/id/snow"
)

func GenerateID() int64 {
	if globalApp == nil || globalApp.idGenerator == nil {
		return snow.DefaultSnowflake().Generate()
	}
	return globalApp.GetIdGenerator().Generate()
}
func GenerateIDStr() string {
	id := GenerateID()
	return fmt.Sprintf("%d", id)
}
