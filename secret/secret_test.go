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

func TestLen(t *testing.T) {
	secrets := Secrets{
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

	expected := 3
	actual := secrets.Len()

	if actual != expected {
		t.Errorf("Length of secrets is wrong. expected: %d, actual: %d", expected, actual)
	}
}

func TestLess(t *testing.T) {
	secrets := Secrets{
		&Secret{
			Key:   "BAZ",
			Value: "1",
		},
		&Secret{
			Key:   "FOO",
			Value: "bar",
		},
		&Secret{
			Key:   "FOO",
			Value: "fuga",
		},
	}
	testcases := []struct {
		i        int
		j        int
		expected bool
	}{
		{
			i:        0,
			j:        1,
			expected: true,
		},
		{
			i:        1,
			j:        0,
			expected: false,
		},
		{
			i:        0,
			j:        0,
			expected: false,
		},
		{
			i:        1,
			j:        2,
			expected: true,
		},
	}

	for _, tc := range testcases {
		actual := secrets.Less(tc.i, tc.j)
		if actual != tc.expected {
			t.Errorf("Comparison result is wrong. src: %q, dst: %q, expected: %t, actual: %t", secrets[tc.i], secrets[tc.j], tc.expected, actual)
		}
	}
}

func TestSwap(t *testing.T) {
	secrets := Secrets{
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

	i, j := 0, 1
	expected := Secrets{
		&Secret{
			Key:   "BAZ",
			Value: "1",
		},
		&Secret{
			Key:   "FOO",
			Value: "bar",
		},
		&Secret{
			Key:   "HOGE",
			Value: "fuga",
		},
	}

	secrets.Swap(i, j)

	if !reflect.DeepEqual(secrets, expected) {
		t.Errorf("Swap result is wrong. expected: %q, actual: %q", expected, secrets)
	}
}

func TestCompareList(t *testing.T) {
	newSecrets := Secrets{
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
	oldSecrets := Secrets{
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

	expectAdded := Secrets{
		&Secret{
			Key:   "HOGE",
			Value: "fuga",
		},
	}
	expectUpdated := Secrets{
		&Secret{
			Key:   "FOO",
			Value: "bar",
		},
	}
	expectDeleted := Secrets{
		&Secret{
			Key:   "QUX",
			Value: "true",
		},
		&Secret{
			Key:   "PIYO",
			Value: "piyo",
		},
	}

	added, updated, deleted := newSecrets.CompareList(oldSecrets)

	if !secretListsEqual(added, expectAdded) {
		t.Errorf("Returned added secrets are wrong. expected: %s, actual: %s", stringifySecretList(expectAdded), stringifySecretList(added))
	}

	if !secretListsEqual(updated, expectUpdated) {
		t.Errorf("Returned updated secrets are wrong. expected: %s, actual: %s", stringifySecretList(expectUpdated), stringifySecretList(updated))
	}

	if !secretListsEqual(deleted, expectDeleted) {
		t.Errorf("Returned deleted secrets are wrong. expected: %s, actual: %s", stringifySecretList(expectDeleted), stringifySecretList(deleted))
	}
}

func secretListsEqual(a, b Secrets) bool {
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

func stringifySecretList(secrets Secrets) string {
	ss := []string{}

	for _, secret := range secrets {
		ss = append(ss, fmt.Sprintf("%#v", secret))
	}

	return fmt.Sprintf("[%s]", strings.Join(ss, ", "))
}

func TestListToMap(t *testing.T) {
	secrets := Secrets{
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

	secretMap := secrets.ListToMap()

	if !reflect.DeepEqual(secretMap, expected) {
		t.Errorf("Secret map does not match. expected: %q, actual:%q", expected, secretMap)
	}
}

func testdataPath(name string) string {
	return filepath.Join("..", "testdata", name)
}

func TestSaveAsYAML(t *testing.T) {
	secrets := Secrets{
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

	if err := secrets.SaveAsYAML(filename); err != nil {
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

func TestLoadFromYAML_valid(t *testing.T) {
	testcases := []struct {
		filename string
		kmsKey   string
	}{
		{
			filename: "test_valid.yaml",
			kmsKey:   "valec-qa",
		},
		{
			filename: "test_no_kmskey.yaml",
			kmsKey:   "",
		},
	}

	expectedSecrets := Secrets{
		&Secret{
			Key:   "FOO",
			Value: "bar",
		},
		&Secret{
			Key:   "BAZ",
			Value: "1",
		},
		&Secret{
			Key:   "QUX",
			Value: "true",
		},
	}

	for _, tc := range testcases {
		filepath := testdataPath(tc.filename)
		kmsKey, secrets, err := LoadFromYAML(filepath)
		if err != nil {
			t.Errorf("Error should not be raised. error: %s, filename: %s", err, tc.filename)
		}

		if kmsKey != tc.kmsKey {
			t.Errorf("kmsKey does not match. expected: %s, actual: %s, filename: %s", tc.kmsKey, kmsKey, tc.filename)
		}

		if !reflect.DeepEqual(secrets, expectedSecrets) {
			t.Errorf("Secrets does not match. expected: %v, actual: %v, filename: %s", expectedSecrets, secrets, tc.filename)
		}
	}
}

func TestLoadFromYAML_invalid(t *testing.T) {
	testcases := []struct {
		filename  string
		errPrefix string
	}{
		{
			filename:  "test_invalid.yaml",
			errPrefix: fmt.Sprintf("Failed to parse secret file as YAML. filename=%s", testdataPath("test_invalid.yaml")),
		},
		{
			filename:  "test_notexist.yaml",
			errPrefix: fmt.Sprintf("Failed to read secret file. filename=%s", testdataPath("test_notexist.yaml")),
		},
	}

	for _, tc := range testcases {
		filepath := testdataPath(tc.filename)
		_, _, err := LoadFromYAML(filepath)
		if err == nil {
			t.Errorf("Error should be raised. filename: %s", tc.filename)
		}

		if !strings.HasPrefix(err.Error(), tc.errPrefix) {
			t.Fatalf("Error message prefix does not match. expected prefix: %q, actual message: %q, filename: %s", tc.errPrefix, err.Error(), tc.filename)
		}
	}
}

func TestMapToList(t *testing.T) {
	secretMap := map[string]string{
		"FOO":  "bar",
		"BAZ":  "1",
		"HOGE": "fuga",
	}
	expected := Secrets{
		&Secret{
			Key:   "BAZ",
			Value: "1",
		},
		&Secret{
			Key:   "FOO",
			Value: "bar",
		},
		&Secret{
			Key:   "HOGE",
			Value: "fuga",
		},
	}

	secrets := MapToList(secretMap)

	if !reflect.DeepEqual(secrets, expected) {
		t.Errorf("Secret list does not match. expected: %q, actual:%q", expected, secrets)
	}
}
