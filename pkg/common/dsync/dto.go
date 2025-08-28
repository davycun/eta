package dsync

import (
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/iface"
)

type SyncOption struct {
	ConsumerGoSize    int  `json:"consumer_go_size,omitempty"`
	ChanSize          int  `json:"chan_size,omitempty"`
	UpdateDbBatchSize int  `json:"update_db_batch_size,omitempty"` // 保存到数据库的批量大小
	Restore           bool `json:"clean,omitempty"`                //是否对相关数据先清空
	Merge             bool `json:"merge,omitempty"`                //同步的时候是否采用Merge的方式，比如同步宽表，如果以前表已经有数据就传入true

	UpdateDbRaContent bool `json:"update_db_ra_content,omitempty"`
	UpdateDbEncrypt   bool `json:"update_db_encrypt,omitempty"`
	UpdateDbSign      bool `json:"update_db_sign,omitempty"`
	SyncToEs          bool `json:"sync_to_es,omitempty"`
	Upsert            bool `json:"upsert,omitempty"` // db2es/es2db 时, 是否采用 upsert 方式。如果为 true，那么会先根据 id 进行查询，如果存在就更新，不存在就插入。如果为 false，那么会先根据 id 进行查询，如果存在就跳过，不存在就插入。

	//StartId  string `json:"start_id,omitempty"` // 从哪个 id 开始，升序
	StartEid int64 `json:"start_eid,omitempty"`

	// 同步的时候，用于存储一些临时数据
	TempData any `json:"-"`
}

func (o SyncOption) GetUpdateDbBatchSize() int {
	if o.UpdateDbBatchSize > 0 {
		return o.UpdateDbBatchSize
	}
	return 100
}

type SyncArgs struct {
	Srv  iface.Service
	Args *dto.Param
}
