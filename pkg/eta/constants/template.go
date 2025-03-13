package constants

import "github.com/davycun/eta/pkg/common/utils"

const (
	TemplateColumnPrefix   = "F"
	TemplateTableNameField = "TableName"
	TemplateRaDbFields     = "RaDbFields"
)

func Column2StructFieldName(name string) string {
	return TemplateColumnPrefix + utils.UnderlineToHump(name)
}
