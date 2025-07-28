package config

import (
	"errors"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/hashicorp/consul/api"
	jsoniter "github.com/json-iterator/go"
	"os"
	"strings"
)

const (
	envConfigServerForce = "MDT_CONFIG_SERVER_FORCE"
	envConfigServer      = "MDT_CONFIG_SERVER"
	split                = "/"
)

var (
	csNs = "eta" // config server namespace
)

func SetConsulNamespace(namespace string) {
	if namespace != "" {
		csNs = namespace
	} else {
		logger.Debug("consul namespace is empty, use default namespace: %s", csNs)
	}
}

func LoadFromConsul() (*Configuration, error) {
	var (
		err error
		cf  = &Configuration{}
	)
	forceConsulUrl := os.Getenv(envConfigServerForce)
	consulUrl := os.Getenv(envConfigServer)
	if forceConsulUrl != "" {
		consulUrl = forceConsulUrl
	}
	if consulUrl == "" {
		logger.Warnf("没有配置consul")
		return cf, errors.New("没有配置consul")
	}
	_, after, found := strings.Cut(consulUrl, "://")
	if found {
		consulUrl = strings.TrimSpace(after)
	}
	logger.Infof("consul url: %s", consulUrl)
	cm := NewConsulManager(strings.Split(consulUrl, ";"))
	GetFromConsul(cm, cf)

	return cf, err
}

func getDatabase(cm *ConsulManager, cf *Configuration, cfgName string, db *dorm.Database) {
	cm.Get(cm.Path(split, csNs, cfgName, "host"), &db.Host, "")
	cm.Get(cm.Path(split, csNs, cfgName, "port"), &db.Port, 0)
	cm.Get(cm.Path(split, csNs, cfgName, "user"), &db.User, "")
	cm.Get(cm.Path(split, csNs, cfgName, "password"), &db.Password, "")
	cm.Get(cm.Path(split, csNs, cfgName, "dbname"), &db.DBName, "")
	cm.Get(cm.Path(split, csNs, cfgName, "schema"), &db.Schema, "")
	cm.Get(cm.Path(split, csNs, cfgName, "type"), &db.Type, "")
	cm.Get(cm.Path(split, csNs, cfgName, "key"), &db.Key, "")
	cm.Get(cm.Path(split, csNs, cfgName, "log_level"), &db.LogLevel, 0)
	cm.Get(cm.Path(split, csNs, cfgName, "slow_threshold"), &db.SlowThreshold, 0)
	cm.Get(cm.Path(split, csNs, cfgName, "max_open_cons"), &db.MaxOpenCons, 0)
	cm.Get(cm.Path(split, csNs, cfgName, "max_idle_cons"), &db.MaxIdleCons, 0)
	cm.Get(cm.Path(split, csNs, cfgName, "conn_max_lifetime"), &db.ConMaxLifetime, 0)
	cm.Get(cm.Path(split, csNs, cfgName, "conn_max_idle_time"), &db.ConMaxIdleTime, 0)
}

