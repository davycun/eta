package template_srv

import (
	"errors"
	"fmt"
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/id/nanoid"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/history"
	"github.com/davycun/eta/pkg/core/ra"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/core/updater"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/template"
	"github.com/duke-git/lancet/v2/slice"
	"gorm.io/gorm"
	"regexp"
	"sort"
	"strings"
)

func beforeCreateFillField(cfg *hook.SrvConfig, newValues []template.Template) error {
	regex := regexp.MustCompile(`^[\w\p{Han}]+$`)
	for i, _ := range newValues {
		t := &newValues[i]
		tbFields := t.Table.GetFields()
		for j, _ := range tbFields {
			f := &tbFields[j]
			if !regex.MatchString(f.Name) {
				return errors.New(fmt.Sprintf("field name %s is invalid", f.Name))
			}
		}
		if t.Code == "" {
			t.Code = nanoid.New()
		}
		if t.Table.TableName == "" {
			t.Table.TableName = constants.TableTemplatePrefix + t.Code
		}
	}
	return nil
}

func afterCreate(cfg *hook.SrvConfig, newValues []template.Template) error {
	var (
		err error
	)
	for i, _ := range newValues {
		p := &newValues[i]
		err = caller.NewCaller().
			Call(func(cl *caller.Caller) error {
				return createTable(cfg.TxDB, p)
			}).
			Call(func(cl *caller.Caller) error {
				return createIndex(cfg.TxDB, p)
			}).
			Call(func(cl *caller.Caller) error {
				return createHistoryTrigger(cfg.TxDB, cfg.Ctx, p)
			}).
			Call(func(cl *caller.Caller) error {
				return createFieldUpdaterTrigger(cfg.TxDB, p)
			}).
			Call(func(cl *caller.Caller) error {
				return createRaTrigger(cfg.TxDB, p)
			}).Err
	}

	return err
}

func createTable(db *gorm.DB, p *template.Template) error {

	fs := make([]entity.TableField, len(p.Table.GetFields()))
	fields := getBaseField("")
	copy(fs, p.Table.GetFields())
	fields = append(fields, fs...)
	scm := dorm.GetDbSchema(db)

	sql, commentSql, err := buildTableAndCommentSql(db, scm, p.GetTableName(), fields)
	if err != nil {
		return err
	}

	tx1 := db.Exec(sql)
	if tx1.Error == nil {
		for _, v := range commentSql {
			if v == "" {
				continue
			}
			tx1 = db.Exec(v)
			if tx1.Error != nil {
				break
			}
		}
	}

	return tx1.Error
}

func createIndex(db *gorm.DB, p *template.Template) error {
	var (
		indexSql = make([]string, 0, 10)
	)

	if p.Table.Indexes != nil && len(p.Table.Indexes) > 0 {
		err := validateIndex(p)
		if err != nil {
			return err
		}
		for _, v := range p.Table.Indexes {
			//indexSql = append(indexSql, buildIndexSql(p.Schema+"."+p.Name, v))
			indexSqlStr, err1 := buildIndexSql(db, p.GetTableName(), v)
			if err1 != nil {
				return err1
			}
			indexSql = append(indexSql, indexSqlStr)
		}
	}
	for _, v := range indexSql {
		tx := db.Exec(v)
		if tx.Error != nil {
			return tx.Error
		}
	}
	return nil
}

