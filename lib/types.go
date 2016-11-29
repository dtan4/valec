package lib

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// Config represents key=value pair
type Config struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

// CompareConfigList compares two config lists and returns the differences between them
func CompareConfigList(src, dst []*Config) ([]*Config, []*Config) {
	added, deleted := []*Config{}, []*Config{}

	for _, c := range src {
		if !configExists(c, dst) {
			added = append(added, c)
		}
	}

	for _, c := range dst {
		if !configExists(c, src) {
			deleted = append(deleted, c)
		}
	}

	return added, deleted
}

// ConfigsToMap converts config list to map
func ConfigsToMap(configs []*Config) map[string]string {
	configMap := map[string]string{}

	for _, config := range configs {
		configMap[config.Key] = config.Value
	}

	return configMap
}

func configExists(config *Config, configs []*Config) bool {
	for _, c := range configs {
		if config.Key == c.Key && config.Value == c.Value {
			return true
		}
	}

	return false
}

// LoadConfigYAML loads configs from the given YAML file
func LoadConfigYAML(filename string) ([]*Config, error) {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return []*Config{}, errors.Wrapf(err, "Failed to read config file. filename=%s", filename)
	}

	var configs []*Config

	if err := yaml.Unmarshal(body, &configs); err != nil {
		return []*Config{}, errors.Wrapf(err, "Failed to parse config file as YAML. filename=%s", filename)
	}

	return configs, nil
}
