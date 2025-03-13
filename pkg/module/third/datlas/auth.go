package datlas

import (
	"github.com/davycun/eta/pkg/common/logger"
	"time"
)

type Auth struct {
	Datlas
}

type LoginByNameParam struct {
	Name     string `json:"name"`
	Password string `json:"password"` // md5("xxxx".encode('utf-8')).hexdigest()
}
type LoginByNameResp struct {
	Rc     int `json:"rc"`
	Result struct {
		Auth       string      `json:"auth"`
		MdtUser    string      `json:"mdt_user"`
		Product    interface{} `json:"product"`
		MdtProduct []int       `json:"mdt_product"`
		Customer   struct {
			Id             int         `json:"id"`
			Uuid           string      `json:"uuid"`
			AppId          int         `json:"app_id"`
			Name           string      `json:"name"`
			Nickname       interface{} `json:"nickname"`
			Dtnickname     interface{} `json:"dtnickname"`
			Wechatbinded   bool        `json:"wechatbinded"`
			Dingtalkbinded bool        `json:"dingtalkbinded"`
			Profile        interface{} `json:"profile"`
			Phone          string      `json:"phone"`
			AreaCode       int         `json:"area_code"`
			Email          string      `json:"email"`
			EmailConfirmed bool        `json:"email_confirmed"`
			PhoneConfirmed bool        `json:"phone_confirmed"`
			VerifyCode     string      `json:"verify_code"`
			DtverifyCode   string      `json:"dtverify_code"`
			Expired        bool        `json:"expired"`
			NoExpireTime   bool        `json:"no_expire_time"`
			ExpireTime     time.Time   `json:"expire_time"`
			Logo           interface{} `json:"logo"`
			AdminMenu      struct {
			} `json:"admin_menu"`
			DetailMenu struct {
			} `json:"detail_menu"`
			Permission         interface{}   `json:"permission"`
			ManagedApp         []interface{} `json:"managed_app"`
			Roles              []int         `json:"roles"`
			PasswordUpdateTime time.Time     `json:"password_update_time"`
			Role               []string      `json:"role"`
		} `json:"customer"`
		Apps struct {
			All []struct {
				Id          int         `json:"id"`
				Name        string      `json:"name"`
				Maintainer  interface{} `json:"maintainer"`
				ExpireTime  time.Time   `json:"expire_time"`
				Description interface{} `json:"description"`
				Pid         interface{} `json:"pid"`
				Logo        interface{} `json:"logo"`
				DetailMenu  interface{} `json:"detail_menu"`
				Permission  interface{} `json:"permission"`
			} `json:"all"`
			Default struct {
				Id          int         `json:"id"`
				Name        string      `json:"name"`
				Maintainer  interface{} `json:"maintainer"`
				ExpireTime  time.Time   `json:"expire_time"`
				Description interface{} `json:"description"`
				Pid         interface{} `json:"pid"`
				Logo        interface{} `json:"logo"`
				DetailMenu  interface{} `json:"detail_menu"`
				Permission  interface{} `json:"permission"`
			} `json:"default"`
		} `json:"apps"`
		Warnings struct {
			WillExpireIn1Month bool `json:"will_expire_in_1_month"`
		} `json:"warnings"`
	} `json:"result"`
}

type CreateDatlasUserParam struct {
	Email    string                 `json:"email,omitempty" binding:"required"`
	Name     string                 `json:"name,omitempty" binding:"required"`
	Password string                 `json:"password,omitempty"` //md5以后的密码
	Admin    []CreateUserParamAdmin `json:"admin,omitempty"`
	Phone    string                 `json:"phone,omitempty" binding:"required"`
}
type CreateUserParamAdmin struct {
	AppId      int      `json:"app_id,omitempty"`
	Enable     bool     `json:"enable,omitempty"`
	ExpireTime int64    `json:"expire_time,omitempty"`
	Role       []string `json:"role,omitempty"`
}

type RegisterResult struct {
	UserId   string `json:"user_id"`
	UserUuid string `json:"user_uuid"`
}

