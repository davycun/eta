package ding

import (
	"github.com/davycun/eta/pkg/common/logger"
)

type User struct {
	Ding
}

type GetUserResp struct {
	Nick      string `json:"nick"`
	AvatarUrl string `json:"avatarUrl"`
	Mobile    string `json:"mobile"`
	OpenId    string `json:"openId"`
	UnionId   string `json:"unionId"`
	Email     string `json:"email"`
	StateCode string `json:"stateCode"`
}

/*
GetUser 获取用户通讯录个人信息

https://open.dingtalk.com/document/orgapp/dingtalk-retrieve-user-information
*/
func (o *User) GetUser(unionId, userAccessToken string) *GetUserResp {
	res := &GetUserResp{}
	if o.Err != nil {
		return res
	}
	path := "/v1.0/contact/users/{unionId}"
	resp, err := o.client.R().
		SetPathParam("unionId", unionId).
		SetHeader("Content-Type", "application/json").
		SetHeader("x-acs-dingtalk-access-token", userAccessToken).
		SetError(&GetUserResp{}).
		SetResult(&GetUserResp{}).
		Get(path)

	if err != nil {
		o.Err = err
		return res
	}
	logger.Debugf("Ding User.GetUser resp: %s", resp)
	// {"nick":"zhangsan","avatarUrl":"https://xxx","mobile":"150xxxx9144","openId":"123","unionId":"z21HjQliSzpw0Yxxxx","email":"zhangsan@alibaba-inc.com","stateCode":"86"}
	if resp.IsError() {
		return resp.Error().(*GetUserResp)
	}
	return resp.Result().(*GetUserResp)
}
