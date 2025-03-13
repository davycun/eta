package sms_test

import (
	_ "github.com/davycun/eta/pkg/common/http_tes"
	"github.com/davycun/eta/pkg/module/sms"
	"github.com/davycun/eta/pkg/module/sms/sms_sender"
	"github.com/stretchr/testify/assert"
	"testing"
)

type aliSms struct {
}

func (a aliSms) SendCustom(task sms.Task) ([]sms.Target, error) {
	return []sms.Target{sms.Target{}}, nil
}
func (a aliSms) SendTemplate(target sms.Target) (sms.Target, error) {

	return target, nil
}

func TestSendVerifyCode(t *testing.T) {
	// 1. 初始化客户端

	sd := &aliSms{}
	err := sms_sender.Registry("aliTest", sd)
	assert.Nil(t, err)
	_, err = sms_sender.SendSms(nil, sms.Task{})
	assert.Nil(t, err)
	err = sms_sender.SendVerifyCode(nil, sms.Target{Mobile: "13312345678"}, "123")
	assert.Nil(t, err)
	assert.True(t, sms_sender.VerifyVerifyCode("13312345678", "123"))
}
