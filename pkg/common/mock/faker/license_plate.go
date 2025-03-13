package faker

import (
	"fmt"
	"strings"
)

// LicensePlate 传统车牌
func LicensePlate() string {
	return fmt.Sprintf("%s%s%s",
		getRandValue([]string{"license_plate", "provinces"}),
		strings.ToUpper(Letter()),
		multiRandValue([]string{"license_plate", "num"}, 5),
	)
}

// LicensePlateSpecial 特种车牌
func LicensePlateSpecial() string {
	return fmt.Sprintf("%s%s%s%s",
		getRandValue([]string{"license_plate", "provinces"}),
		strings.ToUpper(Letter()),
		multiRandValue([]string{"license_plate", "num"}, 4),
		getRandValue([]string{"license_plate", "last"}),
	)
}

// LicensePlateCustom 自定义车牌
func LicensePlateCustom(prov, org, last string) string {
	if last == "" {
		return fmt.Sprintf("%s%s%s", prov, org, multiRandValue([]string{"license_plate", "num"}, 5))
	} else {
		return fmt.Sprintf("%s%s%s%s", prov, org, multiRandValue([]string{"license_plate", "num"}, 4), last)
	}
}

// LicensePlateNewEnergy 新能源车牌
func LicensePlateNewEnergy(carModel int) string {
	// 大型车
	if carModel == 1 {
		return fmt.Sprintf("%s%s%s%s",
			getRandValue([]string{"license_plate", "provinces"}),
			strings.ToUpper(Letter()),
			multiRandValue([]string{"license_plate", "num"}, 5),
			getRandValue([]string{"license_plate", "new_energy"}),
		)
	}
	// 小型车
	return fmt.Sprintf("%s%s%s%s",
		getRandValue([]string{"license_plate", "provinces"}),
		strings.ToUpper(Letter()),
		getRandValue([]string{"license_plate", "new_energy"}),
		multiRandValue([]string{"license_plate", "num"}, 5),
	)
}
