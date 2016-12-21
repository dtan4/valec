package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/dtan4/valec/aws"
	"github.com/dtan4/valec/secret"
	"github.com/pkg/errors"
)

func dumpAll(secrets secret.Secrets, quote bool) ([]string, error) {
	dotenv := []string{}

	for _, secret := range secrets {
		plainValue, err := aws.KMS.DecryptBase64(secret.Key, secret.Value)
		if err != nil {
			return []string{}, errors.Wrap(err, "Failed to decrypt value.")
		}

		if quote {
			dotenv = append(dotenv, fmt.Sprintf("%s=%q", secret.Key, plainValue))
		} else {
			dotenv = append(dotenv, fmt.Sprintf("%s=%s", secret.Key, plainValue))
		}
	}

	return dotenv, nil
}

func dumpWithTemplate(secrets secret.Secrets, quote bool) ([]string, error) {
	fp, err := os.Open(dotenvTemplate)
	if err != nil {
		return []string{}, errors.Wrapf(err, "Failed to open dotenv template. filename=%s", dotenvTemplate)
	}
	defer fp.Close()

	secretMap := secrets.ListToMap()
	sc := bufio.NewScanner(fp)
	dotenv := []string{}

	for sc.Scan() {
		line := sc.Text()

		if strings.HasPrefix(line, "#") {
			dotenv = append(dotenv, line)
			continue
		}

		ss := strings.SplitN(line, "=", 2)
		if len(ss) != 2 {
			dotenv = append(dotenv, line)
			continue
		}

		key, value := ss[0], ss[1]

		if override || value == "" {
			v, ok := secretMap[key]
			if ok {
				plainValue, err := aws.KMS.DecryptBase64(key, v)
				if err != nil {
					return []string{}, errors.Wrap(err, "Failed to decrypt value.")
				}

				value = plainValue
			}
		}

		if quote {
			dotenv = append(dotenv, fmt.Sprintf("%s=%q", key, value))
		} else {
			dotenv = append(dotenv, fmt.Sprintf("%s=%s", key, value))
		}
	}

	return dotenv, nil
}

func scanFromStdin(r io.Reader) []string {
	lines := []string{}
	sc := bufio.NewScanner(r)

	for sc.Scan() {
		lines = append(lines, sc.Text())
	}

	return lines
}
