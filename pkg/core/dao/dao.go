package dao

import (
	"gorm.io/gorm"
)

func FetchById(id string, db *gorm.DB, data any, columns ...string) error {
	tx := db.Model(data)
	if len(columns) > 0 {
		tx = tx.Select(columns)
	}
	return tx.Where(map[string]any{"id": id}).First(data).Error
}
