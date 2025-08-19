package reload

import (
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/service"
)

type Service struct {
	service.DefaultService
}

func (s *Service) Db2Es(args *dto.Param, rs *dto.Result) error {

	return nil
}
