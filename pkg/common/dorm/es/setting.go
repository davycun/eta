package es

// DefaultAnalysis
// 不用变量的原因是，确保默认的返回都一直，不会存在变量被修改的情况
func DefaultAnalysis() map[string]interface{} {
	return map[string]interface{}{
		"tokenizer": map[string]interface{}{
			"digit_tokenizer": map[string]interface{}{
				"type":    "pattern",
				"pattern": "",
			},
		},
		"analyzer": map[string]interface{}{
			"digit_analyzer": map[string]interface{}{
				"type":      "custom",
				"tokenizer": "digit_tokenizer",
			},
		},
	}
}

// DefaultSetting
// 数据量不大，ES基本都是单节点，所以replicas未0
func DefaultSetting() map[string]interface{} {
	return map[string]interface{}{
		"number_of_shards":   1,
		"number_of_replicas": 0,
		"analysis":           DefaultAnalysis(),
	}
}

// DefaultKeyword
// 字段默认的 keyword field
func DefaultKeyword() map[string]interface{} {
	return map[string]interface{}{
		"type": "keyword",
		//"ignore_above": 256,
	}
}
