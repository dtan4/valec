package lib

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestCompareConfigList(t *testing.T) {
	src := []*Config{
		&Config{
			Key:   "FOO",
			Value: "bar",
		},
		&Config{
			Key:   "BAZ",
			Value: "1",
		},
		&Config{
			Key:   "HOGE",
			Value: "fuga",
		},
	}
	dst := []*Config{
		&Config{
			Key:   "FOO",
			Value: "baz",
		},
		&Config{
			Key:   "BAZ",
			Value: "1",
		},
		&Config{
			Key:   "QUX",
			Value: "true",
		},
		&Config{
			Key:   "PIYO",
			Value: "piyo",
		},
	}

	expectAdded := []*Config{
		&Config{
			Key:   "FOO",
			Value: "bar",
		},
		&Config{
			Key:   "HOGE",
			Value: "fuga",
		},
	}
	expectDeleted := []*Config{
		&Config{
			Key:   "FOO",
			Value: "baz",
		},
		&Config{
			Key:   "QUX",
			Value: "true",
		},
		&Config{
			Key:   "PIYO",
			Value: "piyo",
		},
	}

	added, deleted := CompareConfigList(src, dst)

	if !configListsEqual(added, expectAdded) {
		t.Errorf("Returned added configs are wrong. expected: %s, actual: %s", stringifyConfigList(expectAdded), stringifyConfigList(added))
	}

	if !configListsEqual(deleted, expectDeleted) {
		t.Errorf("Returned deleted configs are wrong. expected: %s, actual: %s", stringifyConfigList(expectDeleted), stringifyConfigList(deleted))
	}
}

func configListsEqual(a, b []*Config) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i].Key != b[i].Key {
			return false
		}

		if a[i].Value != b[i].Value {
			return false
		}
	}

	return true
}

func stringifyConfigList(configs []*Config) string {
	ss := []string{}

	for _, config := range configs {
		ss = append(ss, fmt.Sprintf("%#v", config))
	}

	return fmt.Sprintf("[%s]", strings.Join(ss, ", "))
}

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

func TestConfigsToMap(t *testing.T) {
	configs := []*Config{
		&Config{
			Key:   "FOO",
			Value: "bar",
		},
		&Config{
			Key:   "BAZ",
			Value: "1",
		},
		&Config{
			Key:   "HOGE",
			Value: "fuga",
		},
	}
	expected := map[string]string{
		"FOO":  "bar",
		"BAZ":  "1",
		"HOGE": "fuga",
	}

	configMap := ConfigsToMap(configs)

	if !reflect.DeepEqual(configMap, expected) {
		t.Errorf("Config map does not match. expected: %q, actual:%q", expected, configMap)
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
