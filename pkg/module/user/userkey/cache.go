package userkey

import "github.com/davycun/eta/pkg/common/global"

func LoadByKey(key string) (UserKey, error) {
	var (
		db     = global.GetLocalGorm()
		ukList []UserKey
	)
	err := db.Model(&ukList).Where(map[string]any{"access_key": key}).Limit(1).Find(&ukList).Error

	if err != nil || len(ukList) < 1 {
		return UserKey{}, err
	}
	return ukList[0], err
}
func LoadByFixToken(token string) (UserKey, error) {
	var (
		db     = global.GetLocalGorm()
		ukList []UserKey
	)
	err := db.Model(&ukList).Where(map[string]any{"fixed_token": token}).Limit(1).Find(&ukList).Error

	if err != nil || len(ukList) < 1 {
		return UserKey{}, err
	}
	return ukList[0], err
}
