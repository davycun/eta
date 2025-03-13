package service

import (
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/eta/constants"
)

const (
	MsgKey = constants.WsKeyDataExport
)

// Export query 接口导出
func (s *DefaultService) Export(param *dto.Param, rs *dto.Result) error {

	return nil
}
