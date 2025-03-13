package http_tes

import (
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/module/integration"
	"net/http"
	"testing"
)

type Header struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type HttpCaseOption func(hc *HttpCase)

func Query[T any](t *testing.T, uri string, args dto.RetrieveParam, hco ...HttpCaseOption) ([]T, int64) {
	var (
		rs  = make([]T, 0, 2)
		qRs = dto.Result{
			Data: &rs,
		}
		resp = dto.ControllerResponse{
			Result: &qRs,
		}
	)

	hc := HttpCase{
		Path:         uri,
		ResponseDest: &resp,
		Body: &dto.Param{
			RetrieveParam: args,
		},
	}
	for _, f := range hco {
		f(&hc)
	}

	Call(t, hc)

	return rs, qRs.Total
}

// Modify
// T类型是ModifyResult中data的类型
func Modify[T any](t *testing.T, uri string, args dto.ModifyParam, hco ...HttpCaseOption) ([]T, int64) {
	var (
		rs  = make([]T, 0, 2)
		qRs = dto.Result{
			Data: &rs,
		}
		resp = dto.ControllerResponse{
			Result: &qRs,
		}
	)

	hc := HttpCase{
		Path:         uri,
		ResponseDest: &resp,
		Body:         &args,
	}
	for _, f := range hco {
		f(&hc)
	}

	Call(t, hc)

	return rs, qRs.RowsAffected
}
func Integration(t *testing.T, uri string, args integration.CommandParam, hco ...HttpCaseOption) (rs integration.CommandResult) {
	var (
		resp = dto.ControllerResponse{
			Result: &rs,
		}
	)

	hc := HttpCase{
		Path:         uri,
		ResponseDest: &resp,
		Body:         &args,
	}
	for _, f := range hco {
		f(&hc)
	}

	Call(t, hc)

	return
}
func Delete[T any](t *testing.T, uri string, args dto.ModifyParam, hco ...HttpCaseOption) ([]T, int64) {
	var (
		rs  = make([]T, 0, 2)
		qRs = dto.Result{
			Data: &rs,
		}
		resp = dto.ControllerResponse{
			Result: &qRs,
		}
	)

	hc := HttpCase{
		Method:       http.MethodDelete,
		Path:         uri,
		ResponseDest: &resp,
		Body:         &args,
	}
	for _, f := range hco {
		f(&hc)
	}

	Call(t, hc)

	return rs, qRs.RowsAffected
}
