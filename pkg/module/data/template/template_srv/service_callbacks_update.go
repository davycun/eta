package template_srv

import (
	"errors"
	"fmt"
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/history"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/core/ra"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/core/updater"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/data/template"
	"gorm.io/gorm"
	"regexp"
	"strings"
	"time"
)

func beforeUpdateValidate(cfg *hook.SrvConfig, oldValues []template.Template, newValues []template.Template) error {
	regex := regexp.MustCompile(`^[\w\p{Han}]+$`)
	if cfg.Method == iface.MethodUpdate {
		newValueMap := make(map[string]template.Template, len(newValues))
		for _, v := range newValues {
			tbFields := v.GetTable().GetFields()
			for j, _ := range tbFields {
				f := &tbFields[j]
				if !regex.MatchString(f.Name) {
					return errors.New(fmt.Sprintf("field name %s is invalid", f.Name))
				}
			}
			newValueMap[v.ID] = v
		}
		for _, v := range oldValues {
			if v.GetTableName() != newValueMap[v.ID].GetTableName() && newValueMap[v.ID].GetTableName() != "" {
				return errors.New("template name can not be changed")
			}
			if v.Code != newValueMap[v.ID].Code && newValueMap[v.ID].Code != "" {
				return errors.New("template code can not be changed")
			}
		}
	} else {
		newValue := newValues[0]
		tbFields := newValue.Table.GetFields()
		for j, _ := range tbFields {
			f := &tbFields[j]
			if !regex.MatchString(f.Name) {
				return errors.New(fmt.Sprintf("field name %s is invalid", f.Name))
			}
		}
		for _, v := range oldValues {
			if v.GetTableName() != newValue.GetTableName() && newValue.GetTableName() != "" {
				return errors.New("template name can not be changed")
			}
			if v.Code != newValue.Code && newValue.Code != "" {
				return errors.New("template code can not be changed")
			}
		}
	}
	return nil
}

func afterUpdate(cfg *hook.SrvConfig, oldValues []template.Template, newValues []template.Template) error {
	var (
		db          = cfg.TxDB
		scm         = dorm.GetDbSchema(db)
		dbType      = dorm.GetDbType(db)
		oldValueMap = make(map[string]template.Template)
	)
	for _, v := range oldValues {
		oldValueMap[v.ID] = v
	}
	hasHistoryTableMap, err := buildHasHistoryTableMap(db, dbType, scm, newValues)
	if err != nil {
		return err
	}
	for _, v := range newValues {
		hasHistoryTable := hasHistoryTableMap[v.HistoryTableName()]
		err = caller.NewCaller().
			Call(func(cl *caller.Caller) error {
				return handleFields(db, dbType, scm, oldValueMap[v.ID], v, hasHistoryTable)
			}).
			Call(func(cl *caller.Caller) error {
				return handleIndexes(db, dbType, scm, oldValueMap[v.ID], v)
			}).
			Call(func(cl *caller.Caller) error {
				return updateHistoryTrigger(db, cfg.Ctx, scm, oldValueMap[v.ID], v, hasHistoryTable)
			}).
			Call(func(cl *caller.Caller) error {
				return updateFieldUpdaterTrigger(db, oldValueMap[v.ID], v)
			}).
			Call(func(cl *caller.Caller) error {
				return updateRaTrigger(db, oldValueMap[v.ID], v)
			}).Err

		if err != nil {
			return err
		}
	}

	return err
}

func buildHasHistoryTableMap(db *gorm.DB, dbType dorm.DbType, scm string, templates []template.Template) (map[string]bool, error) {
	tableNames := make([]string, 0, len(templates))
	for _, v := range templates {
		historyTableName := v.HistoryTableName()
		tableNames = append(tableNames, fmt.Sprintf(`'%s'`, historyTableName))
	}
	res, err := getHasHistoryTableMap(db, dbType, scm, tableNames)
	return res, err
}

