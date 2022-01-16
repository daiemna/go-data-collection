package config

import (
	"fmt"
	"io/ioutil"

	"github.com/daiemna/go-data-collection/internal/bigdb"
	validator "github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

var (
	//ErrConfigUninitalized is thrown whenever a configuration is not initialized.
	ErrConfigUninitalized = fmt.Errorf("Config uninitialized. Please call InitConfig first")
)

type DatabaseInfo struct {
	Hosts                    []string `yaml:"hosts" validate:"required,dive,hostname_port,min=1"`
	Username                 string   `yaml:"username" validate:"required,min=1"`
	Password                 string   `yaml:"password" validate:"required,min=1"`
	KeySpace                 string   `yaml:"key_space" validate:"required,min=1"`
	TableSpace               string   `yaml:"table_space" validate:"required,min=1"`
	ProtoVersion             int      `yaml:"protocol_ver"`
	DisableInitialHostLookup bool     `yaml:"initial_host_look_up"`
}

type ServerInfo struct {
	HostPort string `yaml:"host_port" validate:"hostname_port"`
}

type Config struct {
	Database DatabaseInfo `yaml:"database"`
	Server   ServerInfo   `yaml:"server"`
}

type ClientConfig struct {
	Server ServerInfo `yaml:"server"`
}

//ToCassConfig converts DatabaseInfo to CassandraConfig
func (di *DatabaseInfo) ToCassConfig() *bigdb.CassandraConfig {
	return &bigdb.CassandraConfig{
		IP:                       di.Hosts,
		Username:                 di.Username,
		Password:                 di.Password,
		Keyspace:                 di.KeySpace,
		Table:                    di.TableSpace,
		ProtoVersion:             di.ProtoVersion,
		DisableInitialHostLookup: di.DisableInitialHostLookup,
	}
}

// UnmarshalYAML implements a custom unmarshaller to ensure setting proper default values
func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type rawConfig Config

	// Put your defaults here *******************
	raw := rawConfig{
		Database: DatabaseInfo{
			ProtoVersion:             4,
			DisableInitialHostLookup: false,
		},
		Server: ServerInfo{
			HostPort: "localhost:5005",
		},
	}

	if err := unmarshal(&raw); err != nil {
		return err
	}

	*c = Config(raw)
	return nil
}

// UnmarshalYAML implements a custom unmarshaller to ensure setting proper default values
func (c *ClientConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type rawConfig ClientConfig

	// Put your defaults here *******************
	raw := rawConfig{
		Server: ServerInfo{
			HostPort: "localhost:5005",
		},
	}

	if err := unmarshal(&raw); err != nil {
		return err
	}

	*c = ClientConfig(raw)
	return nil
}

// NewServerFromFile returns the config object for this project
func NewServerFromFile(configPath string) (*Config, error) {

	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to read the configuration file from %s", configPath)
	}
	// if confType == "server" {
	var C Config
	err = yaml.Unmarshal(yamlFile, &C)
	if err != nil {
		return nil, fmt.Errorf("unable to decode configuration into struct, %v", err)
	}
	validate := validator.New()
	err = validate.Struct(C)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			log.Error().Msgf("Failed in config file validation %s, %s", err.Namespace(), err.Kind())
		}
		return nil, err
	}
	return &C, nil
	// } else {

	// }
}

func NewClientFromFile(configPath string) (*ClientConfig, error) {
	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to read the configuration file from %s", configPath)
	}
	var C ClientConfig
	err = yaml.Unmarshal(yamlFile, &C)
	if err != nil {
		return nil, fmt.Errorf("unable to decode configuration into struct, %v", err)
	}
	validate := validator.New()
	err = validate.Struct(C)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			log.Error().Msgf("Failed in config file validation %s, %s", err.Namespace(), err.Kind())
		}
		return nil, err
	}
	return &C, nil
}
