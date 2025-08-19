package dict_srv

import (
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/dict"
)

var (

	//组织类型
	labelCategoriesName       = "标签大类"
	labelCategoriesDictionary = []dict.Dictionary{
		{BaseEntity: entity.BaseEntity{ID: "167556962316193792", Remark: "default"}, Namespace: constants.NamespaceEta, Category: ctype.NewStringPrt(labelCategoriesName), Name: ctype.NewStringPrt("预设标签"), Order: 1},
		{BaseEntity: entity.BaseEntity{ID: "167556962316193793", Remark: "custom"}, Namespace: constants.NamespaceEta, Category: ctype.NewStringPrt(labelCategoriesName), Name: ctype.NewStringPrt("自定义标签"), Order: 2},
	}
	labelColorName       = "标签颜色"
	labelColorDictionary = []dict.Dictionary{
		{BaseEntity: entity.BaseEntity{ID: "167556962316193794", Remark: "blue"}, Namespace: constants.NamespaceEta, Category: ctype.NewStringPrt(labelColorName), Name: ctype.NewStringPrt("蓝色"), Order: 1},
		{BaseEntity: entity.BaseEntity{ID: "167556962316193795", Remark: "#ECEFFA"}, Namespace: constants.NamespaceEta, Category: ctype.NewStringPrt(labelColorName), Name: ctype.NewStringPrt("background"), Order: 1, ParentId: "167556962316193794"},
		{BaseEntity: entity.BaseEntity{ID: "167556962316193796", Remark: "#3D7FE9"}, Namespace: constants.NamespaceEta, Category: ctype.NewStringPrt(labelColorName), Name: ctype.NewStringPrt("color"), Order: 2, ParentId: "167556962316193794"},
		{BaseEntity: entity.BaseEntity{ID: "167556962316193797", Remark: "gray"}, Namespace: constants.NamespaceEta, Category: ctype.NewStringPrt(labelColorName), Name: ctype.NewStringPrt("灰色"), Order: 2},
		{BaseEntity: entity.BaseEntity{ID: "167556962316193798", Remark: "#EBECEF"}, Namespace: constants.NamespaceEta, Category: ctype.NewStringPrt(labelColorName), Name: ctype.NewStringPrt("background"), Order: 1, ParentId: "167556962316193797"},
		{BaseEntity: entity.BaseEntity{ID: "167556962316193799", Remark: "#6E798C"}, Namespace: constants.NamespaceEta, Category: ctype.NewStringPrt(labelColorName), Name: ctype.NewStringPrt("color"), Order: 2, ParentId: "167556962316193797"},
		{BaseEntity: entity.BaseEntity{ID: "167556962316193800", Remark: "green"}, Namespace: constants.NamespaceEta, Category: ctype.NewStringPrt(labelColorName), Name: ctype.NewStringPrt("绿色"), Order: 3},
		{BaseEntity: entity.BaseEntity{ID: "167556962316193801", Remark: "rgba(67, 154, 124, 0.15)"}, Namespace: constants.NamespaceEta, Category: ctype.NewStringPrt(labelColorName), Name: ctype.NewStringPrt("background"), Order: 1, ParentId: "167556962316193800"},
		{BaseEntity: entity.BaseEntity{ID: "167556962316193802", Remark: "#2BB480"}, Namespace: constants.NamespaceEta, Category: ctype.NewStringPrt(labelColorName), Name: ctype.NewStringPrt("color"), Order: 2, ParentId: "167556962316193800"},
		{BaseEntity: entity.BaseEntity{ID: "167556962316193803", Remark: "purple"}, Namespace: constants.NamespaceEta, Category: ctype.NewStringPrt(labelColorName), Name: ctype.NewStringPrt("紫色"), Order: 4},
		{BaseEntity: entity.BaseEntity{ID: "167556962316193804", Remark: "rgba(108, 85, 176, 0.15)"}, Namespace: constants.NamespaceEta, Category: ctype.NewStringPrt(labelColorName), Name: ctype.NewStringPrt("background"), Order: 1, ParentId: "167556962316193803"},
		{BaseEntity: entity.BaseEntity{ID: "167556962316193805", Remark: "#7860CA"}, Namespace: constants.NamespaceEta, Category: ctype.NewStringPrt(labelColorName), Name: ctype.NewStringPrt("color"), Order: 2, ParentId: "167556962316193803"},
		{BaseEntity: entity.BaseEntity{ID: "167556962316193806", Remark: "red"}, Namespace: constants.NamespaceEta, Category: ctype.NewStringPrt(labelColorName), Name: ctype.NewStringPrt("红色"), Order: 5},
		{BaseEntity: entity.BaseEntity{ID: "167556962316193807", Remark: "#EC7281"}, Namespace: constants.NamespaceEta, Category: ctype.NewStringPrt(labelColorName), Name: ctype.NewStringPrt("background"), Order: 1, ParentId: "167556962316193806"},
		{BaseEntity: entity.BaseEntity{ID: "167556962316193808", Remark: "#fff"}, Namespace: constants.NamespaceEta, Category: ctype.NewStringPrt(labelColorName), Name: ctype.NewStringPrt("color"), Order: 2, ParentId: "167556962316193806"},
		{BaseEntity: entity.BaseEntity{ID: "167556962316193809", Remark: "orange"}, Namespace: constants.NamespaceEta, Category: ctype.NewStringPrt(labelColorName), Name: ctype.NewStringPrt("橙色"), Order: 6},
		{BaseEntity: entity.BaseEntity{ID: "167556962316193810", Remark: "rgba(193, 71, 79, 0.6)"}, Namespace: constants.NamespaceEta, Category: ctype.NewStringPrt(labelColorName), Name: ctype.NewStringPrt("background"), Order: 1, ParentId: "167556962316193809"},
		{BaseEntity: entity.BaseEntity{ID: "167556962316193811", Remark: "#fff"}, Namespace: constants.NamespaceEta, Category: ctype.NewStringPrt(labelColorName), Name: ctype.NewStringPrt("color"), Order: 2, ParentId: "167556962316193809"},
		{BaseEntity: entity.BaseEntity{ID: "167556962316193812", Remark: "yellow"}, Namespace: constants.NamespaceEta, Category: ctype.NewStringPrt(labelColorName), Name: ctype.NewStringPrt("黄色"), Order: 7},
		{BaseEntity: entity.BaseEntity{ID: "167556962316193813", Remark: "rgba(229, 155, 0, 0.6)"}, Namespace: constants.NamespaceEta, Category: ctype.NewStringPrt(labelColorName), Name: ctype.NewStringPrt("background"), Order: 1, ParentId: "167556962316193812"},
		{BaseEntity: entity.BaseEntity{ID: "167556962316193814", Remark: "#fff"}, Namespace: constants.NamespaceEta, Category: ctype.NewStringPrt(labelColorName), Name: ctype.NewStringPrt("color"), Order: 2, ParentId: "167556962316193812"},
	}
)
