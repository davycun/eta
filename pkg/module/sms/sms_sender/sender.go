package sms_sender

import (
	"errors"
	"fmt"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/module/setting"
	"github.com/davycun/eta/pkg/module/sms"
	"gorm.io/gorm"
	"maps"
)

var (
	senderMap = make(map[string]sms.Sender)
)

func Registry(key string, sd sms.Sender) error {
	if _, ok := senderMap[key]; ok {
		return errors.New(fmt.Sprintf("短信供应商[%s]已经注册", key))
	}
	if sd == nil {
		return errors.New(fmt.Sprintf("短信供应商[%s]的实现为空", key))
	}
	senderMap[key] = sd
	return nil
}

// SendSms
// 发送短信
func SendSms(txDb *gorm.DB, taskList ...sms.Task) ([]sms.Target, error) {

	var (
		targetResult = make([]sms.Target, 0, len(taskList)*5)
	)

	sd, err := NewSender(txDb)
	if err != nil {
		return targetResult, err
	}

	for _, v := range taskList {
		tgList, err1 := sd.SendCustom(v)
		if err1 != nil {
			return targetResult, err1
		}
		for i, _ := range tgList {
			tgList[i].TaskId = v.ID
		}
		targetResult = append(targetResult, tgList...)
	}
	if len(targetResult) < 1 {
		return targetResult, errs.NewClientError("没有发送成功")
	}
	return targetResult, nil
}

func NewSender(txDb *gorm.DB) (sms.Sender, error) {
	cfg, b := setting.GetSmsConfig(txDb)
	if !b && len(senderMap) != 1 {
		return nil, errs.NewServerError("没有配置短信供应商!")
	}
	s, ok := senderMap[cfg.Vendor]
	if !ok {
		//虽然没有找到配置的供应商，但是实际只有一个供应商，那么就取第一个
		if len(senderMap) == 1 {
			for v := range maps.Values(senderMap) {
				s = v
				break
			}
		} else {
			return nil, errs.NewServerError(fmt.Sprintf("没有找到短信供应商[%s]的发送器", cfg.Vendor))
		}
	}
	return s, nil
}
