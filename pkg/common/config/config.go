package config

import (
	"errors"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/common/utils/bean"
	link "github.com/duke-git/lancet/v2/datastructure/link"
	"github.com/duke-git/lancet/v2/slice"
	jsoniter "github.com/json-iterator/go"
	"gopkg.in/yaml.v3"
)

const (
	IdGeneratorTypeConfig = "config"
	IdGeneratorTypeRedis  = "redis"
	IdGeneratorRedisKey   = "eta:idgenerator"
)

type Configuration struct {
	Database  dorm.Database     `json:"database" yaml:"database"`
	Doris     dorm.Database     `json:"doris" yaml:"doris"`
	Server    Server            `json:"server" yaml:"server"`
	Monitor   Monitor           `json:"monitor" yaml:"monitor"`
	Redis     Redis             `json:"redis" yaml:"redis"`
	Nebula    Nebula            `json:"nebula" yaml:"nebula"`
	Migrate   bool              `json:"migrate" yaml:"migrate"`
	Variables map[string]string `json:"variables" yaml:"variables"`
	ID        IdGenerator       `json:"id" yaml:"id"`
	Third     Third             `json:"third" yaml:"third"`
	EsConfig  ElasticConfig     `json:"es_config" yaml:"es_config"`
}

func LoadConfig(cfgFile string, argConfig *Configuration) *Configuration {

	confLink := link.NewSinglyLink[*Configuration]()
	//从环境变量加载配置，优先级4
	confLink.InsertAtTail(LoadFromEnv())
	logger.Debugf("env config added")

	//从配置文件加载配置，优先级3
	ymlConfig, err := LoadFromFile(cfgFile)
	if err != nil {
		logger.Warnf("从配置文件[%s]读取配置异常,%v", cfgFile, err)
	} else {
		confLink.InsertAtTail(ymlConfig)
		logger.Debugf("file config added")
	}

	//从consul加载配置，优先级2
	consulConfig, err := LoadFromConsul()
	if err != nil {
		logger.Warnf("从配置中心读取配置失败,%s", err)
	} else {
		confLink.InsertAtTail(consulConfig)
		logger.Debugf("consul config added")
	}

	//从命令行加载配置，优先级1
	if argConfig != nil {
		confLink.InsertAtTail(argConfig)
		logger.Debugf("arg config added")
	}

	destConfig := mergeConfiguration(confLink)
	checkRequiredCfg(destConfig)
	marshal, err := jsoniter.Marshal(destConfig)
	if err == nil {
		logger.Debugf("config: %s", marshal)
	}
	dorm.RawFetchCache = destConfig.Server.RawFetchCache

	return destConfig
}

func LoadFromContent(content string, c *Configuration) error {
	bs := utils.StringToBytes(content)
	return yaml.Unmarshal(bs, c)
}

func mergeConfiguration(lk *link.SinglyLink[*Configuration]) *Configuration {

	//先保留一份默认配置
	dc := &Configuration{}
	Copy(dc, &defaultConf)
	for _, c := range lk.Values() {
		Copy(dc, c)
	}
	dc.Server.Middleware = utils.Merge(slice.Unique(dc.Server.Middleware), defaultConf.Server.Middleware...)
	dc.Server.RouterPkg = utils.Merge(slice.Unique(dc.Server.RouterPkg), defaultConf.Server.RouterPkg...)
	dc.Server.MigratePkg = utils.Merge(slice.Unique(dc.Server.MigratePkg), defaultConf.Server.MigratePkg...)
	return dc
}

func checkRequiredCfg(cfg *Configuration) {
	if cfg.Database.Host == "" || cfg.Redis.Host == "" {
		panic(errors.New("配置错误，请检查配置"))
	}
}

func Copy(dst, src *Configuration) {
	err := bean.Copy(dst, src)
	if err != nil {
		logger.Errorf("copy config error, %v", err)
	}
}
