package config

import (
	"errors"
	"github.com/davycun/eta/pkg/common/logger"
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

func LoadFromFile(file string) (*Configuration, error) {
	c := &Configuration{}
	if file == "" {
		return c, errors.New("the file path is empty when LoadFromFile")
	}
	logger.Infof("config file: %s", file)
	open, err := os.Open(file)
	defer func() {
		_ = open.Close()
	}()
	if err != nil {
		return c, err
	}
	all, err := io.ReadAll(open)
	if err != nil {
		return c, err
	}
	err = yaml.Unmarshal(all, c)
	return c, err
}
