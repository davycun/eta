package mig_type

import (
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/logger"
	"gorm.io/gorm"
)

func initDm(db *gorm.DB) error {
	var (
		err           error
		geoFlag       int
		geo2Flag      int
		alterTableOpt int
	)
	err = caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			//这个其实是需要重启数据库的,已经在镜像中支持配置dm.ini 参数来解决
			err = db.Exec(`alter system set 'COMPATIBLE_MODE'=7 spfile;`).Error
			if err != nil {
				logger.Errorf("set COMPATIBLE_MODE=7 error: %s", err)
			}
			return nil
		}).
		Call(func(cl *caller.Caller) error {
			return db.Raw(`SELECT SF_CHECK_GEO_SYS()`).Find(&geoFlag).Error
		}).
		Call(func(cl *caller.Caller) error {
			return db.Raw(`SELECT SF_CHECK_GEO2_SYS()`).Find(&geo2Flag).Error
		}).
		Call(func(cl *caller.Caller) error {
			if geoFlag == 0 {
				return db.Exec(`SP_INIT_GEO_SYS(1)`).Error
			}
			return nil
		}).
		Call(func(cl *caller.Caller) error {
			if geo2Flag == 0 {
				return db.Exec(`SP_INIT_GEO2_SYS(1)`).Error
			}
			return nil
		}).
		Call(func(cl *caller.Caller) error {
			return db.Raw(`SELECT PARA_VALUE FROM V$DM_INI WHERE PARA_NAME='ALTER_TABLE_OPT'`).Find(&alterTableOpt).Error
		}).
		Call(func(cl *caller.Caller) error {
			if alterTableOpt != 3 {
				return db.Exec(`ALTER SYSTEM SET 'ALTER_TABLE_OPT'=3 BOTH`).Error
			}
			return nil
		}).Err

	return err
}

func arrContainsStr(db *gorm.DB) error {
	dbUser := dorm.GetDbUser(db)
	sq := `CREATE OR REPLACE FUNCTION ` + dbUser + `.arr_contains_str(arr_cls ` + dbUser + `.ARR_STR_CLS,target varchar,regex BIT)  RETURN BIT AS
DECLARE
    item VARCHAR2;
	sz INT;
	arr ` + dbUser + `.ARR_STR;
BEGIN
	if arr_cls is null then
		return 0;
	end if;
	arr := arr_cls.V;

	sz := arr.COUNT();
	for i in 1..sz loop
		item := arr(i);
		if regex then
			if item like target then
				return 1;
			end if;
		else 
			if item = target then
				return 1;
			end if;
		end if;
	end loop;
	return 0;
END;`
	return db.Exec(sq).Error
}
func arrContainsInt(db *gorm.DB) error {
	dbUser := dorm.GetDbUser(db)
	sq := `CREATE OR REPLACE FUNCTION ` + dbUser + `.arr_contains_int(arr_cls ` + dbUser + `.ARR_INT_CLS,target INT)  RETURN BIT AS
DECLARE
    item INT;
	sz INT;
	arr ` + dbUser + `.ARR_INT;
BEGIN
	if arr_cls is null then
		return 0;
	end if;
	arr := arr_cls.V;
	sz := arr.COUNT();
	for i in 1..sz loop
		item := arr(i);
		if item = target then
			return 1;
		end if;
	end loop;
	return 0;
END;`
	return db.Exec(sq).Error
}
func arrIntAnyBetween(db *gorm.DB) error {
	dbUser := dorm.GetDbUser(db)
	sq := `CREATE OR REPLACE FUNCTION ` + dbUser + `.arr_int_any_between(cls ` + dbUser + `.ARR_INT_CLS,st INT,ed INT) RETURN BIT AS
DECLARE
    item INT;
    arr ` + dbUser + `.ARR_INT;
	total INT;
BEGIN
    if cls is null then
		return 0;
	end if;
    arr := cls.V;
	total := arr.count();
	if arr.count() < 1 then
		return 0;
	end if;
    for i in 1..total loop 
        item := arr(i);
        if item >= st and item <= ed then
            return 1;
        end if;
    end loop;
	return 0;
end;`
	return db.Exec(sq).Error
}

