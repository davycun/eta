package user_srv

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/id/nanoid"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/xuri/excelize/v2"
	"time"
)

var (
	defaultSheetName = "sheet1"
)

// 文件命名
func createFileName() string {
	name := time.Now().Format("2006-01-02")
	return fmt.Sprintf("t_user/%s-%v.xlsx", nanoid.New(), name)
}

// 设置首行
func writeTop(columns []string, styleID int) []interface{} {
	var fisrtRow = make([]interface{}, 0)
	for _, conf := range columns {
		title := conf
		fisrtRow = append(fisrtRow, excelize.Cell{
			Value:   title,
			StyleID: styleID,
		})
	}
	return fisrtRow
}

func initStreamWriter() (*excelize.File, *excelize.StreamWriter, int) {
	// 初始化文件
	file := excelize.NewFile()
	// 创建一个默认的sheet
	sheetName := defaultSheetName
	// 生成一个流式处理的写入器
	streamWriter, _ := file.NewStreamWriter(sheetName)
	// 设置样式
	styleID, err := file.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		logger.Errorf("设置样式失败:%s", err) //错误处理
	}
	return file, streamWriter, styleID
}
