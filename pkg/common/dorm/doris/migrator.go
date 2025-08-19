package doris

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/tag"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/migrator"
	"gorm.io/gorm/schema"
	"strings"
)

type Migrator struct {
	mysql.Migrator
	Dialector
}

func (m Migrator) CurrentDatabase() (name string) {
	var (
		scm = dorm.GetDbSchema(m.DB)
	)
	return scm
}

// HasTable returns table exists or not for value, value could be a struct or string
func (m Migrator) HasTable(value interface{}) bool {
	var count int64

	err := m.RunWithValue(value, func(stmt *gorm.Statement) error {
		var (
			scm = dorm.GetDbSchema(m.DB)
		)
		return dorm.RawFetch(fmt.Sprintf("select count(*) from `information_schema`.`tables` where `TABLE_SCHEMA`='%s' and table_type='%s' and table_name='%s'", scm, "BASE TABLE", stmt.Table), m.DB, &count)
	})
	if err != nil {
		logger.Errorf("error while checking table existence: %v", err)
	}

	return count > 0
}

// AutoMigrate
// TODO 已经修改的表暂不支持修改
func (m Migrator) AutoMigrate(values ...interface{}) error {
	for _, value := range m.ReorderModels(values, true) {
		queryTx, execTx := m.GetQueryAndExecTx()
		if !queryTx.Migrator().HasTable(value) {
			if err := execTx.Migrator().CreateTable(value); err != nil {
				return err
			}
		}
	}

	return nil
}
func (m Migrator) FullDataTypeOf(field *schema.Field) clause.Expr {

	var (
		tg      = tag.ParseDorisTag(field.Tag.Get(tag.DorisTagName))
		expr    = clause.Expr{}
		aggType = tg.Get("agg_type")
	)

	expr.SQL = m.Migrator.Migrator.DataTypeOf(field)
	if aggType != "" {
		expr.SQL += " " + aggType
	}

	if field.NotNull {
		expr.SQL += " NOT NULL"
	}

	if field.HasDefaultValue && (field.DefaultValueInterface != nil || field.DefaultValue != "") {
		if field.DefaultValueInterface != nil {
			defaultStmt := &gorm.Statement{Vars: []interface{}{field.DefaultValueInterface}}
			m.Dialector.BindVarTo(defaultStmt, defaultStmt, field.DefaultValueInterface)
			expr.SQL += " DEFAULT " + m.Dialector.Explain(defaultStmt.SQL.String(), field.DefaultValueInterface)
		} else if field.DefaultValue != "(-)" {
			expr.SQL += " DEFAULT " + field.DefaultValue
		}
	}

	if value, ok := field.TagSettings["COMMENT"]; ok {
		expr.SQL += " COMMENT " + m.Dialector.Explain("?", value)
	}

	return expr
}

func (m Migrator) CreateTable(values ...interface{}) error {
	for _, value := range m.ReorderModels(values, false) {
		tx := m.DB.Session(&gorm.Session{})
		if err := m.RunWithValue(value, func(stmt *gorm.Statement) (err error) {
			var (
				createTableSQL = "CREATE TABLE ? ("
				values         = []interface{}{m.CurrentTable(stmt)}
				//hasPrimaryKeyInDataType bool
			)

			for _, dbName := range stmt.Schema.DBNames {
				field := stmt.Schema.FieldsByDBName[dbName]
				if !field.IgnoreMigration {
					createTableSQL += "? ?"
					//hasPrimaryKeyInDataType = hasPrimaryKeyInDataType || strings.Contains(strings.ToUpper(m.DataTypeOf(field)), "PRIMARY KEY")
					values = append(values, clause.Column{Name: dbName}, m.DB.Migrator().FullDataTypeOf(field))
					createTableSQL += ","
				}
			}

			//if !hasPrimaryKeyInDataType && len(stmt.Schema.PrimaryFields) > 0 {
			//	createTableSQL += "PRIMARY KEY ?,"
			//	primaryKeys := make([]interface{}, 0, len(stmt.Schema.PrimaryFields))
			//	for _, field := range stmt.Schema.PrimaryFields {
			//		primaryKeys = append(primaryKeys, clause.Column{Name: field.DBName})
			//	}
			//
			//	values = append(values, primaryKeys)
			//}

			for _, idx := range stmt.Schema.ParseIndexes() {
				if m.CreateIndexAfterCreateTable {
					defer func(value interface{}, name string) {
						if err == nil {
							err = tx.Migrator().CreateIndex(value, name)
						}
					}(value, idx.Name)
				} else {
					if idx.Class != "" {
						createTableSQL += idx.Class + " "
					}
					createTableSQL += "INDEX ? ?"

					if idx.Comment != "" {
						createTableSQL += fmt.Sprintf(" COMMENT '%s'", idx.Comment)
					}

					if idx.Option != "" {
						createTableSQL += " " + idx.Option
					}

					createTableSQL += ","
					values = append(values, clause.Column{Name: idx.Name}, tx.Migrator().(migrator.BuildIndexOptionsInterface).BuildIndexOptions(idx.Fields, stmt))
				}
			}

			if !m.DB.DisableForeignKeyConstraintWhenMigrating && !m.DB.IgnoreRelationshipsWhenMigrating {
				for _, rel := range stmt.Schema.Relationships.Relations {
					if rel.Field.IgnoreMigration {
						continue
					}
					if constraint := rel.ParseConstraint(); constraint != nil {
						if constraint.Schema == stmt.Schema {
							sql, vars := constraint.Build()
							createTableSQL += sql + ","
							values = append(values, vars...)
						}
					}
				}
			}

			for _, uni := range stmt.Schema.ParseUniqueConstraints() {
				createTableSQL += "CONSTRAINT ? UNIQUE (?),"
				values = append(values, clause.Column{Name: uni.Name}, clause.Expr{SQL: stmt.Quote(uni.Field.DBName)})
			}

			for _, chk := range stmt.Schema.ParseCheckConstraints() {
				createTableSQL += "CONSTRAINT ? CHECK (?),"
				values = append(values, clause.Column{Name: chk.Name}, clause.Expr{SQL: chk.Constraint})
			}

			createTableSQL = strings.TrimSuffix(createTableSQL, ",")

			createTableSQL += ")"

			if tableOption, ok := m.DB.Get("gorm:table_options"); ok {
				createTableSQL += fmt.Sprint(tableOption)
			}

			err = tx.Exec(createTableSQL, values...).Error
			return err
		}); err != nil {
			return err
		}
	}
	return nil
}