func getHasHistoryTableMap(db *gorm.DB, dbType dorm.DbType, scm string, tableNames []string) (map[string]bool, error) {
	tableNamesStr := strings.Join(tableNames, ",")
	res := make([]string, 0, len(tableNames))
	switch dbType {
	case dorm.PostgreSQL:
		err := db.Raw(fmt.Sprintf(`select table_name as count from information_schema.tables where table_schema = '%s' and table_name in (%s)`, scm, tableNamesStr)).Scan(&res).Error
		if err != nil {
			return nil, err
		}
	case dorm.DaMeng:
		err := db.Raw(fmt.Sprintf(`select table_name from dba_tables where owner = '%s' and table_name in (%s)`, scm, tableNamesStr)).Scan(&res).Error
		if err != nil {
			return nil, err
		}
	case dorm.Mysql:
		err := db.Raw(fmt.Sprintf(`select table_name from information_schema.tables where table_schema = '%s' and table_name in (%s)`, scm, tableNamesStr)).Scan(&res).Error
		if err != nil {
			return nil, err
		}
	default:
		//not support
	}
	resMap := make(map[string]bool, len(res))
	for _, v := range res {
		resMap[v] = true
	}
	return resMap, nil
}

func getTableTriggerNames(db *gorm.DB, dbType dorm.DbType, scm, tableName string) ([]string, error) {
	var (
		triggerNames = make([]string, 0, 10)
		triggerSql   string
	)
	switch dbType {
	case dorm.PostgreSQL:
		triggerSql = fmt.Sprintf(`select tgname from pg_trigger where tgname like '%%trigger_%s_%%'`, tableName)
	case dorm.Mysql:
		triggerSql = fmt.Sprintf(`select trigger_name from information_schema.triggers where trigger_schema = '%s' and trigger_name like '%%trigger_%s_%%'`, scm, tableName)
	case dorm.DaMeng:
		triggerSql = fmt.Sprintf(`SELECT trigger_name FROM DBA_TRIGGERS where table_owner = '%s' and table_name = '%s' and trigger_name like '%%trigger_%s_%%'`, scm, tableName, tableName)
	default:
		//not support
	}

	err := db.Raw(triggerSql).Scan(&triggerNames).Error
	if err != nil {
		return nil, err
	}
	return triggerNames, nil
}

func getTableIndexNames(db *gorm.DB, dbType dorm.DbType, scm, tableName string) ([]string, error) {
	var (
		indexNames = make([]string, 0, 10)
		indexSql   string
	)
	switch dbType {
	case dorm.PostgreSQL:
		indexSql = fmt.Sprintf(`select indexname from pg_indexes where schemaname = '%s' and tablename = '%s' and indexname like '%%idx_%s_%%'`, scm, tableName, tableName)
	case dorm.DaMeng:
		indexSql = fmt.Sprintf(`select index_name as indexname from dba_ind_columns where table_owner = '%s' and table_name = '%s' and index_name like '%%idx_%s_%%'`, scm, tableName, tableName)
	case dorm.Mysql:
		indexSql = fmt.Sprintf(`select index_name  as indexname from information_schema.statistics where table_schema = '%s' and table_name = '%s' and index_name like '%%idx_%s_%%'`, scm, tableName, tableName)
	default:
		//not support
	}

	err := db.Raw(indexSql).Scan(&indexNames).Error
	if err != nil {
		return nil, err
	}
	return indexNames, nil
}

