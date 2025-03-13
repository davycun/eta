package db_table

import (
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/logger"
	"gorm.io/gorm"
)

// TableExists
// TODO 需要做缓存，但是需要注意缓存的key的定义，应该是需要 Database.Key
func TableExists(db *gorm.DB, scm, tableName string) bool {

	var (
		err    error
		count  int64
		dbType = dorm.GetDbType(db)
	)

	switch dbType {
	case dorm.PostgreSQL, dorm.Mysql:
		err = db.Raw("SELECT count(*) FROM information_schema.tables WHERE table_schema = ? AND table_name = ? AND table_type = ?", scm, tableName, "BASE TABLE").Scan(&count).Error
	case dorm.DaMeng:
		tableSql := `SELECT /*+ MAX_OPT_N_TABLES(5) */ COUNT(TABS.NAME) FROM
				(SELECT ID, PID FROM SYS.SYSOBJECTS WHERE TYPE$ = 'SCH' AND NAME = ?) SCHEMAS,
				(SELECT ID, SCHID, NAME FROM SYS.SYSOBJECTS WHERE
				NAME = ? AND TYPE$ = 'SCHOBJ' AND SUBTYPE$ IN ('UTAB', 'STAB', 'VIEW', 'SYNOM')
				AND ((SUBTYPE$ ='UTAB' AND CAST((INFO3 & 0x00FF & 0x003F) AS INT) not in (9, 27, 29, 25, 12, 7, 21, 23, 18, 5))
				OR SUBTYPE$ in ('STAB', 'VIEW', 'SYNOM'))) TABS
				WHERE TABS.SCHID = SCHEMAS.ID AND SF_CHECK_PRIV_OPT(UID(), CURRENT_USERTYPE(), TABS.ID, SCHEMAS.PID, -1, TABS.ID) = 1;`

		err = db.Raw(tableSql, scm, tableName).Row().Scan(&count)
	}

	if err != nil {
		logger.Errorf("tableExists err %s", err)
	}

	return count > 0
}
