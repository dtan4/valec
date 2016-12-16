package util

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

var (
	separatorRegExp = regexp.MustCompile(`^#+\s*[-=]{3,}`)
	yamlExtRegexp   = regexp.MustCompile(`\.[yY][aA]?[mM][lL]$`)
)

// IsSecretFile returns whether the given file is secret file or not
func IsSecretFile(filename string) bool {
	base := filepath.Base(filename)

	return !strings.HasPrefix(base, ".") && yamlExtRegexp.MatchString(filepath.Ext(base))
}

// NamespaceFromPath returns namespace from the given path
func NamespaceFromPath(path, basedir string) string {
	var namespace string

	namespace = strings.Replace(path, basedir, "", 1)
	namespace = filepath.ToSlash(namespace)
	namespace = yamlExtRegexp.ReplaceAllString(namespace, "")

	if strings.HasPrefix(namespace, "/") {
		namespace = namespace[1:len(namespace)]
	}

	return namespace
}

// ListYAMLFiles parses and executes function recursively
func ListYAMLFiles(dirname string) ([]string, error) {
	files := []string{}

	fs, err := ioutil.ReadDir(dirname)
	if err != nil {
		return []string{}, errors.Wrapf(err, "Failed to open directory. dirname=%s")
	}

	for _, file := range fs {
		if file.IsDir() {
			childDir := filepath.Join(dirname, file.Name())

			childFiles, err := ListYAMLFiles(childDir)
			if err != nil {
				return []string{}, errors.Wrapf(err, "failed to parse directory. dirname=%s", childDir)
			}

			files = append(files, childFiles...)

			continue
		}

		filename := filepath.Join(dirname, file.Name())

		if !IsSecretFile(filename) {
			continue
		}

		files = append(files, filename)
	}

	return files, nil
}

// WriteFile writes body to file
func WriteFile(filename string, body []byte) error {
	if err := ioutil.WriteFile(filename, body, 0644); err != nil {
		return errors.Wrapf(err, "Failed to save file. filename=%s", filename)
	}

	return nil
}

// WriteFileWithoutSection writes body to file keeping preserved section
func WriteFileWithoutSection(filename string, body []byte) error {
	fp, err := os.Open(filename)
	if err != nil {
		return errors.Wrapf(err, "Failed to open file. filename=%s", filename)
	}

	sc := bufio.NewScanner(fp)
	preserve := false
	preserveLines := []string{""}

	for sc.Scan() {
		line := sc.Text()

		if !preserve && separatorRegExp.Match([]byte(line)) {
			preserve = true
		}

		if preserve {
			preserveLines = append(preserveLines, line)
		}
	}

	body = append(body, []byte(strings.Join(preserveLines, "\n")+"\n")...)

	if err := ioutil.WriteFile(filename, body, 0644); err != nil {
		return errors.Wrapf(err, "Failed to save file. filename=%s", filename)
	}

	return nil
}
