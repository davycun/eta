package errs

var (
	NotFound          = NewClientError("Not Found")
	NoRecordAffected  = NewClientError("没有匹配到需要操作的数据")
	NoPermission      = NewAuthErrorCode("403", "No Permission")
	NoPermissionNoErr = NewAuthErrorCode("200", "无权限获取相关数据") //当数据权限校验是遇到这个错误，还是返回给前端200，只是没有数据而已
)

type BaseError struct {
	Code    string
	Message string
	Cause   error
}

func (c *BaseError) Error() string {
	return c.Message
}

type ClientError struct {
	BaseError
}

func NewClientError(msg string) error {
	return &ClientError{BaseError{Message: msg}}
}
func NewClientErrorCode(code, msg string) error {
	return &ClientError{BaseError{Code: code, Message: msg}}
}
func NewClientErrorCause(msg string, cause error) error {
	return &ClientError{BaseError{Message: msg, Cause: cause}}
}

type ServerError struct {
	BaseError
}

func NewServerError(msg string) error {
	return &ServerError{BaseError{Message: msg}}
}
func NewServerErrorCode(code, msg string) error {
	return &ServerError{BaseError{Code: code, Message: msg}}
}
func NewServerErrorCause(code, msg string, cause error) error {
	return &ServerError{BaseError{Message: msg, Cause: cause}}
}

type AuthError struct {
	BaseError
}

func NewAuthError(msg string) error {
	return &AuthError{BaseError{Message: msg}}
}
func NewAuthErrorCode(code, msg string) error {
	return &AuthError{BaseError{Message: msg}}
}
func NewAuthErrorCause(code, msg string, cause error) error {
	return &AuthError{BaseError{Message: msg}}
}
