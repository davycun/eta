package oauth2

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/user"
	"gorm.io/gorm"
)

var (
	loginFuncMap = map[string]LoginFunc{
		constants.LoginTypeAccount:       LoginByAccount,
		constants.LoginTypeDingService:   LoginByDingCode,
		constants.LoginTypeDingQrcode:    LoginByDingCode,
		constants.LoginTypeZzdService:    LoginByZzDingCode,
		constants.LoginTypeZzdQrcode:     LoginByZzDingQrCode,
		constants.LoginTypeWechatService: LoginByWechatCode,
		constants.LoginTypeWechatQrcode:  LoginByWechatCode,
		constants.LoginTypeWeComService:  LoginByWeComCode,
		constants.LoginTypeWeComQrcode:   LoginByWeComCode,
		constants.LoginTypeSmsService:    LoginBySmsCode,
		constants.LoginTypeAccessToken:   LoginByAccessToken,
	}
)

type LoginFunc func(c *ctx.Context, args any) (user.User, error)

func RegistryLoginFunc(loginType string, fc LoginFunc) {
	if _, ok := loginFuncMap[loginType]; ok {
		logger.Errorf("The login method %s already exists and will be overwritten", loginType)
	}
	loginFuncMap[loginType] = fc
}

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
