package builder_test

import (
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/filter"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/core/builder"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestCteBuilder(t *testing.T) {

	bd := builder.NewCteSqlBuilder(dorm.Mysql, "eta_dev_backend", constants.TablePeople)

	bd.AddColumn("id").
		AddFilter(filter.Filter{Column: entity.IdDbName, Operator: filter.Eq, Value: "1077862"}).
		Join("eta_dev_backend", constants.TablePep2RoomLive, entity.FromIdDbName, constants.TablePeople, entity.IdDbName)

	listSql, countSql, err := bd.Build()
	assert.Nil(t, err)
	logger.Info(listSql)
	logger.Info(countSql)
	assert.Equal(t, "select count( *) from `eta_dev_backend`.`t_people` join `delta_dev_backend`.`r_pep2room_live` on `r_pep2room_live`.`from_id` = `t_people`.`id` where  ( `t_people`.`id`  = '1077862')", strings.TrimSpace(countSql))
	assert.Equal(t, "select  `t_people`.`id` from `delta_dev_backend`.`t_people` join `delta_dev_backend`.`r_pep2room_live` on `r_pep2room_live`.`from_id` = `t_people`.`id` where  ( `t_people`.`id`  = '1077862')", strings.TrimSpace(listSql))
}
func TestCteBuilder2(t *testing.T) {

	bd := builder.NewCteSqlBuilder(dbType, scm, constants.TablePeople)

	bd1 := builder.NewSqlBuilder(dorm.DaMeng, scm, constants.TableBd2Addr).
		AddColumn("from_id").
		AddFilter(filter.Filter{Column: "level4_id", Operator: filter.Eq, Value: "46"})
	bd2 := builder.NewSqlBuilder(dorm.DaMeng, scm, constants.TablePep2RoomLive).
		AddColumn("from_id").
		AddFilter(filter.Filter{Column: "live_type", Operator: filter.Eq, Value: "租住"}).
		Join("", "bd2addr", "from_id", constants.TablePep2RoomLive, "building_id")

	bd.With("bd2addr", bd1).
		With("live", bd2).
		Join("", "live", "from_id", constants.TablePeople, "id")

	listSql, countSql, err := bd.Build()
	assert.Nil(t, err)
	logger.Info(listSql)
	logger.Info(countSql)
	assert.Equal(t, `with "bd2addr" as (select  "r_bd2addr"."from_id" from "delta_dev_backend"."r_bd2addr"  where  ( "r_bd2addr"."level4_id"  = '46')   ),"live" as (select  "r_pep2room_live"."from_id" from "delta_dev_backend"."r_pep2room_live" join "bd2addr" on "bd2addr"."from_id" = "r_pep2room_live"."building_id" where  ( "r_pep2room_live"."live_type"  = '租住')   ) select  "t_people".* from "delta_dev_backend"."t_people" join "live" on "live"."from_id" = "t_people"."id"`, strings.TrimSpace(listSql))
	assert.Equal(t, `with "bd2addr" as (select  "r_bd2addr"."from_id" from "delta_dev_backend"."r_bd2addr"  where  ( "r_bd2addr"."level4_id"  = '46')   ),"live" as (select  "r_pep2room_live"."from_id" from "delta_dev_backend"."r_pep2room_live" join "bd2addr" on "bd2addr"."from_id" = "r_pep2room_live"."building_id" where  ( "r_pep2room_live"."live_type"  = '租住')   ) select count( *) from "delta_dev_backend"."t_people" join "live" on "live"."from_id" = "t_people"."id"`, strings.TrimSpace(countSql))
}
func TestCteBuilder3(t *testing.T) {

	bd := builder.NewCteSqlBuilder(dbType, "", constants.TablePeople)

	bd1 := builder.NewCteSqlBuilder(dbType, scm, constants.TablePeople)

	bd11 := builder.NewSqlBuilder(dorm.DaMeng, scm, constants.TableBd2Addr).
		AddColumn("from_id").
		AddFilter(filter.Filter{Column: "level4_id", Operator: filter.Eq, Value: "46"}).SetDistinct(true)
	bd12 := builder.NewSqlBuilder(dorm.DaMeng, scm, constants.TablePep2RoomLive).
		AddColumn("from_id").
		AddFilter(filter.Filter{Column: "live_type", Operator: filter.Eq, Value: "租住"}).
		Join("", "bd2addr", "from_id", constants.TablePep2RoomLive, "building_id").SetDistinct(true)

	bd1.With("bd2addr", bd11).
		With("live", bd12).
		Join("", "live", "from_id", constants.TablePeople, "id")

	bd.With("r", bd1).AddColumn("*")
	bd.SetTableName("r")

	listSql, countSql, err := bd.Build()
	assert.Nil(t, err)
	logger.Info(listSql)
	logger.Info(countSql)
	assert.Equal(t, `with "r" as (with "bd2addr" as (select distinct "r_bd2addr"."from_id" from "delta_dev_backend"."r_bd2addr"  where  ( "r_bd2addr"."level4_id"  = '46')   ),"live" as (select distinct "r_pep2room_live"."from_id" from "delta_dev_backend"."r_pep2room_live" join "bd2addr" on "bd2addr"."from_id" = "r_pep2room_live"."building_id" where  ( "r_pep2room_live"."live_type"  = '租住')   ) select  "t_people".* from "delta_dev_backend"."t_people" join "live" on "live"."from_id" = "t_people"."id"   ) select  "r".* from "r"`, strings.TrimSpace(listSql))
	assert.Equal(t, `with "r" as (with "bd2addr" as (select distinct "r_bd2addr"."from_id" from "delta_dev_backend"."r_bd2addr"  where  ( "r_bd2addr"."level4_id"  = '46')   ),"live" as (select distinct "r_pep2room_live"."from_id" from "delta_dev_backend"."r_pep2room_live" join "bd2addr" on "bd2addr"."from_id" = "r_pep2room_live"."building_id" where  ( "r_pep2room_live"."live_type"  = '租住')   ) select  "t_people".* from "delta_dev_backend"."t_people" join "live" on "live"."from_id" = "t_people"."id"   ) select count( *) from "r"`, strings.TrimSpace(countSql))
}
func TestBuilderSetTableName(t *testing.T) {

	bd := builder.NewCteSqlBuilder(dbType, "", constants.TablePeople)

	bd1 := builder.NewSqlBuilder(dbType, scm, constants.TableBd2Addr).
		AddColumn("from_id").SetDistinct(true).
		AddFilter(filter.Filter{Column: "level4_id", Operator: filter.Eq, Value: "46"})
	bd.With("bd2addr", bd1)

	bd2 := builder.NewSqlBuilder(dbType, scm, constants.TableShop2Bd).
		AddColumn("from_id").SetDistinct(true).
		Join("", "bd2addr", "from_id", constants.TableShop2Bd, "to_id")
	bd.With("shop2bd", bd2)

	bd3 := builder.NewSqlBuilder(dbType, scm, constants.TableShop).
		AddColumn("id").SetDistinct(true).
		Join("", "shop2bd", "from_id", constants.TableShop, "id")
	bd.With("shop", bd3)

	bd4 := builder.NewSqlBuilder(dbType, scm, constants.TablePep2Shop).
		AddColumn("from_id").SetDistinct(true).
		Join("", "shop", "id", constants.TablePep2Shop, "to_id")

	bd.With("pep2shop", bd4)
	bd.SetTableName("pep2shop").AddColumn("from_id")

	listSql, countSql, err := bd.Build()
	assert.Nil(t, err)
	logger.Info(listSql)
	logger.Info(countSql)
	assert.Equal(t, `with "bd2addr" as (select distinct "r_bd2addr"."from_id" from "delta_dev_backend"."r_bd2addr"  where  ( "r_bd2addr"."level4_id"  = '46')   ),"shop2bd" as (select distinct "r_shop2bd"."from_id" from "delta_dev_backend"."r_shop2bd" join "bd2addr" on "bd2addr"."from_id" = "r_shop2bd"."to_id"   ),"shop" as (select distinct "t_shop"."id" from "delta_dev_backend"."t_shop" join "shop2bd" on "shop2bd"."from_id" = "t_shop"."id"   ),"pep2shop" as (select distinct "r_pep2shop"."from_id" from "delta_dev_backend"."r_pep2shop" join "shop" on "shop"."id" = "r_pep2shop"."to_id"   ) select  "pep2shop"."from_id" from "pep2shop"`, strings.TrimSpace(listSql))
	assert.Equal(t, `with "bd2addr" as (select distinct "r_bd2addr"."from_id" from "delta_dev_backend"."r_bd2addr"  where  ( "r_bd2addr"."level4_id"  = '46')   ),"shop2bd" as (select distinct "r_shop2bd"."from_id" from "delta_dev_backend"."r_shop2bd" join "bd2addr" on "bd2addr"."from_id" = "r_shop2bd"."to_id"   ),"shop" as (select distinct "t_shop"."id" from "delta_dev_backend"."t_shop" join "shop2bd" on "shop2bd"."from_id" = "t_shop"."id"   ),"pep2shop" as (select distinct "r_pep2shop"."from_id" from "delta_dev_backend"."r_pep2shop" join "shop" on "shop"."id" = "r_pep2shop"."to_id"   ) select count( *) from "pep2shop"`, strings.TrimSpace(countSql))
}

func TestUnion(t *testing.T) {

	var (
		flt = filter.Filter{Column: "building_id", Operator: filter.Eq, Value: "4645"}
	)

	bd := builder.NewSqlBuilder(dorm.DaMeng, scm, constants.TablePep2RoomLive).AddColumn("from_id").AddFilter(flt)
	bd1 := builder.NewSqlBuilder(dorm.DaMeng, scm, constants.TablePep2RoomHuJi).AddColumn("from_id").AddFilter(flt)
	bd2 := builder.NewSqlBuilder(dorm.DaMeng, scm, constants.TablePep2EntJob).AddColumn("from_id")
	bd3 := builder.NewSqlBuilder(dorm.DaMeng, scm, constants.TablePep2Label).AddColumn("from_id")

	listSql, countSql, err := bd.Union(bd1).UnionIntersect(bd2).UnionExcept(bd3).Build()
	assert.Nil(t, err)
	logger.Info(listSql)
	logger.Info(countSql)
	assert.Equal(t, `select  "r_pep2room_live"."from_id" from "delta_dev_backend"."r_pep2room_live"  where  ( "r_pep2room_live"."building_id"  = '4645')    union (select  "r_pep2room_huji"."from_id" from "delta_dev_backend"."r_pep2room_huji"  where  ( "r_pep2room_huji"."building_id"  = '4645')   ) intersect (select  "r_pep2ent_job"."from_id" from "delta_dev_backend"."r_pep2ent_job"    ) except (select  "r_pep2label"."from_id" from "delta_dev_backend"."r_pep2label"    )`, listSql)
	assert.Equal(t, `with "r" as (select  "r_pep2room_live"."from_id" from "delta_dev_backend"."r_pep2room_live"  where  ( "r_pep2room_live"."building_id"  = '4645')    union (select  "r_pep2room_huji"."from_id" from "delta_dev_backend"."r_pep2room_huji"  where  ( "r_pep2room_huji"."building_id"  = '4645')   ) intersect (select  "r_pep2ent_job"."from_id" from "delta_dev_backend"."r_pep2ent_job"    ) except (select  "r_pep2label"."from_id" from "delta_dev_backend"."r_pep2label"    )) select count(*) from "r"`, countSql)
}
