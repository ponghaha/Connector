package config

import (
	"encoding/json"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server       ServerConfig           `yaml:"server" json:"server"`
	Logger       LoggerConfig           `yaml:"logger" json:"logger"`
	APIKeys      []APIKey               `yaml:"apiKeys" json:"apiKeys"`
	Destinations map[string]Destination `yaml:"destinations" json:"destinations"`
	Routes       map[string]Route       `yaml:"routes" json:"routes"`
	ELKPath      string                 `yaml:"elkPath"`
}
type ServerConfig struct {
	Port string `yaml:"port"`
	Mode string `yaml:"mode"`
}
type LoggerConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}
type APIKey struct {
	Key         []string   `yaml:"key"`
	ClientName  string   `yaml:"clientName"`
	Status      string   `yaml:"status"`
	Permissions []string `yaml:"permissions"`
}
type Destination struct {
	Type   string              `json:"type"`
	IP     string              `json:"ip"`
	Ports  map[string][]string `json:"ports"`
	APIKey string              `json:"apiKey"`
}
type Route struct {
	System  		string `json:"System"`
	SystemV1  		string `json:"SystemV1"`
	SystemV2  		string `json:"SystemV2"`
	Service 		string `json:"Service"`
	Format  		string `json:"Format"`
	FormatV1  		string `json:"FormatV1"`
	FormatV2  		string `json:"FormatV2"`
	RequestLength   string `json:"RequestLength"`
}
type DestinationsAndRoutes struct {
	Destinations map[string]Destination `json:"destinations"`
	Routes       map[string]Route       `json:"routes"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config Config
	if err = yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func LoadAPIKeys(path string) ([]APIKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var apiKeys []APIKey
	if err := json.Unmarshal(data, &apiKeys); err != nil {
		return nil, err
	}
	return apiKeys, nil
}

func LoadDestinationsAndRoutes(path string) (*DestinationsAndRoutes, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var dr DestinationsAndRoutes
	if err := json.Unmarshal(data, &dr); err != nil {
		return nil, err
	}
	return &dr, nil
}
