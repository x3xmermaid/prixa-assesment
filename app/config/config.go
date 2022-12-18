package config

import (
	"fmt"
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// YAMLConfigLoader is the loader of YAML file configuration.
type YAMLConfigLoader struct {
	fileLocation string
}

// NewYamlConfigLoader return the YAML Configuration loader.
func NewYamlConfigLoader(fileLocation string) *YAMLConfigLoader {
	return &YAMLConfigLoader{
		fileLocation: fileLocation,
	}
}

// ServiceConfig stores the whole configuration for service.
type ServiceConfig struct {
	ServiceData ServiceDataConfig `yaml:"service_data"`
	SourceData  SourceDataConfig  `yaml:"source_data"`
}

// ServiceDataConfig contains the service data configuration.
type ServiceDataConfig struct {
	Address string `yaml:"address"`
}

// SourceDataConfig contains the source data configuration.
type SourceDataConfig struct {
	RedisNetwork           string `yaml:"redis_network"`
	RedisAddress           string `yaml:"redis_address"`
	RedisPassword          string `yaml:"redis_password"`
	RedisTimeout           int64  `yaml:"redis_timeout"`
	RedisKeyExpireDuration int64  `yaml:"redis_key_expire_duration"`
	CacheDuration          int    `yaml:"cache_duration"`
}

func getRawConfig(fileLocation string) (*ServiceConfig, error) {
	configByte, err := ioutil.ReadFile(fileLocation)
	if err != nil {
		logrus.Errorf("Error Read File Raw Config: %v", err)
		return nil, err
	}
	config := &ServiceConfig{}
	err = yaml.Unmarshal(configByte, config)
	if err != nil {
		logrus.Errorf("Error Unmarshal Raw Config: %v", err)
		return nil, err
	}
	return config, nil
}

// GetServiceConfig parse the configuration from YAML file.
func (c *YAMLConfigLoader) GetServiceConfig() (*ServiceConfig, error) {
	config, err := getRawConfig(c.fileLocation)
	if err != nil {
		return nil, fmt.Errorf("unable to get raw config content: %v", err)
	}

	return config, nil
}
