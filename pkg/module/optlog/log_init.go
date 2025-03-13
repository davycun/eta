package optlog

var (
	uriOptLog = make(map[string]OptLog)
)

func init() {
	uriOptLog["/oauth2/login:POST"] = OptLog{OptType: Login, OptTarget: "用户", OptContent: "用户登录"}
	uriOptLog["/oauth2/logout:POST"] = OptLog{OptType: Logout, OptTarget: "用户", OptContent: "用户退出"}

	uriOptLog["/user/create:POST"] = OptLog{OptType: CreateData, OptTarget: "用户", OptContent: "新增用户"}
	uriOptLog["/user/update:POST"] = OptLog{OptType: UpdateData, OptTarget: "用户", OptContent: "更新用户"}
	uriOptLog["/user/delete:POST"] = OptLog{OptType: DeleteData, OptTarget: "用户", OptContent: "删除用户"}
	uriOptLog["/user/query:POST"] = OptLog{OptType: QueryData, OptTarget: "用户", OptContent: "查询用户"}
	//uriOptLog["/user/detail:POST"] = OptLog{OptType: "详情数据", OptTarget: "用户", OptContent: "查询用户"}
	uriOptLog["/user/count:POST"] = OptLog{OptType: AggregateData, OptTarget: "用户", OptContent: "查询用户"}
	uriOptLog["/user/aggregate:POST"] = OptLog{OptType: AggregateData, OptTarget: "用户", OptContent: "查询用户"}
	uriOptLog["/user/partition:POST"] = OptLog{OptType: QueryData, OptTarget: "用户", OptContent: "查询用户"}

	// app
	uriOptLog["/app/create:POST"] = OptLog{OptType: CreateData, OptTarget: "App", OptContent: "新增App"}
	uriOptLog["/app/update:POST"] = OptLog{OptType: UpdateData, OptTarget: "App", OptContent: "更新App"}
	uriOptLog["/app/delete:POST"] = OptLog{OptType: DeleteData, OptTarget: "App", OptContent: "删除App"}
	uriOptLog["/app/query:POST"] = OptLog{OptType: QueryData, OptTarget: "App", OptContent: "查询App"}
	uriOptLog["/app/count:POST"] = OptLog{OptType: AggregateData, OptTarget: "App", OptContent: "查询App"}
	uriOptLog["/app/aggregate:POST"] = OptLog{OptType: AggregateData, OptTarget: "App", OptContent: "查询App"}
	uriOptLog["/app/partition:POST"] = OptLog{OptType: QueryData, OptTarget: "App", OptContent: "查询App"}
	// task
	uriOptLog["/task/create:POST"] = OptLog{OptType: CreateData, OptTarget: "任务", OptContent: "新增任务"}
	uriOptLog["/task/update:POST"] = OptLog{OptType: UpdateData, OptTarget: "任务", OptContent: "更新任务"}
	uriOptLog["/task/query:POST"] = OptLog{OptType: QueryData, OptTarget: "任务", OptContent: "查询任务"}
	// template
	uriOptLog["/template/create:POST"] = OptLog{OptType: CreateData, OptTarget: "模板", OptContent: "新增模板"}
	uriOptLog["/template/update:POST"] = OptLog{OptType: UpdateData, OptTarget: "模板", OptContent: "更新模板"}
	uriOptLog["/template/delete:POST"] = OptLog{OptType: DeleteData, OptTarget: "模板", OptContent: "删除模板"}
	uriOptLog["/template/query:POST"] = OptLog{OptType: QueryData, OptTarget: "模板", OptContent: "查询模板"}
	uriOptLog["/template/aggregate:POST"] = OptLog{OptType: AggregateData, OptTarget: "模板", OptContent: "查询模板"}
	uriOptLog["/template/partition:POST"] = OptLog{OptType: QueryData, OptTarget: "模板", OptContent: "查询模板"}
	// people
	uriOptLog["/citizen/people/create:POST"] = OptLog{OptType: CreateData, OptTarget: "居民", OptContent: "新增人员"}
	uriOptLog["/citizen/people/update:POST"] = OptLog{OptType: UpdateData, OptTarget: "居民", OptContent: "更新人员"}
	uriOptLog["/citizen/people/delete:POST"] = OptLog{OptType: DeleteData, OptTarget: "居民", OptContent: "删除人员"}
	uriOptLog["/citizen/people/query:POST"] = OptLog{OptType: QueryData, OptTarget: "居民", OptContent: "查询人员"}
	uriOptLog["/citizen/people/count:POST"] = OptLog{OptType: AggregateData, OptTarget: "居民", OptContent: "查询人员"}
	uriOptLog["/citizen/people/aggregate:POST"] = OptLog{OptType: AggregateData, OptTarget: "居民", OptContent: "查询人员"}
	uriOptLog["/citizen/people/partition:POST"] = OptLog{OptType: QueryData, OptTarget: "居民", OptContent: "查询人员"}
	uriOptLog["/citizen/people/list:POST"] = OptLog{OptType: QueryData, OptTarget: "居民", OptContent: "查询全部居民列表"}
	uriOptLog["/citizen/people/list_his:POST"] = OptLog{OptType: QueryData, OptTarget: "居民", OptContent: "查询历史居民列表"}
	uriOptLog["/citizen/people/list_huji_his:POST"] = OptLog{OptType: QueryData, OptTarget: "居民", OptContent: "查询历史户籍居民列表"}
	uriOptLog["/citizen/people/list_stats:POST"] = OptLog{OptType: AggregateData, OptTarget: "居民", OptContent: "统计定制接口居民数据"}
	uriOptLog["/citizen/people/list_addr_agg:POST"] = OptLog{OptType: AggregateData, OptTarget: "居民", OptContent: "透视统计居民数据"}
	uriOptLog["/citizen/people/list_shop:POST"] = OptLog{OptType: QueryData, OptTarget: "居民", OptContent: "查询商业从业人员"}
	uriOptLog["/citizen/people/list_ent:POST"] = OptLog{OptType: QueryData, OptTarget: "居民", OptContent: "查询企业从业人员"}
	uriOptLog["/citizen/people/list_org:POST"] = OptLog{OptType: QueryData, OptTarget: "居民", OptContent: "查询组织人员"}
	uriOptLog["/citizen/people/list_export:POST"] = OptLog{OptType: ExportData, OptTarget: "居民", OptContent: "导出全部居民"}
	uriOptLog["/citizen/people/list_his_export:POST"] = OptLog{OptType: ExportData, OptTarget: "居民", OptContent: "导出历史居民"}
	uriOptLog["/citizen/people/list_addr_agg_export:POST"] = OptLog{OptType: ExportData, OptTarget: "居民", OptContent: "导出透视统计居民"}
	uriOptLog["/citizen/people/list_shop_export:POST"] = OptLog{OptType: ExportData, OptTarget: "居民", OptContent: "导出商业从业人员"}
	uriOptLog["/citizen/people/list_ent_export:POST"] = OptLog{OptType: ExportData, OptTarget: "居民", OptContent: "导出企业从业人员"}
	uriOptLog["/citizen/people/data_import:POST"] = OptLog{OptType: ImportData, OptTarget: "居民", OptContent: "导入居民"}
	// room
	uriOptLog["/citizen/room/create:POST"] = OptLog{OptType: CreateData, OptTarget: "房间", OptContent: "新增房间"}
	uriOptLog["/citizen/room/update:POST"] = OptLog{OptType: UpdateData, OptTarget: "房间", OptContent: "更新房间"}
	uriOptLog["/citizen/room/delete:POST"] = OptLog{OptType: DeleteData, OptTarget: "房间", OptContent: "删除房间"}
	uriOptLog["/citizen/room/query:POST"] = OptLog{OptType: QueryData, OptTarget: "房间", OptContent: "查询房间"}
	uriOptLog["/citizen/room/count:POST"] = OptLog{OptType: AggregateData, OptTarget: "房间", OptContent: "查询房间"}
	uriOptLog["/citizen/room/aggregate:POST"] = OptLog{OptType: AggregateData, OptTarget: "房间", OptContent: "查询房间"}
	uriOptLog["/citizen/room/partition:POST"] = OptLog{OptType: QueryData, OptTarget: "房间", OptContent: "查询房间"}
	uriOptLog["/citizen/room/list:POST"] = OptLog{OptType: QueryData, OptTarget: "房间", OptContent: "定制化查询全部房间列表"}
	uriOptLog["/citizen/room/list_stats:POST"] = OptLog{OptType: AggregateData, OptTarget: "房间", OptContent: "定制化统计房间数据"}
	uriOptLog["/citizen/room/list_export:POST"] = OptLog{OptType: ExportData, OptTarget: "房间", OptContent: "导出房间"}
	// shop
	uriOptLog["/citizen/shop/create:POST"] = OptLog{OptType: CreateData, OptTarget: "商铺", OptContent: "新增商铺"}
	uriOptLog["/citizen/shop/update:POST"] = OptLog{OptType: UpdateData, OptTarget: "商铺", OptContent: "更新商铺"}
	uriOptLog["/citizen/shop/delete:POST"] = OptLog{OptType: DeleteData, OptTarget: "商铺", OptContent: "删除商铺"}
	uriOptLog["/citizen/shop/query:POST"] = OptLog{OptType: QueryData, OptTarget: "商铺", OptContent: "查询商铺"}
	uriOptLog["/citizen/shop/count:POST"] = OptLog{OptType: AggregateData, OptTarget: "商铺", OptContent: "查询商铺"}
	uriOptLog["/citizen/shop/aggregate:POST"] = OptLog{OptType: AggregateData, OptTarget: "商铺", OptContent: "查询商铺"}
	uriOptLog["/citizen/shop/partition:POST"] = OptLog{OptType: QueryData, OptTarget: "商铺", OptContent: "查询商铺"}
	uriOptLog["/citizen/shop/list:POST"] = OptLog{OptType: QueryData, OptTarget: "商铺", OptContent: "定制化查询全部商铺列表"}
	uriOptLog["/citizen/shop/list_addr_agg:POST"] = OptLog{OptType: AggregateData, OptTarget: "商铺", OptContent: "定制化透视统计商铺"}
	uriOptLog["/citizen/shop/list_export:POST"] = OptLog{OptType: ExportData, OptTarget: "商铺", OptContent: "导出商铺"}
	// enterprise
	uriOptLog["/citizen/enterprise/create:POST"] = OptLog{OptType: CreateData, OptTarget: "企业", OptContent: "新增企业"}
	uriOptLog["/citizen/enterprise/update:POST"] = OptLog{OptType: UpdateData, OptTarget: "企业", OptContent: "更新企业"}
	uriOptLog["/citizen/enterprise/delete:POST"] = OptLog{OptType: DeleteData, OptTarget: "企业", OptContent: "删除企业"}
	uriOptLog["/citizen/enterprise/query:POST"] = OptLog{OptType: QueryData, OptTarget: "企业", OptContent: "查询企业"}
	uriOptLog["/citizen/enterprise/count:POST"] = OptLog{OptType: AggregateData, OptTarget: "企业", OptContent: "查询企业"}
	uriOptLog["/citizen/enterprise/aggregate:POST"] = OptLog{OptType: AggregateData, OptTarget: "企业", OptContent: "查询企业"}
	uriOptLog["/citizen/enterprise/partition:POST"] = OptLog{OptType: QueryData, OptTarget: "企业", OptContent: "查询企业"}
	uriOptLog["/citizen/enterprise/list:POST"] = OptLog{OptType: QueryData, OptTarget: "企业", OptContent: "定制化查询全部企业列表"}
	uriOptLog["/citizen/enterprise/list_addr_agg:POST"] = OptLog{OptType: AggregateData, OptTarget: "企业", OptContent: "定制化透视统计企业"}
	uriOptLog["/citizen/enterprise/list_export:POST"] = OptLog{OptType: ExportData, OptTarget: "企业", OptContent: "导出企业"}
	// building
	uriOptLog["/citizen/building/create:POST"] = OptLog{OptType: CreateData, OptTarget: "楼栋", OptContent: "新增楼栋"}
	uriOptLog["/citizen/building/update:POST"] = OptLog{OptType: UpdateData, OptTarget: "楼栋", OptContent: "更新楼栋"}
	uriOptLog["/citizen/building/delete:POST"] = OptLog{OptType: DeleteData, OptTarget: "楼栋", OptContent: "删除楼栋"}
	uriOptLog["/citizen/building/query:POST"] = OptLog{OptType: QueryData, OptTarget: "楼栋", OptContent: "查询楼栋"}
	uriOptLog["/citizen/building/count:POST"] = OptLog{OptType: AggregateData, OptTarget: "楼栋", OptContent: "查询楼栋"}
	uriOptLog["/citizen/building/aggregate:POST"] = OptLog{OptType: AggregateData, OptTarget: "楼栋", OptContent: "查询楼栋"}
	uriOptLog["/citizen/building/partition:POST"] = OptLog{OptType: QueryData, OptTarget: "楼栋", OptContent: "查询楼栋"}
	uriOptLog["/citizen/building/list:POST"] = OptLog{OptType: QueryData, OptTarget: "楼栋", OptContent: "查询全部楼栋列表"}
	uriOptLog["/citizen/building/list_stats:POST"] = OptLog{OptType: AggregateData, OptTarget: "楼栋", OptContent: "定制化统计楼栋"}
	uriOptLog["/citizen/building/list_export:POST"] = OptLog{OptType: ExportData, OptTarget: "楼栋", OptContent: "导出楼栋"}
	// address
	uriOptLog["/citizen/address/create:POST"] = OptLog{OptType: CreateData, OptTarget: "地址", OptContent: "新增地址"}
	uriOptLog["/citizen/address/update:POST"] = OptLog{OptType: UpdateData, OptTarget: "地址", OptContent: "更新地址"}
	uriOptLog["/citizen/address/delete:POST"] = OptLog{OptType: DeleteData, OptTarget: "地址", OptContent: "删除地址"}
	uriOptLog["/citizen/address/query:POST"] = OptLog{OptType: QueryData, OptTarget: "地址", OptContent: "查询地址"}
	uriOptLog["/citizen/address/count:POST"] = OptLog{OptType: AggregateData, OptTarget: "地址", OptContent: "查询地址"}
	uriOptLog["/citizen/address/aggregate:POST"] = OptLog{OptType: AggregateData, OptTarget: "地址", OptContent: "查询地址"}
	uriOptLog["/citizen/address/partition:POST"] = OptLog{OptType: QueryData, OptTarget: "地址", OptContent: "查询地址"}
	uriOptLog["/citizen/address/list:POST"] = OptLog{OptType: QueryData, OptTarget: "地址", OptContent: "定制化查询全部地址列表"}
	uriOptLog["/citizen/address/list_aggregate:POST"] = OptLog{OptType: AggregateData, OptTarget: "地址", OptContent: "定制化地址聚合统计"}
	uriOptLog["/citizen/address/list_export:POST"] = OptLog{OptType: ExportData, OptTarget: "地址", OptContent: "导出地址"}
	// floor
	uriOptLog["/citizen/floor/create:POST"] = OptLog{OptType: CreateData, OptTarget: "楼层", OptContent: "新增楼层"}
	uriOptLog["/citizen/floor/update:POST"] = OptLog{OptType: UpdateData, OptTarget: "楼层", OptContent: "更新楼层"}
	uriOptLog["/citizen/floor/delete:POST"] = OptLog{OptType: DeleteData, OptTarget: "楼层", OptContent: "删除楼层"}
	uriOptLog["/citizen/floor/query:POST"] = OptLog{OptType: QueryData, OptTarget: "楼层", OptContent: "查询楼层"}
	uriOptLog["/citizen/floor/count:POST"] = OptLog{OptType: AggregateData, OptTarget: "楼层", OptContent: "查询楼层"}
	uriOptLog["/citizen/floor/aggregate:POST"] = OptLog{OptType: AggregateData, OptTarget: "楼层", OptContent: "查询楼层"}
	uriOptLog["/citizen/floor/partition:POST"] = OptLog{OptType: QueryData, OptTarget: "楼层", OptContent: "查询楼层"}
	// label
	uriOptLog["/citizen/label/create:POST"] = OptLog{OptType: CreateData, OptTarget: "标签", OptContent: "新增标签"}
	uriOptLog["/citizen/label/update:POST"] = OptLog{OptType: UpdateData, OptTarget: "标签", OptContent: "更新标签"}
	uriOptLog["/citizen/label/delete:POST"] = OptLog{OptType: DeleteData, OptTarget: "标签", OptContent: "删除标签"}
	uriOptLog["/citizen/label/query:POST"] = OptLog{OptType: QueryData, OptTarget: "标签", OptContent: "查询标签"}
	uriOptLog["/citizen/label/count:POST"] = OptLog{OptType: AggregateData, OptTarget: "标签", OptContent: "查询标签"}
	uriOptLog["/citizen/label/aggregate:POST"] = OptLog{OptType: AggregateData, OptTarget: "标签", OptContent: "查询标签"}
	uriOptLog["/citizen/label/partition:POST"] = OptLog{OptType: QueryData, OptTarget: "标签", OptContent: "查询标签"}
	uriOptLog["/citizen/label/tree:POST"] = OptLog{OptType: QueryData, OptTarget: "标签", OptContent: "定制化查询标签树状数据"}
	uriOptLog["/citizen/label/tree_update:POST"] = OptLog{OptType: UpdateData, OptTarget: "标签", OptContent: "定制化更新标签树状数据"}
	uriOptLog["/citizen/label/tree_delete:POST"] = OptLog{OptType: DeleteData, OptTarget: "标签", OptContent: "定制化删除标签树状数据"}
	// organization
	uriOptLog["/citizen/organization/create:POST"] = OptLog{OptType: CreateData, OptTarget: "组织", OptContent: "新增组织"}
	uriOptLog["/citizen/organization/update:POST"] = OptLog{OptType: UpdateData, OptTarget: "组织", OptContent: "更新组织"}
	uriOptLog["/citizen/organization/delete:POST"] = OptLog{OptType: DeleteData, OptTarget: "组织", OptContent: "删除组织"}
	uriOptLog["/citizen/organization/query:POST"] = OptLog{OptType: QueryData, OptTarget: "组织", OptContent: "查询组织"}
	uriOptLog["/citizen/organization/count:POST"] = OptLog{OptType: AggregateData, OptTarget: "组织", OptContent: "查询组织"}
	uriOptLog["/citizen/organization/aggregate:POST"] = OptLog{OptType: AggregateData, OptTarget: "组织", OptContent: "查询组织"}
	uriOptLog["/citizen/organization/partition:POST"] = OptLog{OptType: QueryData, OptTarget: "组织", OptContent: "查询组织"}
	uriOptLog["/citizen/organization/list:POST"] = OptLog{OptType: QueryData, OptTarget: "组织", OptContent: "定制化查询全部组织列表"}
	uriOptLog["/citizen/organization/list_aggregate:POST"] = OptLog{OptType: AggregateData, OptTarget: "组织", OptContent: "定制化聚合组织列表统计"}

}
