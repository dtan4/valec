package util

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/Songmu/prompter"
	"github.com/pkg/errors"
)

var (
	separatorRegExp = regexp.MustCompile(`^#+\s*[-=]{3,}`)
	yamlExtRegexp   = regexp.MustCompile(`\.[yY][aA]?[mM][lL]$`)
)

// CompareStrings compares two string slices
func CompareStrings(src, dst []string) ([]string, []string) {
	added, deleted := []string{}, []string{}

	ss := map[string]int{}

	for _, s := range src {
		ss[s] = 1
	}

	for _, s := range dst {
		ss[s] += 2
	}

	for k, v := range ss {
		switch v {
		case 1:
			deleted = append(deleted, k)
		case 2:
			added = append(added, k)
		}
	}

	sort.Strings(added)
	sort.Strings(deleted)

	return added, deleted
}

// IsExist returns whether the given file / directory exists or not
func IsExist(name string) bool {
	_, err := os.Stat(name)
	return err == nil || os.IsExist(err)
}

// IsSecretFile returns whether the given file is secret file or not
func IsSecretFile(filename string) bool {
	base := filepath.Base(filename)

	return !strings.HasPrefix(base, ".") && yamlExtRegexp.MatchString(filepath.Ext(base))
}

// NamespaceFromPath returns namespace from the given path
func NamespaceFromPath(path, basedir string) (string, error) {
	var namespace string

	fullpath, err := filepath.Abs(path)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to get full path of file. filename=%s", path)
	}

	fulldir, err := filepath.Abs(basedir)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to get full path of directory. dirname=%s", basedir)
	}

	namespace = strings.Replace(fullpath, fulldir, "", 1)
	namespace = filepath.ToSlash(namespace)
	namespace = yamlExtRegexp.ReplaceAllString(namespace, "")

	if strings.HasPrefix(namespace, "/") {
		namespace = namespace[1:len(namespace)]
	}

	return namespace, nil
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

// ScanLines reads text stream and return
func ScanLines(r io.Reader) []string {
	lines := []string{}
	sc := bufio.NewScanner(r)

	for sc.Scan() {
		lines = append(lines, sc.Text())
	}

	return lines
}

// ScanNoecho reads password without printing password text in console
func ScanNoecho(key string) string {
	return prompter.Password(key)
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
