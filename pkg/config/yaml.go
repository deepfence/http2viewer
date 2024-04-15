package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type PCAPConfig struct {
}

type JSONConfig struct {
	StreamIDKey string `yaml:"stream_id_key"`
	DataKey     string `yaml:"data_key"`
}

type ConfigSource struct {
	FilePath   string      `yaml:"file_path"`
	PCAPConfig *PCAPConfig `yaml:"pcap_config"`
	JSONCondig *JSONConfig `yaml:"json_config"`
}

type Config struct {
	InputSources []ConfigSource `yaml:"input_sources"`
}

func LoadConfig(config_path string) (Config, error) {
	b, err := os.ReadFile(config_path)
	var res Config
	if err != nil {
		return res, err
	}
	err = yaml.Unmarshal(b, &res)
	if err != nil {
		return res, err
	}
	for i, s := range res.InputSources {
		if s.PCAPConfig == nil && s.JSONCondig == nil {
			return res, fmt.Errorf("entry %d invalid, no types configured", i)
		} else if s.PCAPConfig != nil && s.JSONCondig != nil {
			return res, fmt.Errorf("entry %d invalid, too many configured types, only one is allowed", i)
		}
	}
	return res, err
}
