package oauth2

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/errs"
	"gorm.io/gorm"
)

type Service struct {
	C  *ctx.Context
	DB *gorm.DB
}

func NewService(c *ctx.Context, db *gorm.DB) *Service {
	return &Service{C: c, DB: db}
}

func (s *Service) Login(args *LoginParam, result *LoginResult) error {

	if args.LoginType == "" {
		return errs.NewClientError("login type is empty")
	}
	fc, ok := loginFuncMap[args.LoginType]
	if !ok {
		return errs.NewClientError(fmt.Sprintf("the login type of %s not supported", args.LoginType))
	}
	us, err := fc(s.C, args.Param)

	if err != nil {
		return err
	}
	if us.ID == "" {
		return errs.NewClientError("user not exists")
	}

	return ProcessResult(s.C, &us, args.LoginType, result)
}
