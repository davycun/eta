package cache

/*
keys pattern支持3个通配符*，?，[]：
	*：通配任意多个字符
	?：通配单个字符
	[]：通配括号内的某一个字符
*/

type ScanParam struct {
	Cursor uint64 `json:"cursor,omitempty"`
	Match  string `json:"match,omitempty"`
	Count  int64  `json:"count,omitempty"`
}

type ScanResult struct {
	Keys   []string `json:"keys,omitempty" `
	Cursor uint64   `json:"cursor,omitempty" `
}

type DelParam struct {
	Keys []string `json:"keys,omitempty" binding:"required"` // 删除的key, 支持通配符
}

type SetParam struct {
	Key        string `json:"key" binding:"required"`
	Value      any    `json:"value,omitempty"`
	Expiration int64  `json:"expiration,omitempty"` // 秒
}

type DetailResult struct {
	Ttl   *int64 `json:"ttl,omitempty"`
	Value *any   `json:"value,omitempty"`
}
