package errs

import (
	"github.com/davycun/eta/pkg/common/logger"
	"net/http"
)

func Cover(src error, target error) error {
	if src != nil {
		return src
	}
	if target != nil {
		return target
	}
	return src
}

func HttpStatus(err error) (httpStatus int, baseError BaseError) {

	var e interface{}
	e = err

	switch cErr := e.(type) {
	case *ClientError:
		return http.StatusBadRequest, cErr.BaseError
	case *ServerError:
		return http.StatusInternalServerError, cErr.BaseError
	case *AuthError:
		return http.StatusForbidden, cErr.BaseError
	case *BaseError:
		return http.StatusBadRequest, *cErr
	}

	if err != nil {
		return http.StatusInternalServerError, BaseError{Code: "500", Message: err.Error()}
	}
	return http.StatusOK, BaseError{}
}

// TryUnwrap if err is nil then it returns a valid value
// If err is not nil, Unwrap panics with err.
// Play: https://go.dev/play/p/w84d7Mb3Afk
func TryUnwrap[T any](val T, err error) T {
	if err != nil {
		logger.Errorf("error occured:%s", err)
	}
	return val
}
