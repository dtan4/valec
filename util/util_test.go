package util

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
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