func createHistoryTrigger(db *gorm.DB, ctx *ctx.Context, p *template.Template) error {

	if !ctype.Bool(p.Table.History) {
		return nil
	}
	var (
		err        error
		hisFields  = getBaseHistoryField(db)
		baseFields = getBaseField(constants.TableHistoryFieldPrefix)
		ps         = make([]entity.TableField, len(p.Table.GetFields()))
		scm        = dorm.GetDbSchema(db)
		tbName     = p.GetTableName()
	)
	for i, v := range p.Table.GetFields() {
		v.Name = constants.TableHistoryFieldPrefix + v.Name
		ps[i] = v
	}

	hisFields = append(hisFields, baseFields...)
	hisFields = append(hisFields, ps...)

	// 添加已删除的字段（以del_开头的字段）
	// 因为可能先创建表，然后删除过字段（删除字段非物理删除），最后在启用History，这时候需要把非物理删除的字段也加入
	cols := dorm.FetchTableColumns(db, scm, tbName)
	sort.Strings(cols)
	deleteCols := make(map[string][]string, len(cols))
	for _, col := range cols {
		if strings.HasPrefix(col, "del_") {
			// del_20230101235959_ 开头的字段是已删除的字段
			prefix := col[:19]
			colName := col[19:]
			deleteCols[colName] = append(deleteCols[colName], prefix)
		}
	}
	if len(deleteCols) > 0 {
		var (
			wh  = fmt.Sprintf("%s = ?", dorm.Quote(dorm.GetDbType(db), "h_code"))
			ord = fmt.Sprintf(`%s asc`, dorm.Quote(dorm.GetDbType(db), "created_at"))
			rs  []template.History
		)
		err = dorm.Table(db, constants.TableTemplateHistory).Where(wh, p.Code).Limit(1000).Order(ord).Find(&rs).Error
		if err != nil {
			return err
		}
		for _, v := range rs {
			for _, col := range v.Entity.Table.GetFields() {
				if _, ok := deleteCols[col.Name]; ok {
					tmpName := col.Name
					col.Name = constants.TableHistoryFieldPrefix + deleteCols[col.Name][0] + col.Name
					hisFields = append(hisFields, col)
					deleteCols[tmpName] = append(deleteCols[tmpName][:0], deleteCols[tmpName][1:]...)
					if len(deleteCols[tmpName]) == 0 {
						delete(deleteCols, tmpName)
					}
				}
			}
		}
	}
	sql, commentSql, err := buildTableAndCommentSql(db, scm, p.HistoryTableName(), hisFields)
	if err != nil {
		return err
	}

	err = db.Exec(sql).Error
	if err != nil {
		return err
	}
	for _, v := range commentSql {
		if v == "" {
			continue
		}
		err = db.Exec(v).Error
		if err != nil {
			return err
		}
	}

	err = history.CreateTrigger(db, scm, tbName)
	if err != nil {
		return err
	}
	return nil
}
func createFieldUpdaterTrigger(db *gorm.DB, p *template.Template) error {
	if p.Table.FieldUpdater.Data {
		return updater.CreateUpdaterTrigger(db, p.GetTableName())
	}
	return nil
}
func createRaTrigger(db *gorm.DB, p *template.Template) error {
	if len(p.Table.RaDbFields) > 0 {
		return ra.CreateTrigger(db, p.GetTableName(), p.Table.RaDbFields)
	}
	return nil
}

func buildTableAndCommentSql(db *gorm.DB, scm, tableName string, fields []entity.TableField) (string, []string, error) {

	var (
		bd            = strings.Builder{}
		commentSql    = make([]string, 0, len(fields))
		varcharFields = []string{entity.IdDbName, entity.CreatorIdDbName, entity.UpdaterIdDbName, entity.CreatorDeptIdDbName, entity.UpdaterDeptIdDbName}
	)
	bd.WriteString("CREATE TABLE IF NOT EXISTS ")
	bd.WriteString(dorm.GetScmTableName(db, tableName))
	bd.WriteString(" (")
	for i, v := range fields {
		if i > 0 {
			bd.WriteByte(',')
		}
		name, err := ctype.GetDbTypeName(db, v.Type)
		if err != nil {
			return "", nil, err
		}
		if slice.Contain(varcharFields, v.Name) {
			name = "varchar(255)"
		}

		bd.WriteString(fmt.Sprintf(`%s %s`, dorm.Quote(dorm.GetDbType(db), v.Name), name))
		if v.Default != "" {
			bd.WriteString(fmt.Sprintf(` default %s`, handleDefault(v.Type, v.Default)))
		}
		if v.Comment != "" {
			commentSql = append(commentSql, buildCommentSql(db, tableName, v.Name, v.Comment))
		}
	}
	bd.WriteString(fmt.Sprintf(`,PRIMARY KEY (%s) )`, dorm.Quote(dorm.GetDbType(db), "id")))
	return bd.String(), commentSql, nil
}