type CreateDatlasUserResp struct {
	Rc     int `json:"rc"`
	Result struct {
		UserId   int    `json:"user_id"`
		UserUuid string `json:"user_uuid"`
	} `json:"result"`
}

type GetDatlasUserResp struct {
	Rc     int `json:"rc"`
	Result struct {
		Auth       string        `json:"auth"`
		MdtUser    string        `json:"mdt_user"`
		Product    interface{}   `json:"product"`
		MdtProduct []interface{} `json:"mdt_product"`
		Customer   struct {
			Id             int       `json:"id"`
			Uuid           string    `json:"uuid"`
			AppId          int       `json:"app_id"`
			Name           string    `json:"name"`
			Nickname       string    `json:"nickname"`
			Dtnickname     string    `json:"dtnickname"`
			Wechatbinded   bool      `json:"wechatbinded"`
			Dingtalkbinded bool      `json:"dingtalkbinded"`
			Profile        string    `json:"profile"`
			Phone          string    `json:"phone"`
			AreaCode       int       `json:"area_code"`
			Email          string    `json:"email"`
			EmailConfirmed bool      `json:"email_confirmed"`
			PhoneConfirmed bool      `json:"phone_confirmed"`
			VerifyCode     string    `json:"verify_code"`
			DtverifyCode   string    `json:"dtverify_code"`
			Expired        bool      `json:"expired"`
			NoExpireTime   bool      `json:"no_expire_time"`
			ExpireTime     time.Time `json:"expire_time"`
			Logo           string    `json:"logo"`
			AdminMenu      struct {
			} `json:"admin_menu"`
			DetailMenu struct {
			} `json:"detail_menu"`
			ManagedApp         []interface{} `json:"managed_app"`
			Roles              []int         `json:"roles"`
			PasswordUpdateTime time.Time     `json:"password_update_time"`
			Role               []string      `json:"role"`
		} `json:"customer"`
		Warnings struct {
			WillExpireIn1Month bool `json:"will_expire_in_1_month"`
		} `json:"warnings"`
	} `json:"result"`
}

func (a *Auth) LoginByName(param *LoginByNameParam) *LoginByNameResp {
	resp, err := a.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader(RequestDatlasHeader, "1").
		SetBody(param).
		SetError(&LoginByNameResp{}).
		SetResult(&LoginByNameResp{}).
		Post("/auth/login/byname")

	if err != nil {
		logger.Errorf("LoginByName resp error: %v", err)
		return nil
	}
	logger.Debugf("LoginByName resp: %s", resp)
	if resp.IsError() {
		return resp.Error().(*LoginByNameResp)
	}
	return resp.Result().(*LoginByNameResp)
}

func (a *Auth) GetUser(datlasToken string) *GetDatlasUserResp {
	// 如果指定token，则直接使用
	if datlasToken == "" {
		datlasToken = a.GetToken()
	}

	resp, err := a.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", datlasToken).
		SetHeader(RequestDatlasHeader, "1").
		SetError(&GetDatlasUserResp{}).
		SetResult(&GetDatlasUserResp{}).
		Get("/auth/v2/user/bytoken")

	if err != nil {
		logger.Errorf("GetUser resp error: %v", err)
		return nil
	}
	logger.Debugf("GetUser resp: %s", resp)

	if resp.IsError() {
		return resp.Error().(*GetDatlasUserResp)
	}
	return resp.Result().(*GetDatlasUserResp)
}

func (a *Auth) CreateUser(param *CreateDatlasUserParam) *CreateDatlasUserResp {
	res, err := a.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", a.GetToken()).
		SetHeader(RequestDatlasHeader, "1").
		SetBody(param).
		SetResult(&CreateDatlasUserResp{}).
		SetError(&CreateDatlasUserResp{}).
		Post("/auth/v2/user")

	if err != nil {
		logger.Errorf("createDatlasUser resp error: %v", err)
		return nil
	}
	logger.Debugf("createDatlasUser resp: %s", res)
	if res.IsError() {
		return res.Error().(*CreateDatlasUserResp)
	}
	return res.Result().(*CreateDatlasUserResp)
}
