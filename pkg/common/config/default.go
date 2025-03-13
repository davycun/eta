package config

import "github.com/davycun/eta/pkg/common/dorm"

var (
	defaultConf = Configuration{ // 配置，并设置默认值
		Database: dorm.Database{Host: "127.0.0.1", Port: 5432, Schema: "eta"},
		Server: Server{
			Port: 8080,
			Env:  "unknown",
			IgnoreUri: []string{
				"/oauth2/*",
				"/storage/public/*",
				"/storage/download/*",
				"/storage/upload/*",
				"/public_key",
				"/api_doc/*",
			},
		},
		Monitor: Monitor{
			Port:         6060,
			SentryEnable: false,
			//SentryDsn:              "",
			SentrySendDefaultPii:   true,
			SentryTracesSampleRate: 0.1,
			SentryEnableTracing:    true,
		},
		ID: IdGenerator{Type: IdGeneratorTypeConfig, NodeId: 1},
		EsConfig: ElasticConfig{
			LogLevel:           4,
			InsecureSkipVerify: true,
			NumberOfShards:     1,
			NumberOfReplicas:   0,
		},
	}
)

func UpdateDefaultConfig(fc func(cfg *Configuration)) {
	fc(&defaultConf)
}
