package util

import (
	"bufio"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

var (
	separatorRegExp = regexp.MustCompile(`^#+\s*[-=]{3,}`)
)

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