func jsonSize(db *gorm.DB) error {
	dbUser := dorm.GetDbUser(db)
	sq := `CREATE OR REPLACE FUNCTION ` + dbUser + `.json_size(json_arr varchar2)  RETURN INT AS
DECLARE
    arr jdom_t;
BEGIN
	if json_arr is null or json_arr='' then
		return 0;
	end if;
	arr := jdom_t.parse(json_arr);
	return arr.get_size();
END;`
	return db.Exec(sq).Error
}
func jsonContainsAny(db *gorm.DB) error {
	dbUser := dorm.GetDbUser(db)
	sq := `CREATE OR REPLACE FUNCTION  ` + dbUser + `.json_contains_any(json_str varchar2,fd varchar,target varchar,regex BIT)  RETURN BIT AS
DECLARE
    json_obj JDOM_T;
BEGIN

	if json_str is null or json_str='' then
		return 0;
	end if;
	json_obj := jdom_t.parse(json_str);

	if json_obj.is_array() then
		return "` + dbUser + `".json_array_contains_any(json_str,fd,target,regex);
	elseif json_obj.is_object() then
		return "` + dbUser + `".json_object_contains_any(json_str,fd,target,regex);
	elseif json_obj.is_scalar() then
		return "` + dbUser + `".json_scalar_contains_any(json_str,target,regex);
	end if;
	return 0;
END;`
	return db.Exec(sq).Error
}
func jsonArrayContainsAny(db *gorm.DB) error {
	dbUser := dorm.GetDbUser(db)
	sq := `CREATE OR REPLACE FUNCTION  "` + dbUser + `".json_array_contains_any(json_str varchar2,fd varchar,target varchar,regex BIT)  RETURN BIT AS
DECLARE
    json_obj JDOM_T;
	item JDOM_T;
	sz INT;
	rs BIT;
BEGIN

	if json_str is null or json_str='' then
		return 0;
	end if;
	json_obj := jdom_t.parse(json_str);

	if json_obj.is_array() then
		sz := json_obj.get_size();
		for i in 1..sz loop
			item := json_obj."GET"(i-1);
			if item.is_object() then
				rs := "` + dbUser + `".json_object_contains_any(item.to_string(),fd,target,regex);
				if rs then 
					return 1;
				end if;
			elseif item.is_scalar() then
				rs := "` + dbUser + `".json_scalar_contains_any(item.to_string(),target,regex);
				if rs then 
					return 1;
				end if;
			elseif item.is_array() then
				rs := "` + dbUser + `".json_array_contains_any(item.to_string(),fd,target,regex);
				if rs then 
					return 1;
				end if;
			end if;
		end loop;
	end if;
	return 0;
END;`
	return db.Exec(sq).Error
}
func jsonObjectContainsAny(db *gorm.DB) error {
	dbUser := dorm.GetDbUser(db)
	sq := `CREATE OR REPLACE FUNCTION  "` + dbUser + `".json_object_contains_any(json_str varchar2,fd varchar,target varchar,regex BIT)  RETURN BIT AS
DECLARE
    json_obj JDOM_T;
	item JDOM_T;
	item2 JDOM_T;
	tmp varchar2;
	sz int;
	rs BIT;
BEGIN

	if json_str is null or json_str='' then
		return 0;
	end if;
	json_obj := jdom_t.parse(json_str);

	if json_obj.is_object() then
		if fd = '' or fd is null then
			select json_query(json_str,'$.*' with wrapper) into tmp;
			if tmp = '' or tmp is null then
				return 0;
			end if;
			item := jdom_t.parse(tmp);
			if item.is_array() then
				sz := item.get_size();
				for i in 1..sz loop
					item2 := item."GET"(i-1);
					if item2.is_scalar() then
						rs := "` + dbUser + `".json_scalar_contains_any(item2.to_string(),target,regex);
						if rs then 
							return 1;
						end if;
					end if;
				end loop;
			end if;
		else
			tmp := json_obj.get_string(fd);
			if regex then
				if tmp like target then
					return 1;
				end if;
			else 
				if tmp = target then
					return 1;
				end if;
			end if;
		end if;
	end if;
	return 0;
END;`
	return db.Exec(sq).Error
}
func jsonScalarContainsAny(db *gorm.DB) error {
	dbUser := dorm.GetDbUser(db)
	sq := `CREATE OR REPLACE FUNCTION  "` + dbUser + `".json_scalar_contains_any(json_str varchar2,target varchar,regex BIT)  RETURN BIT AS
DECLARE
    json_obj JDOM_T;
	tmp varchar2;
BEGIN

	if json_str is null or json_str='' then
		return 0;
	end if;
	json_obj := jdom_t.parse(json_str);

	if json_obj.is_Scalar() then
		tmp := json_obj.to_string();
		if regex then
			if tmp like target then
				return 1;
			end if;
		else 
			if tmp = target then
				return 1;
			end if;
		end if;
	end if;
	return 0;
END;`
	return db.Exec(sq).Error
}

func jsonFieldToArrayInt(db *gorm.DB) error {
	scm := dorm.GetDbUser(db)
	sq := `CREATE OR REPLACE FUNCTION "` + scm + `".json_field2arr_int(json_arr varchar2,fd varchar)  RETURN "` + scm + `".ARR_INT AS
DECLARE
    arr JDOM_T;
	item JDOM_T;
	target "` + scm + `".ARR_INT;
	jType varchar;
	sz INT;

BEGIN
	if json_arr is null or json_arr='' or fd='' then
		return target;
	end if;
	arr := jdom_t.parse(json_arr);
	if not arr.is_array() then
		return target;
    end if;

	sz := arr.get_size();
	if sz < 1 then
		return target;
	end if;

	jType := arr.get_type(0);
	
	if jType != 'OBJECT' then
		return target;
	end if;

	target := "` + scm + `".ARR_INT();
	target.extend(sz);

	for i in 1..sz loop
		item := arr."GET"(i-1);
		target(i) := item.get_number(fd);
	end loop;

	return target;
END;`
	return db.Exec(sq).Error
}
func jsonFieldToArrayStr(db *gorm.DB) error {
	scm := dorm.GetDbUser(db)
	sq := `CREATE OR REPLACE FUNCTION "` + scm + `".json_field2arr_str(json_arr varchar2,fd varchar)  RETURN "` + scm + `".ARR_STR_CLS AS
DECLARE
    arr JDOM_T;
	item JDOM_T;
	target "` + scm + `".ARR_STR;
	arr_cls "` + scm + `".ARR_STR_CLS;
	jType varchar;
	sz INT;
BEGIN

	if json_arr is null or json_arr='' or fd='' then
		return  arr_cls;
	end if;
	arr := jdom_t.parse(json_arr);
	if not arr.is_array() then
		return arr_cls;
    end if;

	sz := arr.get_size();
	if sz < 1 then
		return arr_cls;
	end if;

	jType := arr.get_type(0);
	
	if jType != 'OBJECT' then
		return arr_cls;
	end if;

	target := "` + scm + `".ARR_STR();
	target.extend(sz);

	for i in 1..sz loop
		item := arr."GET"(i-1);
		target(i) := item.get_string(fd);
	end loop;

	return "` + scm + `".ARR_STR_CLS(target);

END;`
	return db.Exec(sq).Error
}
