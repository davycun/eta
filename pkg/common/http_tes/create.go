package http_tes

import (
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Create[T any](t *testing.T, uri string, data []T) []entity.BaseEntity {

	var (
		rs  = make([]entity.BaseEntity, 0, 2)
		mRs = dto.Param{
			ModifyParam: dto.ModifyParam{
				Data: &rs,
			},
		}
		resp = dto.ControllerResponse{
			Result: &mRs,
		}
	)
	Call(t, HttpCase{
		Method:       "POST",
		Path:         uri,
		ResponseDest: &resp,
		Body: dto.ModifyParam{
			Data: data,
		},
	})

	assert.Equal(t, "200", resp.Code)
	return rs
}
