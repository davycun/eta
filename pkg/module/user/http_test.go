package user_test

import (
	"github.com/davycun/eta/pkg/common/http_tes"
	"github.com/davycun/eta/pkg/module/user"
	"testing"
)

func TestCreate(t *testing.T) {

	us := user.NewTestData()

	ls := http_tes.Create(t, "/user/create", []user.User{us})

}
