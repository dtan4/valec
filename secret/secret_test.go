package secret

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestCompareList(t *testing.T) {
	src := []*Secret{
		&Secret{
			Key:   "FOO",
			Value: "bar",
		},
		&Secret{
			Key:   "BAZ",
			Value: "1",
		},
		&Secret{
			Key:   "HOGE",
			Value: "fuga",
		},
	}
	dst := []*Secret{
		&Secret{
			Key:   "FOO",
			Value: "baz",
		},
		&Secret{
			Key:   "BAZ",
			Value: "1",
		},
		&Secret{
			Key:   "QUX",
			Value: "true",
		},
		&Secret{
			Key:   "PIYO",
			Value: "piyo",
		},
	}

	expectAdded := []*Secret{
		&Secret{
			Key:   "FOO",
			Value: "bar",
		},
		&Secret{
			Key:   "HOGE",
			Value: "fuga",
		},
	}
	expectDeleted := []*Secret{
		&Secret{
			Key:   "FOO",
			Value: "baz",
		},
		&Secret{
			Key:   "QUX",
			Value: "true",
		},
		&Secret{
			Key:   "PIYO",
			Value: "piyo",
		},
	}

	added, deleted := CompareList(src, dst)

	if !secretListsEqual(added, expectAdded) {
		t.Errorf("Returned added secrets are wrong. expected: %s, actual: %s", stringifySecretList(expectAdded), stringifySecretList(added))
	}

	if !secretListsEqual(deleted, expectDeleted) {
		t.Errorf("Returned deleted secrets are wrong. expected: %s, actual: %s", stringifySecretList(expectDeleted), stringifySecretList(deleted))
	}
}

func secretListsEqual(a, b []*Secret) bool {
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

func stringifySecretList(secrets []*Secret) string {
	ss := []string{}

	for _, secret := range secrets {
		ss = append(ss, fmt.Sprintf("%#v", secret))
	}

	return fmt.Sprintf("[%s]", strings.Join(ss, ", "))
}

func TestLoadFromFromYAML_valid(t *testing.T) {
	filepath := testdataPath("test_valid.yaml")
	secrets, err := LoadFromYAML(filepath)
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

	if len(secrets) != len(expects) {
		t.Fatalf("Secrets does not loaded correctly. expected length: %d, actual length: %d", len(expects), len(secrets))
	}

	for i, secret := range secrets {
		if secret.Key != expects[i].key {
			t.Errorf("Secret key does not match. expected: %s, actual: %s", expects[i].key, secret.Key)
		}

		if secret.Value != expects[i].value {
			t.Errorf("Secret value does not match. expected: %s, actual: %s", expects[i].value, secret.Value)
		}
	}
}

func TestListToMap(t *testing.T) {
	secrets := []*Secret{
		&Secret{
			Key:   "FOO",
			Value: "bar",
		},
		&Secret{
			Key:   "BAZ",
			Value: "1",
		},
		&Secret{
			Key:   "HOGE",
			Value: "fuga",
		},
	}
	expected := map[string]string{
		"FOO":  "bar",
		"BAZ":  "1",
		"HOGE": "fuga",
	}

	secretMap := ListToMap(secrets)

	if !reflect.DeepEqual(secretMap, expected) {
		t.Errorf("Secret map does not match. expected: %q, actual:%q", expected, secretMap)
	}
}

func TestLoadFromFromYAML_invalid(t *testing.T) {
	filepath := testdataPath("test_invalid.yaml")
	_, err := LoadFromYAML(filepath)
	if err == nil {
		t.Fatalf("Error should be raised. error: %s", err)
	}

	expected := fmt.Sprintf("Failed to parse secret file as YAML. filename=%s", filepath)

	if !strings.HasPrefix(err.Error(), expected) {
		t.Fatalf("Error message prefix does not match. expected prefix: %q, actual message: %q", expected, err.Error())
	}
}

func TestLoadFromFromYAML_notexist(t *testing.T) {
	filepath := testdataPath("test_notexist.yaml")
	_, err := LoadFromYAML(filepath)
	if err == nil {
		t.Fatalf("Error should be raised. error: %s", err)
	}

	expected := fmt.Sprintf("Failed to read secret file. filename=%s", filepath)

	if !strings.HasPrefix(err.Error(), expected) {
		t.Fatalf("Error message prefix does not match. expected prefix: %q, actual message: %q", expected, err.Error())
	}
}

func testdataPath(name string) string {
	return filepath.Join("..", "testdata", name)
}

func TestMapToSecrets(t *testing.T) {
	secretMap := map[string]string{
		"FOO":  "bar",
		"BAZ":  "1",
		"HOGE": "fuga",
	}
	expected := []*Secret{
		&Secret{
			Key:   "FOO",
			Value: "bar",
		},
		&Secret{
			Key:   "BAZ",
			Value: "1",
		},
		&Secret{
			Key:   "HOGE",
			Value: "fuga",
		},
	}

	secrets := MapToSecrets(secretMap)

	if !reflect.DeepEqual(secrets, expected) {
		t.Errorf("Secret list does not match. expected: %q, actual:%q", expected, secrets)
	}
}

func TestSaveAsYAML(t *testing.T) {
	secrets := []*Secret{
		&Secret{
			Key:   "FOO",
			Value: "bar",
		},
		&Secret{
			Key:   "BAZ",
			Value: "1",
		},
		&Secret{
			Key:   "HOGE",
			Value: "fuga",
		},
	}

	dir, err := ioutil.TempDir("", "test-save-as-yaml")
	if err != nil {
		t.Fatalf("Failed to create tempdir. dir: %s", dir)
	}
	defer os.RemoveAll(dir)

	filename := filepath.Join(dir, "secret.yaml")

	if err := SaveAsYAML(secrets, filename); err != nil {
		t.Fatalf("Error should not be raised. err: %s", err)
	}

	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			t.Fatalf("File is not created. err: %s", err)
		} else {
			t.Fatalf("Saved file has something wrong. err: %s", err)
		}
	}
}
