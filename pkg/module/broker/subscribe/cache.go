package subscribe

import (
	"github.com/davycun/eta/pkg/core/loader"
	"github.com/davycun/eta/pkg/eta/constants"
	"gorm.io/gorm"
)

var (
	//存储的内容是appId -> []Subscriber
	allData = loader.NewCacheLoader[Subscriber, Subscriber](constants.TableSubscriber, constants.CacheAllDataSubscriber).SetKeyName("id")
)

func DelCache(db *gorm.DB, dataList ...Subscriber) {
	for _, v := range dataList {
		if v.ID == "" {
			continue
		}
		allData.Delete(db, v.ID)
	}
}

func LoadAllSubscriber(db *gorm.DB) (dtMap map[string]Subscriber, err error) {
	return allData.LoadAll(db)
}

func LoadSubscriberByTarget(db *gorm.DB, target string) ([]Subscriber, error) {
	subs, err := LoadAllSubscriber(db)
	if err != nil {
		return nil, err
	}

	rs := make([]Subscriber, 0, 3)
	for _, v := range subs {
		if v.Target == target || v.Target == ("d_"+target) {
			rs = append(rs, v)
			continue
		}
	}
	return rs, nil
}
