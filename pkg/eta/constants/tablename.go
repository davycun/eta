package constants

const (
	TableApp                = "t_app"
	TableAppHistory         = TableApp + TableHistorySubFix
	TableUser               = "t_user"
	TableUserHistory        = TableUser + TableHistorySubFix
	TableUserKey            = "t_user_key"
	TableUserKeyHistory     = TableUserKey + TableHistorySubFix
	TablePermission         = "t_permission"
	TableRole               = "t_role"
	TableAuth2Role          = "r_auth2role"
	TableRole2Role          = "r_role2role"
	TableUser2App           = "r_user2app"
	TableUser2AppHistory    = TableUser2App + TableHistorySubFix
	TableUser2Dept          = "r_user2dept"
	TableUser2Role          = "r_user2role"
	TableMenu               = "t_menu"
	TableDept               = "t_department"
	TableDeptHistory        = TableDept + TableHistorySubFix
	TableSetting            = "t_setting"
	TableConfig             = "t_config"
	TableTemplate           = "t_template"
	TableTemplateHistory    = "t_template_history"
	TableOperateLog         = "t_operate_log"
	TableTransferKey        = "t_transfer_key"
	TableTransferKeyHistory = TableTransferKey + TableHistorySubFix
	TableDictionary         = "t_dictionary"
	TableSubscriber         = "t_subscriber"
	TablePublishRecord      = "t_publish_record"
	TableDataTask           = "t_data_task"

	// TableSmsTask 短信相关的表
	TableSmsTask   = "t_sms_task"
	TableSmsTarget = "t_sms_target"

	TableTemplatePrefix     = "d_"       //模板表的统一前缀
	TableHistorySubFix      = "_history" //历史表的名称后缀
	TableTriggerPrefix      = "trigger_" //历史记录触发器名称前缀
	TableHistoryFieldPrefix = "h_"       //历史表的业务字段前缀
)
