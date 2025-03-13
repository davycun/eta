package dorm_test

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	dameng "github.com/davycun/dm8-gorm"
	"github.com/davycun/eta/pkg/common/dorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
	"time"
)

type User struct {
	ID           uint           `gorm:"primaryKey"` // Standard field for the primary key
	Name         string         // 一个常规字符串字段
	Email        *string        // 一个指向字符串的指针, allowing for null values
	Age          uint8          // 一个未签名的8位整数
	Birthday     *time.Time     // A pointer to time.Time, can be null
	MemberNumber sql.NullString // Uses sql.NullString to handle nullable strings
	ActivatedAt  sql.NullTime   // Uses sql.NullTime for nullable time fields
	CreatedAt    time.Time      // 创建时间（由GORM自动管理）
	UpdatedAt    time.Time      // 最后一次更新时间（由GORM自动管理）
}

func TestBatchUpdate(t *testing.T) {
	var (
		cols  = []string{"name", "email"}
		name  = "jinzhu"
		email = "jinzhu@qq.com"
		id    = 1
	)
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()
	// pg
	pgDialector := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})
	pgDB, _ := gorm.Open(pgDialector, &gorm.Config{})
	// Test case 1: pg single
	mock.ExpectExec(`^UPDATE "users" set "name"="excluded"\."name","email"="excluded"\."email","id"="excluded"\."id" from \(VALUES \(.+\)\) AS "excluded" \("name","email","created_at","updated_at","id"\) where "users"\."id" = "excluded"\."id"$`).WillReturnResult(sqlmock.NewResult(1, 1))
	err := dorm.BatchUpdate(pgDB, &User{
		ID:    uint(id),
		Name:  name,
		Email: &email,
	}, cols...)
	if err != nil {
		t.Errorf("Error occurred: %v", err)
	}
	// Test case 2: pg multiple
	mock.ExpectExec(`^UPDATE "users" set "name"="excluded"\."name","email"="excluded"\."email","id"="excluded"\."id" from \(VALUES \(.+\), \(.+\)\) AS "excluded" \("name","email","created_at","updated_at","id"\) where "users"\."id" = "excluded"\."id"$`).WillReturnResult(sqlmock.NewResult(1, 1))
	err1 := dorm.BatchUpdate(pgDB, &[]User{
		{
			ID:    uint(id),
			Name:  name,
			Email: &email,
		},
		{
			ID:    uint(id),
			Name:  name,
			Email: &email,
		},
	}, cols...)
	if err1 != nil {
		t.Errorf("Error occurred: %v", err1)
	}

	//dm
	dmDialector := dameng.New(dameng.Config{
		Conn:       mockDB,
		DriverName: "dm",
	})
	dmDB, _ := gorm.Open(dmDialector, &gorm.Config{})

	// Test case 3: pg single
	mock.ExpectExec(`^UPDATE "users" set "name"="excluded"\."name","email"="excluded"\."email","id"="excluded"\."id" from \(SELECT .+ FROM DUAL\) AS "excluded" \("name","email","created_at","updated_at","id"\) where "users"\."id" = "excluded"\."id"$`).WillReturnResult(sqlmock.NewResult(1, 1))
	err2 := dorm.BatchUpdate(dmDB, &User{
		ID:    uint(id),
		Name:  name,
		Email: &email,
	}, cols...)
	if err2 != nil {
		t.Errorf("Error occurred: %v", err2)
	}
	// Test case 4: pg multiple
	mock.ExpectExec(`^UPDATE "users" set "name"="excluded"\."name","email"="excluded"\."email","id"="excluded"\."id" from \(SELECT .+ FROM DUAL UNION ALL SELECT .+ FROM DUAL\) AS "excluded" \("name","email","created_at","updated_at","id"\) where "users"\."id" = "excluded"\."id"$`).WillReturnResult(sqlmock.NewResult(1, 1))
	err3 := dorm.BatchUpdate(dmDB, &[]User{
		{
			ID:    uint(id),
			Name:  name,
			Email: &email,
		},
		{
			ID:    uint(id),
			Name:  name,
			Email: &email,
		},
	}, cols...)
	if err3 != nil {
		t.Errorf("Error occurred: %v", err3)
	}
}
