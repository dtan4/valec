package util

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestSeparatorRegexp(t *testing.T) {
	positiveCases := []string{
		"#--------",
		"#========",
		"# ------",
		"# =====",
		"### =====",
		"# ===== production",
		"# ===== production =====",
	}

	for _, s := range positiveCases {
		if !separatorRegExp.Match([]byte(s)) {
			t.Errorf("String %q should be matched to regexp.", s)
		}
	}

	negativeCases := []string{
		"#FOO=bar",
		"FOO=#===bar",
	}

	for _, s := range negativeCases {
		if separatorRegExp.Match([]byte(s)) {
			t.Errorf("String %q should not be matched to regexp.", s)
		}
	}
}

func TestCompareStrings(t *testing.T) {
	src := []string{
		"foo",
		"bar",
		"baz",
	}
	dst := []string{
		"bar",
		"qux",
	}
	expectedAdded := []string{
		"qux",
	}
	expectedDeleted := []string{
		"baz",
		"foo",
	}

	added, deleted := CompareStrings(src, dst)

	if !reflect.DeepEqual(added, expectedAdded) {
		t.Errorf("Added strings does not match. expected: %q, actual: %q", expectedAdded, added)
	}

	if !reflect.DeepEqual(deleted, expectedDeleted) {
		t.Errorf("Deleted strings does not match. expected: %q, actual: %q", expectedDeleted, deleted)
	}
}

func TestIsExist(t *testing.T) {
	dir, err := ioutil.TempDir("", "test-save-as-dotenv")
	if err != nil {
		t.Fatalf("Failed to create tempdir. dir: %s", dir)
	}
	defer os.RemoveAll(dir)

	name := filepath.Join(dir, "sample")
	if err := ioutil.WriteFile(name, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create tempfile. name: %s", name)
	}

	testcases := []struct {
		name     string
		expected bool
	}{
		{
			name:     "",
			expected: true,
		},
		{
			name:     "sample",
			expected: true,
		},
		{
			name:     "nonexist",
			expected: false,
		},
	}

	for _, tc := range testcases {
		if IsExist(filepath.Join(dir, tc.name)) != tc.expected {
			t.Errorf("IsExist result is wrong. filename: %s, expected: %t", tc.name, tc.expected)
		}
	}
}

func TestIsSecretFile(t *testing.T) {
	testcases := []struct {
		filename string
		expected bool
	}{
		{
			filename: filepath.Join("secrets", "foo.yml"),
			expected: true,
		},
		{
			filename: filepath.Join("secrets", "foo", "bar.yaml"),
			expected: true,
		},
		{
			filename: filepath.Join("secrets", "foo", "bar"),
			expected: false,
		},
		{
			filename: filepath.Join("secrets", "foo", ".env"),
			expected: false,
		},
	}

	for _, tc := range testcases {
		actual := IsSecretFile(tc.filename)
		if actual != tc.expected {
			t.Errorf("IsSecretFile result is wrong. filename: %s, expected: %t, actual: %t", tc.filename, tc.expected, actual)
		}
	}
}

func TestNamespaceFromPath(t *testing.T) {
	testcases := []struct {
		path     string
		basedir  string
		expected string
	}{
		{
			path:     filepath.Join("secrets", "foo.yml"),
			basedir:  "secrets",
			expected: "foo",
		},
		{
			path:     filepath.Join("secrets", "foo.yml"),
			basedir:  ".",
			expected: "secrets/foo",
		},
		{
			path:     filepath.Join("secrets", "foo", "bar.yml"),
			basedir:  "secrets",
			expected: "foo/bar",
		},
		{
			path:     filepath.Join("secrets", "foo", "bar.yml"),
			basedir:  filepath.Join("secrets", "foo"),
			expected: "bar",
		},
		{
			path:     filepath.Join("secrets", "foo", "bar", "baz.yaml"),
			basedir:  "secrets",
			expected: "foo/bar/baz",
		},
	}

	for _, tc := range testcases {
		actual, err := NamespaceFromPath(tc.path, tc.basedir)
		if err != nil {
			t.Errorf("Error should not be raised. error: %s", err)
		}

		if actual != tc.expected {
			t.Errorf("Namespace does not match. expected: %q, actual: %q", tc.expected, actual)
		}
	}
}

func TestListYAMLFiles(t *testing.T) {
	dirname := filepath.Join("..", "testdata", "foo")
	expected := []string{
		filepath.Join("..", "testdata", "foo", "bar", "test_valid.yaml"),
		filepath.Join("..", "testdata", "foo", "test.yml"),
	}

	actual, err := ListYAMLFiles(dirname)
	if err != nil {
		t.Errorf("Error should not raised. err: %s", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("YAML file list is wrong. expected: %q, actual: %q", expected, actual)
	}
}

func TestScanLines(t *testing.T) {
	body := `FOO=bar
BAZ=1
HOGE=fuga`
	r := bytes.NewBufferString(body)
	expected := []string{
		"FOO=bar",
		"BAZ=1",
		"HOGE=fuga",
	}

	actual := ScanLines(r)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Lines are different. expected: %q, actual: %q", expected, actual)
	}
}

func TestWriteFile(t *testing.T) {
	body := []byte(`FOO=bar
BAZ=1
HOGE=fuga
`)
	dir, err := ioutil.TempDir("", "test-save-as-dotenv")
	if err != nil {
		t.Fatalf("Failed to create tempdir. dir: %s", dir)
	}
	defer os.RemoveAll(dir)

	filename := filepath.Join(dir, "secret.yaml")

	if err := WriteFile(filename, body); err != nil {
		t.Errorf("Error should not be raised. err: %s", err)
	}

	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			t.Errorf("File is not created. filename: %s", filename)
		} else {
			t.Errorf("Saved file has something wrong. err: %s", err)
		}
	}

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to open file. filename: %s", filename)
	}

	actual := string(b)

	if actual != string(body) {
		t.Errorf("File body does not match. expected: %q, actual: %q", string(body), actual)
	}
}

func TestWriteFileWithoutSection(t *testing.T) {
	dir, err := ioutil.TempDir("", "test-save-as-dotenv-preserve")
	if err != nil {
		t.Fatalf("Failed to create tempdir. dir: %s", dir)
	}
	defer os.RemoveAll(dir)

	srcName := filepath.Join("..", "testdata", "test.env")
	dstName := filepath.Join(dir, ".env")
	if err := copyFile(srcName, dstName); err != nil {
		t.Fatalf("Failed to copy file. src: %s, dst: %s, err: %s", srcName, dstName, err)
	}

	body := []byte(`FOO=bar
BAZ=1
HOGE=fuga
`)

	if err := WriteFileWithoutSection(dstName, body); err != nil {
		t.Errorf("Error should not be raised. err: %s", err)
	}

	if _, err := os.Stat(dstName); err != nil {
		if os.IsNotExist(err) {
			t.Errorf("File is not created. filename: %s", dstName)
		} else {
			t.Errorf("Saved file has something wrong. err: %s", err)
		}
	}

	b, err := ioutil.ReadFile(dstName)
	if err != nil {
		t.Fatalf("Failed to open file. filename: %s", dstName)
	}

	expected := `FOO=bar
BAZ=1
HOGE=fuga

#-------------

FOO=production-foo
`
	actual := string(b)

	if actual != expected {
		t.Errorf("File body does not match. expected: %q, actual: %q", expected, actual)
	}
}

func copyFile(srcName, dstName string) error {
	src, err := os.Open(srcName)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(dstName)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	return nil
}
