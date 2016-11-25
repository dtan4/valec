package lib

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadConfigFromYAML_valid(t *testing.T) {
	filepath := testdataPath("test_valid.yaml")
	configs, err := LoadConfigYAML(filepath)
	if err != nil {
		t.Fatalf("Error should not be raised. error: %s", err)
	}

	expects := []struct {
		key   string
		value string
	}{
		{"FOO", "bar"},
		{"BAZ", "1"},
		{"QUX", "true"},
	}

	if len(configs) != len(expects) {
		t.Fatalf("Configs does not loaded correctly. expected length: %d, actual length: %d", len(expects), len(configs))
	}

	for i, config := range configs {
		if config.Key != expects[i].key {
			t.Errorf("Config key does not match. expected: %s, actual: %s", expects[i].key, config.Key)
		}

		if config.Value != expects[i].value {
			t.Errorf("Config value does not match. expected: %s, actual: %s", expects[i].value, config.Value)
		}
	}
}

func TestLoadConfigFromYAML_invalid(t *testing.T) {
	filepath := testdataPath("test_invalid.yaml")
	_, err := LoadConfigYAML(filepath)
	if err == nil {
		t.Fatalf("Error should be raised. error: %s", err)
	}

	expected := fmt.Sprintf("Failed to parse config file as YAML. filename=%s", filepath)

	if !strings.HasPrefix(err.Error(), expected) {
		t.Fatalf("Error message prefix does not match. expected prefix: %q, actual message: %q", expected, err.Error())
	}
}

func TestLoadConfigFromYAML_notexist(t *testing.T) {
	filepath := testdataPath("test_notexist.yaml")
	_, err := LoadConfigYAML(filepath)
	if err == nil {
		t.Fatalf("Error should be raised. error: %s", err)
	}

	expected := fmt.Sprintf("Failed to read config file. filename=%s", filepath)

	if !strings.HasPrefix(err.Error(), expected) {
		t.Fatalf("Error message prefix does not match. expected prefix: %q, actual message: %q", expected, err.Error())
	}
}

func testdataPath(name string) string {
	return filepath.Join("..", "testdata", name)
}