// handleIndexes 处理索引
func handleIndexes(db *gorm.DB, dbType dorm.DbType, scm string, origin, target template.Template) error {
	// 查询已经存在的索引
	originIndexNames, err := getTableIndexNames(db, dbType, scm, target.GetTableName())
	if err != nil {
		return err
	}
	// 构造更新后的索引
	targetIndexNames := make([]string, 0, 10)
	if len(target.Table.Indexes) > 0 {
		err = validateIndex(&target)
		if err != nil {
			return err
		}
		for _, v := range target.Table.Indexes {
			targetIndexNames = append(targetIndexNames, v.IndexName(target.GetTableName()))
		}
	}

	redundantIndexNames := utils.DifferenceOfStringSlices(originIndexNames, targetIndexNames)
	err = removeIndexes(db, dbType, scm, target.GetTableName(), redundantIndexNames)
	if err != nil {
		return err
	}
	newIndexNames := utils.DifferenceOfStringSlices(targetIndexNames, originIndexNames)
	err = addIndexes(db, target, newIndexNames)
	if err != nil {
		return err
	}
	intersectionIndexNames := utils.IntersectionOfStringSlices(targetIndexNames, originIndexNames)
	err = updateIndexes(db, dbType, scm, origin, target, intersectionIndexNames)
	if err != nil {
		return err
	}
	return nil
}

