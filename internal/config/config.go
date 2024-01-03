package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type (
	DB struct { // type for db connection string
		Host     string `yaml:"host"` // read from config, in yaml according field will have the name specified here
		Port     int    `yaml:"port"`
		Name     string `yaml:"name"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	}
	NATS struct { // linking struct for arguments to connect to NATS streaming
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	}
	STAN struct { // sane fir STAN,
		ClusterID string `yaml:"clusterId"`
		ClientID  string `yaml:"clientId"`
	}
	Service struct { // main service settings
		QueueName   string `yaml:"queueName"`
		QueueGroup  string `yaml:"queueGroup"`
		StartSeq    uint64 `yaml:"startSeq"`
		DeliverLast bool   `yaml:"deliverLast"`
		DeliverAll  bool   `yaml:"deliverAll"`
		NewOnly     bool   `yaml:"newOnly"`
		StartDelta  string `yaml:"startDelta"`
	}
	HTTP struct {
		Port int `yaml:"port"`
	}
	Config struct {
		DB      DB      `yaml:"db"`
		NATS    NATS    `yaml:"nats"`
		STAN    STAN    `yaml:"stan"`
		Service Service `yaml:"service"`
		HTTP    HTTP    `yaml:"http"`
	}
)

// Parse config from yaml file.
// it's fields will be used as parameters
// to connect and configure services, db, streaming etc.
func Parse(fileName string) (*Config, error) { // *Config pointer to be able to return nil if error is present
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to read config from yml file: %w", err)
	}
	cfg := Config{}
	err = yaml.Unmarshal(data, &cfg) // try to convert data from yml format to go-struct
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config data: %w", err)
	}
	return &cfg, nil
}
