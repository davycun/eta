package faker

import "fmt"

func Xiaoqu() string {
	return fmt.Sprintf("%s%s",
		getRandValue([]string{"address", "xiaoqu_prefix"}), getRandValue([]string{"address", "xiaoqu_suffix"}),
	)
}

func Building() string {
	return fmt.Sprintf("%s%d%s",
		Xiaoqu(),
		Number(0, 100),
		RandString([]string{"栋", "幢", "座", "号", "号楼"}),
	)
}

func Road() string {
	return fmt.Sprintf("%s%s",
		getRandValue([]string{"address", "road_prefix"}), getRandValue([]string{"address", "road_suffix"}),
	)
}

func ShorAddress() string {
	return fmt.Sprintf("%s%d%s",
		Road(),
		Number(0, 100),
		RandString([]string{"弄", "号"}),
	)
}
