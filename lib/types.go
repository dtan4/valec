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
