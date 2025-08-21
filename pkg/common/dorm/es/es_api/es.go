package es_api

import (
	"github.com/davycun/eta/pkg/common/config"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/typedapi"
	"net/http"
)

type Api struct {
	EsClient      *elasticsearch.Client
	EsTypedClient *elasticsearch.TypedClient
	EsApi         *esapi.API
	EsTypedApi    *typedapi.API
	LogLevel      int
}

func New(cfg config.ElasticConfig) *Api {
	if len(cfg.Addresses) > 0 {
		esCfg := elasticsearch.Config{
			Addresses:              cfg.Addresses,
			Username:               cfg.Username,
			Password:               cfg.Password,
			ServiceToken:           cfg.ServiceToken,
			CertificateFingerprint: cfg.CertificateFingerprint,
		}

		//下面可以取消掉一些自签名证书的校验
		if cfg.InsecureSkipVerify {
			if defaultTransport, ok := http.DefaultTransport.(*http.Transport); ok {
				transport := defaultTransport.Clone()
				transport.TLSClientConfig.InsecureSkipVerify = true
				esCfg.Transport = transport
			} else {
				//ignore
			}
		} else {
			esCfg.CertificateFingerprint = cfg.CertificateFingerprint
			if cfg.CACert != "" {
				esCfg.CACert = []byte(cfg.CACert)
			}
		}

		es, err := elasticsearch.NewClient(esCfg)
		if err != nil {
			logger.Errorf("Elasticsearch NewClient err %s", err)
			return nil
		}
		typedES, err := elasticsearch.NewTypedClient(esCfg)
		if err != nil {
			logger.Errorf("Elasticsearch NewClient err %s", err)
			return nil
		}

		if cfg.LogLevel < 1 {
			cfg.LogLevel = 4
		}

		return &Api{
			EsClient:      es,
			EsTypedClient: typedES,
			EsApi:         esapi.New(es),
			EsTypedApi:    typedapi.New(typedES),
			LogLevel:      cfg.LogLevel,
		}
	}
	return nil
}