func removeIndexes(db *gorm.DB, dbType dorm.DbType, scm, tableName string, idxNames []string) error {
	// 删除多余的索引
	for _, v := range idxNames {
		var dropIndexSql string
		switch dbType {
		case dorm.PostgreSQL, dorm.DaMeng:
			dropIndexSql = fmt.Sprintf(`drop index %s`, dorm.GetScmTableName(db, v))
		case dorm.Mysql:
			dropIndexSql = fmt.Sprintf(`drop index %s on %s`, dorm.Quote(dbType, v), dorm.GetScmTableName(db, tableName))
		default:
			//not support
		}
		err := db.Exec(dropIndexSql).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func addIndexes(db *gorm.DB, template template.Template, idxNames []string) error {
	// 添加新增的索引
	for _, v := range template.Table.Indexes {
		indexName := v.IndexName(template.GetTableName())
		if isInSlice(idxNames, indexName) {
			indexSqlStr, err := buildIndexSql(db, template.GetTableName(), v)
			if err != nil {
				return err
			}
			err = db.Exec(indexSqlStr).Error
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func updateIndexes(db *gorm.DB, dbType dorm.DbType, scm string, origin, target template.Template, idxNames []string) error {
	// 更新已经存在的索引
	if len(idxNames) == 0 {
		return nil
	}
	var (
		originIndexMap = make(map[string]entity.TableIndex)
		targetIndexMap = make(map[string]entity.TableIndex)
	)
	for _, v := range origin.Table.Indexes {
		indexName := v.IndexName(target.GetTableName())
		if isInSlice(idxNames, indexName) {
			originIndexMap[indexName] = v
		}
	}
	for _, v := range target.Table.Indexes {
		indexName := v.IndexName(target.GetTableName())
		if isInSlice(idxNames, indexName) {
			targetIndexMap[indexName] = v
		}
	}

	updateIndexNames := make([]string, 0, 10)
	for key, val := range originIndexMap {
		if targetIndexMap[key].Type != val.Type || targetIndexMap[key].Class != val.Class || targetIndexMap[key].Option != val.Option {
			updateIndexNames = append(updateIndexNames, key)
		}
	}
	err := removeIndexes(db, dbType, scm, origin.GetTableName(), updateIndexNames)
	if err != nil {
		return err
	}
	err = addIndexes(db, target, updateIndexNames)
	if err != nil {
		return err
	}
	return nil
}

func isInSlice(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func handleFields(db *gorm.DB, dbType dorm.DbType, scm string, origin, target template.Template, hasHistoryTable bool) error {
	// 获取已经存在的字段
	originFieldNames := getTableFieldNames(origin)
	targetFieldNames := getTableFieldNames(target)

	// 重命名删除的字段
	redundantFieldNames := utils.DifferenceOfStringSlices(originFieldNames, targetFieldNames)
	err := removeFields(db, scm, origin.GetTableName(), redundantFieldNames, false)
	if err != nil {
		return err
	}
	// 添加新增字段
	newFieldNames := utils.DifferenceOfStringSlices(targetFieldNames, originFieldNames)
	err = addFields(db, dbType, scm, target, newFieldNames, false)
	if err != nil {
		return err
	}
	// 更新字段
	intersectionFieldNames := utils.IntersectionOfStringSlices(targetFieldNames, originFieldNames)
	err = updateFields(db, dbType, scm, origin, target, intersectionFieldNames, false)
	if err != nil {
		return err
	}
	// 如果存在历史表 则处理历史表相关字段
	if hasHistoryTable {
		err = removeFields(db, scm, origin.HistoryTableName(), redundantFieldNames, true)
		if err != nil {
			return err
		}
		err = addFields(db, dbType, scm, target, newFieldNames, true)
		if err != nil {
			return err
		}
		err = updateFields(db, dbType, scm, origin, target, intersectionFieldNames, true)
		if err != nil {
			return err
		}
	}
	return nil
}

func getTableFieldNames(template template.Template) []string {
	var fieldNames []string

	for _, v := range template.GetTable().GetFields() {
		fieldNames = append(fieldNames, v.Name)
	}
	return fieldNames
}

func removeFields(db *gorm.DB, scm, tableName string, fieldNames []string, isHistory bool) error {
	for _, v := range fieldNames {
		originColName := fmt.Sprintf("%s", v)
		targetColName := fmt.Sprintf("del_%s_%s", time.Now().Format("20060102150405"), v)
		if isHistory {
			targetColName = fmt.Sprintf("%sdel_%s_%s", constants.TableHistoryFieldPrefix, time.Now().Format("20060102150405"), v)
			originColName = constants.TableHistoryFieldPrefix + v
		}
		dbType := dorm.GetDbType(db)
		var err error
		switch dorm.GetDbType(db) {
		case dorm.PostgreSQL, dorm.DaMeng, dorm.Mysql:
			dropFieldSql := fmt.Sprintf(`ALTER TABLE %s RENAME COLUMN %s TO %s`, dorm.GetScmTableName(db, tableName), dorm.Quote(dbType, originColName), dorm.Quote(dbType, targetColName))
			err = db.Exec(dropFieldSql).Error
		default:
			err = errors.New("不支持的数据库")
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func addFields(db *gorm.DB, dbType dorm.DbType, scm string, template template.Template, fieldNames []string, isHistory bool) error {
	if len(fieldNames) == 0 {
		return nil
	}
	var (
		bd         = strings.Builder{}
		fieldMap   = make(map[string]entity.TableField)
		commentSql = make([]string, 0, len(fieldNames))
	)
	for _, v := range template.Table.GetFields() {
		fieldMap[v.Name] = v
	}

	tableName := template.GetTableName()
	if isHistory {
		tableName = template.HistoryTableName()
	}

	bd.WriteString(fmt.Sprintf("alter table %s ", dorm.GetScmTableName(db, tableName)))
	for i, v := range fieldNames {
		field := fieldMap[v]
		if i > 0 {
			bd.WriteByte(',')
		}
		colName := v
		if isHistory {
			colName = constants.TableHistoryFieldPrefix + v
		}
		typeName, err := ctype.GetDbTypeName(db, field.Type)
		if err != nil {
			return err
		}
		switch dbType {
		case dorm.PostgreSQL:
			bd.WriteString(fmt.Sprintf(` add column %s %s`, dorm.Quote(dbType, colName), typeName))
			if field.Comment != "" {
				commentSql = append(commentSql, buildCommentSql(db, tableName, colName, field.Comment))
			}
		case dorm.DaMeng:
			if i == 0 {
				bd.WriteString(" add column (")
			}
			bd.WriteString(fmt.Sprintf(` %s %s`, dorm.Quote(dbType, colName), typeName))
			if i == (len(fieldNames) - 1) {
				bd.WriteString(" )")
			}
			if field.Comment != "" {
				commentSql = append(commentSql, buildCommentSql(db, tableName, colName, field.Comment))
			}
		case dorm.Mysql:
			if i == 0 {
				bd.WriteString(" add column (")
			}
			bd.WriteString(fmt.Sprintf(` %s %s`, dorm.Quote(dbType, colName), typeName))
			if field.Comment != "" {
				bd.WriteString(" comment '" + field.Comment + "'")
			}
			if i == (len(fieldNames) - 1) {
				bd.WriteString(" )")
			}
		default:
			//not support
		}
	}
	tx1 := db.Exec(bd.String())
	if tx1.Error == nil {
		for _, v := range commentSql {
			if v == "" {
				continue
			}
			tx1 = db.Exec(v)
			if tx1.Error != nil {
				return tx1.Error
			}
		}
	}
	return nil
}

func updateFields(db *gorm.DB, dbType dorm.DbType, scm string, origin, target template.Template, fieldNames []string, isHistory bool) error {
	var (
		originFieldMap    = make(map[string]entity.TableField)
		targetFieldMap    = make(map[string]entity.TableField)
		sqlBd             = strings.Builder{}
		changedFieldNames = make([]string, 0, len(fieldNames))
		commentSql        = make([]string, 0, len(fieldNames))
		addDefaultSql     = make([]string, 0, len(fieldNames))
		dropDefaultSql    = make([]string, 0, len(fieldNames))
	)
	for _, v := range origin.GetTable().GetFields() {
		originFieldMap[v.Name] = v
	}
	for _, v := range target.GetTable().GetFields() {
		targetFieldMap[v.Name] = v
	}
	for _, v := range fieldNames {
		var (
			originField = originFieldMap[v]
			targetField = targetFieldMap[v]
		)
		if originField.Type == targetField.Type && originField.Comment == targetField.Comment && originField.Default == targetField.Default {
			continue
		}
		changedFieldNames = append(changedFieldNames, v)
	}
	if len(changedFieldNames) == 0 {
		return nil
	}

	tableName := target.GetTableName()
	if isHistory {
		tableName = target.HistoryTableName()
	}
	sqlBd.WriteString(fmt.Sprintf("alter table %s ", dorm.GetScmTableName(db, tableName)))
	sqlBdChanged := false
	for i, v := range changedFieldNames {
		var (
			originField = originFieldMap[v]
			targetField = targetFieldMap[v]
		)
		colName := v
		if isHistory {
			colName = constants.TableHistoryFieldPrefix + v
		}
		if i > 0 {
			sqlBd.WriteByte(',')
		}
		typeName, err := ctype.GetDbTypeName(db, targetField.Type)
		if err != nil {
			return err
		}
		switch dbType {
		case dorm.PostgreSQL:
			if originField.Type != targetField.Type {
				sqlBd.WriteString(fmt.Sprintf(` alter column %s type %s using %s::%s`, dorm.Quote(dbType, colName), typeName, dorm.Quote(dbType, colName), typeName))
				sqlBdChanged = true
			}
			if targetField.Default != originField.Default {
				if targetField.Default != "" {
					dropDefaultSql = append(dropDefaultSql, buildDefaultSql(db, scm, tableName, colName, "", targetField.Type))
					addDefaultSql = append(addDefaultSql, buildDefaultSql(db, scm, tableName, colName, targetField.Default, targetField.Type))
				} else {
					dropDefaultSql = append(dropDefaultSql, buildDefaultSql(db, scm, tableName, colName, targetField.Default, targetField.Type))
				}
			}
			if targetField.Comment != originField.Comment {
				commentSql = append(commentSql, buildCommentSql(db, tableName, colName, targetField.Comment))
			}
		case dorm.DaMeng:
			if i == 0 {
				sqlBd.WriteString(" modify (")
			}
			sqlBd.WriteString(fmt.Sprintf(`%s`, dorm.Quote(dbType, colName)))
			if originField.Type != targetField.Type {
				sqlBd.WriteString(fmt.Sprintf(` %s`, typeName))
				sqlBdChanged = true
			}
			if targetField.Default != originField.Default {
				if targetField.Default != "" {
					sqlBd.WriteString(" default " + handleDefault(targetField.Type, targetField.Default))
					sqlBdChanged = true
				} else {
					dropDefaultSql = append(dropDefaultSql, buildDefaultSql(db, scm, tableName, colName, targetField.Default, targetField.Type))
				}
			}
			if i == (len(changedFieldNames) - 1) {
				sqlBd.WriteString(")")
			}
			if targetField.Comment != originField.Comment {
				commentSql = append(commentSql, buildCommentSql(db, tableName, colName, targetField.Comment))
			}
		case dorm.Mysql:
			sqlBd.WriteString(fmt.Sprintf(` modify %s %s`, dorm.Quote(dbType, colName), typeName))
			if targetField.Default != originField.Default {
				if targetField.Default != "" {
					sqlBd.WriteString(" default " + handleDefault(targetField.Type, targetField.Default))
				} else {
					sqlBd.WriteString(" default null")
				}
			}
			if targetField.Comment != originField.Comment {
				if targetField.Comment != "" {
					sqlBd.WriteString(" comment '" + targetField.Comment + "'")
				} else {
					sqlBd.WriteString(" null")
				}
			}
			sqlBdChanged = true
		default:
			//not support
		}
	}
	// 先删除默认值 否则有些字段无法改变type
	for _, v := range dropDefaultSql {
		err := db.Exec(v).Error
		if err != nil {
			return err
		}
	}
	if sqlBdChanged {
		err := db.Exec(sqlBd.String()).Error
		if err != nil {
			return err
		}
	}
	for _, v := range commentSql {
		if v == "" {
			continue
		}
		tx1 := db.Exec(v)
		if tx1.Error != nil {
			return tx1.Error
		}
	}
	for _, v := range addDefaultSql {
		tx1 := db.Exec(v)
		if tx1.Error != nil {
			return tx1.Error
		}
	}
	return nil
}

func buildDefaultSql(db *gorm.DB, scm, tableName string, columnName, columnDefaultValue, columnType string) string {
	dbType := dorm.GetDbType(db)
	switch dbType {
	case dorm.PostgreSQL, dorm.DaMeng, dorm.Mysql:
		if columnDefaultValue == "" {
			return fmt.Sprintf(`alter table %s alter column %s set default NULL`, dorm.GetScmTableName(db, tableName), dorm.Quote(dbType, columnName))
		}
		return fmt.Sprintf(`alter table %s alter column %s set default %s`, dorm.GetScmTableName(db, tableName), dorm.Quote(dbType, columnName), handleDefault(columnType, columnDefaultValue))
	default:
		//not support
	}
	return ""
}

func handleDefault(columnType string, defaultVal string) string {
	if columnType == ctype.TypeBoolName || columnType == ctype.TypeArrayIntName || columnType == ctype.TypeIntegerName || columnType == ctype.TypeBigIntegerName || columnType == ctype.TypeTimeName || columnType == ctype.TypeNumericName {
		return fmt.Sprintf("%s", defaultVal)
	} else {
		return fmt.Sprintf(`'%s'`, defaultVal)
	}
}

func handleTableName(db *gorm.DB, dbType dorm.DbType, scm string, origin, target template.Template) error {
	if origin.GetTableName() != target.GetTableName() {
		err := renameIndexNames(db, dbType, scm, origin, target)
		if err != nil {
			return err
		}
		err = renameTriggerNames(db, dbType, scm, origin, target)
		if err != nil {
			return err
		}
		err = renameTableNames(db, origin, target)
		if err != nil {
			return err
		}
	}
	return nil
}

func renameTriggerNames(db *gorm.DB, dbType dorm.DbType, scm string, origin, target template.Template) error {
	originTriggerNames, err := getTableTriggerNames(db, dbType, scm, origin.GetTableName())
	if err != nil {
		return err
	}
	for _, v := range originTriggerNames {
		err = db.Exec(buildRenameTriggerSql(dbType, scm, origin.GetTableName(), v, strings.Replace(v, origin.GetTableName(), target.GetTableName(), 1))).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func buildRenameTriggerSql(dbType dorm.DbType, scm, tableName, oldName, newName string) string {
	switch dbType {
	case dorm.PostgreSQL:
		return fmt.Sprintf(`alter trigger "%s" on "%s" rename to "%s"`, oldName, tableName, newName)
	case dorm.DaMeng:
		//not support
	default:
		//not support
	}
	return ""
}

func renameTableNames(db *gorm.DB, origin, target template.Template) (err error) {
	dbType := dorm.GetDbType(db)
	err = db.Exec(fmt.Sprintf(`alter table %s rename to %s`, dorm.GetScmTableName(db, origin.GetTableName()), dorm.Quote(dbType, target.GetTableName()))).Error
	if err != nil {
		return err
	}
	err = db.Exec(fmt.Sprintf(`alter table %s rename to %s`, dorm.GetScmTableName(db, origin.GetTableName()+"_history"), dorm.Quote(dbType, target.GetTableName()+"_history"))).Error
	if err != nil {
		return err
	}
	return nil
}

func renameIndexNames(db *gorm.DB, dbType dorm.DbType, scm string, origin, target template.Template) error {
	originIndexNames, err := getTableIndexNames(db, dbType, scm, origin.GetTableName())
	if err != nil {
		return err
	}
	//originHistoryIndexNames, err := getTableIndexNames(db, dbType, scm, origin.Name+"_history")
	//if err != nil {
	//	return err
	//}
	for _, v := range originIndexNames {
		err = db.Exec(buildRenameIndexSql(db, scm, origin.GetTableName(), v, strings.Replace(v, origin.GetTableName(), target.GetTableName(), 1))).Error
		if err != nil {
			return err
		}
	}
	//for _, v := range originHistoryIndexNames {
	//	err = db.Exec(buildRenameIndexSql(dbType, scm, origin.Name, v, strings.Replace(v, origin.Name, target.Name+"_history", 1))).Error
	//	if err != nil {
	//		return err
	//	}
	//}
	return nil
}

func buildRenameIndexSql(db *gorm.DB, scm, tableName, oldName, newName string) string {
	dbType := dorm.GetDbType(db)
	switch dbType {
	case dorm.PostgreSQL, dorm.DaMeng:
		return fmt.Sprintf(`alter index %s.%s rename to %s`, dorm.Quote(dbType, scm), dorm.Quote(dbType, oldName), dorm.Quote(dbType, newName))
	case dorm.Mysql:
		return fmt.Sprintf(`alter table %s rename index %s to %s`, dorm.GetScmTableName(db, tableName), dorm.Quote(dbType, oldName), dorm.Quote(dbType, newName))
	default:
		//not support
	}
	return ""
}

func updateHistoryTrigger(db *gorm.DB, ctx *ctx.Context, scm string, origin, target template.Template, hasHistoryTable bool) error {
	if ctype.Bool(target.Table.History) {
		return history.CreateTrigger(db, scm, target.GetTableName())
	} else {
		return history.DropTrigger(db, target.GetTableName())
	}
}

func updateFieldUpdaterTrigger(db *gorm.DB, origin, target template.Template) error {
	if ctype.Bool(target.Table.FieldUpdater) {
		return updater.CreateUpdaterTrigger(db, target.GetTableName())
	} else {
		return updater.DropUpdaterTrigger(db, target.GetTableName())
	}
}

func updateRaTrigger(db *gorm.DB, origin, target template.Template) error {
	if len(target.Table.RaDbFields) > 0 {
		return ra.CreateTrigger(db, target.GetTableName(), target.Table.RaDbFields)
	} else {
		return ra.DropTrigger(db, target.GetTableName())
	}
}
