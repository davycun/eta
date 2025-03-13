package http_tes

import (
	"errors"
	"fmt"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/module/user"
	jsoniter "github.com/json-iterator/go"
	"net/http"
)

// func Before() {
func init() {
	logger.Infof("http_test 初始化。。。。。")
	err := initServer()
	if err != nil {
		panic(err)
	}
	initLogin()
}

func initLogin() {
	var err error
	AppId, UserId, LoginToken, err = Login(user.RootUserAccount, user.RootUserPassword)
	if err != nil {
		panic(err)
	}
}

func Login(account, password string) (appId, userId, token string, err error) {
	_, _, w := PerformRequest("POST",
		"/oauth2/login", map[string]string{"Content-Type": "application/json"},
		fmt.Sprintf(`{"username":"%s","password":"%s"}`, account, password))

	type LoginResult struct {
		Code    string `json:"code,omitempty"`
		Success bool   `json:"success,omitempty"`
		Message string `json:"message,omitempty"`
		Result  struct {
			Authorization string `json:"authorization"`
			Data          any    `json:"data"`
		} `json:"result,omitempty"`
	}

	rs := LoginResult{}

	err = jsoniter.Unmarshal(w.Body.Bytes(), &rs)
	if err != nil {
		return
	}
	if !rs.Success {
		err = errs.NewClientError(rs.Message)
		return
	}
	//LoginToken = rs.Result.Authorization
	token = rs.Result.Authorization
	d := rs.Result.Data.(map[string]interface{})
	if _, ok := d["app"]; !ok {
		err = errors.New("app 不存在或存在多个")
		return
	}
	appId = d["app"].(map[string]interface{})["id"].(string)
	userId = d["user"].(map[string]interface{})["id"].(string)
	logger.Infof("登录返回的token是：%s\n", rs.Result.Authorization)
	return
}

func migrateApp() (err error) {
	PerformRequest(http.MethodPost, "/app/migrate", map[string]string{"Content-Type": "application/json", "Authorization": LoginToken}, `{}`)
	return
}
