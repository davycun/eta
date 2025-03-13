package config

import (
	"os"
	"strconv"
)

const (
	EtaDbHost        = "ETA_DB_HOST"
	EtaDbPort        = "ETA_DB_PORT"
	EtaDbUser        = "ETA_DB_USER"
	EtaDbPassword    = "ETA_DB_PASSWORD"
	EtaDbName        = "ETA_DB_NAME"
	EtaDbSchema      = "ETA_DB_SCHEMA"
	EtaSrvPort       = "ETA_SRV_PORT"
	EtaMonitorPort   = "ETA_MONITOR_PORT"
	EtaRedisHost     = "ETA_REDIS_HOST"
	EtaRedisPort     = "ETA_REDIS_PORT"
	EtaRedisUsername = "ETA_REDIS_USERNAME"
	EtaRedisPassword = "ETA_REDIS_PASSWORD"
	EtaRedisDB       = "ETA_REDIS_DB"
)

func LoadFromEnv() *Configuration {
	cf := &Configuration{}
	cf.Database.Host = os.Getenv(EtaDbHost)
	cf.Database.User = os.Getenv(EtaDbUser)
	cf.Database.Password = os.Getenv(EtaDbPassword)
	cf.Database.DBName = os.Getenv(EtaDbName)
	cf.Database.Schema = os.Getenv(EtaDbSchema)
	if os.Getenv(EtaDbPort) != "" {
		port, err := strconv.Atoi(os.Getenv(EtaDbPort))
		if err == nil {
			cf.Database.Port = port
		}
	}

	cf.Redis.Host = os.Getenv(EtaRedisHost)
	cf.Redis.Username = os.Getenv(EtaRedisUsername)
	cf.Redis.Password = os.Getenv(EtaRedisPassword)
	if os.Getenv(EtaRedisPort) != "" {
		port, err := strconv.Atoi(os.Getenv(EtaRedisPort))
		if err == nil {
			cf.Redis.Port = port
		}
	}
	if os.Getenv(EtaRedisDB) != "" {
		db, err := strconv.Atoi(os.Getenv(EtaRedisDB))
		if err == nil {
			cf.Redis.DB = db
		}
	}

	if os.Getenv(EtaSrvPort) != "" {
		port, err := strconv.Atoi(os.Getenv(EtaSrvPort))
		if err == nil {
			cf.Server.Port = port
		}
	}
	if os.Getenv(EtaMonitorPort) != "" {
		port, err := strconv.Atoi(os.Getenv(EtaMonitorPort))
		if err == nil {
			cf.Monitor.Port = port
		}
	}
	return cf
}