// 公共字段，prefix主要是为了处理 历史表中所有的基础字段和原表的基础字段冲突的问题
func getBaseField(prefix string) []entity.TableField {

	fs := entity.GetTableFields(&entity.BaseEntity{})
	if prefix != "" {
		for i := range fs {
			fs[i].Name = prefix + fs[i].Name
		}
	}
	return fs
	//如果有前缀表示是历史表的基础字段，ID就不能是自增的逐渐
	//fs = append(fs, entity.TableField{Name: prefix + entity.IdDbName, Type: ctype.TypeStringName})
	//fs = append(fs, entity.TableField{Name: prefix + entity.CreatedAtDbName, Type: ctype.TypeTimeName})
	//fs = append(fs, entity.TableField{Name: prefix + entity.UpdatedAtDbName, Type: ctype.TpBigInteger})
	//fs = append(fs, entity.TableField{Name: prefix + entity.CreatorIdDbName, Type: ctype.TypeStringName})
	//fs = append(fs, entity.TableField{Name: prefix + entity.UpdaterIdDbName, Type: ctype.TypeStringName})
	//fs = append(fs, entity.TableField{Name: prefix + entity.CreatorDeptIdDbName, Type: ctype.TypeStringName})
	//fs = append(fs, entity.TableField{Name: prefix + entity.UpdaterDeptIdDbName, Type: ctype.TypeStringName})
	//fs = append(fs, entity.TableField{Name: prefix + "deleted", Type: ctype.TypeBoolName})
	//fs = append(fs, entity.TableField{Name: prefix + entity.FieldUpdaterDbName, Type: ctype.TypeJsonName})
	//fs = append(fs, entity.TableField{Name: prefix + entity.FieldUpdaterIdsDbName, Type: ctype.TypeArrayStringName})
	//fs = append(fs, entity.TableField{Name: prefix + entity.ExtraDbName, Type: ctype.TypeJsonName})
	//fs = append(fs, entity.TableField{Name: prefix + entity.EtlExtraDbName, Type: ctype.TypeJsonName})
	//fs = append(fs, entity.TableField{Name: prefix + entity.RemarkDbName, Type: ctype.TypeStringName})
	//fs = append(fs, entity.TableField{Name: prefix + entity.RaContentDbName, Type: ctype.TpText})
	//fs = append(fs, entity.TableField{Name: prefix + entity.RaContentDbName, Type: ctype.TpText})
	//return fs
}
func getBaseHistoryField(db *gorm.DB) []entity.TableField {
	var (
		fs     = make([]entity.TableField, 5)
		idType = ctype.TypeIdName
	)
	fs[0] = entity.TableField{Name: entity.IdDbName, Type: idType}
	fs[1] = entity.TableField{Name: entity.CreatedAtDbName, Type: ctype.TypeTimeName}
	fs[2] = entity.TableField{Name: "op_type", Type: ctype.TypeIntegerName}
	fs[3] = entity.TableField{Name: "opt_user_id", Type: ctype.TypeStringName}
	fs[4] = entity.TableField{Name: "opt_dept_id", Type: ctype.TypeStringName}
	return fs
}

func buildIndexSql(db *gorm.DB, tableName string, idx entity.TableIndex) (string, error) {
	return idx.Build(db, tableName)
}

// 判断切片 src 中的元素是否都在切片 target 中
func isSubset(src, target []string) bool {
	// 遍历 src 的每个元素
	for _, srcElem := range src {
		found := false
		// 查看 src 的元素是否在 target 中
		for _, targetElem := range target {
			if srcElem == targetElem {
				found = true
				break
			}
		}
		// 如果 src 中的元素不在 target 中，返回 false
		if !found {
			return false
		}
	}
	// 所有的元素都在 target 中，返回 true
	return true
}

func buildCommentSql(db *gorm.DB, tableName, column, comment string) string {
	var (
		dbType = dorm.GetDbType(db)
	)
	if dbType == dorm.Mysql {
		return ""
	}
	return fmt.Sprintf(` COMMENT ON COLUMN %s IS '%s'`, dorm.Quote(dbType, tableName, column), comment)
}

// validateIndex 判断索引字段是否都在表中
func validateIndex(p *template.Template) error {
	fieldNames := make([]string, 0, len(p.Table.GetFields()))
	for _, v := range p.Table.GetFields() {
		fieldNames = append(fieldNames, v.Name)
	}
	allIndexCols := make([]string, 0, 10)
	for _, v := range p.Table.Indexes {
		allIndexCols = append(allIndexCols, v.Fields...)
	}
	if !isSubset(allIndexCols, fieldNames) {
		return errors.New("index fields not in table fields")
	}
	return nil
}

func getAllHistoryTableData(db *gorm.DB, scm, tableID string) ([]template.History, error) {
	var (
		rs = make([]template.History, 0, 10)
	)
	tx := db.Raw(fmt.Sprintf(`select * from "%s"."%s" where "h_id" = '%s'`, scm, constants.TableTemplateHistory, tableID)).Find(&rs)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return rs, nil
}
