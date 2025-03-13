package mysql

import (
	"fmt"
	"gorm.io/gorm"
)

func CreateIndexIfNotExists(db *gorm.DB, tableName, indexName, createIndexSql string) error {
	sql := fmt.Sprintf(`SELECT COUNT(*) FROM information_schema.statistics 
                WHERE TABLE_SCHEMA = DATABASE() 
                  AND TABLE_NAME = '%s' 
                  AND INDEX_NAME = '%s'`, tableName, indexName)

	count := 0
	err := db.Raw(sql).Scan(&count).Error
	if err != nil {
		return err
	}
	if count <= 0 {
		return db.Exec(createIndexSql).Error
	}
	return nil
}
