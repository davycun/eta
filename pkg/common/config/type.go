package config

import (
	"fmt"
	"strconv"
	"strings"
)

type IdGenerator struct {
	Type   string `json:"type" yaml:"type"` // config/redis
	NodeId int64  `json:"node_id" yaml:"node_id"`
	Epoch  string `json:"epoch" yaml:"epoch"`
}

type Login struct {
	PwdValidateReg          string `json:"pwd_validate_reg" yaml:"pwd_validate_reg"`                     // 密码校验正则表达式
	FailLockDurationMinutes int64  `json:"fail_lock_duration_minutes" yaml:"fail_lock_duration_minutes"` // 登录失败连续时间
	FailLockMaxTimes        int64  `json:"fail_lock_max_times" yaml:"fail_lock_max_times"`               // 登录失败次数，达到后会锁定
	FailLockLockMinutes     int64  `json:"fail_lock_lock_minutes" yaml:"fail_lock_lock_minutes"`         // 登录锁定时间
}
type Third struct {
	DatlasBaseUrl       string `json:"datlas_base_url" yaml:"datlas_base_url"`             // datlas base url
	DatlasUserName      string `json:"datlas_user_name" yaml:"datlas_user_name"`           // datlas user name
	DatlasPassword      string `json:"datlas_password" yaml:"datlas_password"`             // datlas password
	CollectorBaseUrl    string `json:"collector_base_url" yaml:"collector_base_url"`       // collector base url
	CollectorFixedToken string `json:"collector_fixed_token" yaml:"collector_fixed_token"` // collector fixed token
}

type Monitor struct {
	Port                   int     `yaml:"port"`
	SentryEnable           bool    `yaml:"sentry_enable"`
	SentryDsn              string  `yaml:"sentry_dsn"`
	SentrySendDefaultPii   bool    `yaml:"sentry_send_default_pii"`
	SentryTracesSampleRate float64 `yaml:"sentry_traces_sample_rate"`
	SentryEnableTracing    bool    `yaml:"sentry_enable_tracing"`
}

func (m Monitor) GetPort() int {
	if m.Port < 1 {
		return 6060
	}
	return m.Port
}

type Redis struct {
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	Username string `json:"user_name" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	DB       int    `json:"db" yaml:"db"`
}

func (r Redis) Addr() string {
	//return r.Host + ":" + strconv.Itoa(r.Port)
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

type Server struct {
	Host          string   `yaml:"host"`
	Port          int      `yaml:"port"`
	GinMode       string   `yaml:"gin_mode"` //release、debug
	Middleware    []string `json:"middleware" yaml:"middleware"`
	RouterPkg     []string `json:"router_pkg" yaml:"router_pkg"`
	MigratePkg    []string `json:"migrate_pkg" yaml:"migrate_pkg"`
	RawFetchCache bool     `json:"raw_fetch_cache" yaml:"raw_fetch_cache"`
	Env           string   `json:"env" yaml:"env"`
	ApiDocEnable  bool     `json:"api_doc_enable" yaml:"api_doc_enable"`
}

type ElasticConfig struct {
	Addresses              []string `json:"addresses,omitempty" yaml:"addresses"`                   // A list of ElasticConfig nodes to use.
	Username               string   `json:"username" yaml:"username"`                               // Username for HTTP Basic Authentication.
	Password               string   `json:"password" yaml:"password"`                               // Password for HTTP Basic Authentication.
	ServiceToken           string   `json:"service_token" yaml:"service_token"`                     // Service token for authorization; if set, overrides username/password.
	CertificateFingerprint string   `json:"certificate_fingerprint" yaml:"certificate_fingerprint"` // SHA256 hex fingerprint given by ElasticConfig on first launch.
	CACert                 string   `json:"ca_cert" yaml:"ca_cert"`
	InsecureSkipVerify     bool     `json:"insecure_skip_verify" yaml:"insecure_skip_verify"`
	LogLevel               int      `json:"log_level" yaml:"log_level"`
	NumberOfShards         int      `json:"number_of_shards" yaml:"number_of_shards"`
	NumberOfReplicas       int      `json:"number_of_replicas" yaml:"number_of_replicas"`
}

type Nebula struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Space    string `yaml:"space"`
	//unit second
	IdleTime int    `yaml:"idle_time"`
	Timeout  int    `yaml:"timeout"`
	MaxSize  int    `yaml:"max_size"`
	MinSize  int    `yaml:"min_size"`
	Key      string `json:"key" yaml:"key"`
}

func (n *Nebula) GetKey() string {
	if n.Key != "" {
		return n.Key
	}
	n.Key = strings.Join([]string{n.Host, strconv.Itoa(n.Port), n.Space}, ":")
	return n.Key
}
