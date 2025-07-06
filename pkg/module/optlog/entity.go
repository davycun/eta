package optlog

import (
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/eta/constants"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"strings"
)

const (
	AutoCollect   = "自动采集" //运用到OptCategory
	ManualCollect = "手动采集"
)

type OptLog struct {
	entity.BaseEntity
	ReqId         string          `json:"req_id,omitempty" gorm:"column:req_id;comment:请求ID" es:"type:text"`
	ReqUri        string          `json:"req_uri,omitempty" gorm:"column:req_uri;comment:请求地址" es:"type:text"`
	OptUserId     string          `json:"opt_user_id,omitempty" gorm:"column:opt_user_id;comment:操作用户ID" es:"type:text"`        //操作用户ID
	OptUserKey    string          `json:"opt_user_key,omitempty" gorm:"column:opt_user_key;comment:操作用户的key" es:"type:text"`    //通过key方式访问的时候进行记录
	OptDeptId     string          `json:"opt_dept_id,omitempty" gorm:"column:opt_dept_id;comment:操作用户所属部门ID" es:"type:text"`    //操作用户所属部门ID
	OptCategory   string          `json:"opt_category,omitempty" gorm:"column:opt_category;comment:日志来源分类手动或自动" es:"type:text"` //自动采集，手动采集
	OptType       string          `json:"opt_type,omitempty" gorm:"column:opt_type;comment:操作类型" es:"type:text"`                //操作类型，登录、退出、数据查询、数据打标、打标、导出
	OptContent    string          `json:"opt_content,omitempty" gorm:"column:opt_content;comment:操作内容" es:"type:text"`          //操作内容
	OptTarget     string          `json:"opt_target,omitempty" gorm:"column:opt_target;comment:操作目标" es:"type:text"`            //操作目标， 居民信息、房间信息
	OptTime       ctype.LocalTime `json:"opt_time,omitempty" gorm:"column:opt_time;comment:操作时间" es:"type:date"`                //操作时间
	ClientIp      string          `json:"client_ip,omitempty" gorm:"column:client_ip;comment:客户端IP" es:"type:text"`             //客户端IP
	ClientType    string          `json:"client_type,omitempty" gorm:"column:client_type;comment:客户端类型" es:"type:text"`         //PC、Mobile、Service
	ClientTrigger string          `json:"client_trigger,omitempty" gorm:"column:client_trigger;comment:客户端触发点" es:"type:text"`  //客户端的触发点（比如页面上的菜单按钮的名称)
	RsStatus      string          `json:"rs_status,omitempty" gorm:"column:rs_status;comment:操作结果状态" es:"type:text"`            //调用返回的http状态，
	RsRemark      string          `json:"rs_remark,omitempty" gorm:"column:rs_remark;comment:操作结果说明" es:"type:text"`            //调用返回的说明，脱敏/非脱敏，加密或者非加密
	Latency       int64           `json:"latency,omitempty" gorm:"column:latency;comment:响应时长毫秒" es:"type:text"`
}

func (o OptLog) TableName(namer schema.Namer) string {

	if namer != nil {
		return namer.TableName(constants.TableOperateLog)
	}

	return constants.TableOperateLog
}
func (o OptLog) EsIndexName() string {
	return o.TableName(nil)
}
func (o OptLog) RaDbFields() []string {
	return []string{
		"req_id",
		"req_uri",
		"opt_user_id",
		"opt_dept_id",
		"opt_category",
		"opt_type",
		"opt_content",
		"opt_target",
		"opt_time",
		"client_ip",
		"client_type",
		"client_trigger",
		"rs_status",
		"rs_remark",
		"latency",
	}
}

func (o *OptLog) BeforeCreate(db *gorm.DB) error {

	err := o.BaseEntity.BeforeCreate(db)
	if o.OptUserId == "" {
		c, b := ctx.GetCurrentContext()
		if b && c != nil {
			o.OptUserId = c.GetContextUserId()
			o.OptDeptId = c.GetContextCurrentDeptId()
		}
	}
	//设置为手动采集
	if o.OptCategory == "" {
		o.OptCategory = ManualCollect
	}
	if strings.Contains(o.OptType, "?") {
		o.OptType = strings.Split(o.OptType, "?")[0]
	}
	if strings.Contains(o.OptTarget, "?") {
		o.OptTarget = strings.Split(o.OptTarget, "?")[0]
	}

	return err
}
