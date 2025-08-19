package utils

import "github.com/davycun/eta/pkg/eta/constants"

func Column2StructFieldName(name string) string {
	return constants.TemplateColumnPrefix + UnderlineToHump(name)
}
