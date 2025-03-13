package dorm_test

import (
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateIndex(t *testing.T) {
	db, execSql, _, err := dorm.NewTestDB(dorm.DaMeng, "eta")
	assert.Nil(t, err)

	err = dorm.CreateIndex(db, "my_test", "id_no")
	assert.Nil(t, err)
	assert.True(t, execSql.Exists(`CREATE  INDEX IF NOT EXISTS "idx_my_test_id_no" ON "eta"."my_test"("id_no")`))

	err = dorm.CreateIndex(db, "my_user", "app_id", "mobile")
	assert.Nil(t, err)
	assert.True(t, execSql.Exists(`CREATE  INDEX IF NOT EXISTS "idx_my_user_app_id_mobile" ON "eta"."my_user"("app_id","mobile")`))
	execSql.Reset()

	err = dorm.CreateContextIndex(db, "my_test", "content")
	assert.Nil(t, err)
	assert.True(t, execSql.Exists(`CREATE CONTEXT INDEX IF NOT EXISTS "idx_my_test_content" ON "eta"."my_test"("content") SYNC TRANSACTION`))

	err = dorm.CreateArrayIndex(db, "my_test", "addr_ids")
	assert.Nil(t, err)
	assert.True(t, execSql.Exists(`CREATE ARRAY INDEX IF NOT EXISTS "idx_my_test_addr_ids" ON "eta"."my_test"("addr_ids")`))

	db, execSql, _, err = dorm.NewTestDB(dorm.Mysql, "eta")
	assert.Nil(t, err)
	err = dorm.CreateIndex(db, "my_test", "name(256)", "post(128)")
	assert.Nil(t, err)
	assert.True(t, execSql.Exists("CREATE  INDEX `idx_my_test_name_post` ON `eta`.`my_test`(`name`(256),`post`(128))"))

	err = dorm.CreateUniqueIndex(db, "my_test", "name(256)", "post(128)")
	assert.Nil(t, err)
	assert.True(t, execSql.Exists("CREATE UNIQUE INDEX `idx_my_test_name_post` ON `eta`.`my_test`(`name`(256),`post`(128))"))
}