func GetFromConsul(cm *ConsulManager, cf *Configuration) {

	getDatabase(cm, cf, "database", &cf.Database)
	getDatabase(cm, cf, "doris", &cf.Doris)

	svr := "server"
	cm.Get(cm.Path(split, csNs, svr, "host"), &cf.Server.Host, "")
	cm.Get(cm.Path(split, csNs, svr, "port"), &cf.Server.Port, 0)
	cm.Get(cm.Path(split, csNs, svr, "gin_mode"), &cf.Server.GinMode, "")
	cm.Get(cm.Path(split, csNs, svr, "middleware"), &cf.Server.Middleware, nil)
	cm.Get(cm.Path(split, csNs, svr, "router_pkg"), &cf.Server.RouterPkg, nil)
	cm.Get(cm.Path(split, csNs, svr, "migrate_pkg"), &cf.Server.MigratePkg, nil)
	cm.Get(cm.Path(split, csNs, svr, "raw_fetch_cache"), &cf.Server.RawFetchCache, false)
	cm.Get(cm.Path(split, csNs, svr, "env"), &cf.Server.Env, "")
	cm.Get(cm.Path(split, csNs, svr, "api_doc_enable"), &cf.Server.ApiDocEnable, false)

	scConfig := "es_config"
	cm.Get(cm.Path(split, csNs, scConfig, "addresses"), &cf.EsConfig.Addresses, nil)
	cm.Get(cm.Path(split, csNs, scConfig, "username"), &cf.EsConfig.Username, "")
	cm.Get(cm.Path(split, csNs, scConfig, "password"), &cf.EsConfig.Password, "")
	cm.Get(cm.Path(split, csNs, scConfig, "service_token"), &cf.EsConfig.ServiceToken, "")
	cm.Get(cm.Path(split, csNs, scConfig, "certificate_fingerprint"), &cf.EsConfig.CertificateFingerprint, "")
	cm.Get(cm.Path(split, csNs, scConfig, "ca_cert"), &cf.EsConfig.CACert, "")
	cm.Get(cm.Path(split, csNs, scConfig, "insecure_skip_verify"), &cf.EsConfig.InsecureSkipVerify, true)
	cm.Get(cm.Path(split, csNs, scConfig, "log_level"), &cf.EsConfig.LogLevel, 4)
	cm.Get(cm.Path(split, csNs, scConfig, "number_of_shards"), &cf.EsConfig.NumberOfShards, 1)
	cm.Get(cm.Path(split, csNs, scConfig, "number_of_replicas"), &cf.EsConfig.NumberOfReplicas, 0)

	m := "monitor"
	cm.Get(cm.Path(split, csNs, m, "port"), &cf.Monitor.Port, 0)
	cm.Get(cm.Path(split, csNs, m, "sentry_enable"), &cf.Monitor.SentryEnable, false)
	cm.Get(cm.Path(split, csNs, m, "sentry_dsn"), &cf.Monitor.SentryDsn, "")
	cm.Get(cm.Path(split, csNs, m, "sentry_send_default_pii"), &cf.Monitor.SentrySendDefaultPii, true)
	cm.Get(cm.Path(split, csNs, m, "sentry_traces_sample_rate"), &cf.Monitor.SentryTracesSampleRate, 0.1)
	cm.Get(cm.Path(split, csNs, m, "sentry_enable_tracing"), &cf.Monitor.SentryEnableTracing, true)

	rds := "redis"
	cm.Get(cm.Path(split, csNs, rds, "host"), &cf.Redis.Host, "")
	cm.Get(cm.Path(split, csNs, rds, "port"), &cf.Redis.Port, 0)
	cm.Get(cm.Path(split, csNs, rds, "user_name"), &cf.Redis.Username, "")
	cm.Get(cm.Path(split, csNs, rds, "password"), &cf.Redis.Password, "")
	cm.Get(cm.Path(split, csNs, rds, "db"), &cf.Redis.DB, 0)

	mig := "migrate"
	cm.Get(cm.Path(split, csNs, mig), &cf.Migrate, false)

	v := "variables"
	cm.Get(cm.Path(split, csNs, v), &cf.Variables, make(map[string]string))

	i := "id"
	cm.Get(cm.Path(split, csNs, i, "node_id"), &cf.ID.NodeId, 0)
	cm.Get(cm.Path(split, csNs, i, "epoch"), &cf.ID.Epoch, "")

	t := "third"
	cm.Get(cm.Path(split, csNs, t, "datlas_base_url"), &cf.Third.DatlasBaseUrl, "")
	cm.Get(cm.Path(split, csNs, t, "datlas_user_name"), &cf.Third.DatlasUserName, "")
	cm.Get(cm.Path(split, csNs, t, "datlas_password"), &cf.Third.DatlasPassword, "")
}

type ConsulManager struct {
	kv *api.KV
}

func NewConsulManager(machines []string) *ConsulManager {
	conf := api.DefaultConfig()
	if len(machines) > 0 {
		conf.Address = machines[0]
	}
	client, err := api.NewClient(conf)
	if err != nil {
		logger.Errorf("初始化consul客户端失败。%v", err)
	}
	return &ConsulManager{
		kv: client.KV(),
	}
}

func (c *ConsulManager) Path(split string, names ...string) string {
	return strings.Join(names, split)
}

func (c *ConsulManager) Get(key string, cfg any, defaultV any) {
	v, _, err := c.kv.Get(key, nil)
	if err == nil && v != nil {
		cv := v.Value
		cnLen := len(cv)
		if cv[0] == '"' && cv[1] == '{' && cv[cnLen-1] == '"' && cv[cnLen-2] == '}' {
			logger.Infof("consul 特殊字符处理, %s", key)
			cv = cv[1 : cnLen-1]
		}
		err = jsoniter.Unmarshal(cv, cfg)
	}
	if err != nil {
		cfg = &defaultV
	}
}
