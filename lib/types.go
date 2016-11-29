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
	srcMap, dstMap := ConfigsToMap(src), ConfigsToMap(dst)

	for _, c := range src {
		v, ok := dstMap[c.Key]
		if !ok || v != c.Value {
			added = append(added, c)
		}
	}

	for _, c := range dst {
		v, ok := srcMap[c.Key]
		if !ok || v != c.Value {
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
