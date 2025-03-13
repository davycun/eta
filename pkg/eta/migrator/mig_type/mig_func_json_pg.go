package mig_type

import (
	"gorm.io/gorm"
)

func catArray(db *gorm.DB) error {
	sql := `create or replace function array_distinct_cat(a1 anycompatiblearray, a2 anycompatiblearray) returns anycompatiblearray
			language plpgsql as
		$$
		BEGIN
			return ARRAY(select distinct unnest(array_cat(a1,a2)));
		end;
		$$;`
	return db.Transaction(func(tx *gorm.DB) error {
		return tx.Exec(sql).Error
	})
}
